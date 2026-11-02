package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"lumini-hub/api.auth/internal/domain"
	"lumini-hub/api.auth/internal/repository"
	"lumini-hub/api.auth/internal/service"
	"lumini-hub/api.auth/internal/validator"
	"lumini-hub/common/config"
	"lumini-hub/common/internalmail"
	"lumini-hub/common/utils"
	"lumini-hub/common/utils/path"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UserHandler gerencia as requisições relacionadas a usuários
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler cria um novo handler de usuários
func NewUserHandler(db *gorm.DB, cfg *config.Config) *UserHandler {
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	permRepo := repository.NewPermissionRepository(db)

	return &UserHandler{
		userService: service.NewUserService(userRepo, roleRepo, permRepo,
			service.NewPasswordResetService(db, internalmail.NewClient("api.auth", cfg))),
	}
}

// GetUsers retorna uma lista paginada de usuários
// @Summary      Lista usuários com paginação
// @Description  Retorna uma lista paginada de usuários cadastrados no ERP
// @Tags         Usuários
// @Accept       json
// @Produce      json
// @Param        page   query      int     false  "Página ativa"
// @Param        limit  query      int     false  "Limite por página"
// @Param        sort   query      string  false  "Ordenação (ex: username asc)"
// @Success      200    {object}  utils.Response{data=domain.ApiUserListPaginated}
// @Failure      401    {object}  utils.Response
// @Failure      500    {object}  utils.Response
// @Router       /users [get]
// @Security     ApiKeyAuth
func (h *UserHandler) GetUsers(c *gin.Context) {
	pagination := utils.GetPaginationParams(c)

	users, err := h.userService.GetUsers(&pagination, utils.RoleFromGinContext(c))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao buscar usuários", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Usuários encontrados", users, nil)
}

// PostUserFilter busca usuários aplicando os filtros recebidos no corpo (PostFilter)
// @Summary      Filtra usuários com paginação (PostFilter)
// @Description  Busca usuários por texto (nome, usuário ou e-mail), perfil e status, com paginação e ordenação
// @Tags         Usuários
// @Accept       json
// @Produce      json
// @Param        filter  body      domain.UserFilterRequest  true  "Filtros e dados de paginação"
// @Success      200     {object}  utils.Response{data=domain.ApiUserListPaginated}
// @Failure      400     {object}  utils.Response
// @Failure      401     {object}  utils.Response
// @Failure      500     {object}  utils.Response
// @Router       /users/filter [post]
// @Security     ApiKeyAuth
func (h *UserHandler) PostUserFilter(c *gin.Context) {
	var req domain.UserFilterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Parâmetros de filtro inválidos", err)
		return
	}

	users, err := h.userService.GetUsersByFilter(req, utils.RoleFromGinContext(c))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao buscar usuários", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Usuários encontrados", users, nil)
}

// GetUser retorna um usuário específico
// @Summary      Busca usuário por ID
// @Description  Retorna os detalhes de um único usuário
// @Tags         Usuários
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID do Usuário"
// @Success      200  {object}  utils.Response{data=domain.ApiUserDetail}
// @Failure      401  {object}  utils.Response
// @Failure      404  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /users/{id} [get]
// @Security     ApiKeyAuth
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := path.IdFromPathParamOrSendError(c)
	if err != nil {
		return
	}

	user, err := h.userService.GetUserByID(id, utils.RoleFromGinContext(c))
	if err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Usuário não encontrado", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao buscar usuário", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Usuário encontrado", user, nil)
}

// CreateUser cria um novo usuário
// @Summary      Cria novo usuário
// @Description  Cadastra um usuário com perfil (Role) associado
// @Tags         Usuários
// @Accept       json
// @Produce      json
// @Param        user  body      domain.CreateUserRequest  true  "Dados do Usuário"
// @Success      201   {object}  utils.Response{data=domain.ApiUserDetail}
// @Failure      400   {object}  utils.Response
// @Failure      401   {object}  utils.Response
// @Failure      500   {object}  utils.Response
// @Router       /users [post]
// @Security     ApiKeyAuth
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req domain.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Dados inválidos", err)
		return
	}

	user, err := h.userService.CreateUser(req, utils.RoleFromGinContext(c), usernameFromContext(c))
	if err != nil {
		if validator.IsValidationError(err) {
			utils.ValidationErrorResponse(c, "Dados inválidos", err)
		} else {
			utils.ErrorResponse(c, http.StatusBadRequest, "Erro ao criar usuário", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Usuário criado com sucesso", user, nil)
}

// usernameFromContext devolve o usuário logado (para citar no convite quem cadastrou)
func usernameFromContext(c *gin.Context) string {
	username, _ := c.Get("username")
	name, _ := username.(string)
	return name
}

// ResendInvite reenvia o convite de um usuário que ainda não entrou no sistema
// @Summary      Reenvia o convite do usuário
// @Description  Gera um novo link de convite (o anterior deixa de valer) e reenvia o e-mail para o usuário que ainda não entrou no sistema
// @Tags         Usuários
// @Produce      json
// @Param        id   path      int  true  "ID do Usuário"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      404  {object}  utils.Response
// @Failure      502  {object}  utils.Response
// @Router       /users/{id}/invite [post]
// @Security     ApiKeyAuth
func (h *UserHandler) ResendInvite(c *gin.Context) {
	id, err := path.IdFromPathParamOrSendError(c)
	if err != nil {
		return
	}

	err = h.userService.ResendInvite(id, utils.RoleFromGinContext(c), usernameFromContext(c))
	switch {
	case err == nil:
		utils.SuccessResponse(c, http.StatusOK, "Convite reenviado com sucesso", nil, nil)
	case errors.Is(err, utils.ErrNotFound):
		utils.ErrorResponse(c, http.StatusNotFound, "Usuário não encontrado", err.Error())
	case errors.Is(err, service.ErrInviteNotAllowed):
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), err.Error())
	default:
		// Falha do e-mail (servidor recusou, e-mail não configurado...): o texto vem do api.core e é útil a quem cadastra.
		utils.ErrorResponse(c, http.StatusBadGateway, "Não foi possível enviar o convite", err.Error())
	}
}

// UpdateUser atualiza um usuário existente
// @Summary      Atualiza usuário
// @Description  Altera os dados cadastrais e o perfil de um usuário existente
// @Tags         Usuários
// @Accept       json
// @Produce      json
// @Param        id    path      int                       true  "ID do Usuário"
// @Param        user  body      domain.UpdateUserRequest  true  "Novos dados do usuário"
// @Success      200   {object}  utils.Response{data=domain.ApiUserDetail}
// @Failure      400   {object}  utils.Response
// @Failure      401   {object}  utils.Response
// @Failure      404   {object}  utils.Response
// @Failure      500   {object}  utils.Response
// @Router       /users/{id} [put]
// @Security     ApiKeyAuth
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := path.IdFromPathParamOrSendError(c)
	if err != nil {
		return
	}

	var req domain.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Dados inválidos", err)
		return
	}

	user, err := h.userService.UpdateUser(id, req, utils.RoleFromGinContext(c))
	if err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Usuário não encontrado", err.Error())
		} else if validator.IsValidationError(err) {
			utils.ValidationErrorResponse(c, "Dados inválidos", err)
		} else {
			utils.ErrorResponse(c, http.StatusBadRequest, "Erro ao atualizar usuário", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Usuário atualizado com sucesso", user, nil)
}

// ChangePassword altera a senha de um usuário
// @Summary      Altera senha do usuário
// @Description  Altera a credencial de senha de um usuário específico (requer validação se não for Admin)
// @Tags         Usuários
// @Accept       json
// @Produce      json
// @Param        id    path      int                          true  "ID do Usuário"
// @Param        data  body      domain.ChangePasswordRequest  true  "Dados de alteração de senha"
// @Success      200   {object}  utils.Response
// @Failure      400   {object}  utils.Response
// @Failure      401   {object}  utils.Response
// @Failure      403   {object}  utils.Response
// @Router       /users/{id}/password [put]
// @Security     ApiKeyAuth
func (h *UserHandler) ChangePassword(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido", err.Error())
		return
	}

	var req domain.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Dados inválidos", err)
		return
	}

	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	isAdmin := role == "ADMIN"
	isSelf := userID.(uint) == uint(id)

	if !isAdmin && !isSelf {
		utils.ErrorResponse(c, http.StatusForbidden, "Acesso negado", "Você não tem permissão para alterar a senha de outro usuário")
		return
	}

	if isAdmin && !isSelf {
		err = h.userService.ChangePassword(uint(id), "", req.NewPassword, true)
	} else {
		err = h.userService.ChangePassword(uint(id), req.CurrentPassword, req.NewPassword, false)
	}

	if err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Usuário não encontrado", err.Error())
		} else if err == utils.ErrInvalidCredentials {
			utils.ErrorResponse(c, http.StatusBadRequest, "Senha atual incorreta", err.Error())
		} else if validator.IsValidationError(err) {
			utils.ValidationErrorResponse(c, "Dados inválidos", err)
		} else {
			utils.ErrorResponse(c, http.StatusBadRequest, "Erro ao alterar senha", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Senha alterada com sucesso", nil, nil)
}

// DeleteUser exclui um usuário
// @Summary      Exclui usuário
// @Description  Remove logicamente um usuário do sistema (não é permitido excluir a si mesmo)
// @Tags         Usuários
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID do Usuário"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      401  {object}  utils.Response
// @Failure      404  {object}  utils.Response
// @Router       /users/{id} [delete]
// @Security     ApiKeyAuth
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := path.IdFromPathParamOrSendError(c)
	if err != nil {
		return
	}

	userID, _ := c.Get("userID")
	if userID.(uint) == uint(id) {
		utils.ErrorResponse(c, http.StatusBadRequest, "Operação inválida", "Você não pode excluir seu próprio usuário")
		return
	}

	if err := h.userService.DeleteUser(id, utils.RoleFromGinContext(c)); err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Usuário não encontrado", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusBadRequest, "Erro ao excluir usuário", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Usuário excluído com sucesso", nil, nil)
}

// UpdateUserPermissions atualiza as permissões diretas de um usuário
// @Summary      Atualiza permissões do usuário
// @Description  Associa permissões diretas a um usuário específico (sobrescrevendo as anteriores)
// @Tags         Usuários
// @Accept       json
// @Produce      json
// @Param        id    path      int                                  true  "ID do Usuário"
// @Param        data  body      domain.UpdateUserPermissionsRequest  true  "IDs das novas permissões"
// @Success      200   {object}  utils.Response{data=domain.ApiUserDetail}
// @Failure      400   {object}  utils.Response
// @Failure      401   {object}  utils.Response
// @Failure      404   {object}  utils.Response
// @Failure      500   {object}  utils.Response
// @Router       /users/{id}/permissions [put]
// @Security     ApiKeyAuth
func (h *UserHandler) UpdateUserPermissions(c *gin.Context) {
	id, err := path.IdFromPathParamOrSendError(c)
	if err != nil {
		return
	}

	var req domain.UpdateUserPermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Dados inválidos", err)
		return
	}

	user, err := h.userService.UpdateUserPermissions(id, req.PermissionIDs, utils.RoleFromGinContext(c))
	if err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Usuário não encontrado", err.Error())
		} else if validator.IsValidationError(err) {
			utils.ValidationErrorResponse(c, "Dados inválidos", err)
		} else {
			utils.ErrorResponse(c, http.StatusBadRequest, "Erro ao atualizar permissões do usuário", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Permissões do usuário atualizadas com sucesso", user, nil)
}

