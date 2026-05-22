package repository

import (
	"errors"
	"simple-erp-service/internal/data-structure/models"
	"simple-erp-service/internal/utils"

	"gorm.io/gorm"
)

// CustomerRepository define as operações de acesso a dados para clientes
type CustomerRepository interface {
	Repository
	FindAll(pagination *models.Pagination) ([]models.Customer, error)
	FindByID(id uint) (*models.Customer, error)
	FindByDocument(document string) (*models.Customer, error)
	Create(customer *models.Customer) error
	Update(customer *models.Customer) error
	Delete(id uint) error
	ExistsByDocument(document string) (bool, error)
	ExistsByDocumentExcept(document string, id uint) (bool, error)
	ExistsByEmail(email string) (bool, error)
	ExistsByEmailExcept(email string, id uint) (bool, error)
}

// GormCustomerRepository implementa CustomerRepository usando GORM
type GormCustomerRepository struct {
	*BaseRepository
}

// NewCustomerRepository cria um novo repository de clientes
func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &GormCustomerRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// FindAll retorna todos os clientes com paginação
func (r *GormCustomerRepository) FindAll(pagination *models.Pagination) ([]models.Customer, error) {
	var customers []models.Customer

	query := r.GetDB().Model(&models.Customer{})
	query, err := utils.Paginate(&models.Customer{}, pagination, query)
	if err != nil {
		return nil, err
	}

	if err := query.Find(&customers).Error; err != nil {
		return nil, err
	}

	return customers, nil
}

// FindByID busca um cliente pelo ID
func (r *GormCustomerRepository) FindByID(id uint) (*models.Customer, error) {
	var customer models.Customer
	if err := r.GetDB().First(&customer, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &customer, nil
}

// FindByDocument busca um cliente pelo documento
func (r *GormCustomerRepository) FindByDocument(document string) (*models.Customer, error) {
	var customer models.Customer
	if err := r.GetDB().Where("document = ?", document).First(&customer).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &customer, nil
}

// Create cria um novo cliente
func (r *GormCustomerRepository) Create(customer *models.Customer) error {
	return r.GetDB().Create(customer).Error
}

// Update atualiza um cliente existente
func (r *GormCustomerRepository) Update(customer *models.Customer) error {
	return r.GetDB().Save(customer).Error
}

// Delete exclui um cliente (soft delete)
func (r *GormCustomerRepository) Delete(id uint) error {
	return r.GetDB().Delete(&models.Customer{}, id).Error
}

// ExistsByDocument verifica se existe um cliente com o documento especificado
func (r *GormCustomerRepository) ExistsByDocument(document string) (bool, error) {
	var count int64
	err := r.GetDB().Model(&models.Customer{}).Where("document = ?", document).Count(&count).Error
	return count > 0, err
}

// ExistsByDocumentExcept verifica se existe um cliente com o documento especificado, exceto o cliente com o ID especificado
func (r *GormCustomerRepository) ExistsByDocumentExcept(document string, id uint) (bool, error) {
	var count int64
	err := r.GetDB().Model(&models.Customer{}).Where("document = ? AND id != ?", document, id).Count(&count).Error
	return count > 0, err
}

// ExistsByEmail verifica se existe um cliente com o email especificado
func (r *GormCustomerRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := r.GetDB().Model(&models.Customer{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

// ExistsByEmailExcept verifica se existe um cliente com o email especificado, exceto o cliente com o ID especificado
func (r *GormCustomerRepository) ExistsByEmailExcept(email string, id uint) (bool, error) {
	var count int64
	err := r.GetDB().Model(&models.Customer{}).Where("email = ? AND id != ?", email, id).Count(&count).Error
	return count > 0, err
}
