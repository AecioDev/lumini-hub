package repository

import (
	"testing"

	"lumini-hub/common/config"
	"lumini-hub/common/database"
)

func TestSequenceRepository_Next(t *testing.T) {
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Erro ao carregar configurações: %v", err)
	}

	if cfg.SQLServer.Password == "" {
		t.Skip("Pulando teste de sequência SQL Server: senha não configurada")
	}

	db, err := database.InitSQLServerDB(cfg)
	if err != nil {
		t.Fatalf("Erro ao conectar ao SQL Server: %v", err)
	}

	repo := NewSequenceRepository(db)

	first, err := repo.Next("TEST_SEQ_API_INTEGRATIONS")
	if err != nil {
		t.Fatalf("Erro ao obter primeira sequência: %v", err)
	}

	second, err := repo.Next("TEST_SEQ_API_INTEGRATIONS")
	if err != nil {
		t.Fatalf("Erro ao obter segunda sequência: %v", err)
	}

	if second != first+1 {
		t.Errorf("Esperado que a sequência incrementasse em 1 (obtido %d após %d)", second, first)
	}
}
