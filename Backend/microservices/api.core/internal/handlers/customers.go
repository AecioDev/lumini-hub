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
func (h *CustomerHandler) CreateCustomer(c *gin.Context) {
	var req domain.CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		return
	}

	// Pegar o ID do usuário autenticado via middleware comum de auth
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

// DeleteCustomer exclui um cliente
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
