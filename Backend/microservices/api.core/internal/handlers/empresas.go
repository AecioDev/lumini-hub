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
)

// EmpresaHandler gerencia as requisições relacionadas a empresas
type EmpresaHandler struct {
	empresaService *service.EmpresaService
}

// NewEmpresaHandler cria um novo handler de empresas recebendo o Unit of Work
func NewEmpresaHandler(uow repository.UnitOfWork) *EmpresaHandler {
	return &EmpresaHandler{
		empresaService: service.NewEmpresaService(uow),
	}
}

// GetEmpresas retorna todas as empresas cadastradas
// @Summary      Lista empresas
// @Description  Retorna todas as empresas cadastradas (Matriz e vinculadas)
// @Tags         Empresas
// @Accept       json
// @Produce      json
// @Success      200  {object}  utils.Response{data=domain.ApiEmpresaList}
// @Failure      401  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /empresas [get]
// @Security     ApiKeyAuth
func (h *EmpresaHandler) GetEmpresas(c *gin.Context) {
	empresas, err := h.empresaService.GetEmpresas()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao buscar empresas", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Empresas encontradas", empresas, nil)
}

// GetEmpresa retorna uma empresa específica pelo ID
// @Summary      Busca empresa por ID
// @Description  Retorna os detalhes de uma única empresa cadastrada
// @Tags         Empresas
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID da Empresa"
// @Success      200  {object}  utils.Response{data=domain.ApiEmpresaDetail}
// @Failure      401  {object}  utils.Response
// @Failure      404  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /empresas/{id} [get]
// @Security     ApiKeyAuth
func (h *EmpresaHandler) GetEmpresa(c *gin.Context) {
	id, err := path.IdFromPathParamOrSendError(c)
	if err != nil {
		return
	}

	empresa, err := h.empresaService.GetEmpresaByID(id)
	if err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Empresa não encontrada", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao buscar empresa", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Empresa encontrada", empresa, nil)
}

// CreateEmpresa cria uma nova empresa
// @Summary      Cria nova empresa
// @Description  Cadastra uma empresa (Matriz ou vinculada a outra) no sistema
// @Tags         Empresas
// @Accept       json
// @Produce      json
// @Param        empresa  body      domain.CreateEmpresaRequest  true  "Dados da Empresa"
// @Success      201      {object}  utils.Response{data=domain.ApiEmpresa}
// @Failure      400      {object}  utils.Response
// @Failure      401      {object}  utils.Response
// @Failure      500      {object}  utils.Response
// @Router       /empresas [post]
// @Security     ApiKeyAuth
func (h *EmpresaHandler) CreateEmpresa(c *gin.Context) {
	var req domain.CreateEmpresaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		return
	}

	userID, exists := utils.GetUserIDFromContext(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Usuário não autenticado", "")
		return
	}

	empresa, err := h.empresaService.CreateEmpresa(req, userID)
	if err != nil {
		if validator.IsValidationError(err) {
			utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusBadRequest, "Erro ao criar empresa", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Empresa criada com sucesso", empresa, nil)
}

// UpdateEmpresa atualiza uma empresa existente
// @Summary      Atualiza empresa
// @Description  Altera dados cadastrais de uma empresa pelo ID
// @Tags         Empresas
// @Accept       json
// @Produce      json
// @Param        id       path      int                           true  "ID da Empresa"
// @Param        empresa  body      domain.UpdateEmpresaRequest  true  "Novos dados da empresa"
// @Success      200      {object}  utils.Response{data=domain.ApiEmpresa}
// @Failure      400      {object}  utils.Response
// @Failure      401      {object}  utils.Response
// @Failure      404      {object}  utils.Response
// @Failure      500      {object}  utils.Response
// @Router       /empresas/{id} [put]
// @Security     ApiKeyAuth
func (h *EmpresaHandler) UpdateEmpresa(c *gin.Context) {
	id, err := path.IdFromPathParamOrSendError(c)
	if err != nil {
		return
	}

	var req domain.UpdateEmpresaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		return
	}

	userID, exists := utils.GetUserIDFromContext(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Usuário não autenticado", "")
		return
	}

	empresa, err := h.empresaService.UpdateEmpresa(id, req, userID)
	if err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Empresa não encontrada", err.Error())
		} else if validator.IsValidationError(err) {
			utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusBadRequest, "Erro ao atualizar empresa", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Empresa atualizada com sucesso", empresa, nil)
}

// DeleteEmpresa exclui uma empresa pelo ID
// @Summary      Exclui empresa
// @Description  Realiza soft delete de uma empresa (bloqueado se houver outras empresas vinculadas a ela)
// @Tags         Empresas
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID da Empresa"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      401  {object}  utils.Response
// @Failure      404  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /empresas/{id} [delete]
// @Security     ApiKeyAuth
func (h *EmpresaHandler) DeleteEmpresa(c *gin.Context) {
	id, err := path.IdFromPathParamOrSendError(c)
	if err != nil {
		return
	}

	if err := h.empresaService.DeleteEmpresa(id); err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Empresa não encontrada", err.Error())
		} else if validator.IsValidationError(err) {
			utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusBadRequest, "Erro ao excluir empresa", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Empresa excluída com sucesso", nil, nil)
}
