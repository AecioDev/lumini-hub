import { useEffect, useState } from "react";
import { DeleteOutlined, EditOutlined, PlusOutlined } from "@ant-design/icons";
import { Button, Card, Input, Select, Space, Tag } from "antd";
import { Link } from "react-router-dom";
import { useAuth } from "@/contexts/AuthContext";
import { permissionService } from "@/services/permissions/permission-service";
import { getApiErrorMessage } from "@/utils/api-error";
import { getTagColor } from "@/utils/avatar";
import { useFeedback } from "@/hooks/useFeedback";
import type { ApiPermission } from "@/types/permission";
import { BREAKPOINT_TABLET } from "@/components/layout/layout-constants";
import { useMediaQuery } from "@/hooks/useMediaQuery";
import { ResponsiveTable, RowActionButtons, type RowAction } from "@/components/common/ResponsiveTable";

const PAGE_SIZE = 10;

export function PermissionsTable() {
  const isPhone = !useMediaQuery(`(min-width: ${BREAKPOINT_TABLET}px)`);
  const { hasPermission, isDeveloper } = useAuth();
  const feedback = useFeedback();

  const [permissions, setPermissions] = useState<ApiPermission[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [name, setName] = useState<string | undefined>();
  const [module, setModule] = useState<string | undefined>();
  const [modules, setModules] = useState<string[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    permissionService.modules().then(setModules).catch(() => undefined);
  }, []);

  const load = () => {
    setLoading(true);
    permissionService
      .list({ page, limit: PAGE_SIZE, name, module })
      .then((result) => {
        setPermissions(result.data);
        setTotal(result.pagination.totalRows);
      })
      .catch(() => feedback.error("Erro ao carregar permissões."))
      .finally(() => setLoading(false));
  };

  useEffect(load, [page, name, module, feedback]);

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

  // No celular os filtros não cabem no cabeçalho do cartão: vão para uma linha que quebra, acima da lista.
  const controls = (
    <Space wrap>
      <Select
        placeholder="Módulo"
        allowClear
        style={{ width: "min(160px, 40vw)" }}
        value={module}
        onChange={(value) => {
          setModule(value);
          setPage(1);
        }}
        options={modules.map((m) => ({ value: m, label: m }))}
      />
      <Input.Search
        placeholder="Buscar por código/descrição..."
        allowClear
        style={{ width: "min(240px, 42vw)" }}
        onSearch={(value) => {
          setName(value || undefined);
          setPage(1);
        }}
      />
      {hasPermission("permissions.create") && (
        <Link to="/settings/permissions/create">
          <Button type="primary" icon={<PlusOutlined />}>
            Nova Permissão
          </Button>
        </Link>
      )}
    </Space>
  );

  return (
    <Card
      title={
        <span>
          Permissões
          <Tag bordered={false} style={{ marginLeft: 8 }}>
            {total}
          </Tag>
        </span>
      }
      extra={isPhone ? undefined : controls}
    >
      {isPhone && <div style={{ marginBottom: 12 }}>{controls}</div>}
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
    </Card>
  );
}
