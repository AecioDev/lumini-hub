import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { PageHeader } from "@/components/common/PageHeader";
import { PermissionForm } from "@/components/permissions/forms/PermissionForm";
import { permissionService } from "@/services/permissions/permission-service";
import { getApiErrorMessage } from "@/utils/api-error";
import { useFeedback } from "@/hooks/useFeedback";
import type { PermissionFormValues } from "@/schemas/permission-schema";

export function CreatePermissionPage() {
  const navigate = useNavigate();
  const feedback = useFeedback();
  const [submitting, setSubmitting] = useState(false);

  const handleSubmit = async (values: PermissionFormValues) => {
    setSubmitting(true);
    try {
      await permissionService.create(values);
      feedback.success("Permissão criada com sucesso.");
      navigate("/settings/roles?tab=permissoes");
    } catch (error) {
      feedback.error(getApiErrorMessage(error, "Erro ao criar permissão."));
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div>
      <PageHeader
        title="Nova Permissão"
        subtitle="Cadastre uma nova permissão no catálogo"
        backTo="/settings/roles?tab=permissoes"
      />
      <PermissionForm
        submitting={submitting}
        submitLabel="Criar Permissão"
        onSubmit={(values) => void handleSubmit(values)}
        onCancel={() => navigate("/settings/roles?tab=permissoes")}
      />
    </div>
  );
}
