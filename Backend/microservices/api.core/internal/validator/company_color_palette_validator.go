package validator

import (
	"lumini-hub/api.core/internal/domain"
	"lumini-hub/api.core/internal/repository"
)

// CompanyColorPaletteValidator valida regras de negócio relacionadas às
// paletas de cores personalizadas de uma empresa. Reaproveita
// hexColorPattern/validateColors já definidos em
// company_visual_config_validator.go (mesmo pacote) — mesmo formato de cor
// exigido em CompanyVisualConfig.
type CompanyColorPaletteValidator struct {
	paletteRepo repository.CompanyColorPaletteRepository
	companyRepo repository.CompanyRepository
}

// NewCompanyColorPaletteValidator cria um novo validador de paletas de cores personalizadas
func NewCompanyColorPaletteValidator(
	paletteRepo repository.CompanyColorPaletteRepository,
	companyRepo repository.CompanyRepository,
) *CompanyColorPaletteValidator {
	return &CompanyColorPaletteValidator{
		paletteRepo: paletteRepo,
		companyRepo: companyRepo,
	}
}

// ValidateForCreation valida os dados para salvar uma nova paleta personalizada
func (v *CompanyColorPaletteValidator) ValidateForCreation(req domain.CreateCompanyColorPaletteRequest) error {
	var errors ValidationErrors

	company, err := v.companyRepo.FindByID(req.CompanyID)
	if err != nil {
		return err
	}
	if company == nil {
		errors.AddError("company_id", "empresa não encontrada")
		return errors
	}

	validateColors(&errors, req.PrimaryColor, req.SecondaryColor, req.AccentColor, "")

	if errors.HasErrors() {
		return errors
	}
	return nil
}

// ValidateForDeletion valida se a paleta existe antes de excluir
func (v *CompanyColorPaletteValidator) ValidateForDeletion(id uint) error {
	palette, err := v.paletteRepo.FindByID(id)
	if err != nil {
		return err
	}
	if palette == nil {
		var errors ValidationErrors
		errors.AddError("id", "paleta não encontrada")
		return errors
	}
	return nil
}
