package handlers

import (
	"net/http"

	legacyrepo "lumini-hub/api.integrations/internal/legacy/repository"
	"lumini-hub/common/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// LegacyLookupHandler expõe consultas de apoio ao ERP legado (SQL Server), usadas para
// alimentar os selects de local/empresa na tela de configuração. A conexão pode ser nil se o
// SQL Server estiver indisponível na subida do serviço.
type LegacyLookupHandler struct {
	sqlServerDB *gorm.DB
}

// NewLegacyLookupHandler cria um novo handler de consultas ao legado
func NewLegacyLookupHandler(sqlServerDB *gorm.DB) *LegacyLookupHandler {
	return &LegacyLookupHandler{sqlServerDB: sqlServerDB}
}

// GetLocations lista os locais de armazenamento cadastrados no ERP legado
// @Summary      Lista os locais de armazenamento do ERP legado
// @Description  Usado para alimentar os selects de codlocarm_oficial/codlocarm_reserva na tela de configuração
// @Tags         Integrações - Legado
// @Produce      json
// @Success      200  {object}  utils.Response
// @Failure      401  {object}  utils.Response
// @Failure      503  {object}  utils.Response
// @Router       /integrations/legacy/locations [get]
// @Security     ApiKeyAuth
func (h *LegacyLookupHandler) GetLocations(c *gin.Context) {
	if h.sqlServerDB == nil {
		utils.ErrorResponse(c, http.StatusServiceUnavailable, "SQL Server indisponível", "conexão com o ERP legado não configurada")
		return
	}

	repo := legacyrepo.NewLocArmRepository(h.sqlServerDB)
	locations, err := repo.FindAll()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao buscar locais de armazenamento", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Locais de armazenamento encontrados", locations, nil)
}

// GetCompanies lista as empresas cadastradas no ERP legado
// @Summary      Lista as empresas do ERP legado
// @Description  Usado para alimentar o select de codemp na tela de configuração
// @Tags         Integrações - Legado
// @Produce      json
// @Success      200  {object}  utils.Response
// @Failure      401  {object}  utils.Response
// @Failure      503  {object}  utils.Response
// @Router       /integrations/legacy/companies [get]
// @Security     ApiKeyAuth
func (h *LegacyLookupHandler) GetCompanies(c *gin.Context) {
	if h.sqlServerDB == nil {
		utils.ErrorResponse(c, http.StatusServiceUnavailable, "SQL Server indisponível", "conexão com o ERP legado não configurada")
		return
	}

	repo := legacyrepo.NewCencusRepository(h.sqlServerDB)
	companies, err := repo.FindAll()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao buscar empresas", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Empresas encontradas", companies, nil)
}
