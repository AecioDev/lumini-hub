package domain

import (
	"time"

	"lumini-hub/common/database"
	"lumini-hub/common/utils"

	"gorm.io/gorm"
)

// Status possíveis de um WebhookEvent
const (
	WebhookStatusReceived  = "received"
	WebhookStatusProcessed = "processed"
	WebhookStatusError     = "error"
)

// WebhookEvent registra o payload cru recebido de um webhook externo (ex.: Loja Integrada),
// antes de ser processado (Fase 2 processa e vira um pedido/SyncLog).
type WebhookEvent struct {
	gorm.Model

	Source       string         `gorm:"size:50;not null" json:"source"`
	EventType    string         `gorm:"size:50" json:"event_type"`
	Headers      database.JSONB `gorm:"type:jsonb" json:"headers"`
	Payload      database.JSONB `gorm:"type:jsonb" json:"payload"`
	Status       string         `gorm:"size:20;not null" json:"status"`
	ReceivedAt   time.Time      `json:"received_at"`
	ProcessedAt  *time.Time     `json:"processed_at"`
	ErrorMessage string         `gorm:"size:500" json:"error_message"`
}

// TableName especifica o nome da tabela
func (WebhookEvent) TableName() string {
	return "webhook_events"
}

// WebhookEventFilterRequest define os parâmetros para busca filtrada e paginada de webhook events via HTTP POST
type WebhookEventFilterRequest struct {
	Source    string `json:"source"`
	EventType string `json:"event_type"`
	Status    string `json:"status"`

	PageNo        *int   `json:"page_no"`
	PageSize      *int   `json:"page_size"`
	OrderByColumn string `json:"order_by_column"`
	IsAsc         *bool  `json:"is_asc"`
}

// ApiWebhookEvent representa um WebhookEvent para exibição via API
type ApiWebhookEvent struct {
	ID           uint           `json:"id"`
	Source       string         `json:"source"`
	EventType    string         `json:"event_type"`
	Payload      database.JSONB `json:"payload"`
	Status       string         `json:"status"`
	ReceivedAt   time.Time      `json:"received_at"`
	ProcessedAt  *time.Time     `json:"processed_at"`
	ErrorMessage string         `json:"error_message"`
}

// ApiWebhookEventListPaginated representa uma lista paginada de webhook events
type ApiWebhookEventListPaginated struct {
	WebhookEvents []ApiWebhookEvent   `json:"data"`
	Pagination    utils.ApiPagination `json:"pagination"`
}

// ApiWebhookEventFromModel converte um WebhookEvent para ApiWebhookEvent
func ApiWebhookEventFromModel(w WebhookEvent) ApiWebhookEvent {
	return ApiWebhookEvent{
		ID:           w.ID,
		Source:       w.Source,
		EventType:    w.EventType,
		Payload:      w.Payload,
		Status:       w.Status,
		ReceivedAt:   w.ReceivedAt,
		ProcessedAt:  w.ProcessedAt,
		ErrorMessage: w.ErrorMessage,
	}
}
