import { Tabs } from "antd";
import { useSearchParams } from "react-router-dom";
import { RolesTable } from "@/components/roles/RolesTable";
import { PermissionsTable } from "@/components/permissions/PermissionsTable";
import { PageHeader } from "@/components/common/PageHeader";


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
      <PageHeader
        title="Perfis e Permissões"
        subtitle="Catálogo de perfis de acesso e permissões do sistema"
      />

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
