package dbversion

import (
	"embed"
	"fmt"
	"io/fs"
	"sort"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Migrações SQL versionadas: cobrem o que o AutoMigrate não faz (remover constraint ou índice
// antigo, corrigir dados existentes). Cada arquivo de migrations/ roda UMA vez por banco, em ordem
// de nome (NNNN_descricao.sql), dentro de uma transação, e fica registrado em schema_migrations.
// O SQL deve ser idempotente e tolerar tabelas ainda inexistentes (o api.auth roda antes do
// api.core num banco novo). Não ponha BEGIN/COMMIT nos arquivos. Uma migração aplicada nunca
// deve ser editada: crie outra.
//
//go:embed migrations/*.sql
var embeddedMigrations embed.FS

// SchemaMigration registra uma migração já aplicada.
type SchemaMigration struct {
	Name      string    `gorm:"primaryKey;size:150"`
	AppliedAt time.Time `gorm:"not null"`
}

// TableName fixa o nome da tabela.
func (SchemaMigration) TableName() string { return "schema_migrations" }

// lockKey é a chave do advisory lock que impede duas instâncias de migrarem ao mesmo tempo.
const lockKey = 7424100050

func listMigrations(fsys fs.FS) ([]string, error) {
	entries, err := fs.Glob(fsys, "migrations/*.sql")
	if err != nil {
		return nil, err
	}
	sort.Strings(entries)
	return entries, nil
}

func migrationName(path string) string { return strings.TrimPrefix(path, "migrations/") }

// ApplyMigrations aplica as migrações embutidas que ainda não rodaram neste banco e devolve os
// nomes aplicadas. Só o api.auth chama (com DB_AUTO_MIGRATE=true), depois do AutoMigrate.
func ApplyMigrations(db *gorm.DB) ([]string, error) { return applyMigrations(db, embeddedMigrations) }

func applyMigrations(db *gorm.DB, fsys fs.FS) ([]string, error) {
	if err := db.AutoMigrate(&SchemaMigration{}); err != nil {
		return nil, fmt.Errorf("criando schema_migrations: %w", err)
	}
	files, err := listMigrations(fsys)
	if err != nil {
		return nil, err
	}
	var applied []string
	for _, path := range files {
		name := migrationName(path)
		ran, err := applyOne(db, fsys, path, name)
		if err != nil {
			return applied, fmt.Errorf("migração %s: %w", name, err)
		}
		if ran {
			applied = append(applied, name)
		}
	}
	return applied, nil
}

func applyOne(db *gorm.DB, fsys fs.FS, path, name string) (bool, error) {
	body, err := fs.ReadFile(fsys, path)
	if err != nil {
		return false, err
	}
	ran := false
	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SELECT pg_advisory_xact_lock(?)", lockKey).Error; err != nil {
			return err
		}
		var n int64
		if err := tx.Model(&SchemaMigration{}).Where("name = ?", name).Count(&n).Error; err != nil {
			return err
		}
		if n > 0 {
			return nil
		}
		if err := tx.Exec(string(body)).Error; err != nil {
			return err
		}
		ran = true
		return tx.Create(&SchemaMigration{Name: name, AppliedAt: time.Now()}).Error
	})
	return ran, err
}

// PendingMigrations lista as migrações embutidas que ainda não rodaram neste banco. Sem a tabela
// schema_migrations, todas estão pendentes.
func PendingMigrations(db *gorm.DB) ([]string, error) {
	return pendingMigrations(db, embeddedMigrations)
}

func pendingMigrations(db *gorm.DB, fsys fs.FS) ([]string, error) {
	files, err := listMigrations(fsys)
	if err != nil {
		return nil, err
	}
	done := map[string]bool{}
	if db.Migrator().HasTable(&SchemaMigration{}) {
		var rows []SchemaMigration
		if err := db.Find(&rows).Error; err != nil {
			return nil, fmt.Errorf("lendo schema_migrations: %w", err)
		}
		for _, r := range rows {
			done[r.Name] = true
		}
	}
	var pending []string
	for _, path := range files {
		if name := migrationName(path); !done[name] {
			pending = append(pending, name)
		}
	}
	return pending, nil
}
