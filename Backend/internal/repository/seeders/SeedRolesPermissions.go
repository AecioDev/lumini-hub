package seeders

import (
	"log"

	"simple-erp-service/internal/data-structure/models"

	"gorm.io/gorm"
)

// SeedRolesPermissions insere perfis e permissões no banco
func SeedRolesPermissions(db *gorm.DB) {
	err := db.Transaction(func(tx *gorm.DB) error {
		roles := []models.Role{
			{Name: "ADMIN", Description: "Administrador do sistema com acesso completo"},
			{Name: "GESTOR", Description: "Acesso gerencial a múltiplos módulos"},
			{Name: "VENDAS", Description: "Acesso ao módulo de vendas"},
			{Name: "ESTOQUE", Description: "Acesso ao módulo de estoque"},
			{Name: "FINANCEIRO", Description: "Acesso ao módulo financeiro"},
		}

		for _, role := range roles {
			if err := tx.Where("name = ?", role.Name).FirstOrCreate(&role).Error; err != nil {
				return err
			}
		}

		// --- Seed de Permissions (Permissões) ---
		permissions := []models.Permission{
			// Vendas
			{Permission: "sales.view", Description: "Visualizar vendas", Module: "sales"},
			{Permission: "sales.create", Description: "Criar vendas", Module: "sales"},
			{Permission: "sales.edit", Description: "Editar vendas", Module: "sales"},
			{Permission: "sales.delete", Description: "Excluir vendas", Module: "sales"},
			{Permission: "sales.reports", Description: "Gerar relatórios de vendas", Module: "sales"},
			// Novas permissões para módulos de vendas (ex: clientes, planos de pagamento)
			{Permission: "customers.view", Description: "Visualizar clientes", Module: "sales.cadastros"},
			{Permission: "payment_plans.view", Description: "Visualizar planos de pagamento", Module: "sales.cadastros"},
			{Permission: "orders.view", Description: "Visualizar pedidos de vendas", Module: "sales"},

			// Estoque
			{Permission: "inventory.view", Description: "Visualizar estoque", Module: "inventory"},
			{Permission: "inventory.create", Description: "Adicionar itens ao estoque", Module: "inventory"},
			{Permission: "inventory.edit", Description: "Editar itens do estoque", Module: "inventory"},
			{Permission: "inventory.delete", Description: "Remover itens do estoque", Module: "inventory"},
			{Permission: "inventory.reports", Description: "Gerar relatórios de estoque", Module: "inventory"},
			// Novas permissões para módulos de estoque (ex: produtos, fornecedores, locais)
			{Permission: "products.view", Description: "Visualizar produtos", Module: "inventory.cadastros"},
			{Permission: "supplier_codes.view", Description: "Visualizar códigos por fornecedor", Module: "inventory.cadastros"},
			{Permission: "stock_locations.view", Description: "Visualizar locais de estoque", Module: "inventory.cadastros"},
			{Permission: "product_location.view", Description: "Visualizar localização de produtos", Module: "inventory.cadastros"},
			{Permission: "prices_promotions.view", Description: "Visualizar preços e promoções", Module: "inventory.cadastros"},
			{Permission: "taxation.view", Description: "Visualizar tributação", Module: "inventory.cadastros"},

			// Financeiro
			{Permission: "finance.view", Description: "Visualizar finanças", Module: "finance"},
			{Permission: "finance.create", Description: "Criar transações financeiras", Module: "finance"},
			{Permission: "finance.edit", Description: "Editar transações financeiras", Module: "finance"},
			{Permission: "finance.delete", Description: "Excluir transações financeiras", Module: "finance"},
			{Permission: "finance.reports", Description: "Gerar relatórios financeiros", Module: "finance"},
			// Novas permissões Financeiro granular (receber boleto, ver pendências)
			{Permission: "finance.receive_boleto", Description: "Permissão para receber boleto financeiro", Module: "finance.contas_a_receber"},
			{Permission: "finance.view_pendencies", Description: "Visualizar pendências financeiras", Module: "finance.contas_a_receber"},

			// Permissões de Dashboard (novas)
			{Permission: "dashboard.view_default", Description: "Visualizar o dashboard padrão do perfil", Module: "dashboard"},
			{Permission: "dashboard.admin.view", Description: "Visualizar o dashboard do administrador", Module: "dashboard"},
			{Permission: "dashboard.sales.view", Description: "Visualizar o dashboard de vendas", Module: "dashboard"},
			{Permission: "dashboard.finance.view", Description: "Visualizar o dashboard financeiro", Module: "dashboard"},
			{Permission: "dashboard.manager.view", Description: "Visualizar o dashboard gerencial", Module: "dashboard"},
			// Nova permissão: Dashboard de Estoque
			{Permission: "dashboard.inventory.view", Description: "Visualizar o dashboard de estoque", Module: "dashboard"},

			// Configurações admin.create_permissions
			{Permission: "admin.create_permissions", Description: "Cadastrar Novas Permissões", Module: "admin"},

			// Usuários
			{Permission: "users.view", Description: "Visualizar usuários", Module: "users"},
			{Permission: "users.create", Description: "Criar usuários", Module: "users"},
			{Permission: "users.edit", Description: "Editar usuários", Module: "users"},
			{Permission: "users.delete", Description: "Excluir usuários", Module: "users"},

			// Permissões
			{Permission: "permissions.view", Description: "Visualizar Permissões", Module: "permissions"},
			{Permission: "permissions.create", Description: "Criar Permissões", Module: "permissions"},
			{Permission: "permissions.edit", Description: "Editar Permissões", Module: "permissions"},
			{Permission: "permissions.delete", Description: "Excluir Permissões", Module: "permissions"},

			// Perfis
			{Permission: "roles.view", Description: "Visualizar usuários", Module: "roles"},
			{Permission: "roles.create", Description: "Criar usuários", Module: "roles"},
			{Permission: "roles.edit", Description: "Editar usuários", Module: "roles"},
			{Permission: "roles.delete", Description: "Excluir usuários", Module: "roles"},
		}

		for _, perm := range permissions {
			if err := tx.Where("permission = ?", perm.Permission).FirstOrCreate(&perm).Error; err != nil {
				return err
			}
		}

		// --- Atribuir permissões aos Perfis (Roles) ---

		// Mapear roles por nome para facilitar a atribuição
		rolesMap := make(map[string]models.Role)
		for _, role := range roles {
			if err := tx.Where("name = ?", role.Name).First(&role).Error; err != nil {
				return err // Erro se a role não for encontrada (não deveria acontecer)
			}
			rolesMap[role.Name] = role
		}

		// ADMIN: Atribuir todas as permissões
		if adminRole, ok := rolesMap["ADMIN"]; ok {
			var allPermissions []models.Permission
			if err := tx.Find(&allPermissions).Error; err != nil {
				return err
			}
			// Use assignPermissionToRole em loop para garantir que não haja duplicatas de associações
			for _, perm := range allPermissions {
				if err := assignPermissionToRole(tx, adminRole.ID, perm.Permission); err != nil {
					return err
				}
			}
		}

		// GESTOR: Permissões de visualização, relatórios e dashboard gerencial/default
		if managerRole, ok := rolesMap["GESTOR"]; ok {
			var viewAndReportPermissions []models.Permission
			if err := tx.Where("permission LIKE ?", "%.view").Or("permission LIKE ?", "%.reports").Find(&viewAndReportPermissions).Error; err != nil {
				return err
			}
			for _, perm := range viewAndReportPermissions {
				if err := assignPermissionToRole(tx, managerRole.ID, perm.Permission); err != nil {
					return err
				}
			}
			// Adicionar permissões de dashboard específicas do gerente
			assignPermissionToRole(tx, managerRole.ID, "dashboard.manager.view")
			assignPermissionToRole(tx, managerRole.ID, "dashboard.view_default")
		}

		// VENDAS: Todas as permissões do módulo 'sales' + dashboards de vendas/default + financeiro granular
		if salesRole, ok := rolesMap["VENDAS"]; ok {
			// Atribui todas as permissões do módulo 'sales' (ex: sales.view, sales.create, etc.)
			assignRolePermissionsByModule(tx, salesRole.ID, "sales")
			// Atribui permissões de cadastro de vendas (ex: customers.view, payment_plans.view)
			assignRolePermissionsByModule(tx, salesRole.ID, "sales.cadastros")
			// Atribui permissões de dashboard
			assignPermissionToRole(tx, salesRole.ID, "dashboard.sales.view")
			assignPermissionToRole(tx, salesRole.ID, "dashboard.view_default")
			// Permissões extras para VENDAS
			assignPermissionToRole(tx, salesRole.ID, "finance.receive_boleto")
			assignPermissionToRole(tx, salesRole.ID, "finance.view_pendencies")
		}

		// ESTOQUE: Todas as permissões do módulo 'inventory' + dashboards de estoque/default
		if stockRole, ok := rolesMap["ESTOQUE"]; ok {
			// Atribui todas as permissões do módulo 'inventory'
			assignRolePermissionsByModule(tx, stockRole.ID, "inventory")
			// Atribui permissões de cadastro de estoque
			assignRolePermissionsByModule(tx, stockRole.ID, "inventory.cadastros")
			// Atribui permissões de dashboard
			assignPermissionToRole(tx, stockRole.ID, "dashboard.inventory.view")
			assignPermissionToRole(tx, stockRole.ID, "dashboard.view_default")
		}

		// FINANCEIRO: Todas as permissões do módulo 'finance' + dashboards de financeiro/default
		if financeRole, ok := rolesMap["FINANCEIRO"]; ok {
			// Atribui todas as permissões do módulo 'finance'
			assignRolePermissionsByModule(tx, financeRole.ID, "finance")
			// Atribui permissões de contas a receber (ex: finance.receive_boleto, finance.view_pendencies)
			assignRolePermissionsByModule(tx, financeRole.ID, "finance.contas_a_receber")
			// Atribui permissões de dashboard
			assignPermissionToRole(tx, financeRole.ID, "dashboard.finance.view")
			assignPermissionToRole(tx, financeRole.ID, "dashboard.view_default")
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Erro ao executar SeedRolesPermissions: %v", err)
	}

	log.Println("Seeder de Roles e Permissions executado com sucesso!")
}

// assignPermissionToRole atribui uma única permissão a um perfil, se a associação não existir.
func assignPermissionToRole(tx *gorm.DB, roleID uint, permissionName string) error {
	var permission models.Permission
	if err := tx.Where("permission = ?", permissionName).First(&permission).Error; err != nil {
		// Se a permissão não existe, isso pode ser um aviso dependendo da sua estratégia.
		// Aqui, vamos apenas logar e retornar nil para não interromper o seeder.
		log.Printf("Aviso: Permissão '%s' não encontrada para atribuição à role ID %d. Erro: %v", permissionName, roleID, err)
		return nil
	}

	var role models.Role
	if err := tx.First(&role, roleID).Error; err != nil {
		return err // Erro fatal: a role não foi encontrada, o que não deveria acontecer aqui.
	}

	// Consulta direta para verificar se a associação já existe
	// Usa a tabela de junção (role_permissions) implicitamente pelo GORM
	var existingAssociationCount int64
	// MUDANÇA AQUI: Corrigindo a forma de contar associações para GORM
	err := tx.Table("role_permissions").
		Where("role_id = ? AND permission_id = ?", role.ID, permission.ID).
		Count(&existingAssociationCount).Error
	if err != nil {
		return err
	}

	if existingAssociationCount == 0 {
		// MUDANÇA AQUI: Usando o modelo carregado para a associação
		log.Printf("Atribuindo permissão '%s' (ID %d) à role '%s' (ID %d).", permission.Permission, permission.ID, role.Name, role.ID)
		return tx.Model(&role).Association("Permissions").Append(&permission)
	}
	log.Printf("Permissão '%s' já atribuída à role '%s'. Ignorando.", permission.Permission, role.Name)
	return nil
}

// assignRolePermissionsByModule atribui todas as permissões de um módulo
// a um perfil específico, se a associação não existir.
func assignRolePermissionsByModule(tx *gorm.DB, roleID uint, moduleName string) error {
	var role models.Role
	if err := tx.First(&role, roleID).Error; err != nil {
		return err // Erro fatal: a role não foi encontrada.
	}

	var modulePermissions []models.Permission
	// Busca todas as permissões onde o nome da permissão começa com o nome do módulo (ex: "sales.%")
	// OU onde o campo 'Module' da permissão é exatamente o nome do módulo (ex: "sales")
	if err := tx.Where("permission LIKE ? OR module = ?", moduleName+".%", moduleName).Find(&modulePermissions).Error; err != nil {
		return err
	}

	for _, perm := range modulePermissions {
		// Reutiliza a função auxiliar para cada permissão individual, garantindo que não haja duplicatas
		if err := assignPermissionToRole(tx, role.ID, perm.Permission); err != nil {
			return err
		}
	}
	return nil
}

// Removendo a função antiga para evitar confusão.
// Se você ainda tiver chamadas para essa função no seu main.go,
// substitua-as pelas novas funções mais específicas.
/*
func assignRolePermissions(tx *gorm.DB, roleName string, moduleName string) error {
	var role models.Role
	if err := tx.Where("name = ?", roleName).First(&role).Error; err != nil {
		return err
	}

	var modulePermissions []models.Permission
	if err := tx.Where("module = ?", moduleName).Find(&modulePermissions).Error; err != nil {
		return err
	}

	return tx.Model(&role).Association("Permissions").Append(&modulePermissions)
}
*/
