package handlers

import (
	"errors"
	"net/http"

	"lumini-hub/api.auth/internal/domain"
	"lumini-hub/api.auth/internal/service"
	"lumini-hub/common/config"
	"lumini-hub/common/internalmail"
	"lumini-hub/common/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PasswordResetHandler gerencia a recuperação de senha por e-mail (rotas públicas)
type PasswordResetHandler struct {
	resetService *service.PasswordResetService
}

// NewPasswordResetHandler cria um novo handler de recuperação de senha
func NewPasswordResetHandler(db *gorm.DB, cfg *config.Config) *PasswordResetHandler {
	return &PasswordResetHandler{
		resetService: service.NewPasswordResetService(db, internalmail.NewClient("api.auth", cfg)),
	}
}

// ForgotPassword pede o envio do link de redefinição de senha
// @Summary      Esqueci a senha
// @Description  Envia, para o e-mail informado, um link de redefinição de senha válido por 30 minutos. A resposta é sempre a mesma, exista ou não uma conta com esse e-mail.
// @Tags         Autenticação
// @Accept       json
// @Produce      json
// @Param        body  body      domain.ForgotPasswordRequest  true  "E-mail da conta"
// @Success      200   {object}  utils.Response
// @Failure      400   {object}  utils.Response
// @Failure      429   {object}  utils.Response
// @Router       /auth/forgot-password [post]
func (h *PasswordResetHandler) ForgotPassword(c *gin.Context) {
	var req domain.ForgotPasswordRequest
	if err := utils.BindJsonOrSendErrorRes(c, &req); err != nil {
		return
	}

	if err := h.resetService.ForgotPassword(req.Email, c.ClientIP()); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao processar a recuperação de senha", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Se o e-mail estiver cadastrado, enviamos um link para redefinir a senha.", nil, nil)
}

// ResetPassword define a nova senha com o token recebido por e-mail
// @Summary      Redefine a senha
// @Description  Troca a senha da conta dona do token do link. O token vale uma vez só e por 30 minutos.
// @Tags         Autenticação
// @Accept       json
// @Produce      json
// @Param        body  body      domain.ResetPasswordRequest  true  "Token do link e nova senha"
// @Success      200   {object}  utils.Response
// @Failure      400   {object}  utils.Response
// @Failure      429   {object}  utils.Response
// @Router       /auth/reset-password [post]
func (h *PasswordResetHandler) ResetPassword(c *gin.Context) {
	var req domain.ResetPasswordRequest
	if err := utils.BindJsonOrSendErrorRes(c, &req); err != nil {
		return
	}

	if err := h.resetService.ResetPassword(req.Token, req.NewPassword); err != nil {
		if errors.Is(err, service.ErrInvalidResetToken) {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao redefinir a senha", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Senha redefinida com sucesso. Entre com a nova senha.", nil, nil)
}
