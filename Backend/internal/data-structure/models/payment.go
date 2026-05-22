package models

import (
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	gorm.Model

	Amount      float64   `gorm:"not null" json:"amount"`
	Currency    string    `gorm:"size:3;not null" json:"currency"`
	PaymentDate time.Time `gorm:"not null" json:"payment_date"`
	Description *string   `json:"description"`

	// Relacionamentos
	PaymentMethodID uint          `gorm:"not null" json:"payment_method_id"`
	PaymentMethod   PaymentMethod `gorm:"foreignKey:PaymentMethodID" json:"payment_method"`

	TransactionID uint        `json:"transaction_id"`
	Transaction   Transaction `gorm:"foreignKey:TransactionID" json:"transaction"`

	AccountID uint    `gorm:"not null" json:"account_id"`
	Account   Account `gorm:"foreignKey:AccountID" json:"account"`
}
