package domain

import "gorm.io/gorm"

// IntegrationConfig representa uma configuração chave/valor usada pelo api.integrations
// (chaves da Loja Integrada, códigos do ERP legado, etc.)
type IntegrationConfig struct {
	gorm.Model

	Key         string `gorm:"size:100;unique;not null" json:"key"`
	Value       string `gorm:"type:text" json:"value"`
	Description string `gorm:"size:255" json:"description"`
}

// TableName especifica o nome da tabela
func (IntegrationConfig) TableName() string {
	return "integration_configs"
}

// UpdateIntegrationSettingsRequest representa o upsert em lote das configurações
type UpdateIntegrationSettingsRequest struct {
	LiApiKey         *string `json:"li_api_key"`
	LiAppKey         *string `json:"li_app_key"`
	LiWebhookSecret  *string `json:"li_webhook_secret"`
	CodTipNot        *string `json:"codtipnot"`
	CodLocArmOficial *string `json:"codlocarm_oficial"`
	CodLocArmReserva *string `json:"codlocarm_reserva"`
	CodEmp           *string `json:"codemp"`
}

// ApiIntegrationSettings representa as configurações expostas via API
type ApiIntegrationSettings struct {
	LiApiKey         string `json:"li_api_key"`
	LiAppKey         string `json:"li_app_key"`
	LiWebhookSecret  string `json:"li_webhook_secret"`
	CodTipNot        string `json:"codtipnot"`
	CodLocArmOficial string `json:"codlocarm_oficial"`
	CodLocArmReserva string `json:"codlocarm_reserva"`
	CodEmp           string `json:"codemp"`
}
