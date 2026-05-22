// @title Simple ERP Service API
// @version 1.0
// @description API para o sistema Simple ERP Service
// @host localhost:4000
// @BasePath /api

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Descrição: Insira seu token JWT no formato "Bearer TOKEN"

package main

import (
	"log"

	"simple-erp-service/config"
	"simple-erp-service/internal/api/server"
	"simple-erp-service/internal/repository/db"
)

func main() {
	// Carregar configurações
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Erro ao carregar configurações: %v", err)
	}

	// Inicializar banco de dados
	database, err := db.InitDB(cfg)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}

	// Inicializar e executar o servidor
	s := server.NewServer(cfg, database)
	if err := s.Run(); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
