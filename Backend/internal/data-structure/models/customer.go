package models

import (
	"gorm.io/gorm"
)

// Customer representa um cliente
type Customer struct {
	gorm.Model

	FirstName      string `gorm:"size:100" json:"first_name"`
	LastName       string `gorm:"size:100" json:"last_name"`
	PersonType     string `gorm:"default:F" json:"person_type"`          // F: Física, J: Jurídica
	DocumentNumber string `gorm:"size:20;unique" json:"document_number"` // CPF ou CNPJ
	CompanyName    string `gorm:"size:100" json:"company_name"`
	IsActive       bool   `gorm:"default:true" json:"is_active"`
	Notes          string `gorm:"size:255" json:"notes"`

	CreatedByID *uint `gorm:"column:created_by" json:"created_by_id"`
	CreatedBy   *User `gorm:"foreignKey:CreatedByID" json:"created_by,omitempty"`

	UpdatedByID *uint `gorm:"column:updated_by" json:"updated_by_id"`
	UpdatedBy   *User `gorm:"foreignKey:UpdatedByID" json:"updated_by,omitempty"`

	// Relacionamentos (não retornados por padrão)
	Documents []Document `gorm:"foreignKey:CustomerID" json:"-"`
	Addresses []Address  `gorm:"foreignKey:CustomerID" json:"-"`
	Contacts  []Contact  `gorm:"foreignKey:CustomerID" json:"-"`
	Sales     []Sale     `gorm:"foreignKey:CustomerID" json:"-"`
}

// TableName especifica o nome da tabela
func (Customer) TableName() string {
	return "customers"
}

// CreateCustomerRequest representa os dados para criar um novo cliente
// Especificar no Front um Estilo Passo a Passo para Gravar o Cliente no Primeiro Passo, e os Demais serem Update com os Demais dados.
type CreateCustomerRequest struct {
	FirstName      string `json:"first_name" binding:"required"`
	LastName       string `json:"last_name" binding:"omitempty"`
	PersonType     string `json:"person_type" binding:"required"`
	DocumentNumber string `json:"document_number" binding:"required"`
	CompanyName    string `json:"company_name" binding:"omitempty"`
	IsActive       bool   `json:"is_active"`
	Notes          string `json:"notes"`
}

// UpdateCustomerRequest representa os dados para atualizar um cliente
type UpdateCustomerRequest struct {
	FirstName      string `json:"first_name" binding:"required"`
	LastName       string `json:"last_name" binding:"omitempty"`
	DocumentNumber string `json:"document_number" binding:"required"`
	CompanyName    string `json:"company_name" binding:"omitempty"`
	IsActive       bool   `json:"is_active" binding:"required"`
	Notes          string `json:"notes"`
}
