import { useState } from "react";
import { App as AntdApp, Typography } from "antd";
import { useNavigate } from "react-router-dom";
import { PermissionForm } from "@/components/permissions/forms/PermissionForm";
import { permissionService } from "@/services/permissions/permission-service";
import { getApiErrorMessage } from "@/utils/api-error";
import type { PermissionFormValues } from "@/schemas/permission-schema";

const { Title, Paragraph } = Typography;

export function CreatePermissionPage() {
  const navigate = useNavigate();
  const { message } = AntdApp.useApp();
  const [submitting, setSubmitting] = useState(false);

  const handleSubmit = async (values: PermissionFormValues) => {
    setSubmitting(true);
    try {
      await permissionService.create(values);
      message.success("Permissão criada com sucesso.");
      navigate("/settings/roles?tab=permissoes");
    } catch (error) {
      message.error(getApiErrorMessage(error, "Erro ao criar permissão."));
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div>
      <Title level={4} style={{ margin: "0 0 4px" }}>
        Nova Permissão
      </Title>
      <Paragraph type="secondary" style={{ marginBottom: 20 }}>
        Cadastre uma nova permissão no catálogo
      </Paragraph>
      <PermissionForm
        submitting={submitting}
        submitLabel="Criar Permissão"
        onSubmit={(values) => void handleSubmit(values)}
        onCancel={() => navigate("/settings/roles?tab=permissoes")}
      />
    </div>
  );
}
