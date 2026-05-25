package handlers

import (
	"net/http"
	"strconv"

	"lumini-hub/api.auth/internal/domain"
	"lumini-hub/api.auth/internal/repository"
	"lumini-hub/api.auth/internal/service"
	"lumini-hub/api.auth/internal/validator"
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
func NewUserHandler(db *gorm.DB) *UserHandler {
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)

	return &UserHandler{
		userService: service.NewUserService(userRepo, roleRepo),
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

	users, err := h.userService.GetUsers(&pagination)
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

	user, err := h.userService.GetUserByID(id)
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
		utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		return
	}

	user, err := h.userService.CreateUser(req)
	if err != nil {
		if validator.IsValidationError(err) {
			utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusBadRequest, "Erro ao criar usuário", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Usuário criado com sucesso", user, nil)
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
		utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		return
	}

	user, err := h.userService.UpdateUser(id, req)
	if err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Usuário não encontrado", err.Error())
		} else if validator.IsValidationError(err) {
			utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
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
		utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
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
			utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
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

	if err := h.userService.DeleteUser(id); err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Usuário não encontrado", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusBadRequest, "Erro ao excluir usuário", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Usuário excluído com sucesso", nil, nil)
}
