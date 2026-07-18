package repository

import (
	"lumini-hub/api.integrations/internal/legacy/domain"

	"gorm.io/gorm"
)

// CencusRepository lê as empresas cadastradas no ERP legado
type CencusRepository interface {
	FindAll() ([]domain.Cencus, error)
}

// GormCencusRepository implementa CencusRepository usando GORM sobre o SQL Server
type GormCencusRepository struct {
	db *gorm.DB
}

// NewCencusRepository cria um novo repository de empresas
func NewCencusRepository(db *gorm.DB) CencusRepository {
	return &GormCencusRepository{db: db}
}

// FindAll retorna todas as empresas cadastradas
func (r *GormCencusRepository) FindAll() ([]domain.Cencus, error) {
	var companies []domain.Cencus
	err := r.db.Order("codcencus asc").Find(&companies).Error
	if err != nil {
		return nil, err
	}
	return companies, nil
}
