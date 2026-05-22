package handlers

import (
	"net/http"
	"strconv"

	"simple-erp-service/internal/data-structure/models"
	"simple-erp-service/internal/repository"
	"simple-erp-service/internal/service"
	"simple-erp-service/internal/utils"
	"simple-erp-service/internal/validator"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CustomerHandler gerencia as requisições relacionadas a clientes
type CustomerHandler struct {
	customerService *service.CustomerService
}

// NewCustomerHandler cria um novo handler de clientes
func NewCustomerHandler(db *gorm.DB) *CustomerHandler {
	customerRepo := repository.NewCustomerRepository(db)

	return &CustomerHandler{
		customerService: service.NewCustomerService(customerRepo),
	}
}

// GetCustomers retorna uma lista paginada de clientes
// @Summary Listar clientes
// @Description Retorna uma lista paginada de clientes
// @Tags customers
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "Número da página" default(1)
// @Param limit query int false "Limite de itens por página" default(10)
// @Param sort query string false "Campo para ordenação" default(id)
// @Param order query string false "Direção da ordenação (asc/desc)" default(asc)
// @Success 200 {object} utils.Response "Clientes encontrados"
// @Failure 401 {object} utils.Response "Não autorizado"
// @Failure 500 {object} utils.Response "Erro ao buscar clientes"
// @Router /customers [get]
func (h *CustomerHandler) GetCustomers(c *gin.Context) {
	pagination := utils.GetPaginationParams(c)

	customers, err := h.customerService.GetCustomers(&pagination)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao buscar clientes", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Clientes encontrados", customers, nil)
}

// GetCustomer retorna um cliente específico
// @Summary Buscar cliente
// @Description Retorna um cliente específico pelo ID
// @Tags customers
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "ID do cliente"
// @Success 200 {object} utils.Response "Cliente encontrado"
// @Failure 400 {object} utils.Response "ID inválido"
// @Failure 401 {object} utils.Response "Não autorizado"
// @Failure 404 {object} utils.Response "Cliente não encontrado"
// @Router /customers/{id} [get]
func (h *CustomerHandler) GetCustomer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido", err.Error())
		return
	}

	customer, err := h.customerService.GetCustomerByID(uint(id))
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
// @Summary Criar cliente
// @Description Cria um novo cliente
// @Tags customers
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body models.CreateCustomerRequest true "Dados do cliente"
// @Success 201 {object} utils.Response "Cliente criado com sucesso"
// @Failure 400 {object} utils.Response "Dados inválidos"
// @Failure 401 {object} utils.Response "Não autorizado"
// @Router /customers [post]
func (h *CustomerHandler) CreateCustomer(c *gin.Context) {
	var req models.CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		return
	}

	// Aqui você pega o ID do usuário autenticado via middleware
	userID, exists := utils.GetUserIDFromContext(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Usuário não autenticado", "")
		return
	}

	// Passa o ID do usuário separadamente
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
// @Summary Atualizar cliente
// @Description Atualiza um cliente existente
// @Tags customers
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "ID do cliente"
// @Param request body models.UpdateCustomerRequest true "Dados do cliente"
// @Success 200 {object} utils.Response "Cliente atualizado com sucesso"
// @Failure 400 {object} utils.Response "Dados inválidos"
// @Failure 401 {object} utils.Response "Não autorizado"
// @Failure 404 {object} utils.Response "Cliente não encontrado"
// @Router /customers/{id} [put]
func (h *CustomerHandler) UpdateCustomer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido", err.Error())
		return
	}

	var req models.UpdateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		return
	}

	customer, err := h.customerService.UpdateCustomer(uint(id), req)
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

// DeleteCustomer exclui um cliente
// @Summary Excluir cliente
// @Description Exclui um cliente
// @Tags customers
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "ID do cliente"
// @Success 200 {object} utils.Response "Cliente excluído com sucesso"
// @Failure 400 {object} utils.Response "ID inválido"
// @Failure 401 {object} utils.Response "Não autorizado"
// @Failure 404 {object} utils.Response "Cliente não encontrado"
// @Router /customers/{id} [delete]
func (h *CustomerHandler) DeleteCustomer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido", err.Error())
		return
	}

	if err := h.customerService.DeleteCustomer(uint(id)); err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Cliente não encontrado", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusBadRequest, "Erro ao excluir cliente", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Cliente excluído com sucesso", nil, nil)
}
