// src/app/(system)/settings/integrations/page.tsx
"use client";

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { PagePermissionGuard } from "@/components/layout/PagePermissionGuard";
import { IntegrationSettingsForm } from "@/components/settings/integrations/integration-settings-form";
import { SyncLogsTable } from "@/components/settings/integrations/tables/sync-logs-table";
import { WebhookEventsTable } from "@/components/settings/integrations/tables/webhook-events-table";

export default function IntegrationsSettingsPage() {
  return (
    <PagePermissionGuard
      requiredPermissions={["integrations.view"]}
      accessDeniedMessage="Você não tem permissão para visualizar as configurações de integrações."
    >
      <div className="space-y-6 px-4 py-6 md:px-8 md:py-8 bg-background text-foreground min-h-screen">
        <div>
          <h1 className="text-3xl font-bold tracking-tight lg:text-4xl">
            Integrações
          </h1>
          <p className="text-lg text-muted-foreground mt-1">
            Configure a sincronização entre a Loja Integrada e o ERP legado
          </p>
        </div>

        <Card className="shadow-lg border-primary/20">
          <CardHeader>
            <CardTitle className="text-xl font-semibold">Configurações</CardTitle>
            <CardDescription>
              Chaves da Loja Integrada e códigos do ERP legado usados na
              sincronização de pedidos e estoque.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <IntegrationSettingsForm />
          </CardContent>
        </Card>

        <Card className="shadow-lg border-primary/20">
          <CardHeader>
            <CardTitle className="text-xl font-semibold">
              Logs de Sincronização
            </CardTitle>
            <CardDescription>
              Histórico das sincronizações entre a Loja Integrada e o SQL Server
              legado.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <SyncLogsTable />
          </CardContent>
        </Card>

        <Card className="shadow-lg border-primary/20">
          <CardHeader>
            <CardTitle className="text-xl font-semibold">
              Eventos de Webhook
            </CardTitle>
            <CardDescription>
              Notificações recebidas da Loja Integrada, ainda pendentes de
              processamento em pedido de venda.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <WebhookEventsTable />
          </CardContent>
        </Card>
      </div>
    </PagePermissionGuard>
  );
}
