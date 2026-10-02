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

	// AutoMigrate das tabelas do api.auth (DB_AUTO_MIGRATE, ligado por padrão pra
	// um banco vazio subir sozinho; em produção com dados reais, desligar).
	// Ordem de dependência: Permission e Role antes de User (FK role_id e as
	// tabelas de junção role_permissions/user_permissions, criadas pelo GORM a
	// partir dos many2many) e de MenuItem (FK permission_id).
	if cfg.Database.AutoMigrate {
		if err := db.AutoMigrate(
			&domain.Permission{},
			&domain.Role{},
			&domain.User{},
			&domain.MenuItem{},
		); err != nil {
			log.Fatalf("Erro ao migrar tabelas do api.auth: %v", err)
		}
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
