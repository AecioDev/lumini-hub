package main

import (
	"fmt"

	"lumini-hub/api.auth/internal/domain"
	"lumini-hub/api.auth/internal/seeder"
	"lumini-hub/api.auth/seed"
	"lumini-hub/common/config"

	"gorm.io/gorm"
)

// prepareDatabase deixa o banco pronto pro api.auth subir. Tudo roda sob
// DB_AUTO_MIGRATE (ligado por padrão pra um banco vazio subir sozinho): com a
// flag desligada — produção com dados reais — o serviço não migra NEM semeia
// nada ao subir, e o operador assume schema e dados (o catálogo de permissões
// precisa então ser aplicado por outro meio, ex.: rodar uma vez com a flag
// ligada). Na ordem:
//
//  1. AutoMigrate das tabelas. Ordem de dependência: Permission e Role antes de
//     User (FK role_id e as tabelas de junção role_permissions/user_permissions,
//     criadas pelo GORM a partir dos many2many) e de MenuItem (FK permission_id).
//  2. seed.SyncCatalog: catálogo de permissões + role ADMIN (permissão nova no
//     código entra na tabela e no ADMIN no boot seguinte).
//  3. seed.BootstrapAdmin: o primeiro usuário ADMIN, a partir de BOOTSTRAP_ADMIN_*, só
//     com a tabela users vazia e depois do SyncCatalog (a role ADMIN precisa existir).
//  4. seeder.SeedMenuItems: depois do catálogo, porque o menu referencia as
//     permissões pelo código; com elas ausentes o item nasceria sem permissão
//     própria, e isso nunca é corrigido depois (itens existentes são pulados).
func prepareDatabase(db *gorm.DB, cfg *config.Config) error {
	if !cfg.Database.AutoMigrate {
		return nil
	}

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
	if _, err := seed.BootstrapAdmin(db, seed.AdminBootstrap{
		Username: cfg.Bootstrap.Username,
		Password: cfg.Bootstrap.Password,
		Email:    cfg.Bootstrap.Email,
	}); err != nil {
		return fmt.Errorf("criando o admin inicial: %w", err)
	}
	if err := seeder.SeedMenuItems(db); err != nil {
		return fmt.Errorf("semeando itens de menu: %w", err)
	}
	return nil
}
