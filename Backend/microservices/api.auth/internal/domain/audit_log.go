package domain

import (
	"time"

	"lumini-hub/common/utils"
)

// AuditLogFilterRequest define os parâmetros da busca paginada da auditoria (POST /audit-logs/filter).
// Search procura na ação e no nome do usuário; DateFrom e DateTo são instantes (RFC 3339), o
// fim já exclusivo (o frontend manda o início do dia seguinte).
type AuditLogFilterRequest struct {
	Search     string     `json:"search" binding:"omitempty,max=100"`
	EntityType string     `json:"entity_type" binding:"omitempty,max=50"`
	Method     string     `json:"method" binding:"omitempty,oneof=POST PUT PATCH DELETE"`
	UserID     uint       `json:"user_id"`
	DateFrom   *time.Time `json:"date_from"`
	DateTo     *time.Time `json:"date_to"`

	PageNo        *int   `json:"page_no"`
	PageSize      *int   `json:"page_size" binding:"omitempty,max=100"`
	OrderByColumn string `json:"order_by_column" binding:"omitempty,max=50"`
	IsAsc         *bool  `json:"is_asc"`
}

// InGetAuditLogsFilters são os filtros já tratados que o repository aplica.
type InGetAuditLogsFilters struct {
	Search     string
	EntityType string
	Method     string
	UserID     uint
	DateFrom   *time.Time
	DateTo     *time.Time
}

// ApiAuditLog é um registro de auditoria como o frontend o recebe.
type ApiAuditLog struct {
	ID         uint      `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	UserID     *uint     `json:"user_id"`
	Username   string    `json:"username"`
	Action     string    `json:"action"`
	EntityType string    `json:"entity_type"`
	EntityID   string    `json:"entity_id"`
	IPAddress  string    `json:"ip_address"`
	Status     int       `json:"status"`
	// Detalhes do que mudou: "changes" (campo -> de/para, numa alteração), "after" (objeto criado),
	// "before" (objeto excluído) ou "truncated" (campos de um conteúdo grande demais para guardar).
	Details map[string]any `json:"details,omitempty"`
}

// ApiAuditLogListPaginated é a página de registros devolvida pela busca.
type ApiAuditLogListPaginated struct {
	AuditLogs  []ApiAuditLog       `json:"data"`
	Pagination utils.ApiPagination `json:"pagination"`
}
