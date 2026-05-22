package repository

import (
	"errors"
	"simple-erp-service/internal/data-structure/dto"
	"simple-erp-service/internal/data-structure/models"
	"simple-erp-service/internal/utils"

	"gorm.io/gorm"
)

// PermissionRepository define as operações de acesso a dados para permissões
type PermissionRepository interface {
	Repository
	FindAll() ([]models.Permission, error)
	FindByID(id uint) (*models.Permission, error)
	FindByIDs(ids []uint) ([]models.Permission, error)
	FindAllFiltered(pagination *models.Pagination, filters dto.InGetPermissionsFilters) ([]models.Permission, error)
	GroupByModule() (map[string][]models.Permission, error)
	FindAllModules() ([]string, error)
	Create(permission *models.Permission) error
	Update(permission *models.Permission) error
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
func (r *GormPermissionRepository) FindAll() ([]models.Permission, error) {
	var permissions []models.Permission
	if err := r.GetDB().Find(&permissions).Error; err != nil {
		return nil, err
	}

	return permissions, nil
}

// FindAllFiltered retorna todas as permissões paginadas e filtradas
func (r *GormPermissionRepository) FindAllFiltered(pagination *models.Pagination, filters dto.InGetPermissionsFilters) ([]models.Permission, error) {
	var permissions []models.Permission
	query := r.GetDB().Model(&models.Permission{})

	// Aplicar filtros
	if filters.Name != "" {
		query = query.Where("permission ILIKE ?", "%"+filters.Name+"%") // ILIKE para case-insensitive no PostgreSQL/MySQL
	}
	if filters.Module != "" {
		query = query.Where("module = ?", filters.Module)
	}
	if filters.RoleId != 0 {
		if filters.IsLinkedToRole != nil && *filters.IsLinkedToRole == false { // Se isLinkedToRole foi explicitamente setado para FALSE
			// Queremos permissões NÃO vinculadas a este RoleId.
			// Usamos NOT EXISTS para encontrar permissões que NÃO TÊM uma associação com o RoleId dado.
			query = query.Where("NOT EXISTS (?)",
				r.GetDB().Table("role_permissions").
					Select("1").
					Where("role_permissions.permission_id = permissions.id").
					Where("role_permissions.role_id = ?", filters.RoleId),
			)
		} else { // Se isLinkedToRole é TRUE, ou não foi fornecido (e o padrão é buscar vinculadas)
			// Queremos permissões VINCULADAS a este RoleId.
			// Fazemos um JOIN com a tabela de associação para garantir que a permissão está ligada ao RoleId.
			query = query.Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
				Where("role_permissions.role_id = ?", filters.RoleId)
		}
	}

	// Aplicar paginação
	query, err := utils.Paginate(&models.Permission{}, pagination, query)
	if err != nil {
		return nil, err
	}

	if err := query.Find(&permissions).Error; err != nil {
		return nil, err
	}

	return permissions, nil
}

// FindByID busca um perfil pelo ID
func (r *GormPermissionRepository) FindByID(id uint) (*models.Permission, error) {
	var permission models.Permission
	if err := r.GetDB().First(&permission, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &permission, nil
}

// FindByIDs busca permissões pelos IDs
func (r *GormPermissionRepository) FindByIDs(ids []uint) ([]models.Permission, error) {
	var permissions []models.Permission
	if err := r.GetDB().Where("id IN ?", ids).Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

// GroupByModule retorna permissões agrupadas por módulo
func (r *GormPermissionRepository) GroupByModule() (map[string][]models.Permission, error) {
	var permissions []models.Permission
	if err := r.GetDB().Find(&permissions).Error; err != nil {
		return nil, err
	}

	// Agrupar permissões por módulo
	moduleMap := make(map[string][]models.Permission)
	for _, perm := range permissions {
		moduleMap[perm.Module] = append(moduleMap[perm.Module], perm)
	}

	return moduleMap, nil
}

// Retorna a Lista de Modulos Cadastrados
func (r *GormPermissionRepository) FindAllModules() ([]string, error) {
	var modules []string
	// Seleciona distintamente a coluna 'module' da tabela 'permissions'
	if err := r.GetDB().Model(&models.Permission{}).Distinct("module").Pluck("module", &modules).Error; err != nil {
		return nil, err
	}
	return modules, nil
}

// Create cria uma nova permissão
func (r *GormPermissionRepository) Create(permission *models.Permission) error {
	return r.GetDB().Create(permission).Error
}

// Update atualiza uma permissão existente
func (r *GormPermissionRepository) Update(permission *models.Permission) error {
	return r.GetDB().Save(permission).Error
}

// Delete exclui uma permissão
func (r *GormPermissionRepository) Delete(id uint) error {
	return r.GetDB().Delete(&models.Permission{}, id).Error
}

// ExistsByName verifica se existe uma permissão com o nome especificado
func (r *GormPermissionRepository) ExistsByName(name string) (bool, error) {
	var count int64
	err := r.GetDB().Model(&models.Permission{}).Where("permission = ?", name).Count(&count).Error
	return count > 0, err
}

// ExistsByNameExcept verifica se existe um perfil com o nome especificado, exceto o perfil com o ID especificado
func (r *GormPermissionRepository) ExistsByNameExcept(name string, id uint) (bool, error) {
	var count int64
	err := r.GetDB().Model(&models.Permission{}).Where("permission = ? AND id != ?", name, id).Count(&count).Error
	return count > 0, err
}
