package routes

import (
	"github.com/gin-gonic/gin"
	"simple-erp-service/config"
	"simple-erp-service/internal/api/handlers"
	"simple-erp-service/internal/api/middlewares"
	"gorm.io/gorm"
)

// SetupSupplierRoutes configura as rotas de fornecedores
func SetupSupplierRoutes(router *gin.RouterGroup, db *gorm.DB) {
	supplierHandler := handlers.NewSupplierHandler(db)

	// Obter configuração para middleware de autenticação
	cfg, _ := config.Load()

	// Grupo de rotas de fornecedores (todas protegidas)
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