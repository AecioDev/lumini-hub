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

// CompanyFiscalConfigHandler gerencia as requisições relacionadas à configuração fiscal de empresas
type CompanyFiscalConfigHandler struct {
	fiscalConfigService *service.CompanyFiscalConfigService
}

// NewCompanyFiscalConfigHandler cria um novo handler de configuração fiscal de empresa
func NewCompanyFiscalConfigHandler(uow repository.UnitOfWork, encryptionKey string) *CompanyFiscalConfigHandler {
	return &CompanyFiscalConfigHandler{
		fiscalConfigService: service.NewCompanyFiscalConfigService(uow, encryptionKey),
	}
}

// GetCompanyFiscalConfigByCompany retorna a configuração fiscal de uma empresa
// @Summary      Busca configuração fiscal por empresa
// @Description  Retorna a configuração fiscal (certificado, tributação, contador) de uma empresa pelo ID dela
// @Tags         CompanyFiscalConfig
// @Accept       json
// @Produce      json
// @Param        companyId  path      int  true  "ID da Empresa"
// @Success      200        {object}  utils.Response{data=domain.ApiCompanyFiscalConfigDetail}
// @Failure      401        {object}  utils.Response
// @Failure      404        {object}  utils.Response
// @Failure      500        {object}  utils.Response
// @Router       /company-fiscal-configs/by-company/{companyId} [get]
// @Security     ApiKeyAuth
func (h *CompanyFiscalConfigHandler) GetCompanyFiscalConfigByCompany(c *gin.Context) {
	companyID, err := path.UintFromPathParam(c, "companyId")
	if err != nil {
		utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		return
	}

	config, err := h.fiscalConfigService.GetByCompanyID(companyID)
	if err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Configuração fiscal não encontrada", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao buscar configuração fiscal", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Configuração fiscal encontrada", config, nil)
}

// CreateCompanyFiscalConfig cria a configuração fiscal de uma empresa
// @Summary      Cria configuração fiscal
// @Description  Cadastra a configuração fiscal (tributação, contador) de uma empresa. O certificado é enviado separadamente via multipart.
// @Tags         CompanyFiscalConfig
// @Accept       json
// @Produce      json
// @Param        config  body      domain.CreateCompanyFiscalConfigRequest  true  "Dados da Configuração Fiscal"
// @Success      201     {object}  utils.Response{data=domain.ApiCompanyFiscalConfig}
// @Failure      400     {object}  utils.Response
// @Failure      401     {object}  utils.Response
// @Failure      500     {object}  utils.Response
// @Router       /company-fiscal-configs [post]
// @Security     ApiKeyAuth
func (h *CompanyFiscalConfigHandler) CreateCompanyFiscalConfig(c *gin.Context) {
	var req domain.CreateCompanyFiscalConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		return
	}

	config, err := h.fiscalConfigService.CreateCompanyFiscalConfig(req)
	if err != nil {
		if validator.IsValidationError(err) {
			utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusBadRequest, "Erro ao criar configuração fiscal", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Configuração fiscal criada com sucesso", config, nil)
}

// UpdateCompanyFiscalConfig atualiza a configuração fiscal de uma empresa
// @Summary      Atualiza configuração fiscal
// @Description  Altera tributação/dados do contador de uma configuração fiscal pelo ID (não mexe no certificado)
// @Tags         CompanyFiscalConfig
// @Accept       json
// @Produce      json
// @Param        id      path      int                                       true  "ID da Configuração Fiscal"
// @Param        config  body      domain.UpdateCompanyFiscalConfigRequest  true  "Novos dados"
// @Success      200     {object}  utils.Response{data=domain.ApiCompanyFiscalConfig}
// @Failure      400     {object}  utils.Response
// @Failure      401     {object}  utils.Response
// @Failure      404     {object}  utils.Response
// @Failure      500     {object}  utils.Response
// @Router       /company-fiscal-configs/{id} [put]
// @Security     ApiKeyAuth
func (h *CompanyFiscalConfigHandler) UpdateCompanyFiscalConfig(c *gin.Context) {
	id, err := path.IdFromPathParamOrSendError(c)
	if err != nil {
		return
	}

	var req domain.UpdateCompanyFiscalConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		return
	}

	config, err := h.fiscalConfigService.UpdateCompanyFiscalConfig(id, req)
	if err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Configuração fiscal não encontrada", err.Error())
		} else if validator.IsValidationError(err) {
			utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusBadRequest, "Erro ao atualizar configuração fiscal", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Configuração fiscal atualizada com sucesso", config, nil)
}

// Sem endpoint de exclusão de propósito — ver comentário em
// CompanyFiscalConfigService sobre configs 1:1 não terem Delete.

// UploadCertificate grava o certificado digital (arquivo + senha) de uma configuração fiscal já existente
// @Summary      Envia certificado digital
// @Description  Upload multipart do certificado A1 (.pfx) e sua senha, separado do PUT de configuração — a senha é criptografada antes de persistir, nunca devolvida em nenhuma resposta
// @Tags         CompanyFiscalConfig
// @Accept       multipart/form-data
// @Produce      json
// @Param        id           path      int     true  "ID da Configuração Fiscal"
// @Param        certificate  formData  file    true  "Arquivo do certificado (.pfx)"
// @Param        password     formData  string  true  "Senha do certificado"
// @Success      200          {object}  utils.Response{data=domain.ApiCompanyFiscalConfig}
// @Failure      400          {object}  utils.Response
// @Failure      401          {object}  utils.Response
// @Failure      404          {object}  utils.Response
// @Failure      500          {object}  utils.Response
// @Router       /company-fiscal-configs/{id}/certificate [post]
// @Security     ApiKeyAuth
func (h *CompanyFiscalConfigHandler) UploadCertificate(c *gin.Context) {
	id, err := path.IdFromPathParamOrSendError(c)
	if err != nil {
		return
	}

	fileHeader, err := c.FormFile("certificate")
	if err != nil {
		utils.ValidationErrorResponse(c, "Dados inválidos", "arquivo do certificado (\"certificate\") é obrigatório")
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao ler certificado", err.Error())
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao ler certificado", err.Error())
		return
	}

	password := c.PostForm("password")
	if password == "" {
		utils.ValidationErrorResponse(c, "Dados inválidos", "senha do certificado (\"password\") é obrigatória")
		return
	}

	config, err := h.fiscalConfigService.SetCertificate(id, fileBytes, password)
	if err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Configuração fiscal não encontrada", err.Error())
		} else if validator.IsValidationError(err) {
			utils.ValidationErrorResponse(c, "Não foi possível processar o certificado", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusBadRequest, "Erro ao salvar certificado", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Certificado salvo com sucesso", config, nil)
}
