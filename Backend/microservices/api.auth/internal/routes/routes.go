package routes

import (
	"time"

	"lumini-hub/api.auth/internal/handlers"
	"lumini-hub/common/config"
	"lumini-hub/common/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRoutes configura todas as rotas do microsserviço de autenticação e perfis
func SetupRoutes(router *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
	// Handlers
	authHandler := handlers.NewAuthHandler(db, cfg)
	passwordResetHandler := handlers.NewPasswordResetHandler(db, cfg)
	userHandler := handlers.NewUserHandler(db)
	roleHandler := handlers.NewRoleHandler(db)
	permHandler := handlers.NewPermissionHandler(db)
	menuItemHandler := handlers.NewMenuItemHandler(db)
	auditHandler := handlers.NewAuditLogHandler(db)

	// Rotas de Autenticação
	auth := router.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh-token", authHandler.RefreshToken)

		// Recuperação de senha por e-mail (públicas): freio por IP aqui, e o limite por usuário fica no service.
		auth.POST("/forgot-password", middlewares.RateLimitByIP(5, time.Hour), passwordResetHandler.ForgotPassword)
		auth.POST("/reset-password", middlewares.RateLimitByIP(10, 15*time.Minute), passwordResetHandler.ResetPassword)

		// Rotas de autenticação protegidas
		protected := auth.Group("")
		protected.Use(middlewares.AuthMiddleware(cfg))
		{
			protected.POST("/logout", authHandler.Logout)
			protected.GET("/me", authHandler.GetMe)
			protected.GET("/version", authHandler.GetVersion)
			protected.PUT("/active-company", authHandler.SetActiveCompany)
			protected.PUT("/favorites", authHandler.SetMenuFavorites)
		}
	}

	// Rotas de Usuários (todas protegidas)
	users := router.Group("/users")
	users.Use(middlewares.AuthMiddleware(cfg))
	{
		users.GET("", middlewares.RequirePermission("users.view"), userHandler.GetUsers)
		users.POST("/filter", middlewares.RequirePermission("users.view"), userHandler.PostUserFilter)
		users.GET("/:id", middlewares.RequirePermission("users.view"), userHandler.GetUser)
		users.POST("", middlewares.RequirePermission("users.create"), userHandler.CreateUser)
		users.PUT("/:id", middlewares.RequirePermission("users.edit"), userHandler.UpdateUser)
		users.DELETE("/:id", middlewares.RequirePermission("users.delete"), userHandler.DeleteUser)
		users.PUT("/:id/password", userHandler.ChangePassword) // A permissão é verificada internamente no handler
		users.PUT("/:id/permissions", middlewares.RequirePermission("users.edit"), userHandler.UpdateUserPermissions)
	}

	// Rotas de Perfis/Roles (todas protegidas)
	roles := router.Group("/roles")
	roles.Use(middlewares.AuthMiddleware(cfg))
	{
		roles.GET("", middlewares.RequirePermission("roles.view"), roleHandler.GetRoles)
		roles.GET("/:id", middlewares.RequirePermission("roles.view"), roleHandler.GetRole)
		roles.POST("", middlewares.RequirePermission("roles.create"), roleHandler.CreateRole)
		roles.PUT("/:id", middlewares.RequirePermission("roles.edit"), roleHandler.UpdateRole)
		roles.PUT("/:id/permissions", middlewares.RequirePermission("roles.edit"), roleHandler.UpdateRolePermissions)
		roles.DELETE("/:id", middlewares.RequirePermission("roles.delete"), roleHandler.DeleteRole)
	}

	// Rotas de Permissões (todas protegidas)
	permissions := router.Group("/permissions")
	permissions.Use(middlewares.AuthMiddleware(cfg))
	{
		permissions.GET("", middlewares.RequirePermission("permissions.view"), permHandler.GetPermissions)
		permissions.POST("/filter", middlewares.RequirePermission("permissions.view"), permHandler.PostPermissionFilter)
		permissions.GET("/:id", middlewares.RequirePermission("permissions.view"), permHandler.GetPermission)
		permissions.GET("/by-module", middlewares.RequirePermission("permissions.view"), permHandler.GetPermissionsByModule)
		permissions.GET("/modules", middlewares.RequirePermission("permissions.view"), permHandler.GetAvailableModules)
		permissions.POST("", middlewares.RequirePermission("permissions.create"), permHandler.CreatePermission)
		permissions.PUT("/:id", middlewares.RequirePermission("permissions.edit"), permHandler.UpdatePermission)
		permissions.DELETE("/:id", middlewares.RequirePermission("permissions.delete"), permHandler.DeletePermission)
	}

	// Rotas de Itens de Menu (todas protegidas, dev-only)
	// Auditoria (somente leitura): quem fez qual escrita e quando
	audit := router.Group("/audit-logs")
	audit.Use(middlewares.AuthMiddleware(cfg))
	{
		audit.POST("/filter", middlewares.RequirePermission("audit.view"), auditHandler.PostAuditLogFilter)
		audit.GET("/entity-types", middlewares.RequirePermission("audit.view"), auditHandler.GetAuditEntityTypes)
	}

	menuItems := router.Group("/menu-items")
	menuItems.Use(middlewares.AuthMiddleware(cfg))
	{
		menuItems.GET("", middlewares.RequirePermission("admin.create_permissions"), menuItemHandler.GetMenuTree)
		menuItems.GET("/:id", middlewares.RequirePermission("admin.create_permissions"), menuItemHandler.GetMenuItem)
		menuItems.POST("", middlewares.RequirePermission("admin.create_permissions"), menuItemHandler.CreateMenuItem)
		menuItems.PUT("/:id", middlewares.RequirePermission("admin.create_permissions"), menuItemHandler.UpdateMenuItem)
		menuItems.DELETE("/:id", middlewares.RequirePermission("admin.create_permissions"), menuItemHandler.DeleteMenuItem)
	}
}
