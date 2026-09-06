import { useEffect, useMemo, useState } from "react";
import { DeleteOutlined, EditOutlined, PlusOutlined } from "@ant-design/icons";
import { Button, Card, Input, Popconfirm, Space, Table, Tag } from "antd";
import { Link } from "react-router-dom";
import { useAuth } from "@/contexts/AuthContext";
import { companyService } from "@/services/companies/company-service";
import { getApiErrorMessage } from "@/utils/api-error";
import { useFeedback } from "@/hooks/useFeedback";
import type { ApiCompany } from "@/types/company";

interface CompanyTreeNode extends ApiCompany {
  children?: CompanyTreeNode[];
}

// Monta a árvore self-referencing (Matriz -> vinculadas) pra exibição
// hierárquica na grid — o antd Table já sabe renderizar recuo/expand-
// collapse quando os itens do dataSource têm `children`. Só usada quando
// não há termo de busca (buscar "achata" a árvore, ver `dataSource` mais
// abaixo — senão um item que bate com o termo mas cujo pai não bate ficaria
// escondido).
function buildCompanyTree(companies: ApiCompany[]): CompanyTreeNode[] {
  const byParent = new Map<number | null, ApiCompany[]>();
  for (const company of companies) {
    const key = company.parent_id;
    if (!byParent.has(key)) byParent.set(key, []);
    byParent.get(key)!.push(company);
  }

  const build = (parentId: number | null): CompanyTreeNode[] =>
    (byParent.get(parentId) ?? []).map((company) => {
      const children = build(company.id);
      return children.length > 0 ? { ...company, children } : { ...company };
    });

  return build(null);
}

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
        (c.cnpj ?? "").includes(term)
    );
  }, [companies, search]);

  // Sem busca: árvore hierárquica (Matriz -> vinculadas). Com busca: lista
  // achatada filtrada — senão um resultado cujo pai não bate com o termo
  // ficaria escondido dentro de uma árvore que não teria motivo pra expandir.
  const treeData = useMemo(() => buildCompanyTree(companies), [companies]);
  const dataSource = search.trim() ? filtered : treeData;

  // `defaultExpandAllRows` só decide o estado inicial contra o dataSource
  // que existe na montagem (vazio, ainda carregando) e nunca reconsidera
  // quando os dados chegam — por isso controlado aqui, recalculado sempre
  // que a lista mudar, sempre com toda hierarquia visível.
  const expandedRowKeys = useMemo(
    () => [...new Set(companies.map((c) => c.parent_id).filter((id): id is number => id !== null))],
    [companies]
  );

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
      <Table<CompanyTreeNode>
        rowKey="id"
        loading={loading}
        dataSource={dataSource}
        pagination={false}
        expandable={{ expandedRowKeys }}
        columns={[
          { title: "ID", dataIndex: "id", width: 70 },
          { title: "Razão Social", dataIndex: "legal_name" },
          { title: "Nome Fantasia", dataIndex: "trade_name" },
          {
            title: "CNPJ",
            dataIndex: "cnpj",
            render: (cnpj: string | null) => cnpj || <span style={{ opacity: 0.5 }}>—</span>,
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
