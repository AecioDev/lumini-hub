import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { PageHeader } from "@/components/common/PageHeader";
import { CompanyForm } from "@/components/companies/forms/CompanyForm";
import { companyService } from "@/services/companies/company-service";
import { getApiErrorMessage } from "@/utils/api-error";
import { useFeedback } from "@/hooks/useFeedback";
import type { CompanyFormValues } from "@/schemas/company-schema";

export function CreateCompanyPage() {
  const navigate = useNavigate();
  const feedback = useFeedback();
  const [submitting, setSubmitting] = useState(false);

  const handleSubmit = async (values: CompanyFormValues) => {
    setSubmitting(true);
    try {
      await companyService.create({
        parent_id: values.parent_id,
        legal_name: values.legal_name,
        trade_name: values.trade_name,
        cnpj: values.cnpj ?? "",
      });
      feedback.success("Empresa criada com sucesso.");
      navigate("/settings/companies", { replace: true });
    } catch (error) {
      feedback.error(getApiErrorMessage(error, "Erro ao criar empresa."));
    } finally {
      setSubmitting(false);
    }
  };

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
        onSubmit={(values) => void handleSubmit(values)}
        onCancel={() => navigate("/settings/companies")}
      />
    </div>
  );
}
