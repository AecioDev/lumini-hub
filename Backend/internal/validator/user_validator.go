package validator

import (
	"simple-erp-service/internal/data-structure/models"
	"simple-erp-service/internal/repository"
)

// UserValidator valida regras de negócio relacionadas a usuários
type UserValidator struct {
	userRepo repository.UserRepository
	roleRepo repository.RoleRepository
}

// NewUserValidator cria um novo validador de usuários
func NewUserValidator(userRepo repository.UserRepository, roleRepo repository.RoleRepository) *UserValidator {
	return &UserValidator{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

// ValidateForCreation valida os dados para criação de um usuário
func (v *UserValidator) ValidateForCreation(req models.CreateUserRequest) error {
	var errors ValidationErrors

	// Verificar se o username já existe
	exists, err := v.userRepo.ExistsByUsername(req.Username)
	if err != nil {
		return err
	}
	if exists {
		errors.AddError("username", "nome de usuário já está em uso")
	}

	// Verificar se o email já existe (se fornecido)
	if req.Email != "" {
		exists, err := v.userRepo.ExistsByEmail(req.Email)
		if err != nil {
			return err
		}
		if exists {
			errors.AddError("email", "email já está em uso")
		}
	}

	// Verificar se o perfil existe
	role, err := v.roleRepo.FindByID(req.RoleID)
	if err != nil {
		return err
	}
	if role == nil {
		errors.AddError("role_id", "perfil não encontrado")
	}

	// Verificar senha
	if len(req.Password) < 6 {
		errors.AddError("password", "senha deve ter pelo menos 6 caracteres")
	}

	if errors.HasErrors() {
		return errors
	}
	return nil
}

// ValidateForUpdate valida os dados para atualização de um usuário
func (v *UserValidator) ValidateForUpdate(id uint, req models.UpdateUserRequest) error {
	var errors ValidationErrors

	// Verificar se o usuário existe
	user, err := v.userRepo.FindByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		errors.AddError("id", "usuário não encontrado")
		return errors
	}

	// Verificar se o email já está em uso por outro usuário (se fornecido)
	if req.Email != "" && req.Email != user.Email {
		exists, err := v.userRepo.ExistsByEmailExcept(req.Email, id)
		if err != nil {
			return err
		}
		if exists {
			errors.AddError("email", "email já está em uso")
		}
	}

	// Verificar se o perfil existe (se fornecido)
	if req.RoleID != 0 && req.RoleID != user.RoleID {
		role, err := v.roleRepo.FindByID(req.RoleID)
		if err != nil {
			return err
		}
		if role == nil {
			errors.AddError("role_id", "perfil não encontrado")
		}
	}

	if errors.HasErrors() {
		return errors
	}
	return nil
}

// ValidatePasswordChange valida os dados para alteração de senha
func (v *UserValidator) ValidatePasswordChange(id uint, req models.ChangePasswordRequest, isAdmin bool) error {
	var errors ValidationErrors

	// Verificar se o usuário existe
	user, err := v.userRepo.FindByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		errors.AddError("id", "usuário não encontrado")
		return errors
	}

	// Verificar nova senha
	if len(req.NewPassword) < 6 {
		errors.AddError("new_password", "nova senha deve ter pelo menos 6 caracteres")
	}

	if errors.HasErrors() {
		return errors
	}
	return nil
}
