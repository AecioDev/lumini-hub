package main

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"

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
	if first.permissions == 0 || first.roles != 1 || first.links == 0 || first.menus == 0 {
		t.Errorf("1º boot incompleto: %+v", first)
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

func TestPrepareDatabase_FlagDesligadaEmBancoPopuladoNaoSemeiaNada(t *testing.T) {
	db := testutil.OpenTestDB(t)
	cfg := &config.Config{}
	cfg.Database.AutoMigrate = true
	if err := prepareDatabase(db, cfg); err != nil {
		t.Fatal(err)
	}
	if err := db.Exec("DELETE FROM menu_items").Error; err != nil {
		t.Fatal(err)
	}

	cfg.Database.AutoMigrate = false
	if err := prepareDatabase(db, cfg); err != nil {
		t.Fatal(err)
	}
	if got := takeBoot(t, db); got.menus != 0 {
		t.Errorf("a flag desligada semeou %d itens de menu", got.menus)
	}
}
