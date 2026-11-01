package dbversion

import (
	"strings"
	"testing"

	"lumini-hub/common/testutil"
)

func TestCurrent_SemTabelaDevolveZero(t *testing.T) {
	db := testutil.OpenTestDB(t)
	got, err := Current(db)
	if err != nil || got != 0 {
		t.Fatalf("Current = %d, %v; esperado 0 sem erro", got, err)
	}
}

func TestRecord_PrimeiraVezIdempotenteENovaPublicacao(t *testing.T) {
	db := testutil.OpenTestDB(t)

	inserted, err := Record(db)
	if err != nil || !inserted {
		t.Fatalf("1º Record = %v, %v; esperado inserir", inserted, err)
	}
	inserted, err = Record(db)
	if err != nil || inserted {
		t.Fatalf("2º Record = %v, %v; reiniciar não pode duplicar", inserted, err)
	}
	var n int64
	db.Model(&SchemaVersion{}).Count(&n)
	if n != 1 {
		t.Fatalf("%d linhas, esperado 1", n)
	}
	if got, _ := Current(db); got != Required {
		t.Fatalf("Current = %d, esperado %d", got, Required)
	}

	// simula uma publicação antiga já registrada: outra versão de sistema, mesmo banco
	if err := db.Model(&SchemaVersion{}).Where("1 = 1").Update("app_version", "2026.2.1.0").Error; err != nil {
		t.Fatal(err)
	}
	inserted, err = Record(db)
	if err != nil || !inserted {
		t.Fatalf("nova publicação = %v, %v; esperado inserir linha nova", inserted, err)
	}
	db.Model(&SchemaVersion{}).Count(&n)
	if n != 2 {
		t.Fatalf("%d linhas, esperado 2 (histórico)", n)
	}
}

func TestCheck(t *testing.T) {
	db := testutil.OpenTestDB(t)

	err := Check(db)
	if err == nil || !strings.Contains(err.Error(), "não tem versão registrada") {
		t.Fatalf("banco sem tabela: erro = %v", err)
	}

	if err := db.AutoMigrate(&SchemaVersion{}); err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&SchemaVersion{DBVersion: Required - 1, AppVersion: "x"}).Error; err != nil {
		t.Fatal(err)
	}
	err = Check(db)
	if err == nil || !strings.Contains(err.Error(), "está na versão") {
		t.Fatalf("banco atrás: erro = %v", err)
	}

	if err := db.Create(&SchemaVersion{DBVersion: Required, AppVersion: AppVersion}).Error; err != nil {
		t.Fatal(err)
	}
	err = Check(db)
	if err == nil || !strings.Contains(err.Error(), "migrações SQL pendentes") {
		t.Fatalf("versão certa mas com migração pendente: erro = %v", err)
	}
	if _, err := ApplyMigrations(db); err != nil {
		t.Fatalf("aplicando as migrações: %v", err)
	}
	if err := Check(db); err != nil {
		t.Fatalf("banco na versão exigida: %v", err)
	}

	if err := db.Create(&SchemaVersion{DBVersion: Required + 1, AppVersion: "futuro"}).Error; err != nil {
		t.Fatal(err)
	}
	if err := Check(db); err != nil {
		t.Fatalf("banco à frente do código não é erro: %v", err)
	}
}
