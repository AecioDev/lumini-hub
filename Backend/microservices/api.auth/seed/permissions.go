package seed

import (
	"errors"
	"fmt"

	"lumini-hub/api.auth/internal/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// SyncPermissions cria no banco as permissões do PermissionCatalog que ainda
// não existem (chave natural: o código em `permission`) e devolve só as que
// acabou de criar. Idempotente: permissão que já existe não é tocada — nem
// atualizada, nem apagada — e uma permissão excluída (soft delete) também conta
// como existente, porque `permission` é unique e recriá-la violaria a constraint.
// As permissões criadas à mão pelo operador ficam intactas.
func SyncPermissions(db *gorm.DB) ([]domain.Permission, error) {
	var created []domain.Permission

	err := db.Transaction(func(tx *gorm.DB) error {
		var err error
		created, err = syncPermissions(tx)
		return err
	})
	if err != nil {
		return nil, err
	}
	return created, nil
}

func syncPermissions(tx *gorm.DB) ([]domain.Permission, error) {
	var created []domain.Permission

	for _, entry := range PermissionCatalog {
		var existing domain.Permission
		err := tx.Unscoped().Where("permission = ?", entry.Permission).First(&existing).Error
		if err == nil {
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("buscando permissão %q: %w", entry.Permission, err)
		}

		permission := domain.Permission{
			Permission:  entry.Permission,
			Description: entry.Description,
			Module:      entry.Module,
		}
		// DoNothing: se outra instância criou a mesma permissão entre a busca e o
		// insert (boot simultâneo), o conflito no unique não derruba a transação.
		result := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&permission)
		if result.Error != nil {
			return nil, fmt.Errorf("criando permissão %q: %w", entry.Permission, result.Error)
		}
		if result.RowsAffected == 0 {
			continue
		}
		created = append(created, permission)
	}
	return created, nil
}
