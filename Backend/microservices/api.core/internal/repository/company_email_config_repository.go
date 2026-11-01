package repository

import (
	"errors"

	"lumini-hub/api.core/internal/domain"
	commonrepo "lumini-hub/common/repository"

	"gorm.io/gorm"
)

// CompanyEmailConfigRepository define as operações de acesso a dados para o servidor SMTP próprio de uma empresa (1:1).
type CompanyEmailConfigRepository interface {
	commonrepo.Repository[domain.CompanyEmailConfig]
	FindByCompanyID(companyID uint) (*domain.CompanyEmailConfig, error)
}

// GormCompanyEmailConfigRepository implementa CompanyEmailConfigRepository usando GORM e Generics
type GormCompanyEmailConfigRepository struct {
	*commonrepo.GormRepository[domain.CompanyEmailConfig]
}

// NewCompanyEmailConfigRepository cria um novo repository de configuração de e-mail de empresa
func NewCompanyEmailConfigRepository(db *gorm.DB) CompanyEmailConfigRepository {
	return &GormCompanyEmailConfigRepository{
		GormRepository: commonrepo.NewGormRepository[domain.CompanyEmailConfig](db),
	}
}

// FindByCompanyID busca a configuração de e-mail de uma empresa; nil quando ela não tem uma própria
func (r *GormCompanyEmailConfigRepository) FindByCompanyID(companyID uint) (*domain.CompanyEmailConfig, error) {
	var config domain.CompanyEmailConfig
	err := r.GetDB().Where("company_id = ?", companyID).First(&config).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &config, nil
}
