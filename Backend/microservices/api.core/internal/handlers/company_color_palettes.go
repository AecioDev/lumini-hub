package handlers

import (
	"net/http"

	"lumini-hub/api.core/internal/domain"
	"lumini-hub/api.core/internal/repository"
	"lumini-hub/api.core/internal/service"
	"lumini-hub/api.core/internal/validator"
	"lumini-hub/common/utils"
	"lumini-hub/common/utils/path"

	"github.com/gin-gonic/gin"
)

// CompanyColorPaletteHandler gerencia as requisições relacionadas às paletas de cores personalizadas de empresas
type CompanyColorPaletteHandler struct {
	paletteService *service.CompanyColorPaletteService
}

// NewCompanyColorPaletteHandler cria um novo handler de paletas de cores personalizadas
func NewCompanyColorPaletteHandler(uow repository.UnitOfWork) *CompanyColorPaletteHandler {
	return &CompanyColorPaletteHandler{
		paletteService: service.NewCompanyColorPaletteService(uow),
	}
}

// GetCompanyColorPalettesByCompany lista as paletas personalizadas de uma empresa
// @Summary      Lista paletas de cores personalizadas por empresa
// @Description  Retorna as paletas de cores salvas pelo usuário pra uma empresa, mais recentes primeiro
// @Tags         CompanyColorPalette
// @Accept       json
// @Produce      json
// @Param        companyId  path      int  true  "ID da Empresa"
// @Success      200        {object}  utils.Response{data=[]domain.ApiCompanyColorPalette}
// @Failure      401        {object}  utils.Response
// @Failure      500        {object}  utils.Response
// @Router       /company-color-palettes/by-company/{companyId} [get]
// @Security     ApiKeyAuth
func (h *CompanyColorPaletteHandler) GetCompanyColorPalettesByCompany(c *gin.Context) {
	companyID, err := path.UintFromPathParam(c, "companyId")
	if err != nil {
		utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		return
	}

	palettes, err := h.paletteService.GetByCompanyID(companyID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao buscar paletas personalizadas", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Paletas personalizadas encontradas", palettes, nil)
}

// CreateCompanyColorPalette salva uma nova paleta de cores personalizada
// @Summary      Salva paleta de cores personalizada
// @Description  Cria uma nova paleta de cores personalizada (nome + 3 cores) pra uma empresa
// @Tags         CompanyColorPalette
// @Accept       json
// @Produce      json
// @Param        palette  body      domain.CreateCompanyColorPaletteRequest  true  "Dados da Paleta"
// @Success      201      {object}  utils.Response{data=domain.ApiCompanyColorPalette}
// @Failure      400      {object}  utils.Response
// @Failure      401      {object}  utils.Response
// @Failure      500      {object}  utils.Response
// @Router       /company-color-palettes [post]
// @Security     ApiKeyAuth
func (h *CompanyColorPaletteHandler) CreateCompanyColorPalette(c *gin.Context) {
	var req domain.CreateCompanyColorPaletteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		return
	}

	palette, err := h.paletteService.CreateCompanyColorPalette(req)
	if err != nil {
		if validator.IsValidationError(err) {
			utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusBadRequest, "Erro ao salvar paleta personalizada", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Paleta personalizada salva com sucesso", palette, nil)
}

// DeleteCompanyColorPalette exclui uma paleta de cores personalizada
// @Summary      Exclui paleta de cores personalizada
// @Description  Remove uma paleta de cores personalizada pelo ID
// @Tags         CompanyColorPalette
// @Accept       json
// @Produce      json
// @Param        id   path  int  true  "ID da Paleta"
// @Success      200  {object}  utils.Response
// @Failure      401  {object}  utils.Response
// @Failure      404  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /company-color-palettes/{id} [delete]
// @Security     ApiKeyAuth
func (h *CompanyColorPaletteHandler) DeleteCompanyColorPalette(c *gin.Context) {
	id, err := path.IdFromPathParamOrSendError(c)
	if err != nil {
		return
	}

	if err := h.paletteService.DeleteCompanyColorPalette(id); err != nil {
		if validator.IsValidationError(err) {
			utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao excluir paleta personalizada", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Paleta personalizada excluída com sucesso", nil, nil)
}
