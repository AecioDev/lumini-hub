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
func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/", "", h.cfg.App.Env == "production", true)
	c.SetCookie("refresh_token", "", -1, "/", "", h.cfg.App.Env == "production", true)

	utils.SuccessResponse(c, http.StatusOK, "Logout realizado com sucesso", nil, nil)
}

// GetMe retorna informações do usuário logado
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

	userResponse := domain.LoginSuccessResponse{
		User: domain.ApiUserDetailFromModel(*user),
	}

	utils.SuccessResponse(c, http.StatusOK, "Usuário encontrado", userResponse, nil)
}
