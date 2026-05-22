package validator

import (
	"simple-erp-service/internal/data-structure/dto"
	"simple-erp-service/internal/repository"
)

// PermissionValidator valida regras de negócio relacionadas a perfis
type PermissionValidator struct {
	roleRepo repository.RoleRepository
	permRepo repository.PermissionRepository
}

// NewRoleValidator cria um novo validador de perfis
func NewPermissionValidator(roleRepo repository.RoleRepository, permRepo repository.PermissionRepository) *PermissionValidator {
	return &PermissionValidator{
		roleRepo: roleRepo,
		permRepo: permRepo,
	}
}

// ValidateForCreation valida os dados para criação de um perfil
func (v *PermissionValidator) ValidateForCreation(req dto.InCreatePermission) error {
	var errors ValidationErrors

	// Verificar se o nome já existe
	exists, err := v.permRepo.ExistsByName(req.Permission)
	if err != nil {
		return err
	}
	if exists {
		errors.AddError("permission", "permissão informada já está em uso")
	}

	if errors.HasErrors() {
		return errors
	}
	return nil
}

// ValidateForUpdate valida os dados para atualização de um perfil
func (v *PermissionValidator) ValidateForUpdate(id uint, req dto.InUpdatePermission) error {
	var errors ValidationErrors

	// Verificar se a permissão existe
	permisison, err := v.permRepo.FindByID(id)
	if err != nil {
		return err
	}
	if permisison == nil {
		errors.AddError("id", "permissão não encontrada")
		return errors
	}

	// Verificar se o nome já está em uso por outra permissão
	if req.Permission != nil && req.Permission != &permisison.Permission {
		exists, err := v.permRepo.ExistsByNameExcept(*req.Permission, id)
		if err != nil {
			return err
		}
		if exists {
			errors.AddError("permission", "permissão informada já está em uso")
		}
	}

	if errors.HasErrors() {
		return errors
	}
	return nil
}

// ValidateForDeletion valida se uma permissão pode ser excluída
func (v *PermissionValidator) ValidateForDeletion(id uint) error {
	var errors ValidationErrors

	// Verificar se a permissão existe
	permisison, err := v.permRepo.FindByID(id)
	if err != nil {
		return err
	}
	if permisison == nil {
		errors.AddError("id", "permissão não encontrada")
		return errors
	}

	// Verificar se a permissão está sendo usada em algum papel
	count, err := v.roleRepo.CountByPermissionID(id)
	if err != nil {
		return err
	}
	if count > 0 {
		errors.AddError("id", "não é possível excluir uma permissão que está sendo usada em um papel")
	}

	if errors.HasErrors() {
		return errors
	}
	return nil
}
