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

// PermissionHandler gerencia as requisições relacionadas a permissões
type PermissionHandler struct {
	permService *service.PermissionService
}

// NewPermissionHandler cria um novo handler de permissões
func NewPermissionHandler(db *gorm.DB) *PermissionHandler {
	roleRepo := repository.NewRoleRepository(db)
	permRepo := repository.NewPermissionRepository(db)

	return &PermissionHandler{
		permService: service.NewPermissionService(roleRepo, permRepo),
	}
}

// GetPermissions retorna todas as permissões cadastradas com paginação e filtros
func (h *PermissionHandler) GetPermissions(c *gin.Context) {
	pagination := utils.GetPaginationParams(c)

	var filters domain.InGetPermissionsFilters
	if err := utils.BindQueryOrSendErrorRes(c, &filters); err != nil {
		return
	}

	permissions, err := h.permService.GetPermissions(&pagination, filters)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao buscar permissões", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Permissões encontradas", permissions, nil)
}

// GetPermission retorna uma permissão específica pelo ID
func (h *PermissionHandler) GetPermission(c *gin.Context) {
	id, err := path.IdFromPathParamOrSendError(c)
	if err != nil {
		return
	}

	permission, err := h.permService.GetPermissionByID(id)
	if err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Permissão não encontrada", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao buscar permissão", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Permissão encontrada", permission, nil)
}

// GetPermissionsByModule retorna permissões agrupadas por módulo
func (h *PermissionHandler) GetPermissionsByModule(c *gin.Context) {
	permissions, err := h.permService.GetPermissionsByModule()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao buscar permissões", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Permissões encontradas", permissions, nil)
}

// GetAvailableModules retorna todos os módulos cadastrados
func (h *PermissionHandler) GetAvailableModules(c *gin.Context) {
	modules, err := h.permService.GetAvailableModules()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao buscar módulos", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Módulos encontrados", modules, nil)
}

// CreatePermission cria uma nova permissão
func (h *PermissionHandler) CreatePermission(c *gin.Context) {
	var body domain.InCreatePermission
	if err := utils.BindJsonOrSendErrorRes(c, &body); err != nil {
		return
	}

	permission, err := h.permService.CreatePermission(body)
	if err != nil {
		if validator.IsValidationError(err) {
			utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		} else if err.Error() == "permissão informada já está em uso" {
			utils.ErrorResponse(c, http.StatusConflict, "Conflito ao criar permissão", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao criar permissão", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Permissão criada com sucesso", permission, nil)
}

// UpdatePermission atualiza uma permissão existente
func (h *PermissionHandler) UpdatePermission(c *gin.Context) {
	id, err := path.IdFromPathParamOrSendError(c)
	if err != nil {
		return
	}

	var body domain.InUpdatePermission
	if err := utils.BindJsonOrSendErrorRes(c, &body); err != nil {
		return
	}

	permission, err := h.permService.UpdatePermission(id, body)
	if err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Permissão não encontrada", err.Error())
		} else if validator.IsValidationError(err) {
			utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		} else if err.Error() == "permissão informada já está em uso" {
			utils.ErrorResponse(c, http.StatusConflict, "Conflito ao atualizar permissão", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao atualizar permissão", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Permissão atualizada com sucesso", permission, nil)
}

// DeletePermission exclui uma permissão
func (h *PermissionHandler) DeletePermission(c *gin.Context) {
	id, err := path.IdFromPathParamOrSendError(c)
	if err != nil {
		return
	}

	if err := h.permService.DeletePermission(id); err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Permissão não encontrada", err.Error())
		} else if validator.IsValidationError(err) {
			utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao excluir permissão", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Permissão excluída com sucesso", nil, nil)
}
