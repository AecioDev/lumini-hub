import { useEffect, useState } from "react";
import { App as AntdApp, Skeleton, Tabs, Typography } from "antd";
import { useAuth } from "@/contexts/AuthContext";
import { integrationConfigService } from "@/services/integrations/integration-config-service";
import { legacyLookupService } from "@/services/integrations/legacy-lookup-service";
import { IntegrationsOverview } from "@/components/integrations/IntegrationsOverview";
import { IntegrationSettingsForm } from "@/components/integrations/IntegrationSettingsForm";
import { SyncLogsTable } from "@/components/integrations/SyncLogsTable";
import { WebhookEventsTable } from "@/components/integrations/WebhookEventsTable";
import type { ApiIntegrationSettings, LegacyCompany, LegacyLocation } from "@/types/integration";

const { Title, Paragraph } = Typography;

export function IntegrationsPage() {
  const { hasPermission } = useAuth();
  const { message } = AntdApp.useApp();

  const [settings, setSettings] = useState<ApiIntegrationSettings | null>(null);
  const [loadingSettings, setLoadingSettings] = useState(true);
  const [locations, setLocations] = useState<LegacyLocation[]>([]);
  const [companies, setCompanies] = useState<LegacyCompany[]>([]);
  const [legacyAvailable, setLegacyAvailable] = useState<boolean | null>(null);

  useEffect(() => {
    let cancelled = false;

    integrationConfigService
      .getSettings()
      .then((result) => {
        if (!cancelled) setSettings(result);
      })
      .catch(() => message.error("Erro ao carregar configurações de integrações."))
      .finally(() => {
        if (!cancelled) setLoadingSettings(false);
      });

    Promise.all([legacyLookupService.getLocations(), legacyLookupService.getCompanies()])
      .then(([locationsResult, companiesResult]) => {
        if (cancelled) return;
        setLocations(locationsResult);
        setCompanies(companiesResult);
        setLegacyAvailable(true);
      })
      .catch(() => {
        if (!cancelled) setLegacyAvailable(false);
      });

    return () => {
      cancelled = true;
    };
  }, [message]);

  const canEdit = hasPermission("integrations.edit");

  return (
    <div>
      <Title level={4} style={{ margin: "0 0 4px" }}>
        Integrações
      </Title>
      <Paragraph type="secondary" style={{ marginBottom: 20 }}>
        Sincronização com a Loja Integrada e o ERP legado
      </Paragraph>

      <Tabs
        items={[
          {
            key: "overview",
            label: "Visão Geral",
            children: <IntegrationsOverview legacyAvailable={legacyAvailable} />,
          },
          {
            key: "settings",
            label: "Configurações",
            children: loadingSettings || !settings ? (
              <Skeleton active paragraph={{ rows: 8 }} />
            ) : (
              <IntegrationSettingsForm
                settings={settings}
                locations={locations}
                companies={companies}
                legacyAvailable={legacyAvailable}
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
