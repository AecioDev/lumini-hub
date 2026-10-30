import { theme } from "antd";
import { ActiveCompanySwitcher } from "@/components/layout/ActiveCompanySwitcher";
import { Logo } from "@/components/layout/Logo";
import { SidebarMenu } from "@/components/layout/SidebarMenu";
import { TOP_BAR_HEIGHT } from "@/components/layout/layout-constants";
import type { ApiUserMenuItem } from "@/types/menu";

interface SiderContentProps {
  /** Só ícones/logo (menu lateral recolhido). Na gaveta é sempre `false`. */
  collapsed: boolean;
  menuItems: ApiUserMenuItem[];
}

// Conteúdo do menu: logo da marca no topo (mesma altura do cabeçalho), o menu, que é a única parte que rola,
// e no pé a empresa ativa (clicar abre a troca de empresa). É o mesmo no menu lateral fixo (AppSider) e na
// gaveta que o substitui em telas menores (MobileMenu).
export function SiderContent({ collapsed, menuItems }: SiderContentProps) {
  const { token } = theme.useToken();

  return (
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
  );
}
