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
// @Summary Listar fornecedores
// @Description Retorna uma lista paginada de fornecedores
// @Tags suppliers
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "Número da página" default(1)
// @Param limit query int false "Limite de itens por página" default(10)
// @Param sort query string false "Campo para ordenação" default(id)
// @Param order query string false "Direção da ordenação (asc/desc)" default(asc)
// @Success 200 {object} utils.Response "Fornecedores encontrados"
// @Failure 401 {object} utils.Response "Não autorizado"
// @Failure 500 {object} utils.Response "Erro ao buscar fornecedores"
// @Router /suppliers [get]
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
// @Summary Buscar fornecedor
// @Description Retorna um fornecedor específico pelo ID
// @Tags suppliers
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "ID do fornecedor"
// @Success 200 {object} utils.Response "Fornecedor encontrado"
// @Failure 400 {object} utils.Response "ID inválido"
// @Failure 401 {object} utils.Response "Não autorizado"
// @Failure 404 {object} utils.Response "Fornecedor não encontrado"
// @Router /suppliers/{id} [get]
func (h *SupplierHandler) GetSupplier(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido", err.Error())
		return
	}

	supplier, err := h.supplierService.GetSupplierByID(uint(id))
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
// @Summary Criar fornecedor
// @Description Cria um novo fornecedor
// @Tags suppliers
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body models.CreateSupplierRequest true "Dados do fornecedor"
// @Success 201 {object} utils.Response "Fornecedor criado com sucesso"
// @Failure 400 {object} utils.Response "Dados inválidos"
// @Failure 401 {object} utils.Response "Não autorizado"
// @Router /suppliers [post]
func (h *SupplierHandler) CreateSupplier(c *gin.Context) {
	var req models.CreateSupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		return
	}

	supplier, err := h.supplierService.CreateSupplier(req)
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
// @Summary Atualizar fornecedor
// @Description Atualiza um fornecedor existente
// @Tags suppliers
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "ID do fornecedor"
// @Param request body models.UpdateSupplierRequest true "Dados do fornecedor"
// @Success 200 {object} utils.Response "Fornecedor atualizado com sucesso"
// @Failure 400 {object} utils.Response "Dados inválidos"
// @Failure 401 {object} utils.Response "Não autorizado"
// @Failure 404 {object} utils.Response "Fornecedor não encontrado"
// @Router /suppliers/{id} [put]
func (h *SupplierHandler) UpdateSupplier(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido", err.Error())
		return
	}

	var req models.UpdateSupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		return
	}

	supplier, err := h.supplierService.UpdateSupplier(uint(id), req)
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
// @Summary Excluir fornecedor
// @Description Exclui um fornecedor
// @Tags suppliers
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "ID do fornecedor"
// @Success 200 {object} utils.Response "Fornecedor excluído com sucesso"
// @Failure 400 {object} utils.Response "ID inválido"
// @Failure 401 {object} utils.Response "Não autorizado"
// @Failure 404 {object} utils.Response "Fornecedor não encontrado"
// @Router /suppliers/{id} [delete]
func (h *SupplierHandler) DeleteSupplier(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido", err.Error())
		return
	}

	if err := h.supplierService.DeleteSupplier(uint(id)); err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Fornecedor não encontrado", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusBadRequest, "Erro ao excluir fornecedor", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Fornecedor excluído com sucesso", nil, nil)
}
