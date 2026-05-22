package domain

import (
	"time"

	"lumini-hub/common/utils"

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
	Notes          string `gorm:"size:255" json:"notes"`

	CreatedByID *uint `gorm:"column:created_by" json:"created_by"`
	UpdatedByID *uint `gorm:"column:updated_by" json:"updated_by"`

	// Relacionamentos
	Documents []Document `gorm:"foreignKey:SupplierID" json:"-"`
	Addresses []Address  `gorm:"foreignKey:SupplierID" json:"-"`
	Contacts  []Contact  `gorm:"foreignKey:SupplierID" json:"-"`
}

// TableName especifica o nome da tabela
func (Supplier) TableName() string {
	return "suppliers"
}

// CreateSupplierRequest representa os dados para criar um novo fornecedor
type CreateSupplierRequest struct {
	FirstName      string `json:"first_name" binding:"required"`
	LastName       string `json:"last_name" binding:"omitempty"`
	PersonType     string `json:"person_type" binding:"required"`
	DocumentNumber string `json:"document_number" binding:"required"`
	CompanyName    string `json:"company_name" binding:"omitempty"`
	IsActive       bool   `json:"is_active" binding:"required"`
	Notes          string `json:"notes"`
	CreatedByID    uint   `json:"created_by"`
}

// UpdateSupplierRequest representa os dados para atualizar um fornecedor
type UpdateSupplierRequest struct {
	FirstName      string `json:"first_name" binding:"required"`
	LastName       string `json:"last_name" binding:"omitempty"`
	DocumentNumber string `json:"document_number" binding:"required"`
	CompanyName    string `json:"company_name" binding:"omitempty"`
	IsActive       bool   `json:"is_active" binding:"required"`
	Notes          string `json:"notes"`

	// Relacionamentos
	DocumentsIDs []uint `json:"documents_ids" binding:"omitempty"`
	AdressesIDs  []uint `json:"adresses_ids" binding:"omitempty"`
	ContactsIDs  []uint `json:"contacts_ids" binding:"omitempty"`
}

// ApiSupplier representa os dados de fornecedor para exibição
type ApiSupplier struct {
	ID             uint       `json:"id"`
	FirstName      string     `json:"first_name"`
	LastName       string     `json:"last_name"`
	PersonType     string     `json:"person_type"`
	DocumentNumber string     `json:"document_number"`
	CompanyName    string     `json:"company_name"`
	IsActive       bool       `json:"is_active"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	Addresses      *[]Address `json:"addresses,omitempty"`
	Contacts       *[]Contact `json:"contacts,omitempty"`
}

// ApiSupplierDetail representa os dados detalhados de fornecedor para exibição
type ApiSupplierDetail struct {
	ID             uint       `json:"id"`
	FirstName      string     `json:"first_name"`
	LastName       string     `json:"last_name"`
	PersonType     string     `json:"person_type"`
	DocumentNumber string     `json:"document_number"`
	CompanyName    string     `json:"company_name"`
	IsActive       bool       `json:"is_active"`
	Addresses      *[]Address `json:"addresses,omitempty"`
	Contacts       *[]Contact `json:"contacts,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	CreatedBy      *ApiUser   `json:"created_by,omitempty"`
	UpdatedBy      *ApiUser   `json:"updated_by,omitempty"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// ApiSupplierListPaginated representa uma lista paginada de fornecedores
type ApiSupplierListPaginated struct {
	Supplier   []ApiSupplier        `json:"data"`
	Pagination utils.ApiPagination `json:"pagination"`
}

// ApiSupplierFromModel converte um Supplier para ApiSupplier
func ApiSupplierFromModel(s Supplier) ApiSupplier {
	dto := ApiSupplier{
		ID:             s.ID,
		FirstName:      s.FirstName,
		LastName:       s.LastName,
		PersonType:     s.PersonType,
		DocumentNumber: s.DocumentNumber,
		CompanyName:    s.CompanyName,
		IsActive:       s.IsActive,
		CreatedAt:      s.CreatedAt,
		UpdatedAt:      s.UpdatedAt,
	}

	if len(s.Addresses) > 0 {
		dto.Addresses = &s.Addresses
	}
	if len(s.Contacts) > 0 {
		dto.Contacts = &s.Contacts
	}

	return dto
}

// ApiSupplierDetailFromModel converte um Supplier para ApiSupplierDetail
func ApiSupplierDetailFromModel(s Supplier) ApiSupplierDetail {
	dto := ApiSupplierDetail{
		ID:             s.ID,
		FirstName:      s.FirstName,
		LastName:       s.LastName,
		PersonType:     s.PersonType,
		DocumentNumber: s.DocumentNumber,
		CompanyName:    s.CompanyName,
		IsActive:       s.IsActive,
		CreatedAt:      s.CreatedAt,
		UpdatedAt:      s.UpdatedAt,
	}

	if len(s.Addresses) > 0 {
		dto.Addresses = &s.Addresses
	}
	if len(s.Contacts) > 0 {
		dto.Contacts = &s.Contacts
	}

	if s.CreatedByID != nil && *s.CreatedByID > 0 {
		dto.CreatedBy = &ApiUser{
			ID: *s.CreatedByID,
		}
	}
	if s.UpdatedByID != nil && *s.UpdatedByID > 0 {
		dto.UpdatedBy = &ApiUser{
			ID: *s.UpdatedByID,
		}
	}

	return dto
}
