import { useEffect, useState } from "react";
import { EditOutlined, EyeOutlined, PlusOutlined } from "@ant-design/icons";
import {
  Avatar,
  Button,
  Card,
  Select,
  Space,
  Tag,
} from "antd";
import { Link } from "react-router-dom";
import { userService } from "@/services/users/user-service";
import { roleService } from "@/services/roles/role-service";
import { useAuth } from "@/contexts/AuthContext";
import { UserDetailDrawer } from "@/components/users/UserDetailDrawer";
import { useFeedback } from "@/hooks/useFeedback";
import { useDebouncedValue } from "@/hooks/useDebouncedValue";
import { getAvatarGradient, getInitials, getTagColor } from "@/utils/avatar";
import type { ApiUser } from "@/types/auth";
import type { ApiRole } from "@/types/role";
import { PageHeader } from "@/components/common/PageHeader";
import { FilterButton, FilterDrawer } from "@/components/common/FilterDrawer";
import { FormField } from "@/components/common/FormField";
import { ListToolbar } from "@/components/common/ListToolbar";
import { ResponsiveTable, RowActionButtons, type RowAction } from "@/components/common/ResponsiveTable";


const PAGE_SIZE = 10;

// Filtros da gaveta (a busca por texto fica no cabeçalho da lista).
interface UserFilters {
  role_id?: number;
  is_active?: boolean;
}

const NO_FILTERS: UserFilters = {};

export function UsersListPage() {
  const { hasPermission } = useAuth();
  const feedback = useFeedback();

  const [users, setUsers] = useState<ApiUser[]>([]);
  const [totalRows, setTotalRows] = useState(0);
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState("");
  const [filters, setFilters] = useState<UserFilters>(NO_FILTERS);
  const [filterOpen, setFilterOpen] = useState(false);
  const [roles, setRoles] = useState<ApiRole[]>([]);
  const [detailUserId, setDetailUserId] = useState<number | null>(null);

  const debouncedSearch = useDebouncedValue(search.trim());
  const activeFilters = Object.values(filters).filter((value) => value !== undefined).length;

  // Perfis da gaveta de filtros; sem permissão de ver perfis, o filtro por perfil fica sem opções.
  useEffect(() => {
    roleService.list().then(setRoles).catch(() => undefined);
  }, []);

  // Mudou a busca ou os filtros: volta para a primeira página.
  useEffect(() => {
    setPage(1);
  }, [debouncedSearch, filters]);

  useEffect(() => {
    let cancelled = false;
    setLoading(true);
    userService
      .filter({
        search: debouncedSearch || undefined,
        role_id: filters.role_id,
        is_active: filters.is_active,
        page_no: page,
        page_size: PAGE_SIZE,
      })
      .then((result) => {
        if (cancelled) return;
        setUsers(result.users);
        setTotalRows(result.pagination?.totalRows ?? result.users.length);
      })
      .catch(() => {
        if (!cancelled) feedback.error("Erro ao carregar usuários.");
      })
      .finally(() => {
        if (!cancelled) setLoading(false);
      });
    return () => {
      cancelled = true;
    };
  }, [page, debouncedSearch, filters, feedback]);

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

      <Card>
        <ListToolbar
          searchPlaceholder="Buscar por nome, usuário ou e-mail..."
          onSearch={setSearch}
          filter={<FilterButton count={activeFilters} onClick={() => setFilterOpen(true)} />}
        />
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
          dataSource={users}
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

      <FilterDrawer
        open={filterOpen}
        onClose={() => setFilterOpen(false)}
        value={filters}
        empty={NO_FILTERS}
        onApply={setFilters}
      >
        {(draft, setField) => (
          <>
            <FormField label="Perfil" stacked>
              <Select
                allowClear
                placeholder="Todos os perfis"
                value={draft.role_id}
                onChange={(value) => setField("role_id", value)}
                options={roles.map((role) => ({ value: role.id, label: role.name }))}
              />
            </FormField>
            <FormField label="Status" stacked>
              <Select
                allowClear
                placeholder="Todos"
                value={draft.is_active}
                onChange={(value) => setField("is_active", value)}
                options={[
                  { value: true, label: "Ativo" },
                  { value: false, label: "Inativo" },
                ]}
              />
            </FormField>
          </>
        )}
      </FilterDrawer>

      <UserDetailDrawer userId={detailUserId} onClose={() => setDetailUserId(null)} />
    </div>
  );
}
