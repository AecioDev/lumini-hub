import { useState } from "react";
import { App as AntdApp, Typography } from "antd";
import { useNavigate } from "react-router-dom";
import { RoleForm } from "@/components/roles/forms/RoleForm";
import { roleService } from "@/services/roles/role-service";
import { getApiErrorMessage } from "@/utils/api-error";
import type { RoleFormValues } from "@/schemas/role-schema";

const { Title, Paragraph } = Typography;

export function CreateRolePage() {
  const navigate = useNavigate();
  const { message } = AntdApp.useApp();
  const [submitting, setSubmitting] = useState(false);

  const handleSubmit = async (values: RoleFormValues) => {
    setSubmitting(true);
    try {
      const created = await roleService.create(values);
      message.success("Perfil criado com sucesso. Agora configure as permissões.");
      // Vai direto pra edição — permissões só podem ser definidas depois que
      // o perfil existe (PUT /roles/:id/permissions é uma chamada separada).
      navigate(`/settings/roles/${created.id}/edit`, { replace: true });
    } catch (error) {
      message.error(getApiErrorMessage(error, "Erro ao criar perfil."));
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div>
      <Title level={4} style={{ margin: "0 0 4px" }}>
        Novo Perfil
      </Title>
      <Paragraph type="secondary" style={{ marginBottom: 20 }}>
        Cadastre um novo perfil de acesso
      </Paragraph>
      <RoleForm
        submitting={submitting}
        submitLabel="Criar Perfil"
        onSubmit={(values) => void handleSubmit(values)}
        onCancel={() => navigate("/settings/roles")}
      />
    </div>
  );
}
