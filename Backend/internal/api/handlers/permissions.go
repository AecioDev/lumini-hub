package handlers

import (
	"net/http"

	"simple-erp-service/internal/data-structure/dto"
	"simple-erp-service/internal/repository"
	"simple-erp-service/internal/service"
	"simple-erp-service/internal/utils"
	"simple-erp-service/internal/utils/path"
	"simple-erp-service/internal/validator" // Certifique-se que seu validator tenha IsValidationError e lide com a msg de duplicidade

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

// GetPermissions retorna todas as permissões cadastradas
// @Summary Obter todas as permissões
// @Description Retorna uma lista de todas as permissões cadastradas no sistema.
// @Tags permissions
// @Accept json
// @Produce json
// @Security ApiKeyAuth // Mantido para geração automática, mas representa JWT em cookie
// @Success 200 {object} utils.Response{data=[]models.Permission} "Lista de permissões encontrada com sucesso"
// @Failure 500 {object} utils.Response "Erro interno do servidor"
// @Router /permissions [get]
func (h *PermissionHandler) GetPermissions(c *gin.Context) {
	pagination := utils.GetPaginationParams(c)

	var filters dto.InGetPermissionsFilters
	if err := utils.BindQueryOrSendErrorRes(c, &filters); err != nil {
		return
	}

	// Log para depuração: verifique o que foi bindado
	// log.Printf("GetPermissions: Filtros recebidos: PermissionName='%s', Module='%s', RoleId=%d",
	// 	filters.Name, filters.Module, filters.RoleId)
	// log.Printf("GetPermissions: Paginação recebida: Page=%d, Limit=%d",
	// 	pagination.Page, pagination.Limit)

	permissions, err := h.permService.GetPermissions(&pagination, filters)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao buscar permissões", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Permissões encontradas", permissions, nil)
}

// GetPermission retorna uma permissão específica pelo ID
// @Summary Obter uma permissão por ID
// @Description Retorna os detalhes de uma permissão específica usando seu ID.
// @Tags permissions
// @Accept json
// @Produce json
// @Security ApiKeyAuth // Mantido para geração automática, mas representa JWT em cookie
// @Param id path int true "ID da permissão"
// @Success 200 {object} utils.Response{data=models.Permission} "Permissão encontrada com sucesso"
// @Failure 400 {object} utils.Response "ID da permissão inválido"
// @Failure 404 {object} utils.Response "Permissão não encontrada"
// @Failure 500 {object} utils.Response "Erro interno do servidor"
// @Router /permissions/{id} [get]
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
// @Summary Obter permissões agrupadas por módulo
// @Description Retorna uma lista de permissões, agrupadas por seus respectivos módulos.
// @Tags permissions
// @Accept json
// @Produce json
// @Security ApiKeyAuth // Mantido para geração automática, mas representa JWT em cookie
// @Success 200 {object} utils.Response{data=map[string][]models.Permission} "Permissões agrupadas por módulo encontradas com sucesso"
// @Failure 500 {object} utils.Response "Erro interno do servidor"
// @Router /permissions/by-module [get]
func (h *PermissionHandler) GetPermissionsByModule(c *gin.Context) {
	permissions, err := h.permService.GetPermissionsByModule()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao buscar permissões", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Permissões encontradas", permissions, nil)
}

// GetAvailableModules retorna todos os módulos únicos de permissões
// @Summary Obter todos os módulos de permissões
// @Description Retorna uma lista de strings contendo todos os nomes de módulos únicos de permissões.
// @Tags permissions
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} utils.Response{data=[]string} "Lista de módulos encontrada com sucesso"
// @Failure 500 {object} utils.Response "Erro interno do servidor"
// @Router /permissions/modules [get]
func (h *PermissionHandler) GetAvailableModules(c *gin.Context) {
	modules, err := h.permService.GetAvailableModules()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao buscar modulos", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Modulos encontrados", modules, nil)
}

// CreatePermission cria uma nova permissão
// @Summary Criar uma nova permissão
// @Description Cria uma nova permissão no sistema. Requer um nome único, descrição e módulo.
// @Tags permissions
// @Accept json
// @Produce json
// @Security ApiKeyAuth // Mantido para geração automática, mas representa JWT em cookie
// @Param request body dto.InCreatePermission true "Dados da nova permissão"
// @Success 201 {object} utils.Response{data=models.Permission} "Permissão criada com sucesso"
// @Failure 400 {object} utils.Response "Dados inválidos"
// @Failure 409 {object} utils.Response "Conflito: Permissão informada já está em uso"
// @Failure 500 {object} utils.Response "Erro interno do servidor"
// @Router /permissions [post]
func (h *PermissionHandler) CreatePermission(c *gin.Context) {
	var body dto.InCreatePermission
	if err := utils.BindJsonOrSendErrorRes(c, &body); err != nil {
		return
	}

	permission, err := h.permService.CreatePermission(body)
	if err != nil {
		if validator.IsValidationError(err) {
			// Se o validador retorna um erro de validação (ex: formato incorreto)
			utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		} else if err.Error() == "permissão informada já está em uso" { // Verificação da mensagem específica para 409
			utils.ErrorResponse(c, http.StatusConflict, "Conflito ao criar permissão", err.Error())
		} else {
			// Para outros erros internos do serviço
			utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao criar permissão", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Permissão criada com sucesso", permission, nil)
}

// UpdatePermission atualiza uma permissão existente
// @Summary Atualizar permissão
// @Description Atualiza os dados de uma permissão existente pelo seu ID.
// @Tags permissions
// @Accept json
// @Produce json
// @Security ApiKeyAuth // Mantido para geração automática, mas representa JWT em cookie
// @Param id path int true "ID da permissão"
// @Param request body dto.InUpdatePermission true "Dados da permissão para atualização"
// @Success 200 {object} utils.Response{data=models.Permission} "Permissão atualizada com sucesso"
// @Failure 400 {object} utils.Response "Dados inválidos ou ID da permissão inválido"
// @Failure 404 {object} utils.Response "Permissão não encontrada"
// @Failure 409 {object} utils.Response "Conflito: Permissão informada já está em uso"
// @Failure 500 {object} utils.Response "Erro interno do servidor"
// @Router /permissions/{id} [put]
func (h *PermissionHandler) UpdatePermission(c *gin.Context) {
	id, err := path.IdFromPathParamOrSendError(c)
	if err != nil {
		return
	}

	var body dto.InUpdatePermission
	if err := utils.BindJsonOrSendErrorRes(c, &body); err != nil {
		return
	}

	permission, err := h.permService.UpdatePermission(id, body)
	if err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Permissão não encontrada", err.Error())
		} else if validator.IsValidationError(err) {
			utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		} else if err.Error() == "permissão informada já está em uso" { // Verificação da mensagem específica para 409
			utils.ErrorResponse(c, http.StatusConflict, "Conflito ao atualizar permissão", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao atualizar permissão", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Permissão atualizada com sucesso", permission, nil)
}

// DeletePermission exclui uma permissão
// @Summary Excluir uma permissão
// @Description Exclui uma permissão existente pelo seu ID.
// @Tags permissions
// @Accept json
// @Produce json
// @Security ApiKeyAuth // Mantido para geração automática, mas representa JWT em cookie
// @Param id path int true "ID da permissão"
// @Success 200 {object} utils.Response "Permissão excluída com sucesso"
// @Failure 400 {object} utils.Response "ID da permissão inválido"
// @Failure 404 {object} utils.Response "Permissão não encontrada"
// @Failure 500 {object} utils.Response "Erro interno do servidor"
// @Router /permissions/{id} [delete]
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
