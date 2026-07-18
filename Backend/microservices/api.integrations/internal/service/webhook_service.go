package service

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"lumini-hub/api.integrations/internal/domain"
	"lumini-hub/api.integrations/internal/repository"
	"lumini-hub/common/database"
)

// ErrInvalidWebhookSecret é retornado quando o segredo do webhook está configurado e não confere
var ErrInvalidWebhookSecret = errors.New("segredo do webhook inválido")

// WebhookService recebe e persiste eventos de webhook crus (ex.: da Loja Integrada). O
// processamento do evento em um pedido de venda fica para a Fase 2, quando o payload real da
// Loja Integrada e o codtipnot estiverem confirmados.
type WebhookService struct {
	uow           repository.UnitOfWork
	configService *ConfigService
}

// NewWebhookService cria um novo WebhookService
func NewWebhookService(uow repository.UnitOfWork, configService *ConfigService) *WebhookService {
	return &WebhookService{uow: uow, configService: configService}
}

// ReceiveLojaIntegradaWebhook valida o segredo compartilhado (se já configurado) e grava o
// payload cru como um WebhookEvent pendente de processamento.
func (s *WebhookService) ReceiveLojaIntegradaWebhook(headers http.Header, rawBody []byte, providedSecret string) (*domain.WebhookEvent, error) {
	expectedSecret := s.configService.Get(KeyLiWebhookSecret)
	if expectedSecret != "" && providedSecret != expectedSecret {
		return nil, ErrInvalidWebhookSecret
	}

	payload := database.JSONB{}
	if err := json.Unmarshal(rawBody, &payload); err != nil {
		payload = database.JSONB{"raw": string(rawBody)}
	}

	headerSnapshot := database.JSONB{}
	for key, values := range headers {
		if len(values) > 0 {
			headerSnapshot[key] = values[0]
		}
	}

	event := &domain.WebhookEvent{
		Source:     "loja_integrada",
		Headers:    headerSnapshot,
		Payload:    payload,
		Status:     domain.WebhookStatusReceived,
		ReceivedAt: time.Now(),
	}

	err := s.uow.Execute(func(uow repository.UnitOfWork) error {
		return uow.WebhookEvents().Create(event)
	})
	if err != nil {
		return nil, err
	}

	return event, nil
}
