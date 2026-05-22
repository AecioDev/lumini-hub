package service

import (
	dto "simple-erp-service/internal/data-structure/dto"
	"simple-erp-service/internal/data-structure/models"
	"simple-erp-service/internal/repository"
	"simple-erp-service/internal/utils"
	"simple-erp-service/internal/validator"
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
func (s *CustomerService) GetCustomers(pagination *models.Pagination) (*dto.ApiCustomerListPaginated, error) {
	customers, err := s.customerRepo.FindAll(pagination)
	if err != nil {
		return nil, err
	}

	// Converter para DTOs
	customerDTOs := make([]dto.ApiCustomer, 0, len(customers))
	for _, customer := range customers {
		customerDTOs = append(customerDTOs, dto.ApiCustomerFromModel(customer))
	}

	return &dto.ApiCustomerListPaginated{
		Customer:   customerDTOs,
		Pagination: *dto.ApiPaginationFromModel(pagination),
	}, nil
}

// GetCustomerByID busca um cliente pelo ID
func (s *CustomerService) GetCustomerByID(id uint) (*dto.ApiCustomerDetail, error) {
	customer, err := s.customerRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, utils.ErrNotFound
	}

	// Converter para DTO
	customerDetailDTO := dto.ApiCustomerDetailFromModel(*customer)
	return &customerDetailDTO, nil
}

// CreateCustomer cria um novo cliente
func (s *CustomerService) CreateCustomer(req models.CreateCustomerRequest, userID uint) (*dto.ApiCustomer, error) {
	// Validar dados
	if err := s.validator.ValidateForCreation(req); err != nil {
		return nil, err
	}

	// Limpar CPF/CNPJ
	document := strings.NewReplacer(".", "", "-", "", "/", "").Replace(req.DocumentNumber)

	// Criar cliente
	customer := models.Customer{
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
	customerDTO := dto.ApiCustomerFromModel(customer)
	return &customerDTO, nil
}

// UpdateCustomer atualiza um cliente existente
func (s *CustomerService) UpdateCustomer(id uint, req models.UpdateCustomerRequest) (*dto.ApiCustomer, error) {
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

	// TODO Implementar os Métodos
	//// Atualizar Relacionamentos (se IDs foram enviados)
	// if len(req.DocumentsIDs) > 0 {
	// 	documents, err := s.documentRepo.FindByIDs(req.DocumentsIDs)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	customer.Documents = documents
	// }

	// if len(req.AdressesIDs) > 0 {
	// 	addresses, err := s.addressRepo.FindByIDs(req.AdressesIDs)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	customer.Addresses = addresses
	// }

	// if len(req.ContactsIDs) > 0 {
	// 	contacts, err := s.contactRepo.FindByIDs(req.ContactsIDs)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	customer.Contacts = contacts
	// }

	// Salvar alterações
	if err := s.customerRepo.Update(customer); err != nil {
		return nil, err
	}

	// Converter para DTO
	customerDTO := dto.ApiCustomerFromModel(*customer)
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
