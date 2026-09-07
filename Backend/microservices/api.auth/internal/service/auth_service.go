package service

import (
	"errors"
	"time"

	"lumini-hub/api.auth/internal/domain"
	"lumini-hub/api.auth/internal/repository"
	"lumini-hub/common/config"
	"lumini-hub/common/utils"

	"gorm.io/gorm"
)

// AuthService gerencia a autenticação de usuários
type AuthService struct {
	db  *gorm.DB
	cfg *config.Config
}

// NewAuthService cria um novo serviço de autenticação
func NewAuthService(db *gorm.DB, cfg *config.Config) *AuthService {
	return &AuthService{
		db:  db,
		cfg: cfg,
	}
}

// LoginResponse representa a resposta do login
type LoginResponse struct {
	User         domain.ApiUserDetail `json:"user"`
	AccessToken  string               `json:"access_token"`
	RefreshToken string               `json:"refresh_token"`
	ExpiresIn    int                  `json:"expires_in"`
}

// GetMenuItemsForUser monta a árvore de menu já filtrada pelas permissões do usuário
// (ADMIN vê a árvore completa, exceto o catálogo de Perfis e Permissões — ver
// domain.FilterMenuTreeForUser / utils.IsDeveloperOnlyPermission).
func (s *AuthService) GetMenuItemsForUser(user domain.User) ([]domain.ApiUserMenuItem, error) {
	menuItemRepo := repository.NewMenuItemRepository(s.db)
	flat, err := menuItemRepo.FindAllFlat()
	if err != nil {
		return nil, err
	}
	tree := domain.BuildMenuTree(flat)

	userPermissionCodes := make(map[string]bool, len(user.Permissions))
	for _, perm := range user.Permissions {
		userPermissionCodes[perm.Permission] = true
	}

	role := ""
	if user.Role != nil {
		role = user.Role.Name
	}
	return domain.FilterMenuTreeForUser(tree, userPermissionCodes, role), nil
}

// Login autentica um usuário e retorna tokens JWT
func (s *AuthService) Login(username, password string) (*LoginResponse, error) {
	var user domain.User

	// Buscar usuário pelo username
	result := s.db.Preload("Role").Preload("Role.Permissions").Preload("Permissions").Where("LOWER(username) = LOWER(?)", username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("usuário não encontrado")
		}
		return nil, result.Error
	}

	// Verificar se o usuário está ativo
	if !user.IsActive {
		return nil, errors.New("usuário inativo")
	}

	// Verifica se a senha informada corresponde ao hash armazenado
	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return nil, errors.New("senha incorreta")
	}

	// Gerar tokens
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Username, user.RoleID, user.Role.Name, s.cfg)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID, user.Username, s.cfg)
	if err != nil {
		return nil, err
	}

	// Atualizar último login
	now := time.Now()
	user.LastLogin = &now
	s.db.Save(&user)

	userDetail := domain.ApiUserDetailFromModel(user)
	if menuItems, err := s.GetMenuItemsForUser(user); err == nil {
		userDetail.MenuItems = menuItems
	}
	if err := s.ResolveActiveCompany(&userDetail); err != nil {
		return nil, err
	}

	return &LoginResponse{
		User:         userDetail,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(s.cfg.JWT.AccessTokenExp.Minutes()),
	}, nil
}

// RefreshToken renova o token de acesso usando um token de refresh
func (s *AuthService) RefreshToken(refreshToken string) (*LoginResponse, error) {
	// Validar token de refresh
	claims, err := utils.ValidateToken(refreshToken, s.cfg)
	if err != nil {
		return nil, err
	}

	// Buscar usuário pelo Subject (username)
	var user domain.User
	result := s.db.Preload("Role").Preload("Role.Permissions").Preload("Permissions").Where("LOWER(username) = LOWER(?)", claims.Subject).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	// Verificar se o usuário está ativo
	if !user.IsActive {
		return nil, errors.New("usuário inativo")
	}

	// Gerar novo token de acesso
	newAccessToken, err := utils.GenerateAccessToken(user.ID, user.Username, user.RoleID, user.Role.Name, s.cfg)
	if err != nil {
		return nil, err
	}

	// Gerar novo token de refresh
	newRefreshToken, err := utils.GenerateRefreshToken(user.ID, user.Username, s.cfg)
	if err != nil {
		return nil, err
	}

	userDetail := domain.ApiUserDetailFromModel(user)
	if menuItems, err := s.GetMenuItemsForUser(user); err == nil {
		userDetail.MenuItems = menuItems
	}
	if err := s.ResolveActiveCompany(&userDetail); err != nil {
		return nil, err
	}

	return &LoginResponse{
		User:         userDetail,
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int(s.cfg.JWT.AccessTokenExp.Minutes()),
	}, nil
}

// ResolveActiveCompany preenche ActiveCompanyID/ActiveCompanyName/
// RequiresCompanySelection/VisibleCompanies na resposta com o valor
// revalidado (nunca confiar cegamente no que está gravado em
// users.active_company_id — ver utils.ResolveActiveCompany). Chamado só nos
// pontos em que a resposta é sobre o PRÓPRIO usuário logado (login/refresh/
// me) — não faz sentido (nem vale o round-trip extra) pra quando um admin
// só está consultando o cadastro de outro usuário.
//
// VisibleCompanies/ActiveCompanyName vêm de utils.ResolveCompanyOptions, não
// de GET /companies — de propósito: saber em qual empresa o usuário está e
// trocar entre as que ele enxerga é identidade de sessão, não administração
// do cadastro, então não pode depender de companies.view (a maioria dos
// perfis operacionais nunca tem essa permission).
func (s *AuthService) ResolveActiveCompany(userDetail *domain.ApiUserDetail) error {
	activeID, requiresSelection, err := utils.ResolveActiveCompany(s.db, userDetail.ID, userDetail.ActiveCompanyID)
	if err != nil {
		return err
	}
	userDetail.ActiveCompanyID = activeID
	userDetail.RequiresCompanySelection = requiresSelection

	visibleIDs, unrestricted, err := utils.ResolveVisibleCompanyIDs(s.db, userDetail.ID)
	if err != nil {
		return err
	}
	options, err := utils.ResolveCompanyOptions(s.db, visibleIDs, unrestricted)
	if err != nil {
		return err
	}
	userDetail.VisibleCompanies = options

	if activeID != nil {
		for _, option := range options {
			if option.ID == *activeID {
				userDetail.ActiveCompanyName = option.Name
				break
			}
		}

		visualConfig, err := utils.ResolveActiveCompanyVisualConfig(s.db, *activeID)
		if err != nil {
			return err
		}
		userDetail.ActiveCompanyVisualConfig = visualConfig
	}

	return nil
}

// SetActiveCompany grava a empresa que o usuário está operando no momento
// (utils.SetActiveCompany valida que ele realmente a enxerga antes de
// gravar) e devolve o ApiUserDetail já atualizado, pronto pra resposta —
// mesma forma de login/refresh, pra o frontend poder trocar o contexto sem
// precisar dar reload na sessão inteira.
func (s *AuthService) SetActiveCompany(userID uint, companyID uint) (*domain.ApiUserDetail, error) {
	if err := utils.SetActiveCompany(s.db, userID, companyID); err != nil {
		return nil, err
	}

	var user domain.User
	if err := s.db.Preload("Role").Preload("Role.Permissions").Preload("Permissions").First(&user, userID).Error; err != nil {
		return nil, err
	}

	userDetail := domain.ApiUserDetailFromModel(user)
	if menuItems, err := s.GetMenuItemsForUser(user); err == nil {
		userDetail.MenuItems = menuItems
	}
	if err := s.ResolveActiveCompany(&userDetail); err != nil {
		return nil, err
	}

	return &userDetail, nil
}

// GetUserByID busca um usuário pelo ID
func (s *AuthService) GetUserByID(userID uint) (*domain.User, error) {
	var user domain.User
	result := s.db.Preload("Role").Preload("Role.Permissions").Preload("Permissions").First(&user, userID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
