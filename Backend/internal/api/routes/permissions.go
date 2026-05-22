package routes

import (
	"simple-erp-service/config"
	"simple-erp-service/internal/api/handlers"
	"simple-erp-service/internal/api/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupPermissionRoutes configura as rotas de perfis
func SetupPermissionRoutes(router *gin.RouterGroup, db *gorm.DB) {
	PermissionHandler := handlers.NewPermissionHandler(db)

	// Obter configuração para middleware de autenticação
	cfg, _ := config.Load()

	// Grupo de rotas de permissões (todas protegidas)
	permissions := router.Group("/permissions")
	permissions.Use(middlewares.AuthMiddleware(cfg))
	{
		// Rotas de permissões
		permissions.GET("", middlewares.RequirePermission("permissions.view"), PermissionHandler.GetPermissions)
		permissions.GET("/:id", middlewares.RequirePermission("permissions.view"), PermissionHandler.GetPermission)
		permissions.GET("/by-module", middlewares.RequirePermission("permissions.view"), PermissionHandler.GetPermissionsByModule)
		permissions.GET("/modules", middlewares.RequirePermission("permissions.view"), PermissionHandler.GetAvailableModules)
		permissions.POST("", middlewares.RequirePermission("permissions.create"), PermissionHandler.CreatePermission)
		permissions.PUT("/:id", middlewares.RequirePermission("permissions.edit"), PermissionHandler.UpdatePermission)
		permissions.DELETE("/:id", middlewares.RequirePermission("permissions.delete"), PermissionHandler.DeletePermission)
	}
}
