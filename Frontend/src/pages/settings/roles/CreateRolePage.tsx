import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { PageHeader } from "@/components/common/PageHeader";
import { RoleForm } from "@/components/roles/forms/RoleForm";
import { roleService } from "@/services/roles/role-service";
import { getApiErrorMessage } from "@/utils/api-error";
import { useFeedback } from "@/hooks/useFeedback";
import type { RoleFormValues } from "@/schemas/role-schema";

export function CreateRolePage() {
  const navigate = useNavigate();
  const feedback = useFeedback();
  const [submitting, setSubmitting] = useState(false);

  const handleSubmit = async (values: RoleFormValues) => {
    setSubmitting(true);
    try {
      const created = await roleService.create(values);
      feedback.success("Perfil criado com sucesso. Agora configure as permissões.");
      // Vai direto pra edição — permissões só podem ser definidas depois que
      // o perfil existe (PUT /roles/:id/permissions é uma chamada separada).
      navigate(`/settings/roles/${created.id}/edit`, { replace: true });
    } catch (error) {
      feedback.error(getApiErrorMessage(error, "Erro ao criar perfil."));
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div>
      <PageHeader
        title="Novo Perfil"
        subtitle="Cadastre um novo perfil de acesso"
        backTo="/settings/roles"
      />
      <RoleForm
        submitting={submitting}
        submitLabel="Criar Perfil"
        onSubmit={(values) => void handleSubmit(values)}
        onCancel={() => navigate("/settings/roles")}
      />
    </div>
  );
}
