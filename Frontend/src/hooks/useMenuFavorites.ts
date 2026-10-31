import { useCallback, useMemo, useState } from "react";
import { App } from "antd";
import { useAuth } from "@/contexts/AuthContext";
import type { ApiUserMenuItem } from "@/types/menu";

// Limite de favoritos por usuário; o backend repete este número (domain.MaxMenuFavorites).
export const MAX_MENU_FAVORITES = 5;

function indexById(items: ApiUserMenuItem[], into: Map<number, ApiUserMenuItem>) {
  for (const item of items) {
    into.set(item.id, item);
    if (item.children) indexById(item.children, into);
  }
  return into;
}

// Favoritos do menu do usuário logado (gravados no banco, valem em qualquer máquina). `favorites` vem na
// ordem do cabeçalho; `toggle` marca ou desmarca uma tela e avisa quando o limite é atingido.
export function useMenuFavorites() {
  const { menuItems, favoriteMenuItemIds, saveFavorites } = useAuth();
  const { message } = App.useApp();
  const [saving, setSaving] = useState(false);

  const favorites = useMemo(() => {
    const byId = indexById(menuItems, new Map());
    return favoriteMenuItemIds.flatMap((id) => byId.get(id) ?? []);
  }, [menuItems, favoriteMenuItemIds]);

  const isFavorite = useCallback((id: number) => favoriteMenuItemIds.includes(id), [favoriteMenuItemIds]);

  const toggle = useCallback(
    async (item: ApiUserMenuItem) => {
      if (saving) return;
      const marked = favoriteMenuItemIds.includes(item.id);
      if (!marked && favoriteMenuItemIds.length >= MAX_MENU_FAVORITES) {
        message.warning(`Você pode ter no máximo ${MAX_MENU_FAVORITES} favoritos. Remova um para adicionar outro.`);
        return;
      }
      const next = marked
        ? favoriteMenuItemIds.filter((id) => id !== item.id)
        : [...favoriteMenuItemIds, item.id];
      setSaving(true);
      try {
        await saveFavorites(next);
      } catch {
        message.error("Não foi possível salvar os favoritos. Tente novamente.");
      } finally {
        setSaving(false);
      }
    },
    [saving, favoriteMenuItemIds, saveFavorites, message]
  );

  return { favorites, isFavorite, toggle };
}
