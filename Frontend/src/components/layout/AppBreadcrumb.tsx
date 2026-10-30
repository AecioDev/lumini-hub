import { Breadcrumb } from "antd";
import { Link, useLocation } from "react-router-dom";
import { buildBreadcrumb } from "@/utils/breadcrumb";
import type { ApiUserMenuItem } from "@/types/menu";

interface AppBreadcrumbProps {
  menuItems: ApiUserMenuItem[];
}

// Faixa com o caminho da tela atual, entre o cabeçalho e o corpo (que é o que rola): fica sempre à vista.
// O caminho vem do menu do usuário; telas de cadastro/edição acrescentam o próprio nome ao final.
export function AppBreadcrumb({ menuItems }: AppBreadcrumbProps) {
  const { pathname } = useLocation();
  const crumbs = buildBreadcrumb(menuItems, pathname);

  if (crumbs.length === 0) return null;

  return (
    <div style={{ flexShrink: 0, padding: "12px 24px 0" }}>
      <Breadcrumb
        items={crumbs.map((crumb) => ({
          title: crumb.href ? <Link to={crumb.href}>{crumb.title}</Link> : crumb.title,
        }))}
      />
    </div>
  );
}
