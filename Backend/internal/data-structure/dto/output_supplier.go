package dto

import (
	"simple-erp-service/internal/data-structure/models"
	"time"
)

// ApiSupplier representa os dados de fornecedor para exibição
type ApiSupplier struct {
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

// ApiSupplierDetail representa os dados detalhados de fornecedor para exibição
type ApiSupplierDetail struct {
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

	// Purchases    *[]ApiPurchase    `json:"purchases"`
	// Transactions *[]ApiTransaction `json:"transactions"`
}

// SupplierListDTO representa uma lista paginada de fornecedores
type ApiSupplierListPaginated struct {
	Supplier   []ApiSupplier `json:"data"`
	Pagination ApiPagination `json:"pagination"`
}

// ApiSupplierFromModel converte um Supplier para ApiSupplier
func ApiSupplierFromModel(s models.Supplier) ApiSupplier {
	return ApiSupplier{
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
}

// ApiSupplierDetailFromModel converte um Supplier para ApiSupplierDetail
func ApiSupplierDetailFromModel(s models.Supplier) ApiSupplierDetail {
	return ApiSupplierDetail{
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
}
