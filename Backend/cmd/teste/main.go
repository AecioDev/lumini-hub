package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar o arquivo .env:", err)
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL não encontrada no .env")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Erro ao conectar no banco:", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Erro ao obter instancia DB:", err)
	}

	defer sqlDB.Close()

	err = sqlDB.Ping()
	if err != nil {
		log.Fatal("Erro ao dar ping no banco:", err)
	}

	fmt.Println("✅ Conectado com sucesso ao banco da Neon!")
}
