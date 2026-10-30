import { useEffect, useState } from "react";
import { Skeleton, Tabs } from "antd";
import { useAuth } from "@/contexts/AuthContext";
import { integrationConfigService } from "@/services/integrations/integration-config-service";
import { IntegrationsOverview } from "@/components/integrations/IntegrationsOverview";
import { IntegrationSettingsForm } from "@/components/integrations/IntegrationSettingsForm";
import { SyncLogsTable } from "@/components/integrations/SyncLogsTable";
import { WebhookEventsTable } from "@/components/integrations/WebhookEventsTable";
import { useFeedback } from "@/hooks/useFeedback";
import type { ApiIntegrationSettings } from "@/types/integration";
import { PageHeader } from "@/components/common/PageHeader";


export function IntegrationsPage() {
  const { hasPermission } = useAuth();
  const feedback = useFeedback();

  const [settings, setSettings] = useState<ApiIntegrationSettings | null>(null);
  const [loadingSettings, setLoadingSettings] = useState(true);

  useEffect(() => {
    let cancelled = false;

    integrationConfigService
      .getSettings()
      .then((result) => {
        if (!cancelled) setSettings(result);
      })
      .catch(() => feedback.error("Erro ao carregar configurações de integrações."))
      .finally(() => {
        if (!cancelled) setLoadingSettings(false);
      });

    return () => {
      cancelled = true;
    };
  }, [feedback]);

  const canEdit = hasPermission("integrations.edit");

  return (
    <div>
      <PageHeader
        title="Integrações"
        subtitle="Sincronização com a Loja Integrada"
      />

      <Tabs
        items={[
          {
            key: "overview",
            label: "Visão Geral",
            children: <IntegrationsOverview />,
          },
          {
            key: "settings",
            label: "Configurações",
            children: loadingSettings || !settings ? (
              <Skeleton active paragraph={{ rows: 8 }} />
            ) : (
              <IntegrationSettingsForm
                settings={settings}
                canEdit={canEdit}
                onSaved={setSettings}
              />
            ),
          },
          {
            key: "logs",
            label: "Logs",
            children: (
              <div style={{ display: "flex", flexDirection: "column", gap: 16 }}>
                <SyncLogsTable />
                <WebhookEventsTable />
              </div>
            ),
          },
        ]}
      />
    </div>
  );
}
