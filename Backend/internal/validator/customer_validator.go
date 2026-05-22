package validator

import (
	"simple-erp-service/internal/data-structure/models"
	"simple-erp-service/internal/repository"
	"simple-erp-service/internal/utils"
)

// CustomerValidator valida regras de negócio relacionadas a clientes
type CustomerValidator struct {
	customerRepo repository.CustomerRepository
}

// NewCustomerValidator cria um novo validador de clientes
func NewCustomerValidator(customerRepo repository.CustomerRepository) *CustomerValidator {
	return &CustomerValidator{
		customerRepo: customerRepo,
	}
}

// ValidateForCreation valida os dados para criação de um cliente
func (v *CustomerValidator) ValidateForCreation(req models.CreateCustomerRequest) error {
	var errors ValidationErrors

	// Verificar se o documento já existe
	if req.DocumentNumber != "" {
		// Remover caracteres não numéricos para comparação
		document := utils.RemoveMask(req.DocumentNumber)

		exists, err := v.customerRepo.ExistsByDocument(document)
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

// ValidateForUpdate valida os dados para atualização de um cliente
func (v *CustomerValidator) ValidateForUpdate(id uint, req models.UpdateCustomerRequest) error {
	var errors ValidationErrors

	// Verificar se o cliente existe
	customer, err := v.customerRepo.FindByID(id)
	if err != nil {
		return err
	}
	if customer == nil {
		errors.AddError("id", "cliente não encontrado")
		return errors
	}

	// Verificar se o documento já está em uso por outro cliente (se fornecido)
	if req.DocumentNumber != "" && req.DocumentNumber != customer.DocumentNumber {
		// Remover caracteres não numéricos para comparação
		document := utils.RemoveMask(req.DocumentNumber)

		exists, err := v.customerRepo.ExistsByDocumentExcept(document, id)
		if err != nil {
			return err
		}
		if exists {
			errors.AddError("document_number", "documento já está em uso")
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
