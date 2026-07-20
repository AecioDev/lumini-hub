import { useEffect, useMemo, useState } from "react";
import { EditOutlined, EyeOutlined, PlusOutlined } from "@ant-design/icons";
import {
  Avatar,
  Button,
  Card,
  Input,
  Space,
  Table,
  Tag,
  Tooltip,
  Typography,
} from "antd";
import { Link } from "react-router-dom";
import { userService } from "@/services/users/user-service";
import { useAuth } from "@/contexts/AuthContext";
import { UserDetailDrawer } from "@/components/users/UserDetailDrawer";
import { useFeedback } from "@/hooks/useFeedback";
import { getAvatarGradient, getInitials, getTagColor } from "@/utils/avatar";
import type { ApiUser } from "@/types/auth";

const { Title, Paragraph } = Typography;

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

  return (
    <div>
      <div
        style={{
          display: "flex",
          alignItems: "center",
          justifyContent: "space-between",
          marginBottom: 20,
          flexWrap: "wrap",
          gap: 12,
        }}
      >
        <div>
          <Title level={4} style={{ margin: "0 0 4px" }}>
            Usuários do Sistema
          </Title>
          <Paragraph type="secondary" style={{ margin: 0 }}>
            Gerencie contas, perfis e permissões da sua equipe
          </Paragraph>
        </div>
        {hasPermission("users.create") && (
          <Link to="/settings/users/create">
            <Button type="primary" icon={<PlusOutlined />}>
              Novo Usuário
            </Button>
          </Link>
        )}
      </div>

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
            style={{ width: 240 }}
            onChange={(e) => setSearch(e.target.value)}
          />
        }
      >
        <Table<ApiUser>
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
            {
              title: "Usuário",
              dataIndex: "name",
              render: (_, user) => (
                <div style={{ display: "flex", alignItems: "center", gap: 10 }}>
                  <Avatar style={{ background: getAvatarGradient(user.name) }}>
                    {getInitials(user.name)}
                  </Avatar>
                  <div>
                    <div style={{ fontSize: 13, fontWeight: 500 }}>{user.name}</div>
                    <div style={{ fontSize: 12, opacity: 0.6 }}>{user.email}</div>
                  </div>
                </div>
              ),
            },
            {
              title: "Perfil",
              dataIndex: "role",
              render: (role?: string) =>
                role ? <Tag color={getTagColor(role)}>{role}</Tag> : "—",
            },
            {
              title: "Status",
              dataIndex: "is_active",
              render: (isActive: boolean) => (
                <Tag color={isActive ? "success" : "default"}>
                  {isActive ? "Ativo" : "Inativo"}
                </Tag>
              ),
            },
            {
              title: "Ações",
              key: "actions",
              render: (_, user) => (
                <Space size={6}>
                  <Tooltip title="Ver detalhes">
                    <Button
                      size="small"
                      icon={<EyeOutlined />}
                      onClick={() => setDetailUserId(user.id)}
                    />
                  </Tooltip>
                  {hasPermission("users.edit") && (
                    <Tooltip title="Editar">
                      <Link to={`/settings/users/${user.id}/edit`}>
                        <Button size="small" icon={<EditOutlined />} />
                      </Link>
                    </Tooltip>
                  )}
                </Space>
              ),
            },
          ]}
        />
      </Card>

      <UserDetailDrawer userId={detailUserId} onClose={() => setDetailUserId(null)} />
    </div>
  );
}
