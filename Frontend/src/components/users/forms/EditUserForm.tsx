import { useEffect, useState } from "react";
import { zodResolver } from "@hookform/resolvers/zod";
import { Controller, useForm, type FieldErrors } from "react-hook-form";
import { Button, Card, Col, Input, Row, Skeleton, Switch, Tabs, Typography } from "antd";
import { useNavigate } from "react-router-dom";
import { FormField } from "@/components/common/FormField";
import { AccessDeniedResult } from "@/components/common/AccessDeniedResult";
import { ModulePermissionsPanel, pickModuleForRole } from "./ModulePermissionsPanel";
import { RolePresetPicker } from "./RolePresetPicker";
import { updateUserSchema, type UpdateUserFormValues } from "@/schemas/update-user-schema";
import { roleService } from "@/services/roles/role-service";
import { permissionService } from "@/services/permissions/permission-service";
import { userService } from "@/services/users/user-service";
import { useFeedback } from "@/hooks/useFeedback";
import { isForbiddenError } from "@/utils/api-error";
import type { ApiRole } from "@/types/role";
import type { ApiPermissionsByModule } from "@/types/permission";

const { Text } = Typography;

const DADOS_FIELDS: (keyof UpdateUserFormValues)[] = ["name", "email", "phone"];

function setsEqual(a: Set<number>, b: Set<number>): boolean {
  if (a.size !== b.size) return false;
  for (const value of a) if (!b.has(value)) return false;
  return true;
}

export function EditUserForm({ userId }: { userId: number }) {
  const navigate = useNavigate();
  const feedback = useFeedback();

  const [roles, setRoles] = useState<ApiRole[]>([]);
  const [permissionsByModule, setPermissionsByModule] = useState<
    ApiPermissionsByModule[]
  >([]);
  const [loading, setLoading] = useState(true);
  const [forbidden, setForbidden] = useState(false);
  const [username, setUsername] = useState("");

  const [selectedRoleId, setSelectedRoleId] = useState<number | null>(null);
  const [selectedPermissionIds, setSelectedPermissionIds] = useState<Set<number>>(
    new Set()
  );
  const [isCustom, setIsCustom] = useState(false);
  const [activeModule, setActiveModule] = useState<string | null>(null);
  const [activeTab, setActiveTab] = useState<"dados" | "permissoes">("dados");
  const [submitting, setSubmitting] = useState(false);

  const {
    control,
    handleSubmit,
    setValue,
    formState: { errors },
  } = useForm<UpdateUserFormValues>({
    resolver: zodResolver(updateUserSchema),
    defaultValues: { name: "", email: "", phone: "", role_id: 0, is_active: true },
  });

  useEffect(() => {
    let cancelled = false;

    Promise.all([
      userService.getById(userId),
      roleService.list(),
      permissionService.byModule(),
    ])
      .then(async ([user, rolesResult, permsResult]) => {
        if (cancelled) return;

        setUsername(user.username);
        setValue("name", user.name);
        setValue("email", user.email ?? "");
        setValue("phone", user.phone ?? "");
        setValue("role_id", user.role_id, { shouldValidate: true });
        setValue("is_active", user.is_active);

        setRoles(rolesResult);
        setPermissionsByModule(permsResult);
        setSelectedRoleId(user.role_id);

        const userPermissionIds = new Set((user.permissions ?? []).map((p) => p.id));
        setSelectedPermissionIds(userPermissionIds);
        setActiveModule(pickModuleForRole(user.role.name, permsResult, userPermissionIds));

        const roleDetail = await roleService.getById(user.role_id);
        const rolePermissionIds = new Set(roleDetail.permissions.map((p) => p.id));
        if (!cancelled) {
          setIsCustom(!setsEqual(userPermissionIds, rolePermissionIds));
        }
      })
      .catch((error) => {
        if (cancelled) return;
        if (isForbiddenError(error)) {
          setForbidden(true);
        } else {
          feedback.error("Erro ao carregar o usuário.");
        }
      })
      .finally(() => {
        if (!cancelled) setLoading(false);
      });

    return () => {
      cancelled = true;
    };
  }, [userId, setValue, feedback]);

  const handleSelectRole = async (roleId: number) => {
    setSelectedRoleId(roleId);
    setValue("role_id", roleId, { shouldValidate: true });
    try {
      const detail = await roleService.getById(roleId);
      const ids = new Set(detail.permissions.map((p) => p.id));
      setSelectedPermissionIds(ids);
      setActiveModule(pickModuleForRole(detail.name, permissionsByModule, ids));
      setIsCustom(false);
    } catch {
      feedback.error("Erro ao carregar as permissões do perfil selecionado.");
    }
  };

  const handlePermissionsChange = (next: Set<number>) => {
    setSelectedPermissionIds(next);
    setIsCustom(true);
  };

  const selectedRoleName = roles.find((r) => r.id === selectedRoleId)?.name ?? "";
  const perfilLabel = isCustom ? "Personalizado" : selectedRoleName;

  const onSubmit = async (values: UpdateUserFormValues) => {
    setSubmitting(true);
    try {
      await userService.update(userId, values);
      if (isCustom) {
        await userService.updatePermissions(userId, {
          permission_ids: [...selectedPermissionIds],
        });
      }
      feedback.success("Usuário atualizado com sucesso.");
      navigate("/settings/users");
    } catch {
      feedback.error("Erro ao salvar usuário. Confira os dados e tente novamente.");
    } finally {
      setSubmitting(false);
    }
  };

  const onInvalid = (formErrors: FieldErrors<UpdateUserFormValues>) => {
    if (DADOS_FIELDS.some((field) => formErrors[field])) {
      setActiveTab("dados");
    } else if (formErrors.role_id) {
      setActiveTab("permissoes");
    }
  };

  if (loading) {
    return <Skeleton active paragraph={{ rows: 10 }} />;
  }

  if (forbidden) {
    return (
      <AccessDeniedResult subTitle="Você não tem permissão para editar usuários ou visualizar perfis/permissões." />
    );
  }

  return (
    <form onSubmit={handleSubmit(onSubmit, onInvalid)} noValidate>
      <Tabs
        activeKey={activeTab}
        onChange={(key) => setActiveTab(key as "dados" | "permissoes")}
        items={[
          {
            key: "dados",
            label: "Dados do Usuário",
            children: (
              <Card>
                <Row gutter={16}>
                  <Col xs={24} sm={12}>
                    <FormField label="Usuário">
                      <Input value={username} disabled />
                    </FormField>
                  </Col>
                  <Col xs={24} sm={12}>
                    <FormField label="Nome completo" error={errors.name}>
                      <Controller
                        name="name"
                        control={control}
                        render={({ field }) => (
                          <Input {...field} placeholder="Ex: Carla Souza" />
                        )}
                      />
                    </FormField>
                  </Col>
                  <Col xs={24} sm={12}>
                    <FormField label="E-mail" error={errors.email}>
                      <Controller
                        name="email"
                        control={control}
                        render={({ field }) => (
                          <Input {...field} placeholder="carla@empresa.com" />
                        )}
                      />
                    </FormField>
                  </Col>
                  <Col xs={24} sm={12}>
                    <FormField label="Telefone" error={errors.phone}>
                      <Controller
                        name="phone"
                        control={control}
                        render={({ field }) => (
                          <Input {...field} placeholder="(00) 00000-0000" />
                        )}
                      />
                    </FormField>
                  </Col>
                  <Col xs={24} sm={12}>
                    <FormField label="Ativo">
                      <Controller
                        name="is_active"
                        control={control}
                        render={({ field: { value, onChange } }) => (
                          <Switch
                            checked={value}
                            onChange={onChange}
                            checkedChildren="Ativo"
                            unCheckedChildren="Inativo"
                          />
                        )}
                      />
                    </FormField>
                  </Col>
                </Row>
              </Card>
            ),
          },
          {
            key: "permissoes",
            label: "Permissões",
            children: (
              <Row gutter={[16, 16]}>
                <Col xs={24} sm={8} lg={6}>
                  <Card title="Perfil">
                    <RolePresetPicker
                      roles={roles}
                      selectedRoleId={selectedRoleId}
                      isCustom={isCustom}
                      onSelectRole={(id) => void handleSelectRole(id)}
                    />
                    {errors.role_id && (
                      <Text
                        type="danger"
                        style={{ fontSize: 12, marginTop: 8, display: "block" }}
                      >
                        {errors.role_id.message}
                      </Text>
                    )}
                  </Card>
                </Col>

                <Col xs={24} sm={16} lg={18}>
                  <Card title={`Permissões — ${perfilLabel}`}>
                    <ModulePermissionsPanel
                      permissionsByModule={permissionsByModule}
                      selectedIds={selectedPermissionIds}
                      onChange={handlePermissionsChange}
                      activeModule={activeModule}
                    />
                  </Card>
                </Col>
              </Row>
            ),
          },
        ]}
      />

      <div style={{ display: "flex", justifyContent: "flex-end", gap: 10, marginTop: 4 }}>
        <Button onClick={() => navigate("/settings/users")}>Cancelar</Button>
        <Button type="primary" htmlType="submit" loading={submitting}>
          Salvar Usuário
        </Button>
      </div>
    </form>
  );
}
