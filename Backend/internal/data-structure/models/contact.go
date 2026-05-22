package models

import (
	"gorm.io/gorm"
)

type Contact struct {
	gorm.Model

	CustomerID *uint     `gorm:"index" json:"customer_id,omitempty"`
	Customer   *Customer `gorm:"foreignKey:CustomerID" json:"-"`

	SupplierID *uint     `gorm:"index" json:"supplier_id,omitempty"`
	Supplier   *Supplier `gorm:"foreignKey:SupplierID" json:"-"`

	Type    string `gorm:"not null" json:"type"`           // Tipo de Contato: email, Telefone, Celular
	Contact string `gorm:"unique;not null" json:"contact"` // Contato em si o email, o telefone etc
	Name    string `json:"name"`                           // Um Nome para contato se for o caso
}

// TableName especifica o nome da tabela
func (Contact) TableName() string {
	return "contact"
}
