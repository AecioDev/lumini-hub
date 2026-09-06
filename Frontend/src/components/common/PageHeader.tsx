import { ArrowLeftOutlined } from "@ant-design/icons";
import { Button, Typography, theme } from "antd";
import { useNavigate } from "react-router-dom";

const { Title, Paragraph } = Typography;

interface PageHeaderProps {
  title: string;
  subtitle?: string;
  /** Rota de destino do botão "Voltar" — geralmente a listagem que originou esta página. */
  backTo: string;
}

// Cabeçalho padrão das páginas de Cadastro/Edição — título + descrição à
// esquerda, botão "Voltar" à direita. Existe pra sempre ter uma saída da
// página mesmo quando o formulário abaixo não tem seu próprio botão
// Cancelar (ex.: aba "Permissões" de Editar Perfil, que não tem).
//
// Fixo no topo (position: sticky) em relação ao <Content> rolável de
// AppLayout.tsx — pedido do usuário pra não perder o título/Voltar ao
// rolar formulários longos (ex.: Identidade Visual). `background` precisa
// ser explícito com o token colorBgLayout, senão o conteúdo por trás
// aparece atravessando o cabeçalho fixo durante o scroll.
export function PageHeader({ title, subtitle, backTo }: PageHeaderProps) {
  const navigate = useNavigate();
  const { token } = theme.useToken();

  return (
    <div
      style={{
        display: "flex",
        justifyContent: "space-between",
        alignItems: "flex-start",
        gap: 16,
        marginBottom: 20,
        position: "sticky",
        top: 0,
        zIndex: 10,
        background: token.colorBgLayout,
        paddingBottom: 12,
      }}
    >
      <div>
        <Title level={4} style={{ margin: "0 0 4px" }}>
          {title}
        </Title>
        {subtitle && (
          <Paragraph type="secondary" style={{ margin: 0 }}>
            {subtitle}
          </Paragraph>
        )}
      </div>
      <Button icon={<ArrowLeftOutlined />} onClick={() => navigate(backTo)}>
        Voltar
      </Button>
    </div>
  );
}
