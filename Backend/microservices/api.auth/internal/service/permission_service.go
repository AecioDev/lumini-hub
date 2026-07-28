package service

import (
	"lumini-hub/api.auth/internal/domain"
	"lumini-hub/api.auth/internal/repository"
	"lumini-hub/api.auth/internal/validator"
	"lumini-hub/common/utils"
)

// PermissionService gerencia operações relacionadas a permissões de usuário
type PermissionService struct {
	roleRepo  repository.RoleRepository
	permRepo  repository.PermissionRepository
	validator *validator.PermissionValidator
}

// NewPermissionService cria um novo serviço de permissões
func NewPermissionService(
	roleRepo repository.RoleRepository,
	permRepo repository.PermissionRepository,
) *PermissionService {
	return &PermissionService{
		roleRepo:  roleRepo,
		permRepo:  permRepo,
		validator: validator.NewPermissionValidator(roleRepo, permRepo),
	}
}

// GetPermissions retorna todas as permissões. O módulo utils.DeveloperModule
// (roles.*/permissions.*/admin.create_permissions) fica escondido de quem
// não é DEVELOP.
func (s *PermissionService) GetPermissions(pagination *utils.Pagination, filters domain.InGetPermissionsFilters, requesterRole string) (*domain.ApiPermissionListPaginated, error) {
	permissions, err := s.permRepo.FindAllFiltered(pagination, filters, excludeModuleFor(requesterRole))
	if err != nil {
		return nil, err
	}

	// Converter para DTOs
	permissionDTOs := make([]domain.ApiPermission, 0, len(permissions))
	for _, perm := range permissions {
		permissionDTOs = append(permissionDTOs, domain.ApiPermissionFromModel(perm))
	}

	return &domain.ApiPermissionListPaginated{
		Permission: permissionDTOs,
		Pagination: *utils.ApiPaginationFromModel(pagination),
	}, nil
}

// GetPermissionByID busca uma permissão pelo ID
func (s *PermissionService) GetPermissionByID(id uint) (*domain.ApiPermissionDetail, error) {
	permission, err := s.permRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if permission == nil {
		return nil, utils.ErrNotFound
	}

	// Converter para DTO
	permissionDetailDTO := domain.ApiPermissionDetailFromModel(*permission)
	return &permissionDetailDTO, nil
}

// GetPermissionsByModule retorna permissões agrupadas por módulo. O módulo
// utils.DeveloperModule fica escondido de quem não é DEVELOP — inclusive
// aqui, já que este endpoint alimenta o seletor de permissões usado ao
// editar um Perfil (senão um ADMIN poderia atribuir permissões de DEVELOP a
// um Perfil customizado).
func (s *PermissionService) GetPermissionsByModule(requesterRole string) ([]domain.ApiPermissionsByModule, error) {
	moduleMap, err := s.permRepo.GroupByModule(excludeModuleFor(requesterRole))
	if err != nil {
		return nil, err
	}

	// Converter mapa para slice
	var result []domain.PermissionsByModule
	for module, perms := range moduleMap {
		result = append(result, domain.PermissionsByModule{
			Module:      module,
			Permissions: perms,
		})
	}

	// Converter para DTOs
	resultDTOs := make([]domain.ApiPermissionsByModule, 0, len(result))
	for _, permission := range result {
		resultDTOs = append(resultDTOs, domain.ApiPermissionByModuleFromModel(permission))
	}

	return resultDTOs, nil
}

// GetAvailableModules retorna uma lista de módulos cadastrados. O módulo
// utils.DeveloperModule fica escondido de quem não é DEVELOP.
func (s *PermissionService) GetAvailableModules(requesterRole string) ([]string, error) {
	modules, err := s.permRepo.FindAllModules(excludeModuleFor(requesterRole))
	if err != nil {
		return nil, err
	}
	return modules, nil
}

// excludeModuleFor devolve utils.DeveloperModule pra qualquer requisitante
// que não seja DEVELOP (nenhuma exclusão pro próprio DEVELOP).
func excludeModuleFor(requesterRole string) string {
	if requesterRole == utils.RoleDeveloper {
		return ""
	}
	return utils.DeveloperModule
}

// CreatePermission cria uma nova permissão
func (s *PermissionService) CreatePermission(req domain.InCreatePermission) (*domain.ApiPermission, error) {
	// Validar dados
	if err := s.validator.ValidateForCreation(req); err != nil {
		return nil, err
	}

	// Criar permissão
	permission := domain.Permission{
		Permission:  req.Permission,
		Description: req.Description,
		Module:      req.Module,
	}

	if err := s.permRepo.Create(&permission); err != nil {
		return nil, err
	}

	// Converter para DTO
	permissionDTO := domain.ApiPermissionFromModel(permission)
	return &permissionDTO, nil
}

// UpdatePermission atualiza uma permissão existente
func (s *PermissionService) UpdatePermission(id uint, req domain.InUpdatePermission) (*domain.ApiPermission, error) {
	// Validar dados
	if err := s.validator.ValidateForUpdate(id, req); err != nil {
		return nil, err
	}

	// Buscar a permissão
	permission, err := s.permRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if permission == nil {
		return nil, utils.ErrNotFound
	}

	// Atualizar campos
	if req.Permission != nil {
		permission.Permission = *req.Permission
	}

	if req.Description != nil {
		permission.Description = *req.Description
	}

	if req.Module != nil {
		permission.Module = *req.Module
	}

	// Salvar alterações
	if err := s.permRepo.Update(permission); err != nil {
		return nil, err
	}

	// Converter para DTO
	permDTO := domain.ApiPermissionFromModel(*permission)
	return &permDTO, nil
}

// DeletePermission exclui uma permissão
func (s *PermissionService) DeletePermission(id uint) error {
	// Validar se a permissão pode ser excluída
	if err := s.validator.ValidateForDeletion(id); err != nil {
		return err
	}

	// Excluir permissão
	return s.permRepo.Delete(id)
}
