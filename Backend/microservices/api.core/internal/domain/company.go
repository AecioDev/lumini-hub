package domain

import (
	"time"

	"gorm.io/gorm"
)

// Company representa uma empresa (Matriz ou vinculada) dentro do tenant.
// Self-referencing: uma "Matriz" é simplesmente qualquer empresa com
// ParentID == nil — o tenant pode ter vários grupos empresariais
// independentes, sem holding entre si (decidido 2026-09-05, ver
// plano_empresa.md). Não é multi-tenant — multi-tenant é banco físico
// separado por cliente da Lumini Hub; isso é a estrutura organizacional
// DENTRO do banco de um único cliente.
type Company struct {
	gorm.Model

	ParentID  *uint     `json:"parent_id"`
	Children  []Company `gorm:"foreignKey:ParentID" json:"-"`
	LegalName string    `gorm:"size:150;not null" json:"legal_name"`
	TradeName string    `gorm:"size:150" json:"trade_name"`
	// TaxID (CNPJ) é obrigatório só pra quem é Matriz (ParentID == nil) — uma
	// empresa vinculada pode não ter CNPJ próprio quando a responsabilidade
	// fiscal fica com a Matriz (ela existe só pra personalizar relatórios:
	// logo, endereço etc., caso real do usuário, decidido 2026-09-05).
	// Nullable (não string vazia) de propósito: várias empresas sem CNPJ não
	// podem colidir no índice único (NULL nunca é igual a NULL pro Postgres,
	// "" seria igual a "" e quebraria a segunda empresa sem CNPJ criada).
	TaxID    *string `gorm:"size:20" json:"tax_id"`
	IsActive bool    `gorm:"default:true" json:"is_active"`

	CreatedByID *uint `gorm:"column:created_by" json:"created_by_id"`
	UpdatedByID *uint `gorm:"column:updated_by" json:"updated_by_id"`
}

// TableName especifica o nome da tabela
func (Company) TableName() string {
	return "companies"
}

// CreateCompanyRequest representa os dados para criar uma nova empresa.
// TaxID sem `binding:"required"` de propósito — só é obrigatório pra quem é
// Matriz (ParentID == nil), checado em CompanyValidator, não aqui (uma
// empresa vinculada pode não ter CNPJ próprio).
type CreateCompanyRequest struct {
	ParentID  *uint  `json:"parent_id"`
	LegalName string `json:"legal_name" binding:"required,min=3"`
	TradeName string `json:"trade_name"`
	TaxID     string `json:"tax_id"`
}

// UpdateCompanyRequest representa os dados para atualizar uma empresa
type UpdateCompanyRequest struct {
	ParentID  *uint  `json:"parent_id"`
	LegalName string `json:"legal_name" binding:"required,min=3"`
	TradeName string `json:"trade_name"`
	TaxID     string `json:"tax_id"`
	// Sem `binding:"required"` de propósito: num bool, "required" exige true
	// (zero value é false), o que impediria desativar uma empresa via PUT.
	IsActive bool `json:"is_active"`
}

// ApiCompany representa os dados de empresa para exibição em lista
type ApiCompany struct {
	ID        uint      `json:"id"`
	ParentID  *uint     `json:"parent_id"`
	LegalName string    `json:"legal_name"`
	TradeName string    `json:"trade_name"`
	TaxID     *string   `json:"tax_id"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ApiCompanyDetail representa os dados detalhados de uma empresa
type ApiCompanyDetail struct {
	ID        uint      `json:"id"`
	ParentID  *uint     `json:"parent_id"`
	LegalName string    `json:"legal_name"`
	TradeName string    `json:"trade_name"`
	TaxID     *string   `json:"tax_id"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy *ApiUser  `json:"created_by,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy *ApiUser  `json:"updated_by,omitempty"`
}

// ApiCompanyList representa a lista (não paginada, ver nota em CompanyRepository) de empresas
type ApiCompanyList struct {
	Companies []ApiCompany `json:"data"`
}

// ApiCompanyFromModel converte um Company para ApiCompany
func ApiCompanyFromModel(e Company) ApiCompany {
	return ApiCompany{
		ID:        e.ID,
		ParentID:  e.ParentID,
		LegalName: e.LegalName,
		TradeName: e.TradeName,
		TaxID:     e.TaxID,
		IsActive:  e.IsActive,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

// ApiCompanyDetailFromModel converte um Company para ApiCompanyDetail
func ApiCompanyDetailFromModel(e Company) ApiCompanyDetail {
	dto := ApiCompanyDetail{
		ID:        e.ID,
		ParentID:  e.ParentID,
		LegalName: e.LegalName,
		TradeName: e.TradeName,
		TaxID:     e.TaxID,
		IsActive:  e.IsActive,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}

	if e.CreatedByID != nil && *e.CreatedByID > 0 {
		dto.CreatedBy = &ApiUser{ID: *e.CreatedByID}
	}
	if e.UpdatedByID != nil && *e.UpdatedByID > 0 {
		dto.UpdatedBy = &ApiUser{ID: *e.UpdatedByID}
	}

	return dto
}
