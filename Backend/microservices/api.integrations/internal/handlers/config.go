package handlers

import (
	"net/http"

	"lumini-hub/api.integrations/internal/domain"
	"lumini-hub/api.integrations/internal/service"
	"lumini-hub/common/utils"

	"github.com/gin-gonic/gin"
)

// ConfigHandler gerencia as requisições relacionadas às configurações do api.integrations
type ConfigHandler struct {
	configService *service.ConfigService
}

// NewConfigHandler cria um novo handler de configurações
func NewConfigHandler(configService *service.ConfigService) *ConfigHandler {
	return &ConfigHandler{configService: configService}
}

// GetSettings retorna as configurações atuais do api.integrations
// @Summary      Obtém as configurações do api.integrations
// @Description  Retorna as chaves da Loja Integrada e os códigos do ERP legado configurados
// @Tags         Integrações - Configurações
// @Produce      json
// @Success      200  {object}  utils.Response{data=domain.ApiIntegrationSettings}
// @Failure      401  {object}  utils.Response
// @Router       /integrations/settings [get]
// @Security     ApiKeyAuth
func (h *ConfigHandler) GetSettings(c *gin.Context) {
	settings := h.configService.GetSettings()
	utils.SuccessResponse(c, http.StatusOK, "Configurações obtidas com sucesso", settings, nil)
}

// UpdateSettings atualiza em lote as configurações do api.integrations
// @Summary      Atualiza as configurações do api.integrations
// @Description  Faz o upsert em lote das chaves informadas (campos omitidos não são alterados)
// @Tags         Integrações - Configurações
// @Accept       json
// @Produce      json
// @Param        settings  body      domain.UpdateIntegrationSettingsRequest  true  "Configurações a atualizar"
// @Success      200       {object}  utils.Response{data=domain.ApiIntegrationSettings}
// @Failure      400       {object}  utils.Response
// @Failure      401       {object}  utils.Response
// @Failure      500       {object}  utils.Response
// @Router       /integrations/settings [put]
// @Security     ApiKeyAuth
func (h *ConfigHandler) UpdateSettings(c *gin.Context) {
	var req domain.UpdateIntegrationSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		return
	}

	if err := h.configService.UpdateSettings(req); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao atualizar configurações", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Configurações atualizadas com sucesso", h.configService.GetSettings(), nil)
}
