package routes

import (
	"lumini-hub/api.core/internal/handlers"
	"lumini-hub/common/config"
	"lumini-hub/common/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRoutes configura todas as rotas do microsserviço Core (Clientes e Fornecedores)
func SetupRoutes(router *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
	// Handlers
	customerHandler := handlers.NewCustomerHandler(db)
	supplierHandler := handlers.NewSupplierHandler(db)

	// Rotas de Clientes (todas protegidas)
	customers := router.Group("/customers")
	customers.Use(middlewares.AuthMiddleware(cfg))
	{
		customers.GET("", middlewares.RequirePermission("customers.view"), customerHandler.GetCustomers)
		customers.GET("/:id", middlewares.RequirePermission("customers.view"), customerHandler.GetCustomer)
		customers.POST("", middlewares.RequirePermission("customers.create"), customerHandler.CreateCustomer)
		customers.PUT("/:id", middlewares.RequirePermission("customers.edit"), customerHandler.UpdateCustomer)
		customers.DELETE("/:id", middlewares.RequirePermission("customers.delete"), customerHandler.DeleteCustomer)
	}

	// Rotas de Fornecedores (todas protegidas)
	suppliers := router.Group("/suppliers")
	suppliers.Use(middlewares.AuthMiddleware(cfg))
	{
		suppliers.GET("", middlewares.RequirePermission("suppliers.view"), supplierHandler.GetSuppliers)
		suppliers.GET("/:id", middlewares.RequirePermission("suppliers.view"), supplierHandler.GetSupplier)
		suppliers.POST("", middlewares.RequirePermission("suppliers.create"), supplierHandler.CreateSupplier)
		suppliers.PUT("/:id", middlewares.RequirePermission("suppliers.edit"), supplierHandler.UpdateSupplier)
		suppliers.DELETE("/:id", middlewares.RequirePermission("suppliers.delete"), supplierHandler.DeleteSupplier)
	}
}
