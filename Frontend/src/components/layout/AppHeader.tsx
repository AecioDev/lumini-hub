import {
  BellOutlined,
  LogoutOutlined,
  MenuOutlined,
  UserOutlined,
} from "@ant-design/icons";
import { Avatar, Badge, Button, Dropdown, Layout, theme, type MenuProps } from "antd";
import { ThemeToggle } from "@/components/layout/ThemeToggle";
import { TOP_BAR_HEIGHT } from "@/components/layout/layout-constants";
import { useAuth } from "@/contexts/AuthContext";
import { getAvatarGradient, getInitials } from "@/utils/avatar";

const { Header } = Layout;

interface AppHeaderProps {
  onToggleMenu: () => void;
}

// Cabeçalho do app: botão do menu à esquerda; tema, avisos e usuário à direita. Não há busca global por
// enquanto. O nome da tela não fica aqui (está no caminho, AppBreadcrumb, e no título da página) nem a
// empresa ativa (fica no pé do menu lateral e na barra de status).
export function AppHeader({ onToggleMenu }: AppHeaderProps) {
  const { user, logout } = useAuth();
  const { token } = theme.useToken();

  const userMenuItems: MenuProps["items"] = [
    {
      key: "logout",
      label: "Sair",
      icon: <LogoutOutlined />,
      onClick: () => void logout(),
    },
  ];

  return (
    <Header
      style={{
        height: TOP_BAR_HEIGHT,
        padding: "0 24px",
        display: "flex",
        alignItems: "center",
        justifyContent: "space-between",
        borderBottom: `1px solid ${token.colorBorderSecondary}`,
      }}
    >
      <Button type="text" icon={<MenuOutlined />} onClick={onToggleMenu} />

      <div style={{ display: "flex", alignItems: "center", gap: 16 }}>
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
  );
}
