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

// SupplierHandler gerencia as requisições relacionadas a fornecedores
type SupplierHandler struct {
	supplierService *service.SupplierService
}

// NewSupplierHandler cria um novo handler de fornecedores recebendo o Unit of Work
func NewSupplierHandler(uow repository.UnitOfWork) *SupplierHandler {
	return &SupplierHandler{
		supplierService: service.NewSupplierService(uow),
	}
}

// GetSuppliers retorna uma lista paginada de fornecedores (GET legada)
// @Summary      Lista fornecedores com paginação (GET)
// @Description  Retorna uma lista de fornecedores paginada a partir de query strings
// @Tags         Fornecedores
// @Accept       json
// @Produce      json
// @Param        page   query      int  false  "Página ativa"
// @Param        limit  query      int  false  "Limite por página"
// @Param        sort   query      string  false  "Ordenação (ex: first_name asc)"
// @Success      200    {object}  utils.Response{data=domain.ApiSupplierListPaginated}
// @Failure      401    {object}  utils.Response
// @Failure      500    {object}  utils.Response
// @Router       /suppliers [get]
// @Security     ApiKeyAuth
func (h *SupplierHandler) GetSuppliers(c *gin.Context) {
	pagination := utils.GetPaginationParams(c)

	suppliers, err := h.supplierService.GetSuppliers(&pagination)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao buscar fornecedores", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Fornecedores encontrados", suppliers, nil)
}

// PostSupplierFilter busca fornecedores aplicando filtros complexos recebidos no Body (PostFilter)
// @Summary      Filtra fornecedores com paginação (PostFilter)
// @Description  Busca e filtra fornecedores dinamicamente enviando as opções no JSON do Body
// @Tags         Fornecedores
// @Accept       json
// @Produce      json
// @Param        filter  body      domain.SupplierFilterRequest  true  "Filtros e Dados de Paginação"
// @Success      200     {object}  utils.Response{data=domain.ApiSupplierListPaginated}
// @Failure      400     {object}  utils.Response
// @Failure      401     {object}  utils.Response
// @Failure      500     {object}  utils.Response
// @Router       /suppliers/filter [post]
// @Security     ApiKeyAuth
func (h *SupplierHandler) PostSupplierFilter(c *gin.Context) {
	var req domain.SupplierFilterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Parâmetros de filtro inválidos", err.Error())
		return
	}

	suppliers, err := h.supplierService.GetSuppliersByFilter(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao buscar fornecedores com filtros", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Fornecedores encontrados", suppliers, nil)
}

// GetSupplier retorna um fornecedor específico pelo ID
// @Summary      Busca fornecedor por ID
// @Description  Retorna os detalhes de um único fornecedor cadastrado
// @Tags         Fornecedores
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID do Fornecedor"
// @Success      200  {object}  utils.Response{data=domain.ApiSupplierDetail}
// @Failure      401  {object}  utils.Response
// @Failure      404  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /suppliers/{id} [get]
// @Security     ApiKeyAuth
func (h *SupplierHandler) GetSupplier(c *gin.Context) {
	id, err := path.IdFromPathParamOrSendError(c)
	if err != nil {
		return
	}

	supplier, err := h.supplierService.GetSupplierByID(id)
	if err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Fornecedor não encontrado", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao buscar fornecedor", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Fornecedor encontrado", supplier, nil)
}

// CreateSupplier cria um novo fornecedor
// @Summary      Cria novo fornecedor
// @Description  Cadastra um fornecedor no sistema ERP
// @Tags         Fornecedores
// @Accept       json
// @Produce      json
// @Param        supplier  body      domain.CreateSupplierRequest  true  "Dados do Fornecedor"
// @Success      201       {object}  utils.Response{data=domain.ApiSupplier}
// @Failure      400       {object}  utils.Response
// @Failure      401       {object}  utils.Response
// @Failure      500       {object}  utils.Response
// @Router       /suppliers [post]
// @Security     ApiKeyAuth
func (h *SupplierHandler) CreateSupplier(c *gin.Context) {
	var req domain.CreateSupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		return
	}

	userID, _ := utils.GetUserIDFromContext(c)

	supplier, err := h.supplierService.CreateSupplier(req, userID)
	if err != nil {
		if validator.IsValidationError(err) {
			utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusBadRequest, "Erro ao criar fornecedor", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Fornecedor criado com sucesso", supplier, nil)
}

// UpdateSupplier atualiza um fornecedor existente
// @Summary      Atualiza fornecedor
// @Description  Altera dados cadastrais de um fornecedor pelo ID
// @Tags         Fornecedores
// @Accept       json
// @Produce      json
// @Param        id        path      int                           true  "ID do Fornecedor"
// @Param        supplier  body      domain.UpdateSupplierRequest  true  "Novos dados do fornecedor"
// @Success      200       {object}  utils.Response{data=domain.ApiSupplier}
// @Failure      400       {object}  utils.Response
// @Failure      401       {object}  utils.Response
// @Failure      404       {object}  utils.Response
// @Failure      500       {object}  utils.Response
// @Router       /suppliers/{id} [put]
// @Security     ApiKeyAuth
func (h *SupplierHandler) UpdateSupplier(c *gin.Context) {
	id, err := path.IdFromPathParamOrSendError(c)
	if err != nil {
		return
	}

	var req domain.UpdateSupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		return
	}

	userID, _ := utils.GetUserIDFromContext(c)

	supplier, err := h.supplierService.UpdateSupplier(id, req, userID)
	if err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Fornecedor não encontrado", err.Error())
		} else if validator.IsValidationError(err) {
			utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusBadRequest, "Erro ao atualizar fornecedor", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Fornecedor atualizado com sucesso", supplier, nil)
}

// DeleteSupplier exclui um fornecedor pelo ID
// @Summary      Exclui fornecedor
// @Description  Realiza soft delete de um fornecedor
// @Tags         Fornecedores
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID do Fornecedor"
// @Success      200  {object}  utils.Response
// @Failure      401  {object}  utils.Response
// @Failure      404  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /suppliers/{id} [delete]
// @Security     ApiKeyAuth
func (h *SupplierHandler) DeleteSupplier(c *gin.Context) {
	id, err := path.IdFromPathParamOrSendError(c)
	if err != nil {
		return
	}

	if err := h.supplierService.DeleteSupplier(id); err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Fornecedor não encontrado", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusBadRequest, "Erro ao excluir fornecedor", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Fornecedor excluído com sucesso", nil, nil)
}
