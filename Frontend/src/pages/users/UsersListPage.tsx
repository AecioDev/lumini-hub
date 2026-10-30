import { useEffect, useMemo, useState } from "react";
import { EditOutlined, EyeOutlined, PlusOutlined } from "@ant-design/icons";
import {
  Avatar,
  Button,
  Card,
  Input,
  Space,
  Tag,
} from "antd";
import { Link } from "react-router-dom";
import { userService } from "@/services/users/user-service";
import { useAuth } from "@/contexts/AuthContext";
import { UserDetailDrawer } from "@/components/users/UserDetailDrawer";
import { useFeedback } from "@/hooks/useFeedback";
import { getAvatarGradient, getInitials, getTagColor } from "@/utils/avatar";
import type { ApiUser } from "@/types/auth";
import { PageHeader } from "@/components/common/PageHeader";
import { ResponsiveTable, RowActionButtons, type RowAction } from "@/components/common/ResponsiveTable";


const PAGE_SIZE = 10;

export function UsersListPage() {
  const { hasPermission } = useAuth();
  const feedback = useFeedback();

  const [users, setUsers] = useState<ApiUser[]>([]);
  const [totalRows, setTotalRows] = useState(0);
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState("");
  const [detailUserId, setDetailUserId] = useState<number | null>(null);

  useEffect(() => {
    setLoading(true);
    userService
      .list({ page, limit: PAGE_SIZE })
      .then((result) => {
        setUsers(result.users);
        setTotalRows(result.pagination?.totalRows ?? result.users.length);
      })
      .catch(() => feedback.error("Erro ao carregar usuários."))
      .finally(() => setLoading(false));
  }, [page, feedback]);

  const filteredUsers = useMemo(() => {
    const term = search.trim().toLowerCase();
    if (!term) return users;
    return users.filter(
      (user) =>
        user.name.toLowerCase().includes(term) ||
        user.username.toLowerCase().includes(term) ||
        (user.email ?? "").toLowerCase().includes(term)
    );
  }, [users, search]);

  const rowActions = (user: ApiUser): RowAction[] => [
    {
      key: "view",
      label: "Ver detalhes",
      icon: <EyeOutlined />,
      onClick: () => setDetailUserId(user.id),
    },
    ...(hasPermission("users.edit")
      ? [
          {
            key: "edit",
            label: "Editar",
            icon: <EditOutlined />,
            to: `/settings/users/${user.id}/edit`,
          },
        ]
      : []),
  ];

  const renderRole = (role?: string) =>
    role ? <Tag color={getTagColor(role)}>{role}</Tag> : "—";

  const renderStatus = (isActive: boolean) => (
    <Tag color={isActive ? "success" : "default"}>{isActive ? "Ativo" : "Inativo"}</Tag>
  );

  const renderIdentity = (user: ApiUser) => (
    <div style={{ display: "flex", alignItems: "center", gap: 10 }}>
      <Avatar style={{ background: getAvatarGradient(user.name) }}>
        {getInitials(user.name)}
      </Avatar>
      <div>
        <div style={{ fontSize: 13, fontWeight: 500 }}>{user.name}</div>
        <div style={{ fontSize: 12, opacity: 0.6 }}>{user.email}</div>
      </div>
    </div>
  );

  return (
    <div>
      <PageHeader
        title="Usuários do Sistema"
        subtitle="Gerencie contas, perfis e permissões da sua equipe"
        actions={hasPermission("users.create") && (
          <Link to="/settings/users/create">
            <Button type="primary" icon={<PlusOutlined />}>
              Novo Usuário
            </Button>
          </Link>
        )}
      />

      <Card
        title={
          <span>
            Todos os Usuários
            <Tag bordered={false} style={{ marginLeft: 8 }}>
              {totalRows}
            </Tag>
          </span>
        }
        extra={
          <Input.Search
            placeholder="Buscar na página atual..."
            allowClear
            style={{ width: "min(240px, 42vw)" }}
            onChange={(e) => setSearch(e.target.value)}
          />
        }
      >
        <ResponsiveTable<ApiUser>
          card={{
            title: renderIdentity,
            tags: (user) => (
              <Space size={4} wrap>
                {renderRole(user.role)}
                {renderStatus(user.is_active)}
              </Space>
            ),
            actions: rowActions,
          }}
          scroll={{ x: "max-content" }}
          rowKey="id"
          loading={loading}
          dataSource={filteredUsers}
          pagination={{
            current: page,
            pageSize: PAGE_SIZE,
            total: totalRows,
            onChange: setPage,
            showSizeChanger: false,
          }}
          columns={[
            { title: "ID", dataIndex: "id", width: 70 },
            { title: "Usuário", dataIndex: "name", render: (_, user) => renderIdentity(user) },
            { title: "Perfil", dataIndex: "role", render: renderRole },
            { title: "Status", dataIndex: "is_active", render: renderStatus },
            { title: "Ações", key: "actions", render: (_, user) => <RowActionButtons actions={rowActions(user)} /> },
          ]}
        />
      </Card>

      <UserDetailDrawer userId={detailUserId} onClose={() => setDetailUserId(null)} />
    </div>
  );
}
