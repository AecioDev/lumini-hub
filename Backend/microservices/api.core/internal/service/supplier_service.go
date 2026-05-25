package service

import (
	"lumini-hub/api.core/internal/domain"
	"lumini-hub/api.core/internal/repository"
	"lumini-hub/api.core/internal/validator"
	"lumini-hub/common/utils"
	"strings"
)

// SupplierService gerencia operações de negócios relacionadas a fornecedores
type SupplierService struct {
	uow       repository.UnitOfWork
	validator *validator.SupplierValidator
}

// NewSupplierService cria um novo serviço de fornecedores recebendo o Unit of Work
func NewSupplierService(uow repository.UnitOfWork) *SupplierService {
	return &SupplierService{
		uow:       uow,
		validator: validator.NewSupplierValidator(uow.Suppliers()),
	}
}

// GetSuppliers retorna uma lista paginada de fornecedores (legada via GET)
func (s *SupplierService) GetSuppliers(pagination *utils.Pagination) (*domain.ApiSupplierListPaginated, error) {
	// Utiliza o FindByFilter genérico sem filtros para obter todos paginados
	suppliers, err := s.uow.Suppliers().FindByFilter(domain.SupplierFilterRequest{}, pagination)
	if err != nil {
		return nil, err
	}

	supplierDTOs := make([]domain.ApiSupplier, 0, len(suppliers))
	for _, supplier := range suppliers {
		supplierDTOs = append(supplierDTOs, domain.ApiSupplierFromModel(supplier))
	}

	return &domain.ApiSupplierListPaginated{
		Supplier:   supplierDTOs,
		Pagination: *utils.ApiPaginationFromModel(pagination),
	}, nil
}

// GetSuppliersByFilter retorna fornecedores com filtros avançados e paginação (PostFilter)
func (s *SupplierService) GetSuppliersByFilter(filter domain.SupplierFilterRequest) (*domain.ApiSupplierListPaginated, error) {
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

	suppliers, err := s.uow.Suppliers().FindByFilter(filter, &pagination)
	if err != nil {
		return nil, err
	}

	supplierDTOs := make([]domain.ApiSupplier, 0, len(suppliers))
	for _, supplier := range suppliers {
		supplierDTOs = append(supplierDTOs, domain.ApiSupplierFromModel(supplier))
	}

	return &domain.ApiSupplierListPaginated{
		Supplier:   supplierDTOs,
		Pagination: *utils.ApiPaginationFromModel(&pagination),
	}, nil
}

// GetSupplierByID busca um fornecedor pelo ID
func (s *SupplierService) GetSupplierByID(id uint) (*domain.ApiSupplierDetail, error) {
	supplier, err := s.uow.Suppliers().FindByID(id)
	if err != nil {
		return nil, err
	}
	if supplier == nil {
		return nil, utils.ErrNotFound
	}

	supplierDetailDTO := domain.ApiSupplierDetailFromModel(*supplier)
	return &supplierDetailDTO, nil
}

// CreateSupplier cria um novo fornecedor sob transação atômica do Unit of Work
func (s *SupplierService) CreateSupplier(req domain.CreateSupplierRequest, userID uint) (*domain.ApiSupplier, error) {
	if err := s.validator.ValidateForCreation(req); err != nil {
		return nil, err
	}

	document := strings.NewReplacer(".", "", "-", "", "/", "").Replace(req.DocumentNumber)

	creatorID := userID
	if req.CreatedByID > 0 {
		creatorID = req.CreatedByID
	}

	supplier := domain.Supplier{
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		PersonType:     req.PersonType,
		DocumentNumber: document,
		CompanyName:    req.CompanyName,
		Notes:          req.Notes,
		CreatedByID:    &creatorID,
		IsActive:       true,
	}

	err := s.uow.Execute(func(uow repository.UnitOfWork) error {
		return uow.Suppliers().Create(&supplier)
	})
	if err != nil {
		return nil, err
	}

	supplierDTO := domain.ApiSupplierFromModel(supplier)
	return &supplierDTO, nil
}

// UpdateSupplier atualiza um fornecedor existente sob transação do Unit of Work
func (s *SupplierService) UpdateSupplier(id uint, req domain.UpdateSupplierRequest, userID uint) (*domain.ApiSupplier, error) {
	if err := s.validator.ValidateForUpdate(id, req); err != nil {
		return nil, err
	}

	supplier, err := s.uow.Suppliers().FindByID(id)
	if err != nil {
		return nil, err
	}
	if supplier == nil {
		return nil, utils.ErrNotFound
	}

	supplier.FirstName = req.FirstName
	supplier.LastName = req.LastName
	supplier.CompanyName = req.CompanyName
	supplier.IsActive = req.IsActive
	supplier.Notes = req.Notes
	supplier.UpdatedByID = &userID

	document := strings.NewReplacer(".", "", "-", "", "/", "").Replace(req.DocumentNumber)
	supplier.DocumentNumber = document

	err = s.uow.Execute(func(uow repository.UnitOfWork) error {
		return uow.Suppliers().Update(supplier)
	})
	if err != nil {
		return nil, err
	}

	supplierDTO := domain.ApiSupplierFromModel(*supplier)
	return &supplierDTO, nil
}

// DeleteSupplier exclui um fornecedor sob transação do Unit of Work
func (s *SupplierService) DeleteSupplier(id uint) error {
	supplier, err := s.uow.Suppliers().FindByID(id)
	if err != nil {
		return err
	}
	if supplier == nil {
		return utils.ErrNotFound
	}

	return s.uow.Execute(func(uow repository.UnitOfWork) error {
		return uow.Suppliers().Delete(id)
	})
}
