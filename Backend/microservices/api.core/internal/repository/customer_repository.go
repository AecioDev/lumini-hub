package repository

import (
	"errors"
	"lumini-hub/api.core/internal/domain"
	commonrepo "lumini-hub/common/repository"
	"lumini-hub/common/utils"

	"gorm.io/gorm"
)

// CustomerRepository define as operações de acesso a dados para clientes
type CustomerRepository interface {
	commonrepo.Repository[domain.Customer]
	FindByID(id uint) (*domain.Customer, error) // Sobrescreve para incluir preloads
	FindByDocument(document string) (*domain.Customer, error)
	FindByFilter(filter domain.CustomerFilterRequest, pagination *utils.Pagination) ([]domain.Customer, error)
	ExistsByDocument(document string) (bool, error)
	ExistsByDocumentExcept(document string, id uint) (bool, error)
}

// GormCustomerRepository implementa CustomerRepository usando GORM e Generics
type GormCustomerRepository struct {
	*commonrepo.GormRepository[domain.Customer]
}

// NewCustomerRepository cria um novo repository de clientes
func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &GormCustomerRepository{
		GormRepository: commonrepo.NewGormRepository[domain.Customer](db),
	}
}

// FindByID busca um cliente pelo ID pré-carregando endereços e contatos
func (r *GormCustomerRepository) FindByID(id uint) (*domain.Customer, error) {
	var customer domain.Customer
	err := r.GetDB().Preload("Addresses.City.State.Country").Preload("Contacts").First(&customer, id).Error
	if err != nil {
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
	err := r.GetDB().Where("document_number = ?", document).First(&customer).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &customer, nil
}

// FindByFilter busca clientes aplicando filtros dinâmicos e paginação
func (r *GormCustomerRepository) FindByFilter(filter domain.CustomerFilterRequest, pagination *utils.Pagination) ([]domain.Customer, error) {
	var customers []domain.Customer

	query := r.GetDB().Model(&domain.Customer{})

	// Aplicar filtros dinâmicos se estiverem preenchidos
	if filter.ID != nil && *filter.ID > 0 {
		query = query.Where("id = ?", *filter.ID)
	}
	if filter.FirstName != "" {
		query = query.Where("first_name ILIKE ?", "%"+filter.FirstName+"%")
	}
	if filter.LastName != "" {
		query = query.Where("last_name ILIKE ?", "%"+filter.LastName+"%")
	}
	if filter.PersonType != "" {
		query = query.Where("person_type = ?", filter.PersonType)
	}
	if filter.DocumentNumber != "" {
		query = query.Where("document_number = ?", filter.DocumentNumber)
	}
	if filter.CompanyName != "" {
		query = query.Where("company_name ILIKE ?", "%"+filter.CompanyName+"%")
	}
	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}

	// Executar paginação e contagem usando o helper comum
	query, err := utils.Paginate(&domain.Customer{}, pagination, query)
	if err != nil {
		return nil, err
	}

	// Buscar registros com os relacionamentos pré-carregados
	err = query.Preload("Addresses.City.State.Country").Preload("Contacts").Find(&customers).Error
	if err != nil {
		return nil, err
	}

	return customers, nil
}

// ExistsByDocument verifica se existe um cliente com o documento especificado
func (r *GormCustomerRepository) ExistsByDocument(document string) (bool, error) {
	var count int64
	err := r.GetDB().Model(&domain.Customer{}).Where("document_number = ?", document).Count(&count).Error
	return count > 0, err
}

// ExistsByDocumentExcept verifica se existe um cliente com o documento especificado, exceto o de ID informado
func (r *GormCustomerRepository) ExistsByDocumentExcept(document string, id uint) (bool, error) {
	var count int64
	err := r.GetDB().Model(&domain.Customer{}).Where("document_number = ? AND id != ?", document, id).Count(&count).Error
	return count > 0, err
}
