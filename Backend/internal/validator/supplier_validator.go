package validator

import (
	"simple-erp-service/internal/data-structure/models"
	"simple-erp-service/internal/repository"
	"simple-erp-service/internal/utils"
)

// SupplierValidator valida regras de negócio relacionadas a fornecedores
type SupplierValidator struct {
	supplierRepo repository.SupplierRepository
}

// NewSupplierValidator cria um novo validador de fornecedores
func NewSupplierValidator(supplierRepo repository.SupplierRepository) *SupplierValidator {
	return &SupplierValidator{
		supplierRepo: supplierRepo,
	}
}

// ValidateForCreation valida os dados para criação de um fornecedor
func (v *SupplierValidator) ValidateForCreation(req models.CreateSupplierRequest) error {
	var errors ValidationErrors

	// Verificar se o documento já existe
	if req.DocumentNumber != "" {
		// Remover caracteres não numéricos para comparação
		document := utils.RemoveMask(req.DocumentNumber)

		exists, err := v.supplierRepo.ExistsByDocument(document)
		if err != nil {
			return err
		}
		if exists {
			errors.AddError("document", "documento já está em uso")
		}
	}

	// Validar nome (se fornecido)
	if req.FirstName != "" && len(req.FirstName) < 3 {
		errors.AddError("first_name", "o primeiro nome deve ter pelo menos 3 caracteres")
	}
	if req.LastName != "" && len(req.LastName) < 3 {
		errors.AddError("last_name", "o sobrenome deve ter pelo menos 3 caracteres")
	}

	if errors.HasErrors() {
		return errors
	}
	return nil
}

// ValidateForUpdate valida os dados para atualização de um fornecedor
func (v *SupplierValidator) ValidateForUpdate(id uint, req models.UpdateSupplierRequest) error {
	var errors ValidationErrors

	// Verificar se o fornecedor existe
	supplier, err := v.supplierRepo.FindByID(id)
	if err != nil {
		return err
	}
	if supplier == nil {
		errors.AddError("id", "fornecedor não encontrado")
		return errors
	}

	// Verificar se o documento já está em uso por outro fornecedor (se fornecido)
	if req.DocumentNumber != "" && req.DocumentNumber != supplier.DocumentNumber {
		// Remover caracteres não numéricos para comparação
		document := utils.RemoveMask(req.DocumentNumber)

		exists, err := v.supplierRepo.ExistsByDocumentExcept(document, id)
		if err != nil {
			return err
		}
		if exists {
			errors.AddError("document", "documento já está em uso")
		}
	}

	// Validar nome (se fornecido)
	if req.FirstName != "" && len(req.FirstName) < 3 {
		errors.AddError("first_name", "o primeiro nome deve ter pelo menos 3 caracteres")
	}
	if req.LastName != "" && len(req.LastName) < 3 {
		errors.AddError("last_name", "o sobrenome deve ter pelo menos 3 caracteres")
	}

	if errors.HasErrors() {
		return errors
	}
	return nil
}
