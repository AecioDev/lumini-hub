package repository

import (
	"errors"

	"lumini-hub/api.core/internal/domain"
	commonrepo "lumini-hub/common/repository"

	"gorm.io/gorm"
)

// CompanyRepository define as operações de acesso a dados para empresas.
//
// Sem método de busca paginada/filtrada de propósito: o volume esperado de
// empresas por tenant é baixo (um grupo com sua Matriz e algumas filiais),
// não justifica o padrão POST /filter usado em entidades de alto volume
// (Customer/Supplier) — ver lumini_hub_backend_architecture skill.
type CompanyRepository interface {
	commonrepo.Repository[domain.Company]
	FindByID(id uint) (*domain.Company, error) // Sobrescreve para incluir preloads
	FindByTaxID(taxID string) (*domain.Company, error)
	ExistsByTaxID(taxID string) (bool, error)
	ExistsByTaxIDExcept(taxID string, id uint) (bool, error)
	CountByParentID(parentID uint) (int64, error)
}

// GormCompanyRepository implementa CompanyRepository usando GORM e Generics
type GormCompanyRepository struct {
	*commonrepo.GormRepository[domain.Company]
}

// NewCompanyRepository cria um novo repository de empresas
func NewCompanyRepository(db *gorm.DB) CompanyRepository {
	return &GormCompanyRepository{
		GormRepository: commonrepo.NewGormRepository[domain.Company](db),
	}
}

// FindByID busca uma empresa pelo ID pré-carregando as empresas vinculadas diretamente abaixo dela
func (r *GormCompanyRepository) FindByID(id uint) (*domain.Company, error) {
	var company domain.Company
	err := r.GetDB().Preload("Children").First(&company, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &company, nil
}

// FindByTaxID busca uma empresa pelo CNPJ
func (r *GormCompanyRepository) FindByTaxID(taxID string) (*domain.Company, error) {
	var company domain.Company
	err := r.GetDB().Where("tax_id = ?", taxID).First(&company).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &company, nil
}

// ExistsByTaxID verifica se existe uma empresa com o CNPJ especificado
func (r *GormCompanyRepository) ExistsByTaxID(taxID string) (bool, error) {
	var count int64
	err := r.GetDB().Model(&domain.Company{}).Where("tax_id = ?", taxID).Count(&count).Error
	return count > 0, err
}

// ExistsByTaxIDExcept verifica se existe uma empresa com o CNPJ especificado, exceto a de ID informado
func (r *GormCompanyRepository) ExistsByTaxIDExcept(taxID string, id uint) (bool, error) {
	var count int64
	err := r.GetDB().Model(&domain.Company{}).Where("tax_id = ? AND id != ?", taxID, id).Count(&count).Error
	return count > 0, err
}

// CountByParentID conta quantas empresas estão vinculadas diretamente abaixo da empresa informada
func (r *GormCompanyRepository) CountByParentID(parentID uint) (int64, error) {
	var count int64
	err := r.GetDB().Model(&domain.Company{}).Where("parent_id = ?", parentID).Count(&count).Error
	return count, err
}
