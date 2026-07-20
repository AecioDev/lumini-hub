import { useEffect, useState } from "react";
import { App as AntdApp, Result, Skeleton, Tabs, Typography } from "antd";
import { useNavigate, useParams } from "react-router-dom";
import { RoleForm } from "@/components/roles/forms/RoleForm";
import { RolePermissionsEditor } from "@/components/roles/forms/RolePermissionsEditor";
import { roleService } from "@/services/roles/role-service";
import { getApiErrorMessage } from "@/utils/api-error";
import type { ApiRoleDetail } from "@/types/role";
import type { RoleFormValues } from "@/schemas/role-schema";

const { Title, Paragraph } = Typography;

export function EditRolePage() {
  const { id } = useParams<{ id: string }>();
  const roleId = Number(id);
  const navigate = useNavigate();
  const { message } = AntdApp.useApp();

  const [role, setRole] = useState<ApiRoleDetail | null>(null);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);

  useEffect(() => {
    if (!id || Number.isNaN(roleId)) return;
    let cancelled = false;
    roleService
      .getById(roleId)
      .then((result) => {
        if (!cancelled) setRole(result);
      })
      .catch(() => message.error("Erro ao carregar perfil."))
      .finally(() => {
        if (!cancelled) setLoading(false);
      });
    return () => {
      cancelled = true;
    };
  }, [id, roleId, message]);

  if (!id || Number.isNaN(roleId)) {
    return <Result status="404" title="Perfil não encontrado" />;
  }

  const handleSubmit = async (values: RoleFormValues) => {
    setSubmitting(true);
    try {
      const updated = await roleService.update(roleId, values);
      setRole((prev) => (prev ? { ...prev, ...updated } : prev));
      message.success("Perfil atualizado com sucesso.");
    } catch (error) {
      message.error(getApiErrorMessage(error, "Erro ao atualizar perfil."));
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div>
      <Title level={4} style={{ margin: "0 0 4px" }}>
        Editar Perfil
      </Title>
      <Paragraph type="secondary" style={{ marginBottom: 20 }}>
        Atualize os dados e permissões do perfil
      </Paragraph>

      {loading || !role ? (
        <Skeleton active paragraph={{ rows: 10 }} />
      ) : (
        <Tabs
          items={[
            {
              key: "dados",
              label: "Dados do Perfil",
              children: (
                <RoleForm
                  initialValues={{ name: role.name, description: role.description }}
                  submitting={submitting}
                  submitLabel="Salvar Perfil"
                  onSubmit={(values) => void handleSubmit(values)}
                  onCancel={() => navigate("/settings/roles")}
                />
              ),
            },
            {
              key: "permissoes",
              label: "Permissões",
              children: (
                <RolePermissionsEditor
                  roleId={roleId}
                  initialPermissionIds={role.permissions.map((p) => p.id)}
                />
              ),
            },
          ]}
        />
      )}
    </div>
  );
}
