package service

import (
	"lumini-hub/api.core/internal/domain"
	"lumini-hub/api.core/internal/repository"
	"lumini-hub/api.core/internal/validator"
	"lumini-hub/common/utils"
)

// allowedLogoMimeTypes restringe o upload de logo aos formatos que fazem
// sentido pra exibir num ConfigProvider/sidebar — mesmo espírito de
// plano_empresa.md.
var allowedLogoMimeTypes = map[string]bool{
	"image/png":     true,
	"image/svg+xml": true,
	"image/jpeg":    true,
}

// CompanyVisualConfigService gerencia operações de negócio relacionadas à
// configuração visual de empresas (logo + paleta de cores).
//
// Sem operação de exclusão de propósito: é uma configuração 1:1 intrínseca
// à Company, mesmo raciocínio já aplicado em CompanyFiscalConfigService —
// ver Documentos/Planejamento/Modulo_0_Configuracao/plano_empresa.md
// § "Configs 1:1 não têm Delete".
type CompanyVisualConfigService struct {
	uow       repository.UnitOfWork
	validator *validator.CompanyVisualConfigValidator
}

// NewCompanyVisualConfigService cria um novo serviço de configuração visual de empresa
func NewCompanyVisualConfigService(uow repository.UnitOfWork) *CompanyVisualConfigService {
	return &CompanyVisualConfigService{
		uow:       uow,
		validator: validator.NewCompanyVisualConfigValidator(uow.CompanyVisualConfigs(), uow.Companies()),
	}
}

// GetByCompanyID busca a configuração visual de uma empresa pelo ID da empresa
func (s *CompanyVisualConfigService) GetByCompanyID(companyID uint) (*domain.ApiCompanyVisualConfigDetail, error) {
	config, err := s.uow.CompanyVisualConfigs().FindByCompanyID(companyID)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, utils.ErrNotFound
	}

	dto := domain.ApiCompanyVisualConfigFromModel(*config)
	return &dto, nil
}

// CreateCompanyVisualConfig cria a configuração visual de uma empresa sob
// transação do Unit of Work. O logo não entra por aqui — é definido depois
// via SetLogo, chamado pelo endpoint de upload multipart.
func (s *CompanyVisualConfigService) CreateCompanyVisualConfig(req domain.CreateCompanyVisualConfigRequest) (*domain.ApiCompanyVisualConfig, error) {
	if err := s.validator.ValidateForCreation(req); err != nil {
		return nil, err
	}

	config := domain.CompanyVisualConfig{
		CompanyID:      req.CompanyID,
		PrimaryColor:   req.PrimaryColor,
		SecondaryColor: nilIfEmpty(req.SecondaryColor),
		AccentColor:    nilIfEmpty(req.AccentColor),
	}

	err := s.uow.Execute(func(uow repository.UnitOfWork) error {
		return uow.CompanyVisualConfigs().Create(&config)
	})
	if err != nil {
		return nil, err
	}

	dto := domain.ApiCompanyVisualConfigFromModel(config)
	return &dto, nil
}

// UpdateCompanyVisualConfig atualiza a paleta de cores de uma empresa sob transação do Unit of Work
func (s *CompanyVisualConfigService) UpdateCompanyVisualConfig(id uint, req domain.UpdateCompanyVisualConfigRequest) (*domain.ApiCompanyVisualConfig, error) {
	if err := s.validator.ValidateForUpdate(id, req); err != nil {
		return nil, err
	}

	config, err := s.uow.CompanyVisualConfigs().FindByID(id)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, utils.ErrNotFound
	}

	config.PrimaryColor = req.PrimaryColor
	config.SecondaryColor = nilIfEmpty(req.SecondaryColor)
	config.AccentColor = nilIfEmpty(req.AccentColor)

	err = s.uow.Execute(func(uow repository.UnitOfWork) error {
		return uow.CompanyVisualConfigs().Update(config)
	})
	if err != nil {
		return nil, err
	}

	dto := domain.ApiCompanyVisualConfigFromModel(*config)
	return &dto, nil
}

// SetLogo grava o logo (arquivo + mimetype) de uma configuração visual já
// existente. Pensado pra ser chamado pelo endpoint de upload multipart,
// separado do CRUD principal — mesmo padrão de SetCertificate em
// CompanyFiscalConfigService.
func (s *CompanyVisualConfigService) SetLogo(id uint, fileBytes []byte, mimeType string) (*domain.ApiCompanyVisualConfig, error) {
	config, err := s.uow.CompanyVisualConfigs().FindByID(id)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, utils.ErrNotFound
	}

	if !allowedLogoMimeTypes[mimeType] {
		var validationErrors validator.ValidationErrors
		validationErrors.AddError("logo", "formato de arquivo não suportado — envie PNG, JPEG ou SVG")
		return nil, validationErrors
	}

	config.LogoFile = fileBytes
	config.LogoMimeType = mimeType

	err = s.uow.Execute(func(uow repository.UnitOfWork) error {
		return uow.CompanyVisualConfigs().Update(config)
	})
	if err != nil {
		return nil, err
	}

	dto := domain.ApiCompanyVisualConfigFromModel(*config)
	return &dto, nil
}

// GetLogo devolve os bytes crus do logo + mimetype de uma configuração
// visual, pra o handler servir a imagem diretamente (fora do envelope
// utils.Response padrão — é um binário, não JSON).
func (s *CompanyVisualConfigService) GetLogo(id uint) (fileBytes []byte, mimeType string, err error) {
	config, err := s.uow.CompanyVisualConfigs().FindByID(id)
	if err != nil {
		return nil, "", err
	}
	if config == nil || len(config.LogoFile) == 0 {
		return nil, "", utils.ErrNotFound
	}
	return config.LogoFile, config.LogoMimeType, nil
}

func nilIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
