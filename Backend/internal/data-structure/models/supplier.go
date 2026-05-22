package models

import (
	"time"

	"gorm.io/gorm"
)

// Supplier representa um fornecedor
type Supplier struct {
	gorm.Model

	FirstName      string `gorm:"size:100;" json:"first_name"`
	LastName       string `gorm:"size:100;" json:"last_name"`
	PersonType     string `gorm:"default:F" json:"person_type"`
	DocumentNumber string `gorm:"size:20;unique" json:"document_number"`
	CompanyName    string `gorm:"size:100;" json:"company_name"`
	IsActive       bool   `gorm:"default:true" json:"is_active"`

	CreatedByID *uint `gorm:"column:created_by" json:"created_by"`
	CreatedBy   *User `gorm:"foreignKey:CreatedByID" json:"created_by_user,omitempty"`

	UpdatedByID *uint `gorm:"column:updated_by" json:"updated_by"`
	UpdatedBy   *User `gorm:"foreignKey:UpdatedByID" json:"updated_by_user,omitempty"`

	Notes     string `gorm:"size:255" json:"notes"`
	CreatedAt time.Time
	UpdatedAt time.Time

	// Relacionamentos
	Documents    []Document    `gorm:"foreignKey:SupplierID" json:"-"`
	Addresses    []Address     `gorm:"foreignKey:SupplierID" json:"-"`
	Contacts     []Contact     `gorm:"foreignKey:SupplierID" json:"-"`
	Purchases    []Purchase    `gorm:"foreignKey:SupplierID" json:"-"` // Compras do Fornecedor
	Transactions []Transaction `gorm:"foreignKey:SupplierID" json:"-"`
}

// TableName especifica o nome da tabela
func (Supplier) TableName() string {
	return "suppliers"
}

// CreateSupplierRequest representa os dados para criar um novo fornecedor
type CreateSupplierRequest struct {
	FirstName      string `json:"first_name" binding:"required"`
	LastName       string `json:"last_name" binding:"omitempty"`
	PersonType     string `json:"person_type" binding:"required"` // Solicitar Atenção no Front pq não poderá alterar mais depois
	DocumentNumber string `json:"document_number" binding:"required"`
	CompanyName    string `json:"company_name" binding:"omitempty"`
	IsActive       bool   `json:"is_active" binding:"required"`
	Notes          string `json:"notes"`
	CreatedByID    uint   `json:"created_by" binding:"required"`
}

// UpdateSupplierRequest representa os dados para atualizar um fornecedor
type UpdateSupplierRequest struct {
	FirstName      string `json:"first_name" binding:"required"`
	LastName       string `json:"last_name" binding:"omitempty"`
	DocumentNumber string `json:"document_number" binding:"required"`
	CompanyName    string `json:"company_name" binding:"omitempty"`
	IsActive       bool   `json:"is_active" binding:"required"`
	Notes          string `json:"notes"`

	// Relacionamentos O Front manda uma Lista de IDs válidos para cada entidade relacionada.
	DocumentsIDs []uint `json:"documents_ids" binding:"omitempty"`
	AdressesIDs  []uint `json:"adresses_ids" binding:"omitempty"`
	ContactsIDs  []uint `json:"contacts_ids" binding:"omitempty"`
}
