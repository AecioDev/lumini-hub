import { Descriptions, Drawer, Empty, Space, Table, Tag, Typography } from "antd";
import dayjs from "dayjs";
import type { ApiAuditLog } from "@/types/audit";

interface AuditDetailDrawerProps {
  log: ApiAuditLog | null;
  onClose: () => void;
}

// Mudança de uma lista de objetos com id (ex.: permissões): o que entrou e o que saiu.
interface ListChange {
  adicionados: string[];
  removidos: string[];
}

const isListChange = (value: unknown): value is ListChange =>
  typeof value === "object" && value !== null && "adicionados" in value && "removidos" in value;

const isBeforeAfter = (value: unknown): value is { de: unknown; para: unknown } =>
  typeof value === "object" && value !== null && "de" in value && "para" in value;

function formatValue(value: unknown) {
  if (value === null || value === undefined || value === "") return <Typography.Text type="secondary">vazio</Typography.Text>;
  if (typeof value === "boolean") return value ? "Sim" : "Não";
  if (typeof value === "object") return <code style={{ wordBreak: "break-all" }}>{JSON.stringify(value)}</code>;
  return String(value);
}

function ChangesTable({ changes }: { changes: Record<string, unknown> }) {
  const rows = Object.entries(changes).map(([field, change]) => ({ field, change }));
  return (
    <Table
      size="small"
      pagination={false}
      rowKey="field"
      dataSource={rows}
      scroll={{ x: "max-content" }}
      columns={[
        { title: "Campo", dataIndex: "field", render: (field: string) => <code>{field}</code> },
        {
          title: "Antes",
          key: "before",
          render: (_, row) => (isBeforeAfter(row.change) ? formatValue(row.change.de) : "—"),
        },
        {
          title: "Depois",
          key: "after",
          render: (_, row) => {
            if (isListChange(row.change)) {
              return (
                <Space size={[4, 4]} wrap>
                  {row.change.adicionados.map((label) => (
                    <Tag key={`+${label}`} color="success">+ {label}</Tag>
                  ))}
                  {row.change.removidos.map((label) => (
                    <Tag key={`-${label}`} color="error">− {label}</Tag>
                  ))}
                </Space>
              );
            }
            return isBeforeAfter(row.change) ? formatValue(row.change.para) : formatValue(row.change);
          },
        },
      ]}
    />
  );
}

function SnapshotTable({ values }: { values: Record<string, unknown> }) {
  const rows = Object.entries(values).map(([field, value]) => ({ field, value }));
  return (
    <Table
      size="small"
      pagination={false}
      rowKey="field"
      dataSource={rows}
      scroll={{ x: "max-content" }}
      columns={[
        { title: "Campo", dataIndex: "field", render: (field: string) => <code>{field}</code> },
        { title: "Valor", dataIndex: "value", render: formatValue },
      ]}
    />
  );
}

// Gaveta com o detalhe de um registro de auditoria: quem, quando, e o que mudou (antes e depois).
export function AuditDetailDrawer({ log, onClose }: AuditDetailDrawerProps) {
  const details = log?.details;
  const changes = details?.changes as Record<string, unknown> | undefined;
  const after = details?.after as Record<string, unknown> | undefined;
  const before = details?.before as Record<string, unknown> | undefined;
  const truncated = details?.truncated as string[] | undefined;

  return (
    <Drawer title="Detalhes da auditoria" open={log !== null} onClose={onClose} width={560}>
      {log && (
        <Space direction="vertical" size={20} style={{ width: "100%" }}>
          <Descriptions column={1} size="small" bordered>
            <Descriptions.Item label="Data e hora">{dayjs(log.created_at).format("DD/MM/YYYY HH:mm:ss")}</Descriptions.Item>
            <Descriptions.Item label="Usuário">{log.username || (log.user_id ? `Usuário ${log.user_id}` : "—")}</Descriptions.Item>
            <Descriptions.Item label="Rota"><code>{log.action}</code></Descriptions.Item>
            <Descriptions.Item label="Registro">{log.entity_id ? `${log.entity_type} #${log.entity_id}` : log.entity_type}</Descriptions.Item>
            <Descriptions.Item label="IP">{log.ip_address}</Descriptions.Item>
          </Descriptions>

          {changes && (
            <div>
              <Typography.Title level={5}>O que mudou</Typography.Title>
              <ChangesTable changes={changes} />
            </div>
          )}
          {after && (
            <div>
              <Typography.Title level={5}>Dados criados</Typography.Title>
              <SnapshotTable values={after} />
            </div>
          )}
          {before && (
            <div>
              <Typography.Title level={5}>Dados excluídos</Typography.Title>
              <SnapshotTable values={before} />
            </div>
          )}
          {truncated && (
            <Typography.Paragraph type="secondary">
              O conteúdo era grande demais para ser guardado. Campos envolvidos: {truncated.join(", ")}.
            </Typography.Paragraph>
          )}
          {!changes && !after && !before && !truncated && (
            <Empty
              image={Empty.PRESENTED_IMAGE_SIMPLE}
              description="Sem detalhes do conteúdo para esta ação (o sistema não guarda o antes e o depois dela, ou nada mudou)."
            />
          )}
        </Space>
      )}
    </Drawer>
  );
}
