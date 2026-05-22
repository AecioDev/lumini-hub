package seeders

import (
	"log"

	"simple-erp-service/internal/data-structure/models"

	"gorm.io/gorm"
)

func SeedPaymentMethod(db *gorm.DB) {

	// Seed de métodos de pagamento
	paymentMethods := []models.PaymentMethod{
		{Name: "Dinheiro", Description: "Pagamento em espécie", IsActive: true},
		{Name: "Cartão de Crédito", Description: "Pagamento com cartão de crédito", IsActive: true},
		{Name: "Cartão de Débito", Description: "Pagamento com cartão de débito", IsActive: true},
		{Name: "Transferência Bancária", Description: "Pagamento via transferência bancária", IsActive: true},
		{Name: "Pix", Description: "Pagamento via Pix", IsActive: true},
		{Name: "Boleto", Description: "Pagamento via boleto bancário", IsActive: true},
	}

	for _, method := range paymentMethods {
		if err := db.Where("name = ?", method.Name).FirstOrCreate(&method).Error; err != nil {
			log.Printf("Erro ao inserir estado %s: %v", method.Name, err)
		}
	}
}
