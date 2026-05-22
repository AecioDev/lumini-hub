package seeders

import (
	"log"

	"simple-erp-service/internal/data-structure/models"

	"gorm.io/gorm"
)

func SeedMeasurementUnit(db *gorm.DB) {

	units := []models.MeasurementUnit{
		{Name: "Unidade", Abbreviation: "un"},
		{Name: "Caixa", Abbreviation: "cx"},
		{Name: "Pacote", Abbreviation: "pct"},
		{Name: "Quilograma", Abbreviation: "kg"},
		{Name: "Litro", Abbreviation: "l"},
		{Name: "Metro", Abbreviation: "m"},
	}

	for _, unit := range units {
		if err := db.Where("abbreviation = ?", unit.Abbreviation).FirstOrCreate(&unit).Error; err != nil {
			log.Printf("Erro ao inserir estado %s: %v", unit.Name, err)
		}
	}
}
