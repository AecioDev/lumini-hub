import { Result } from "antd";
import { useLocation } from "react-router-dom";

interface PlaceholderPageProps {
  title?: string;
  subTitle?: string;
}

// Fallback genérico para qualquer rota que ainda não tem tela própria —
// inclusive hrefs vindos do menu dinâmico (tabela menu_items) que apontam
// para módulos ainda não implementados no frontend (Clientes, Vendas,
// Tarefas, Configurações, etc.).
export function PlaceholderPage({
  title = "Em construção",
  subTitle = "Esta tela ainda não foi implementada.",
}: PlaceholderPageProps) {
  const location = useLocation();

  return (
    <Result
      status="info"
      title={title}
      subTitle={`${subTitle} (${location.pathname})`}
    />
  );
}
