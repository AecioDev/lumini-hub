package repository

import (
	"errors"
	"simple-erp-service/internal/data-structure/models"

	"gorm.io/gorm"
)

// RoleRepository define as operações de acesso a dados para perfis
type RoleRepository interface {
	Repository
	FindAll() ([]models.Role, error)
	FindByID(id uint) (*models.Role, error)
	FindByIDWithPermissions(id uint) (*models.Role, error)
	FindByName(name string) (*models.Role, error)
	Create(role *models.Role) error
	Update(role *models.Role) error
	Delete(id uint) error
	ExistsByName(name string) (bool, error)
	ExistsByNameExcept(name string, id uint) (bool, error)
	UpdatePermissions(role *models.Role, permissionIDs []uint) error
	CountByPermissionID(permissionID uint) (int64, error)
}

// GormRoleRepository implementa RoleRepository usando GORM
type GormRoleRepository struct {
	*BaseRepository
}

// NewRoleRepository cria um novo repository de perfis
func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &GormRoleRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// FindAll retorna todos os perfis com paginação
func (r *GormRoleRepository) FindAll() ([]models.Role, error) {
	var roles []models.Role
	if err := r.GetDB().Find(&roles).Error; err != nil {
		return nil, err
	}

	return roles, nil
}

// FindByID busca um perfil pelo ID
func (r *GormRoleRepository) FindByID(id uint) (*models.Role, error) {
	var role models.Role
	if err := r.GetDB().First(&role, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

// FindByIDWithPermissions busca um perfil pelo ID e carrega o relacionamento com Permissions
func (r *GormRoleRepository) FindByIDWithPermissions(id uint) (*models.Role, error) {
	var role models.Role
	if err := r.GetDB().Preload("Permissions").First(&role, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

// FindByName busca um perfil pelo nome
func (r *GormRoleRepository) FindByName(name string) (*models.Role, error) {
	var role models.Role
	if err := r.GetDB().Where("name = ?", name).First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

// Create cria um novo perfil
func (r *GormRoleRepository) Create(role *models.Role) error {
	return r.GetDB().Create(role).Error
}

// Update atualiza um perfil existente
func (r *GormRoleRepository) Update(role *models.Role) error {
	return r.GetDB().Save(role).Error
}

// Delete exclui um perfil
func (r *GormRoleRepository) Delete(id uint) error {
	return r.GetDB().Delete(&models.Role{}, id).Error
}

// ExistsByName verifica se existe um perfil com o nome especificado
func (r *GormRoleRepository) ExistsByName(name string) (bool, error) {
	var count int64
	err := r.GetDB().Model(&models.Role{}).Where("name = ?", name).Count(&count).Error
	return count > 0, err
}

// ExistsByNameExcept verifica se existe um perfil com o nome especificado, exceto o perfil com o ID especificado
func (r *GormRoleRepository) ExistsByNameExcept(name string, id uint) (bool, error) {
	var count int64
	err := r.GetDB().Model(&models.Role{}).Where("name = ? AND id != ?", name, id).Count(&count).Error
	return count > 0, err
}

// UpdatePermissions atualiza as permissões de um perfil
func (r *GormRoleRepository) UpdatePermissions(role *models.Role, permissionIDs []uint) error {
	var permissions []models.Permission
	if err := r.GetDB().Where("id IN ?", permissionIDs).Find(&permissions).Error; err != nil {
		return err
	}

	// Verificar se todas as permissões solicitadas existem
	if len(permissions) != len(permissionIDs) {
		return errors.New("uma ou mais permissões não existem")
	}

	// Atualizar permissões do perfil
	return r.GetDB().Model(role).Association("Permissions").Replace(&permissions)
}

// CountByPermissionID conta quantos Papeis estão usando uma determinada permissão
func (r *GormRoleRepository) CountByPermissionID(permissionID uint) (int64, error) {
	var count int64
	err := r.GetDB().Model(&models.RolePermissions{}).Where("permission_id = ?", permissionID).Count(&count).Error
	return count, err
}
