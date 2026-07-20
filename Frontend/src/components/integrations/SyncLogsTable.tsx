import { useEffect, useState } from "react";
import { App as AntdApp, Card, Select, Table, Tag } from "antd";
import { syncLogService } from "@/services/integrations/sync-log-service";
import type { ApiSyncLog, SyncStatus } from "@/types/integration";

const STATUS_COLOR: Record<SyncStatus, string> = {
  success: "success",
  error: "error",
  pending: "processing",
};

const DIRECTION_LABEL: Record<string, string> = {
  li_to_sql: "Loja Integrada → ERP",
  sql_to_li: "ERP → Loja Integrada",
};

const PAGE_SIZE = 10;

export function SyncLogsTable() {
  const { message } = AntdApp.useApp();
  const [logs, setLogs] = useState<ApiSyncLog[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [status, setStatus] = useState<string | undefined>();
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    setLoading(true);
    syncLogService
      .filter({
        page_no: page,
        page_size: PAGE_SIZE,
        status,
        order_by_column: "created_at",
        is_asc: false,
      })
      .then((result) => {
        setLogs(result.data);
        setTotal(result.pagination.totalRows);
      })
      .catch(() => message.error("Erro ao carregar logs de sincronização."))
      .finally(() => setLoading(false));
  }, [page, status, message]);

  return (
    <Card
      title="Logs de Sincronização"
      extra={
        <Select
          placeholder="Status"
          allowClear
          style={{ width: 160 }}
          value={status}
          onChange={(value) => {
            setStatus(value);
            setPage(1);
          }}
          options={[
            { value: "success", label: "Sucesso" },
            { value: "error", label: "Erro" },
            { value: "pending", label: "Pendente" },
          ]}
        />
      }
    >
      <Table<ApiSyncLog>
        rowKey="id"
        loading={loading}
        dataSource={logs}
        pagination={{
          current: page,
          pageSize: PAGE_SIZE,
          total,
          onChange: setPage,
          showSizeChanger: false,
        }}
        scroll={{ x: "max-content" }}
        columns={[
          {
            title: "Data",
            dataIndex: "created_at",
            render: (value: string) => new Date(value).toLocaleString("pt-BR"),
          },
          {
            title: "Direção",
            dataIndex: "direction",
            render: (value: string) => DIRECTION_LABEL[value] ?? value,
          },
          { title: "Entidade", dataIndex: "entity_type" },
          { title: "Referência", dataIndex: "reference_id" },
          {
            title: "Status",
            dataIndex: "status",
            render: (value: SyncStatus) => <Tag color={STATUS_COLOR[value]}>{value}</Tag>,
          },
          { title: "Mensagem", dataIndex: "message", ellipsis: true },
        ]}
      />
    </Card>
  );
}
