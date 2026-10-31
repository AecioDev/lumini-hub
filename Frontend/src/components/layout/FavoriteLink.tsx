import { StarFilled } from "@ant-design/icons";
import { Icon } from "@iconify/react";
import { Button, Tooltip } from "antd";
import { useLocation, useNavigate } from "react-router-dom";
import type { ApiUserMenuItem } from "@/types/menu";

interface FavoriteLinkProps {
  item: ApiUserMenuItem;
  /** Mostra o nome ao lado do ícone; sem ele fica só o ícone, com o nome numa dica. */
  showLabel: boolean;
  /** `header`: botão com a animação do cabeçalho; `list`: linha larga, usada na gaveta do menu. */
  variant: "header" | "list";
}

// Atalho de uma tela favorita (cabeçalho e topo da gaveta do menu).
export function FavoriteLink({ item, showLabel, variant }: FavoriteLinkProps) {
  const navigate = useNavigate();
  const { pathname } = useLocation();
  const active = item.href === pathname;
  const icon = item.icon ? <Icon icon={item.icon} width={18} height={18} /> : <StarFilled />;

  return (
    <Tooltip title={showLabel ? undefined : item.name}>
      <Button
        type="text"
        className={variant === "header" ? "header-favorite" : undefined}
        data-active={active || undefined}
        icon={<span className="header-favorite-icon">{icon}</span>}
        onClick={() => navigate(item.href)}
        aria-label={item.name}
        aria-current={active ? "page" : undefined}
        style={variant === "list" ? { width: "100%", justifyContent: "flex-start" } : undefined}
      >
        {showLabel ? item.name : null}
      </Button>
    </Tooltip>
  );
}
