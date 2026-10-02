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
      } catch (error) {
        feedback.error(getApiErrorMessage(error, "Erro ao criar empresa."));
        setSubmitting(false);
        return;
      }

      // A empresa já existe a partir daqui: uma falha do `onSuccess` (ex.:
      // refreshUser do Gate) não pode virar "Erro ao criar empresa" — o
      // usuário tentaria recriar e esbarraria no CNPJ duplicado.
      feedback.success("Empresa criada com sucesso.");
      try {
        await onSuccess();
      } catch {
        feedback.error(
          "A empresa foi criada, mas não foi possível continuar. Recarregue a página.",
          "Atenção"
        );
      } finally {
        setSubmitting(false);
      }
    },
    [feedback, onSuccess]
  );

  return { submit, submitting };
}
