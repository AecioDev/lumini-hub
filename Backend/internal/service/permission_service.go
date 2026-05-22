package service

import (
	dto "simple-erp-service/internal/data-structure/dto"
	"simple-erp-service/internal/data-structure/models"
	"simple-erp-service/internal/repository"
	"simple-erp-service/internal/utils"
	"simple-erp-service/internal/validator"
)

// PermissionService gerencia operações relacionadas a permissões de usuário
type PermissionService struct {
	roleRepo  repository.RoleRepository
	permRepo  repository.PermissionRepository
	validator *validator.PermissionValidator
}

// NewRoleService cria um novo serviço de permissões
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

// GetPermissions retorna todas as permissões
func (s *PermissionService) GetPermissions(pagination *models.Pagination, filters dto.InGetPermissionsFilters) (*dto.ApiPermissionListPaginated, error) {
	permissions, err := s.permRepo.FindAllFiltered(pagination, filters)
	if err != nil {
		return nil, err
	}

	// Converter para DTOs
	permissionDTOs := make([]dto.ApiPermission, 0, len(permissions))
	for _, perm := range permissions {
		permissionDTOs = append(permissionDTOs, dto.ApiPermissionFromModel(perm))
	}

	return &dto.ApiPermissionListPaginated{
		Permission: permissionDTOs,
		Pagination: *dto.ApiPaginationFromModel(pagination),
	}, nil
}

// GetPermissionByID busca um perfil pelo ID
func (s *PermissionService) GetPermissionByID(id uint) (*dto.ApiPermissionDetail, error) {
	permission, err := s.permRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if permission == nil {
		return nil, utils.ErrNotFound
	}

	// Converter para DTO
	permissionDetailDTO := dto.ApiPermissionDetailFromModel(*permission)
	return &permissionDetailDTO, nil
}

// GetPermissionsByModule retorna permissões agrupadas por módulo
func (s *PermissionService) GetPermissionsByModule() ([]dto.ApiPermissionsByModule, error) {
	moduleMap, err := s.permRepo.GroupByModule()
	if err != nil {
		return nil, err
	}

	// Converter mapa para slice
	var result []models.PermissionsByModule
	for module, perms := range moduleMap {
		result = append(result, models.PermissionsByModule{
			Module:      module,
			Permissions: perms,
		})
	}

	// Converter para DTOs
	resultDTOs := make([]dto.ApiPermissionsByModule, 0, len(result))
	for _, permission := range result {
		resultDTOs = append(resultDTOs, dto.ApiPermissionByModuleFromModel(permission))
	}

	return resultDTOs, nil
}

// GetAvailableModules retorna uma lista de módulos cadastrados
func (s *PermissionService) GetAvailableModules() ([]string, error) {
	modules, err := s.permRepo.FindAllModules()
	if err != nil {
		return nil, err
	}
	return modules, nil
}

// CreatePermission cria uma novo permissão
func (s *PermissionService) CreatePermission(req dto.InCreatePermission) (*dto.ApiPermission, error) {
	// Validar dados
	if err := s.validator.ValidateForCreation(req); err != nil {
		return nil, err
	}

	// Criar permissão
	permission := models.Permission{
		Permission:  req.Permission,
		Description: req.Description,
		Module:      req.Module,
	}

	if err := s.permRepo.Create(&permission); err != nil {
		return nil, err
	}

	// Converter para DTO
	permissionDTO := dto.ApiPermissionFromModel(permission)
	return &permissionDTO, nil
}

// UpdateRole atualiza um perfil existente
func (s *PermissionService) UpdatePermission(id uint, req dto.InUpdatePermission) (*dto.ApiPermission, error) {
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
	permDTO := dto.ApiPermissionFromModel(*permission)
	return &permDTO, nil
}

// DeletePermission exclui uma Permissão
func (s *PermissionService) DeletePermission(id uint) error {
	// Validar se o permissão pode ser excluída
	if err := s.validator.ValidateForDeletion(id); err != nil {
		return err
	}

	// Excluir perfil
	return s.permRepo.Delete(id)
}
