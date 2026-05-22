package repository

import (
	"errors"
	"simple-erp-service/internal/data-structure/models"
	"simple-erp-service/internal/utils"

	"gorm.io/gorm"
)

// SupplierRepository define as operações de acesso a dados para fornecedores
type SupplierRepository interface {
	Repository
	FindAll(pagination *models.Pagination) ([]models.Supplier, error)
	FindByID(id uint) (*models.Supplier, error)
	FindByDocument(document string) (*models.Supplier, error)
	Create(supplier *models.Supplier) error
	Update(supplier *models.Supplier) error
	Delete(id uint) error
	ExistsByDocument(document string) (bool, error)
	ExistsByDocumentExcept(document string, id uint) (bool, error)
	ExistsByEmail(email string) (bool, error)
	ExistsByEmailExcept(email string, id uint) (bool, error)
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
func (r *GormSupplierRepository) FindAll(pagination *models.Pagination) ([]models.Supplier, error) {
	var suppliers []models.Supplier

	query := r.GetDB().Model(&models.Supplier{})
	query, err := utils.Paginate(&models.Supplier{}, pagination, query)
	if err != nil {
		return nil, err
	}

	if err := query.Find(&suppliers).Error; err != nil {
		return nil, err
	}

	return suppliers, nil
}

// FindByID busca um fornecedor pelo ID
func (r *GormSupplierRepository) FindByID(id uint) (*models.Supplier, error) {
	var supplier models.Supplier
	if err := r.GetDB().First(&supplier, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &supplier, nil
}

// FindByDocument busca um fornecedor pelo documento
func (r *GormSupplierRepository) FindByDocument(document string) (*models.Supplier, error) {
	var supplier models.Supplier
	if err := r.GetDB().Where("document = ?", document).First(&supplier).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &supplier, nil
}

// Create cria um novo fornecedor
func (r *GormSupplierRepository) Create(supplier *models.Supplier) error {
	return r.GetDB().Create(supplier).Error
}

// Update atualiza um fornecedor existente
func (r *GormSupplierRepository) Update(supplier *models.Supplier) error {
	return r.GetDB().Save(supplier).Error
}

// Delete exclui um fornecedor (soft delete)
func (r *GormSupplierRepository) Delete(id uint) error {
	return r.GetDB().Delete(&models.Supplier{}, id).Error
}

// ExistsByDocument verifica se existe um fornecedor com o documento especificado
func (r *GormSupplierRepository) ExistsByDocument(document string) (bool, error) {
	var count int64
	err := r.GetDB().Model(&models.Supplier{}).Where("document = ?", document).Count(&count).Error
	return count > 0, err
}

// ExistsByDocumentExcept verifica se existe um fornecedor com o documento especificado, exceto o fornecedor com o ID especificado
func (r *GormSupplierRepository) ExistsByDocumentExcept(document string, id uint) (bool, error) {
	var count int64
	err := r.GetDB().Model(&models.Supplier{}).Where("document = ? AND id != ?", document, id).Count(&count).Error
	return count > 0, err
}

// ExistsByEmail verifica se existe um fornecedor com o email especificado
func (r *GormSupplierRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := r.GetDB().Model(&models.Supplier{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

// ExistsByEmailExcept verifica se existe um fornecedor com o email especificado, exceto o fornecedor com o ID especificado
func (r *GormSupplierRepository) ExistsByEmailExcept(email string, id uint) (bool, error) {
	var count int64
	err := r.GetDB().Model(&models.Supplier{}).Where("email = ? AND id != ?", email, id).Count(&count).Error
	return count > 0, err
}
