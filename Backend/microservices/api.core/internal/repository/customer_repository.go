package repository

import (
	"errors"
	"lumini-hub/api.core/internal/domain"
	"lumini-hub/common/utils"

	"gorm.io/gorm"
)

// CustomerRepository define as operações de acesso a dados para clientes
type CustomerRepository interface {
	Repository
	FindAll(pagination *utils.Pagination) ([]domain.Customer, error)
	FindByID(id uint) (*domain.Customer, error)
	FindByDocument(document string) (*domain.Customer, error)
	Create(customer *domain.Customer) error
	Update(customer *domain.Customer) error
	Delete(id uint) error
	ExistsByDocument(document string) (bool, error)
	ExistsByDocumentExcept(document string, id uint) (bool, error)
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
func (r *GormCustomerRepository) FindAll(pagination *utils.Pagination) ([]domain.Customer, error) {
	var customers []domain.Customer

	query := r.GetDB().Model(&domain.Customer{})
	query, err := utils.Paginate(&domain.Customer{}, pagination, query)
	if err != nil {
		return nil, err
	}

	// Carregar relacionamentos de endereço e contato para exibição completa se o frontend requerer
	if err := query.Preload("Addresses.City.State.Country").Preload("Contacts").Find(&customers).Error; err != nil {
		return nil, err
	}

	return customers, nil
}

// FindByID busca um cliente pelo ID
func (r *GormCustomerRepository) FindByID(id uint) (*domain.Customer, error) {
	var customer domain.Customer
	if err := r.GetDB().Preload("Addresses.City.State.Country").Preload("Contacts").First(&customer, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &customer, nil
}

// FindByDocument busca um cliente pelo documento
func (r *GormCustomerRepository) FindByDocument(document string) (*domain.Customer, error) {
	var customer domain.Customer
	if err := r.GetDB().Where("document_number = ?", document).First(&customer).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &customer, nil
}

// Create cria um novo cliente
func (r *GormCustomerRepository) Create(customer *domain.Customer) error {
	return r.GetDB().Create(customer).Error
}

// Update atualiza um cliente existente
func (r *GormCustomerRepository) Update(customer *domain.Customer) error {
	return r.GetDB().Save(customer).Error
}

// Delete exclui um cliente (soft delete)
func (r *GormCustomerRepository) Delete(id uint) error {
	return r.GetDB().Delete(&domain.Customer{}, id).Error
}

// ExistsByDocument verifica se existe um cliente com o documento especificado
func (r *GormCustomerRepository) ExistsByDocument(document string) (bool, error) {
	var count int64
	err := r.GetDB().Model(&domain.Customer{}).Where("document_number = ?", document).Count(&count).Error
	return count > 0, err
}

// ExistsByDocumentExcept verifica se existe um cliente com o documento especificado, exceto o cliente com o ID especificado
func (r *GormCustomerRepository) ExistsByDocumentExcept(document string, id uint) (bool, error) {
	var count int64
	err := r.GetDB().Model(&domain.Customer{}).Where("document_number = ? AND id != ?", document, id).Count(&count).Error
	return count > 0, err
}
