import { useCallback, useState } from "react";
import { companyService } from "@/services/companies/company-service";
import { getApiErrorMessage } from "@/utils/api-error";
import { useFeedback } from "@/hooks/useFeedback";
import type { CompanyFormValues } from "@/schemas/company-schema";

// Envio do cadastro de Empresa compartilhado entre a CreateCompanyPage
// (navega pra lista ao salvar) e o CompanySelectionGate (primeiro acesso,
// sem nenhuma empresa — recarrega o usuário e segue o login). O destino pós
// sucesso é de quem chama, via `onSuccess`; aqui fica só o que é igual nos
// dois: chamada ao serviço, mapeamento dos campos, feedback e `submitting`.
export function useCreateCompany(onSuccess: () => void | Promise<void>) {
  const feedback = useFeedback();
  const [submitting, setSubmitting] = useState(false);

  const submit = useCallback(
    async (values: CompanyFormValues) => {
      setSubmitting(true);
      try {
        await companyService.create({
          parent_id: values.parent_id,
          legal_name: values.legal_name,
          trade_name: values.trade_name,
          cnpj: values.cnpj ?? "",
        });
        feedback.success("Empresa criada com sucesso.");
        await onSuccess();
      } catch (error) {
        feedback.error(getApiErrorMessage(error, "Erro ao criar empresa."));
      } finally {
        setSubmitting(false);
      }
    },
    [feedback, onSuccess]
  );

  return { submit, submitting };
}
