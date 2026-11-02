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
	mailConfigHandler := handlers.NewMailConfigHandler(uow, cfg.Security.CertificateEncryptionKey)

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
		companies.POST("/:id/logo", middlewares.RequirePermission("companies.edit"), companyHandler.UploadLogo)
		companies.DELETE("/:id/logo", middlewares.RequirePermission("companies.edit"), companyHandler.ClearLogo)
		// Sem RequirePermission de propósito (mesmo raciocínio já usado antes
		// pro logo da CompanyVisualConfig e pro CompanyOption/VisibleCompanies,
		// ver CFG-6.1.1/CFG-7): exibir o logo de uma empresa que o usuário já
		// enxerga é identidade de sessão, não administração do cadastro (que
		// segue exigindo companies.edit pra criar/trocar/remover). AuthMiddleware
		// do grupo já exige login.
		companies.GET("/:id/logo", companyHandler.GetLogo)
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
	// Rotas de E-mail do sistema (parâmetros globais, valem para qualquer empresa)
	systemSettings := router.Group("/system-settings")
	systemSettings.Use(middlewares.AuthMiddleware(cfg))
	{
		systemSettings.GET("/email", middlewares.RequirePermission("system_settings.view"), mailConfigHandler.GetSystemEmail)
		systemSettings.PUT("/email", middlewares.RequirePermission("system_settings.edit"), mailConfigHandler.SaveSystemEmail)
		systemSettings.POST("/email/test", middlewares.RequirePermission("system_settings.edit"), mailConfigHandler.SendSystemTestEmail)
	}

	// Rotas de E-mail próprio de cada empresa (sem ele, a empresa herda da Matriz e depois do sistema).
	// Mesmo desenho por empresa da configuração fiscal; aqui o DELETE existe porque apagar é "voltar a herdar".
	companyEmailConfigs := router.Group("/company-email-configs")
	companyEmailConfigs.Use(middlewares.AuthMiddleware(cfg))
	{
		companyEmailConfigs.GET("/by-company/:companyId", middlewares.RequirePermission("companies.email_config.view"), mailConfigHandler.GetCompanyEmail)
		companyEmailConfigs.PUT("/by-company/:companyId", middlewares.RequirePermission("companies.email_config.edit"), mailConfigHandler.SaveCompanyEmail)
		companyEmailConfigs.DELETE("/by-company/:companyId", middlewares.RequirePermission("companies.email_config.edit"), mailConfigHandler.DeleteCompanyEmail)
		companyEmailConfigs.POST("/by-company/:companyId/test", middlewares.RequirePermission("companies.email_config.edit"), mailConfigHandler.SendCompanyTestEmail)
	}
	// Rotas internas (chamadas só por outros microsserviços, com token de serviço). O gateway não repassa
	// /api/internal: elas não são alcançáveis de fora do servidor.
	internal := router.Group("/internal")
	internal.Use(middlewares.ServiceAuthMiddleware(cfg))
	{
		internal.POST("/mail/password-reset", mailConfigHandler.SendPasswordResetEmail)
	}
}
