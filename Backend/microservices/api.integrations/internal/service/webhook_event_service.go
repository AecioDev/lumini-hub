package service

import (
	"lumini-hub/api.integrations/internal/domain"
	"lumini-hub/api.integrations/internal/repository"
	"lumini-hub/common/utils"
)

// WebhookEventService expõe consultas sobre os eventos de webhook recebidos
type WebhookEventService struct {
	uow repository.UnitOfWork
}

// NewWebhookEventService cria um novo WebhookEventService
func NewWebhookEventService(uow repository.UnitOfWork) *WebhookEventService {
	return &WebhookEventService{uow: uow}
}

// GetWebhookEventsByFilter retorna eventos de webhook com filtros e paginação (PostFilter)
func (s *WebhookEventService) GetWebhookEventsByFilter(filter domain.WebhookEventFilterRequest) (*domain.ApiWebhookEventListPaginated, error) {
	pageNo := 1
	if filter.PageNo != nil && *filter.PageNo > 0 {
		pageNo = *filter.PageNo
	}
	pageSize := 10
	if filter.PageSize != nil && *filter.PageSize > 0 {
		pageSize = *filter.PageSize
	}

	orderBy := "id desc"
	if filter.OrderByColumn != "" {
		orderBy = filter.OrderByColumn
		isAsc := true
		if filter.IsAsc != nil {
			isAsc = *filter.IsAsc
		}
		if isAsc {
			orderBy += " asc"
		} else {
			orderBy += " desc"
		}
	}

	pagination := utils.Pagination{
		Page:  pageNo,
		Limit: pageSize,
		Sort:  orderBy,
	}

	events, err := s.uow.WebhookEvents().FindByFilter(filter, &pagination)
	if err != nil {
		return nil, err
	}

	eventDTOs := make([]domain.ApiWebhookEvent, 0, len(events))
	for _, event := range events {
		eventDTOs = append(eventDTOs, domain.ApiWebhookEventFromModel(event))
	}

	return &domain.ApiWebhookEventListPaginated{
		WebhookEvents: eventDTOs,
		Pagination:    *utils.ApiPaginationFromModel(&pagination),
	}, nil
}
