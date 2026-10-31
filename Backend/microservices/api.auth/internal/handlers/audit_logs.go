package handlers

import (
	"net/http"

	"lumini-hub/api.auth/internal/domain"
	"lumini-hub/api.auth/internal/repository"
	"lumini-hub/api.auth/internal/service"
	"lumini-hub/common/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuditLogHandler gerencia a consulta da auditoria
type AuditLogHandler struct {
	auditService *service.AuditLogService
}

// NewAuditLogHandler cria um novo handler de auditoria
func NewAuditLogHandler(db *gorm.DB) *AuditLogHandler {
	return &AuditLogHandler{
		auditService: service.NewAuditLogService(repository.NewAuditLogRepository(db)),
	}
}

// PostAuditLogFilter busca registros de auditoria aplicando os filtros recebidos no corpo
//
//	@Summary		Busca a auditoria
//	@Description	Lista paginada das escritas feitas no sistema (quem, o quê, quando), com filtros por texto, entidade, método, usuário e período.
//	@Tags			Auditoria
//	@Accept			json
//	@Produce		json
//	@Param			filter	body		domain.AuditLogFilterRequest	true	"Filtros e paginação"
//	@Success		200		{object}	domain.ApiAuditLogListPaginated
//	@Router			/audit-logs/filter [post]
func (h *AuditLogHandler) PostAuditLogFilter(c *gin.Context) {
	var req domain.AuditLogFilterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Parâmetros de filtro inválidos", err)
		return
	}

	logs, err := h.auditService.GetAuditLogs(req, utils.RoleFromGinContext(c))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao buscar a auditoria", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Registros de auditoria encontrados", logs, nil)
}

// GetAuditEntityTypes lista os tipos de entidade presentes na auditoria
//
//	@Summary		Tipos de entidade da auditoria
//	@Description	Valores distintos de entidade (users, companies...) para o filtro da tela.
//	@Tags			Auditoria
//	@Produce		json
//	@Success		200	{array}	string
//	@Router			/audit-logs/entity-types [get]
func (h *AuditLogHandler) GetAuditEntityTypes(c *gin.Context) {
	types, err := h.auditService.GetEntityTypes()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao buscar os tipos de entidade", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Tipos de entidade encontrados", types, nil)
}
