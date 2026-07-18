package repository

import (
	"lumini-hub/api.integrations/internal/legacy/domain"

	"gorm.io/gorm"
)

// LocArmRepository lê os locais de armazenamento cadastrados no ERP legado
type LocArmRepository interface {
	FindAll() ([]domain.LocArm, error)
}

// GormLocArmRepository implementa LocArmRepository usando GORM sobre o SQL Server
type GormLocArmRepository struct {
	db *gorm.DB
}

// NewLocArmRepository cria um novo repository de locais de armazenamento
func NewLocArmRepository(db *gorm.DB) LocArmRepository {
	return &GormLocArmRepository{db: db}
}

// FindAll retorna todos os locais de armazenamento cadastrados
func (r *GormLocArmRepository) FindAll() ([]domain.LocArm, error) {
	var locations []domain.LocArm
	err := r.db.Order("CodLocAm asc").Find(&locations).Error
	if err != nil {
		return nil, err
	}
	return locations, nil
}
