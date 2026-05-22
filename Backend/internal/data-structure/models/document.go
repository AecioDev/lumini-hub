package models

import (
	"time"

	"gorm.io/gorm"
)

type Document struct {
	gorm.Model

	CustomerID *uint     `gorm:"index" json:"customer_id,omitempty"`
	Customer   *Customer `gorm:"foreignKey:CustomerID" json:"-"`

	SupplierID *uint     `gorm:"index" json:"supplier_id,omitempty"`
	Supplier   *Supplier `gorm:"foreignKey:SupplierID" json:"-"`

	Type         string     `gorm:"not null" json:"type"`            // Tipo de : CPF, RG, CNPJ, IE, IM, CNH etc.
	Number       string     `gorm:"unique;not null" json:"number"`   // Número do o
	Validate     *time.Time `json:"validate"`                        // Data de Validade (nullable)
	EmissionDate *time.Time `json:"emission_date"`                   // Data de Emissão (nullable)
	Department   string     `json:"department"`                      // Órgão Emissor
	StateID      uint       `json:"state_id"`                        // FK para State (UF de Emissão)
	State        State      `gorm:"foreignKey:StateID" json:"state"` // Associação com a UF
}

// TableName especifica o nome da tabela
func (Document) TableName() string {
	return "document"
}
