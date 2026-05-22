package seeders

import (
	"log"

	"simple-erp-service/internal/data-structure/models"

	"gorm.io/gorm"
)

func SeedProductCategory(db *gorm.DB) {

	// Seed de categorias de produtos
	categories := []models.ProductCategory{
		{Name: "Geral", Description: "Categoria geral de produtos"},
		{Name: "Eletrônicos", Description: "Produtos eletrônicos"},
		{Name: "Alimentos", Description: "Produtos alimentícios"},
		{Name: "Vestuário", Description: "Roupas e acessórios"},
		{Name: "Papelaria", Description: "Material de escritório"},
	}

	for _, category := range categories {
		if err := db.Where("name = ?", category.Name).FirstOrCreate(&category).Error; err != nil {
			log.Printf("Erro ao inserir estado %s: %v", category.Name, err)
		}
	}
}
