package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"lumini-hub/api.auth/internal/domain"
	"lumini-hub/api.auth/internal/routes"
	"lumini-hub/api.auth/internal/seeder"
	"lumini-hub/api.auth/seed"
	"lumini-hub/common/config"
	"lumini-hub/common/database"
	"lumini-hub/common/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

	// Migração do schema e seed de dados base (ver prepareDatabase)
	if err := prepareDatabase(db, cfg); err != nil {
		log.Fatalf("Erro ao preparar o banco de dados: %v", err)
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

// ATENÇÃO: prepareDatabase fica neste arquivo de propósito. Os serviços são iniciados com
// `go run main.go` (run_services.bat), que compila SÓ o main.go: uma função em outro arquivo
// do pacote daria "undefined" nesse comando (e passaria despercebido em `go build .`).
// prepareDatabase deixa o banco pronto pro api.auth subir. Tudo roda sob
// DB_AUTO_MIGRATE (ligado por padrão pra um banco vazio subir sozinho): com a
// flag desligada — produção com dados reais — o serviço não migra NEM semeia
// nada ao subir, e o operador assume schema e dados (o catálogo de permissões
// precisa então ser aplicado por outro meio, ex.: rodar uma vez com a flag
// ligada). Na ordem:
//
//  1. AutoMigrate das tabelas. Ordem de dependência: Permission e Role antes de
//     User (FK role_id e as tabelas de junção role_permissions/user_permissions,
//     criadas pelo GORM a partir dos many2many) e de MenuItem (FK permission_id).
//  2. seed.SyncCatalog: catálogo de permissões + role ADMIN (permissão nova no
//     código entra na tabela e no ADMIN no boot seguinte).
//  3. seed.BootstrapAdmin: o primeiro usuário ADMIN, a partir de BOOTSTRAP_ADMIN_*, só
//     com a tabela users vazia e depois do SyncCatalog (a role ADMIN precisa existir).
//  4. seeder.SeedMenuItems: depois do catálogo, porque o menu referencia as
//     permissões pelo código; com elas ausentes o item nasceria sem permissão
//     própria, e isso nunca é corrigido depois (itens existentes são pulados).
func prepareDatabase(db *gorm.DB, cfg *config.Config) error {
	if !cfg.Database.AutoMigrate {
		return nil
	}

	if err := db.AutoMigrate(
		&domain.Permission{},
		&domain.Role{},
		&domain.User{},
		&domain.MenuItem{},
	); err != nil {
		return fmt.Errorf("migrando tabelas do api.auth: %w", err)
	}
	if err := seed.SyncCatalog(db); err != nil {
		return fmt.Errorf("semeando permissões e a role %s: %w", seed.AdminRoleName, err)
	}
	if _, err := seed.BootstrapAdmin(db, seed.AdminBootstrap{
		Username: cfg.Bootstrap.Username,
		Password: cfg.Bootstrap.Password,
		Email:    cfg.Bootstrap.Email,
	}); err != nil {
		return fmt.Errorf("criando o admin inicial: %w", err)
	}
	if err := seeder.SeedMenuItems(db); err != nil {
		return fmt.Errorf("semeando itens de menu: %w", err)
	}
	return nil
}
