package routes

import (
	"simple-erp-service/config"
	"simple-erp-service/internal/api/handlers"
	"simple-erp-service/internal/api/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRoleRoutes configura as rotas de perfis
func SetupRoleRoutes(router *gin.RouterGroup, db *gorm.DB) {
	roleHandler := handlers.NewRoleHandler(db)

	// Obter configuração para middleware de autenticação
	cfg, _ := config.Load()

	// Grupo de rotas de perfis (todas protegidas)
	roles := router.Group("/roles")
	roles.Use(middlewares.AuthMiddleware(cfg))
	{
		// Rotas de perfis
		roles.GET("", middlewares.RequirePermission("users.view"), roleHandler.GetRoles)
		roles.GET("/:id", middlewares.RequirePermission("users.view"), roleHandler.GetRole)
		roles.POST("", middlewares.RequirePermission("users.create"), roleHandler.CreateRole)
		roles.PUT("/:id", middlewares.RequirePermission("users.edit"), roleHandler.UpdateRole)
		roles.PUT("/:id/permissions", middlewares.RequirePermission("users.edit"), roleHandler.UpdateRolePermissions)
		roles.DELETE("/:id", middlewares.RequirePermission("users.delete"), roleHandler.DeleteRole)
	}
}
