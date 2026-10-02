package main

import (
	"fmt"

	"lumini-hub/api.auth/internal/domain"
	"lumini-hub/api.auth/internal/seeder"
	"lumini-hub/api.auth/seed"
	"lumini-hub/common/config"

	"gorm.io/gorm"
)

// prepareDatabase deixa o banco pronto pro api.auth subir, na ordem:
//
//  1. AutoMigrate das tabelas (DB_AUTO_MIGRATE, ligado por padrão pra um banco
//     vazio subir sozinho; em produção com dados reais, desligar). Ordem de
//     dependência: Permission e Role antes de User (FK role_id e as tabelas de
//     junção role_permissions/user_permissions, criadas pelo GORM a partir dos
//     many2many) e de MenuItem (FK permission_id).
//  2. seed.SyncCatalog: catálogo de permissões + role ADMIN. Roda sob a mesma
//     flag porque é parte do mecanismo de migração (permissão nova no código
//     entra na tabela e no ADMIN no boot seguinte). Vem antes do menu porque o
//     menu referencia as permissões pelo código.
//  3. seeder.SeedMenuItems: sempre roda, como antes desta flag existir.
func prepareDatabase(db *gorm.DB, cfg *config.Config) error {
	if cfg.Database.AutoMigrate {
		if err := db.AutoMigrate(
			&domain.Permission{},
			&domain.Role{},
			&domain.User{},
			&domain.MenuItem{},
		); err != nil {
			return fmt.Errorf("migrando tabelas do api.auth: %w", err)
		}
		if err := seed.SyncCatalog(db); err != nil {
			return fmt.Errorf("semeando permissões e a role %s: %w", seed.AdminRoleName, err)
		}
	}

	if err := seeder.SeedMenuItems(db); err != nil {
		return fmt.Errorf("semeando itens de menu: %w", err)
	}
	return nil
}
