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
	companyHandler := handlers.NewCompanyHandler(uow)
	companyFiscalConfigHandler := handlers.NewCompanyFiscalConfigHandler(uow, cfg.Security.CertificateEncryptionKey)
	companyVisualConfigHandler := handlers.NewCompanyVisualConfigHandler(uow)

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
	companies := router.Group("/companies")
	companies.Use(middlewares.AuthMiddleware(cfg))
	{
		companies.GET("", middlewares.RequirePermission("companies.view"), companyHandler.GetCompanies)
		companies.GET("/:id", middlewares.RequirePermission("companies.view"), companyHandler.GetCompany)
		companies.POST("", middlewares.RequirePermission("companies.create"), companyHandler.CreateCompany)
		companies.PUT("/:id", middlewares.RequirePermission("companies.edit"), companyHandler.UpdateCompany)
		companies.DELETE("/:id", middlewares.RequirePermission("companies.delete"), companyHandler.DeleteCompany)
	}

	// Rotas de Configuração Fiscal de Empresas (todas protegidas).
	// GET por empresa fica aqui, não aninhado em /companies/:id/..., porque o
	// Gin não permite dois nomes de wildcard diferentes (:id vs :companyId)
	// na mesma posição da árvore de rotas dentro do mesmo grupo.
	companyFiscalConfigs := router.Group("/company-fiscal-configs")
	companyFiscalConfigs.Use(middlewares.AuthMiddleware(cfg))
	{
		companyFiscalConfigs.GET("/by-company/:companyId", middlewares.RequirePermission("companies.fiscal_config.view"), companyFiscalConfigHandler.GetCompanyFiscalConfigByCompany)
		companyFiscalConfigs.POST("", middlewares.RequirePermission("companies.fiscal_config.create"), companyFiscalConfigHandler.CreateCompanyFiscalConfig)
		companyFiscalConfigs.PUT("/:id", middlewares.RequirePermission("companies.fiscal_config.edit"), companyFiscalConfigHandler.UpdateCompanyFiscalConfig)
		// Sem DELETE de propósito — config 1:1 da empresa, ver comentário em
		// CompanyFiscalConfigService.
		companyFiscalConfigs.POST("/:id/certificate", middlewares.RequirePermission("companies.fiscal_config.edit"), companyFiscalConfigHandler.UploadCertificate)
	}

	// Rotas de Configuração Visual de Empresas (todas protegidas). Mesmo
	// motivo de /by-company/:companyId ficar num grupo próprio, ver
	// comentário acima em companyFiscalConfigs.
	companyVisualConfigs := router.Group("/company-visual-configs")
	companyVisualConfigs.Use(middlewares.AuthMiddleware(cfg))
	{
		companyVisualConfigs.GET("/by-company/:companyId", middlewares.RequirePermission("companies.visual_config.view"), companyVisualConfigHandler.GetCompanyVisualConfigByCompany)
		companyVisualConfigs.POST("", middlewares.RequirePermission("companies.visual_config.create"), companyVisualConfigHandler.CreateCompanyVisualConfig)
		companyVisualConfigs.PUT("/:id", middlewares.RequirePermission("companies.visual_config.edit"), companyVisualConfigHandler.UpdateCompanyVisualConfig)
		// Sem DELETE de propósito — config 1:1 da empresa, ver comentário em
		// CompanyVisualConfigService.
		companyVisualConfigs.POST("/:id/logo", middlewares.RequirePermission("companies.visual_config.edit"), companyVisualConfigHandler.UploadLogo)
		companyVisualConfigs.GET("/:id/logo", middlewares.RequirePermission("companies.visual_config.view"), companyVisualConfigHandler.GetLogo)
	}
}
