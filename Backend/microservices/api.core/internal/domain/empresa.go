package domain

import (
	"time"

	"gorm.io/gorm"
)

// Empresa representa uma empresa (Matriz ou vinculada) dentro do tenant.
// Self-referencing: a Matriz é a empresa com ParentID == nil (só pode haver
// uma por tenant, garantido em validator.EmpresaValidator). Não é
// multi-tenant — multi-tenant é banco físico separado por cliente da Lumini
// Hub; isso é a estrutura organizacional DENTRO do banco de um único
// cliente (ver Documentos/Planejamento/Modulo_0_Configuracao/plano_empresa.md).
type Empresa struct {
	gorm.Model

	ParentID     *uint     `json:"parent_id"`
	Children     []Empresa `gorm:"foreignKey:ParentID" json:"-"`
	RazaoSocial  string    `gorm:"size:150;not null" json:"razao_social"`
	NomeFantasia string    `gorm:"size:150" json:"nome_fantasia"`
	CNPJ         string    `gorm:"size:20;not null;unique" json:"cnpj"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`

	CreatedByID *uint `gorm:"column:created_by" json:"created_by_id"`
	UpdatedByID *uint `gorm:"column:updated_by" json:"updated_by_id"`
}

// TableName especifica o nome da tabela
func (Empresa) TableName() string {
	return "empresas"
}

// CreateEmpresaRequest representa os dados para criar uma nova empresa
type CreateEmpresaRequest struct {
	ParentID     *uint  `json:"parent_id"`
	RazaoSocial  string `json:"razao_social" binding:"required,min=3"`
	NomeFantasia string `json:"nome_fantasia"`
	CNPJ         string `json:"cnpj" binding:"required"`
}

// UpdateEmpresaRequest representa os dados para atualizar uma empresa
type UpdateEmpresaRequest struct {
	ParentID     *uint  `json:"parent_id"`
	RazaoSocial  string `json:"razao_social" binding:"required,min=3"`
	NomeFantasia string `json:"nome_fantasia"`
	CNPJ         string `json:"cnpj" binding:"required"`
	// Sem `binding:"required"` de propósito: num bool, "required" exige true
	// (zero value é false), o que impediria desativar uma empresa via PUT.
	IsActive bool `json:"is_active"`
}

// ApiEmpresa representa os dados de empresa para exibição em lista
type ApiEmpresa struct {
	ID           uint      `json:"id"`
	ParentID     *uint     `json:"parent_id"`
	RazaoSocial  string    `json:"razao_social"`
	NomeFantasia string    `json:"nome_fantasia"`
	CNPJ         string    `json:"cnpj"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ApiEmpresaDetail representa os dados detalhados de uma empresa
type ApiEmpresaDetail struct {
	ID           uint      `json:"id"`
	ParentID     *uint     `json:"parent_id"`
	RazaoSocial  string    `json:"razao_social"`
	NomeFantasia string    `json:"nome_fantasia"`
	CNPJ         string    `json:"cnpj"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	CreatedBy    *ApiUser  `json:"created_by,omitempty"`
	UpdatedAt    time.Time `json:"updated_at"`
	UpdatedBy    *ApiUser  `json:"updated_by,omitempty"`
}

// ApiEmpresaList representa a lista (não paginada, ver nota em EmpresaRepository) de empresas
type ApiEmpresaList struct {
	Empresas []ApiEmpresa `json:"data"`
}

// ApiEmpresaFromModel converte um Empresa para ApiEmpresa
func ApiEmpresaFromModel(e Empresa) ApiEmpresa {
	return ApiEmpresa{
		ID:           e.ID,
		ParentID:     e.ParentID,
		RazaoSocial:  e.RazaoSocial,
		NomeFantasia: e.NomeFantasia,
		CNPJ:         e.CNPJ,
		IsActive:     e.IsActive,
		CreatedAt:    e.CreatedAt,
		UpdatedAt:    e.UpdatedAt,
	}
}

// ApiEmpresaDetailFromModel converte um Empresa para ApiEmpresaDetail
func ApiEmpresaDetailFromModel(e Empresa) ApiEmpresaDetail {
	dto := ApiEmpresaDetail{
		ID:           e.ID,
		ParentID:     e.ParentID,
		RazaoSocial:  e.RazaoSocial,
		NomeFantasia: e.NomeFantasia,
		CNPJ:         e.CNPJ,
		IsActive:     e.IsActive,
		CreatedAt:    e.CreatedAt,
		UpdatedAt:    e.UpdatedAt,
	}

	if e.CreatedByID != nil && *e.CreatedByID > 0 {
		dto.CreatedBy = &ApiUser{ID: *e.CreatedByID}
	}
	if e.UpdatedByID != nil && *e.UpdatedByID > 0 {
		dto.UpdatedBy = &ApiUser{ID: *e.UpdatedByID}
	}

	return dto
}
