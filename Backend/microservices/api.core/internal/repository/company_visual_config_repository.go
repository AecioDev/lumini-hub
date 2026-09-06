package repository

import (
	"errors"

	"lumini-hub/api.core/internal/domain"
	commonrepo "lumini-hub/common/repository"

	"gorm.io/gorm"
)

// CompanyVisualConfigRepository define as operações de acesso a dados para a
// configuração visual de uma empresa (relação 1:1 com Company).
type CompanyVisualConfigRepository interface {
	commonrepo.Repository[domain.CompanyVisualConfig]
	FindByCompanyID(companyID uint) (*domain.CompanyVisualConfig, error)
	ExistsByCompanyID(companyID uint) (bool, error)
}

// GormCompanyVisualConfigRepository implementa CompanyVisualConfigRepository usando GORM e Generics
type GormCompanyVisualConfigRepository struct {
	*commonrepo.GormRepository[domain.CompanyVisualConfig]
}

// NewCompanyVisualConfigRepository cria um novo repository de configuração visual de empresa
func NewCompanyVisualConfigRepository(db *gorm.DB) CompanyVisualConfigRepository {
	return &GormCompanyVisualConfigRepository{
		GormRepository: commonrepo.NewGormRepository[domain.CompanyVisualConfig](db),
	}
}

// FindByCompanyID busca a configuração visual de uma empresa pelo ID da empresa
func (r *GormCompanyVisualConfigRepository) FindByCompanyID(companyID uint) (*domain.CompanyVisualConfig, error) {
	var config domain.CompanyVisualConfig
	err := r.GetDB().Where("company_id = ?", companyID).First(&config).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &config, nil
}

// ExistsByCompanyID verifica se já existe configuração visual cadastrada para a empresa
func (r *GormCompanyVisualConfigRepository) ExistsByCompanyID(companyID uint) (bool, error) {
	var count int64
	err := r.GetDB().Model(&domain.CompanyVisualConfig{}).Where("company_id = ?", companyID).Count(&count).Error
	return count > 0, err
}
