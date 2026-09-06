import { useState } from "react";
import { Alert, Button, Card, Select, Typography } from "antd";
import { CHROME_BG } from "@/theme/antd-theme";
import { useThemeMode } from "@/contexts/ThemeContext";
import { useAuth } from "@/contexts/AuthContext";
import { Logo } from "@/components/layout/Logo";

const { Text } = Typography;

// Bloqueia o acesso ao resto do sistema até o usuário escolher com qual
// Empresa quer operar — obrigatório pra quem enxerga mais de uma (usuário
// master, ou com companies.hierarchy.view; ver ResolveActiveCompany no
// backend). Renderizado pelo ProtectedRoute no lugar do <Outlet/> enquanto
// user.requires_company_selection for true.
//
// Lista vem de user.visible_companies (ApiUserDetail, login/refresh/me),
// NÃO de GET /companies — de propósito: aquele endpoint exige
// companies.view, uma permission de administração do cadastro que a
// maioria dos perfis operacionais nunca tem. Um usuário master sem essa
// permission ficaria preso aqui pra sempre (esta tela bloqueia qualquer
// rota até a escolha ser feita, sem saída visível) se dependesse dela —
// mesmo bug já corrigido no ActiveCompanySwitcher.tsx, replicado aqui.
export function CompanySelectionGate() {
  const { mode } = useThemeMode();
  const { user, setActiveCompany } = useAuth();

  const [selectedId, setSelectedId] = useState<number | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const companies = user?.visible_companies ?? [];

  const handleConfirm = async () => {
    if (!selectedId) return;
    setSubmitting(true);
    setError(null);
    try {
      await setActiveCompany(selectedId);
    } catch {
      setError("Erro ao definir a empresa ativa. Tente novamente.");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div
      style={{
        minHeight: "100vh",
        width: "100%",
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        background: CHROME_BG[mode],
      }}
    >
      <Card
        style={{ width: 420, borderRadius: 12, boxShadow: "0 8px 32px rgba(0,0,0,0.2)" }}
        styles={{ body: { padding: 40 } }}
      >
        <div style={{ marginBottom: 24 }}>
          <Logo />
        </div>

        <Text strong style={{ display: "block", fontSize: 16, marginBottom: 4 }}>
          Escolha uma empresa
        </Text>
        <Text type="secondary" style={{ display: "block", marginBottom: 20 }}>
          Selecione com qual empresa você deseja trabalhar nesta sessão.
        </Text>

        <Select
          style={{ width: "100%", marginBottom: 16 }}
          placeholder="Selecione uma empresa"
          value={selectedId ?? undefined}
          onChange={setSelectedId}
          options={companies.map((company) => ({
            value: company.id,
            label: company.name,
          }))}
        />

        {error && <Alert type="error" message={error} showIcon style={{ marginBottom: 16 }} />}

        <Button
          type="primary"
          block
          disabled={!selectedId}
          loading={submitting}
          onClick={() => void handleConfirm()}
        >
          Continuar
        </Button>
      </Card>
    </div>
  );
}
