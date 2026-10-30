import { useState } from "react";
import { Outlet } from "react-router-dom";
import { Layout } from "antd";
import { useAuth } from "@/contexts/AuthContext";
import { AppBreadcrumb } from "@/components/layout/AppBreadcrumb";
import { AppHeader } from "@/components/layout/AppHeader";
import { AppSider } from "@/components/layout/AppSider";
import { StatusBar } from "@/components/layout/StatusBar";

const { Content } = Layout;

const SIDEBAR_STORAGE_KEY = "lumini-hub:sidebar-collapsed";

// O menu recolhido é uma preferência de cada usuário: lembrada entre recarregamentos, como o tema.
function readCollapsed(): boolean {
  try {
    return window.localStorage.getItem(SIDEBAR_STORAGE_KEY) === "1";
  } catch {
    return false;
  }
}

// Casca do app: menu lateral, cabeçalho, caminho da tela, corpo (a única área que rola) e a barra de status.
export function AppLayout() {
  const [collapsed, setCollapsed] = useState(readCollapsed);
  const { menuItems } = useAuth();

  const toggleMenu = () =>
    setCollapsed((prev) => {
      const next = !prev;
      try {
        window.localStorage.setItem(SIDEBAR_STORAGE_KEY, next ? "1" : "0");
      } catch {
        // sem localStorage (janela privada etc.): o menu só não é lembrado
      }
      return next;
    });

  return (
    <Layout style={{ height: "100vh", overflow: "hidden" }}>
      <Layout hasSider style={{ flex: 1, minHeight: 0 }}>
        <AppSider collapsed={collapsed} menuItems={menuItems} />

        <Layout>
          <AppHeader onToggleMenu={toggleMenu} />
          <AppBreadcrumb menuItems={menuItems} />

          <Content style={{ padding: 24, overflowY: "auto", minHeight: 0 }}>
            <Outlet />
          </Content>
        </Layout>
      </Layout>
      <StatusBar />
    </Layout>
  );
}
