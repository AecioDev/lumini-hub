package seed

import (
	"testing"

	"lumini-hub/api.auth/internal/domain"
	"lumini-hub/common/testutil"

	"gorm.io/gorm"
)

// Testes de banco: precisam de TEST_DATABASE_URL (ver testutil.OpenTestDB) e são
// pulados sem ela. Não usam t.Parallel porque alguns trocam o PermissionCatalog global.

func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db := testutil.OpenTestDB(t)
	if err := db.AutoMigrate(&domain.Permission{}, &domain.Role{}); err != nil {
		t.Fatalf("migrando tabelas: %v", err)
	}
	return db
}

// withCatalog troca o catálogo por uma cópia com entradas extras e o restaura no fim.
func withCatalog(t *testing.T, extra ...PermissionSeed) {
	t.Helper()
	original := PermissionCatalog
	t.Cleanup(func() { PermissionCatalog = original })
	PermissionCatalog = append(append([]PermissionSeed{}, original...), extra...)
}

type snapshot struct{ permissions, roles, links, developLinks int64 }

func take(t *testing.T, db *gorm.DB) snapshot {
	t.Helper()
	var s snapshot
	must := func(err error) {
		t.Helper()
		if err != nil {
			t.Fatal(err)
		}
	}
	must(db.Model(&domain.Permission{}).Count(&s.permissions).Error)
	must(db.Model(&domain.Role{}).Count(&s.roles).Error)
	must(db.Table("role_permissions").Count(&s.links).Error)
	must(db.Table("role_permissions AS rp").
		Joins("JOIN permissions p ON p.id = rp.permission_id").
		Where("p.module = ?", DevelopModule).Count(&s.developLinks).Error)
	return s
}

// adminGrantable é quantas permissões do catálogo atual o ADMIN deve receber
// (todas fora do módulo Develop).
func adminGrantable() int64 {
	var n int64
	for _, entry := range PermissionCatalog {
		if entry.Module != DevelopModule {
			n++
		}
	}
	return n
}

func mustSync(t *testing.T, db *gorm.DB) {
	t.Helper()
	if err := SyncCatalog(db); err != nil {
		t.Fatalf("SyncCatalog: %v", err)
	}
}

func mustPermission(t *testing.T, db *gorm.DB, code string) domain.Permission {
	t.Helper()
	var permission domain.Permission
	if err := db.Where("permission = ?", code).First(&permission).Error; err != nil {
		t.Fatalf("permissão %s não encontrada: %v", code, err)
	}
	return permission
}

func adminRole(t *testing.T, db *gorm.DB) domain.Role {
	t.Helper()
	var role domain.Role
	if err := db.Where("name = ?", AdminRoleName).First(&role).Error; err != nil {
		t.Fatalf("role %s não encontrada: %v", AdminRoleName, err)
	}
	return role
}

// O catálogo não precisa de banco: roda sempre.
func TestPermissionCatalog_Integridade(t *testing.T) {
	if len(PermissionCatalog) == 0 {
		t.Fatal("catálogo vazio")
	}
	seen := map[string]bool{}
	for _, entry := range PermissionCatalog {
		if entry.Permission == "" || entry.Description == "" || entry.Module == "" {
			t.Errorf("entrada com campo vazio: %+v", entry)
		}
		if len(entry.Permission) > 100 || len(entry.Module) > 50 {
			t.Errorf("entrada estoura o tamanho da coluna: %+v", entry)
		}
		if seen[entry.Permission] {
			t.Errorf("permissão duplicada no catálogo: %s", entry.Permission)
		}
		seen[entry.Permission] = true
	}
}

// Os demais testes derivam os números do próprio catálogo (pra não quebrar quando o
// EPIC CFG-10 incluir permissões). Esta âncora fixa um PISO: apagar entradas por engano
// ou tirar o módulo Develop do catálogo faz este teste falhar. Atualizar os pisos junto
// com o CFG-10 se o catálogo crescer de propósito.
func TestPermissionCatalog_Ancora(t *testing.T) {
	const minTotal, minDevelop, minAdmin = 59, 9, 50

	var develop int
	for _, entry := range PermissionCatalog {
		if entry.Module == DevelopModule {
			develop++
		}
	}
	if len(PermissionCatalog) < minTotal {
		t.Errorf("catálogo com %d permissões, piso %d", len(PermissionCatalog), minTotal)
	}
	if develop < minDevelop {
		t.Errorf("módulo %s com %d permissões, piso %d", DevelopModule, develop, minDevelop)
	}
	if adminGrantable() < minAdmin {
		t.Errorf("o ADMIN receberia %d permissões, piso %d", adminGrantable(), minAdmin)
	}
}

func TestSyncCatalog_PrimeiroBoot(t *testing.T) {
	db := newTestDB(t)
	mustSync(t, db)

	got := take(t, db)
	if got.permissions != int64(len(PermissionCatalog)) {
		t.Errorf("permissões = %d, esperado %d (catálogo todo)", got.permissions, len(PermissionCatalog))
	}
	if got.roles != 1 {
		t.Errorf("roles = %d, esperado só a %s", got.roles, AdminRoleName)
	}
	if got.links != adminGrantable() {
		t.Errorf("vínculos = %d, esperado %d (catálogo fora do módulo %s)", got.links, adminGrantable(), DevelopModule)
	}
	if got.developLinks != 0 {
		t.Errorf("%d permissões do módulo %s vinculadas ao ADMIN, esperado 0", got.developLinks, DevelopModule)
	}
	if admin := adminRole(t, db); admin.Description != AdminRoleDescription {
		t.Errorf("descrição do ADMIN = %q", admin.Description)
	}
}

func TestSyncCatalog_Idempotente(t *testing.T) {
	db := newTestDB(t)
	mustSync(t, db)
	first := take(t, db)

	mustSync(t, db)
	if second := take(t, db); second != first {
		t.Errorf("2ª execução mudou o estado: antes %+v, depois %+v", first, second)
	}
}

func TestSyncCatalog_PermissaoManualEEditadaPermanecem(t *testing.T) {
	db := newTestDB(t)
	mustSync(t, db)

	if err := db.Create(&domain.Permission{Permission: "custom.manual", Description: "feita a mão", Module: "Custom"}).Error; err != nil {
		t.Fatal(err)
	}
	edited := PermissionCatalog[0]
	if err := db.Model(&domain.Permission{}).Where("permission = ?", edited.Permission).
		Update("description", "EDITADA PELO OPERADOR").Error; err != nil {
		t.Fatal(err)
	}

	mustSync(t, db)

	var manual domain.Permission
	if err := db.Where("permission = ?", "custom.manual").First(&manual).Error; err != nil {
		t.Errorf("permissão criada à mão sumiu: %v", err)
	}
	after := mustPermission(t, db, edited.Permission)
	if after.Description != "EDITADA PELO OPERADOR" {
		t.Errorf("o seed sobrescreveu a descrição editada: %q", after.Description)
	}
}

func TestSyncCatalog_VinculoRemovidoNaoVoltaEPermissaoExcluidaNaoRenasce(t *testing.T) {
	db := newTestDB(t)
	mustSync(t, db)
	before := take(t, db)

	// duas permissões do catálogo que o ADMIN recebeu (fora do módulo Develop)
	var removed, deleted domain.Permission
	for _, entry := range PermissionCatalog {
		if entry.Module == DevelopModule {
			continue
		}
		if removed.ID == 0 {
			removed = mustPermission(t, db, entry.Permission)
			continue
		}
		deleted = mustPermission(t, db, entry.Permission)
		break
	}
	if removed.ID == 0 || deleted.ID == 0 {
		t.Fatal("não achei permissões de teste no catálogo")
	}

	admin := adminRole(t, db)
	if err := db.Model(&admin).Association("Permissions").Delete(&removed); err != nil {
		t.Fatal(err)
	}
	if err := db.Delete(&deleted).Error; err != nil { // exclusão lógica
		t.Fatal(err)
	}

	mustSync(t, db)

	after := take(t, db)
	if after.links != before.links-1 {
		t.Errorf("vínculo removido foi recolocado: %d vínculos, esperado %d", after.links, before.links-1)
	}
	var alive int64
	db.Model(&domain.Permission{}).Where("permission = ?", deleted.Permission).Count(&alive)
	if alive != 0 {
		t.Error("permissão excluída logicamente foi recriada")
	}
}

func TestSyncCatalog_PermissaoNovaVinculaSoSeNaoForDevelop(t *testing.T) {
	db := newTestDB(t)
	mustSync(t, db)
	before := take(t, db)

	withCatalog(t,
		PermissionSeed{Permission: "novo.view", Description: "nova comum", Module: "Novo"},
		PermissionSeed{Permission: "novo.tecnico", Description: "nova técnica", Module: DevelopModule},
	)
	mustSync(t, db)

	after := take(t, db)
	if after.permissions != before.permissions+2 {
		t.Errorf("permissões = %d, esperado %d", after.permissions, before.permissions+2)
	}
	if after.links != before.links+1 {
		t.Errorf("vínculos = %d, esperado %d (só a que não é do módulo %s)", after.links, before.links+1, DevelopModule)
	}
	if after.developLinks != 0 {
		t.Errorf("permissão do módulo %s foi vinculada ao ADMIN", DevelopModule)
	}
}

func TestSyncCatalog_AdminExcluidoNaoERecriadoEOutraRoleImpedeCriacao(t *testing.T) {
	t.Run("ADMIN excluído", func(t *testing.T) {
		db := newTestDB(t)
		mustSync(t, db)
		admin := adminRole(t, db)
		if err := db.Delete(&admin).Error; err != nil {
			t.Fatal(err)
		}
		before := take(t, db)

		withCatalog(t, PermissionSeed{Permission: "outro.view", Description: "outra", Module: "Outro"})
		mustSync(t, db)

		var active int64
		db.Model(&domain.Role{}).Where("name = ?", AdminRoleName).Count(&active)
		if active != 0 {
			t.Error("o ADMIN excluído foi recriado")
		}
		after := take(t, db)
		if after.permissions != before.permissions+1 {
			t.Errorf("a permissão nova não entrou: %d, esperado %d", after.permissions, before.permissions+1)
		}
		if after.links != before.links {
			t.Errorf("vínculos mudaram sem ADMIN: %d, esperado %d", after.links, before.links)
		}
	})

	t.Run("outra role já existe e não há ADMIN", func(t *testing.T) {
		db := newTestDB(t)
		if err := db.Create(&domain.Role{Name: "Vendas", Description: "x"}).Error; err != nil {
			t.Fatal(err)
		}
		mustSync(t, db)

		var admins int64
		db.Model(&domain.Role{}).Where("name = ?", AdminRoleName).Count(&admins)
		got := take(t, db)
		if admins != 0 || got.roles != 1 {
			t.Errorf("ADMIN=%d roles=%d: o seed não pode criar o ADMIN fora do 1º boot", admins, got.roles)
		}
		if got.permissions != int64(len(PermissionCatalog)) || got.links != 0 {
			t.Errorf("permissões=%d vínculos=%d, esperado %d e 0", got.permissions, got.links, len(PermissionCatalog))
		}
	})
}

func TestSyncCatalog_PrimeiroBootComPermissoesJaExistentes(t *testing.T) {
	db := newTestDB(t)
	if _, err := SyncPermissions(db); err != nil {
		t.Fatal(err)
	}
	mustSync(t, db)

	got := take(t, db)
	if got.roles != 1 || got.links != adminGrantable() {
		t.Errorf("roles=%d vínculos=%d, esperado 1 e %d", got.roles, got.links, adminGrantable())
	}
}
