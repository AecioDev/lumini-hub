package domain

import (
	"time"

	"lumini-hub/common/utils"

	"gorm.io/gorm"
)

// ApiUser representa uma estrutura mínima de usuário para compatibilidade com o frontend
type ApiUser struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
}

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
	UpdatedByID *uint `gorm:"column:updated_by" json:"updated_by_id"`

	// Relacionamentos (não retornados por padrão)
	Documents []Document `gorm:"foreignKey:CustomerID" json:"-"`
	Addresses []Address  `gorm:"foreignKey:CustomerID" json:"-"`
	Contacts  []Contact  `gorm:"foreignKey:CustomerID" json:"-"`
}

// TableName especifica o nome da tabela
func (Customer) TableName() string {
	return "customers"
}

// CreateCustomerRequest representa os dados para criar um novo cliente
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

// ApiCustomer representa os dados de cliente para exibição
type ApiCustomer struct {
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

// ApiCustomerDetail representa os dados detalhados de cliente para exibição
type ApiCustomerDetail struct {
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

// ApiCustomerListPaginated representa uma lista paginada de clientes
type ApiCustomerListPaginated struct {
	Customer   []ApiCustomer        `json:"data"`
	Pagination utils.ApiPagination `json:"pagination"`
}

// ApiCustomerFromModel converte um Customer para ApiCustomer
func ApiCustomerFromModel(c Customer) ApiCustomer {
	dto := ApiCustomer{
		ID:             c.ID,
		FirstName:      c.FirstName,
		LastName:       c.LastName,
		PersonType:     c.PersonType,
		DocumentNumber: c.DocumentNumber,
		CompanyName:    c.CompanyName,
		IsActive:       c.IsActive,
		CreatedAt:      c.CreatedAt,
		UpdatedAt:      c.UpdatedAt,
	}

	if len(c.Addresses) > 0 {
		dto.Addresses = &c.Addresses
	}
	if len(c.Contacts) > 0 {
		dto.Contacts = &c.Contacts
	}

	return dto
}

// ApiCustomerDetailFromModel converte um Customer para ApiCustomerDetail
func ApiCustomerDetailFromModel(c Customer) ApiCustomerDetail {
	dto := ApiCustomerDetail{
		ID:             c.ID,
		FirstName:      c.FirstName,
		LastName:       c.LastName,
		PersonType:     c.PersonType,
		DocumentNumber: c.DocumentNumber,
		CompanyName:    c.CompanyName,
		IsActive:       c.IsActive,
		CreatedAt:      c.CreatedAt,
		UpdatedAt:      c.UpdatedAt,
	}

	if len(c.Addresses) > 0 {
		dto.Addresses = &c.Addresses
	}
	if len(c.Contacts) > 0 {
		dto.Contacts = &c.Contacts
	}

	// Criar instâncias fictícias ou preenchidas apenas com ID para created_by/updated_by
	// já que o relacionamento de banco não existe mais de forma física (com preloads de User).
	// Em produção real, esses dados poderiam vir de uma consulta ao api.auth ou
	// o banco compartilhado pode ser lido caso queiramos fazer um query manual (não preload).
	// Por simplicidade e compatibilidade com o front-end, retornamos structs se os IDs forem válidos.
	if c.CreatedByID != nil && *c.CreatedByID > 0 {
		dto.CreatedBy = &ApiUser{
			ID: *c.CreatedByID,
		}
	}
	if c.UpdatedByID != nil && *c.UpdatedByID > 0 {
		dto.UpdatedBy = &ApiUser{
			ID: *c.UpdatedByID,
		}
	}

	return dto
}
