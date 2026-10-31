import {
  BellOutlined,
  LogoutOutlined,
  MenuOutlined,
  UserOutlined,
} from "@ant-design/icons";
import { Avatar, Badge, Button, Dropdown, Layout, theme, type MenuProps } from "antd";
import { HeaderFavorites } from "@/components/layout/HeaderFavorites";
import { ThemeToggle } from "@/components/layout/ThemeToggle";
import { TOP_BAR_HEIGHT } from "@/components/layout/layout-constants";
import { useAuth } from "@/contexts/AuthContext";
import { getAvatarGradient, getInitials } from "@/utils/avatar";

const { Header } = Layout;

interface AppHeaderProps {
  onToggleMenu: () => void;
}

// Cabeçalho do app: botão do menu à esquerda; favoritos do menu no centro; tema, avisos e usuário à direita. Não há busca global por
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
      <Button
        type="text"
        icon={<MenuOutlined />}
        onClick={onToggleMenu}
        aria-label="Abrir ou recolher o menu"
      />

      {/* Favoritos centralizados no espaço livre entre o botão do menu e o grupo da direita. */}
      <div style={{ flex: 1, minWidth: 0, display: "flex", justifyContent: "center", padding: "0 16px" }}>
        <HeaderFavorites />
      </div>

      <div style={{ display: "flex", alignItems: "center", gap: 16 }}>
        <ThemeToggle />
        <Badge dot offset={[-2, 2]}>
          <Button shape="circle" icon={<BellOutlined />} aria-label="Notificações" />
        </Badge>
        <Dropdown menu={{ items: userMenuItems }} trigger={["click"]}>
          {/* Um <button> de verdade (e não só o Avatar): recebe foco pelo teclado e abre com Enter. */}
          <button
            type="button"
            aria-label="Menu do usuário"
            style={{ padding: 0, border: "none", background: "none", borderRadius: "50%", cursor: "pointer", lineHeight: 0 }}
          >
            <Avatar
              style={{
                background: getAvatarGradient(user?.name ?? ""),
                cursor: "pointer",
              }}
              icon={!user?.name ? <UserOutlined /> : undefined}
            >
              {user?.name ? getInitials(user.name) : undefined}
            </Avatar>
          </button>
        </Dropdown>
      </div>
    </Header>
  );
}
