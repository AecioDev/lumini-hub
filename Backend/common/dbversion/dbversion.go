// Package dbversion registra, dentro do próprio banco, em que versão ele está, e
// confere no boot se o banco é novo o bastante para o código que está subindo.
//
// A versão do banco (2º número da versão do sistema, ano.versao_do_banco.sequencia.patch)
// sobe UMA unidade por publicação que traz mudanças de banco (tabela, coluna, índice,
// seed), não por mudança individual. Patches e publicações sem mudança de banco mantêm
// o número. Ver BACKLOG.md (CFG-020).
package dbversion

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// A cada publicação, o Claude/desenvolvedor atualiza estas duas constantes à mão:
//   - AppVersion: sempre, é a versão completa que está sendo publicada.
//   - Required: só quando a publicação traz mudanças de banco (sobe uma unidade).
const (
	// Required é a versão do banco exigida por este código.
	Required = 3
	// AppVersion é a versão do sistema, no formato ano.versao_do_banco.sequencia_versao.sequencia_patch.
	AppVersion = "2026.3.4.0"
)

// SchemaVersion é uma linha do histórico: uma por publicação que chegou a este banco.
type SchemaVersion struct {
	ID         uint      `gorm:"primaryKey"`
	DBVersion  int       `gorm:"not null;index"`
	AppVersion string    `gorm:"size:30;not null"`
	AppliedAt  time.Time `gorm:"not null"`
}

// TableName fixa o nome da tabela.
func (SchemaVersion) TableName() string { return "schema_versions" }

// Current devolve a maior versão de banco registrada. Tabela inexistente ou vazia
// significa banco sem versão registrada: devolve 0, sem erro.
func Current(db *gorm.DB) (int, error) {
	if !db.Migrator().HasTable(&SchemaVersion{}) {
		return 0, nil
	}
	var latest SchemaVersion
	err := db.Order("db_version DESC, id DESC").Take(&latest).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("lendo schema_versions: %w", err)
	}
	return latest.DBVersion, nil
}

// Record garante a tabela e registra a publicação atual (Required + AppVersion) se ela
// ainda não for a mais recente do histórico. Idempotente: reiniciar o serviço não
// duplica a linha. Devolve true quando inseriu uma linha nova. Só o api.auth chama,
// depois do AutoMigrate e do seed terminarem com sucesso.
func Record(db *gorm.DB) (bool, error) {
	if err := db.AutoMigrate(&SchemaVersion{}); err != nil {
		return false, fmt.Errorf("criando schema_versions: %w", err)
	}

	var latest SchemaVersion
	err := db.Order("id DESC").Take(&latest).Error
	if err == nil && latest.DBVersion == Required && latest.AppVersion == AppVersion {
		return false, nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, fmt.Errorf("lendo schema_versions: %w", err)
	}

	row := SchemaVersion{DBVersion: Required, AppVersion: AppVersion, AppliedAt: time.Now()}
	if err := db.Create(&row).Error; err != nil {
		return false, fmt.Errorf("registrando a versão do banco: %w", err)
	}
	return true, nil
}

// Check recusa o boot quando o banco está atrás do que o código exige (ou sem versão
// registrada). Usado pelos serviços com DB_AUTO_MIGRATE=false, onde o operador cuida
// do schema. Banco à frente do código não é erro (rollback de código é possível).
func Check(db *gorm.DB) error {
	current, err := Current(db)
	if err != nil {
		return err
	}
	if current >= Required {
		return nil
	}
	if current == 0 {
		return fmt.Errorf("o banco não tem versão registrada e este sistema (%s) exige a versão %d: "+
			"suba uma vez o api.auth com DB_AUTO_MIGRATE=true para atualizar e registrar o banco, ou aplique a atualização por outro meio",
			AppVersion, Required)
	}
	return fmt.Errorf("o banco está na versão %d e este sistema (%s) exige a versão %d: "+
		"suba uma vez o api.auth com DB_AUTO_MIGRATE=true para atualizar o banco, ou aplique a atualização por outro meio",
		current, AppVersion, Required)
}
