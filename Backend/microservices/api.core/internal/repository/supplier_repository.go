package repository

import (
	"errors"
	"lumini-hub/api.core/internal/domain"
	commonrepo "lumini-hub/common/repository"
	"lumini-hub/common/utils"

	"gorm.io/gorm"
)

// SupplierRepository define as operações de acesso a dados para fornecedores
type SupplierRepository interface {
	commonrepo.Repository[domain.Supplier]
	FindByID(id uint) (*domain.Supplier, error) // Sobrescreve para incluir preloads
	FindByDocument(document string) (*domain.Supplier, error)
	FindByFilter(filter domain.SupplierFilterRequest, pagination *utils.Pagination) ([]domain.Supplier, error)
	ExistsByDocument(document string) (bool, error)
	ExistsByDocumentExcept(document string, id uint) (bool, error)
}

// GormSupplierRepository implementa SupplierRepository usando GORM e Generics
type GormSupplierRepository struct {
	*commonrepo.GormRepository[domain.Supplier]
}

// NewSupplierRepository cria um novo repository de fornecedores
func NewSupplierRepository(db *gorm.DB) SupplierRepository {
	return &GormSupplierRepository{
		GormRepository: commonrepo.NewGormRepository[domain.Supplier](db),
	}
}

// FindByID busca um fornecedor pelo ID pré-carregando endereços e contatos
func (r *GormSupplierRepository) FindByID(id uint) (*domain.Supplier, error) {
	var supplier domain.Supplier
	err := r.GetDB().Preload("Addresses.City.State.Country").Preload("Contacts").First(&supplier, id).Error
	if err != nil {
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
	err := r.GetDB().Where("document_number = ?", document).First(&supplier).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &supplier, nil
}

// FindByFilter busca fornecedores aplicando filtros dinâmicos e paginação
func (r *GormSupplierRepository) FindByFilter(filter domain.SupplierFilterRequest, pagination *utils.Pagination) ([]domain.Supplier, error) {
	var suppliers []domain.Supplier

	query := r.GetDB().Model(&domain.Supplier{})

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
	query, err := utils.Paginate(&domain.Supplier{}, pagination, query)
	if err != nil {
		return nil, err
	}

	// Buscar registros com os relacionamentos pré-carregados
	err = query.Preload("Addresses.City.State.Country").Preload("Contacts").Find(&suppliers).Error
	if err != nil {
		return nil, err
	}

	return suppliers, nil
}

// ExistsByDocument verifica se existe um fornecedor com o documento especificado
func (r *GormSupplierRepository) ExistsByDocument(document string) (bool, error) {
	var count int64
	err := r.GetDB().Model(&domain.Supplier{}).Where("document_number = ?", document).Count(&count).Error
	return count > 0, err
}

// ExistsByDocumentExcept verifica se existe um fornecedor com o documento especificado, exceto o de ID informado
func (r *GormSupplierRepository) ExistsByDocumentExcept(document string, id uint) (bool, error) {
	var count int64
	err := r.GetDB().Model(&domain.Supplier{}).Where("document_number = ? AND id != ?", document, id).Count(&count).Error
	return count > 0, err
}
