import { useEffect, useState } from "react";
import { Icon } from "@iconify/react";
import { Drawer, List, Skeleton, Typography } from "antd";
import { useAuth } from "@/contexts/AuthContext";
import { useFeedback } from "@/hooks/useFeedback";
import { companyService } from "@/services/companies/company-service";
import { getApiErrorMessage } from "@/utils/api-error";
import type { ApiCompany } from "@/types/company";

const { Text } = Typography;

interface ActiveCompanySwitcherProps {
  collapsed: boolean;
}

// Mostra a Company ativa no topo do sidebar (onde antes ficava a marca
// Lumini Hub — movida pro rodapé, ver AppLayout.tsx). Sem logo de verdade
// ainda (CompanyVisualConfig / CFG-4 não existe), usa um ícone genérico no
// lugar. Só quem enxerga mais de uma Company (master, ou com
// companies.hierarchy.view — ver ResolveVisibleCompanyIDs no backend) pode
// clicar pra trocar; pra quem só tem uma, fica só informativo, sem o
// destaque de "clicável".
export function ActiveCompanySwitcher({ collapsed }: ActiveCompanySwitcherProps) {
  const { user, setActiveCompany } = useAuth();
  const feedback = useFeedback();

  const [companies, setCompanies] = useState<ApiCompany[]>([]);
  const [loading, setLoading] = useState(true);
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [switching, setSwitching] = useState<number | null>(null);

  useEffect(() => {
    companyService
      .list()
      .then(setCompanies)
      .catch(() => setCompanies([]))
      .finally(() => setLoading(false));
  }, []);

  const activeCompany = companies.find((c) => c.id === user?.active_company_id);
  const canSwitch = companies.length > 1;

  const handleSelect = async (companyId: number) => {
    if (companyId === user?.active_company_id) {
      setDrawerOpen(false);
      return;
    }
    setSwitching(companyId);
    try {
      await setActiveCompany(companyId);
      feedback.success("Empresa ativa alterada com sucesso.");
      setDrawerOpen(false);
    } catch (error) {
      feedback.error(getApiErrorMessage(error, "Erro ao trocar de empresa."));
    } finally {
      setSwitching(null);
    }
  };

  if (loading) {
    return (
      <div style={{ padding: collapsed ? "20px 0" : "20px 16px" }}>
        <Skeleton.Avatar active size={22} />
      </div>
    );
  }

  return (
    <>
      <div
        onClick={canSwitch ? () => setDrawerOpen(true) : undefined}
        style={{
          display: "flex",
          alignItems: "center",
          justifyContent: collapsed ? "center" : "flex-start",
          gap: 10,
          padding: collapsed ? "20px 0" : "20px 16px",
          borderBottom: "1px solid rgba(255,255,255,0.08)",
          cursor: canSwitch ? "pointer" : "default",
        }}
      >
        <Icon icon="ph:buildings" width={22} height={22} style={{ flexShrink: 0, opacity: 0.85 }} />
        {!collapsed && (
          <>
            <span
              style={{
                fontSize: 14,
                fontWeight: 600,
                whiteSpace: "nowrap",
                overflow: "hidden",
                textOverflow: "ellipsis",
              }}
            >
              {activeCompany?.trade_name || activeCompany?.legal_name || "—"}
            </span>
            {canSwitch && (
              <Icon
                icon="ph:caret-up-down"
                width={14}
                height={14}
                style={{ flexShrink: 0, opacity: 0.5 }}
              />
            )}
          </>
        )}
      </div>

      <Drawer
        title="Trocar de Empresa"
        open={drawerOpen}
        onClose={() => setDrawerOpen(false)}
        width={360}
      >
        <List
          dataSource={companies}
          renderItem={(company) => (
            <List.Item
              onClick={() => void handleSelect(company.id)}
              style={{
                cursor: "pointer",
                borderRadius: 8,
                paddingLeft: 12,
                paddingRight: 12,
                background:
                  company.id === user?.active_company_id ? "rgba(0,0,0,0.04)" : undefined,
              }}
            >
              <List.Item.Meta
                avatar={<Icon icon="ph:buildings" width={20} height={20} />}
                title={company.trade_name || company.legal_name}
                description={company.legal_name}
              />
              {switching === company.id && <Text type="secondary">Trocando...</Text>}
              {company.id === user?.active_company_id && (
                <Icon icon="ph:check-circle-fill" width={18} height={18} color="#52c41a" />
              )}
            </List.Item>
          )}
        />
      </Drawer>
    </>
  );
}
