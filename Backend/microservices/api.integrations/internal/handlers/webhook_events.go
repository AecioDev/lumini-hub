package handlers

import (
	"net/http"

	"lumini-hub/api.integrations/internal/domain"
	"lumini-hub/api.integrations/internal/repository"
	"lumini-hub/api.integrations/internal/service"
	"lumini-hub/common/utils"

	"github.com/gin-gonic/gin"
)

// WebhookEventHandler gerencia as requisições de consulta aos eventos de webhook recebidos
type WebhookEventHandler struct {
	webhookEventService *service.WebhookEventService
}

// NewWebhookEventHandler cria um novo handler de webhook events recebendo o Unit of Work
func NewWebhookEventHandler(uow repository.UnitOfWork) *WebhookEventHandler {
	return &WebhookEventHandler{webhookEventService: service.NewWebhookEventService(uow)}
}

// PostWebhookEventFilter busca eventos de webhook aplicando filtros e paginação (PostFilter)
// @Summary      Filtra eventos de webhook com paginação
// @Description  Busca e filtra os eventos de webhook recebidos (ex.: da Loja Integrada)
// @Tags         Integrações - Webhook Events
// @Accept       json
// @Produce      json
// @Param        filter  body      domain.WebhookEventFilterRequest  true  "Filtros e Dados de Paginação"
// @Success      200     {object}  utils.Response{data=domain.ApiWebhookEventListPaginated}
// @Failure      400     {object}  utils.Response
// @Failure      401     {object}  utils.Response
// @Failure      500     {object}  utils.Response
// @Router       /integrations/webhook-events/filter [post]
// @Security     ApiKeyAuth
func (h *WebhookEventHandler) PostWebhookEventFilter(c *gin.Context) {
	var req domain.WebhookEventFilterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Parâmetros de filtro inválidos", err.Error())
		return
	}

	events, err := h.webhookEventService.GetWebhookEventsByFilter(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao buscar eventos de webhook", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Eventos de webhook encontrados", events, nil)
}
