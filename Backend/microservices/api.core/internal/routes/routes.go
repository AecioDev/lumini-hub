package routes

import (
	"lumini-hub/api.core/internal/handlers"
	"lumini-hub/api.core/internal/repository"
	"lumini-hub/common/config"
	"lumini-hub/common/middlewares"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configura todas as rotas do microsserviço Core (Clientes, Fornecedores e Empresas)
func SetupRoutes(router *gin.RouterGroup, uow repository.UnitOfWork, cfg *config.Config) {
	// Handlers
	customerHandler := handlers.NewCustomerHandler(uow)
	supplierHandler := handlers.NewSupplierHandler(uow)
	empresaHandler := handlers.NewEmpresaHandler(uow)

	// Rotas de Clientes (todas protegidas)
	customers := router.Group("/customers")
	customers.Use(middlewares.AuthMiddleware(cfg))
	{
		customers.GET("", middlewares.RequirePermission("customers.view"), customerHandler.GetCustomers)
		customers.POST("/filter", middlewares.RequirePermission("customers.view"), customerHandler.PostCustomerFilter)
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
		suppliers.POST("/filter", middlewares.RequirePermission("suppliers.view"), supplierHandler.PostSupplierFilter)
		suppliers.GET("/:id", middlewares.RequirePermission("suppliers.view"), supplierHandler.GetSupplier)
		suppliers.POST("", middlewares.RequirePermission("suppliers.create"), supplierHandler.CreateSupplier)
		suppliers.PUT("/:id", middlewares.RequirePermission("suppliers.edit"), supplierHandler.UpdateSupplier)
		suppliers.DELETE("/:id", middlewares.RequirePermission("suppliers.delete"), supplierHandler.DeleteSupplier)
	}

	// Rotas de Empresas (todas protegidas)
	empresas := router.Group("/empresas")
	empresas.Use(middlewares.AuthMiddleware(cfg))
	{
		empresas.GET("", middlewares.RequirePermission("empresas.view"), empresaHandler.GetEmpresas)
		empresas.GET("/:id", middlewares.RequirePermission("empresas.view"), empresaHandler.GetEmpresa)
		empresas.POST("", middlewares.RequirePermission("empresas.create"), empresaHandler.CreateEmpresa)
		empresas.PUT("/:id", middlewares.RequirePermission("empresas.edit"), empresaHandler.UpdateEmpresa)
		empresas.DELETE("/:id", middlewares.RequirePermission("empresas.delete"), empresaHandler.DeleteEmpresa)
	}
}
