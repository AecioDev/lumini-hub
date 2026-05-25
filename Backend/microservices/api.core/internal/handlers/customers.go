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

// CustomerHandler gerencia as requisições relacionadas a clientes
type CustomerHandler struct {
	customerService *service.CustomerService
}

// NewCustomerHandler cria um novo handler de clientes recebendo o Unit of Work
func NewCustomerHandler(uow repository.UnitOfWork) *CustomerHandler {
	return &CustomerHandler{
		customerService: service.NewCustomerService(uow),
	}
}

// GetCustomers retorna uma lista paginada de clientes (GET legada)
// @Summary      Lista clientes com paginação (GET)
// @Description  Retorna uma lista de clientes paginada a partir de query strings
// @Tags         Clientes
// @Accept       json
// @Produce      json
// @Param        page   query      int  false  "Página ativa"
// @Param        limit  query      int  false  "Limite por página"
// @Param        sort   query      string  false  "Ordenação (ex: first_name asc)"
// @Success      200    {object}  utils.Response{data=domain.ApiCustomerListPaginated}
// @Failure      401    {object}  utils.Response
// @Failure      500    {object}  utils.Response
// @Router       /customers [get]
// @Security     ApiKeyAuth
func (h *CustomerHandler) GetCustomers(c *gin.Context) {
	pagination := utils.GetPaginationParams(c)

	customers, err := h.customerService.GetCustomers(&pagination)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao buscar clientes", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Clientes encontrados", customers, nil)
}

// PostCustomerFilter busca clientes aplicando filtros complexos recebidos no Body (PostFilter)
// @Summary      Filtra clientes com paginação (PostFilter)
// @Description  Busca e filtra clientes dinamicamente enviando as opções no JSON do Body
// @Tags         Clientes
// @Accept       json
// @Produce      json
// @Param        filter  body      domain.CustomerFilterRequest  true  "Filtros e Dados de Paginação"
// @Success      200     {object}  utils.Response{data=domain.ApiCustomerListPaginated}
// @Failure      400     {object}  utils.Response
// @Failure      401     {object}  utils.Response
// @Failure      500     {object}  utils.Response
// @Router       /customers/filter [post]
// @Security     ApiKeyAuth
func (h *CustomerHandler) PostCustomerFilter(c *gin.Context) {
	var req domain.CustomerFilterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Parâmetros de filtro inválidos", err.Error())
		return
	}

	customers, err := h.customerService.GetCustomersByFilter(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao buscar clientes com filtros", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Clientes encontrados", customers, nil)
}

// GetCustomer retorna um cliente específico pelo ID
// @Summary      Busca cliente por ID
// @Description  Retorna os detalhes de um único cliente cadastrado
// @Tags         Clientes
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID do Cliente"
// @Success      200  {object}  utils.Response{data=domain.ApiCustomerDetail}
// @Failure      401  {object}  utils.Response
// @Failure      404  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /customers/{id} [get]
// @Security     ApiKeyAuth
func (h *CustomerHandler) GetCustomer(c *gin.Context) {
	id, err := path.IdFromPathParamOrSendError(c)
	if err != nil {
		return
	}

	customer, err := h.customerService.GetCustomerByID(id)
	if err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Cliente não encontrado", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao buscar cliente", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Cliente encontrado", customer, nil)
}

// CreateCustomer cria um novo cliente
// @Summary      Cria novo cliente
// @Description  Cadastra um cliente no sistema ERP
// @Tags         Clientes
// @Accept       json
// @Produce      json
// @Param        customer  body      domain.CreateCustomerRequest  true  "Dados do Cliente"
// @Success      201       {object}  utils.Response{data=domain.ApiCustomer}
// @Failure      400       {object}  utils.Response
// @Failure      401       {object}  utils.Response
// @Failure      500       {object}  utils.Response
// @Router       /customers [post]
// @Security     ApiKeyAuth
func (h *CustomerHandler) CreateCustomer(c *gin.Context) {
	var req domain.CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		return
	}

	userID, exists := utils.GetUserIDFromContext(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Usuário não autenticado", "")
		return
	}

	customer, err := h.customerService.CreateCustomer(req, userID)
	if err != nil {
		if validator.IsValidationError(err) {
			utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusBadRequest, "Erro ao criar cliente", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Cliente criado com sucesso", customer, nil)
}

// UpdateCustomer atualiza um cliente existente
// @Summary      Atualiza cliente
// @Description  Altera dados cadastrais de um cliente pelo ID
// @Tags         Clientes
// @Accept       json
// @Produce      json
// @Param        id        path      int                           true  "ID do Cliente"
// @Param        customer  body      domain.UpdateCustomerRequest  true  "Novos dados do cliente"
// @Success      200       {object}  utils.Response{data=domain.ApiCustomer}
// @Failure      400       {object}  utils.Response
// @Failure      401       {object}  utils.Response
// @Failure      404       {object}  utils.Response
// @Failure      500       {object}  utils.Response
// @Router       /customers/{id} [put]
// @Security     ApiKeyAuth
func (h *CustomerHandler) UpdateCustomer(c *gin.Context) {
	id, err := path.IdFromPathParamOrSendError(c)
	if err != nil {
		return
	}

	var req domain.UpdateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		return
	}

	customer, err := h.customerService.UpdateCustomer(id, req)
	if err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Cliente não encontrado", err.Error())
		} else if validator.IsValidationError(err) {
			utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusBadRequest, "Erro ao atualizar cliente", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Cliente atualizado com sucesso", customer, nil)
}

// DeleteCustomer exclui um cliente pelo ID
// @Summary      Exclui cliente
// @Description  Realiza soft delete de um cliente
// @Tags         Clientes
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID do Cliente"
// @Success      200  {object}  utils.Response
// @Failure      401  {object}  utils.Response
// @Failure      404  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /customers/{id} [delete]
// @Security     ApiKeyAuth
func (h *CustomerHandler) DeleteCustomer(c *gin.Context) {
	id, err := path.IdFromPathParamOrSendError(c)
	if err != nil {
		return
	}

	if err := h.customerService.DeleteCustomer(id); err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Cliente não encontrado", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusBadRequest, "Erro ao excluir cliente", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Cliente excluído com sucesso", nil, nil)
}
