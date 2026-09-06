package validator

import (
	"lumini-hub/api.core/internal/domain"
	"lumini-hub/api.core/internal/repository"
	"lumini-hub/common/utils"
)

// CompanyValidator valida regras de negócio relacionadas a empresas
type CompanyValidator struct {
	companyRepo repository.CompanyRepository
}

// NewCompanyValidator cria um novo validador de empresas
func NewCompanyValidator(companyRepo repository.CompanyRepository) *CompanyValidator {
	return &CompanyValidator{
		companyRepo: companyRepo,
	}
}

// validateParent garante que, se informado, o ParentID aponta pra uma
// empresa existente e não forma ciclo (a própria empresa não pode aparecer
// na cadeia de ancestrais dela mesma) — protege qualquer resolução futura
// de hierarquia (ex.: permission empresa.hierarquia.view) de entrar em loop
// infinito. selfID é 0 na criação (ainda não existe).
func (v *CompanyValidator) validateParent(parentID *uint, selfID uint) error {
	if parentID == nil {
		return nil
	}

	if selfID != 0 && *parentID == selfID {
		var errs ValidationErrors
		errs.AddError("parent_id", "uma empresa não pode ser vinculada a ela mesma")
		return errs
	}

	current := *parentID
	for {
		parent, err := v.companyRepo.FindByID(current)
		if err != nil {
			return err
		}
		if parent == nil {
			var errs ValidationErrors
			errs.AddError("parent_id", "empresa vinculada não encontrada")
			return errs
		}
		if selfID != 0 && parent.ID == selfID {
			var errs ValidationErrors
			errs.AddError("parent_id", "vínculo formaria um ciclo entre empresas")
			return errs
		}
		if parent.ParentID == nil {
			return nil
		}
		current = *parent.ParentID
	}
}

// ValidateForCreation valida os dados para criação de uma empresa
func (v *CompanyValidator) ValidateForCreation(req domain.CreateCompanyRequest) error {
	var errors ValidationErrors

	taxID := utils.RemoveMask(req.TaxID)
	exists, err := v.companyRepo.ExistsByTaxID(taxID)
	if err != nil {
		return err
	}
	if exists {
		errors.AddError("tax_id", "CNPJ já está em uso")
	}

	// ParentID nulo = esta empresa é (mais) uma Matriz — "Matriz" aqui
	// significa apenas "sem empresa pai", não "a única raiz do tenant"
	// (decidido 2026-09-05: o tenant pode ter vários grupos empresariais
	// independentes, sem relação de holding entre si; ver plano_empresa.md).
	if req.ParentID != nil {
		if err := v.validateParent(req.ParentID, 0); err != nil {
			if ve, ok := err.(ValidationErrors); ok {
				errors = append(errors, ve...)
			} else {
				return err
			}
		}
	}

	if errors.HasErrors() {
		return errors
	}
	return nil
}

// ValidateForUpdate valida os dados para atualização de uma empresa
func (v *CompanyValidator) ValidateForUpdate(id uint, req domain.UpdateCompanyRequest) error {
	var errors ValidationErrors

	company, err := v.companyRepo.FindByID(id)
	if err != nil {
		return err
	}
	if company == nil {
		errors.AddError("id", "empresa não encontrada")
		return errors
	}

	taxID := utils.RemoveMask(req.TaxID)
	if taxID != company.TaxID {
		exists, err := v.companyRepo.ExistsByTaxIDExcept(taxID, id)
		if err != nil {
			return err
		}
		if exists {
			errors.AddError("tax_id", "CNPJ já está em uso")
		}
	}

	if req.ParentID != nil {
		if err := v.validateParent(req.ParentID, id); err != nil {
			if ve, ok := err.(ValidationErrors); ok {
				errors = append(errors, ve...)
			} else {
				return err
			}
		}
	}

	if errors.HasErrors() {
		return errors
	}
	return nil
}

// ValidateForDeletion valida se uma empresa pode ser excluída
func (v *CompanyValidator) ValidateForDeletion(id uint) error {
	var errors ValidationErrors

	company, err := v.companyRepo.FindByID(id)
	if err != nil {
		return err
	}
	if company == nil {
		errors.AddError("id", "empresa não encontrada")
		return errors
	}

	count, err := v.companyRepo.CountByParentID(id)
	if err != nil {
		return err
	}
	if count > 0 {
		errors.AddError("id", "não é possível excluir uma empresa que tem outras empresas vinculadas a ela")
	}

	if errors.HasErrors() {
		return errors
	}
	return nil
}
