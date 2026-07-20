import { useEffect, useMemo, useState } from "react";
import { DeleteOutlined, EditOutlined, PlusOutlined } from "@ant-design/icons";
import { Button, Card, Input, Popconfirm, Space, Table, Tag } from "antd";
import { Link } from "react-router-dom";
import { useAuth } from "@/contexts/AuthContext";
import { roleService } from "@/services/roles/role-service";
import { getApiErrorMessage } from "@/utils/api-error";
import { useFeedback } from "@/hooks/useFeedback";
import type { ApiRole } from "@/types/role";

export function RolesTable() {
  const { hasPermission } = useAuth();
  const feedback = useFeedback();

  const [roles, setRoles] = useState<ApiRole[]>([]);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState("");
  const [deletingId, setDeletingId] = useState<number | null>(null);

  const load = () => {
    setLoading(true);
    roleService
      .list()
      .then(setRoles)
      .catch(() => feedback.error("Erro ao carregar perfis."))
      .finally(() => setLoading(false));
  };

  useEffect(load, [feedback]);

  const filtered = useMemo(() => {
    const term = search.trim().toLowerCase();
    if (!term) return roles;
    return roles.filter((role) => role.name.toLowerCase().includes(term));
  }, [roles, search]);

  const handleDelete = async (role: ApiRole) => {
    setDeletingId(role.id);
    try {
      await roleService.remove(role.id);
      feedback.success("Perfil excluído com sucesso.");
      load();
    } catch (error) {
      feedback.error(getApiErrorMessage(error, "Erro ao excluir perfil."));
    } finally {
      setDeletingId(null);
    }
  };

  return (
    <Card
      title={
        <span>
          Perfis
          <Tag bordered={false} style={{ marginLeft: 8 }}>
            {roles.length}
          </Tag>
        </span>
      }
      extra={
        <Space>
          <Input.Search
            placeholder="Buscar..."
            allowClear
            style={{ width: 220 }}
            onChange={(e) => setSearch(e.target.value)}
          />
          {hasPermission("roles.create") && (
            <Link to="/settings/roles/create">
              <Button type="primary" icon={<PlusOutlined />}>
                Novo Perfil
              </Button>
            </Link>
          )}
        </Space>
      }
    >
      <Table<ApiRole>
        rowKey="id"
        loading={loading}
        dataSource={filtered}
        pagination={false}
        columns={[
          { title: "Nome", dataIndex: "name" },
          { title: "Descrição", dataIndex: "description" },
          {
            title: "Ações",
            key: "actions",
            render: (_, role) => (
              <Space size={6}>
                {hasPermission("roles.edit") && (
                  <Link to={`/settings/roles/${role.id}/edit`}>
                    <Button size="small" icon={<EditOutlined />} />
                  </Link>
                )}
                {hasPermission("roles.delete") && (
                  <Popconfirm
                    title="Excluir perfil"
                    description={`Tem certeza que deseja excluir o perfil "${role.name}"?`}
                    okText="Excluir"
                    okButtonProps={{ danger: true, loading: deletingId === role.id }}
                    cancelText="Cancelar"
                    onConfirm={() => void handleDelete(role)}
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
