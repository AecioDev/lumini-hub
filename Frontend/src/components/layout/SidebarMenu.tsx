import { StarFilled, StarOutlined } from "@ant-design/icons";
import { Icon } from "@iconify/react";
import { Menu, type MenuProps } from "antd";
import { useMemo } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { useThemeMode } from "@/contexts/ThemeContext";
import { useMenuFavorites } from "@/hooks/useMenuFavorites";
import type { ApiUserMenuItem } from "@/types/menu";

interface SidebarMenuProps {
  items: ApiUserMenuItem[];
}

type ItemsById = Map<string, ApiUserMenuItem>;

function menuKey(item: ApiUserMenuItem): string {
  // Sempre baseado no id: o mesmo href pode aparecer em mais de um lugar da
  // árvore (ex.: "Clientes" repetido em Vendas > Cadastros e em Financeiro >
  // Contas a Receber) — usar o href como key colidia e o antd Menu exige
  // keys únicas em toda a árvore, não só por nível.
  return `item-${item.id}`;
}

interface FavoriteStar {
  isFavorite: (id: number) => boolean;
  toggle: (item: ApiUserMenuItem) => void;
}

// Nome da tela com a estrela de favorito à direita. A estrela não pode navegar: o clique e o teclado
// ficam nela (stopPropagation), senão o antd trataria como clique no item do menu.
function ItemLabel({ node, favorite }: { node: ApiUserMenuItem; favorite: FavoriteStar }) {
  const marked = favorite.isFavorite(node.id);
  return (
    <span className="menu-item-label">
      <span className="menu-item-name">{node.name}</span>
      <button
        type="button"
        className="menu-favorite-star"
        data-marked={marked || undefined}
        aria-label={marked ? `Remover ${node.name} dos favoritos` : `Adicionar ${node.name} aos favoritos`}
        aria-pressed={marked}
        onClick={(e) => {
          e.stopPropagation();
          favorite.toggle(node);
        }}
        onKeyDown={(e) => e.stopPropagation()}
      >
        {marked ? <StarFilled /> : <StarOutlined />}
      </button>
    </span>
  );
}

function buildAntdMenuItems(
  nodes: ApiUserMenuItem[],
  itemsById: ItemsById,
  favorite: FavoriteStar
): MenuProps["items"] {
  return nodes.map((node) => {
    itemsById.set(menuKey(node), node);
    const icon = node.icon ? (
      <Icon icon={node.icon} width={16} height={16} />
    ) : undefined;

    if (node.children && node.children.length > 0) {
      return {
        key: menuKey(node),
        icon,
        label: node.name,
        children: buildAntdMenuItems(node.children, itemsById, favorite),
      };
    }

    return {
      key: menuKey(node),
      icon,
      // Só telas (com href) podem ser favoritas.
      label: node.href ? <ItemLabel node={node} favorite={favorite} /> : node.name,
    };
  });
}

export function SidebarMenu({ items }: SidebarMenuProps) {
  const navigate = useNavigate();
  const location = useLocation();
  const { mode } = useThemeMode();
  const { isFavorite, toggle } = useMenuFavorites();

  const { menuItems, itemsById } = useMemo(() => {
    const map: ItemsById = new Map();
    const built = buildAntdMenuItems(items, map, { isFavorite, toggle });
    return { menuItems: built, itemsById: map };
  }, [items, isFavorite, toggle]);

  const selectedKeys = useMemo(() => {
    const keys: string[] = [];
    for (const [key, node] of itemsById) {
      if (node.href && node.href === location.pathname) keys.push(key);
    }
    return keys;
  }, [itemsById, location.pathname]);

  const handleClick: MenuProps["onClick"] = ({ key }) => {
    const node = itemsById.get(key);
    if (node?.href) navigate(node.href);
  };

  return (
    <Menu
      mode="inline"
      theme={mode}
      style={{ borderInlineEnd: "none" }}
      items={menuItems}
      selectedKeys={selectedKeys}
      onClick={handleClick}
    />
  );
}
