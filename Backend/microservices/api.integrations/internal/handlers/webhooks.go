package handlers

import (
	"io"
	"net/http"

	"lumini-hub/api.integrations/internal/service"
	"lumini-hub/common/utils"

	"github.com/gin-gonic/gin"
)

// WebhookHandler gerencia o recebimento de webhooks externos (ex.: Loja Integrada)
type WebhookHandler struct {
	webhookService *service.WebhookService
}

// NewWebhookHandler cria um novo handler de webhooks
func NewWebhookHandler(webhookService *service.WebhookService) *WebhookHandler {
	return &WebhookHandler{webhookService: webhookService}
}

// ReceiveLojaIntegradaWebhook recebe notificações da Loja Integrada (ex.: pedido criado)
// @Summary      Recebe webhook da Loja Integrada
// @Description  Endpoint público (sem AuthMiddleware), validado por um segredo compartilhado no header X-Webhook-Secret (configurado em /integrations/settings). Persiste o payload cru; o processamento em pedido de venda no SQL Server fica para a Fase 2.
// @Tags         Integrações - Webhooks
// @Accept       json
// @Produce      json
// @Param        X-Webhook-Secret  header    string  false  "Segredo compartilhado configurado em /integrations/settings"
// @Success      200               {object}  utils.Response
// @Failure      400               {object}  utils.Response
// @Failure      401               {object}  utils.Response
// @Router       /integrations/webhooks/loja-integrada [post]
func (h *WebhookHandler) ReceiveLojaIntegradaWebhook(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Erro ao ler corpo da requisição", err.Error())
		return
	}

	secret := c.GetHeader("X-Webhook-Secret")

	event, err := h.webhookService.ReceiveLojaIntegradaWebhook(c.Request.Header, body, secret)
	if err != nil {
		if err == service.ErrInvalidWebhookSecret {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Segredo do webhook inválido", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao registrar webhook", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Webhook recebido com sucesso", event, nil)
}
