package models

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model

	Type     string    `gorm:"not null" json:"type"`                    // Valores possíveis: "receivable" (A Receber), "payable" (A Pagar)
	Code     string    `gorm:"unique;not null" json:"transaction_code"` // Código único da transação
	Date     time.Time `gorm:"not null" json:"transaction_date"`        // Data da transação
	Currency string    `gorm:"size:3;not null" json:"currency"`         // Moeda da transação (ex: BRL, USD)
	Amount   float64   `gorm:"not null" json:"amount"`                  // Valor total da transação
	DueDate  time.Time `gorm:"not null" json:"due_date"`                // Data de vencimento
	Fees     float64   `gorm:"default:0" json:"fees"`                   // Taxas de operadoras de cartão
	Interest float64   `gorm:"default:0" json:"interest"`               // Juros por atraso
	Penalty  float64   `gorm:"default:0" json:"penalty"`                // Multa por atraso
	Status   string    `gorm:"not null" json:"status"`                  // "Pendente", "Parcialmente Paga", "Liquidada"
	Notes    *string   `json:"notes"`                                   // Observações

	// Relacionamento com Cliente ou Fornecedor
	CustomerID *uint     `json:"customer_id,omitempty"`
	Customer   *Customer `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`

	SupplierID *uint     `json:"supplier_id,omitempty"`
	Supplier   *Supplier `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`

	SaleID *uint `json:"sale_id"` // Referência a uma venda
	Sale   *Sale `gorm:"foreignKey:SaleID" json:"sale,omitempty"`

	PurchaseID *uint     `json:"purchase_id"` // Referência a uma compra
	Purchase   *Purchase `gorm:"foreignKey:PurchaseID" json:"purchase,omitempty"`

	// Relacionamento com pagamentos
	Payments []Payment `gorm:"foreignKey:TransactionID" json:"payments"`
}
