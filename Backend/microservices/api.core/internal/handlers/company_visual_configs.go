package handlers

import (
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

// CompanyVisualConfigHandler gerencia as requisições relacionadas à configuração visual de empresas
type CompanyVisualConfigHandler struct {
	visualConfigService *service.CompanyVisualConfigService
}

// NewCompanyVisualConfigHandler cria um novo handler de configuração visual de empresa
func NewCompanyVisualConfigHandler(uow repository.UnitOfWork) *CompanyVisualConfigHandler {
	return &CompanyVisualConfigHandler{
		visualConfigService: service.NewCompanyVisualConfigService(uow),
	}
}

// GetCompanyVisualConfigByCompany retorna a configuração visual de uma empresa
// @Summary      Busca configuração visual por empresa
// @Description  Retorna a configuração visual (logo, paleta de cores) de uma empresa pelo ID dela
// @Tags         CompanyVisualConfig
// @Accept       json
// @Produce      json
// @Param        companyId  path      int  true  "ID da Empresa"
// @Success      200        {object}  utils.Response{data=domain.ApiCompanyVisualConfigDetail}
// @Failure      401        {object}  utils.Response
// @Failure      404        {object}  utils.Response
// @Failure      500        {object}  utils.Response
// @Router       /company-visual-configs/by-company/{companyId} [get]
// @Security     ApiKeyAuth
func (h *CompanyVisualConfigHandler) GetCompanyVisualConfigByCompany(c *gin.Context) {
	companyID, err := path.UintFromPathParam(c, "companyId")
	if err != nil {
		utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		return
	}

	config, err := h.visualConfigService.GetByCompanyID(companyID)
	if err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Configuração visual não encontrada", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao buscar configuração visual", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Configuração visual encontrada", config, nil)
}

// CreateCompanyVisualConfig cria a configuração visual de uma empresa
// @Summary      Cria configuração visual
// @Description  Cadastra a paleta de cores de uma empresa. O logo é enviado separadamente via multipart.
// @Tags         CompanyVisualConfig
// @Accept       json
// @Produce      json
// @Param        config  body      domain.CreateCompanyVisualConfigRequest  true  "Dados da Configuração Visual"
// @Success      201     {object}  utils.Response{data=domain.ApiCompanyVisualConfig}
// @Failure      400     {object}  utils.Response
// @Failure      401     {object}  utils.Response
// @Failure      500     {object}  utils.Response
// @Router       /company-visual-configs [post]
// @Security     ApiKeyAuth
func (h *CompanyVisualConfigHandler) CreateCompanyVisualConfig(c *gin.Context) {
	var req domain.CreateCompanyVisualConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		return
	}

	config, err := h.visualConfigService.CreateCompanyVisualConfig(req)
	if err != nil {
		if validator.IsValidationError(err) {
			utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusBadRequest, "Erro ao criar configuração visual", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Configuração visual criada com sucesso", config, nil)
}

// UpdateCompanyVisualConfig atualiza a configuração visual de uma empresa
// @Summary      Atualiza configuração visual
// @Description  Altera a paleta de cores de uma configuração visual pelo ID (não mexe no logo)
// @Tags         CompanyVisualConfig
// @Accept       json
// @Produce      json
// @Param        id      path      int                                       true  "ID da Configuração Visual"
// @Param        config  body      domain.UpdateCompanyVisualConfigRequest  true  "Novos dados"
// @Success      200     {object}  utils.Response{data=domain.ApiCompanyVisualConfig}
// @Failure      400     {object}  utils.Response
// @Failure      401     {object}  utils.Response
// @Failure      404     {object}  utils.Response
// @Failure      500     {object}  utils.Response
// @Router       /company-visual-configs/{id} [put]
// @Security     ApiKeyAuth
func (h *CompanyVisualConfigHandler) UpdateCompanyVisualConfig(c *gin.Context) {
	id, err := path.IdFromPathParamOrSendError(c)
	if err != nil {
		return
	}

	var req domain.UpdateCompanyVisualConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		return
	}

	config, err := h.visualConfigService.UpdateCompanyVisualConfig(id, req)
	if err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Configuração visual não encontrada", err.Error())
		} else if validator.IsValidationError(err) {
			utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusBadRequest, "Erro ao atualizar configuração visual", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Configuração visual atualizada com sucesso", config, nil)
}

// Sem endpoint de exclusão de propósito — ver comentário em
// CompanyVisualConfigService sobre configs 1:1 não terem Delete.

// UploadLogo grava o logo de uma configuração visual já existente
// @Summary      Envia logo da empresa
// @Description  Upload multipart do logo (PNG, JPEG ou SVG), separado do PUT de cores
// @Tags         CompanyVisualConfig
// @Accept       multipart/form-data
// @Produce      json
// @Param        id    path      int   true  "ID da Configuração Visual"
// @Param        logo  formData  file  true  "Arquivo do logo (PNG, JPEG ou SVG)"
// @Success      200   {object}  utils.Response{data=domain.ApiCompanyVisualConfig}
// @Failure      400   {object}  utils.Response
// @Failure      401   {object}  utils.Response
// @Failure      404   {object}  utils.Response
// @Failure      500   {object}  utils.Response
// @Router       /company-visual-configs/{id}/logo [post]
// @Security     ApiKeyAuth
func (h *CompanyVisualConfigHandler) UploadLogo(c *gin.Context) {
	id, err := path.IdFromPathParamOrSendError(c)
	if err != nil {
		return
	}

	fileHeader, err := c.FormFile("logo")
	if err != nil {
		utils.ValidationErrorResponse(c, "Dados inválidos", "arquivo do logo (\"logo\") é obrigatório")
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

	config, err := h.visualConfigService.SetLogo(id, fileBytes, mimeType)
	if err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Configuração visual não encontrada", err.Error())
		} else if validator.IsValidationError(err) {
			utils.ValidationErrorResponse(c, "Não foi possível processar o logo", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusBadRequest, "Erro ao salvar logo", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Logo salvo com sucesso", config, nil)
}

// ClearLogo remove o logo de uma configuração visual já existente
// @Summary      Remove o logo da empresa
// @Description  Limpa o logo de uma configuração visual, mantendo o resto da config (cores) intacto — não é o DELETE do registro inteiro
// @Tags         CompanyVisualConfig
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID da Configuração Visual"
// @Success      200 {object}  utils.Response{data=domain.ApiCompanyVisualConfig}
// @Failure      401 {object}  utils.Response
// @Failure      404 {object}  utils.Response
// @Failure      500 {object}  utils.Response
// @Router       /company-visual-configs/{id}/logo [delete]
// @Security     ApiKeyAuth
func (h *CompanyVisualConfigHandler) ClearLogo(c *gin.Context) {
	id, err := path.IdFromPathParamOrSendError(c)
	if err != nil {
		return
	}

	config, err := h.visualConfigService.ClearLogo(id)
	if err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Configuração visual não encontrada", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao remover logo", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Logo removido com sucesso", config, nil)
}

// GetLogo serve os bytes crus do logo de uma configuração visual
// @Summary      Serve o arquivo do logo
// @Description  Devolve o binário do logo (não passa pelo envelope utils.Response — é servido com o Content-Type do próprio arquivo, pra ser usado direto num <img src>)
// @Tags         CompanyVisualConfig
// @Produce      png,jpeg,octet-stream
// @Param        id  path  int  true  "ID da Configuração Visual"
// @Success      200 {file}  binary
// @Failure      401 {object}  utils.Response
// @Failure      404 {object}  utils.Response
// @Router       /company-visual-configs/{id}/logo [get]
// @Security     ApiKeyAuth
func (h *CompanyVisualConfigHandler) GetLogo(c *gin.Context) {
	id, err := path.IdFromPathParamOrSendError(c)
	if err != nil {
		return
	}

	fileBytes, mimeType, err := h.visualConfigService.GetLogo(id)
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
