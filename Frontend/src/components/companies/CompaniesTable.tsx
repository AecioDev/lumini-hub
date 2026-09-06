import { useEffect, useMemo, useState } from "react";
import { DeleteOutlined, EditOutlined, PlusOutlined } from "@ant-design/icons";
import { Button, Card, Input, Popconfirm, Space, Table, Tag } from "antd";
import { Link } from "react-router-dom";
import { useAuth } from "@/contexts/AuthContext";
import { companyService } from "@/services/companies/company-service";
import { getApiErrorMessage } from "@/utils/api-error";
import { useFeedback } from "@/hooks/useFeedback";
import type { ApiCompany } from "@/types/company";

export function CompaniesTable() {
  const { hasPermission } = useAuth();
  const feedback = useFeedback();

  const [companies, setCompanies] = useState<ApiCompany[]>([]);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState("");
  const [deletingId, setDeletingId] = useState<number | null>(null);

  const load = () => {
    setLoading(true);
    companyService
      .list()
      .then(setCompanies)
      .catch(() => feedback.error("Erro ao carregar empresas."))
      .finally(() => setLoading(false));
  };

  useEffect(load, [feedback]);

  const filtered = useMemo(() => {
    const term = search.trim().toLowerCase();
    if (!term) return companies;
    return companies.filter(
      (c) =>
        c.legal_name.toLowerCase().includes(term) ||
        c.trade_name.toLowerCase().includes(term) ||
        (c.tax_id ?? "").includes(term)
    );
  }, [companies, search]);

  const handleDelete = async (company: ApiCompany) => {
    setDeletingId(company.id);
    try {
      await companyService.remove(company.id);
      feedback.success("Empresa excluída com sucesso.");
      load();
    } catch (error) {
      feedback.error(getApiErrorMessage(error, "Erro ao excluir empresa."));
    } finally {
      setDeletingId(null);
    }
  };

  return (
    <Card
      title={
        <span>
          Empresas
          <Tag bordered={false} style={{ marginLeft: 8 }}>
            {companies.length}
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
          {hasPermission("companies.create") && (
            <Link to="/settings/companies/create">
              <Button type="primary" icon={<PlusOutlined />}>
                Nova Empresa
              </Button>
            </Link>
          )}
        </Space>
      }
    >
      <Table<ApiCompany>
        rowKey="id"
        loading={loading}
        dataSource={filtered}
        pagination={false}
        columns={[
          { title: "Razão Social", dataIndex: "legal_name" },
          { title: "Nome Fantasia", dataIndex: "trade_name" },
          {
            title: "CNPJ",
            dataIndex: "tax_id",
            render: (taxId: string | null) => taxId || <span style={{ opacity: 0.5 }}>—</span>,
          },
          {
            title: "Tipo",
            key: "type",
            render: (_, company) =>
              company.parent_id === null ? (
                <Tag color="blue">Matriz</Tag>
              ) : (
                <Tag>Vinculada</Tag>
              ),
          },
          {
            title: "Status",
            key: "status",
            render: (_, company) =>
              company.is_active ? (
                <Tag color="success">Ativa</Tag>
              ) : (
                <Tag color="default">Inativa</Tag>
              ),
          },
          {
            title: "Ações",
            key: "actions",
            render: (_, company) => (
              <Space size={6}>
                {hasPermission("companies.edit") && (
                  <Link to={`/settings/companies/${company.id}/edit`}>
                    <Button size="small" icon={<EditOutlined />} />
                  </Link>
                )}
                {hasPermission("companies.delete") && (
                  <Popconfirm
                    title="Excluir empresa"
                    description={`Tem certeza que deseja excluir "${company.trade_name || company.legal_name}"?`}
                    okText="Excluir"
                    okButtonProps={{ danger: true, loading: deletingId === company.id }}
                    cancelText="Cancelar"
                    onConfirm={() => void handleDelete(company)}
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
