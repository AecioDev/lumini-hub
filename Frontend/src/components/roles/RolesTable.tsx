import { useEffect, useMemo, useState } from "react";
import { DeleteOutlined, EditOutlined, PlusOutlined } from "@ant-design/icons";
import { Button, Card, Input, Space, Tag } from "antd";
import { Link } from "react-router-dom";
import { useAuth } from "@/contexts/AuthContext";
import { roleService } from "@/services/roles/role-service";
import { getApiErrorMessage } from "@/utils/api-error";
import { useFeedback } from "@/hooks/useFeedback";
import type { ApiRole } from "@/types/role";
import { ResponsiveTable, RowActionButtons, type RowAction } from "@/components/common/ResponsiveTable";

export function RolesTable() {
  const { hasPermission } = useAuth();
  const feedback = useFeedback();

  const [roles, setRoles] = useState<ApiRole[]>([]);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState("");

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
    try {
      await roleService.remove(role.id);
      feedback.success("Perfil excluído com sucesso.");
      load();
    } catch (error) {
      feedback.error(getApiErrorMessage(error, "Erro ao excluir perfil."));
    }
  };

  const rowActions = (role: ApiRole): RowAction[] => [
    ...(hasPermission("roles.edit")
      ? [
          {
            key: "edit",
            label: "Editar",
            icon: <EditOutlined />,
            to: `/settings/roles/${role.id}/edit`,
          },
        ]
      : []),
    ...(hasPermission("roles.delete")
      ? [
          {
            key: "delete",
            label: "Excluir",
            icon: <DeleteOutlined />,
            danger: true,
            confirm: {
              title: "Excluir perfil",
              description: `Tem certeza que deseja excluir o perfil "${role.name}"?`,
              onConfirm: () => handleDelete(role),
            },
          },
        ]
      : []),
  ];

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
            style={{ width: "min(220px, 42vw)" }}
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
      <ResponsiveTable<ApiRole>
        card={{
          title: (role) => role.name,
          subtitle: (role) => role.description,
          actions: rowActions,
        }}
        scroll={{ x: "max-content" }}
        rowKey="id"
        loading={loading}
        dataSource={filtered}
        pagination={false}
        columns={[
          { title: "ID", dataIndex: "id", width: 70 },
          { title: "Nome", dataIndex: "name" },
          { title: "Descrição", dataIndex: "description" },
          { title: "Ações", key: "actions", render: (_, role) => <RowActionButtons actions={rowActions(role)} /> },
        ]}
      />
    </Card>
  );
}
