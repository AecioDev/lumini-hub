package routes

import (
	"lumini-hub/api.integrations/internal/handlers"
	"lumini-hub/api.integrations/internal/repository"
	"lumini-hub/api.integrations/internal/service"
	"lumini-hub/common/config"
	"lumini-hub/common/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRoutes registra as rotas do api.integrations. sqlServerDB pode ser nil se a conexão com
// o ERP legado estiver indisponível — os endpoints que dependem dela retornam 503 nesse caso.
func SetupRoutes(router *gin.RouterGroup, uow repository.UnitOfWork, configService *service.ConfigService, sqlServerDB *gorm.DB, cfg *config.Config) {
	configHandler := handlers.NewConfigHandler(configService)
	syncLogHandler := handlers.NewSyncLogHandler(uow)
	webhookEventHandler := handlers.NewWebhookEventHandler(uow)
	legacyLookupHandler := handlers.NewLegacyLookupHandler(sqlServerDB)
	webhookHandler := handlers.NewWebhookHandler(service.NewWebhookService(uow, configService))

	integrations := router.Group("/integrations")

	protected := integrations.Group("")
	protected.Use(middlewares.AuthMiddleware(cfg))
	{
		protected.GET("/settings", middlewares.RequirePermission("integrations.view"), configHandler.GetSettings)
		protected.PUT("/settings", middlewares.RequirePermission("integrations.edit"), configHandler.UpdateSettings)

		protected.POST("/sync-logs/filter", middlewares.RequirePermission("integrations.view"), syncLogHandler.PostSyncLogFilter)
		protected.POST("/webhook-events/filter", middlewares.RequirePermission("integrations.view"), webhookEventHandler.PostWebhookEventFilter)

		protected.GET("/legacy/locations", middlewares.RequirePermission("integrations.view"), legacyLookupHandler.GetLocations)
		protected.GET("/legacy/companies", middlewares.RequirePermission("integrations.view"), legacyLookupHandler.GetCompanies)
	}

	// Grupo público (sem AuthMiddleware) — autenticado por segredo compartilhado dentro do handler
	public := integrations.Group("/webhooks")
	{
		public.POST("/loja-integrada", webhookHandler.ReceiveLojaIntegradaWebhook)
	}
}
