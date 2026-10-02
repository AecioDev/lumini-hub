import { useNavigate } from "react-router-dom";
import { PageHeader } from "@/components/common/PageHeader";
import { CompanyForm } from "@/components/companies/forms/CompanyForm";
import { useCreateCompany } from "@/hooks/useCreateCompany";

export function CreateCompanyPage() {
  const navigate = useNavigate();
  const { submit, submitting } = useCreateCompany(() =>
    navigate("/settings/companies", { replace: true })
  );

  return (
    <div>
      <PageHeader
        title="Nova Empresa"
        subtitle="Cadastre a empresa Matriz ou uma empresa vinculada"
        backTo="/settings/companies"
      />
      <CompanyForm
        submitting={submitting}
        submitLabel="Criar Empresa"
        onSubmit={(values) => void submit(values)}
        onCancel={() => navigate("/settings/companies")}
      />
    </div>
  );
}
