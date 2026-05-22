package routes

import (
	"simple-erp-service/config"
	"simple-erp-service/internal/api/handlers"
	"simple-erp-service/internal/api/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupCustomersRoutes configura as rotas de customers
func SetupCustomersRoutes(router *gin.RouterGroup, db *gorm.DB) {
	customerHandler := handlers.NewCustomerHandler(db)

	// Obter configuração para middleware de autenticação
	cfg, _ := config.Load()

	// Grupo de rotas de usuários (todas protegidas)
	customers := router.Group("/customers")
	customers.Use(middlewares.AuthMiddleware(cfg))
	{
		customers.GET("", middlewares.RequirePermission("customers.view"), customerHandler.GetCustomers)
		customers.GET("/:id", middlewares.RequirePermission("customers.view"), customerHandler.GetCustomer)
		customers.POST("", middlewares.RequirePermission("customers.create"), customerHandler.CreateCustomer)
		customers.PUT("/:id", middlewares.RequirePermission("customers.edit"), customerHandler.UpdateCustomer)
		customers.DELETE("/:id", middlewares.RequirePermission("customers.delete"), customerHandler.DeleteCustomer)
	}
}
