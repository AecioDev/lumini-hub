package handlers

import (
	"net/http"
	"strconv"

	"simple-erp-service/internal/data-structure/models"
	"simple-erp-service/internal/repository"
	"simple-erp-service/internal/service"
	"simple-erp-service/internal/utils"
	"simple-erp-service/internal/utils/path"
	"simple-erp-service/internal/validator"

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

// GetRoles retorna uma lista paginada de perfis
// @Summary Listar perfis
// @Description Retorna uma lista paginada de perfis
// @Tags roles
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "Número da página" default(1)
// @Param limit query int false "Limite de itens por página" default(10)
// @Param sort query string false "Campo para ordenação" default(id)
// @Param order query string false "Direção da ordenação (asc/desc)" default(asc)
// @Success 200 {object} utils.Response "Perfis encontrados"
// @Failure 401 {object} utils.Response "Não autorizado"
// @Failure 500 {object} utils.Response "Erro ao buscar perfis"
// @Router /roles [get]
func (h *RoleHandler) GetRoles(c *gin.Context) {
	roles, err := h.roleService.GetRoles()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao buscar perfis", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Perfis encontrados", roles, nil)
}

// GetRole retorna um perfil específico
// @Summary Buscar perfil
// @Description Retorna um perfil específico pelo ID
// @Tags roles
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "ID do perfil"
// @Success 200 {object} utils.Response "Perfil encontrado"
// @Failure 400 {object} utils.Response "ID inválido"
// @Failure 401 {object} utils.Response "Não autorizado"
// @Failure 404 {object} utils.Response "Perfil não encontrado"
// @Router /roles/{id} [get]
func (h *RoleHandler) GetRole(c *gin.Context) {
	id, err := path.IdFromPathParamOrSendError(c)
	if err != nil {
		return
	}

	role, err := h.roleService.GetRoleByID(uint(id))
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

// CreateRole cria um novo perfil ADMIN, FINANCEIRO, GERENTE, ETC.
// @Summary Criar perfil
// @Description Cria um novo perfil
// @Tags roles
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body models.CreateRoleRequest true "Dados do perfil"
// @Success 201 {object} utils.Response "Perfil criado com sucesso"
// @Failure 400 {object} utils.Response "Dados inválidos"
// @Failure 401 {object} utils.Response "Não autorizado"
// @Router /roles [post]
func (h *RoleHandler) CreateRole(c *gin.Context) {
	var req models.CreateRoleRequest
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
// @Summary Atualizar perfil
// @Description Atualiza um perfil existente
// @Tags roles
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "ID do perfil"
// @Param request body models.UpdateRoleRequest true "Dados do perfil"
// @Success 200 {object} utils.Response "Perfil atualizado com sucesso"
// @Failure 400 {object} utils.Response "Dados inválidos"
// @Failure 401 {object} utils.Response "Não autorizado"
// @Failure 404 {object} utils.Response "Perfil não encontrado"
// @Router /roles/{id} [put]
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido", err.Error())
		return
	}

	var req models.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		return
	}

	role, err := h.roleService.UpdateRole(uint(id), req)
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

	utils.SuccessResponse(c, http.StatusOK, "Perfil atualizado com sucesso", role, nil)
}

// DeleteRole exclui um perfil
// @Summary Excluir perfil
// @Description Exclui um perfil
// @Tags roles
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "ID do perfil"
// @Success 200 {object} utils.Response "Perfil excluído com sucesso"
// @Failure 400 {object} utils.Response "ID inválido"
// @Failure 401 {object} utils.Response "Não autorizado"
// @Failure 404 {object} utils.Response "Perfil não encontrado"
// @Router /roles/{id} [delete]
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido", err.Error())
		return
	}

	if err := h.roleService.DeleteRole(uint(id)); err != nil {
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
// @Summary Atualizar permissões de perfil
// @Description Atualiza as permissões de um perfil
// @Tags roles
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "ID do perfil"
// @Param request body models.UpdateRolePermissionsRequest true "IDs das permissões"
// @Success 200 {object} utils.Response "Permissões atualizadas com sucesso"
// @Failure 400 {object} utils.Response "Dados inválidos"
// @Failure 401 {object} utils.Response "Não autorizado"
// @Failure 404 {object} utils.Response "Perfil não encontrado"
// @Router /roles/{id}/permissions [put]
func (h *RoleHandler) UpdateRolePermissions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido", err.Error())
		return
	}

	var req models.UpdateRolePermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		return
	}

	role, err := h.roleService.UpdateRolePermissions(uint(id), req.PermissionIDs)
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
