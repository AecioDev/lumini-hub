import { Typography } from "antd";
import { CompaniesTable } from "@/components/companies/CompaniesTable";

const { Title, Paragraph } = Typography;

export function CompaniesListPage() {
  return (
    <div>
      <Title level={4} style={{ margin: "0 0 4px" }}>
        Empresas
      </Title>
      <Paragraph type="secondary" style={{ marginBottom: 20 }}>
        Estrutura organizacional do seu tenant — empresa Matriz e empresas vinculadas
      </Paragraph>

      <CompaniesTable />
    </div>
  );
}
