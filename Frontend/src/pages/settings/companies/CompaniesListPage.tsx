import { CompaniesTable } from "@/components/companies/CompaniesTable";
import { PageHeader } from "@/components/common/PageHeader";


export function CompaniesListPage() {
  return (
    <div>
      <PageHeader
        title="Empresas"
        subtitle="Estrutura organizacional do seu tenant — empresa Matriz e empresas vinculadas"
      />

      <CompaniesTable />
    </div>
  );
}
