import { useEffect, useState } from "react";
import { Card, Select, Table, Tag } from "antd";
import { webhookEventService } from "@/services/integrations/webhook-event-service";
import { useFeedback } from "@/hooks/useFeedback";
import type { ApiWebhookEvent, WebhookStatus } from "@/types/integration";

const STATUS_COLOR: Record<WebhookStatus, string> = {
  received: "processing",
  processed: "success",
  error: "error",
};

const PAGE_SIZE = 10;

export function WebhookEventsTable() {
  const feedback = useFeedback();
  const [events, setEvents] = useState<ApiWebhookEvent[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [status, setStatus] = useState<string | undefined>();
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    setLoading(true);
    webhookEventService
      .filter({
        page_no: page,
        page_size: PAGE_SIZE,
        status,
        order_by_column: "received_at",
        is_asc: false,
      })
      .then((result) => {
        setEvents(result.data);
        setTotal(result.pagination.totalRows);
      })
      .catch(() => feedback.error("Erro ao carregar eventos de webhook."))
      .finally(() => setLoading(false));
  }, [page, status, feedback]);

  return (
    <Card
      title="Webhooks Recebidos"
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
            { value: "received", label: "Recebido" },
            { value: "processed", label: "Processado" },
            { value: "error", label: "Erro" },
          ]}
        />
      }
    >
      <Table<ApiWebhookEvent>
        rowKey="id"
        loading={loading}
        dataSource={events}
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
            title: "Recebido em",
            dataIndex: "received_at",
            render: (value: string) => new Date(value).toLocaleString("pt-BR"),
          },
          { title: "Origem", dataIndex: "source" },
          { title: "Tipo", dataIndex: "event_type" },
          {
            title: "Status",
            dataIndex: "status",
            render: (value: WebhookStatus) => <Tag color={STATUS_COLOR[value]}>{value}</Tag>,
          },
          { title: "Erro", dataIndex: "error_message", ellipsis: true },
        ]}
      />
    </Card>
  );
}
