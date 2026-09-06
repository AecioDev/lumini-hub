package domain

import (
	"time"

	"gorm.io/gorm"
)

// CompanyVisualConfig representa a identidade visual de uma Company (1:1):
// logo e paleta de cores, aplicadas no ConfigProvider do antd em runtime.
// Ver Documentos/Planejamento/Modulo_0_Configuracao/plano_empresa.md pra
// decisões fechadas (logo em bytea na própria linha, cores em hex,
// fallback pra paleta padrão da Lumini Hub se a empresa não configurar
// nada — esse fallback é responsabilidade do consumidor/frontend, não
// deste struct).
type CompanyVisualConfig struct {
	gorm.Model

	// Unique parcial (só entre linhas não excluídas) — mesmo raciocínio de
	// CompanyFiscalConfig.CompanyID: sem o `where`, o soft delete do GORM
	// deixaria uma linha excluída ainda ocupando o índice único.
	CompanyID uint `gorm:"not null;uniqueIndex:idx_company_visual_configs_company_id,where:deleted_at IS NULL" json:"company_id"`

	// LogoFile e LogoMimeType nunca são expostos em JSON direto — ver
	// ApiCompanyVisualConfig.HasLogo. Upload é endpoint multipart separado
	// do PUT de cores, não passa por aqui como Create/UpdateRequest (mesmo
	// padrão do certificado em CompanyFiscalConfig).
	LogoFile     []byte `gorm:"type:bytea" json:"-"`
	LogoMimeType string `gorm:"size:50" json:"-"`

	// PrimaryColor é obrigatória (vira colorPrimary do ConfigProvider);
	// SecondaryColor/AccentColor são opcionais, uso pontual. Validação do
	// formato hex é responsabilidade do validator, não deste struct.
	PrimaryColor   string  `gorm:"size:9;not null" json:"primary_color"`
	SecondaryColor *string `gorm:"size:9" json:"secondary_color"`
	AccentColor    *string `gorm:"size:9" json:"accent_color"`

	// TextColor é opcional — cor do texto sobre elementos na cor primária
	// (cabeçalho, botões). Achado testando: texto fixo em branco não fica
	// legível quando a cor primária escolhida é clara (ex.: amarelo).
	// Nulo = consumidor usa branco como fallback (mesmo espírito do
	// fallback de paleta padrão da Lumini Hub).
	TextColor *string `gorm:"size:9" json:"text_color"`
}

// TableName especifica o nome da tabela
func (CompanyVisualConfig) TableName() string {
	return "company_visual_configs"
}

// CreateCompanyVisualConfigRequest representa os dados para criar a configuração visual de uma empresa
type CreateCompanyVisualConfigRequest struct {
	CompanyID      uint   `json:"company_id" binding:"required"`
	PrimaryColor   string `json:"primary_color" binding:"required"`
	SecondaryColor string `json:"secondary_color"`
	AccentColor    string `json:"accent_color"`
	TextColor      string `json:"text_color"`
}

// UpdateCompanyVisualConfigRequest representa os dados para atualizar a configuração visual de uma empresa
type UpdateCompanyVisualConfigRequest struct {
	PrimaryColor   string `json:"primary_color" binding:"required"`
	SecondaryColor string `json:"secondary_color"`
	AccentColor    string `json:"accent_color"`
	TextColor      string `json:"text_color"`
}

// ApiCompanyVisualConfig representa os dados de configuração visual para exibição
type ApiCompanyVisualConfig struct {
	ID             uint      `json:"id"`
	CompanyID      uint      `json:"company_id"`
	HasLogo        bool      `json:"has_logo"`
	PrimaryColor   string    `json:"primary_color"`
	SecondaryColor *string   `json:"secondary_color"`
	AccentColor    *string   `json:"accent_color"`
	TextColor      *string   `json:"text_color"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ApiCompanyVisualConfigDetail representa os dados detalhados de configuração visual
type ApiCompanyVisualConfigDetail = ApiCompanyVisualConfig

// ApiCompanyVisualConfigFromModel converte um CompanyVisualConfig para ApiCompanyVisualConfig
func ApiCompanyVisualConfigFromModel(c CompanyVisualConfig) ApiCompanyVisualConfig {
	return ApiCompanyVisualConfig{
		ID:             c.ID,
		CompanyID:      c.CompanyID,
		HasLogo:        len(c.LogoFile) > 0,
		PrimaryColor:   c.PrimaryColor,
		SecondaryColor: c.SecondaryColor,
		AccentColor:    c.AccentColor,
		TextColor:      c.TextColor,
		CreatedAt:      c.CreatedAt,
		UpdatedAt:      c.UpdatedAt,
	}
}
