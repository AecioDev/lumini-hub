package handlers

import (
	"errors"
	"net/http"

	"lumini-hub/api.core/internal/domain"
	"lumini-hub/api.core/internal/repository"
	"lumini-hub/api.core/internal/service"
	"lumini-hub/api.core/internal/validator"
	"lumini-hub/common/mail"
	"lumini-hub/common/utils"
	"lumini-hub/common/utils/path"

	"github.com/gin-gonic/gin"
)

// MailConfigHandler gerencia as requisições da configuração de e-mail: a padrão do sistema e a de cada empresa
type MailConfigHandler struct {
	mailService *service.MailConfigService
}

// NewMailConfigHandler cria um novo handler de configuração de e-mail
func NewMailConfigHandler(uow repository.UnitOfWork, encryptionKey string) *MailConfigHandler {
	return &MailConfigHandler{
		mailService: service.NewMailConfigService(uow, mail.NewSMTPSender(), encryptionKey),
	}
}

// replyMailError traduz os erros do service para a resposta HTTP padrão
func replyMailError(c *gin.Context, action string, err error) {
	var sendErr *service.MailSendError
	switch {
	case errors.Is(err, service.ErrAppBaseURLNotConfigured):
		utils.ErrorResponse(c, http.StatusUnprocessableEntity, err.Error(), err.Error())
	case errors.Is(err, utils.ErrNotFound):
		utils.ErrorResponse(c, http.StatusNotFound, "Empresa não encontrada", err.Error())
	case validator.IsValidationError(err):
		utils.ValidationErrorResponse(c, "Dados inválidos", err)
	case errors.Is(err, service.ErrMailNotConfigured):
		utils.ErrorResponse(c, http.StatusUnprocessableEntity, err.Error(), err.Error())
	case errors.As(err, &sendErr):
		utils.ErrorResponse(c, http.StatusBadGateway, "O servidor de e-mail recusou o envio", sendErr.Error())
	default:
		utils.ErrorResponse(c, http.StatusInternalServerError, action, err.Error())
	}
}

// GetSystemEmail retorna o e-mail padrão do sistema
// @Summary      Busca o e-mail padrão do sistema
// @Description  Retorna o servidor SMTP padrão do sistema, usado quando a empresa e a Matriz dela não têm um próprio. A senha nunca é devolvida.
// @Tags         SystemSettings
// @Produce      json
// @Success      200  {object}  utils.Response{data=domain.ApiSystemEmailSettings}
// @Failure      401  {object}  utils.Response
// @Router       /system-settings/email [get]
// @Security     ApiKeyAuth
func (h *MailConfigHandler) GetSystemEmail(c *gin.Context) {
	settings, err := h.mailService.GetSystemEmail()
	if err != nil {
		replyMailError(c, "Erro ao buscar o e-mail do sistema", err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "E-mail do sistema encontrado", settings, nil)
}

// SaveSystemEmail grava o e-mail padrão do sistema
// @Summary      Salva o e-mail padrão do sistema
// @Description  Grava o servidor SMTP padrão do sistema e o endereço público do sistema (usado nos links dos e-mails). Senha vazia mantém a atual.
// @Tags         SystemSettings
// @Accept       json
// @Produce      json
// @Param        config  body      domain.SaveSystemEmailRequest  true  "Servidor SMTP e endereço do sistema"
// @Success      200     {object}  utils.Response{data=domain.ApiSystemEmailSettings}
// @Failure      400     {object}  utils.Response
// @Failure      401     {object}  utils.Response
// @Router       /system-settings/email [put]
// @Security     ApiKeyAuth
func (h *MailConfigHandler) SaveSystemEmail(c *gin.Context) {
	var req domain.SaveSystemEmailRequest
	if err := utils.BindJsonOrSendErrorRes(c, &req); err != nil {
		return
	}
	userID, _ := utils.GetUserIDFromContext(c)

	settings, err := h.mailService.SaveSystemEmail(req, userID)
	if err != nil {
		replyMailError(c, "Erro ao salvar o e-mail do sistema", err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "E-mail do sistema salvo com sucesso", settings, nil)
}

// SendSystemTestEmail manda um e-mail de teste pela configuração do sistema
// @Summary      Envia e-mail de teste (sistema)
// @Description  Envia um e-mail de teste ao endereço informado usando o e-mail padrão do sistema já salvo
// @Tags         SystemSettings
// @Accept       json
// @Produce      json
// @Param        body  body      domain.SendTestEmailRequest  true  "Destinatário do teste"
// @Success      200   {object}  utils.Response{data=domain.ApiTestEmailResult}
// @Failure      400   {object}  utils.Response
// @Failure      422   {object}  utils.Response
// @Failure      502   {object}  utils.Response
// @Router       /system-settings/email/test [post]
// @Security     ApiKeyAuth
func (h *MailConfigHandler) SendSystemTestEmail(c *gin.Context) {
	var req domain.SendTestEmailRequest
	if err := utils.BindJsonOrSendErrorRes(c, &req); err != nil {
		return
	}
	result, err := h.mailService.SendTest(nil, req.To)
	if err != nil {
		replyMailError(c, "Erro ao enviar o e-mail de teste", err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "E-mail de teste enviado", result, nil)
}

// GetCompanyEmail retorna a configuração de e-mail de uma empresa
// @Summary      Busca o e-mail de uma empresa
// @Description  Retorna o servidor SMTP próprio da empresa (campos padrão se ela herda) e de onde o envio sairia hoje: a própria empresa, a Matriz ou o sistema
// @Tags         CompanyEmailConfig
// @Produce      json
// @Param        companyId  path      int  true  "ID da Empresa"
// @Success      200        {object}  utils.Response{data=domain.ApiCompanyEmailConfig}
// @Failure      404        {object}  utils.Response
// @Router       /company-email-configs/by-company/{companyId} [get]
// @Security     ApiKeyAuth
func (h *MailConfigHandler) GetCompanyEmail(c *gin.Context) {
	companyID, err := path.UintFromPathParam(c, "companyId")
	if err != nil {
		utils.ValidationErrorResponse(c, "Dados inválidos", err)
		return
	}
	config, err := h.mailService.GetCompanyEmail(companyID)
	if err != nil {
		replyMailError(c, "Erro ao buscar o e-mail da empresa", err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "E-mail da empresa encontrado", config, nil)
}

// SaveCompanyEmail grava o e-mail próprio de uma empresa
// @Summary      Salva o e-mail de uma empresa
// @Description  Cria ou atualiza o servidor SMTP próprio da empresa. Senha vazia mantém a atual.
// @Tags         CompanyEmailConfig
// @Accept       json
// @Produce      json
// @Param        companyId  path      int                     true  "ID da Empresa"
// @Param        config     body      domain.SaveSMTPRequest  true  "Servidor SMTP"
// @Success      200        {object}  utils.Response{data=domain.ApiCompanyEmailConfig}
// @Failure      400        {object}  utils.Response
// @Failure      404        {object}  utils.Response
// @Router       /company-email-configs/by-company/{companyId} [put]
// @Security     ApiKeyAuth
func (h *MailConfigHandler) SaveCompanyEmail(c *gin.Context) {
	companyID, err := path.UintFromPathParam(c, "companyId")
	if err != nil {
		utils.ValidationErrorResponse(c, "Dados inválidos", err)
		return
	}
	var req domain.SaveSMTPRequest
	if err := utils.BindJsonOrSendErrorRes(c, &req); err != nil {
		return
	}
	config, err := h.mailService.SaveCompanyEmail(companyID, req)
	if err != nil {
		replyMailError(c, "Erro ao salvar o e-mail da empresa", err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "E-mail da empresa salvo com sucesso", config, nil)
}

// DeleteCompanyEmail remove o e-mail próprio de uma empresa
// @Summary      Remove o e-mail de uma empresa
// @Description  Apaga o servidor SMTP próprio: a empresa volta a herdar da Matriz ou do sistema
// @Tags         CompanyEmailConfig
// @Produce      json
// @Param        companyId  path      int  true  "ID da Empresa"
// @Success      200        {object}  utils.Response
// @Failure      404        {object}  utils.Response
// @Router       /company-email-configs/by-company/{companyId} [delete]
// @Security     ApiKeyAuth
func (h *MailConfigHandler) DeleteCompanyEmail(c *gin.Context) {
	companyID, err := path.UintFromPathParam(c, "companyId")
	if err != nil {
		utils.ValidationErrorResponse(c, "Dados inválidos", err)
		return
	}
	if err := h.mailService.DeleteCompanyEmail(companyID); err != nil {
		replyMailError(c, "Erro ao remover o e-mail da empresa", err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "E-mail da empresa removido: ela voltou a herdar", nil, nil)
}

// SendCompanyTestEmail manda um e-mail de teste pela configuração que valeria para a empresa
// @Summary      Envia e-mail de teste (empresa)
// @Description  Envia um e-mail de teste pela configuração efetiva da empresa (própria, da Matriz ou do sistema) e informa qual foi usada
// @Tags         CompanyEmailConfig
// @Accept       json
// @Produce      json
// @Param        companyId  path      int                          true  "ID da Empresa"
// @Param        body       body      domain.SendTestEmailRequest  true  "Destinatário do teste"
// @Success      200        {object}  utils.Response{data=domain.ApiTestEmailResult}
// @Failure      404        {object}  utils.Response
// @Failure      422        {object}  utils.Response
// @Failure      502        {object}  utils.Response
// @Router       /company-email-configs/by-company/{companyId}/test [post]
// @Security     ApiKeyAuth
func (h *MailConfigHandler) SendCompanyTestEmail(c *gin.Context) {
	companyID, err := path.UintFromPathParam(c, "companyId")
	if err != nil {
		utils.ValidationErrorResponse(c, "Dados inválidos", err)
		return
	}
	var req domain.SendTestEmailRequest
	if err := utils.BindJsonOrSendErrorRes(c, &req); err != nil {
		return
	}
	result, err := h.mailService.SendTest(&companyID, req.To)
	if err != nil {
		replyMailError(c, "Erro ao enviar o e-mail de teste", err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "E-mail de teste enviado", result, nil)
}

// SendPasswordResetEmail envia o e-mail de redefinição de senha (rota interna, chamada pelo api.auth)
// @Summary      Envia o e-mail de redefinição de senha
// @Description  Rota interna, só com token de serviço (não passa pelo gateway): monta o link, renderiza o template e envia pela configuração da empresa do usuário
// @Tags         Interno
// @Accept       json
// @Produce      json
// @Param        body  body      domain.SendPasswordResetRequest  true  "Destinatário, nome, token e empresa"
// @Success      200   {object}  utils.Response
// @Failure      401   {object}  utils.Response
// @Failure      422   {object}  utils.Response
// @Failure      502   {object}  utils.Response
// @Router       /internal/mail/password-reset [post]
func (h *MailConfigHandler) SendPasswordResetEmail(c *gin.Context) {
	var req domain.SendPasswordResetRequest
	if err := utils.BindJsonOrSendErrorRes(c, &req); err != nil {
		return
	}
	if err := h.mailService.SendPasswordReset(req); err != nil {
		replyMailError(c, "Erro ao enviar o e-mail de redefinição de senha", err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "E-mail enviado", nil, nil)
}
