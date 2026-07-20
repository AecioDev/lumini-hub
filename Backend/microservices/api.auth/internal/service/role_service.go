package service

import (
	"lumini-hub/api.auth/internal/domain"
	"lumini-hub/api.auth/internal/repository"
	"lumini-hub/api.auth/internal/validator"
	"lumini-hub/common/utils"
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

// GetRoles retorna uma lista de papeis. O perfil DEVELOP (oculto no sistema)
// só aparece pra quem já está logado como DEVELOP.
func (s *RoleService) GetRoles(requesterRole string) ([]domain.ApiRole, error) {
	roles, err := s.roleRepo.FindAll()
	if err != nil {
		return nil, err
	}

	// Converter para DTOs
	roleDTOs := make([]domain.ApiRole, 0, len(roles))
	for _, role := range roles {
		if role.Name == utils.RoleDeveloper && requesterRole != utils.RoleDeveloper {
			continue
		}
		roleDTOs = append(roleDTOs, domain.ApiRoleFromModel(role))
	}

	return roleDTOs, nil
}

// GetRoleByID busca um papel pelo ID. Se for o perfil DEVELOP e quem pediu
// não for DEVELOP, trata como não encontrado.
func (s *RoleService) GetRoleByID(id uint, requesterRole string) (*domain.ApiRoleDetail, error) {
	role, err := s.roleRepo.FindByIDWithPermissions(id)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, utils.ErrNotFound
	}
	if role.Name == utils.RoleDeveloper && requesterRole != utils.RoleDeveloper {
		return nil, utils.ErrNotFound
	}

	// Converter para DTO
	roleDetailDTO := domain.ApiRoleDetailFromModel(*role)
	return &roleDetailDTO, nil
}

// CreateRole cria um novo papel
func (s *RoleService) CreateRole(req domain.CreateRoleRequest) (*domain.ApiRole, error) {
	// Validar dados
	if err := s.validator.ValidateForCreation(req); err != nil {
		return nil, err
	}

	// Criar papel
	role := domain.Role{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.roleRepo.Create(&role); err != nil {
		return nil, err
	}

	// Converter para DTO
	roleDTO := domain.ApiRoleFromModel(role)
	return &roleDTO, nil
}

// UpdateRole atualiza um papel existente
func (s *RoleService) UpdateRole(id uint, req domain.UpdateRoleRequest, requesterRole string) (*domain.ApiRole, error) {
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
	if role.Name == utils.RoleDeveloper && requesterRole != utils.RoleDeveloper {
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
	roleDTO := domain.ApiRoleFromModel(*role)
	return &roleDTO, nil
}

// DeleteRole exclui um papel
func (s *RoleService) DeleteRole(id uint, requesterRole string) error {
	// Buscar papel (pra checar visibilidade antes de validar exclusão)
	role, err := s.roleRepo.FindByID(id)
	if err != nil {
		return err
	}
	if role == nil {
		return utils.ErrNotFound
	}
	if role.Name == utils.RoleDeveloper && requesterRole != utils.RoleDeveloper {
		return utils.ErrNotFound
	}

	// Validar se o papel pode ser excluído
	if err := s.validator.ValidateForDeletion(id); err != nil {
		return err
	}

	// Excluir papel
	return s.roleRepo.Delete(id)
}

// UpdateRolePermissions atualiza as permissões de um papel
func (s *RoleService) UpdateRolePermissions(id uint, permissionIDs []uint, requesterRole string) (*domain.ApiRoleDetail, error) {
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
	if role.Name == utils.RoleDeveloper && requesterRole != utils.RoleDeveloper {
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
	roleDetailDTO := domain.ApiRoleDetailFromModel(*updatedRole)
	return &roleDetailDTO, nil
}
