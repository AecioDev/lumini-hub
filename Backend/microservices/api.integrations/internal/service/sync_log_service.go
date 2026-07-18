package service

import (
	"lumini-hub/api.integrations/internal/domain"
	"lumini-hub/api.integrations/internal/repository"
	"lumini-hub/common/utils"
)

// SyncLogService expõe consultas sobre os logs de sincronização
type SyncLogService struct {
	uow repository.UnitOfWork
}

// NewSyncLogService cria um novo SyncLogService
func NewSyncLogService(uow repository.UnitOfWork) *SyncLogService {
	return &SyncLogService{uow: uow}
}

// GetSyncLogsByFilter retorna logs de sincronização com filtros e paginação (PostFilter)
func (s *SyncLogService) GetSyncLogsByFilter(filter domain.SyncLogFilterRequest) (*domain.ApiSyncLogListPaginated, error) {
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

	logs, err := s.uow.SyncLogs().FindByFilter(filter, &pagination)
	if err != nil {
		return nil, err
	}

	logDTOs := make([]domain.ApiSyncLog, 0, len(logs))
	for _, log := range logs {
		logDTOs = append(logDTOs, domain.ApiSyncLogFromModel(log))
	}

	return &domain.ApiSyncLogListPaginated{
		SyncLogs:   logDTOs,
		Pagination: *utils.ApiPaginationFromModel(&pagination),
	}, nil
}
