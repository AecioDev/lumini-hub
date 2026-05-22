package repository

import (
	"errors"

	"lumini-hub/api.auth/internal/domain"
	"lumini-hub/common/utils"

	"gorm.io/gorm"
)

// PermissionRepository define as operações de acesso a dados para permissões
type PermissionRepository interface {
	Repository
	FindAll() ([]domain.Permission, error)
	FindByID(id uint) (*domain.Permission, error)
	FindByIDs(ids []uint) ([]domain.Permission, error)
	FindAllFiltered(pagination *utils.Pagination, filters domain.InGetPermissionsFilters) ([]domain.Permission, error)
	GroupByModule() (map[string][]domain.Permission, error)
	FindAllModules() ([]string, error)
	Create(permission *domain.Permission) error
	Update(permission *domain.Permission) error
	Delete(id uint) error
	ExistsByName(name string) (bool, error)
	ExistsByNameExcept(name string, id uint) (bool, error)
}

// GormPermissionRepository implementa PermissionRepository usando GORM
type GormPermissionRepository struct {
	*BaseRepository
}

// NewPermissionRepository cria um novo repository de permissões
func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &GormPermissionRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// FindAll retorna todas as permissões sem paginação
func (r *GormPermissionRepository) FindAll() ([]domain.Permission, error) {
	var permissions []domain.Permission
	if err := r.GetDB().Find(&permissions).Error; err != nil {
		return nil, err
	}

	return permissions, nil
}

// FindAllFiltered retorna todas as permissões paginadas e filtradas
func (r *GormPermissionRepository) FindAllFiltered(pagination *utils.Pagination, filters domain.InGetPermissionsFilters) ([]domain.Permission, error) {
	var permissions []domain.Permission
	query := r.GetDB().Model(&domain.Permission{})

	// Aplicar filtros
	if filters.Name != "" {
		query = query.Where("permission ILIKE ?", "%"+filters.Name+"%")
	}
	if filters.Module != "" {
		query = query.Where("module = ?", filters.Module)
	}
	if filters.RoleId != 0 {
		if filters.IsLinkedToRole != nil && *filters.IsLinkedToRole == false {
			query = query.Where("NOT EXISTS (?)",
				r.GetDB().Table("role_permissions").
					Select("1").
					Where("role_permissions.permission_id = permissions.id").
					Where("role_permissions.role_id = ?", filters.RoleId),
			)
		} else {
			query = query.Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
				Where("role_permissions.role_id = ?", filters.RoleId)
		}
	}

	// Aplicar paginação
	paginatedQuery, err := utils.Paginate(&domain.Permission{}, pagination, query)
	if err != nil {
		return nil, err
	}

	if err := paginatedQuery.Find(&permissions).Error; err != nil {
		return nil, err
	}

	return permissions, nil
}

// FindByID busca uma permissão pelo ID
func (r *GormPermissionRepository) FindByID(id uint) (*domain.Permission, error) {
	var permission domain.Permission
	if err := r.GetDB().First(&permission, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &permission, nil
}

// FindByIDs busca permissões pelos IDs
func (r *GormPermissionRepository) FindByIDs(ids []uint) ([]domain.Permission, error) {
	var permissions []domain.Permission
	if err := r.GetDB().Where("id IN ?", ids).Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

// GroupByModule retorna permissões agrupadas por módulo
func (r *GormPermissionRepository) GroupByModule() (map[string][]domain.Permission, error) {
	var permissions []domain.Permission
	if err := r.GetDB().Find(&permissions).Error; err != nil {
		return nil, err
	}

	// Agrupar permissões por módulo
	moduleMap := make(map[string][]domain.Permission)
	for _, perm := range permissions {
		moduleMap[perm.Module] = append(moduleMap[perm.Module], perm)
	}

	return moduleMap, nil
}

// FindAllModules retorna a lista de módulos cadastrados
func (r *GormPermissionRepository) FindAllModules() ([]string, error) {
	var modules []string
	if err := r.GetDB().Model(&domain.Permission{}).Distinct("module").Pluck("module", &modules).Error; err != nil {
		return nil, err
	}
	return modules, nil
}

// Create cria uma nova permissão
func (r *GormPermissionRepository) Create(permission *domain.Permission) error {
	return r.GetDB().Create(permission).Error
}

// Update atualiza uma permissão existente
func (r *GormPermissionRepository) Update(permission *domain.Permission) error {
	return r.GetDB().Save(permission).Error
}

// Delete exclui uma permissão
func (r *GormPermissionRepository) Delete(id uint) error {
	return r.GetDB().Delete(&domain.Permission{}, id).Error
}

// ExistsByName verifica se existe uma permissão com o nome especificado
func (r *GormPermissionRepository) ExistsByName(name string) (bool, error) {
	var count int64
	err := r.GetDB().Model(&domain.Permission{}).Where("permission = ?", name).Count(&count).Error
	return count > 0, err
}

// ExistsByNameExcept verifica se existe uma permissão com o nome especificado, exceto a permissão com o ID especificado
func (r *GormPermissionRepository) ExistsByNameExcept(name string, id uint) (bool, error) {
	var count int64
	err := r.GetDB().Model(&domain.Permission{}).Where("permission = ? AND id != ?", name, id).Count(&count).Error
	return count > 0, err
}
