package migrations

import (
	"log"
	"simple-erp-service/internal/data-structure/models"

	"gorm.io/gorm"
)

// MigrateDB executa as migrações do banco de dados
func MigrateDB(db *gorm.DB) error {
	log.Println("Iniciando migrações do banco de dados...")

	// Lista de todos os modelos para migração
	models := []interface{}{
		//&models.Account{},

		&models.User{},
		&models.Permission{},
		&models.Role{},

		&models.Country{},
		&models.State{},
		&models.City{},
		&models.Address{},

		&models.Contact{},
		&models.Document{},

		//&models.SaleItem{},
		//&models.Sale{},

		//&models.PurchaseItem{},
		//&models.Purchase{},

		//&models.Transaction{},
		//&models.PaymentMethod{},
		//&models.Payment{},

		&models.Customer{},
		&models.Supplier{},

		// &models.InventoryMovement{},
		// &models.MeasurementUnit{},
		// &models.ProductCategory{},
		// &models.Product{},
		&models.SystemLog{},
	}

	// Executar migrações
	err := db.AutoMigrate(models...)
	if err != nil {
		log.Printf("Erro ao executar migrações: %v", err)
		return err
	}

	log.Println("Migrações concluídas com sucesso!")
	return nil
}
