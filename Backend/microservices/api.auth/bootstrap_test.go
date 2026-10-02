package main

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"

	"lumini-hub/api.auth/seed"
	"lumini-hub/common/config"
	"lumini-hub/common/testutil"

	"gorm.io/gorm"
)

// Testes de banco: precisam de TEST_DATABASE_URL (ver testutil.OpenTestDB) e são
// pulados sem ela.

type bootSnapshot struct{ permissions, roles, links, menus, users int64 }

func takeBoot(t *testing.T, db *gorm.DB) bootSnapshot {
	t.Helper()
	var s bootSnapshot
	for table, dest := range map[string]*int64{
		"permissions":      &s.permissions,
		"roles":            &s.roles,
		"role_permissions": &s.links,
		"menu_items":       &s.menus,
		"users":            &s.users,
	} {
		if err := db.Table(table).Count(dest).Error; err != nil {
			t.Fatalf("contando %s: %v", table, err)
		}
	}
	return s
}

func tableCount(t *testing.T, db *gorm.DB) int64 {
	t.Helper()
	var n int64
	err := db.Raw("SELECT count(*) FROM information_schema.tables WHERE table_schema = current_schema()").Scan(&n).Error
	if err != nil {
		t.Fatal(err)
	}
	return n
}

// capturaLog devolve o buffer que recebe o log padrão durante o teste.
func capturaLog(t *testing.T) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	log.SetOutput(&buf)
	t.Cleanup(func() { log.SetOutput(os.Stderr) })
	return &buf
}

func TestPrepareDatabase_FlagDesligadaNaoFazNada(t *testing.T) {
	db := testutil.OpenTestDB(t)
	cfg := &config.Config{}
	cfg.Database.AutoMigrate = false

	if err := prepareDatabase(db, cfg); err != nil {
		t.Fatalf("com a flag desligada não deveria dar erro: %v", err)
	}
	if n := tableCount(t, db); n != 0 {
		t.Errorf("a flag desligada criou %d tabelas num banco vazio", n)
	}
}

func TestPrepareDatabase_BancoVazioSobeSozinhoEEhIdempotente(t *testing.T) {
	db := testutil.OpenTestDB(t)
	cfg := &config.Config{}
	cfg.Database.AutoMigrate = true
	logs := capturaLog(t)

	if err := prepareDatabase(db, cfg); err != nil {
		t.Fatalf("1º boot: %v", err)
	}
	first := takeBoot(t, db)
	var grantable int64
	for _, entry := range seed.PermissionCatalog {
		if entry.Module != seed.DevelopModule {
			grantable++
		}
	}
	if first.permissions != int64(len(seed.PermissionCatalog)) || first.roles != 1 || first.links != grantable || first.menus == 0 {
		t.Errorf("1º boot: %+v, esperado %d permissões, 1 role e %d vínculos (catálogo fora do módulo %s), menu não vazio",
			first, len(seed.PermissionCatalog), grantable, seed.DevelopModule)
	}
	var developLinks int64
	if err := db.Table("role_permissions AS rp").
		Joins("JOIN permissions p ON p.id = rp.permission_id").
		Where("p.module = ?", seed.DevelopModule).Count(&developLinks).Error; err != nil {
		t.Fatal(err)
	}
	if developLinks != 0 {
		t.Errorf("%d permissões do módulo %s vinculadas ao ADMIN, esperado 0", developLinks, seed.DevelopModule)
	}
	if first.users != 0 {
		t.Errorf("o boot não deveria criar usuários (isso é o bootstrap do admin): %d", first.users)
	}
	if strings.Contains(logs.String(), "não encontrada") {
		t.Errorf("o seed do menu avisou permissão não encontrada: %s", logs.String())
	}

	if err := prepareDatabase(db, cfg); err != nil {
		t.Fatalf("2º boot: %v", err)
	}
	if second := takeBoot(t, db); second != first {
		t.Errorf("2º boot mudou o estado: antes %+v, depois %+v", first, second)
	}
}

// O ADMIN é dono do catálogo de perfis: o menu "Perfis e Permissões" tem que nascer ligado a
// roles.view (como no banco de dev), e não a admin.create_permissions, que é só do DEVELOP.
// Um banco semeado do zero com o vínculo errado esconde o menu do ADMIN.
func TestPrepareDatabase_MenuDeRolesNasceComRolesView(t *testing.T) {
	db := testutil.OpenTestDB(t)
	cfg := &config.Config{}
	cfg.Database.AutoMigrate = true
	if err := prepareDatabase(db, cfg); err != nil {
		t.Fatal(err)
	}

	var code string
	err := db.Raw(`SELECT p.permission FROM menu_items m JOIN permissions p ON p.id = m.permission_id
		WHERE m.href = '/settings/roles' AND m.deleted_at IS NULL`).Scan(&code).Error
	if err != nil {
		t.Fatal(err)
	}
	if code != "roles.view" {
		t.Errorf("o menu /settings/roles exige %q, esperado roles.view", code)
	}
}

func TestPrepareDatabase_CriaAdminInicialSoUmaVezESoComSenha(t *testing.T) {
	const password = "SenhaForte-2026"

	t.Run("com senha: nasce no 1º boot e o 2º não duplica", func(t *testing.T) {
		db := testutil.OpenTestDB(t)
		cfg := &config.Config{}
		cfg.Database.AutoMigrate = true
		cfg.Bootstrap = config.BootstrapAdminConfig{Username: "admin", Password: password}
		logs := capturaLog(t)

		for boot := 1; boot <= 2; boot++ {
			if err := prepareDatabase(db, cfg); err != nil {
				t.Fatalf("%dº boot: %v", boot, err)
			}
			if got := takeBoot(t, db); got.users != 1 {
				t.Fatalf("%dº boot: %d usuários, esperado 1", boot, got.users)
			}
		}
		var roleName string
		err := db.Raw("SELECT r.name FROM users u JOIN roles r ON r.id = u.role_id WHERE u.username = 'admin'").Scan(&roleName).Error
		if err != nil || roleName != seed.AdminRoleName {
			t.Errorf("perfil do admin = %q (err=%v), esperado %s", roleName, err, seed.AdminRoleName)
		}
		if strings.Contains(logs.String(), password) {
			t.Errorf("a senha apareceu no log do boot: %s", logs.String())
		}
	})

	t.Run("sem senha: o boot funciona e não cria usuário", func(t *testing.T) {
		db := testutil.OpenTestDB(t)
		cfg := &config.Config{}
		cfg.Database.AutoMigrate = true
		if err := prepareDatabase(db, cfg); err != nil {
			t.Fatal(err)
		}
		if got := takeBoot(t, db); got.users != 0 {
			t.Errorf("criou %d usuários sem BOOTSTRAP_ADMIN_PASSWORD", got.users)
		}
	})

	t.Run("flag desligada: não cria nem o usuário", func(t *testing.T) {
		db := testutil.OpenTestDB(t)
		cfg := &config.Config{}
		cfg.Database.AutoMigrate = false
		cfg.Bootstrap = config.BootstrapAdminConfig{Username: "admin", Password: password}
		if err := prepareDatabase(db, cfg); err != nil {
			t.Fatal(err)
		}
		if n := tableCount(t, db); n != 0 {
			t.Errorf("a flag desligada criou %d tabelas", n)
		}
	})
}

func TestPrepareDatabase_FlagDesligadaEmBancoPopuladoNaoMexeEmNada(t *testing.T) {
	db := testutil.OpenTestDB(t)
	cfg := &config.Config{}
	cfg.Database.AutoMigrate = true
	if err := prepareDatabase(db, cfg); err != nil {
		t.Fatal(err)
	}

	// desfaz um pouco de cada coisa que o boot semeia: se a flag desligada ainda
	// rodasse qualquer parte do seed, algum desses voltaria. A permissão apagada é do
	// módulo Develop, que nunca é vinculada ao ADMIN (apagar uma vinculada violaria a FK).
	for _, stmt := range []string{
		"DELETE FROM menu_items",
		"DELETE FROM role_permissions WHERE permission_id = (SELECT min(id) FROM permissions)",
		"DELETE FROM permissions WHERE id = (SELECT min(id) FROM permissions WHERE module = '" + seed.DevelopModule + "')",
	} {
		if err := db.Exec(stmt).Error; err != nil {
			t.Fatalf("%s: %v", stmt, err)
		}
	}
	before := takeBoot(t, db)

	cfg.Database.AutoMigrate = false
	if err := prepareDatabase(db, cfg); err != nil {
		t.Fatal(err)
	}
	if after := takeBoot(t, db); after != before {
		t.Errorf("a flag desligada alterou o banco: antes %+v, depois %+v", before, after)
	}
}
