package database

import (
	"lumini-hub/common/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB inicializa e retorna a conexão GORM com o banco de dados
func InitDB(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.Database.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Habilitar logs lógicos para desenvolvimento
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}
