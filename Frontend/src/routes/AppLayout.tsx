import { useState } from "react";
import { Outlet, useLocation } from "react-router-dom";
import {
  BellOutlined,
  MenuOutlined,
  LogoutOutlined,
  UserOutlined,
} from "@ant-design/icons";
import {
  Avatar,
  Badge,
  Button,
  Dropdown,
  Input,
  Layout,
  theme,
  type MenuProps,
} from "antd";
import { useAuth } from "@/contexts/AuthContext";
import { SidebarMenu } from "@/components/layout/SidebarMenu";
import { ThemeToggle } from "@/components/layout/ThemeToggle";
import { Logo } from "@/components/layout/Logo";
import { getAvatarGradient, getInitials } from "@/utils/avatar";
import type { ApiUserMenuItem } from "@/types/menu";

const { Sider, Header, Content } = Layout;

function findActiveTitle(items: ApiUserMenuItem[], pathname: string): string | null {
  for (const item of items) {
    if (item.href === pathname) return item.name;
    if (item.children) {
      const found = findActiveTitle(item.children, pathname);
      if (found) return found;
    }
  }
  return null;
}

const EXTRA_TITLES: Array<[RegExp, string]> = [
  [/^\/settings\/users\/create$/, "Novo Usuário"],
  [/^\/settings\/users\/\d+\/edit$/, "Editar Usuário"],
  [/^\/settings\/roles\/create$/, "Novo Perfil"],
  [/^\/settings\/roles\/\d+\/edit$/, "Editar Perfil"],
  [/^\/settings\/permissions\/create$/, "Nova Permissão"],
  [/^\/settings\/permissions\/\d+\/edit$/, "Editar Permissão"],
];

export function AppLayout() {
  const [collapsed, setCollapsed] = useState(false);
  const { user, menuItems, logout } = useAuth();
  const location = useLocation();
  const { token } = theme.useToken();

  const pageTitle =
    findActiveTitle(menuItems, location.pathname) ??
    EXTRA_TITLES.find(([pattern]) => pattern.test(location.pathname))?.[1] ??
    "Lumini Hub";

  const userMenuItems: MenuProps["items"] = [
    {
      key: "logout",
      label: "Sair",
      icon: <LogoutOutlined />,
      onClick: () => void logout(),
    },
  ];

  return (
    <Layout style={{ minHeight: "100vh" }}>
      <Sider
        width={240}
        collapsedWidth={76}
        collapsed={collapsed}
        trigger={null}
        style={{ display: "flex", flexDirection: "column" }}
      >
        <div
          style={{
            display: "flex",
            alignItems: "center",
            justifyContent: collapsed ? "center" : "flex-start",
            padding: collapsed ? "20px 0" : "20px 16px",
            borderBottom: "1px solid rgba(255,255,255,0.08)",
          }}
        >
          <Logo size={22} showWordmark={!collapsed} />
        </div>

        <div style={{ flex: 1, overflowY: "auto", padding: "8px 0" }}>
          <SidebarMenu items={menuItems} />
        </div>

        <div
          style={{
            display: "flex",
            alignItems: "center",
            justifyContent: collapsed ? "center" : "flex-start",
            gap: 10,
            padding: collapsed ? "12px 0" : "12px 16px",
            borderTop: "1px solid rgba(255,255,255,0.08)",
          }}
        >
          <Avatar style={{ background: getAvatarGradient(user?.name ?? "") }}>
            {getInitials(user?.name ?? "?")}
          </Avatar>
          {!collapsed && (
            <div style={{ minWidth: 0 }}>
              <div
                style={{
                  fontSize: 13,
                  fontWeight: 500,
                  color: "rgba(255,255,255,0.85)",
                  whiteSpace: "nowrap",
                  overflow: "hidden",
                  textOverflow: "ellipsis",
                }}
              >
                {user?.name}
              </div>
              <div style={{ fontSize: 12, color: "rgba(255,255,255,0.45)" }}>
                {user?.role?.name}
              </div>
            </div>
          )}
        </div>
      </Sider>

      <Layout>
        <Header
          style={{
            height: 64,
            padding: "0 24px",
            display: "flex",
            alignItems: "center",
            justifyContent: "space-between",
            borderBottom: `1px solid ${token.colorBorderSecondary}`,
          }}
        >
          <div style={{ display: "flex", alignItems: "center", gap: 16 }}>
            <Button
              type="text"
              icon={<MenuOutlined />}
              onClick={() => setCollapsed((prev) => !prev)}
            />
            <span style={{ fontSize: 15, fontWeight: 600 }}>{pageTitle}</span>
          </div>

          <div style={{ display: "flex", alignItems: "center", gap: 16 }}>
            <Input placeholder="Buscar..." style={{ width: 220 }} allowClear />
            <ThemeToggle />
            <Badge dot offset={[-2, 2]}>
              <Button shape="circle" icon={<BellOutlined />} />
            </Badge>
            <Dropdown menu={{ items: userMenuItems }} trigger={["click"]}>
              <Avatar
                style={{
                  background: getAvatarGradient(user?.name ?? ""),
                  cursor: "pointer",
                }}
                icon={!user?.name ? <UserOutlined /> : undefined}
              >
                {user?.name ? getInitials(user.name) : undefined}
              </Avatar>
            </Dropdown>
          </div>
        </Header>

        <Content style={{ padding: 24, overflowY: "auto" }}>
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  );
}
