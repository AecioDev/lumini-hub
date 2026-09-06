package repository

import (
	"errors"

	"lumini-hub/api.core/internal/domain"
	commonrepo "lumini-hub/common/repository"

	"gorm.io/gorm"
)

// CompanyFiscalConfigRepository define as operações de acesso a dados para a
// configuração fiscal de uma empresa (relação 1:1 com Company).
type CompanyFiscalConfigRepository interface {
	commonrepo.Repository[domain.CompanyFiscalConfig]
	FindByCompanyID(companyID uint) (*domain.CompanyFiscalConfig, error)
	ExistsByCompanyID(companyID uint) (bool, error)
}

// GormCompanyFiscalConfigRepository implementa CompanyFiscalConfigRepository usando GORM e Generics
type GormCompanyFiscalConfigRepository struct {
	*commonrepo.GormRepository[domain.CompanyFiscalConfig]
}

// NewCompanyFiscalConfigRepository cria um novo repository de configuração fiscal de empresa
func NewCompanyFiscalConfigRepository(db *gorm.DB) CompanyFiscalConfigRepository {
	return &GormCompanyFiscalConfigRepository{
		GormRepository: commonrepo.NewGormRepository[domain.CompanyFiscalConfig](db),
	}
}

// FindByCompanyID busca a configuração fiscal de uma empresa pelo ID da empresa
func (r *GormCompanyFiscalConfigRepository) FindByCompanyID(companyID uint) (*domain.CompanyFiscalConfig, error) {
	var config domain.CompanyFiscalConfig
	err := r.GetDB().Where("company_id = ?", companyID).First(&config).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &config, nil
}

// ExistsByCompanyID verifica se já existe configuração fiscal cadastrada para a empresa
func (r *GormCompanyFiscalConfigRepository) ExistsByCompanyID(companyID uint) (bool, error) {
	var count int64
	err := r.GetDB().Model(&domain.CompanyFiscalConfig{}).Where("company_id = ?", companyID).Count(&count).Error
	return count > 0, err
}
