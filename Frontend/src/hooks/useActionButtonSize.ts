import { BREAKPOINT_TABLET } from "@/components/layout/layout-constants";
import { useMediaQuery } from "@/hooks/useMediaQuery";

// Tamanho dos botões de ícone das tabelas (ver, editar, excluir): compactos com mouse, maiores no celular
// e no tablet em pé, onde o dedo precisa de uma área de toque confortável (o "small" do antd tem só 24 px).
export function useActionButtonSize(): "small" | "large" {
  const hasMouseWidth = useMediaQuery(`(min-width: ${BREAKPOINT_TABLET}px)`);
  return hasMouseWidth ? "small" : "large";
}
