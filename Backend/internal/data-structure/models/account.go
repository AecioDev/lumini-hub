package models

import "gorm.io/gorm"

type Account struct {
	gorm.Model
	Name          string  `gorm:"not null" json:"name"`              // Nome da conta (ex: Caixa Loja 1, Banco XYZ)
	Type          string  `gorm:"not null" json:"type"`              // Tipo: "cash" para caixa, "bank" para conta bancária
	BankCode      string  `json:"bank_code"`                         // Código do banco (caso seja uma conta bancária)
	Agency        string  `json:"agency"`                            // Agência bancária
	AccountNumber string  `json:"account_number"`                    // Número da conta bancária
	Balance       float64 `gorm:"not null;default:0" json:"balance"` // Saldo da conta
}

// TableName especifica o nome da tabela
func (Account) TableName() string {
	return "account"
}
