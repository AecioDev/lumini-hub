package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"lumini-hub/api.core/internal/repository"
	"lumini-hub/api.core/internal/routes"
	"lumini-hub/common/config"
	"lumini-hub/common/database"

	"github.com/gin-gonic/gin"
)

func main() {
	// Carregar configurações usando a fonte de ambiente do common
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Erro ao carregar configurações: %v", err)
	}

	// Inicializar banco de dados compartilhado do common
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}

	// Configurar engine do Gin
	router := gin.Default()

	// Grupo base da API
	api := router.Group("/api")

	// Inicializar o Unit of Work
	uow := repository.NewUnitOfWork(db)

	// Configurar rotas específicas de clientes/fornecedores
	routes.SetupRoutes(api, uow, cfg)

	// O microsserviço Core escuta na porta 4002
	port := "4002"
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Canal para capturar sinais de interrupção
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Iniciar servidor em uma goroutine
	go func() {
		log.Printf("Microsserviço api.core iniciado na porta %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erro ao iniciar servidor api.core: %v", err)
		}
	}()

	// Aguardar sinal de interrupção
	<-quit
	log.Println("Desligando servidor api.core...")

	// Contexto com timeout para desligamento gracioso
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Erro ao desligar servidor api.core: %v", err)
	}

	log.Println("Servidor api.core desligado com sucesso")
}
