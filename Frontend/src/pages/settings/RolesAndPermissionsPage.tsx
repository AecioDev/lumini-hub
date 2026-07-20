import { Tabs, Typography } from "antd";
import { useSearchParams } from "react-router-dom";
import { RolesTable } from "@/components/roles/RolesTable";
import { PermissionsTable } from "@/components/permissions/PermissionsTable";

const { Title, Paragraph } = Typography;

const VALID_TABS = ["perfis", "permissoes"] as const;
type TabKey = (typeof VALID_TABS)[number];

function isTabKey(value: string | null): value is TabKey {
  return VALID_TABS.includes(value as TabKey);
}

// Ferramenta dev-only (gate admin.create_permissions) — cadastro de Perfis e
// de Permissões vivem juntos aqui porque é o mesmo item de menu real
// ("Perfis e Permissões" -> /settings/roles, ver internal/seeder/menu_item_seeder.go).
export function RolesAndPermissionsPage() {
  const [searchParams, setSearchParams] = useSearchParams();
  const activeTab = isTabKey(searchParams.get("tab")) ? searchParams.get("tab")! : "perfis";

  return (
    <div>
      <Title level={4} style={{ margin: "0 0 4px" }}>
        Perfis e Permissões
      </Title>
      <Paragraph type="secondary" style={{ marginBottom: 20 }}>
        Catálogo de perfis de acesso e permissões do sistema
      </Paragraph>

      <Tabs
        activeKey={activeTab}
        onChange={(key) => setSearchParams({ tab: key })}
        items={[
          { key: "perfis", label: "Perfis", children: <RolesTable /> },
          { key: "permissoes", label: "Permissões", children: <PermissionsTable /> },
        ]}
      />
    </div>
  );
}
