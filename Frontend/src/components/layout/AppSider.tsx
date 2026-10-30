import { Layout } from "antd";
import { SiderContent } from "@/components/layout/SiderContent";
import type { ApiUserMenuItem } from "@/types/menu";

const { Sider } = Layout;

interface AppSiderProps {
  collapsed: boolean;
  menuItems: ApiUserMenuItem[];
}

// Menu lateral fixo (telas médias e grandes). Em tela estreita ele some e o menu abre numa gaveta (MobileMenu).
export function AppSider({ collapsed, menuItems }: AppSiderProps) {
  return (
    <Sider width={240} collapsedWidth={76} collapsed={collapsed} trigger={null}>
      {/* O Sider do antd envolve os filhos num .ant-layout-sider-children com display:block:
          o SiderContent tem altura 100% para o menu ocupar o espaço que sobra. */}
      <SiderContent collapsed={collapsed} menuItems={menuItems} />
    </Sider>
  );
}
