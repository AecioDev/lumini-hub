package service

import (
	"testing"

	"lumini-hub/common/config"
	"lumini-hub/common/dbversion"
	"lumini-hub/common/testutil"
)

func TestAuthService_GetVersionInfo(t *testing.T) {
	db := testutil.OpenTestDB(t)
	cfg := &config.Config{}
	cfg.App.Env = "staging"
	svc := NewAuthService(db, cfg)

	// banco sem versão registrada: a do sistema vem do código, a do banco é 0
	info, err := svc.GetVersionInfo()
	if err != nil {
		t.Fatal(err)
	}
	if info.SystemVersion != dbversion.AppVersion || info.RequiredDatabaseVersion != dbversion.Required {
		t.Errorf("versões do código = %q / %d", info.SystemVersion, info.RequiredDatabaseVersion)
	}
	if info.DatabaseVersion != 0 {
		t.Errorf("banco sem registro devolveu versão %d, esperado 0", info.DatabaseVersion)
	}
	if info.Environment != "staging" {
		t.Errorf("ambiente = %q", info.Environment)
	}

	// depois de registrar, a versão do banco passa a ser a exigida
	if _, err := dbversion.Record(db); err != nil {
		t.Fatal(err)
	}
	info, err = svc.GetVersionInfo()
	if err != nil {
		t.Fatal(err)
	}
	if info.DatabaseVersion != dbversion.Required {
		t.Errorf("versão do banco = %d, esperado %d", info.DatabaseVersion, dbversion.Required)
	}
}
