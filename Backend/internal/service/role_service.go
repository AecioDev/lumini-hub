package service

import (
	dto "simple-erp-service/internal/data-structure/dto"
	"simple-erp-service/internal/data-structure/models"
	"simple-erp-service/internal/repository"
	"simple-erp-service/internal/utils"
	"simple-erp-service/internal/validator"
)

// RoleService gerencia operações relacionadas a perfis de usuário
type RoleService struct {
	roleRepo  repository.RoleRepository
	userRepo  repository.UserRepository
	permRepo  repository.PermissionRepository
	validator *validator.RoleValidator
}

// NewRoleService cria um novo serviço de perfis
func NewRoleService(
	roleRepo repository.RoleRepository,
	userRepo repository.UserRepository,
	permRepo repository.PermissionRepository,
) *RoleService {
	return &RoleService{
		roleRepo:  roleRepo,
		userRepo:  userRepo,
		permRepo:  permRepo,
		validator: validator.NewRoleValidator(roleRepo, userRepo, permRepo),
	}
}

// GetRoles retorna uma lista paginada de papeis
func (s *RoleService) GetRoles() ([]dto.ApiRole, error) {
	roles, err := s.roleRepo.FindAll()
	if err != nil {
		return nil, err
	}

	// Converter para DTOs
	roleDTOs := make([]dto.ApiRole, 0, len(roles))
	for _, role := range roles {
		roleDTOs = append(roleDTOs, dto.ApiRoleFromModel(role))
	}

	return roleDTOs, nil
}

// GetRoleByID busca um papel pelo ID
func (s *RoleService) GetRoleByID(id uint) (*dto.ApiRoleDetail, error) {
	role, err := s.roleRepo.FindByIDWithPermissions(id)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, utils.ErrNotFound
	}

	// Converter para DTO
	roleDetailDTO := dto.ApiRoleDetailFromModel(*role)
	return &roleDetailDTO, nil
}

// CreateRole cria um novo papel
func (s *RoleService) CreateRole(req models.CreateRoleRequest) (*dto.ApiRole, error) {
	// Validar dados
	if err := s.validator.ValidateForCreation(req); err != nil {
		return nil, err
	}

	// Criar papel
	role := models.Role{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.roleRepo.Create(&role); err != nil {
		return nil, err
	}

	// Converter para DTO
	roleDTO := dto.ApiRoleFromModel(role)
	return &roleDTO, nil
}

// UpdateRole atualiza um papel existente
func (s *RoleService) UpdateRole(id uint, req models.UpdateRoleRequest) (*dto.ApiRole, error) {
	// Validar dados
	if err := s.validator.ValidateForUpdate(id, req); err != nil {
		return nil, err
	}

	// Buscar papel
	role, err := s.roleRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, utils.ErrNotFound
	}

	// Atualizar campos
	if req.Name != "" {
		role.Name = req.Name
	}
	if req.Description != "" {
		role.Description = req.Description
	}

	// Salvar alterações
	if err := s.roleRepo.Update(role); err != nil {
		return nil, err
	}

	// Converter para DTO
	roleDTO := dto.ApiRoleFromModel(*role)
	return &roleDTO, nil
}

// DeleteRole exclui um papel
func (s *RoleService) DeleteRole(id uint) error {
	// Validar se o papel pode ser excluído
	if err := s.validator.ValidateForDeletion(id); err != nil {
		return err
	}

	// Excluir papel
	return s.roleRepo.Delete(id)
}

// UpdateRolePermissions atualiza as permissões de um papel
func (s *RoleService) UpdateRolePermissions(id uint, permissionIDs []uint) (*dto.ApiRoleDetail, error) {
	// Validar dados
	if err := s.validator.ValidatePermissionUpdate(id, permissionIDs); err != nil {
		return nil, err
	}

	// Buscar papel
	role, err := s.roleRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, utils.ErrNotFound
	}

	// Atualizar permissões
	if err := s.roleRepo.UpdatePermissions(role, permissionIDs); err != nil {
		return nil, err
	}

	// Buscar papel atualizado com permissões
	updatedRole, err := s.roleRepo.FindByIDWithPermissions(id)
	if err != nil {
		return nil, err
	}

	// Converter para DTO
	roleDetailDTO := dto.ApiRoleDetailFromModel(*updatedRole)
	return &roleDetailDTO, nil
}
