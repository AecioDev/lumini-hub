package service

import (
	"lumini-hub/api.core/internal/domain"
	"lumini-hub/api.core/internal/repository"
	"lumini-hub/api.core/internal/validator"
	"lumini-hub/common/utils"
	"strings"
)

// CustomerService gerencia operações relacionadas a clientes
type CustomerService struct {
	customerRepo repository.CustomerRepository
	validator    *validator.CustomerValidator
}

// NewCustomerService cria um novo serviço de clientes
func NewCustomerService(customerRepo repository.CustomerRepository) *CustomerService {
	return &CustomerService{
		customerRepo: customerRepo,
		validator:    validator.NewCustomerValidator(customerRepo),
	}
}

// GetCustomers retorna uma lista paginada de clientes
func (s *CustomerService) GetCustomers(pagination *utils.Pagination) (*domain.ApiCustomerListPaginated, error) {
	customers, err := s.customerRepo.FindAll(pagination)
	if err != nil {
		return nil, err
	}

	// Converter para DTOs
	customerDTOs := make([]domain.ApiCustomer, 0, len(customers))
	for _, customer := range customers {
		customerDTOs = append(customerDTOs, domain.ApiCustomerFromModel(customer))
	}

	return &domain.ApiCustomerListPaginated{
		Customer:   customerDTOs,
		Pagination: *utils.ApiPaginationFromModel(pagination),
	}, nil
}

// GetCustomerByID busca um cliente pelo ID
func (s *CustomerService) GetCustomerByID(id uint) (*domain.ApiCustomerDetail, error) {
	customer, err := s.customerRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, utils.ErrNotFound
	}

	// Converter para DTO
	customerDetailDTO := domain.ApiCustomerDetailFromModel(*customer)
	return &customerDetailDTO, nil
}

// CreateCustomer cria um novo cliente
func (s *CustomerService) CreateCustomer(req domain.CreateCustomerRequest, userID uint) (*domain.ApiCustomer, error) {
	// Validar dados
	if err := s.validator.ValidateForCreation(req); err != nil {
		return nil, err
	}

	// Limpar CPF/CNPJ
	document := strings.NewReplacer(".", "", "-", "", "/", "").Replace(req.DocumentNumber)

	// Criar cliente
	customer := domain.Customer{
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		PersonType:     req.PersonType,
		DocumentNumber: document,
		CompanyName:    req.CompanyName,
		Notes:          req.Notes,
		CreatedByID:    &userID,
		IsActive:       true, // padrão
	}

	// Persistir
	if err := s.customerRepo.Create(&customer); err != nil {
		return nil, err
	}

	// Retornar DTO
	customerDTO := domain.ApiCustomerFromModel(customer)
	return &customerDTO, nil
}

// UpdateCustomer atualiza um cliente existente
func (s *CustomerService) UpdateCustomer(id uint, req domain.UpdateCustomerRequest) (*domain.ApiCustomer, error) {
	// Validar dados
	if err := s.validator.ValidateForUpdate(id, req); err != nil {
		return nil, err
	}

	// Buscar cliente
	customer, err := s.customerRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, utils.ErrNotFound
	}

	// Atualizar campos básicos
	customer.FirstName = req.FirstName
	customer.LastName = req.LastName
	customer.CompanyName = req.CompanyName
	customer.IsActive = req.IsActive
	customer.Notes = req.Notes

	// Atualizar e formatar o DocumentNumber (removendo caracteres especiais)
	document := strings.NewReplacer(".", "", "-", "", "/", "").Replace(req.DocumentNumber)
	customer.DocumentNumber = document

	// Salvar alterações
	if err := s.customerRepo.Update(customer); err != nil {
		return nil, err
	}

	// Converter para DTO
	customerDTO := domain.ApiCustomerFromModel(*customer)
	return &customerDTO, nil
}

// DeleteCustomer exclui um cliente (soft delete)
func (s *CustomerService) DeleteCustomer(id uint) error {
	// Verificar se o cliente existe
	customer, err := s.customerRepo.FindByID(id)
	if err != nil {
		return err
	}
	if customer == nil {
		return utils.ErrNotFound
	}

	// Excluir cliente
	return s.customerRepo.Delete(id)
}
