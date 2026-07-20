import { Icon } from "@iconify/react";
import { Menu, type MenuProps } from "antd";
import { useMemo } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { useThemeMode } from "@/contexts/ThemeContext";
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

function buildAntdMenuItems(
  nodes: ApiUserMenuItem[],
  itemsById: ItemsById
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
        children: buildAntdMenuItems(node.children, itemsById),
      };
    }

    return {
      key: menuKey(node),
      icon,
      label: node.name,
    };
  });
}

export function SidebarMenu({ items }: SidebarMenuProps) {
  const navigate = useNavigate();
  const location = useLocation();
  const { mode } = useThemeMode();

  const { menuItems, itemsById } = useMemo(() => {
    const map: ItemsById = new Map();
    const built = buildAntdMenuItems(items, map);
    return { menuItems: built, itemsById: map };
  }, [items]);

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
