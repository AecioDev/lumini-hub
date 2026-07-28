package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"lumini-hub/api.auth/internal/domain"
	"lumini-hub/api.auth/internal/routes"
	"lumini-hub/api.auth/internal/seeder"
	"lumini-hub/common/config"
	"lumini-hub/common/database"
	"lumini-hub/common/middlewares"

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
	middlewares.InitPermissionChecker(db)

	// Migração pontual da tabela de itens de menu (demais tabelas do api.auth
	// não usam AutoMigrate hoje; escopo restrito só a esta struct nova)
	if err := db.AutoMigrate(&domain.MenuItem{}); err != nil {
		log.Fatalf("Erro ao migrar tabela de itens de menu: %v", err)
	}
	if err := seeder.SeedMenuItems(db); err != nil {
		log.Fatalf("Erro ao semear itens de menu: %v", err)
	}

	// Configurar engine do Gin
	router := gin.Default()

	// Grupo base da API
	api := router.Group("/api")

	// Configurar rotas específicas de auth/users/roles/permissions
	routes.SetupRoutes(api, db, cfg)

	// O microsserviço de autenticação deve escutar explicitamente na porta 4001
	port := "4001"
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
		log.Printf("Microsserviço api.auth iniciado na porta %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erro ao iniciar servidor api.auth: %v", err)
		}
	}()

	// Aguardar sinal de interrupção
	<-quit
	log.Println("Desligando servidor api.auth...")

	// Contexto com timeout para desligamento gracioso
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Erro ao desligar servidor api.auth: %v", err)
	}

	log.Println("Servidor api.auth desligado com sucesso")
}
