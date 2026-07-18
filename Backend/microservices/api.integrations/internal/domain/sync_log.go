package domain

import (
	"time"

	"lumini-hub/common/database"
	"lumini-hub/common/utils"

	"gorm.io/gorm"
)

// Direções e status possíveis de um SyncLog
const (
	SyncDirectionLiToSql = "li_to_sql"
	SyncDirectionSqlToLi = "sql_to_li"

	SyncStatusSuccess = "success"
	SyncStatusError   = "error"
	SyncStatusPending = "pending"
)

// SyncLog registra o resultado de uma operação de sincronização entre a Loja Integrada e o SQL Server legado
type SyncLog struct {
	gorm.Model

	Direction   string         `gorm:"size:20;not null" json:"direction"`
	EntityType  string         `gorm:"size:50;not null" json:"entity_type"`
	ReferenceID string         `gorm:"size:100" json:"reference_id"`
	Status      string         `gorm:"size:20;not null" json:"status"`
	Message     string         `gorm:"size:500" json:"message"`
	Details     database.JSONB `gorm:"type:jsonb" json:"details"`
}

// TableName especifica o nome da tabela
func (SyncLog) TableName() string {
	return "sync_logs"
}

// SyncLogFilterRequest define os parâmetros para busca filtrada e paginada de sync logs via HTTP POST
type SyncLogFilterRequest struct {
	Direction  string `json:"direction"`
	EntityType string `json:"entity_type"`
	Status     string `json:"status"`

	PageNo        *int   `json:"page_no"`
	PageSize      *int   `json:"page_size"`
	OrderByColumn string `json:"order_by_column"`
	IsAsc         *bool  `json:"is_asc"`
}

// ApiSyncLog representa um SyncLog para exibição via API
type ApiSyncLog struct {
	ID          uint           `json:"id"`
	Direction   string         `json:"direction"`
	EntityType  string         `json:"entity_type"`
	ReferenceID string         `json:"reference_id"`
	Status      string         `json:"status"`
	Message     string         `json:"message"`
	Details     database.JSONB `json:"details"`
	CreatedAt   time.Time      `json:"created_at"`
}

// ApiSyncLogListPaginated representa uma lista paginada de sync logs
type ApiSyncLogListPaginated struct {
	SyncLogs   []ApiSyncLog        `json:"data"`
	Pagination utils.ApiPagination `json:"pagination"`
}

// ApiSyncLogFromModel converte um SyncLog para ApiSyncLog
func ApiSyncLogFromModel(s SyncLog) ApiSyncLog {
	return ApiSyncLog{
		ID:          s.ID,
		Direction:   s.Direction,
		EntityType:  s.EntityType,
		ReferenceID: s.ReferenceID,
		Status:      s.Status,
		Message:     s.Message,
		Details:     s.Details,
		CreatedAt:   s.CreatedAt,
	}
}
