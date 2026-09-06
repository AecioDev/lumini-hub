package service

import (
	"lumini-hub/api.core/internal/domain"
	"lumini-hub/api.core/internal/repository"
	"lumini-hub/api.core/internal/validator"
	"lumini-hub/common/utils"
)

// CompanyFiscalConfigService gerencia operações de negócio relacionadas à
// configuração fiscal de empresas — inclusive a criptografia da senha do
// certificado digital A1, que nunca é gravada em texto puro (ver
// utils.EncryptAES/DecryptAES e config.SecurityConfig.CertificateEncryptionKey).
//
// Sem operação de exclusão de propósito: é uma configuração 1:1 intrínseca
// à Company (enquanto a empresa existe, ela sempre tem uma configuração
// fiscal associada, mesmo que vazia) — "excluir" não tem significado de
// negócio aqui, diferente de uma entidade que é uma lista de verdade (ex.:
// Company, ChartOfAccounts). Pra "limpar" o certificado, o caminho é uma
// ação dedicada (nenhuma implementada ainda) ou reenviar/sobrescrever via
// SetCertificate — nunca apagar o registro inteiro. Mesmo raciocínio se
// aplica a qualquer config 1:1 futura da empresa (ex.: CompanyVisualConfig).
type CompanyFiscalConfigService struct {
	uow           repository.UnitOfWork
	validator     *validator.CompanyFiscalConfigValidator
	encryptionKey string
}

// NewCompanyFiscalConfigService cria um novo serviço de configuração fiscal
// de empresa. encryptionKey é a chave AES-256 em hex usada pra
// criptografar/decriptografar CertificatePassword.
func NewCompanyFiscalConfigService(uow repository.UnitOfWork, encryptionKey string) *CompanyFiscalConfigService {
	return &CompanyFiscalConfigService{
		uow:           uow,
		validator:     validator.NewCompanyFiscalConfigValidator(uow.CompanyFiscalConfigs(), uow.Companies()),
		encryptionKey: encryptionKey,
	}
}

// GetByCompanyID busca a configuração fiscal de uma empresa pelo ID da empresa
func (s *CompanyFiscalConfigService) GetByCompanyID(companyID uint) (*domain.ApiCompanyFiscalConfigDetail, error) {
	config, err := s.uow.CompanyFiscalConfigs().FindByCompanyID(companyID)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, utils.ErrNotFound
	}

	dto := domain.ApiCompanyFiscalConfigDetailFromModel(*config)
	return &dto, nil
}

// CreateCompanyFiscalConfig cria a configuração fiscal de uma empresa sob
// transação do Unit of Work. O certificado (arquivo + senha) não entra por
// aqui — é definido depois via SetCertificate, chamado pelo endpoint de
// upload multipart (CFG-1.1.4).
func (s *CompanyFiscalConfigService) CreateCompanyFiscalConfig(req domain.CreateCompanyFiscalConfigRequest) (*domain.ApiCompanyFiscalConfig, error) {
	if err := s.validator.ValidateForCreation(req); err != nil {
		return nil, err
	}

	config := domain.CompanyFiscalConfig{
		CompanyID:          req.CompanyID,
		TaxRegime:          req.TaxRegime,
		AccountantName:     req.AccountantName,
		AccountantDocument: req.AccountantDocument,
		AccountantContact:  req.AccountantContact,
	}

	err := s.uow.Execute(func(uow repository.UnitOfWork) error {
		return uow.CompanyFiscalConfigs().Create(&config)
	})
	if err != nil {
		return nil, err
	}

	dto := domain.ApiCompanyFiscalConfigFromModel(config)
	return &dto, nil
}

// UpdateCompanyFiscalConfig atualiza a configuração fiscal de uma empresa sob transação do Unit of Work
func (s *CompanyFiscalConfigService) UpdateCompanyFiscalConfig(id uint, req domain.UpdateCompanyFiscalConfigRequest) (*domain.ApiCompanyFiscalConfig, error) {
	if err := s.validator.ValidateForUpdate(id, req); err != nil {
		return nil, err
	}

	config, err := s.uow.CompanyFiscalConfigs().FindByID(id)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, utils.ErrNotFound
	}

	config.TaxRegime = req.TaxRegime
	config.AccountantName = req.AccountantName
	config.AccountantDocument = req.AccountantDocument
	config.AccountantContact = req.AccountantContact

	err = s.uow.Execute(func(uow repository.UnitOfWork) error {
		return uow.CompanyFiscalConfigs().Update(config)
	})
	if err != nil {
		return nil, err
	}

	dto := domain.ApiCompanyFiscalConfigFromModel(*config)
	return &dto, nil
}

// SetCertificate grava o certificado digital (arquivo + senha) de uma
// configuração fiscal já existente, criptografando a senha antes de
// persistir. Pensado pra ser chamado pelo endpoint de upload multipart
// (CFG-1.1.4), separado do CRUD principal — nunca decripta pra devolver a
// senha, só grava.
//
// O certificado (.pfx) é decodificado (utils.ParsePKCS12Certificate) pra
// validar a senha e extrair vencimento e identidade (nome/documento) do
// titular — vencimento nunca é digitado pelo usuário, sempre vem do próprio
// arquivo. Certificado vencido é aceito normalmente (só informativo, quem
// avisa é a UI); a comparação do documento extraído com o CNPJ da Company
// também é só pra exibição, feita no frontend — não bloqueia o upload.
func (s *CompanyFiscalConfigService) SetCertificate(id uint, fileBytes []byte, password string) (*domain.ApiCompanyFiscalConfig, error) {
	config, err := s.uow.CompanyFiscalConfigs().FindByID(id)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, utils.ErrNotFound
	}

	certInfo, err := utils.ParsePKCS12Certificate(fileBytes, password)
	if err != nil {
		var validationErrors validator.ValidationErrors
		validationErrors.AddError("password", err.Error())
		return nil, validationErrors
	}

	encryptedPassword := ""
	if password != "" {
		encryptedPassword, err = utils.EncryptAES(password, s.encryptionKey)
		if err != nil {
			return nil, err
		}
	}

	config.CertificateFile = fileBytes
	config.CertificatePassword = encryptedPassword
	config.CertificateExpiry = &certInfo.NotAfter
	config.CertificateSubjectName = certInfo.SubjectName
	config.CertificateSubjectDocument = certInfo.SubjectDocument

	err = s.uow.Execute(func(uow repository.UnitOfWork) error {
		return uow.CompanyFiscalConfigs().Update(config)
	})
	if err != nil {
		return nil, err
	}

	dto := domain.ApiCompanyFiscalConfigFromModel(*config)
	return &dto, nil
}
