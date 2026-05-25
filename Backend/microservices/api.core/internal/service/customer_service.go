package service

import (
	"lumini-hub/api.core/internal/domain"
	"lumini-hub/api.core/internal/repository"
	"lumini-hub/api.core/internal/validator"
	"lumini-hub/common/utils"
	"strings"
)

// CustomerService gerencia operações de negócios relacionadas a clientes
type CustomerService struct {
	uow       repository.UnitOfWork
	validator *validator.CustomerValidator
}

// NewCustomerService cria um novo serviço de clientes recebendo o Unit of Work
func NewCustomerService(uow repository.UnitOfWork) *CustomerService {
	return &CustomerService{
		uow:       uow,
		validator: validator.NewCustomerValidator(uow.Customers()),
	}
}

// GetCustomers retorna uma lista paginada de clientes (legada via GET)
func (s *CustomerService) GetCustomers(pagination *utils.Pagination) (*domain.ApiCustomerListPaginated, error) {
	// Utiliza o FindByFilter genérico sem filtros para obter todos paginados
	customers, err := s.uow.Customers().FindByFilter(domain.CustomerFilterRequest{}, pagination)
	if err != nil {
		return nil, err
	}

	customerDTOs := make([]domain.ApiCustomer, 0, len(customers))
	for _, customer := range customers {
		customerDTOs = append(customerDTOs, domain.ApiCustomerFromModel(customer))
	}

	return &domain.ApiCustomerListPaginated{
		Customer:   customerDTOs,
		Pagination: *utils.ApiPaginationFromModel(pagination),
	}, nil
}

// GetCustomersByFilter retorna clientes com filtros avançados e paginação (PostFilter)
func (s *CustomerService) GetCustomersByFilter(filter domain.CustomerFilterRequest) (*domain.ApiCustomerListPaginated, error) {
	// Normalizar dados de paginação do filtro
	pageNo := 1
	if filter.PageNo != nil && *filter.PageNo > 0 {
		pageNo = *filter.PageNo
	}
	pageSize := 10
	if filter.PageSize != nil && *filter.PageSize > 0 {
		pageSize = *filter.PageSize
	}

	orderBy := "id desc"
	if filter.OrderByColumn != "" {
		orderBy = filter.OrderByColumn
		isAsc := true
		if filter.IsAsc != nil {
			isAsc = *filter.IsAsc
		}
		if !isAsc {
			orderBy += " desc"
		} else {
			orderBy += " asc"
		}
	}

	pagination := utils.Pagination{
		Page:  pageNo,
		Limit: pageSize,
		Sort:  orderBy,
	}

	customers, err := s.uow.Customers().FindByFilter(filter, &pagination)
	if err != nil {
		return nil, err
	}

	customerDTOs := make([]domain.ApiCustomer, 0, len(customers))
	for _, customer := range customers {
		customerDTOs = append(customerDTOs, domain.ApiCustomerFromModel(customer))
	}

	return &domain.ApiCustomerListPaginated{
		Customer:   customerDTOs,
		Pagination: *utils.ApiPaginationFromModel(&pagination),
	}, nil
}

// GetCustomerByID busca um cliente pelo ID
func (s *CustomerService) GetCustomerByID(id uint) (*domain.ApiCustomerDetail, error) {
	customer, err := s.uow.Customers().FindByID(id)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, utils.ErrNotFound
	}

	customerDetailDTO := domain.ApiCustomerDetailFromModel(*customer)
	return &customerDetailDTO, nil
}

// CreateCustomer cria um novo cliente sob transação atômica do Unit of Work
func (s *CustomerService) CreateCustomer(req domain.CreateCustomerRequest, userID uint) (*domain.ApiCustomer, error) {
	if err := s.validator.ValidateForCreation(req); err != nil {
		return nil, err
	}

	document := strings.NewReplacer(".", "", "-", "", "/", "").Replace(req.DocumentNumber)

	customer := domain.Customer{
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		PersonType:     req.PersonType,
		DocumentNumber: document,
		CompanyName:    req.CompanyName,
		Notes:          req.Notes,
		CreatedByID:    &userID,
		IsActive:       true,
	}

	err := s.uow.Execute(func(uow repository.UnitOfWork) error {
		return uow.Customers().Create(&customer)
	})
	if err != nil {
		return nil, err
	}

	customerDTO := domain.ApiCustomerFromModel(customer)
	return &customerDTO, nil
}

// UpdateCustomer atualiza um cliente existente sob transação do Unit of Work
func (s *CustomerService) UpdateCustomer(id uint, req domain.UpdateCustomerRequest) (*domain.ApiCustomer, error) {
	if err := s.validator.ValidateForUpdate(id, req); err != nil {
		return nil, err
	}

	customer, err := s.uow.Customers().FindByID(id)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, utils.ErrNotFound
	}

	customer.FirstName = req.FirstName
	customer.LastName = req.LastName
	customer.CompanyName = req.CompanyName
	customer.IsActive = req.IsActive
	customer.Notes = req.Notes

	document := strings.NewReplacer(".", "", "-", "", "/", "").Replace(req.DocumentNumber)
	customer.DocumentNumber = document

	err = s.uow.Execute(func(uow repository.UnitOfWork) error {
		return uow.Customers().Update(customer)
	})
	if err != nil {
		return nil, err
	}

	customerDTO := domain.ApiCustomerFromModel(*customer)
	return &customerDTO, nil
}

// DeleteCustomer exclui um cliente sob transação do Unit of Work
func (s *CustomerService) DeleteCustomer(id uint) error {
	customer, err := s.uow.Customers().FindByID(id)
	if err != nil {
		return err
	}
	if customer == nil {
		return utils.ErrNotFound
	}

	return s.uow.Execute(func(uow repository.UnitOfWork) error {
		return uow.Customers().Delete(id)
	})
}
