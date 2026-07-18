package repository

import (
	"lumini-hub/api.integrations/internal/domain"
	commonrepo "lumini-hub/common/repository"
	"lumini-hub/common/utils"

	"gorm.io/gorm"
)

// WebhookEventRepository define as operações de acesso a dados para eventos de webhook
type WebhookEventRepository interface {
	commonrepo.Repository[domain.WebhookEvent]
	FindByFilter(filter domain.WebhookEventFilterRequest, pagination *utils.Pagination) ([]domain.WebhookEvent, error)
}

// GormWebhookEventRepository implementa WebhookEventRepository usando GORM
type GormWebhookEventRepository struct {
	*commonrepo.GormRepository[domain.WebhookEvent]
}

// NewWebhookEventRepository cria um novo repository de eventos de webhook
func NewWebhookEventRepository(db *gorm.DB) WebhookEventRepository {
	return &GormWebhookEventRepository{
		GormRepository: commonrepo.NewGormRepository[domain.WebhookEvent](db),
	}
}

// FindByFilter busca eventos de webhook aplicando filtros dinâmicos e paginação
func (r *GormWebhookEventRepository) FindByFilter(filter domain.WebhookEventFilterRequest, pagination *utils.Pagination) ([]domain.WebhookEvent, error) {
	var events []domain.WebhookEvent

	query := r.GetDB().Model(&domain.WebhookEvent{})

	if filter.Source != "" {
		query = query.Where("source = ?", filter.Source)
	}
	if filter.EventType != "" {
		query = query.Where("event_type = ?", filter.EventType)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	query, err := utils.Paginate(&domain.WebhookEvent{}, pagination, query)
	if err != nil {
		return nil, err
	}

	if err := query.Find(&events).Error; err != nil {
		return nil, err
	}

	return events, nil
}
