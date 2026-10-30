import { ArrowLeftOutlined } from "@ant-design/icons";
import { Button, Typography, theme } from "antd";
import type { ReactNode } from "react";
import { useNavigate } from "react-router-dom";
import { BREAKPOINT_TABLET } from "@/components/layout/layout-constants";
import { useMediaQuery } from "@/hooks/useMediaQuery";

const { Title, Paragraph } = Typography;

interface PageHeaderProps {
  title: string;
  subtitle?: string;
  /** Rota do botão "Voltar" (telas de cadastro e edição): geralmente a listagem de origem. */
  backTo?: string;
  /** Botões de ação à direita (ex.: "Novo Usuário"); o "Voltar", se houver, vem depois deles. */
  actions?: ReactNode;
}

// Cabeçalho padrão de toda página do app: título e descrição à esquerda; à direita, os botões de ação
// da página e o "Voltar" (quando a tela tem uma saída óbvia, como o formulário de cadastro).
//
// Fica fixo no topo da área que rola (o Content do AppLayout), para não perder o título e os botões ao
// rolar um formulário longo. O Content não tem padding próprio (o espaçamento fica num wrapper dentro
// dele), por isso `top: 0` gruda exatamente no topo; o `background` precisa ser explícito, senão o
// conteúdo que rola aparece atravessando o cabeçalho.
export function PageHeader({ title, subtitle, backTo, actions }: PageHeaderProps) {
  const navigate = useNavigate();
  const { token } = theme.useToken();
  // No celular o cabeçalho (título, descrição e botões em linhas separadas) ocuparia muito da tela se
  // ficasse fixo: ele só gruda no topo a partir do tablet.
  const sticky = useMediaQuery(`(min-width: ${BREAKPOINT_TABLET}px)`);

  return (
    <div
      style={{
        position: sticky ? "sticky" : "static",
        top: 0,
        zIndex: 10,
        display: "flex",
        justifyContent: "space-between",
        alignItems: "flex-start",
        flexWrap: "wrap",
        gap: 12,
        padding: "12px 0",
        marginBottom: 8,
        background: token.colorBgLayout,
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

      {(actions || backTo) && (
        <div style={{ display: "flex", alignItems: "center", flexWrap: "wrap", gap: 10 }}>
          {actions}
          {backTo && (
            <Button icon={<ArrowLeftOutlined />} onClick={() => navigate(backTo)}>
              Voltar
            </Button>
          )}
        </div>
      )}
    </div>
  );
}
