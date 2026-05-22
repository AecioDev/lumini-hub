package repository

import (
	"errors"
	"lumini-hub/api.core/internal/domain"
	"lumini-hub/common/utils"

	"gorm.io/gorm"
)

// SupplierRepository define as operações de acesso a dados para fornecedores
type SupplierRepository interface {
	Repository
	FindAll(pagination *utils.Pagination) ([]domain.Supplier, error)
	FindByID(id uint) (*domain.Supplier, error)
	FindByDocument(document string) (*domain.Supplier, error)
	Create(supplier *domain.Supplier) error
	Update(supplier *domain.Supplier) error
	Delete(id uint) error
	ExistsByDocument(document string) (bool, error)
	ExistsByDocumentExcept(document string, id uint) (bool, error)
}

// GormSupplierRepository implementa SupplierRepository usando GORM
type GormSupplierRepository struct {
	*BaseRepository
}

// NewSupplierRepository cria um novo repository de fornecedores
func NewSupplierRepository(db *gorm.DB) SupplierRepository {
	return &GormSupplierRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// FindAll retorna todos os fornecedores com paginação
func (r *GormSupplierRepository) FindAll(pagination *utils.Pagination) ([]domain.Supplier, error) {
	var suppliers []domain.Supplier

	query := r.GetDB().Model(&domain.Supplier{})
	query, err := utils.Paginate(&domain.Supplier{}, pagination, query)
	if err != nil {
		return nil, err
	}

	if err := query.Preload("Addresses.City.State.Country").Preload("Contacts").Find(&suppliers).Error; err != nil {
		return nil, err
	}

	return suppliers, nil
}

// FindByID busca um fornecedor pelo ID
func (r *GormSupplierRepository) FindByID(id uint) (*domain.Supplier, error) {
	var supplier domain.Supplier
	if err := r.GetDB().Preload("Addresses.City.State.Country").Preload("Contacts").First(&supplier, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &supplier, nil
}

// FindByDocument busca um fornecedor pelo documento
func (r *GormSupplierRepository) FindByDocument(document string) (*domain.Supplier, error) {
	var supplier domain.Supplier
	if err := r.GetDB().Where("document_number = ?", document).First(&supplier).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &supplier, nil
}

// Create cria um novo fornecedor
func (r *GormSupplierRepository) Create(supplier *domain.Supplier) error {
	return r.GetDB().Create(supplier).Error
}

// Update atualiza um fornecedor existente
func (r *GormSupplierRepository) Update(supplier *domain.Supplier) error {
	return r.GetDB().Save(supplier).Error
}

// Delete exclui um fornecedor (soft delete)
func (r *GormSupplierRepository) Delete(id uint) error {
	return r.GetDB().Delete(&domain.Supplier{}, id).Error
}

// ExistsByDocument verifica se existe um fornecedor com o documento especificado
func (r *GormSupplierRepository) ExistsByDocument(document string) (bool, error) {
	var count int64
	err := r.GetDB().Model(&domain.Supplier{}).Where("document_number = ?", document).Count(&count).Error
	return count > 0, err
}

// ExistsByDocumentExcept verifica se existe um fornecedor com o documento especificado, exceto o fornecedor com o ID especificado
func (r *GormSupplierRepository) ExistsByDocumentExcept(document string, id uint) (bool, error) {
	var count int64
	err := r.GetDB().Model(&domain.Supplier{}).Where("document_number = ? AND id != ?", document, id).Count(&count).Error
	return count > 0, err
}
