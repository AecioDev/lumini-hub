package repository

import (
	"lumini-hub/api.core/internal/domain"
	commonrepo "lumini-hub/common/repository"

	"gorm.io/gorm"
)

// CompanyColorPaletteRepository define as operações de acesso a dados para
// as paletas de cores personalizadas de uma empresa (relação N:1 com Company).
type CompanyColorPaletteRepository interface {
	commonrepo.Repository[domain.CompanyColorPalette]
	FindAllByCompanyID(companyID uint) ([]domain.CompanyColorPalette, error)
}

// GormCompanyColorPaletteRepository implementa CompanyColorPaletteRepository usando GORM e Generics
type GormCompanyColorPaletteRepository struct {
	*commonrepo.GormRepository[domain.CompanyColorPalette]
}

// NewCompanyColorPaletteRepository cria um novo repository de paletas de cores personalizadas
func NewCompanyColorPaletteRepository(db *gorm.DB) CompanyColorPaletteRepository {
	return &GormCompanyColorPaletteRepository{
		GormRepository: commonrepo.NewGormRepository[domain.CompanyColorPalette](db),
	}
}

// FindAllByCompanyID lista as paletas salvas de uma empresa, mais recentes primeiro
func (r *GormCompanyColorPaletteRepository) FindAllByCompanyID(companyID uint) ([]domain.CompanyColorPalette, error) {
	var palettes []domain.CompanyColorPalette
	err := r.GetDB().Where("company_id = ?", companyID).Order("created_at DESC").Find(&palettes).Error
	if err != nil {
		return nil, err
	}
	return palettes, nil
}
