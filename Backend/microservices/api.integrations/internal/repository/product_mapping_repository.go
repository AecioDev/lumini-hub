package repository

import (
	"errors"

	"lumini-hub/api.integrations/internal/domain"
	commonrepo "lumini-hub/common/repository"

	"gorm.io/gorm"
)

// ProductMappingRepository define as operações de acesso a dados para o mapeamento de produtos
type ProductMappingRepository interface {
	commonrepo.Repository[domain.ProductMapping]
	FindByCodPro(codPro string) (*domain.ProductMapping, error)
}

// GormProductMappingRepository implementa ProductMappingRepository usando GORM
type GormProductMappingRepository struct {
	*commonrepo.GormRepository[domain.ProductMapping]
}

// NewProductMappingRepository cria um novo repository de mapeamento de produtos
func NewProductMappingRepository(db *gorm.DB) ProductMappingRepository {
	return &GormProductMappingRepository{
		GormRepository: commonrepo.NewGormRepository[domain.ProductMapping](db),
	}
}

// FindByCodPro busca o mapeamento de um produto pelo código no ERP legado (CADPRO.codpro)
func (r *GormProductMappingRepository) FindByCodPro(codPro string) (*domain.ProductMapping, error) {
	var mapping domain.ProductMapping
	err := r.GetDB().Where("cod_pro = ?", codPro).First(&mapping).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &mapping, nil
}
