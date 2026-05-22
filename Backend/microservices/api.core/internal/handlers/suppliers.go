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

// SupplierHandler gerencia as requisições relacionadas a fornecedores
type SupplierHandler struct {
	supplierService *service.SupplierService
}

// NewSupplierHandler cria um novo handler de fornecedores
func NewSupplierHandler(db *gorm.DB) *SupplierHandler {
	supplierRepo := repository.NewSupplierRepository(db)

	return &SupplierHandler{
		supplierService: service.NewSupplierService(supplierRepo),
	}
}

// GetSuppliers retorna uma lista paginada de fornecedores
func (h *SupplierHandler) GetSuppliers(c *gin.Context) {
	pagination := utils.GetPaginationParams(c)

	suppliers, err := h.supplierService.GetSuppliers(&pagination)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao buscar fornecedores", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Fornecedores encontrados", suppliers, nil)
}

// GetSupplier retorna um fornecedor específico
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
func (h *SupplierHandler) CreateSupplier(c *gin.Context) {
	var req domain.CreateSupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		return
	}

	// Obter o ID do usuário de auditoria (criador)
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

// DeleteSupplier exclui um fornecedor
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
