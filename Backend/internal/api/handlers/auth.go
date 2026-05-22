package handlers

import (
	"net/http"
	"time" // Adicionar import para time

	"simple-erp-service/config"
	"simple-erp-service/internal/data-structure/dto"
	"simple-erp-service/internal/service"
	"simple-erp-service/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuthHandler gerencia as requisições de autenticação
type AuthHandler struct {
	authService *service.AuthService
	cfg         *config.Config // Adicionar configuração aqui para acessar as durações dos tokens
}

// NewAuthHandler cria um novo handler de autenticação
func NewAuthHandler(db *gorm.DB, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		authService: service.NewAuthService(db, cfg),
		cfg:         cfg, // Passar a configuração para o handler
	}
}

// LoginRequest representa os dados de login
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RefreshTokenRequest representa os dados para renovar o token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
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

	// --- PASSO CHAVE 1: DEFINIR COOKIES HTTP-ONLY ---
	// Calcular a duração do Access Token para o cookie
	accessTokenDuration := time.Duration(h.cfg.JWT.AccessTokenExp.Minutes()) * time.Minute
	// Calcular a duração do Refresh Token para o cookie
	refreshTokenDuration := time.Duration(h.cfg.JWT.RefreshTokenExp.Minutes()) * time.Minute

	// Definir o Access Token como cookie HTTP-Only
	c.SetCookie(
		"access_token",                     // Nome do cookie
		response.AccessToken,               // Valor do token
		int(accessTokenDuration.Seconds()), // MaxAge em segundos
		"/",                                // Path: disponível em toda a aplicação
		"",                                 // Domain: vazio para o domínio da requisição
		h.cfg.App.Env == "production",      // Secure: true apenas em produção (HTTPS)
		true,                               // HttpOnly: Essencial! NÃO acessível por JavaScript
	)

	// Definir o Refresh Token como cookie HTTP-Only
	c.SetCookie(
		"refresh_token",                     // Nome do cookie
		response.RefreshToken,               // Valor do token
		int(refreshTokenDuration.Seconds()), // MaxAge em segundos
		"/",                                 // Path
		"",                                  // Domain
		h.cfg.App.Env == "production",       // Secure
		true,                                // HttpOnly
	)

	// --- PASSO CHAVE 2: ENVIAR APENAS DADOS DO USUÁRIO NO JSON ---
	// Criar uma resposta que contém apenas os dados do usuário para o frontend
	// Isso evita que o frontend tenha acesso direto aos tokens (que agora estão nos cookies)
	successResponse := dto.LoginSuccessResponse{
		User: response.User,
	}

	// Retornar a resposta de sucesso com os dados do usuário (e os cookies definidos)
	utils.SuccessResponse(c, http.StatusOK, "Login realizado com sucesso", successResponse, nil)
}

// RefreshToken renova o token de acesso
// PRECISAMOS REFACTORAR ESTE TAMBÉM
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// --- Lógica de RefreshToken para ler do cookie e definir novos cookies ---
	// Em vez de ler do corpo da requisição, tentaremos ler o refresh_token do cookie
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Falha ao renovar token", "Token de refresh ausente no cookie")
		return
	}

	// Usar o token lido do cookie para renovar
	response, err := h.authService.RefreshToken(refreshToken)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Falha ao renovar token", err.Error())
		return
	}

	// --- Definir NOVOS cookies HTTP-Only para Access e Refresh Token ---
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

	// Enviar apenas os dados do usuário no JSON
	successResponse := dto.RefreshTokenSuccessResponse{
		User: response.User,
	}

	utils.SuccessResponse(c, http.StatusOK, "Token renovado com sucesso", successResponse, nil)
}

// Logout realiza o logout do usuário
func (h *AuthHandler) Logout(c *gin.Context) {
	// --- PASSO CHAVE 3: LIMPAR COOKIES NO LOGOUT ---
	// Limpar os cookies de token
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

	userResponse := dto.LoginSuccessResponse{
		User: dto.ApiUserDetailFromModel(*user),
	}

	utils.SuccessResponse(c, http.StatusOK, "Usuário encontrado", userResponse, nil)
}
