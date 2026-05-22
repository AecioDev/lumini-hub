package seeders

import (
	"log"

	"simple-erp-service/internal/data-structure/models"

	"gorm.io/gorm"
)

func SeedStates(db *gorm.DB) {
	var country models.Country
	if err := db.Where("bacen_code = ?", "1058").First(&country).Error; err != nil {
		log.Fatal("País Brasil não encontrado, execute SeedCountries primeiro!")
	}

	states := []models.State{
		{Name: "Acre", UF: "AC", IBGECode: "12", CountryID: country.ID},
		{Name: "Alagoas", UF: "AL", IBGECode: "27", CountryID: country.ID},
		{Name: "Amapá", UF: "AP", IBGECode: "16", CountryID: country.ID},
		{Name: "Amazonas", UF: "AM", IBGECode: "13", CountryID: country.ID},
		{Name: "Bahia", UF: "BA", IBGECode: "29", CountryID: country.ID},
		{Name: "Ceará", UF: "CE", IBGECode: "23", CountryID: country.ID},
		{Name: "Distrito Federal", UF: "DF", IBGECode: "53", CountryID: country.ID},
		{Name: "Espírito Santo", UF: "ES", IBGECode: "32", CountryID: country.ID},
		{Name: "Goiás", UF: "GO", IBGECode: "52", CountryID: country.ID},
		{Name: "Maranhão", UF: "MA", IBGECode: "21", CountryID: country.ID},
		{Name: "Mato Grosso", UF: "MT", IBGECode: "51", CountryID: country.ID},
		{Name: "Mato Grosso do Sul", UF: "MS", IBGECode: "50", CountryID: country.ID},
		{Name: "Minas Gerais", UF: "MG", IBGECode: "31", CountryID: country.ID},
		{Name: "Pará", UF: "PA", IBGECode: "15", CountryID: country.ID},
		{Name: "Paraíba", UF: "PB", IBGECode: "25", CountryID: country.ID},
		{Name: "Paraná", UF: "PR", IBGECode: "41", CountryID: country.ID},
		{Name: "Pernambuco", UF: "PE", IBGECode: "26", CountryID: country.ID},
		{Name: "Piauí", UF: "PI", IBGECode: "22", CountryID: country.ID},
		{Name: "Rio de Janeiro", UF: "RJ", IBGECode: "33", CountryID: country.ID},
		{Name: "Rio Grande do Norte", UF: "RN", IBGECode: "24", CountryID: country.ID},
		{Name: "Rio Grande do Sul", UF: "RS", IBGECode: "43", CountryID: country.ID},
		{Name: "Rondônia", UF: "RO", IBGECode: "11", CountryID: country.ID},
		{Name: "Roraima", UF: "RR", IBGECode: "14", CountryID: country.ID},
		{Name: "Santa Catarina", UF: "SC", IBGECode: "42", CountryID: country.ID},
		{Name: "São Paulo", UF: "SP", IBGECode: "35", CountryID: country.ID},
		{Name: "Sergipe", UF: "SE", IBGECode: "28", CountryID: country.ID},
		{Name: "Tocantins", UF: "TO", IBGECode: "17", CountryID: country.ID},
	}

	for _, state := range states {
		if err := db.Where("ibge_code = ?", state.IBGECode).FirstOrCreate(&state).Error; err != nil {
			log.Printf("Erro ao inserir estado %s: %v", state.Name, err)
		}
	}
}
