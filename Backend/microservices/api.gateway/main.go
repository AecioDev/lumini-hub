package main

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	_ "lumini-hub/api.gateway/docs"
	"lumini-hub/common/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Lumini Hub ERP - API Gateway
// @version         1.0
// @description     Gateway unificado de entrada do ERP Lumini Hub, gerenciando roteamento de microsserviços.
// @host            localhost:4000
// @BasePath        /api
// @securityDefinitions.apikey ApiKeyAuth
// @in              header
// @name            Authorization
func main() {
	// Carregar configurações do common
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Erro ao carregar configurações no Gateway: %v", err)
	}

	// Inicializar engine do Gin
	router := gin.Default()

	// Configurar CORS no Gateway para permitir requisições do frontend e cookies HTTP-Only
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3001"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "Cookie"},
		ExposeHeaders:    []string{"Content-Length", "Set-Cookie"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Definir as URLs dos microsserviços
	authURL, err := url.Parse("http://localhost:4001")
	if err != nil {
		log.Fatalf("URL do serviço de autenticação inválida: %v", err)
	}

	coreURL, err := url.Parse("http://localhost:4002")
	if err != nil {
		log.Fatalf("URL do serviço core inválida: %v", err)
	}

	// Criar proxies
	authProxy := httputil.NewSingleHostReverseProxy(authURL)
	coreProxy := httputil.NewSingleHostReverseProxy(coreURL)

	// Customizar o comportamento dos proxies para garantir compatibilidade
	adjustProxy := func(proxy *httputil.ReverseProxy, target *url.URL) {
		originalDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			originalDirector(req)
			// Modificar o Host do request para corresponder ao destino
			req.Host = target.Host
		}
	}
	adjustProxy(authProxy, authURL)
	adjustProxy(coreProxy, coreURL)

	// Rota do Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Middleware de proxy usando NoRoute para tratar todas as requisições não mapeadas pelo roteador local
	router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		// Roteamento baseado em prefixo
		if strings.HasPrefix(path, "/api/auth") ||
			strings.HasPrefix(path, "/api/users") ||
			strings.HasPrefix(path, "/api/roles") ||
			strings.HasPrefix(path, "/api/permissions") {
			log.Printf("[Gateway] Proxying %s -> api.auth (%s)", path, authURL.String())
			authProxy.ServeHTTP(c.Writer, c.Request)
			return
		}

		if strings.HasPrefix(path, "/api/customers") ||
			strings.HasPrefix(path, "/api/suppliers") {
			log.Printf("[Gateway] Proxying %s -> api.core (%s)", path, coreURL.String())
			coreProxy.ServeHTTP(c.Writer, c.Request)
			return
		}

		// Rota não encontrada
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Serviço ou rota não encontrada no gateway",
			"path":    path,
		})
	})

	// O Gateway escuta na porta original 4000
	port := "4000"
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
		log.Printf("API Gateway iniciado com sucesso na porta %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erro ao iniciar o API Gateway: %v", err)
		}
	}()

	// Aguardar sinal de interrupção
	<-quit
	log.Println("Desligando API Gateway...")

	// Contexto com timeout para desligamento gracioso
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Erro ao desligar o API Gateway: %v", err)
	}

	log.Println("API Gateway desligado com sucesso")
}
