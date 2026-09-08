package service

import (
	"lumini-hub/api.core/internal/domain"
	"lumini-hub/api.core/internal/repository"
	"lumini-hub/api.core/internal/validator"
	"lumini-hub/common/utils"
)

// allowedCompanyLogoMimeTypes restringe o upload do logo da empresa aos
// formatos que fazem sentido pra exibir num cabeçalho/relatório — mesmo
// espírito do allowedLogoMimeTypes que existia em CompanyVisualConfig antes
// da reversão (CFG-7, ver plano_empresa.md § Reversão). Nome próprio (não
// reaproveita o de CompanyVisualConfigService) porque esse arquivo é
// removido em CFG-7.2.2 e este aqui precisa sobreviver sozinho.
var allowedCompanyLogoMimeTypes = map[string]bool{
	"image/png":     true,
	"image/svg+xml": true,
	"image/jpeg":    true,
}

// MaxCompanyLogoFileSize limita o upload a 2MB — decisão do usuário em
// 2026-09-07 (CFG-7.1.2), resolvendo o débito técnico sinalizado na revisão
// da CFG-1.2.2 (upload de certificado A1 sem limite de tamanho, "revisitar
// se o padrão for reaplicado em algo como o logo"). Exportada porque a
// checagem roda no handler (via fileHeader.Size, antes de ler os bytes pra
// memória — rejeitar cedo sem gastar I/O num arquivo grande demais).
const MaxCompanyLogoFileSize = 2 * 1024 * 1024 // 2MB

// CompanyService gerencia operações de negócio relacionadas a empresas
type CompanyService struct {
	uow       repository.UnitOfWork
	validator *validator.CompanyValidator
}

// NewCompanyService cria um novo serviço de empresas recebendo o Unit of Work
func NewCompanyService(uow repository.UnitOfWork) *CompanyService {
	return &CompanyService{
		uow:       uow,
		validator: validator.NewCompanyValidator(uow.Companies()),
	}
}

// normalizedCNPJ limpa a máscara do CNPJ e devolve nil se ficar vazio —
// CNPJ é *string (nullable), não string vazia, pra várias empresas sem
// CNPJ (vinculadas cuja parte fiscal fica com a Matriz) não colidirem no
// índice único (NULL nunca é igual a NULL pro Postgres).
func normalizedCNPJ(raw string) *string {
	cleaned := utils.RemoveMask(raw)
	if cleaned == "" {
		return nil
	}
	return &cleaned
}

// GetCompanies retorna as empresas visíveis pro usuário requisitante (sem
// paginação — ver nota em CompanyRepository), aplicando a regra de
// visibilidade fechada em plano_empresa.md (2026-07-21): usuário "master"
// (sem Company vinculada) vê todas; usuário vinculado a uma Company vê só a
// própria, + subsidiárias se tiver companies.hierarchy.view. Filtro feito em
// memória de propósito (mesma justificativa de não ter /filter paginado
// aqui: volume de empresas por tenant é baixo).
func (s *CompanyService) GetCompanies(userID uint) (*domain.ApiCompanyList, error) {
	companies, err := s.uow.Companies().FindAll()
	if err != nil {
		return nil, err
	}

	// DÉBITO TÉCNICO (2026-09-06): ResolveVisibleCompanyIDs mora em
	// common/utils (compartilhada entre microsserviços) e pede um *gorm.DB
	// cru, não a UnitOfWork deste serviço — por isso o GetDB() aqui, fora do
	// padrão s.uow.Execute(...) que todo Service deveria seguir. É só leitura,
	// sem transação, então não tem risco real hoje. NÃO copiar esse padrão
	// pra escrita de dados — qualquer operação que grave algo tem que passar
	// por s.uow.Execute(...) pra manter atomicidade. Solução futura: expor um
	// método próprio na interface UnitOfWork (ex. ResolveVisibleCompanyIDs)
	// que chama a função compartilhada por dentro, sem vazar *gorm.DB pro
	// Service.
	visibleIDs, unrestricted, err := utils.ResolveVisibleCompanyIDs(s.uow.Companies().GetDB(), userID)
	if err != nil {
		return nil, err
	}
	if !unrestricted {
		visible := make(map[uint]bool, len(visibleIDs))
		for _, id := range visibleIDs {
			visible[id] = true
		}
		filtered := make([]domain.Company, 0, len(companies))
		for _, company := range companies {
			if visible[company.ID] {
				filtered = append(filtered, company)
			}
		}
		companies = filtered
	}

	companyDTOs := make([]domain.ApiCompany, 0, len(companies))
	for _, company := range companies {
		companyDTOs = append(companyDTOs, domain.ApiCompanyFromModel(company))
	}

	return &domain.ApiCompanyList{Companies: companyDTOs}, nil
}

// GetCompanyByID busca uma empresa pelo ID
func (s *CompanyService) GetCompanyByID(id uint) (*domain.ApiCompanyDetail, error) {
	company, err := s.uow.Companies().FindByID(id)
	if err != nil {
		return nil, err
	}
	if company == nil {
		return nil, utils.ErrNotFound
	}

	dto := domain.ApiCompanyDetailFromModel(*company)
	return &dto, nil
}

// CreateCompany cria uma nova empresa sob transação atômica do Unit of Work
func (s *CompanyService) CreateCompany(req domain.CreateCompanyRequest, userID uint) (*domain.ApiCompany, error) {
	if err := s.validator.ValidateForCreation(req); err != nil {
		return nil, err
	}

	company := domain.Company{
		ParentID:    req.ParentID,
		LegalName:   req.LegalName,
		TradeName:   req.TradeName,
		CNPJ:        normalizedCNPJ(req.CNPJ),
		IsActive:    true,
		CreatedByID: &userID,
	}

	err := s.uow.Execute(func(uow repository.UnitOfWork) error {
		return uow.Companies().Create(&company)
	})
	if err != nil {
		return nil, err
	}

	dto := domain.ApiCompanyFromModel(company)
	return &dto, nil
}

// UpdateCompany atualiza uma empresa existente sob transação do Unit of Work
func (s *CompanyService) UpdateCompany(id uint, req domain.UpdateCompanyRequest, userID uint) (*domain.ApiCompany, error) {
	if err := s.validator.ValidateForUpdate(id, req); err != nil {
		return nil, err
	}

	company, err := s.uow.Companies().FindByID(id)
	if err != nil {
		return nil, err
	}
	if company == nil {
		return nil, utils.ErrNotFound
	}

	company.ParentID = req.ParentID
	company.LegalName = req.LegalName
	company.TradeName = req.TradeName
	company.CNPJ = normalizedCNPJ(req.CNPJ)
	company.IsActive = req.IsActive
	company.UpdatedByID = &userID

	err = s.uow.Execute(func(uow repository.UnitOfWork) error {
		return uow.Companies().Update(company)
	})
	if err != nil {
		return nil, err
	}

	dto := domain.ApiCompanyFromModel(*company)
	return &dto, nil
}

// DeleteCompany exclui uma empresa sob transação do Unit of Work
func (s *CompanyService) DeleteCompany(id uint) error {
	company, err := s.uow.Companies().FindByID(id)
	if err != nil {
		return err
	}
	if company == nil {
		return utils.ErrNotFound
	}

	if err := s.validator.ValidateForDeletion(id); err != nil {
		return err
	}

	return s.uow.Execute(func(uow repository.UnitOfWork) error {
		return uow.Companies().Delete(id)
	})
}

// SetLogo grava o logo (arquivo + mimetype) de uma empresa já existente.
// Pensado pra ser chamado pelo endpoint de upload multipart, separado do
// CRUD principal — mesmo padrão de SetCertificate em
// CompanyFiscalConfigService e do antigo SetLogo de CompanyVisualConfig
// (CFG-7.1, ver plano_empresa.md § Reversão). Validação de tamanho máximo
// fica no handler (fileHeader.Size, antes de ler os bytes) — aqui só
// mimetype, igual o padrão antigo.
func (s *CompanyService) SetLogo(id uint, fileBytes []byte, mimeType string) (*domain.ApiCompany, error) {
	company, err := s.uow.Companies().FindByID(id)
	if err != nil {
		return nil, err
	}
	if company == nil {
		return nil, utils.ErrNotFound
	}

	if !allowedCompanyLogoMimeTypes[mimeType] {
		var validationErrors validator.ValidationErrors
		validationErrors.AddError("logo", "formato de arquivo não suportado — envie PNG, JPEG ou SVG")
		return nil, validationErrors
	}

	company.LogoFile = fileBytes
	company.LogoMimeType = mimeType

	err = s.uow.Execute(func(uow repository.UnitOfWork) error {
		return uow.Companies().Update(company)
	})
	if err != nil {
		return nil, err
	}

	dto := domain.ApiCompanyFromModel(*company)
	return &dto, nil
}

// ClearLogo remove o logo de uma empresa, mantendo o resto do cadastro
// intacto — não é o DELETE da empresa como um todo.
func (s *CompanyService) ClearLogo(id uint) (*domain.ApiCompany, error) {
	company, err := s.uow.Companies().FindByID(id)
	if err != nil {
		return nil, err
	}
	if company == nil {
		return nil, utils.ErrNotFound
	}

	company.LogoFile = nil
	company.LogoMimeType = ""

	err = s.uow.Execute(func(uow repository.UnitOfWork) error {
		return uow.Companies().Update(company)
	})
	if err != nil {
		return nil, err
	}

	dto := domain.ApiCompanyFromModel(*company)
	return &dto, nil
}

// GetLogo devolve os bytes crus do logo + mimetype de uma empresa, pra o
// handler servir a imagem diretamente (fora do envelope utils.Response
// padrão — é um binário, não JSON).
func (s *CompanyService) GetLogo(id uint) (fileBytes []byte, mimeType string, err error) {
	company, err := s.uow.Companies().FindByID(id)
	if err != nil {
		return nil, "", err
	}
	if company == nil || len(company.LogoFile) == 0 {
		return nil, "", utils.ErrNotFound
	}
	return company.LogoFile, company.LogoMimeType, nil
}
