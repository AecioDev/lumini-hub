package validator

import (
	"regexp"

	"lumini-hub/api.core/internal/domain"
	"lumini-hub/api.core/internal/repository"
)

// hexColorPattern aceita só o formato #RRGGBB (6 dígitos hex) — sem canal
// alpha, sem forma curta #RGB, pra bater exatamente com o exemplo fechado em
// plano_empresa.md (`#2563EB`) e não precisar normalizar variações no
// consumidor (ConfigProvider do antd).
var hexColorPattern = regexp.MustCompile(`^#[0-9A-Fa-f]{6}$`)

// CompanyVisualConfigValidator valida regras de negócio relacionadas à configuração visual de empresas
type CompanyVisualConfigValidator struct {
	visualConfigRepo repository.CompanyVisualConfigRepository
	companyRepo      repository.CompanyRepository
}

// NewCompanyVisualConfigValidator cria um novo validador de configuração visual de empresa
func NewCompanyVisualConfigValidator(
	visualConfigRepo repository.CompanyVisualConfigRepository,
	companyRepo repository.CompanyRepository,
) *CompanyVisualConfigValidator {
	return &CompanyVisualConfigValidator{
		visualConfigRepo: visualConfigRepo,
		companyRepo:      companyRepo,
	}
}

// validateColors valida o formato hex das cores — SecondaryColor/AccentColor/
// TextColor só são validadas se enviadas (campos opcionais). textColor vazio
// é aceito também por quem não tem esse campo (ex.: CompanyColorPalette).
func validateColors(errors *ValidationErrors, primaryColor, secondaryColor, accentColor, textColor string) {
	if !hexColorPattern.MatchString(primaryColor) {
		errors.AddError("primary_color", "cor primária deve estar no formato hexadecimal #RRGGBB")
	}
	if secondaryColor != "" && !hexColorPattern.MatchString(secondaryColor) {
		errors.AddError("secondary_color", "cor secundária deve estar no formato hexadecimal #RRGGBB")
	}
	if accentColor != "" && !hexColorPattern.MatchString(accentColor) {
		errors.AddError("accent_color", "cor de destaque deve estar no formato hexadecimal #RRGGBB")
	}
	if textColor != "" && !hexColorPattern.MatchString(textColor) {
		errors.AddError("text_color", "cor do texto deve estar no formato hexadecimal #RRGGBB")
	}
}

// ValidateForCreation valida os dados para criar a configuração visual de uma empresa
func (v *CompanyVisualConfigValidator) ValidateForCreation(req domain.CreateCompanyVisualConfigRequest) error {
	var errors ValidationErrors

	company, err := v.companyRepo.FindByID(req.CompanyID)
	if err != nil {
		return err
	}
	if company == nil {
		errors.AddError("company_id", "empresa não encontrada")
		return errors
	}

	exists, err := v.visualConfigRepo.ExistsByCompanyID(req.CompanyID)
	if err != nil {
		return err
	}
	if exists {
		errors.AddError("company_id", "esta empresa já possui configuração visual cadastrada")
	}

	validateColors(&errors, req.PrimaryColor, req.SecondaryColor, req.AccentColor, req.TextColor)

	if errors.HasErrors() {
		return errors
	}
	return nil
}

// ValidateForUpdate valida os dados para atualizar a configuração visual de uma empresa
func (v *CompanyVisualConfigValidator) ValidateForUpdate(id uint, req domain.UpdateCompanyVisualConfigRequest) error {
	var errors ValidationErrors

	config, err := v.visualConfigRepo.FindByID(id)
	if err != nil {
		return err
	}
	if config == nil {
		errors.AddError("id", "configuração visual não encontrada")
		return errors
	}

	validateColors(&errors, req.PrimaryColor, req.SecondaryColor, req.AccentColor, req.TextColor)

	if errors.HasErrors() {
		return errors
	}
	return nil
}
