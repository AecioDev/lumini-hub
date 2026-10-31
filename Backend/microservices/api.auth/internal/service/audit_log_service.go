package service

import (
	"lumini-hub/api.auth/internal/domain"
	"lumini-hub/api.auth/internal/repository"
	"lumini-hub/common/database"
	"lumini-hub/common/utils"
)

// AuditLogService consulta a trilha de auditoria (somente leitura).
type AuditLogService struct {
	repo repository.AuditLogRepository
}

// NewAuditLogService cria um novo serviço de auditoria
func NewAuditLogService(repo repository.AuditLogRepository) *AuditLogService {
	return &AuditLogService{repo: repo}
}

// GetAuditLogs devolve uma página da auditoria. O que foi feito por usuários do Perfil DEVELOP só
// aparece para o próprio DEVELOP, como o resto do que é dele (ver users e roles).
func (s *AuditLogService) GetAuditLogs(filter domain.AuditLogFilterRequest, requesterRole string) (*domain.ApiAuditLogListPaginated, error) {
	pagination := utils.PaginationFromFilter(filter.PageNo, filter.PageSize, filter.OrderByColumn, filter.IsAsc, "id")

	hideRole := ""
	if requesterRole != utils.RoleDeveloper {
		hideRole = utils.RoleDeveloper
	}
	logs, err := s.repo.FindAllFiltered(&pagination, domain.InGetAuditLogsFilters{
		Search:     filter.Search,
		EntityType: filter.EntityType,
		Method:     filter.Method,
		UserID:     filter.UserID,
		DateFrom:   filter.DateFrom,
		DateTo:     filter.DateTo,
	}, hideRole)
	if err != nil {
		return nil, err
	}

	items := make([]domain.ApiAuditLog, 0, len(logs))
	for _, entry := range logs {
		items = append(items, apiAuditLogFromModel(entry))
	}
	return &domain.ApiAuditLogListPaginated{
		AuditLogs:  items,
		Pagination: *utils.ApiPaginationFromModel(&pagination),
	}, nil
}

// GetEntityTypes devolve os tipos de entidade presentes na auditoria.
func (s *AuditLogService) GetEntityTypes() ([]string, error) {
	types, err := s.repo.FindEntityTypes()
	if types == nil {
		types = []string{}
	}
	return types, err
}

func apiAuditLogFromModel(entry database.SystemLog) domain.ApiAuditLog {
	item := domain.ApiAuditLog{
		ID:         entry.ID,
		CreatedAt:  entry.CreatedAt,
		UserID:     entry.UserID,
		Action:     entry.Action,
		EntityType: entry.EntityType,
		EntityID:   entry.EntityID,
		IPAddress:  entry.IPAddress,
	}
	if username, ok := entry.Details["username"].(string); ok {
		item.Username = username
	}
	for _, key := range []string{"changes", "after", "before", "truncated"} {
		if value, ok := entry.Details[key]; ok {
			if item.Details == nil {
				item.Details = map[string]any{}
			}
			item.Details[key] = value
		}
	}
	// O status volta do jsonb como float64.
	if status, ok := entry.Details["status"].(float64); ok {
		item.Status = int(status)
	}
	return item
}
