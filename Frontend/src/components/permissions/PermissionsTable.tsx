import { useEffect, useState } from "react";
import { DeleteOutlined, EditOutlined, PlusOutlined } from "@ant-design/icons";
import { Button, Card, Select, Tag } from "antd";
import { Link } from "react-router-dom";
import { useAuth } from "@/contexts/AuthContext";
import { permissionService } from "@/services/permissions/permission-service";
import { getApiErrorMessage } from "@/utils/api-error";
import { getTagColor } from "@/utils/avatar";
import { useFeedback } from "@/hooks/useFeedback";
import { useDebouncedValue } from "@/hooks/useDebouncedValue";
import type { ApiPermission } from "@/types/permission";
import { ResponsiveTable, RowActionButtons, type RowAction } from "@/components/common/ResponsiveTable";
import { ListToolbar } from "@/components/common/ListToolbar";
import { FilterButton, FilterDrawer } from "@/components/common/FilterDrawer";
import { FormField } from "@/components/common/FormField";

const PAGE_SIZE = 10;

// Filtros da gaveta (a busca por texto fica no cabeçalho da lista).
interface PermissionFilters {
  module?: string;
}

const NO_FILTERS: PermissionFilters = {};

export function PermissionsTable() {
  const { hasPermission, isDeveloper } = useAuth();
  const feedback = useFeedback();

  const [permissions, setPermissions] = useState<ApiPermission[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState("");
  const [filters, setFilters] = useState<PermissionFilters>(NO_FILTERS);
  const [filterOpen, setFilterOpen] = useState(false);
  const [modules, setModules] = useState<string[]>([]);
  const [loading, setLoading] = useState(true);

  const debouncedSearch = useDebouncedValue(search.trim());
  const activeFilters = Object.values(filters).filter((value) => value !== undefined).length;

  useEffect(() => {
    permissionService.modules().then(setModules).catch(() => undefined);
  }, []);

  // Mudou a busca ou os filtros: volta para a primeira página.
  useEffect(() => {
    setPage(1);
  }, [debouncedSearch, filters]);

  const load = () => {
    setLoading(true);
    permissionService
      .filter({
        search: debouncedSearch || undefined,
        module: filters.module,
        page_no: page,
        page_size: PAGE_SIZE,
      })
      .then((result) => {
        setPermissions(result.data);
        setTotal(result.pagination.totalRows);
      })
      .catch(() => feedback.error("Erro ao carregar permissões."))
      .finally(() => setLoading(false));
  };

  useEffect(load, [page, debouncedSearch, filters, feedback]);

  const handleDelete = async (permission: ApiPermission) => {
    try {
      await permissionService.remove(permission.id);
      feedback.success("Permissão excluída com sucesso.");
      load();
    } catch (error) {
      feedback.error(getApiErrorMessage(error, "Erro ao excluir permissão."));
    }
  };

  const rowActions = (permission: ApiPermission): RowAction[] => [
    ...(hasPermission("permissions.edit")
      ? [
          {
            key: "edit",
            label: "Editar",
            icon: <EditOutlined />,
            to: `/settings/permissions/${permission.id}/edit`,
          },
        ]
      : []),
    ...(hasPermission("permissions.delete")
      ? [
          {
            key: "delete",
            label: "Excluir",
            icon: <DeleteOutlined />,
            danger: true,
            confirm: {
              title: "Excluir permissão",
              description: `Tem certeza que deseja excluir "${permission.permission}"?`,
              onConfirm: () => handleDelete(permission),
            },
          },
        ]
      : []),
  ];

  return (
    <Card>
      <ListToolbar
        searchPlaceholder="Buscar por código ou descrição..."
        onSearch={setSearch}
        filter={<FilterButton count={activeFilters} onClick={() => setFilterOpen(true)} />}
        actions={
          hasPermission("permissions.create") && (
            <Link to="/settings/permissions/create">
              <Button type="primary" icon={<PlusOutlined />}>
                Nova Permissão
              </Button>
            </Link>
          )
        }
      />
      <ResponsiveTable<ApiPermission>
        card={{
          title: (permission) => (isDeveloper ? permission.permission : permission.description),
          subtitle: (permission) => (isDeveloper ? permission.description : undefined),
          tags: (permission) => <Tag color={getTagColor(permission.module)}>{permission.module}</Tag>,
          // ADMIN nunca edita nem exclui permissões (só o DEVELOP), então o cartão dele não tem ações.
          actions: isDeveloper ? rowActions : undefined,
        }}
        rowKey="id"
        loading={loading}
        dataSource={permissions}
        pagination={{
          current: page,
          pageSize: PAGE_SIZE,
          total,
          onChange: setPage,
          showSizeChanger: false,
        }}
        scroll={{ x: "max-content" }}
        columns={[
          { title: "Código", dataIndex: "id", width: 90 },
          // Código técnico (ex.: "admin.create_permissions") só faz sentido
          // pra quem mexe no catálogo de verdade — pro ADMIN é ruído.
          ...(isDeveloper ? [{ title: "Permissão", dataIndex: "permission" }] : []),
          { title: "Descrição", dataIndex: "description", ellipsis: true },
          {
            title: "Módulo",
            dataIndex: "module",
            render: (value: string) => <Tag color={getTagColor(value)}>{value}</Tag>,
          },
          // ADMIN nunca tem permissions.edit/permissions.delete (catálogo é
          // exclusivo do DEVELOP), então a coluna ficaria sempre vazia pra ele.
          ...(isDeveloper
            ? [
                {
                  title: "Ações",
                  key: "actions",
                  render: (_: unknown, permission: ApiPermission) => (
                    <RowActionButtons actions={rowActions(permission)} />
                  ),
                },
              ]
            : []),
        ]}
      />

      <FilterDrawer
        open={filterOpen}
        onClose={() => setFilterOpen(false)}
        value={filters}
        empty={NO_FILTERS}
        onApply={setFilters}
      >
        {(draft, setField) => (
          <FormField label="Módulo" stacked>
            <Select
              allowClear
              placeholder="Todos os módulos"
              value={draft.module}
              onChange={(value) => setField("module", value)}
              options={modules.map((m) => ({ value: m, label: m }))}
            />
          </FormField>
        )}
      </FilterDrawer>
    </Card>
  );
}
