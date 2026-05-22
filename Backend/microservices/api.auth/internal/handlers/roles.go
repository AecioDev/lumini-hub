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

// RoleHandler gerencia as requisições relacionadas a perfis
type RoleHandler struct {
	roleService *service.RoleService
}

// NewRoleHandler cria um novo handler de perfis
func NewRoleHandler(db *gorm.DB) *RoleHandler {
	roleRepo := repository.NewRoleRepository(db)
	userRepo := repository.NewUserRepository(db)
	permRepo := repository.NewPermissionRepository(db)

	return &RoleHandler{
		roleService: service.NewRoleService(roleRepo, userRepo, permRepo),
	}
}

// GetRoles retorna todos os perfis
func (h *RoleHandler) GetRoles(c *gin.Context) {
	roles, err := h.roleService.GetRoles()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao buscar perfis", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Perfis encontrados", roles, nil)
}

// GetRole retorna um perfil específico
func (h *RoleHandler) GetRole(c *gin.Context) {
	id, err := path.IdFromPathParamOrSendError(c)
	if err != nil {
		return
	}

	role, err := h.roleService.GetRoleByID(id)
	if err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Perfil não encontrado", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao buscar perfil", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Perfil encontrado", role, nil)
}

// CreateRole cria um novo perfil
func (h *RoleHandler) CreateRole(c *gin.Context) {
	var req domain.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		return
	}

	role, err := h.roleService.CreateRole(req)
	if err != nil {
		if validator.IsValidationError(err) {
			utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusBadRequest, "Erro ao criar perfil", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Perfil criado com sucesso", role, nil)
}

// UpdateRole atualiza um perfil existente
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	id, err := path.IdFromPathParamOrSendError(c)
	if err != nil {
		return
	}

	var req domain.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		return
	}

	role, err := h.roleService.UpdateRole(id, req)
	if err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Perfil não encontrado", err.Error())
		} else if validator.IsValidationError(err) {
			utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusBadRequest, "Erro ao atualizar perfil", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Perfil updated com sucesso", role, nil)
}

// DeleteRole exclui um perfil
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	id, err := path.IdFromPathParamOrSendError(c)
	if err != nil {
		return
	}

	if err := h.roleService.DeleteRole(id); err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Perfil não encontrado", err.Error())
		} else if validator.IsValidationError(err) {
			utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusBadRequest, "Erro ao excluir perfil", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Perfil excluído com sucesso", nil, nil)
}

// UpdateRolePermissions atualiza as permissões de um perfil
func (h *RoleHandler) UpdateRolePermissions(c *gin.Context) {
	id, err := path.IdFromPathParamOrSendError(c)
	if err != nil {
		return
	}

	var req domain.UpdateRolePermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		return
	}

	role, err := h.roleService.UpdateRolePermissions(id, req.PermissionIDs)
	if err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Perfil não encontrado", err.Error())
		} else if validator.IsValidationError(err) {
			utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusBadRequest, "Erro ao atualizar permissões", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Permissões atualizadas com sucesso", role, nil)
}
