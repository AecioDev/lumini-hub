package seed

import (
	"errors"
	"fmt"

	"lumini-hub/api.auth/internal/domain"
	"lumini-hub/common/utils"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	// AdminRoleName é o único perfil criado pelo seed (os demais — Vendas,
	// Gerente etc. — virão depois, por segmento, e o DEVELOP é criado à mão).
	AdminRoleName = "ADMIN"
	// AdminRoleDescription é a descrição do perfil ADMIN, igual à do banco de dev.
	AdminRoleDescription = "Administrador do sistema"
	// DevelopModule é o módulo das permissões de catálogo/administração técnica
	// (roles.*, permissions.*, admin.create_permissions). O ADMIN recebe toda
	// permissão cujo módulo NÃO é este — o que dá as 50 de hoje. É regra de
	// vínculo-modelo, distinta de utils.IsDeveloperOnlyPermission (bypass em
	// tempo de execução): o ADMIN passa pelo bypass de roles.* e permissions.view.
	DevelopModule = utils.DeveloperModule
)

// SyncCatalog sincroniza o catálogo de permissões e o perfil ADMIN. Idempotente:
//
//   - Primeiro boot (nenhuma role no banco): cria a role ADMIN e vincula a ela
//     todas as permissões do catálogo fora do módulo Develop.
//   - Boots seguintes: cria só as permissões novas do catálogo e vincula ao ADMIN
//     as novas fora do módulo Develop. Não recoloca vínculo que o operador
//     removeu, não toca nas permissões existentes e não cria nenhuma outra role.
//
// Se não há role ADMIN fora do primeiro boot (o operador a removeu ou
// renomeou), as permissões novas são criadas mas nada é vinculado.
func SyncCatalog(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		created, err := syncPermissions(tx)
		if err != nil {
			return err
		}

		admin, adminCreated, err := ensureAdminRole(tx)
		if err != nil {
			return err
		}
		if admin == nil {
			return nil
		}

		var toLink []domain.Permission
		if adminCreated {
			// Primeiro boot: o catálogo todo, mesmo que alguma permissão já existisse.
			var codes []string
			for _, entry := range PermissionCatalog {
				if entry.Module != DevelopModule {
					codes = append(codes, entry.Permission)
				}
			}
			if err := tx.Where("permission IN ?", codes).Find(&toLink).Error; err != nil {
				return fmt.Errorf("buscando permissões do catálogo: %w", err)
			}
		} else {
			for _, permission := range created {
				if permission.Module != DevelopModule {
					toLink = append(toLink, permission)
				}
			}
		}
		if len(toLink) == 0 {
			return nil
		}
		if err := tx.Model(admin).Association("Permissions").Append(&toLink); err != nil {
			return fmt.Errorf("vinculando permissões ao %s: %w", AdminRoleName, err)
		}
		return nil
	})
}

// ensureAdminRole devolve a role ADMIN e se ela foi criada agora. Só cria no
// primeiro boot (tabela roles sem nenhuma linha, contando as excluídas); fora
// dele, devolve nil se o ADMIN não existe ou foi excluído — decisão do operador.
func ensureAdminRole(tx *gorm.DB) (*domain.Role, bool, error) {
	var admin domain.Role
	err := tx.Unscoped().Where("name = ?", AdminRoleName).First(&admin).Error
	if err == nil {
		if admin.DeletedAt.Valid {
			return nil, false, nil
		}
		return &admin, false, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, fmt.Errorf("buscando a role %s: %w", AdminRoleName, err)
	}

	var total int64
	if err := tx.Unscoped().Model(&domain.Role{}).Count(&total).Error; err != nil {
		return nil, false, fmt.Errorf("contando roles: %w", err)
	}
	if total > 0 {
		return nil, false, nil
	}

	admin = domain.Role{Name: AdminRoleName, Description: AdminRoleDescription}
	result := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&admin)
	if result.Error != nil {
		return nil, false, fmt.Errorf("criando a role %s: %w", AdminRoleName, result.Error)
	}
	if result.RowsAffected == 0 {
		// Outra instância criou a role entre a busca e o insert: ela cuida dos vínculos.
		return nil, false, nil
	}
	return &admin, true, nil
}
