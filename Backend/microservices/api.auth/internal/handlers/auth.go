package handlers

import (
	"net/http"
	"time"

	"lumini-hub/api.auth/internal/domain"
	"lumini-hub/api.auth/internal/service"
	"lumini-hub/common/config"
	"lumini-hub/common/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuthHandler gerencia as requisições de autenticação
type AuthHandler struct {
	authService *service.AuthService
	cfg         *config.Config
}

// NewAuthHandler cria um novo handler de autenticação
func NewAuthHandler(db *gorm.DB, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		authService: service.NewAuthService(db, cfg),
		cfg:         cfg,
	}
}

// LoginRequest representa os dados de login
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login autentica um usuário
// @Summary      Realiza login de usuário
// @Description  Autentica o usuário no ERP Lumini Hub, definindo cookies seguros HTTP-Only
// @Tags         Autenticação
// @Accept       json
// @Produce      json
// @Param        credentials  body      LoginRequest  true  "Credenciais de acesso"
// @Success      200          {object}  utils.Response{data=domain.LoginSuccessResponse}
// @Failure      400          {object}  utils.Response
// @Failure      401          {object}  utils.Response
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Dados de login inválidos", err.Error())
		return
	}

	response, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Falha na autenticação", err.Error())
		return
	}

	accessTokenDuration := time.Duration(h.cfg.JWT.AccessTokenExp.Minutes()) * time.Minute
	refreshTokenDuration := time.Duration(h.cfg.JWT.RefreshTokenExp.Minutes()) * time.Minute

	c.SetCookie(
		"access_token",
		response.AccessToken,
		int(accessTokenDuration.Seconds()),
		"/",
		"",
		h.cfg.App.Env == "production",
		true,
	)

	c.SetCookie(
		"refresh_token",
		response.RefreshToken,
		int(refreshTokenDuration.Seconds()),
		"/",
		"",
		h.cfg.App.Env == "production",
		true,
	)

	successResponse := domain.LoginSuccessResponse{
		User: response.User,
	}

	utils.SuccessResponse(c, http.StatusOK, "Login realizado com sucesso", successResponse, nil)
}

// RefreshToken renova o token de acesso
// @Summary      Renova token de autenticação
// @Description  Gera novos tokens de acesso e refresh a partir do cookie refresh_token existente
// @Tags         Autenticação
// @Accept       json
// @Produce      json
// @Success      200  {object}  utils.Response{data=domain.RefreshTokenSuccessResponse}
// @Failure      401  {object}  utils.Response
// @Router       /auth/refresh-token [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Falha ao renovar token", "Token de refresh ausente no cookie")
		return
	}

	response, err := h.authService.RefreshToken(refreshToken)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Falha ao renovar token", err.Error())
		return
	}

	accessTokenDuration := time.Duration(h.cfg.JWT.AccessTokenExp.Minutes()) * time.Minute
	refreshTokenDuration := time.Duration(h.cfg.JWT.RefreshTokenExp.Minutes()) * time.Minute

	c.SetCookie(
		"access_token",
		response.AccessToken,
		int(accessTokenDuration.Seconds()),
		"/", "",
		h.cfg.App.Env == "production",
		true,
	)

	c.SetCookie(
		"refresh_token",
		response.RefreshToken,
		int(refreshTokenDuration.Seconds()),
		"/", "",
		h.cfg.App.Env == "production",
		true,
	)

	successResponse := domain.RefreshTokenSuccessResponse{
		User: response.User,
	}

	utils.SuccessResponse(c, http.StatusOK, "Token renovado com sucesso", successResponse, nil)
}

// Logout realiza o logout do usuário
// @Summary      Efetua logout do sistema
// @Description  Limpa os cookies de sessão ativos no navegador
// @Tags         Autenticação
// @Accept       json
// @Produce      json
// @Success      200  {object}  utils.Response
// @Router       /auth/logout [post]
// @Security     ApiKeyAuth
func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/", "", h.cfg.App.Env == "production", true)
	c.SetCookie("refresh_token", "", -1, "/", "", h.cfg.App.Env == "production", true)

	utils.SuccessResponse(c, http.StatusOK, "Logout realizado com sucesso", nil, nil)
}

// GetMe retorna informações do usuário logado
// @Summary      Obtém dados do usuário atual
// @Description  Retorna o perfil detalhado do usuário ativo na sessão
// @Tags         Autenticação
// @Accept       json
// @Produce      json
// @Success      200  {object}  utils.Response{data=domain.LoginSuccessResponse}
// @Failure      401  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /auth/me [get]
// @Security     ApiKeyAuth
func (h *AuthHandler) GetMe(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Não autorizado", "Usuário não autenticado")
		return
	}

	user, err := h.authService.GetUserByID(userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao buscar usuário", err.Error())
		return
	}

	userDetail := domain.ApiUserDetailFromModel(*user)
	if menuItems, err := h.authService.GetMenuItemsForUser(*user); err == nil {
		userDetail.MenuItems = menuItems
	}
	if err := h.authService.ResolveActiveCompany(&userDetail); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao resolver empresa ativa", err.Error())
		return
	}

	userResponse := domain.LoginSuccessResponse{
		User: userDetail,
	}

	utils.SuccessResponse(c, http.StatusOK, "Usuário encontrado", userResponse, nil)
}

// SetActiveCompanyRequest representa o corpo de PUT /auth/active-company
type SetActiveCompanyRequest struct {
	CompanyID uint `json:"company_id" binding:"required"`
}

// SetActiveCompany define a empresa que o usuário logado está "operando
// como" — mecanismo de empresa ativa (plano_empresa.md), obrigatório pra
// quem enxerga mais de uma Company (usuário master, ou com
// companies.hierarchy.view) antes de usar o restante do sistema.
// @Summary      Define a empresa ativa do usuário logado
// @Description  Grava qual Company o usuário está operando no momento — só aceita uma que ele realmente enxergue (nunca confia cegamente no ID vindo do cliente)
// @Tags         Autenticação
// @Accept       json
// @Produce      json
// @Param        body  body      SetActiveCompanyRequest  true  "ID da Empresa"
// @Success      200   {object}  utils.Response{data=domain.LoginSuccessResponse}
// @Failure      400   {object}  utils.Response
// @Failure      401   {object}  utils.Response
// @Failure      403   {object}  utils.Response
// @Failure      404   {object}  utils.Response
// @Router       /auth/active-company [put]
// @Security     ApiKeyAuth
func (h *AuthHandler) SetActiveCompany(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Não autorizado", "Usuário não autenticado")
		return
	}

	var req SetActiveCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		return
	}

	user, err := h.authService.SetActiveCompany(userID.(uint), req.CompanyID)
	if err != nil {
		switch err {
		case utils.ErrNotFound:
			utils.ErrorResponse(c, http.StatusNotFound, "Empresa não encontrada", err.Error())
		case utils.ErrForbidden:
			utils.ErrorResponse(c, http.StatusForbidden, "Você não tem acesso a essa empresa", err.Error())
		default:
			utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao definir empresa ativa", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Empresa ativa definida com sucesso", domain.LoginSuccessResponse{User: *user}, nil)
}
