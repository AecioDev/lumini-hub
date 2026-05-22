package seeders

import (
	"log"

	"simple-erp-service/internal/data-structure/models"

	"gorm.io/gorm"
)

func SeedCities(db *gorm.DB) {
	var sp models.State
	if err := db.Where("uf = ?", "SP").First(&sp).Error; err != nil {
		log.Fatal("Estado de São Paulo não encontrado, execute SeedStates primeiro!")
	}

	cities := []models.City{
		{Name: "São Paulo", IBGECode: "3550308", StateID: sp.ID},
		{Name: "Campinas", IBGECode: "3509502", StateID: sp.ID},
		{Name: "Santos", IBGECode: "3548500", StateID: sp.ID},
	}

	for _, city := range cities {
		if err := db.Where("ibge_code = ?", city.IBGECode).FirstOrCreate(&city).Error; err != nil {
			log.Printf("Erro ao inserir cidade %s: %v", city.Name, err)
		}
	}
}
