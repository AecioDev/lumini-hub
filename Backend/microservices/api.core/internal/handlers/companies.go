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

// CompanyHandler gerencia as requisições relacionadas a empresas
type CompanyHandler struct {
	companyService *service.CompanyService
}

// NewCompanyHandler cria um novo handler de empresas recebendo o Unit of Work
func NewCompanyHandler(uow repository.UnitOfWork) *CompanyHandler {
	return &CompanyHandler{
		companyService: service.NewCompanyService(uow),
	}
}

// GetCompanies retorna todas as empresas cadastradas
// @Summary      Lista empresas
// @Description  Retorna todas as empresas cadastradas (Matriz e vinculadas)
// @Tags         Companies
// @Accept       json
// @Produce      json
// @Success      200  {object}  utils.Response{data=domain.ApiCompanyList}
// @Failure      401  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /companies [get]
// @Security     ApiKeyAuth
func (h *CompanyHandler) GetCompanies(c *gin.Context) {
	companies, err := h.companyService.GetCompanies()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao buscar empresas", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Empresas encontradas", companies, nil)
}

// GetCompany retorna uma empresa específica pelo ID
// @Summary      Busca empresa por ID
// @Description  Retorna os detalhes de uma única empresa cadastrada
// @Tags         Companies
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID da Empresa"
// @Success      200  {object}  utils.Response{data=domain.ApiCompanyDetail}
// @Failure      401  {object}  utils.Response
// @Failure      404  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /companies/{id} [get]
// @Security     ApiKeyAuth
func (h *CompanyHandler) GetCompany(c *gin.Context) {
	id, err := path.IdFromPathParamOrSendError(c)
	if err != nil {
		return
	}

	company, err := h.companyService.GetCompanyByID(id)
	if err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Empresa não encontrada", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao buscar empresa", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Empresa encontrada", company, nil)
}

// CreateCompany cria uma nova empresa
// @Summary      Cria nova empresa
// @Description  Cadastra uma empresa (Matriz ou vinculada a outra) no sistema
// @Tags         Companies
// @Accept       json
// @Produce      json
// @Param        company  body      domain.CreateCompanyRequest  true  "Dados da Empresa"
// @Success      201      {object}  utils.Response{data=domain.ApiCompany}
// @Failure      400      {object}  utils.Response
// @Failure      401      {object}  utils.Response
// @Failure      500      {object}  utils.Response
// @Router       /companies [post]
// @Security     ApiKeyAuth
func (h *CompanyHandler) CreateCompany(c *gin.Context) {
	var req domain.CreateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		return
	}

	userID, exists := utils.GetUserIDFromContext(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Usuário não autenticado", "")
		return
	}

	company, err := h.companyService.CreateCompany(req, userID)
	if err != nil {
		if validator.IsValidationError(err) {
			utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusBadRequest, "Erro ao criar empresa", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Empresa criada com sucesso", company, nil)
}

// UpdateCompany atualiza uma empresa existente
// @Summary      Atualiza empresa
// @Description  Altera dados cadastrais de uma empresa pelo ID
// @Tags         Companies
// @Accept       json
// @Produce      json
// @Param        id       path      int                          true  "ID da Empresa"
// @Param        company  body      domain.UpdateCompanyRequest  true  "Novos dados da empresa"
// @Success      200      {object}  utils.Response{data=domain.ApiCompany}
// @Failure      400      {object}  utils.Response
// @Failure      401      {object}  utils.Response
// @Failure      404      {object}  utils.Response
// @Failure      500      {object}  utils.Response
// @Router       /companies/{id} [put]
// @Security     ApiKeyAuth
func (h *CompanyHandler) UpdateCompany(c *gin.Context) {
	id, err := path.IdFromPathParamOrSendError(c)
	if err != nil {
		return
	}

	var req domain.UpdateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		return
	}

	userID, exists := utils.GetUserIDFromContext(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Usuário não autenticado", "")
		return
	}

	company, err := h.companyService.UpdateCompany(id, req, userID)
	if err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Empresa não encontrada", err.Error())
		} else if validator.IsValidationError(err) {
			utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusBadRequest, "Erro ao atualizar empresa", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Empresa atualizada com sucesso", company, nil)
}

// DeleteCompany exclui uma empresa pelo ID
// @Summary      Exclui empresa
// @Description  Realiza soft delete de uma empresa (bloqueado se houver outras empresas vinculadas a ela)
// @Tags         Companies
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID da Empresa"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      401  {object}  utils.Response
// @Failure      404  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /companies/{id} [delete]
// @Security     ApiKeyAuth
func (h *CompanyHandler) DeleteCompany(c *gin.Context) {
	id, err := path.IdFromPathParamOrSendError(c)
	if err != nil {
		return
	}

	if err := h.companyService.DeleteCompany(id); err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Empresa não encontrada", err.Error())
		} else if validator.IsValidationError(err) {
			utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusBadRequest, "Erro ao excluir empresa", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Empresa excluída com sucesso", nil, nil)
}
