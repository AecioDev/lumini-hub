import { Layout, theme } from "antd";
import { ActiveCompanySwitcher } from "@/components/layout/ActiveCompanySwitcher";
import { Logo } from "@/components/layout/Logo";
import { SidebarMenu } from "@/components/layout/SidebarMenu";
import { TOP_BAR_HEIGHT } from "@/components/layout/layout-constants";
import type { ApiUserMenuItem } from "@/types/menu";

const { Sider } = Layout;

interface AppSiderProps {
  collapsed: boolean;
  menuItems: ApiUserMenuItem[];
}

// Menu lateral: logo da marca no topo (mesma altura do cabeçalho), o menu, que é a única parte que rola,
// e no pé a empresa ativa (clicar abre a troca de empresa).
export function AppSider({ collapsed, menuItems }: AppSiderProps) {
  const { token } = theme.useToken();

  return (
    <Sider width={240} collapsedWidth={76} collapsed={collapsed} trigger={null}>
      {/* O Sider do antd envolve os filhos num .ant-layout-sider-children com display:block:
          o flex column precisa deste wrapper próprio para o menu ocupar o espaço que sobra. */}
      <div style={{ display: "flex", flexDirection: "column", height: "100%" }}>
        <div
          style={{
            height: TOP_BAR_HEIGHT,
            flexShrink: 0,
            boxSizing: "border-box",
            display: "flex",
            alignItems: "center",
            justifyContent: collapsed ? "center" : "flex-start",
            padding: collapsed ? 0 : "0 16px",
            borderBottom: `1px solid ${token.colorBorderSecondary}`,
          }}
        >
          <Logo size={20} showWordmark={!collapsed} />
        </div>

        <div style={{ flex: 1, overflowY: "auto", padding: "8px 0" }}>
          <SidebarMenu items={menuItems} />
        </div>

        <ActiveCompanySwitcher collapsed={collapsed} />
      </div>
    </Sider>
  );
}
