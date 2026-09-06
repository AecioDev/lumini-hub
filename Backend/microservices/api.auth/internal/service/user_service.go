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
	permRepo  repository.PermissionRepository
	validator *validator.UserValidator
}

// NewUserService cria um novo serviço de usuários
func NewUserService(
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	permRepo repository.PermissionRepository,
) *UserService {
	return &UserService{
		userRepo:  userRepo,
		roleRepo:  roleRepo,
		permRepo:  permRepo,
		validator: validator.NewUserValidator(userRepo, roleRepo, permRepo),
	}
}

// GetUsers retorna uma lista paginada de usuários. Usuários do perfil DEVELOP
// ficam escondidos de qualquer requisitante que não seja ele mesmo DEVELOP
// (perfil "oculto no sistema" — ver utils.RoleDeveloper).
func (s *UserService) GetUsers(pagination *utils.Pagination, requesterRole string) (*domain.ApiUserListPaginated, error) {
	hideRoleName := ""
	if requesterRole != utils.RoleDeveloper {
		hideRoleName = utils.RoleDeveloper
	}

	users, err := s.userRepo.FindAll(pagination, hideRoleName)
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

// GetUserByID busca um usuário pelo ID. Se o usuário encontrado for do
// perfil DEVELOP e quem está pedindo não for DEVELOP, trata como não
// encontrado (não revela nem a existência do usuário).
func (s *UserService) GetUserByID(id uint, requesterRole string) (*domain.ApiUserDetail, error) {
	user, err := s.userRepo.FindByIDWithRole(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, utils.ErrNotFound
	}
	if hidden, err := s.isHiddenFromRequester(user, requesterRole); err != nil {
		return nil, err
	} else if hidden {
		return nil, utils.ErrNotFound
	}

	// Converter para DTO
	userDetailDTO := domain.ApiUserDetailFromModel(*user)
	return &userDetailDTO, nil
}

// guardDeveloperRoleAssignment impede que alguém que não seja DEVELOP
// atribua o perfil DEVELOP (oculto) a um usuário, seja na criação ou ao
// trocar o perfil de um usuário existente. Mensagem de erro genérica de
// propósito — não revela nem o nome do perfil escondido.
func (s *UserService) guardDeveloperRoleAssignment(roleID uint, requesterRole string) error {
	if requesterRole == utils.RoleDeveloper || roleID == 0 {
		return nil
	}
	developerRole, err := s.roleRepo.FindByName(utils.RoleDeveloper)
	if err != nil {
		return err
	}
	if developerRole != nil && roleID == developerRole.ID {
		var errs validator.ValidationErrors
		errs.AddError("role_id", "perfil inválido")
		return errs
	}
	return nil
}

// isHiddenFromRequester indica se `user` deve ficar invisível pra quem tem o
// papel `requesterRole` — hoje só o caso do perfil DEVELOP, oculto de
// qualquer requisitante que não seja ele mesmo DEVELOP. Funciona mesmo se
// `user.Role` não estiver pré-carregado (compara por RoleID).
func (s *UserService) isHiddenFromRequester(user *domain.User, requesterRole string) (bool, error) {
	if requesterRole == utils.RoleDeveloper {
		return false, nil
	}
	developerRole, err := s.roleRepo.FindByName(utils.RoleDeveloper)
	if err != nil {
		return false, err
	}
	if developerRole == nil {
		return false, nil
	}
	return user.RoleID == developerRole.ID, nil
}

// CreateUser cria um novo usuário
func (s *UserService) CreateUser(req domain.CreateUserRequest, requesterRole string) (*domain.ApiUser, error) {
	// Validar dados
	if err := s.validator.ValidateForCreation(req); err != nil {
		return nil, err
	}
	if err := s.guardDeveloperRoleAssignment(req.RoleID, requesterRole); err != nil {
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
		CompanyID:    req.CompanyID,
	}

	if err := s.userRepo.Create(&user); err != nil {
		return nil, err
	}

	// Copiar permissões do perfil (template) para o usuário
	role, err := s.roleRepo.FindByIDWithPermissions(req.RoleID)
	if err == nil && role != nil && len(role.Permissions) > 0 {
		var permIDs []uint
		for _, perm := range role.Permissions {
			permIDs = append(permIDs, perm.ID)
		}
		_ = s.userRepo.UpdatePermissions(&user, permIDs)
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
func (s *UserService) UpdateUser(id uint, req domain.UpdateUserRequest, requesterRole string) (*domain.ApiUser, error) {
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
	if hidden, err := s.isHiddenFromRequester(user, requesterRole); err != nil {
		return nil, err
	} else if hidden {
		return nil, utils.ErrNotFound
	}
	if req.RoleID != 0 {
		if err := s.guardDeveloperRoleAssignment(req.RoleID, requesterRole); err != nil {
			return nil, err
		}
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

	roleChanged := false
	if req.RoleID != 0 && req.RoleID != user.RoleID {
		user.RoleID = req.RoleID
		roleChanged = true
	}

	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	// Sempre atribui (não "if != nil, então mexe") — CompanyID é *uint
	// porque o valor de negócio nulo é "usuário master, sem empresa
	// vinculada", não "campo não enviado". O formulário do frontend sempre
	// manda o estado atual do Select, incluindo null explícito quando o
	// usuário limpa pra virar master — se essa atribuição fosse condicional
	// a `!= nil`, nunca daria pra desvincular a empresa de um usuário depois
	// de vinculada (achado testando: usuário limpava o Select, salvava, e a
	// empresa continuava lá).
	user.CompanyID = req.CompanyID

	// Salvar alterações
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	// Se o perfil mudou, resetar as permissões copiando as do novo perfil (template)
	if roleChanged {
		role, err := s.roleRepo.FindByIDWithPermissions(req.RoleID)
		if err == nil && role != nil && len(role.Permissions) > 0 {
			var permIDs []uint
			for _, perm := range role.Permissions {
				permIDs = append(permIDs, perm.ID)
			}
			_ = s.userRepo.UpdatePermissions(user, permIDs)
		}
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
func (s *UserService) DeleteUser(id uint, requesterRole string) error {
	// Verificar se o usuário existe
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return utils.ErrNotFound
	}
	if hidden, err := s.isHiddenFromRequester(user, requesterRole); err != nil {
		return err
	} else if hidden {
		return utils.ErrNotFound
	}

	// Excluir usuário
	return s.userRepo.Delete(id)
}

// UpdateUserPermissions atualiza as permissões diretas de um usuário
func (s *UserService) UpdateUserPermissions(id uint, permissionIDs []uint, requesterRole string) (*domain.ApiUserDetail, error) {
	// Validar a operação
	if err := s.validator.ValidatePermissionUpdate(id, permissionIDs); err != nil {
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
	if hidden, err := s.isHiddenFromRequester(user, requesterRole); err != nil {
		return nil, err
	} else if hidden {
		return nil, utils.ErrNotFound
	}

	// Atualizar permissões
	if err := s.userRepo.UpdatePermissions(user, permissionIDs); err != nil {
		return nil, err
	}

	// Buscar usuário completo com as permissões atualizadas
	completeUser, err := s.userRepo.FindByIDWithPermissions(id)
	if err != nil {
		return nil, err
	}

	userDetailDTO := domain.ApiUserDetailFromModel(*completeUser)
	return &userDetailDTO, nil
}
