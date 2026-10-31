import { theme, Typography } from "antd";
import { FavoriteLink } from "@/components/layout/FavoriteLink";
import { useMenuFavorites } from "@/hooks/useMenuFavorites";

// Favoritos no topo da gaveta do menu (celular e telas estreitas, onde o cabeçalho não tem espaço).
export function FavoritesShortcuts() {
  const { favorites } = useMenuFavorites();
  const { token } = theme.useToken();

  if (favorites.length === 0) return null;

  return (
    <nav
      aria-label="Favoritos"
      style={{ padding: "8px 8px 12px", borderBottom: `1px solid ${token.colorBorderSecondary}` }}
    >
      <Typography.Text type="secondary" style={{ display: "block", padding: "0 8px 4px", fontSize: 12 }}>
        Favoritos
      </Typography.Text>
      {favorites.map((item) => (
        <FavoriteLink key={item.id} item={item} showLabel variant="list" />
      ))}
    </nav>
  );
}
