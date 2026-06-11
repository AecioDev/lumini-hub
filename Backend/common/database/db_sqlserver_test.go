package database

import (
	"lumini-hub/common/config"
	"testing"
)

func TestInitSQLServerDB(t *testing.T) {
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Erro ao carregar configurações: %v", err)
	}

	// Se o host ou banco de dados não estiver configurado de forma que indique uma conexão real, pulamos
	if cfg.SQLServer.Password == "" {
		t.Skip("Pulando teste de conexão SQL Server: senha não configurada")
	}

	db, err := InitSQLServerDB(cfg)
	if err != nil {
		t.Fatalf("Erro ao conectar ao SQL Server: %v", err)
	}

	var result int
	err = db.Raw("SELECT 1").Scan(&result).Error
	if err != nil {
		t.Fatalf("Erro ao executar query de teste (SELECT 1): %v", err)
	}

	if result != 1 {
		t.Errorf("Esperado resultado 1, obtido %d", result)
	}

	t.Log("Conexão e consulta de teste executadas com sucesso no SQL Server!")
}
