package handlers

import (
	"fmt"
	"io"
	"net/http"

	"lumini-hub/api.core/internal/domain"
	"lumini-hub/api.core/internal/repository"
	"lumini-hub/api.core/internal/service"
	"lumini-hub/api.core/internal/validator"
	"lumini-hub/common/utils"
	"lumini-hub/common/utils/path"

	"github.com/gin-gonic/gin"
)

// CompanyHandler gerencia as requisições relacionadas a empresas
type CompanyHandler struct {
	companyService *service.CompanyService
}

// NewCompanyHandler cria um novo handler de empresas recebendo o Unit of Work
func NewCompanyHandler(uow repository.UnitOfWork) *CompanyHandler {
	return &CompanyHandler{
		companyService: service.NewCompanyService(uow),
	}
}

// GetCompanies retorna todas as empresas cadastradas
// @Summary      Lista empresas
// @Description  Retorna todas as empresas cadastradas (Matriz e vinculadas)
// @Tags         Companies
// @Accept       json
// @Produce      json
// @Success      200  {object}  utils.Response{data=domain.ApiCompanyList}
// @Failure      401  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /companies [get]
// @Security     ApiKeyAuth
func (h *CompanyHandler) GetCompanies(c *gin.Context) {
	userID, exists := utils.GetUserIDFromContext(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Usuário não autenticado", "")
		return
	}

	companies, err := h.companyService.GetCompanies(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao buscar empresas", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Empresas encontradas", companies, nil)
}

// GetCompany retorna uma empresa específica pelo ID
// @Summary      Busca empresa por ID
// @Description  Retorna os detalhes de uma única empresa cadastrada
// @Tags         Companies
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID da Empresa"
// @Success      200  {object}  utils.Response{data=domain.ApiCompanyDetail}
// @Failure      401  {object}  utils.Response
// @Failure      404  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /companies/{id} [get]
// @Security     ApiKeyAuth
func (h *CompanyHandler) GetCompany(c *gin.Context) {
	id, err := path.IdFromPathParamOrSendError(c)
	if err != nil {
		return
	}

	company, err := h.companyService.GetCompanyByID(id)
	if err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Empresa não encontrada", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao buscar empresa", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Empresa encontrada", company, nil)
}

// CreateCompany cria uma nova empresa
// @Summary      Cria nova empresa
// @Description  Cadastra uma empresa (Matriz ou vinculada a outra) no sistema
// @Tags         Companies
// @Accept       json
// @Produce      json
// @Param        company  body      domain.CreateCompanyRequest  true  "Dados da Empresa"
// @Success      201      {object}  utils.Response{data=domain.ApiCompany}
// @Failure      400      {object}  utils.Response
// @Failure      401      {object}  utils.Response
// @Failure      500      {object}  utils.Response
// @Router       /companies [post]
// @Security     ApiKeyAuth
func (h *CompanyHandler) CreateCompany(c *gin.Context) {
	var req domain.CreateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		return
	}

	userID, exists := utils.GetUserIDFromContext(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Usuário não autenticado", "")
		return
	}

	company, err := h.companyService.CreateCompany(req, userID)
	if err != nil {
		if validator.IsValidationError(err) {
			utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusBadRequest, "Erro ao criar empresa", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Empresa criada com sucesso", company, nil)
}

// UpdateCompany atualiza uma empresa existente
// @Summary      Atualiza empresa
// @Description  Altera dados cadastrais de uma empresa pelo ID
// @Tags         Companies
// @Accept       json
// @Produce      json
// @Param        id       path      int                          true  "ID da Empresa"
// @Param        company  body      domain.UpdateCompanyRequest  true  "Novos dados da empresa"
// @Success      200      {object}  utils.Response{data=domain.ApiCompany}
// @Failure      400      {object}  utils.Response
// @Failure      401      {object}  utils.Response
// @Failure      404      {object}  utils.Response
// @Failure      500      {object}  utils.Response
// @Router       /companies/{id} [put]
// @Security     ApiKeyAuth
func (h *CompanyHandler) UpdateCompany(c *gin.Context) {
	id, err := path.IdFromPathParamOrSendError(c)
	if err != nil {
		return
	}

	var req domain.UpdateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		return
	}

	userID, exists := utils.GetUserIDFromContext(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Usuário não autenticado", "")
		return
	}

	company, err := h.companyService.UpdateCompany(id, req, userID)
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

	utils.SuccessResponse(c, http.StatusOK, "Empresa atualizada com sucesso", company, nil)
}

// DeleteCompany exclui uma empresa pelo ID
// @Summary      Exclui empresa
// @Description  Realiza soft delete de uma empresa (bloqueado se houver outras empresas vinculadas a ela)
// @Tags         Companies
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID da Empresa"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      401  {object}  utils.Response
// @Failure      404  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /companies/{id} [delete]
// @Security     ApiKeyAuth
func (h *CompanyHandler) DeleteCompany(c *gin.Context) {
	id, err := path.IdFromPathParamOrSendError(c)
	if err != nil {
		return
	}

	if err := h.companyService.DeleteCompany(id); err != nil {
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

// UploadLogo grava o logo de uma empresa já existente
// @Summary      Envia logo da empresa
// @Description  Upload multipart do logo (PNG, JPEG ou SVG, até 2MB), separado do PUT de dados cadastrais
// @Tags         Companies
// @Accept       multipart/form-data
// @Produce      json
// @Param        id    path      int   true  "ID da Empresa"
// @Param        logo  formData  file  true  "Arquivo do logo (PNG, JPEG ou SVG, até 2MB)"
// @Success      200   {object}  utils.Response{data=domain.ApiCompany}
// @Failure      400   {object}  utils.Response
// @Failure      401   {object}  utils.Response
// @Failure      404   {object}  utils.Response
// @Failure      500   {object}  utils.Response
// @Router       /companies/{id}/logo [post]
// @Security     ApiKeyAuth
func (h *CompanyHandler) UploadLogo(c *gin.Context) {
	id, err := path.IdFromPathParamOrSendError(c)
	if err != nil {
		return
	}

	fileHeader, err := c.FormFile("logo")
	if err != nil {
		utils.ValidationErrorResponse(c, "Dados inválidos", "arquivo do logo (\"logo\") é obrigatório")
		return
	}

	if fileHeader.Size > service.MaxCompanyLogoFileSize {
		utils.ValidationErrorResponse(c, "Arquivo muito grande", fmt.Sprintf("o logo não pode ultrapassar %dMB", service.MaxCompanyLogoFileSize/(1024*1024)))
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao ler logo", err.Error())
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao ler logo", err.Error())
		return
	}

	mimeType := fileHeader.Header.Get("Content-Type")

	company, err := h.companyService.SetLogo(id, fileBytes, mimeType)
	if err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Empresa não encontrada", err.Error())
		} else if validator.IsValidationError(err) {
			utils.ValidationErrorResponse(c, "Não foi possível processar o logo", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusBadRequest, "Erro ao salvar logo", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Logo salvo com sucesso", company, nil)
}

// ClearLogo remove o logo de uma empresa já existente
// @Summary      Remove o logo da empresa
// @Description  Limpa o logo de uma empresa, mantendo o resto do cadastro intacto — não é o DELETE da empresa
// @Tags         Companies
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID da Empresa"
// @Success      200 {object}  utils.Response{data=domain.ApiCompany}
// @Failure      401 {object}  utils.Response
// @Failure      404 {object}  utils.Response
// @Failure      500 {object}  utils.Response
// @Router       /companies/{id}/logo [delete]
// @Security     ApiKeyAuth
func (h *CompanyHandler) ClearLogo(c *gin.Context) {
	id, err := path.IdFromPathParamOrSendError(c)
	if err != nil {
		return
	}

	company, err := h.companyService.ClearLogo(id)
	if err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Empresa não encontrada", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao remover logo", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Logo removido com sucesso", company, nil)
}

// GetLogo serve os bytes crus do logo de uma empresa
// @Summary      Serve o arquivo do logo
// @Description  Devolve o binário do logo (não passa pelo envelope utils.Response — é servido com o Content-Type do próprio arquivo, pra ser usado direto num <img src>)
// @Tags         Companies
// @Produce      png,jpeg,octet-stream
// @Param        id  path  int  true  "ID da Empresa"
// @Success      200 {file}  binary
// @Failure      401 {object}  utils.Response
// @Failure      404 {object}  utils.Response
// @Router       /companies/{id}/logo [get]
// @Security     ApiKeyAuth
func (h *CompanyHandler) GetLogo(c *gin.Context) {
	id, err := path.IdFromPathParamOrSendError(c)
	if err != nil {
		return
	}

	fileBytes, mimeType, err := h.companyService.GetLogo(id)
	if err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Logo não encontrado", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao buscar logo", err.Error())
		}
		return
	}

	c.Data(http.StatusOK, mimeType, fileBytes)
}
