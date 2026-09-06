package domain

import (
	"time"

	"gorm.io/gorm"
)

// CompanyColorPalette representa uma paleta de cores salva pelo usuário
// (N:1 com Company) além das paletas fixas do frontend (`PRESET_PALETTES`
// em CompanyVisualConfigForm.tsx). Diferente de CompanyFiscalConfig/
// CompanyVisualConfig, esta é uma lista de verdade — uma empresa pode ter
// zero, uma ou várias paletas salvas — então TEM Delete (o princípio "1:1
// sem Delete" de plano_empresa.md não se aplica aqui). Sem Update de
// propósito: o fluxo é criar uma nova ou apagar uma antiga, não editar
// in-place (decisão de escopo, 2026-09-06).
type CompanyColorPalette struct {
	gorm.Model

	CompanyID uint `gorm:"not null;index" json:"company_id"`

	// Name é como o usuário identifica a paleta na lista (ex: "Campanha de
	// Verão") — sem unicidade forçada, é só um rótulo de conveniência.
	Name string `gorm:"size:100;not null" json:"name"`

	PrimaryColor   string  `gorm:"size:9;not null" json:"primary_color"`
	SecondaryColor *string `gorm:"size:9" json:"secondary_color"`
	AccentColor    *string `gorm:"size:9" json:"accent_color"`
}

// TableName especifica o nome da tabela
func (CompanyColorPalette) TableName() string {
	return "company_color_palettes"
}

// CreateCompanyColorPaletteRequest representa os dados para salvar uma nova paleta personalizada
type CreateCompanyColorPaletteRequest struct {
	CompanyID      uint   `json:"company_id" binding:"required"`
	Name           string `json:"name" binding:"required"`
	PrimaryColor   string `json:"primary_color" binding:"required"`
	SecondaryColor string `json:"secondary_color"`
	AccentColor    string `json:"accent_color"`
}

// ApiCompanyColorPalette representa os dados de uma paleta personalizada para exibição
type ApiCompanyColorPalette struct {
	ID             uint      `json:"id"`
	CompanyID      uint      `json:"company_id"`
	Name           string    `json:"name"`
	PrimaryColor   string    `json:"primary_color"`
	SecondaryColor *string   `json:"secondary_color"`
	AccentColor    *string   `json:"accent_color"`
	CreatedAt      time.Time `json:"created_at"`
}

// ApiCompanyColorPaletteFromModel converte um CompanyColorPalette para ApiCompanyColorPalette
func ApiCompanyColorPaletteFromModel(p CompanyColorPalette) ApiCompanyColorPalette {
	return ApiCompanyColorPalette{
		ID:             p.ID,
		CompanyID:      p.CompanyID,
		Name:           p.Name,
		PrimaryColor:   p.PrimaryColor,
		SecondaryColor: p.SecondaryColor,
		AccentColor:    p.AccentColor,
		CreatedAt:      p.CreatedAt,
	}
}
