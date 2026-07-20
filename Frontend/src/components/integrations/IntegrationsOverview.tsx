import { useEffect, useState } from "react";
import { Alert, Card, Col, Row, Skeleton, Statistic, Tag } from "antd";
import { syncLogService } from "@/services/integrations/sync-log-service";
import { webhookEventService } from "@/services/integrations/webhook-event-service";
import type { ApiSyncLog } from "@/types/integration";

interface Stats {
  totalSyncs: number;
  syncErrors: number;
  totalWebhooks: number;
  webhookErrors: number;
  lastSync: ApiSyncLog | null;
}

const DIRECTION_LABEL: Record<string, string> = {
  li_to_sql: "Loja Integrada → ERP",
  sql_to_li: "ERP → Loja Integrada",
};

const STATUS_COLOR: Record<string, string> = {
  success: "success",
  error: "error",
  pending: "processing",
};

interface IntegrationsOverviewProps {
  legacyAvailable: boolean | null;
}

// "Dashboard simples" da tela de Integrações: nada aqui é mock — os números
// vêm dos próprios endpoints de filtro (POST /sync-logs/filter e
// /webhook-events/filter) usando page_size=1 só para ler o total da
// paginação, sem precisar de um endpoint de estatísticas dedicado.
export function IntegrationsOverview({ legacyAvailable }: IntegrationsOverviewProps) {
  const [stats, setStats] = useState<Stats | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    let cancelled = false;

    Promise.all([
      syncLogService.filter({ page_size: 1 }),
      syncLogService.filter({ status: "error", page_size: 1 }),
      webhookEventService.filter({ page_size: 1 }),
      webhookEventService.filter({ status: "error", page_size: 1 }),
      syncLogService.filter({ page_size: 1, order_by_column: "created_at", is_asc: false }),
    ])
      .then(([total, errors, webhooksTotal, webhookErrors, latest]) => {
        if (cancelled) return;
        setStats({
          totalSyncs: total.pagination.totalRows,
          syncErrors: errors.pagination.totalRows,
          totalWebhooks: webhooksTotal.pagination.totalRows,
          webhookErrors: webhookErrors.pagination.totalRows,
          lastSync: latest.data[0] ?? null,
        });
      })
      .finally(() => {
        if (!cancelled) setLoading(false);
      });

    return () => {
      cancelled = true;
    };
  }, []);

  return (
    <div>
      {legacyAvailable === false && (
        <Alert
          type="warning"
          showIcon
          style={{ marginBottom: 16 }}
          message="ERP legado indisponível"
          description="Não foi possível conectar ao SQL Server do ERP legado agora. Os selects de local/empresa na aba Configurações vão ficar vazios até a conexão voltar."
        />
      )}

      {loading || !stats ? (
        <Skeleton active paragraph={{ rows: 4 }} />
      ) : (
        <>
          <Row gutter={[16, 16]}>
            <Col xs={24} sm={12} lg={6}>
              <Card>
                <Statistic title="Sincronizações" value={stats.totalSyncs} />
              </Card>
            </Col>
            <Col xs={24} sm={12} lg={6}>
              <Card>
                <Statistic
                  title="Sincronizações com erro"
                  value={stats.syncErrors}
                  valueStyle={stats.syncErrors > 0 ? { color: "#ff4d4f" } : undefined}
                />
              </Card>
            </Col>
            <Col xs={24} sm={12} lg={6}>
              <Card>
                <Statistic title="Webhooks recebidos" value={stats.totalWebhooks} />
              </Card>
            </Col>
            <Col xs={24} sm={12} lg={6}>
              <Card>
                <Statistic
                  title="Webhooks com erro"
                  value={stats.webhookErrors}
                  valueStyle={stats.webhookErrors > 0 ? { color: "#ff4d4f" } : undefined}
                />
              </Card>
            </Col>
          </Row>

          <Card title="Última sincronização" style={{ marginTop: 16 }}>
            {stats.lastSync ? (
              <div>
                <Tag color={STATUS_COLOR[stats.lastSync.status] ?? "default"}>
                  {stats.lastSync.status}
                </Tag>{" "}
                <strong>{stats.lastSync.entity_type}</strong> (
                {DIRECTION_LABEL[stats.lastSync.direction] ?? stats.lastSync.direction}) —{" "}
                {new Date(stats.lastSync.created_at).toLocaleString("pt-BR")}
                {stats.lastSync.message && (
                  <div style={{ marginTop: 6, opacity: 0.7 }}>{stats.lastSync.message}</div>
                )}
              </div>
            ) : (
              <span style={{ opacity: 0.6 }}>Nenhuma sincronização registrada ainda.</span>
            )}
          </Card>
        </>
      )}
    </div>
  );
}
