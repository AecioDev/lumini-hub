import { useEffect, useState } from "react";
import { Outlet, useLocation } from "react-router-dom";
import { Layout } from "antd";
import { useAuth } from "@/contexts/AuthContext";
import { useMediaQuery } from "@/hooks/useMediaQuery";
import { AppBreadcrumb } from "@/components/layout/AppBreadcrumb";
import { AppHeader } from "@/components/layout/AppHeader";
import { AppSider } from "@/components/layout/AppSider";
import { MobileMenu } from "@/components/layout/MobileMenu";
import { StatusBar } from "@/components/layout/StatusBar";
import { BREAKPOINT_DESKTOP, BREAKPOINT_TABLET } from "@/components/layout/layout-constants";

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
//
// O menu depende da largura (ver layout-constants.ts): em tela grande é o menu fixo, completo ou recolhido
// por preferência; em tela média é o menu fixo só com ícones; em tela estreita não há menu fixo. Abaixo da
// tela grande, o botão do cabeçalho abre uma gaveta com o menu completo.
export function AppLayout() {
  const [collapsed, setCollapsed] = useState(readCollapsed);
  const [drawerOpen, setDrawerOpen] = useState(false);
  const { menuItems } = useAuth();
  const { pathname } = useLocation();

  const isDesktop = useMediaQuery(`(min-width: ${BREAKPOINT_DESKTOP}px)`);
  const hasSider = useMediaQuery(`(min-width: ${BREAKPOINT_TABLET}px)`);

  // Ao escolher uma tela no menu da gaveta, ela se fecha.
  useEffect(() => {
    setDrawerOpen(false);
  }, [pathname]);

  const handleMenuButton = () => {
    if (!isDesktop) {
      setDrawerOpen(true);
      return;
    }
    setCollapsed((prev) => {
      const next = !prev;
      try {
        window.localStorage.setItem(SIDEBAR_STORAGE_KEY, next ? "1" : "0");
      } catch {
        // sem localStorage (janela privada etc.): o menu só não é lembrado
      }
      return next;
    });
  };

  return (
    <Layout style={{ height: "100vh", overflow: "hidden" }}>
      <Layout hasSider={hasSider} style={{ flex: 1, minHeight: 0 }}>
        {hasSider && <AppSider collapsed={isDesktop ? collapsed : true} menuItems={menuItems} />}

        <Layout>
          <AppHeader onToggleMenu={handleMenuButton} />
          <AppBreadcrumb menuItems={menuItems} />

          {/* Sem padding no Content: o PageHeader é sticky e precisa grudar no topo da área que rola.
              O espaçamento fica neste wrapper, que rola junto com a página. */}
          <Content style={{ overflowY: "auto", minHeight: 0 }}>
            <div style={{ padding: "8px 24px 24px" }}>
              <Outlet />
            </div>
          </Content>
        </Layout>
      </Layout>
      <StatusBar />

      <MobileMenu
        open={drawerOpen && !isDesktop}
        onClose={() => setDrawerOpen(false)}
        menuItems={menuItems}
      />
    </Layout>
  );
}
