import { useEffect, useState } from "react";
import { App as AntdApp, Result, Skeleton, Typography } from "antd";
import { useNavigate, useParams } from "react-router-dom";
import { PermissionForm } from "@/components/permissions/forms/PermissionForm";
import { permissionService } from "@/services/permissions/permission-service";
import { getApiErrorMessage } from "@/utils/api-error";
import type { ApiPermissionDetail } from "@/types/permission";
import type { PermissionFormValues } from "@/schemas/permission-schema";

const { Title, Paragraph } = Typography;

export function EditPermissionPage() {
  const { id } = useParams<{ id: string }>();
  const permissionId = Number(id);
  const navigate = useNavigate();
  const { message } = AntdApp.useApp();

  const [permission, setPermission] = useState<ApiPermissionDetail | null>(null);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);

  useEffect(() => {
    if (!id || Number.isNaN(permissionId)) return;
    let cancelled = false;
    permissionService
      .getById(permissionId)
      .then((result) => {
        if (!cancelled) setPermission(result);
      })
      .catch(() => message.error("Erro ao carregar permissão."))
      .finally(() => {
        if (!cancelled) setLoading(false);
      });
    return () => {
      cancelled = true;
    };
  }, [id, permissionId, message]);

  if (!id || Number.isNaN(permissionId)) {
    return <Result status="404" title="Permissão não encontrada" />;
  }

  const handleSubmit = async (values: PermissionFormValues) => {
    setSubmitting(true);
    try {
      await permissionService.update(permissionId, values);
      message.success("Permissão atualizada com sucesso.");
      navigate("/settings/roles?tab=permissoes");
    } catch (error) {
      message.error(getApiErrorMessage(error, "Erro ao atualizar permissão."));
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div>
      <Title level={4} style={{ margin: "0 0 4px" }}>
        Editar Permissão
      </Title>
      <Paragraph type="secondary" style={{ marginBottom: 20 }}>
        Atualize os dados da permissão
      </Paragraph>

      {loading || !permission ? (
        <Skeleton active paragraph={{ rows: 6 }} />
      ) : (
        <PermissionForm
          initialValues={{
            permission: permission.permission,
            description: permission.description,
            module: permission.module,
          }}
          submitting={submitting}
          submitLabel="Salvar Permissão"
          onSubmit={(values) => void handleSubmit(values)}
          onCancel={() => navigate("/settings/roles?tab=permissoes")}
        />
      )}
    </div>
  );
}
