import { useMemo } from "react";
import { App as AntdApp } from "antd";

// Notification (canto superior direito) no lugar do `message` (toast
// centralizado) — mais legível e não empilha feio quando duas chamadas
// falham juntas (cada uma vira um card próprio, com título). Placement
// default vem de `notification={{ placement: "topRight" }}` no <App> raiz
// (src/App.tsx), então não precisa repetir aqui.
//
// Memoizado (useMemo) pra devolver a mesma referência entre renders — igual
// o `message`/`notification` que o antd já devolve via App.useApp() — assim
// dá pra usar `feedback` em array de dependências de useEffect sem disparar
// re-execução a cada render.
export function useFeedback() {
  const { notification } = AntdApp.useApp();

  return useMemo(
    () => ({
      success: (description: string, message = "Sucesso") =>
        notification.success({ message, description }),
      error: (description: string, message = "Erro") =>
        notification.error({ message, description }),
      info: (description: string, message = "Aviso") =>
        notification.info({ message, description }),
    }),
    [notification]
  );
}
