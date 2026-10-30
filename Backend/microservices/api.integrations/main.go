package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"lumini-hub/api.integrations/internal/domain"
	"lumini-hub/api.integrations/internal/repository"
	"lumini-hub/api.integrations/internal/routes"
	"lumini-hub/api.integrations/internal/service"
	"lumini-hub/common/config"
	"lumini-hub/common/database"
	"lumini-hub/common/dbversion"
	"lumini-hub/common/middlewares"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Erro ao carregar configurações: %v", err)
	}

	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	middlewares.InitPermissionChecker(db)

	// AutoMigrate das tabelas Postgres do api.integrations (DB_AUTO_MIGRATE,
	// ligado por padrão; em produção com dados reais, desligar).
	if cfg.Database.AutoMigrate {
		if err := db.AutoMigrate(
			&domain.IntegrationConfig{},
			&domain.SyncLog{},
			&domain.ProductMapping{},
			&domain.WebhookEvent{},
		); err != nil {
			log.Fatalf("Erro ao migrar tabelas do api.integrations: %v", err)
		}
	} else if err := dbversion.Check(db); err != nil {
		// Flag desligada: o operador cuida do schema; só confere se o banco é novo o bastante.
		log.Fatalf("Banco incompatível com este sistema: %v", err)
	}

	uow := repository.NewUnitOfWork(db)

	configService := service.NewConfigService(uow.IntegrationConfigs())
	if err := configService.Reload(); err != nil {
		log.Fatalf("Erro ao carregar configurações do api.integrations: %v", err)
	}
	if err := configService.SeedDefaults(); err != nil {
		log.Fatalf("Erro ao inicializar configurações padrão do api.integrations: %v", err)
	}

	router := gin.Default()
	api := router.Group("/api")

	routes.SetupRoutes(api, uow, configService, cfg)

	port := "4007"
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Microsserviço api.integrations iniciado na porta %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erro ao iniciar servidor api.integrations: %v", err)
		}
	}()

	<-quit
	log.Println("Desligando servidor api.integrations...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Erro ao desligar servidor api.integrations: %v", err)
	}
	log.Println("Servidor api.integrations desligado com sucesso")
}
