import { Result } from "antd";

interface AccessDeniedResultProps {
  subTitle?: string;
}

// Estado explícito de "sem acesso" pra dados que uma página precisa pra
// funcionar (não a rota inteira — isso é o RequirePermission — mas uma
// chamada específica que voltou 403). Evita deixar a tela meio carregada
// com Skeleton/formulário vazio quando a causa é falta de permissão.
export function AccessDeniedResult({
  subTitle = "Você não tem permissão para acessar esse recurso ou funcionalidade.",
}: AccessDeniedResultProps) {
  return <Result status="403" title="Sem acesso" subTitle={subTitle} />;
}
