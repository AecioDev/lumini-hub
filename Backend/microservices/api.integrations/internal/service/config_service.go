package service

import (
	"sync"

	"lumini-hub/api.integrations/internal/domain"
	"lumini-hub/api.integrations/internal/repository"
)

// Chaves conhecidas de configuração do api.integrations
const (
	KeyLiApiKey         = "li_api_key"
	KeyLiAppKey         = "li_app_key"
	KeyLiWebhookSecret  = "li_webhook_secret"
	KeyCodTipNot        = "codtipnot"
	KeyCodLocArmOficial = "codlocarm_oficial"
	KeyCodLocArmReserva = "codlocarm_reserva"
	KeyCodEmp           = "codemp"
)

var defaultConfigKeys = []string{
	KeyLiApiKey, KeyLiAppKey, KeyLiWebhookSecret,
	KeyCodTipNot, KeyCodLocArmOficial, KeyCodLocArmReserva, KeyCodEmp,
}

// ConfigService mantém em memória as configurações do api.integrations (chaves da Loja
// Integrada e códigos do ERP legado), evitando ir ao banco a cada requisição.
type ConfigService struct {
	repo repository.IntegrationConfigRepository

	mu   sync.RWMutex
	data map[string]string
}

// NewConfigService cria um novo ConfigService
func NewConfigService(repo repository.IntegrationConfigRepository) *ConfigService {
	return &ConfigService{repo: repo, data: make(map[string]string)}
}

// Reload recarrega o cache em memória a partir do banco de dados
func (s *ConfigService) Reload() error {
	entries, err := s.repo.FindAll()
	if err != nil {
		return err
	}

	data := make(map[string]string, len(entries))
	for _, entry := range entries {
		data[entry.Key] = entry.Value
	}

	s.mu.Lock()
	s.data = data
	s.mu.Unlock()

	return nil
}

// Get retorna o valor em cache de uma chave (string vazia se não existir)
func (s *ConfigService) Get(key string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data[key]
}

// SeedDefaults garante que todas as chaves conhecidas existam no banco (com valor vazio se
// ainda não configuradas), para que a tela de configurações sempre tenha o que exibir/editar.
func (s *ConfigService) SeedDefaults() error {
	s.mu.RLock()
	missing := make([]string, 0)
	for _, key := range defaultConfigKeys {
		if _, ok := s.data[key]; !ok {
			missing = append(missing, key)
		}
	}
	s.mu.RUnlock()

	if len(missing) == 0 {
		return nil
	}

	for _, key := range missing {
		if err := s.repo.Create(&domain.IntegrationConfig{Key: key, Value: ""}); err != nil {
			return err
		}
	}

	return s.Reload()
}

// GetSettings monta o DTO de configurações a partir do cache em memória
func (s *ConfigService) GetSettings() domain.ApiIntegrationSettings {
	return domain.ApiIntegrationSettings{
		LiApiKey:         s.Get(KeyLiApiKey),
		LiAppKey:         s.Get(KeyLiAppKey),
		LiWebhookSecret:  s.Get(KeyLiWebhookSecret),
		CodTipNot:        s.Get(KeyCodTipNot),
		CodLocArmOficial: s.Get(KeyCodLocArmOficial),
		CodLocArmReserva: s.Get(KeyCodLocArmReserva),
		CodEmp:           s.Get(KeyCodEmp),
	}
}

// UpdateSettings faz o upsert em lote das chaves informadas (campos nil são ignorados) e
// recarrega o cache ao final.
func (s *ConfigService) UpdateSettings(req domain.UpdateIntegrationSettingsRequest) error {
	updates := map[string]*string{
		KeyLiApiKey:         req.LiApiKey,
		KeyLiAppKey:         req.LiAppKey,
		KeyLiWebhookSecret:  req.LiWebhookSecret,
		KeyCodTipNot:        req.CodTipNot,
		KeyCodLocArmOficial: req.CodLocArmOficial,
		KeyCodLocArmReserva: req.CodLocArmReserva,
		KeyCodEmp:           req.CodEmp,
	}

	for key, value := range updates {
		if value == nil {
			continue
		}
		if err := s.repo.Upsert(key, *value); err != nil {
			return err
		}
	}

	return s.Reload()
}
