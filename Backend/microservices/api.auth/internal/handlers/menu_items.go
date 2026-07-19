package handlers

import (
	"net/http"

	"lumini-hub/api.auth/internal/domain"
	"lumini-hub/api.auth/internal/repository"
	"lumini-hub/api.auth/internal/service"
	"lumini-hub/api.auth/internal/validator"
	"lumini-hub/common/utils"
	"lumini-hub/common/utils/path"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// MenuItemHandler gerencia as requisições relacionadas a itens do menu lateral
type MenuItemHandler struct {
	menuItemService *service.MenuItemService
}

// NewMenuItemHandler cria um novo handler de itens de menu
func NewMenuItemHandler(db *gorm.DB) *MenuItemHandler {
	menuItemRepo := repository.NewMenuItemRepository(db)
	permRepo := repository.NewPermissionRepository(db)

	return &MenuItemHandler{
		menuItemService: service.NewMenuItemService(menuItemRepo, permRepo),
	}
}

// GetMenuTree retorna a árvore completa de itens de menu, sem filtro de permissão
func (h *MenuItemHandler) GetMenuTree(c *gin.Context) {
	tree, err := h.menuItemService.GetMenuTree()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao buscar itens de menu", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Itens de menu encontrados", tree, nil)
}

// GetMenuItem retorna um item de menu específico pelo ID
func (h *MenuItemHandler) GetMenuItem(c *gin.Context) {
	id, err := path.IdFromPathParamOrSendError(c)
	if err != nil {
		return
	}

	item, err := h.menuItemService.GetMenuItemByID(id)
	if err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Item de menu não encontrado", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao buscar item de menu", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Item de menu encontrado", item, nil)
}

// CreateMenuItem cria um novo item de menu
func (h *MenuItemHandler) CreateMenuItem(c *gin.Context) {
	var body domain.CreateMenuItemRequest
	if err := utils.BindJsonOrSendErrorRes(c, &body); err != nil {
		return
	}

	item, err := h.menuItemService.CreateMenuItem(body)
	if err != nil {
		if validator.IsValidationError(err) {
			utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao criar item de menu", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Item de menu criado com sucesso", item, nil)
}

// UpdateMenuItem atualiza um item de menu existente
func (h *MenuItemHandler) UpdateMenuItem(c *gin.Context) {
	id, err := path.IdFromPathParamOrSendError(c)
	if err != nil {
		return
	}

	var body domain.UpdateMenuItemRequest
	if err := utils.BindJsonOrSendErrorRes(c, &body); err != nil {
		return
	}

	item, err := h.menuItemService.UpdateMenuItem(id, body)
	if err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Item de menu não encontrado", err.Error())
		} else if validator.IsValidationError(err) {
			utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao atualizar item de menu", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Item de menu atualizado com sucesso", item, nil)
}

// DeleteMenuItem exclui um item de menu
func (h *MenuItemHandler) DeleteMenuItem(c *gin.Context) {
	id, err := path.IdFromPathParamOrSendError(c)
	if err != nil {
		return
	}

	if err := h.menuItemService.DeleteMenuItem(id); err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Item de menu não encontrado", err.Error())
		} else if validator.IsValidationError(err) {
			utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao excluir item de menu", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Item de menu excluído com sucesso", nil, nil)
}
