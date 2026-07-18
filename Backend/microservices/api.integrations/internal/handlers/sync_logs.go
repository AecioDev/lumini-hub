package handlers

import (
	"net/http"

	"lumini-hub/api.integrations/internal/domain"
	"lumini-hub/api.integrations/internal/repository"
	"lumini-hub/api.integrations/internal/service"
	"lumini-hub/common/utils"

	"github.com/gin-gonic/gin"
)

// SyncLogHandler gerencia as requisições relacionadas aos logs de sincronização
type SyncLogHandler struct {
	syncLogService *service.SyncLogService
}

// NewSyncLogHandler cria um novo handler de sync logs recebendo o Unit of Work
func NewSyncLogHandler(uow repository.UnitOfWork) *SyncLogHandler {
	return &SyncLogHandler{syncLogService: service.NewSyncLogService(uow)}
}

// PostSyncLogFilter busca logs de sincronização aplicando filtros e paginação (PostFilter)
// @Summary      Filtra logs de sincronização com paginação
// @Description  Busca e filtra os logs de sincronização entre a Loja Integrada e o SQL Server legado
// @Tags         Integrações - Sync Logs
// @Accept       json
// @Produce      json
// @Param        filter  body      domain.SyncLogFilterRequest  true  "Filtros e Dados de Paginação"
// @Success      200     {object}  utils.Response{data=domain.ApiSyncLogListPaginated}
// @Failure      400     {object}  utils.Response
// @Failure      401     {object}  utils.Response
// @Failure      500     {object}  utils.Response
// @Router       /integrations/sync-logs/filter [post]
// @Security     ApiKeyAuth
func (h *SyncLogHandler) PostSyncLogFilter(c *gin.Context) {
	var req domain.SyncLogFilterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Parâmetros de filtro inválidos", err.Error())
		return
	}

	logs, err := h.syncLogService.GetSyncLogsByFilter(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao buscar logs de sincronização", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Logs de sincronização encontrados", logs, nil)
}
