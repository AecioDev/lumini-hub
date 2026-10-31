import { useEffect, useState } from "react";
import { EyeOutlined } from "@ant-design/icons";
import { Card, DatePicker, Select, Tag } from "antd";
import dayjs, { type Dayjs } from "dayjs";
import { auditService } from "@/services/audit/audit-service";
import { useFeedback } from "@/hooks/useFeedback";
import { useDebouncedValue } from "@/hooks/useDebouncedValue";
import type { ApiAuditLog, AuditLogFilterRequest } from "@/types/audit";
import { PageHeader } from "@/components/common/PageHeader";
import { FilterButton, FilterDrawer } from "@/components/common/FilterDrawer";
import { FormField } from "@/components/common/FormField";
import { ListToolbar } from "@/components/common/ListToolbar";
import { ResponsiveTable, RowActionButtons, type RowAction } from "@/components/common/ResponsiveTable";
import { AuditDetailDrawer } from "@/components/audit/AuditDetailDrawer";

const PAGE_SIZE = 10;

type Method = NonNullable<AuditLogFilterRequest["method"]>;

// Filtros da gaveta (a busca por texto fica no cabeçalho da lista).
interface AuditFilters {
  entity_type?: string;
  method?: Method;
  range?: [Dayjs, Dayjs];
}

const NO_FILTERS: AuditFilters = {};

const METHOD_LABEL: Record<string, { text: string; color: string }> = {
  POST: { text: "Criou", color: "success" },
  PUT: { text: "Alterou", color: "processing" },
  PATCH: { text: "Alterou", color: "processing" },
  DELETE: { text: "Excluiu", color: "error" },
};

// Nomes das entidades como aparecem na auditoria (o primeiro trecho da rota); o que não está aqui
// aparece como veio.
const ENTITY_LABEL: Record<string, string> = {
  users: "Usuários",
  roles: "Perfis",
  permissions: "Permissões",
  "menu-items": "Itens do menu",
  companies: "Empresas",
  "company-fiscal-configs": "Configuração fiscal",
  customers: "Clientes",
  suppliers: "Fornecedores",
  integrations: "Integrações",
};

const entityLabel = (type: string) => ENTITY_LABEL[type] ?? type;

const formatDate = (iso: string) => dayjs(iso).format("DD/MM/YYYY HH:mm:ss");

export function AuditLogsPage() {
  const feedback = useFeedback();

  const [logs, setLogs] = useState<ApiAuditLog[]>([]);
  const [totalRows, setTotalRows] = useState(0);
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState("");
  const [filters, setFilters] = useState<AuditFilters>(NO_FILTERS);
  const [filterOpen, setFilterOpen] = useState(false);
  const [entityTypes, setEntityTypes] = useState<string[]>([]);
  const [detailLog, setDetailLog] = useState<ApiAuditLog | null>(null);

  const debouncedSearch = useDebouncedValue(search.trim());
  const activeFilters = Object.values(filters).filter((value) => value !== undefined).length;

  useEffect(() => {
    auditService.entityTypes().then(setEntityTypes).catch(() => undefined);
  }, []);

  // Mudou a busca ou os filtros: volta para a primeira página.
  useEffect(() => {
    setPage(1);
  }, [debouncedSearch, filters]);

  useEffect(() => {
    let cancelled = false;
    setLoading(true);
    auditService
      .filter({
        search: debouncedSearch || undefined,
        entity_type: filters.entity_type,
        method: filters.method,
        // O período é de dias inteiros no fuso do usuário; o fim vai exclusivo (início do dia seguinte).
        date_from: filters.range?.[0].startOf("day").toISOString(),
        date_to: filters.range?.[1].add(1, "day").startOf("day").toISOString(),
        page_no: page,
        page_size: PAGE_SIZE,
      })
      .then((result) => {
        if (cancelled) return;
        setLogs(result.data);
        setTotalRows(result.pagination?.totalRows ?? result.data.length);
      })
      .catch(() => {
        if (!cancelled) feedback.error("Erro ao carregar a auditoria.");
      })
      .finally(() => {
        if (!cancelled) setLoading(false);
      });
    return () => {
      cancelled = true;
    };
  }, [page, debouncedSearch, filters, feedback]);

  const renderWhat = (log: ApiAuditLog) => {
    const method = log.action.split(" ")[0];
    const label = METHOD_LABEL[method];
    return (
      <>
        {label && <Tag color={label.color}>{label.text}</Tag>}
        {entityLabel(log.entity_type)}
        {log.entity_id && <span style={{ opacity: 0.6 }}> #{log.entity_id}</span>}
      </>
    );
  };

  const rowActions = (log: ApiAuditLog): RowAction[] => [
    { key: "view", label: "Ver detalhes", icon: <EyeOutlined />, onClick: () => setDetailLog(log) },
  ];

  // Campos que a alteração mexeu (o antes e depois completo fica na gaveta de detalhes).
  const changedFields = (log: ApiAuditLog) => {
    const changes = log.details?.changes;
    return changes ? Object.keys(changes).join(", ") : "";
  };

  const renderUser = (log: ApiAuditLog) => log.username || (log.user_id ? `Usuário ${log.user_id}` : "—");

  return (
    <div>
      <PageHeader title="Auditoria" subtitle="Quem alterou o quê no sistema, e quando" />

      <Card>
        <ListToolbar
          searchPlaceholder="Buscar por usuário ou ação..."
          onSearch={setSearch}
          filter={<FilterButton count={activeFilters} onClick={() => setFilterOpen(true)} />}
        />
        <ResponsiveTable<ApiAuditLog>
          card={{
            title: renderWhat,
            subtitle: (log) => `${renderUser(log)} · ${formatDate(log.created_at)}`,
            actions: rowActions,
            fields: [
              { label: "Campos alterados", value: changedFields },
              { label: "Rota", value: (log) => log.action },
              { label: "IP", value: (log) => log.ip_address },
            ],
          }}
          scroll={{ x: "max-content" }}
          rowKey="id"
          loading={loading}
          dataSource={logs}
          pagination={{
            current: page,
            pageSize: PAGE_SIZE,
            total: totalRows,
            onChange: setPage,
            showSizeChanger: false,
          }}
          columns={[
            { title: "ID", dataIndex: "id", width: 80 },
            { title: "Data e hora", dataIndex: "created_at", render: formatDate },
            { title: "Usuário", key: "user", render: (_, log) => renderUser(log) },
            { title: "O que foi feito", key: "what", render: (_, log) => renderWhat(log) },
            { title: "Campos alterados", key: "fields", render: (_, log) => changedFields(log) || "—" },
            { title: "Rota", dataIndex: "action", render: (action: string) => <code>{action}</code> },
            { title: "IP", dataIndex: "ip_address" },
            { title: "Ações", key: "actions", render: (_, log) => <RowActionButtons actions={rowActions(log)} /> },
          ]}
        />
      </Card>

      <FilterDrawer
        open={filterOpen}
        onClose={() => setFilterOpen(false)}
        value={filters}
        empty={NO_FILTERS}
        onApply={setFilters}
      >
        {(draft, setField) => (
          <>
            <FormField label="Entidade" stacked>
              <Select
                allowClear
                placeholder="Todas as entidades"
                value={draft.entity_type}
                onChange={(value) => setField("entity_type", value)}
                options={entityTypes.map((type) => ({ value: type, label: entityLabel(type) }))}
              />
            </FormField>
            <FormField label="Ação" stacked>
              <Select
                allowClear
                placeholder="Todas"
                value={draft.method}
                onChange={(value) => setField("method", value)}
                options={[
                  { value: "POST", label: "Criou" },
                  { value: "PUT", label: "Alterou" },
                  { value: "DELETE", label: "Excluiu" },
                ]}
              />
            </FormField>
            <FormField label="Período" stacked>
              <DatePicker.RangePicker
                style={{ width: "100%" }}
                format="DD/MM/YYYY"
                value={draft.range ?? null}
                onChange={(range) => setField("range", range?.[0] && range[1] ? [range[0], range[1]] : undefined)}
              />
            </FormField>
          </>
        )}
      </FilterDrawer>

      <AuditDetailDrawer log={detailLog} onClose={() => setDetailLog(null)} />
    </div>
  );
}
