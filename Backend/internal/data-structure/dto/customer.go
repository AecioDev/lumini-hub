package dto

import (
	"simple-erp-service/internal/data-structure/models"
	"time"
)

// ApiCustomer representa os dados de cliente para exibição
type ApiCustomer struct {
	ID             uint              `json:"id"`
	FirstName      string            `json:"first_name"`
	LastName       string            `json:"last_name"`
	PersonType     string            `json:"person_type"`
	DocumentNumber string            `json:"document_number"`
	CompanyName    string            `json:"company_name"`
	IsActive       bool              `json:"is_active"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
	Addresses      *[]models.Address `json:"addresses"`
	Contacts       *[]models.Contact `json:"contacts"`
}

// ApiCustomerDetail representa os dados detalhados de cliente para exibição
type ApiCustomerDetail struct {
	ID             uint   `json:"id"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	PersonType     string `json:"person_type"`
	DocumentNumber string `json:"document_number"`
	CompanyName    string `json:"company_name"`
	IsActive       bool   `json:"is_active"`

	Addresses *[]models.Address `json:"addresses"`
	Contacts  *[]models.Contact `json:"contacts"`

	CreatedAt time.Time `json:"created_at"`
	CreatedBy *ApiUser  `json:"created_by"`

	UpdatedBy *ApiUser  `json:"updated_by"`
	UpdatedAt time.Time `json:"updated_at"`

	// Sales        *[]ApiSale        `json:"sales"`
	// Transactions *[]ApiTransaction `json:"transactions"`
}

// CustomerListDTO representa uma lista paginada de clientes
type ApiCustomerListPaginated struct {
	Customer   []ApiCustomer `json:"data"`
	Pagination ApiPagination `json:"pagination"`
}

// ApiCustomerFromModel converte um Customer para ApiCustomer
func ApiCustomerFromModel(c models.Customer) ApiCustomer {
	return ApiCustomer{
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
}

// ApiCustomerDetailFromModel converte um Customer para ApiCustomerDetail
func ApiCustomerDetailFromModel(c models.Customer) ApiCustomerDetail {
	return ApiCustomerDetail{
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
}
