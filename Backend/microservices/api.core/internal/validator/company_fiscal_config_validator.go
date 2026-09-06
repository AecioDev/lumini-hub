package validator

import (
	"lumini-hub/api.core/internal/domain"
	"lumini-hub/api.core/internal/repository"
)

// validTaxRegimes define os valores aceitos pra CompanyFiscalConfig.TaxRegime.
var validTaxRegimes = map[string]bool{
	"SIMPLES":   true,
	"PRESUMIDO": true,
	"REAL":      true,
}

// CompanyFiscalConfigValidator valida regras de negócio relacionadas à configuração fiscal de empresas
type CompanyFiscalConfigValidator struct {
	fiscalConfigRepo repository.CompanyFiscalConfigRepository
	companyRepo      repository.CompanyRepository
}

// NewCompanyFiscalConfigValidator cria um novo validador de configuração fiscal de empresa
func NewCompanyFiscalConfigValidator(
	fiscalConfigRepo repository.CompanyFiscalConfigRepository,
	companyRepo repository.CompanyRepository,
) *CompanyFiscalConfigValidator {
	return &CompanyFiscalConfigValidator{
		fiscalConfigRepo: fiscalConfigRepo,
		companyRepo:      companyRepo,
	}
}

// ValidateForCreation valida os dados para criar a configuração fiscal de uma empresa
func (v *CompanyFiscalConfigValidator) ValidateForCreation(req domain.CreateCompanyFiscalConfigRequest) error {
	var errors ValidationErrors

	company, err := v.companyRepo.FindByID(req.CompanyID)
	if err != nil {
		return err
	}
	if company == nil {
		errors.AddError("company_id", "empresa não encontrada")
		return errors
	}

	exists, err := v.fiscalConfigRepo.ExistsByCompanyID(req.CompanyID)
	if err != nil {
		return err
	}
	if exists {
		errors.AddError("company_id", "esta empresa já possui configuração fiscal cadastrada")
	}

	if !validTaxRegimes[req.TaxRegime] {
		errors.AddError("tax_regime", "tipo de tributação inválido")
	}

	if errors.HasErrors() {
		return errors
	}
	return nil
}

// ValidateForUpdate valida os dados para atualizar a configuração fiscal de uma empresa
func (v *CompanyFiscalConfigValidator) ValidateForUpdate(id uint, req domain.UpdateCompanyFiscalConfigRequest) error {
	var errors ValidationErrors

	config, err := v.fiscalConfigRepo.FindByID(id)
	if err != nil {
		return err
	}
	if config == nil {
		errors.AddError("id", "configuração fiscal não encontrada")
		return errors
	}

	if !validTaxRegimes[req.TaxRegime] {
		errors.AddError("tax_regime", "tipo de tributação inválido")
	}

	if errors.HasErrors() {
		return errors
	}
	return nil
}
