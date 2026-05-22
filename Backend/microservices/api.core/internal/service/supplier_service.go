package service

import (
	"lumini-hub/api.core/internal/domain"
	"lumini-hub/api.core/internal/repository"
	"lumini-hub/api.core/internal/validator"
	"lumini-hub/common/utils"
	"strings"
)

// SupplierService gerencia operações relacionadas a fornecedores
type SupplierService struct {
	supplierRepo repository.SupplierRepository
	validator    *validator.SupplierValidator
}

// NewSupplierService cria um novo serviço de fornecedores
func NewSupplierService(supplierRepo repository.SupplierRepository) *SupplierService {
	return &SupplierService{
		supplierRepo: supplierRepo,
		validator:    validator.NewSupplierValidator(supplierRepo),
	}
}

// GetSuppliers retorna uma lista paginada de fornecedores
func (s *SupplierService) GetSuppliers(pagination *utils.Pagination) (*domain.ApiSupplierListPaginated, error) {
	suppliers, err := s.supplierRepo.FindAll(pagination)
	if err != nil {
		return nil, err
	}

	// Converter para DTOs
	supplierDTOs := make([]domain.ApiSupplier, 0, len(suppliers))
	for _, supplier := range suppliers {
		supplierDTOs = append(supplierDTOs, domain.ApiSupplierFromModel(supplier))
	}

	return &domain.ApiSupplierListPaginated{
		Supplier:   supplierDTOs,
		Pagination: *utils.ApiPaginationFromModel(pagination),
	}, nil
}

// GetSupplierByID busca um fornecedor pelo ID
func (s *SupplierService) GetSupplierByID(id uint) (*domain.ApiSupplierDetail, error) {
	supplier, err := s.supplierRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if supplier == nil {
		return nil, utils.ErrNotFound
	}

	// Converter para DTO
	supplierDetailDTO := domain.ApiSupplierDetailFromModel(*supplier)
	return &supplierDetailDTO, nil
}

// CreateSupplier cria um novo fornecedor
func (s *SupplierService) CreateSupplier(req domain.CreateSupplierRequest, userID uint) (*domain.ApiSupplier, error) {
	// Validar dados
	if err := s.validator.ValidateForCreation(req); err != nil {
		return nil, err
	}

	// Formatar documento (remover caracteres não numéricos)
	document := strings.NewReplacer(".", "", "-", "", "/", "").Replace(req.DocumentNumber)

	// Resolver ID do criador (body ou context)
	creatorID := userID
	if req.CreatedByID > 0 {
		creatorID = req.CreatedByID
	}

	// Criar fornecedor
	supplier := domain.Supplier{
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		PersonType:     req.PersonType,
		DocumentNumber: document,
		CompanyName:    req.CompanyName,
		Notes:          req.Notes,
		CreatedByID:    &creatorID,
		IsActive:       true, // Por padrão, fornecedores são criados ativos
	}

	if err := s.supplierRepo.Create(&supplier); err != nil {
		return nil, err
	}

	// Converter para DTO
	supplierDTO := domain.ApiSupplierFromModel(supplier)
	return &supplierDTO, nil
}

// UpdateSupplier atualiza um fornecedor existente
func (s *SupplierService) UpdateSupplier(id uint, req domain.UpdateSupplierRequest, userID uint) (*domain.ApiSupplier, error) {
	// Validar dados
	if err := s.validator.ValidateForUpdate(id, req); err != nil {
		return nil, err
	}

	// Buscar fornecedor
	supplier, err := s.supplierRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if supplier == nil {
		return nil, utils.ErrNotFound
	}

	// Atualizar campos básicos
	supplier.FirstName = req.FirstName
	supplier.LastName = req.LastName
	supplier.CompanyName = req.CompanyName
	supplier.IsActive = req.IsActive
	supplier.Notes = req.Notes
	supplier.UpdatedByID = &userID

	// Atualizar e formatar o DocumentNumber (removendo caracteres especiais)
	document := strings.NewReplacer(".", "", "-", "", "/", "").Replace(req.DocumentNumber)
	supplier.DocumentNumber = document

	// Salvar alterações
	if err := s.supplierRepo.Update(supplier); err != nil {
		return nil, err
	}

	// Converter para DTO
	supplierDTO := domain.ApiSupplierFromModel(*supplier)
	return &supplierDTO, nil
}

// DeleteSupplier exclui um fornecedor (soft delete)
func (s *SupplierService) DeleteSupplier(id uint) error {
	// Verificar se o fornecedor existe
	supplier, err := s.supplierRepo.FindByID(id)
	if err != nil {
		return err
	}
	if supplier == nil {
		return utils.ErrNotFound
	}

	// Excluir fornecedor
	return s.supplierRepo.Delete(id)
}
