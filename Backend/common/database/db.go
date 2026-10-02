package database

import (
	"lumini-hub/common/config"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// InitDB inicializa e retorna a conexão GORM com o banco de dados.
// O nível do log depende de APP_ENV (ver newLogger).
func InitDB(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.Database.DSN()), &gorm.Config{
		Logger: newLogger(cfg.App.Env, nil),
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// InitSQLServerDB inicializa e retorna a conexão GORM com o SQL Server.
// O nível do log depende de APP_ENV (ver newLogger).
func InitSQLServerDB(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(sqlserver.Open(cfg.SQLServer.DSN()), &gorm.Config{
		Logger: newLogger(cfg.App.Env, nil),
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}
