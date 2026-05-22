package validator

import (
	"simple-erp-service/internal/data-structure/models"
	"simple-erp-service/internal/repository"
)

// RoleValidator valida regras de negócio relacionadas a perfis
type RoleValidator struct {
	roleRepo repository.RoleRepository
	userRepo repository.UserRepository
	permRepo repository.PermissionRepository
}

// NewRoleValidator cria um novo validador de perfis
func NewRoleValidator(roleRepo repository.RoleRepository, userRepo repository.UserRepository, permRepo repository.PermissionRepository) *RoleValidator {
	return &RoleValidator{
		roleRepo: roleRepo,
		userRepo: userRepo,
		permRepo: permRepo,
	}
}

// ValidateForCreation valida os dados para criação de um perfil
func (v *RoleValidator) ValidateForCreation(req models.CreateRoleRequest) error {
	var errors ValidationErrors

	// Verificar se o nome já existe
	exists, err := v.roleRepo.ExistsByName(req.Name)
	if err != nil {
		return err
	}
	if exists {
		errors.AddError("name", "nome de perfil já está em uso")
	}

	if errors.HasErrors() {
		return errors
	}
	return nil
}

// ValidateForUpdate valida os dados para atualização de um perfil
func (v *RoleValidator) ValidateForUpdate(id uint, req models.UpdateRoleRequest) error {
	var errors ValidationErrors

	// Verificar se o perfil existe
	role, err := v.roleRepo.FindByID(id)
	if err != nil {
		return err
	}
	if role == nil {
		errors.AddError("id", "perfil não encontrado")
		return errors
	}

	// Verificar se o nome já está em uso por outro perfil
	if req.Name != "" && req.Name != role.Name {
		exists, err := v.roleRepo.ExistsByNameExcept(req.Name, id)
		if err != nil {
			return err
		}
		if exists {
			errors.AddError("name", "nome de perfil já está em uso")
		}
	}

	if errors.HasErrors() {
		return errors
	}
	return nil
}

// ValidateForDeletion valida se um perfil pode ser excluído
func (v *RoleValidator) ValidateForDeletion(id uint) error {
	var errors ValidationErrors

	// Verificar se o perfil existe
	role, err := v.roleRepo.FindByID(id)
	if err != nil {
		return err
	}
	if role == nil {
		errors.AddError("id", "perfil não encontrado")
		return errors
	}

	// Verificar se o perfil está sendo usado por usuários
	count, err := v.userRepo.CountByRoleID(id)
	if err != nil {
		return err
	}
	if count > 0 {
		errors.AddError("id", "não é possível excluir um perfil que está sendo usado por usuários")
	}

	if errors.HasErrors() {
		return errors
	}
	return nil
}

// ValidatePermissionUpdate valida a atualização de permissões de um perfil
func (v *RoleValidator) ValidatePermissionUpdate(id uint, permissionIDs []uint) error {
	var errors ValidationErrors

	// Verificar se o perfil existe
	role, err := v.roleRepo.FindByID(id)
	if err != nil {
		return err
	}
	if role == nil {
		errors.AddError("id", "perfil não encontrado")
		return errors
	}

	// Verificar se todas as permissões existem
	permissions, err := v.permRepo.FindByIDs(permissionIDs)
	if err != nil {
		return err
	}
	if len(permissions) != len(permissionIDs) {
		errors.AddError("permission_ids", "uma ou mais permissões não existem")
	}

	if errors.HasErrors() {
		return errors
	}
	return nil
}
