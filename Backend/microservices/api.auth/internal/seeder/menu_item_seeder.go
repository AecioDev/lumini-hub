package seeder

import (
	"log"

	"lumini-hub/api.auth/internal/domain"

	"gorm.io/gorm"
)

// menuItemSeed descreve um nó da árvore de menu a ser semeada, com os filhos aninhados
// e o código de permissão (string), resolvido para PermissionID em tempo de seed.
type menuItemSeed struct {
	Name           string
	Icon           string
	Href           string
	PermissionCode string
	Children       []menuItemSeed
}

// menuItemSeedTree replica a árvore hoje descrita em Frontend/src/config/navigation.ts,
// com ícones Phosphor (via Iconify, prefixo "ph:") para manter a aparência visual atual.
var menuItemSeedTree = []menuItemSeed{
	{
		Name: "Dashboard", Icon: "ph:squares-four", Href: "/dashboard",
		PermissionCode: "dashboard.view_default",
	},
	{
		Name: "Vendas", Icon: "ph:chart-bar",
		Children: []menuItemSeed{
			{Name: "Pedidos", Icon: "ph:clipboard-text", Href: "/sales/orders", PermissionCode: "orders.view"},
			{
				Name: "Cadastros", Icon: "ph:gear",
				Children: []menuItemSeed{
					{Name: "Clientes", Icon: "ph:users", Href: "/customers", PermissionCode: "customers.view"},
				},
			},
		},
	},
	{
		Name: "Compras", Icon: "ph:shopping-cart",
		Children: []menuItemSeed{
			{Name: "Pedidos de Compras", Icon: "ph:clipboard-text", Href: "/purchases/orders", PermissionCode: "purchases.view"},
			{
				Name: "Cadastros", Icon: "ph:gear",
				Children: []menuItemSeed{
					// TODO: não existe permissão dedicada de Fornecedores no catálogo ainda;
					// usando purchases.view como aproximação até criarem suppliers.view.
					{Name: "Fornecedores", Icon: "ph:briefcase", Href: "/suppliers", PermissionCode: "purchases.view"},
				},
			},
		},
	},
	{
		Name: "Estoque", Icon: "ph:package",
		Children: []menuItemSeed{
			{Name: "Produtos", Icon: "ph:archive", Href: "/products", PermissionCode: "products.view"},
			{
				Name: "Movimentação", Icon: "ph:arrows-left-right",
				Children: []menuItemSeed{
					{Name: "Notas de Entrada", Icon: "ph:arrow-down", Href: "/inventory/movements/inbound", PermissionCode: "inventory.view"},
					{Name: "Notas de Saídas", Icon: "ph:arrow-up", Href: "/inventory/movements/outbound", PermissionCode: "inventory.view"},
					{Name: "Ajuste de Estoque", Icon: "ph:sliders", Href: "/inventory/movements/adjustments", PermissionCode: "inventory.edit"},
				},
			},
			{Name: "Níveis de Estoque", Icon: "ph:stack", Href: "/inventory/stock-levels", PermissionCode: "inventory.view"},
			{
				Name: "Cadastro", Icon: "ph:gear",
				Children: []menuItemSeed{
					{Name: "Locations", Icon: "ph:map-pin", Href: "/inventory/setup/locations", PermissionCode: "inventory.view"},
				},
			},
		},
	},
	{
		Name: "Financeiro", Icon: "ph:currency-circle-dollar",
		Children: []menuItemSeed{
			{Name: "Dashboard Financeiro", Icon: "ph:squares-four", Href: "/dashboard/financial", PermissionCode: "dashboard.finance.view"},
			{
				Name: "Contas a Receber", Icon: "ph:receipt", Href: "/financial/accounts-receivable", PermissionCode: "finance.view_pendencies",
				Children: []menuItemSeed{
					{Name: "Clientes", Icon: "ph:users", Href: "/customers", PermissionCode: "customers.view"},
				},
			},
			{
				Name: "Contas a Pagar", Icon: "ph:bank", Href: "/financial/accounts-payable", PermissionCode: "finance.view",
				Children: []menuItemSeed{
					// TODO: mesma ausência de permissão dedicada de Fornecedores, ver acima.
					{Name: "Fornecedores", Icon: "ph:briefcase", Href: "/suppliers", PermissionCode: "purchases.view"},
				},
			},
		},
	},
	{
		Name: "Configurações", Icon: "ph:user-gear",
		Children: []menuItemSeed{
			{Name: "Usuários", Icon: "ph:user-circle", Href: "/settings/users", PermissionCode: "users.view"},
			{Name: "Perfis e Permissões", Icon: "ph:shield-check", Href: "/settings/roles", PermissionCode: "admin.create_permissions"},
			{Name: "Integrações", Icon: "ph:arrows-clockwise", Href: "/settings/integrations", PermissionCode: "integrations.view"},
		},
	},
}

// SeedMenuItems popula a árvore inicial do menu, apenas se a tabela ainda estiver vazia.
func SeedMenuItems(db *gorm.DB) error {
	var count int64
	if err := db.Model(&domain.MenuItem{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	log.Println("[api.auth] Semeando árvore inicial de itens de menu...")
	return insertMenuItemSeeds(db, menuItemSeedTree, nil)
}

func insertMenuItemSeeds(db *gorm.DB, nodes []menuItemSeed, parentID *uint) error {
	for position, node := range nodes {
		item := domain.MenuItem{
			Name:         node.Name,
			Icon:         node.Icon,
			Href:         node.Href,
			ParentID:     parentID,
			PermissionID: resolvePermissionID(db, node.PermissionCode),
			Position:     position,
			IsActive:     true,
		}
		if err := db.Create(&item).Error; err != nil {
			return err
		}
		if len(node.Children) > 0 {
			if err := insertMenuItemSeeds(db, node.Children, &item.ID); err != nil {
				return err
			}
		}
	}
	return nil
}

// resolvePermissionID busca o ID de uma permissão pelo código. Se não existir, devolve nil
// (o item de menu fica sem exigência de permissão própria) em vez de falhar o seed inteiro.
func resolvePermissionID(db *gorm.DB, code string) *uint {
	if code == "" {
		return nil
	}
	var perm domain.Permission
	if err := db.Where("permission = ?", code).First(&perm).Error; err != nil {
		log.Printf("[api.auth] Aviso: permissão %q não encontrada ao semear menu, item ficará sem restrição própria", code)
		return nil
	}
	return &perm.ID
}
