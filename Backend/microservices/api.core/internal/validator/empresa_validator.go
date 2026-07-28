package validator

import (
	"lumini-hub/api.core/internal/domain"
	"lumini-hub/api.core/internal/repository"
	"lumini-hub/common/utils"
)

// EmpresaValidator valida regras de negócio relacionadas a empresas
type EmpresaValidator struct {
	empresaRepo repository.EmpresaRepository
}

// NewEmpresaValidator cria um novo validador de empresas
func NewEmpresaValidator(empresaRepo repository.EmpresaRepository) *EmpresaValidator {
	return &EmpresaValidator{
		empresaRepo: empresaRepo,
	}
}

// validateParent garante que, se informado, o ParentID aponta pra uma
// empresa existente e não forma ciclo (a própria empresa não pode aparecer
// na cadeia de ancestrais dela mesma) — protege qualquer resolução futura
// de hierarquia (ex.: permission empresa.hierarquia.view) de entrar em loop
// infinito. selfID é 0 na criação (ainda não existe).
func (v *EmpresaValidator) validateParent(parentID *uint, selfID uint) error {
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
		parent, err := v.empresaRepo.FindByID(current)
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
func (v *EmpresaValidator) ValidateForCreation(req domain.CreateEmpresaRequest) error {
	var errors ValidationErrors

	cnpj := utils.RemoveMask(req.CNPJ)
	exists, err := v.empresaRepo.ExistsByCNPJ(cnpj)
	if err != nil {
		return err
	}
	if exists {
		errors.AddError("cnpj", "CNPJ já está em uso")
	}

	if req.ParentID == nil {
		hasRoot, err := v.empresaRepo.ExistsRoot()
		if err != nil {
			return err
		}
		if hasRoot {
			errors.AddError("parent_id", "já existe uma empresa Matriz cadastrada — informe a empresa à qual esta se vincula")
		}
	} else if err := v.validateParent(req.ParentID, 0); err != nil {
		if ve, ok := err.(ValidationErrors); ok {
			errors = append(errors, ve...)
		} else {
			return err
		}
	}

	if errors.HasErrors() {
		return errors
	}
	return nil
}

// ValidateForUpdate valida os dados para atualização de uma empresa
func (v *EmpresaValidator) ValidateForUpdate(id uint, req domain.UpdateEmpresaRequest) error {
	var errors ValidationErrors

	empresa, err := v.empresaRepo.FindByID(id)
	if err != nil {
		return err
	}
	if empresa == nil {
		errors.AddError("id", "empresa não encontrada")
		return errors
	}

	cnpj := utils.RemoveMask(req.CNPJ)
	if cnpj != empresa.CNPJ {
		exists, err := v.empresaRepo.ExistsByCNPJExcept(cnpj, id)
		if err != nil {
			return err
		}
		if exists {
			errors.AddError("cnpj", "CNPJ já está em uso")
		}
	}

	if req.ParentID == nil {
		hasRoot, err := v.empresaRepo.ExistsRootExcept(id)
		if err != nil {
			return err
		}
		if hasRoot {
			errors.AddError("parent_id", "já existe uma empresa Matriz cadastrada — informe a empresa à qual esta se vincula")
		}
	} else if err := v.validateParent(req.ParentID, id); err != nil {
		if ve, ok := err.(ValidationErrors); ok {
			errors = append(errors, ve...)
		} else {
			return err
		}
	}

	if errors.HasErrors() {
		return errors
	}
	return nil
}

// ValidateForDeletion valida se uma empresa pode ser excluída
func (v *EmpresaValidator) ValidateForDeletion(id uint) error {
	var errors ValidationErrors

	empresa, err := v.empresaRepo.FindByID(id)
	if err != nil {
		return err
	}
	if empresa == nil {
		errors.AddError("id", "empresa não encontrada")
		return errors
	}

	count, err := v.empresaRepo.CountByParentID(id)
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
