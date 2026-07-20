import { useMemo, useState } from "react";
import { Card, Input, Table, Tag } from "antd";
import type { ClientStatus, RecentClient } from "@/types/crm";

const STATUS_COLOR: Record<ClientStatus, string> = {
  Lead: "warning",
  Proposta: "cyan",
  Cliente: "success",
};

export function RecentClientsTable({ clients }: { clients: RecentClient[] }) {
  const [search, setSearch] = useState("");

  const filtered = useMemo(() => {
    const term = search.trim().toLowerCase();
    if (!term) return clients;
    return clients.filter(
      (client) =>
        client.name.toLowerCase().includes(term) ||
        client.company.toLowerCase().includes(term)
    );
  }, [clients, search]);

  return (
    <Card
      title={
        <span>
          Clientes Recentes
          <Tag bordered={false} style={{ marginLeft: 8 }}>
            {clients.length}
          </Tag>
        </span>
      }
      extra={
        <Input.Search
          placeholder="Buscar..."
          allowClear
          style={{ width: 200 }}
          onChange={(e) => setSearch(e.target.value)}
        />
      }
    >
      <Table<RecentClient>
        rowKey="id"
        dataSource={filtered}
        pagination={false}
        columns={[
          { title: "Cliente", dataIndex: "name" },
          { title: "Empresa", dataIndex: "company" },
          { title: "Responsável", dataIndex: "owner" },
          {
            title: "Status",
            dataIndex: "status",
            render: (status: ClientStatus) => (
              <Tag color={STATUS_COLOR[status]}>{status}</Tag>
            ),
          },
          { title: "Valor", dataIndex: "value" },
        ]}
      />
    </Card>
  );
}
