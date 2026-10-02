import { useCallback, useState, type ReactNode } from "react";
import { Alert, Button, Card, Select, Typography } from "antd";
import { CHROME_BG } from "@/theme/antd-theme";
import { useThemeMode } from "@/contexts/ThemeContext";
import { useAuth } from "@/contexts/AuthContext";
import { useCreateCompany } from "@/hooks/useCreateCompany";
import { CompanyForm } from "@/components/companies/forms/CompanyForm";
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
//
// Primeiro acesso (CFG-8.1.3): com ZERO empresas cadastradas não há o que
// escolher, e o cadastro em /settings/companies fica atrás deste mesmo Gate.
// Por isso, quem tem companies.create ganha o botão "Cadastrar Empresa", que
// troca o seletor pelo CompanyForm aqui dentro (estado local, sem rota). Ao
// salvar, refreshUser() recarrega o /me e o backend, vendo exatamente 1
// empresa ativa, já a devolve como ativa — o Gate some sozinho. Sem a
// permission, o usuário só vê a orientação de procurar um administrador.
export function CompanySelectionGate() {
  const { mode } = useThemeMode();
  const { user, setActiveCompany, refreshUser, hasPermission, logout } = useAuth();

  const [selectedId, setSelectedId] = useState<number | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [creating, setCreating] = useState(false);

  const companies = user?.visible_companies ?? [];
  const hasNoCompanies = companies.length === 0;
  const canCreateCompany = hasNoCompanies && hasPermission("companies.create");

  // O hook já trata a falha deste callback com mensagem própria.
  const handleCompanyCreated = useCallback(async () => {
    await refreshUser();
  }, [refreshUser]);
  const { submit: createCompany, submitting: creatingCompany } =
    useCreateCompany(handleCompanyCreated);

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

  const handleLogout = () => {
    logout().catch(() => undefined);
  };

  const shell = (children: ReactNode) => (
    <div
      style={{
        minHeight: "100vh",
        width: "100%",
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        padding: 24,
        background: CHROME_BG[mode],
      }}
    >
      {children}
    </div>
  );

  if (creating) {
    return shell(
      <div style={{ width: "100%", maxWidth: 600 }}>
        <div style={{ marginBottom: 24, display: "flex", justifyContent: "center" }}>
          <Logo />
        </div>
        <Text strong style={{ display: "block", fontSize: 16, marginBottom: 4, textAlign: "center" }}>
          Cadastre sua empresa
        </Text>
        <Text
          type="secondary"
          style={{ display: "block", marginBottom: 20, textAlign: "center" }}
        >
          Informe os dados da empresa para começar a usar o sistema.
        </Text>
        <CompanyForm
          submitting={creatingCompany}
          submitLabel="Cadastrar Empresa"
          onSubmit={(values) => void createCompany(values)}
          onCancel={() => setCreating(false)}
        />
      </div>
    );
  }

  return shell(
    <Card
      style={{ width: 420, borderRadius: 12, boxShadow: "0 8px 32px rgba(0,0,0,0.2)" }}
      styles={{ body: { padding: 40 } }}
    >
      <div style={{ marginBottom: 24 }}>
        <Logo />
      </div>

      {hasNoCompanies && !canCreateCompany ? (
        <>
          <Text strong style={{ display: "block", fontSize: 16, marginBottom: 4 }}>
            Nenhuma empresa cadastrada
          </Text>
          <Text type="secondary" style={{ display: "block", marginBottom: 20 }}>
            Peça a um administrador para cadastrar a empresa.
          </Text>
          <Button block onClick={handleLogout}>
            Sair
          </Button>
        </>
      ) : (
        <>
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

          {canCreateCompany && (
            <Button block style={{ marginTop: 12 }} onClick={() => setCreating(true)}>
              Cadastrar Empresa
            </Button>
          )}
        </>
      )}
    </Card>
  );
}
