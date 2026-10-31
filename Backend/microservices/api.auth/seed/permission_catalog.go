package seed

// PermissionSeed é uma entrada do catálogo de permissões do sistema.
type PermissionSeed struct {
	Permission  string
	Description string
	Module      string
}

// PermissionCatalog é a fonte única, em código, do catálogo de permissões.
// Transcrito 1:1 do banco de desenvolvimento em 2026-10-01; em 2026-10-20 (CFG-10.1)
// entraram as 8 permissões que RequirePermission exige e o dev nunca cadastrou
// (customers.create/edit/delete, suppliers.view/create/edit/delete e
// integrations.edit). A limpeza de nomes/módulos herdados é tarefa separada (CFG-10.2).
// O seed só ADICIONA entradas ausentes; nunca altera nem apaga uma permissão já existente.
var PermissionCatalog = []PermissionSeed{
	// Admin
	{Permission: "dashboard.admin.view", Description: "Visualizar o dashboard do administrador", Module: "Admin"},
	{Permission: "audit.view", Description: "Visualizar a auditoria do sistema", Module: "Admin"},
	// Compras
	{Permission: "purchases.create", Description: "Criar compras", Module: "Compras"},
	{Permission: "purchases.delete", Description: "Excluir compras", Module: "Compras"},
	{Permission: "purchases.edit", Description: "Editar compras", Module: "Compras"},
	{Permission: "purchases.reports", Description: "Gerar relatórios de compras", Module: "Compras"},
	{Permission: "purchases.view", Description: "Visualizar compras", Module: "Compras"},
	{Permission: "suppliers.create", Description: "Criar fornecedores", Module: "Compras"},
	{Permission: "suppliers.delete", Description: "Excluir fornecedores", Module: "Compras"},
	{Permission: "suppliers.edit", Description: "Editar fornecedores", Module: "Compras"},
	{Permission: "suppliers.view", Description: "Visualizar fornecedores", Module: "Compras"},
	// Dashboard
	{Permission: "dashboard.view_default", Description: "Visualizar o dashboard padrão do perfil", Module: "Dashboard"},
	// Develop
	{Permission: "admin.create_permissions", Description: "Gerenciar perfis e permissões (catálogo)", Module: "Develop"},
	{Permission: "permissions.create", Description: "Criar permissões", Module: "Develop"},
	{Permission: "permissions.delete", Description: "Excluir permissões", Module: "Develop"},
	{Permission: "permissions.edit", Description: "Editar permissões", Module: "Develop"},
	{Permission: "permissions.view", Description: "Visualizar permissões", Module: "Develop"},
	{Permission: "roles.create", Description: "Criar perfis", Module: "Develop"},
	{Permission: "roles.delete", Description: "Excluir perfis", Module: "Develop"},
	{Permission: "roles.edit", Description: "Editar perfis", Module: "Develop"},
	{Permission: "roles.view", Description: "Visualizar perfis", Module: "Develop"},
	// Empresas
	{Permission: "companies.create", Description: "Criar empresas", Module: "Empresas"},
	{Permission: "companies.delete", Description: "Excluir empresas", Module: "Empresas"},
	{Permission: "companies.edit", Description: "Editar empresas", Module: "Empresas"},
	{Permission: "companies.fiscal_config.create", Description: "Criar configuração fiscal de empresas", Module: "Empresas"},
	{Permission: "companies.fiscal_config.edit", Description: "Editar configuração fiscal de empresas (inclui upload de certificado)", Module: "Empresas"},
	{Permission: "companies.fiscal_config.view", Description: "Visualizar configuração fiscal de empresas", Module: "Empresas"},
	{Permission: "companies.hierarchy.view", Description: "Visualizar empresas subsidiárias na hierarquia", Module: "Empresas"},
	{Permission: "companies.view", Description: "Visualizar empresas", Module: "Empresas"},
	// Estoque
	{Permission: "dashboard.inventory.view", Description: "Visualizar o dashboard de estoque", Module: "Estoque"},
	{Permission: "inventory.create", Description: "Adicionar itens ao estoque", Module: "Estoque"},
	{Permission: "inventory.delete", Description: "Remover itens do estoque", Module: "Estoque"},
	{Permission: "inventory.edit", Description: "Editar itens do estoque", Module: "Estoque"},
	{Permission: "inventory.reports", Description: "Gerar relatórios de estoque", Module: "Estoque"},
	{Permission: "inventory.view", Description: "Visualizar estoque", Module: "Estoque"},
	{Permission: "prices_promotions.view", Description: "Visualizar preços e promoções", Module: "Estoque"},
	{Permission: "product_location.view", Description: "Visualizar localização de produtos", Module: "Estoque"},
	{Permission: "products.view", Description: "Visualizar produtos", Module: "Estoque"},
	{Permission: "stock_locations.view", Description: "Visualizar locais de estoque", Module: "Estoque"},
	{Permission: "supplier_codes.view", Description: "Visualizar códigos por fornecedor", Module: "Estoque"},
	{Permission: "taxation.view", Description: "Visualizar tributação", Module: "Estoque"},
	// Financeiro
	{Permission: "dashboard.finance.view", Description: "Visualizar o dashboard financeiro", Module: "Financeiro"},
	{Permission: "finance.create", Description: "Criar transações financeiras", Module: "Financeiro"},
	{Permission: "finance.delete", Description: "Excluir transações financeiras", Module: "Financeiro"},
	{Permission: "finance.edit", Description: "Editar transações financeiras", Module: "Financeiro"},
	{Permission: "finance.receive_boleto", Description: "Receber boletos", Module: "Financeiro"},
	{Permission: "finance.reports", Description: "Gerar relatórios financeiros", Module: "Financeiro"},
	{Permission: "finance.view", Description: "Visualizar finanças", Module: "Financeiro"},
	{Permission: "finance.view_pendencies", Description: "Visualizar pendências financeiras", Module: "Financeiro"},
	// Gerencial
	{Permission: "dashboard.manager.view", Description: "Visualizar o dashboard gerencial", Module: "Gerencial"},
	// Integrações
	{Permission: "integrations.edit", Description: "Editar integrações", Module: "Integrações"},
	{Permission: "integrations.view", Description: "Visualizar integrações", Module: "Integrações"},
	// Usuários
	{Permission: "users.create", Description: "Criar usuários", Module: "Usuários"},
	{Permission: "users.delete", Description: "Excluir usuários", Module: "Usuários"},
	{Permission: "users.edit", Description: "Editar usuários", Module: "Usuários"},
	{Permission: "users.view", Description: "Visualizar usuários", Module: "Usuários"},
	// Vendas
	{Permission: "customers.create", Description: "Criar clientes", Module: "Vendas"},
	{Permission: "customers.delete", Description: "Excluir clientes", Module: "Vendas"},
	{Permission: "customers.edit", Description: "Editar clientes", Module: "Vendas"},
	{Permission: "customers.view", Description: "Visualizar clientes", Module: "Vendas"},
	{Permission: "dashboard.sales.view", Description: "Visualizar o dashboard de vendas", Module: "Vendas"},
	{Permission: "orders.view", Description: "Visualizar pedidos de vendas", Module: "Vendas"},
	{Permission: "payment_plans.view", Description: "Visualizar planos de pagamento", Module: "Vendas"},
	{Permission: "sales.create", Description: "Criar vendas", Module: "Vendas"},
	{Permission: "sales.delete", Description: "Excluir vendas", Module: "Vendas"},
	{Permission: "sales.edit", Description: "Editar vendas", Module: "Vendas"},
	{Permission: "sales.reports", Description: "Gerar relatórios de vendas", Module: "Vendas"},
	{Permission: "sales.view", Description: "Visualizar vendas", Module: "Vendas"},
}
