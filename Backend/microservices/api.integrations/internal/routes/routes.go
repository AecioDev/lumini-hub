package routes

import (
	"lumini-hub/api.integrations/internal/handlers"
	"lumini-hub/api.integrations/internal/repository"
	"lumini-hub/api.integrations/internal/service"
	"lumini-hub/common/config"
	"lumini-hub/common/middlewares"

	"github.com/gin-gonic/gin"
)

// SetupRoutes registra as rotas do api.integrations.
func SetupRoutes(router *gin.RouterGroup, uow repository.UnitOfWork, configService *service.ConfigService, cfg *config.Config) {
	configHandler := handlers.NewConfigHandler(configService)
	syncLogHandler := handlers.NewSyncLogHandler(uow)
	webhookEventHandler := handlers.NewWebhookEventHandler(uow)
	webhookHandler := handlers.NewWebhookHandler(service.NewWebhookService(uow, configService))

	integrations := router.Group("/integrations")

	protected := integrations.Group("")
	protected.Use(middlewares.AuthMiddleware(cfg))
	{
		protected.GET("/settings", middlewares.RequirePermission("integrations.view"), configHandler.GetSettings)
		protected.PUT("/settings", middlewares.RequirePermission("integrations.edit"), configHandler.UpdateSettings)

		protected.POST("/sync-logs/filter", middlewares.RequirePermission("integrations.view"), syncLogHandler.PostSyncLogFilter)
		protected.POST("/webhook-events/filter", middlewares.RequirePermission("integrations.view"), webhookEventHandler.PostWebhookEventFilter)
	}

	// Grupo público (sem AuthMiddleware) — autenticado por segredo compartilhado dentro do handler
	public := integrations.Group("/webhooks")
	{
		public.POST("/loja-integrada", webhookHandler.ReceiveLojaIntegradaWebhook)
	}
}
