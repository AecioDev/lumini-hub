package seeders

import (
	"log"

	"simple-erp-service/internal/data-structure/models"

	"gorm.io/gorm"
)

func SeedCountries(db *gorm.DB) {
	countries := []models.Country{
		{Name: "Brasil", BacenCode: "1058", PhoneCode: "55"},
	}

	for _, country := range countries {
		if err := db.Where("bacen_code = ?", country.BacenCode).FirstOrCreate(&country).Error; err != nil {
			log.Printf("Erro ao inserir país %s: %v", country.Name, err)
		}
	}
}
