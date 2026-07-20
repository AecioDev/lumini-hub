package domain

import (
	"sort"

	"lumini-hub/common/utils"

	"gorm.io/gorm"
)

// MenuItem representa um item do menu lateral, com hierarquia livre (pai/filho/neto)
type MenuItem struct {
	gorm.Model

	Name         string      `gorm:"size:100;not null" json:"name"`
	Icon         string      `gorm:"size:100" json:"icon"`
	Href         string      `gorm:"size:255" json:"href"`
	ParentID     *uint       `json:"parent_id"`
	Children     []MenuItem  `gorm:"foreignKey:ParentID" json:"-"`
	PermissionID *uint       `json:"permission_id"`
	Permission   *Permission `gorm:"foreignKey:PermissionID" json:"-"`
	// Position controla a ordem de exibição entre irmãos.
	// Nome escolhido para não colidir com "order", palavra reservada no Postgres.
	Position int  `gorm:"default:0" json:"position"`
	IsActive bool `gorm:"default:true" json:"is_active"`
}

// TableName especifica o nome da tabela
func (MenuItem) TableName() string {
	return "menu_items"
}

// CreateMenuItemRequest representa os dados para criar um item de menu
type CreateMenuItemRequest struct {
	Name         string `json:"name" binding:"required,min=1,max=100"`
	Icon         string `json:"icon"`
	Href         string `json:"href"`
	ParentID     *uint  `json:"parent_id"`
	PermissionID *uint  `json:"permission_id"`
	Position     int    `json:"position"`
}

// UpdateMenuItemRequest representa os dados para atualizar um item de menu
type UpdateMenuItemRequest struct {
	Name         string `json:"name" binding:"required,min=1,max=100"`
	Icon         string `json:"icon"`
	Href         string `json:"href"`
	ParentID     *uint  `json:"parent_id"`
	PermissionID *uint  `json:"permission_id"`
	Position     int    `json:"position"`
	IsActive     *bool  `json:"is_active"`
}

// ApiMenuItem representa um item de menu para telas administrativas (árvore completa, sem filtro)
type ApiMenuItem struct {
	ID           uint          `json:"id"`
	Name         string        `json:"name"`
	Icon         string        `json:"icon"`
	Href         string        `json:"href"`
	ParentID     *uint         `json:"parent_id"`
	PermissionID *uint         `json:"permission_id"`
	Permission   *ApiPermission `json:"permission,omitempty"`
	Position     int           `json:"position"`
	IsActive     bool          `json:"is_active"`
	Children     []ApiMenuItem `json:"children,omitempty"`
}

// ApiUserMenuItem representa um item de menu já filtrado para um usuário específico
// (árvore enxuta, embutida na resposta de login/refresh/me)
type ApiUserMenuItem struct {
	ID       uint              `json:"id"`
	Name     string            `json:"name"`
	Icon     string            `json:"icon"`
	Href     string            `json:"href"`
	Children []ApiUserMenuItem `json:"children,omitempty"`
}

// ApiMenuItemFromModel converte um MenuItem (flat, sem filhos) para ApiMenuItem
func ApiMenuItemFromModel(m MenuItem) ApiMenuItem {
	dto := ApiMenuItem{
		ID:           m.ID,
		Name:         m.Name,
		Icon:         m.Icon,
		Href:         m.Href,
		ParentID:     m.ParentID,
		PermissionID: m.PermissionID,
		Position:     m.Position,
		IsActive:     m.IsActive,
	}
	if m.Permission != nil {
		p := ApiPermissionFromModel(*m.Permission)
		dto.Permission = &p
	}
	return dto
}

// BuildMenuTree agrupa uma lista plana de itens (com Permission pré-carregada) por ParentID,
// ordenando os irmãos por Position, e devolve só as raízes com filhos aninhados.
func BuildMenuTree(flat []MenuItem) []ApiMenuItem {
	sort.SliceStable(flat, func(i, j int) bool { return flat[i].Position < flat[j].Position })

	childrenByParent := make(map[uint][]MenuItem)
	var roots []MenuItem
	for _, item := range flat {
		if item.ParentID != nil {
			childrenByParent[*item.ParentID] = append(childrenByParent[*item.ParentID], item)
		} else {
			roots = append(roots, item)
		}
	}

	var build func(nodes []MenuItem) []ApiMenuItem
	build = func(nodes []MenuItem) []ApiMenuItem {
		result := make([]ApiMenuItem, 0, len(nodes))
		for _, n := range nodes {
			dto := ApiMenuItemFromModel(n)
			dto.Children = build(childrenByParent[n.ID])
			result = append(result, dto)
		}
		return result
	}

	return build(roots)
}

// FilterMenuTreeForUser filtra a árvore de menu para as permissões de um usuário.
// Nós que só agrupam filhos (sem Href, ex. "Cadastros") nunca ganham acesso por conta
// própria — eles só aparecem se sobrar algum descendente visível. Nós que são destino de
// navegação (com Href) ficam visíveis se: o papel do usuário dá bypass pra permissão
// exigida por aquele nó (ver utils.IsDeveloperOnlyPermission — o catálogo de Perfis e
// Permissões só faz bypass pra "DEVELOP", não pra "ADMIN"), o nó não exige
// permissão própria, o usuário tem a permissão exigida, ou algum descendente ficou visível.
func FilterMenuTreeForUser(tree []ApiMenuItem, userPermissionCodes map[string]bool, role string) []ApiUserMenuItem {
	var visible []ApiUserMenuItem
	for _, node := range tree {
		children := FilterMenuTreeForUser(node.Children, userPermissionCodes, role)

		isDestination := node.Href != ""
		requiredPermission := ""
		if node.Permission != nil {
			requiredPermission = node.Permission.Permission
		}

		bypasses := role == utils.RoleDeveloper ||
			(role == utils.RoleAdmin &&
				(requiredPermission == "" || !utils.IsDeveloperOnlyPermission(requiredPermission)))

		hasOwnAccess := node.IsActive && isDestination &&
			(bypasses || requiredPermission == "" || userPermissionCodes[requiredPermission])

		if hasOwnAccess || len(children) > 0 {
			visible = append(visible, ApiUserMenuItem{
				ID:       node.ID,
				Name:     node.Name,
				Icon:     node.Icon,
				Href:     node.Href,
				Children: children,
			})
		}
	}
	return visible
}
