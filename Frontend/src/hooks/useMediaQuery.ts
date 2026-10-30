import { useSyncExternalStore } from "react";

// Diz se uma media query (ex.: "(min-width: 992px)") vale agora e acompanha as mudanças de largura da janela.
// Usa o valor real já no primeiro render (sem o "piscar" de começar errado e corrigir depois).
export function useMediaQuery(query: string): boolean {
  return useSyncExternalStore(
    (notify) => {
      const list = window.matchMedia(query);
      list.addEventListener("change", notify);
      return () => list.removeEventListener("change", notify);
    },
    () => window.matchMedia(query).matches,
    () => false
  );
}
