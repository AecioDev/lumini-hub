package repository

import (
	"lumini-hub/api.integrations/internal/domain"
	commonrepo "lumini-hub/common/repository"
	"lumini-hub/common/utils"

	"gorm.io/gorm"
)

// SyncLogRepository define as operações de acesso a dados para sync logs
type SyncLogRepository interface {
	commonrepo.Repository[domain.SyncLog]
	FindByFilter(filter domain.SyncLogFilterRequest, pagination *utils.Pagination) ([]domain.SyncLog, error)
}

// GormSyncLogRepository implementa SyncLogRepository usando GORM
type GormSyncLogRepository struct {
	*commonrepo.GormRepository[domain.SyncLog]
}

// NewSyncLogRepository cria um novo repository de sync logs
func NewSyncLogRepository(db *gorm.DB) SyncLogRepository {
	return &GormSyncLogRepository{
		GormRepository: commonrepo.NewGormRepository[domain.SyncLog](db),
	}
}

// FindByFilter busca sync logs aplicando filtros dinâmicos e paginação
func (r *GormSyncLogRepository) FindByFilter(filter domain.SyncLogFilterRequest, pagination *utils.Pagination) ([]domain.SyncLog, error) {
	var logs []domain.SyncLog

	query := r.GetDB().Model(&domain.SyncLog{})

	if filter.Direction != "" {
		query = query.Where("direction = ?", filter.Direction)
	}
	if filter.EntityType != "" {
		query = query.Where("entity_type = ?", filter.EntityType)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	query, err := utils.Paginate(&domain.SyncLog{}, pagination, query)
	if err != nil {
		return nil, err
	}

	if err := query.Find(&logs).Error; err != nil {
		return nil, err
	}

	return logs, nil
}
