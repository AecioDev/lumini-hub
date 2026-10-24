package seeder

import (
	"testing"

	"lumini-hub/api.auth/seed"
)

func walkMenuSeeds(nodes []menuItemSeed, visit func(menuItemSeed)) {
	for _, node := range nodes {
		visit(node)
		walkMenuSeeds(node.Children, visit)
	}
}

// Um item de menu ligado a uma permissão que não existe no catálogo nasceria sem
// permissão própria, e isso nunca é corrigido depois (o seed pula itens existentes).
func TestMenuItemSeedTree_TodaPermissaoExisteNoCatalogo(t *testing.T) {
	catalog := make(map[string]bool, len(seed.PermissionCatalog))
	for _, entry := range seed.PermissionCatalog {
		catalog[entry.Permission] = true
	}

	walkMenuSeeds(menuItemSeedTree, func(node menuItemSeed) {
		if node.PermissionCode == "" {
			return
		}
		if !catalog[node.PermissionCode] {
			t.Errorf("item de menu %q usa a permissão %q, que não está no catálogo", node.Name, node.PermissionCode)
		}
	})
}

// Todo item com tela (href) precisa de permissão própria; só os grupos, sem href, ficam sem.
func TestMenuItemSeedTree_ItemComTelaTemPermissao(t *testing.T) {
	walkMenuSeeds(menuItemSeedTree, func(node menuItemSeed) {
		if node.Href != "" && node.PermissionCode == "" {
			t.Errorf("item de menu %q (%s) tem tela e não tem permissão", node.Name, node.Href)
		}
	})
}

// Fornecedores é o cadastro (suppliers.*); purchases.* é a rotina de Compras.
func TestMenuItemSeedTree_FornecedoresUsaPermissaoDeFornecedores(t *testing.T) {
	found := 0
	walkMenuSeeds(menuItemSeedTree, func(node menuItemSeed) {
		if node.Href != "/suppliers" {
			return
		}
		found++
		if node.PermissionCode != "suppliers.view" {
			t.Errorf("item %q (/suppliers) usa %q, esperado suppliers.view", node.Name, node.PermissionCode)
		}
	})
	if found == 0 {
		t.Error("nenhum item de menu aponta para /suppliers")
	}
}
