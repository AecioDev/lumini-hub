import { useEffect, useState } from "react";
import { DeleteOutlined, EditOutlined, PlusOutlined } from "@ant-design/icons";
import { Button, Card, Input, Popconfirm, Select, Space, Table, Tag } from "antd";
import { Link } from "react-router-dom";
import { useAuth } from "@/contexts/AuthContext";
import { permissionService } from "@/services/permissions/permission-service";
import { getApiErrorMessage } from "@/utils/api-error";
import { getTagColor } from "@/utils/avatar";
import { useFeedback } from "@/hooks/useFeedback";
import type { ApiPermission } from "@/types/permission";

const PAGE_SIZE = 10;

export function PermissionsTable() {
  const { hasPermission } = useAuth();
  const feedback = useFeedback();

  const [permissions, setPermissions] = useState<ApiPermission[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [name, setName] = useState<string | undefined>();
  const [module, setModule] = useState<string | undefined>();
  const [modules, setModules] = useState<string[]>([]);
  const [loading, setLoading] = useState(true);
  const [deletingId, setDeletingId] = useState<number | null>(null);

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
    setDeletingId(permission.id);
    try {
      await permissionService.remove(permission.id);
      feedback.success("Permissão excluída com sucesso.");
      load();
    } catch (error) {
      feedback.error(getApiErrorMessage(error, "Erro ao excluir permissão."));
    } finally {
      setDeletingId(null);
    }
  };

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
      extra={
        <Space>
          <Select
            placeholder="Módulo"
            allowClear
            style={{ width: 160 }}
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
            style={{ width: 240 }}
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
      }
    >
      <Table<ApiPermission>
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
          { title: "Permissão", dataIndex: "permission" },
          { title: "Descrição", dataIndex: "description", ellipsis: true },
          {
            title: "Módulo",
            dataIndex: "module",
            render: (value: string) => <Tag color={getTagColor(value)}>{value}</Tag>,
          },
          {
            title: "Ações",
            key: "actions",
            render: (_, permission) => (
              <Space size={6}>
                {hasPermission("permissions.edit") && (
                  <Link to={`/settings/permissions/${permission.id}/edit`}>
                    <Button size="small" icon={<EditOutlined />} />
                  </Link>
                )}
                {hasPermission("permissions.delete") && (
                  <Popconfirm
                    title="Excluir permissão"
                    description={`Tem certeza que deseja excluir "${permission.permission}"?`}
                    okText="Excluir"
                    okButtonProps={{ danger: true, loading: deletingId === permission.id }}
                    cancelText="Cancelar"
                    onConfirm={() => void handleDelete(permission)}
                  >
                    <Button size="small" danger icon={<DeleteOutlined />} />
                  </Popconfirm>
                )}
              </Space>
            ),
          },
        ]}
      />
    </Card>
  );
}
