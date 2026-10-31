import { FavoriteLink } from "@/components/layout/FavoriteLink";
import { BREAKPOINT_TABLET } from "@/components/layout/layout-constants";
import { useMenuFavorites } from "@/hooks/useMenuFavorites";
import { useMediaQuery } from "@/hooks/useMediaQuery";

// Largura (px) a partir da qual há espaço para o nome ao lado do ícone de cada favorito.
const BREAKPOINT_LABELS = 1280;

// Favoritos no meio do cabeçalho. Abaixo do tablet não cabem: ali aparecem no topo da gaveta do menu
// (FavoritesShortcuts). Até o tamanho "laptop" só os ícones, com o nome numa dica.
export function HeaderFavorites() {
  const { favorites } = useMenuFavorites();
  const hasRoom = useMediaQuery(`(min-width: ${BREAKPOINT_TABLET}px)`);
  const showLabels = useMediaQuery(`(min-width: ${BREAKPOINT_LABELS}px)`);

  if (!hasRoom || favorites.length === 0) return null;

  return (
    <nav aria-label="Favoritos" style={{ display: "flex", alignItems: "center", gap: 4 }}>
      {favorites.map((item) => (
        <FavoriteLink key={item.id} item={item} showLabel={showLabels} variant="header" />
      ))}
    </nav>
  );
}
