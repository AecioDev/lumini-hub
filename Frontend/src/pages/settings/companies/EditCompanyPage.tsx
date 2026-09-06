import { useEffect, useState } from "react";
import { Result, Skeleton, Tabs } from "antd";
import { useNavigate, useParams } from "react-router-dom";
import { PageHeader } from "@/components/common/PageHeader";
import { CompanyForm } from "@/components/companies/forms/CompanyForm";
import { CompanyFiscalConfigForm } from "@/components/companies/forms/CompanyFiscalConfigForm";
import { companyService } from "@/services/companies/company-service";
import { getApiErrorMessage } from "@/utils/api-error";
import { useFeedback } from "@/hooks/useFeedback";
import type { ApiCompanyDetail } from "@/types/company";
import type { CompanyFormValues } from "@/schemas/company-schema";

export function EditCompanyPage() {
  const { id } = useParams<{ id: string }>();
  const companyId = Number(id);
  const navigate = useNavigate();
  const feedback = useFeedback();

  const [company, setCompany] = useState<ApiCompanyDetail | null>(null);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);

  useEffect(() => {
    if (!id || Number.isNaN(companyId)) return;
    let cancelled = false;
    companyService
      .getById(companyId)
      .then((result) => {
        if (!cancelled) setCompany(result);
      })
      .catch(() => feedback.error("Erro ao carregar empresa."))
      .finally(() => {
        if (!cancelled) setLoading(false);
      });
    return () => {
      cancelled = true;
    };
  }, [id, companyId, feedback]);

  if (!id || Number.isNaN(companyId)) {
    return <Result status="404" title="Empresa não encontrada" />;
  }

  const handleSubmit = async (values: CompanyFormValues) => {
    setSubmitting(true);
    try {
      const updated = await companyService.update(companyId, {
        parent_id: values.parent_id,
        legal_name: values.legal_name,
        trade_name: values.trade_name,
        tax_id: values.tax_id ?? "",
        is_active: values.is_active,
      });
      setCompany((prev) => (prev ? { ...prev, ...updated } : prev));
      feedback.success("Empresa atualizada com sucesso.");
    } catch (error) {
      feedback.error(getApiErrorMessage(error, "Erro ao atualizar empresa."));
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div>
      <PageHeader
        title="Editar Empresa"
        subtitle="Atualize os dados cadastrais da empresa"
        backTo="/settings/companies"
      />

      {loading || !company ? (
        <Skeleton active paragraph={{ rows: 8 }} />
      ) : (
        <Tabs
          items={[
            {
              key: "dados",
              label: "Dados da Empresa",
              children: (
                <CompanyForm
                  editingId={company.id}
                  initialValues={{
                    parent_id: company.parent_id,
                    legal_name: company.legal_name,
                    trade_name: company.trade_name,
                    tax_id: company.tax_id ?? "",
                    is_active: company.is_active,
                  }}
                  submitting={submitting}
                  submitLabel="Salvar Empresa"
                  onSubmit={(values) => void handleSubmit(values)}
                  onCancel={() => navigate("/settings/companies")}
                />
              ),
            },
            {
              key: "fiscal",
              label: "Configuração Fiscal",
              children: (
                <CompanyFiscalConfigForm companyId={company.id} companyTaxId={company.tax_id ?? undefined} />
              ),
            },
          ]}
        />
      )}
    </div>
  );
}
