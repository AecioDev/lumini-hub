package service

import (
	"lumini-hub/api.auth/internal/domain"
	"lumini-hub/api.auth/internal/repository"
	"lumini-hub/api.auth/internal/validator"
	"lumini-hub/common/utils"
)

// UserService gerencia operações relacionadas a usuários
type UserService struct {
	userRepo  repository.UserRepository
	roleRepo  repository.RoleRepository
	validator *validator.UserValidator
}

// NewUserService cria um novo serviço de usuários
func NewUserService(
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
) *UserService {
	return &UserService{
		userRepo:  userRepo,
		roleRepo:  roleRepo,
		validator: validator.NewUserValidator(userRepo, roleRepo),
	}
}

// GetUsers retorna uma lista paginada de usuários
func (s *UserService) GetUsers(pagination *utils.Pagination) (*domain.ApiUserListPaginated, error) {
	users, err := s.userRepo.FindAll(pagination)
	if err != nil {
		return nil, err
	}

	// Converter para DTOs
	userDTOs := make([]domain.ApiUser, 0, len(users))
	for _, user := range users {
		userDTOs = append(userDTOs, domain.ApiUserFromModel(user))
	}

	return &domain.ApiUserListPaginated{
		Users:      userDTOs,
		Pagination: utils.ApiPaginationFromModel(pagination),
	}, nil
}

// GetUserByID busca um usuário pelo ID
func (s *UserService) GetUserByID(id uint) (*domain.ApiUserDetail, error) {
	user, err := s.userRepo.FindByIDWithRole(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, utils.ErrNotFound
	}

	// Converter para DTO
	userDetailDTO := domain.ApiUserDetailFromModel(*user)
	return &userDetailDTO, nil
}

// CreateUser cria um novo usuário
func (s *UserService) CreateUser(req domain.CreateUserRequest) (*domain.ApiUser, error) {
	// Validar dados
	if err := s.validator.ValidateForCreation(req); err != nil {
		return nil, err
	}

	// Hash da senha
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Criar usuário
	user := domain.User{
		Username:     req.Username,
		PasswordHash: passwordHash,
		Name:         req.Name,
		Email:        req.Email,
		Phone:        req.Phone,
		RoleID:       req.RoleID,
		IsActive:     true, // Por padrão, usuários são criados ativos
	}

	if err := s.userRepo.Create(&user); err != nil {
		return nil, err
	}

	// Buscar usuário completo com relacionamentos
	completeUser, err := s.userRepo.FindByIDWithRole(user.ID)
	if err != nil {
		return nil, err
	}

	// Converter para DTO
	userDTO := domain.ApiUserFromModel(*completeUser)
	return &userDTO, nil
}

// UpdateUser atualiza um usuário existente
func (s *UserService) UpdateUser(id uint, req domain.UpdateUserRequest) (*domain.ApiUser, error) {
	// Validar dados
	if err := s.validator.ValidateForUpdate(id, req); err != nil {
		return nil, err
	}

	// Buscar usuário
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, utils.ErrNotFound
	}

	// Atualizar campos
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.RoleID != 0 {
		user.RoleID = req.RoleID
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	// Salvar alterações
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	// Buscar usuário completo com relacionamentos
	completeUser, err := s.userRepo.FindByIDWithRole(user.ID)
	if err != nil {
		return nil, err
	}

	// Converter para DTO
	userDTO := domain.ApiUserFromModel(*completeUser)
	return &userDTO, nil
}

// ChangePassword altera a senha de um usuário
func (s *UserService) ChangePassword(id uint, currentPassword, newPassword string, isAdmin bool) error {
	// Validar dados
	req := domain.ChangePasswordRequest{
		CurrentPassword: currentPassword,
		NewPassword:     newPassword,
	}
	if err := s.validator.ValidatePasswordChange(id, req, isAdmin); err != nil {
		return err
	}

	// Buscar usuário
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return utils.ErrNotFound
	}

	// Se não for admin, verificar a senha atual
	if !isAdmin && !utils.CheckPasswordHash(currentPassword, user.PasswordHash) {
		return utils.ErrInvalidCredentials
	}

	// Hash da nova senha
	passwordHash, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Atualizar senha
	user.PasswordHash = passwordHash
	return s.userRepo.Update(user)
}

// DeleteUser exclui um usuário (soft delete)
func (s *UserService) DeleteUser(id uint) error {
	// Verificar se o usuário existe
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return utils.ErrNotFound
	}

	// Excluir usuário
	return s.userRepo.Delete(id)
}
