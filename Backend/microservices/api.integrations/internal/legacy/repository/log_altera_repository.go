package repository

import (
	"time"

	"lumini-hub/api.integrations/internal/legacy/domain"

	"gorm.io/gorm"
)

// LogAlteraRepository lê o log de alterações do ERP legado usado como cursor de polling
type LogAlteraRepository interface {
	FindChangesSince(cursor time.Time, tipos []string) ([]domain.LogAltera, error)
}

// GormLogAlteraRepository implementa LogAlteraRepository usando GORM sobre o SQL Server
type GormLogAlteraRepository struct {
	db *gorm.DB
}

// NewLogAlteraRepository cria um novo repository de leitura do LogAltera
func NewLogAlteraRepository(db *gorm.DB) LogAlteraRepository {
	return &GormLogAlteraRepository{db: db}
}

// FindChangesSince retorna as alterações registradas após o cursor informado, filtrando pelos tipos desejados
// (ex.: TipoAltProduto, TipoAltPrecoVenda, TipoAltCusto)
func (r *GormLogAlteraRepository) FindChangesSince(cursor time.Time, tipos []string) ([]domain.LogAltera, error) {
	var changes []domain.LogAltera
	err := r.db.
		Where("DatFimAlt > ?", cursor).
		Where("TipoAlt IN ?", tipos).
		Order("DatFimAlt asc").
		Find(&changes).Error
	if err != nil {
		return nil, err
	}
	return changes, nil
}
