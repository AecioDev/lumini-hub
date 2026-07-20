import { useEffect, useState } from "react";
import { zodResolver } from "@hookform/resolvers/zod";
import { Controller, useForm, type FieldErrors } from "react-hook-form";
import { App as AntdApp, Button, Card, Col, Input, Row, Skeleton, Tabs, Typography } from "antd";
import { useNavigate } from "react-router-dom";
import { FormField } from "@/components/common/FormField";
import { ModulePermissionsPanel, pickModuleForRole } from "./ModulePermissionsPanel";
import { RolePresetPicker } from "./RolePresetPicker";
import { createUserSchema, type CreateUserFormValues } from "@/schemas/create-user-schema";
import { roleService } from "@/services/roles/role-service";
import { permissionService } from "@/services/permissions/permission-service";
import { userService } from "@/services/users/user-service";
import type { ApiRole } from "@/types/role";
import type { ApiPermissionsByModule } from "@/types/permission";

const { Text } = Typography;

const DADOS_FIELDS: (keyof CreateUserFormValues)[] = [
  "username",
  "password",
  "name",
  "email",
  "phone",
];

export function CreateUserForm() {
  const navigate = useNavigate();
  const { message } = AntdApp.useApp();

  const [roles, setRoles] = useState<ApiRole[]>([]);
  const [permissionsByModule, setPermissionsByModule] = useState<
    ApiPermissionsByModule[]
  >([]);
  const [loadingPresets, setLoadingPresets] = useState(true);

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
  } = useForm<CreateUserFormValues>({
    resolver: zodResolver(createUserSchema),
    defaultValues: {
      username: "",
      password: "",
      name: "",
      email: "",
      phone: "",
      role_id: 0,
    },
  });

  useEffect(() => {
    let cancelled = false;

    Promise.all([roleService.list(), permissionService.byModule()])
      .then(async ([rolesResult, permsResult]) => {
        if (cancelled) return;
        setRoles(rolesResult);
        setPermissionsByModule(permsResult);

        if (rolesResult.length > 0) {
          const first = rolesResult[0];
          setSelectedRoleId(first.id);
          setValue("role_id", first.id, { shouldValidate: true });
          const detail = await roleService.getById(first.id);
          if (!cancelled) {
            const ids = new Set(detail.permissions.map((p) => p.id));
            setSelectedPermissionIds(ids);
            setActiveModule(pickModuleForRole(detail.name, permsResult, ids));
          }
        }
      })
      .catch(() => message.error("Erro ao carregar perfis/permissões."))
      .finally(() => {
        if (!cancelled) setLoadingPresets(false);
      });

    return () => {
      cancelled = true;
    };
  }, [setValue, message]);

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
      message.error("Erro ao carregar as permissões do perfil selecionado.");
    }
  };

  const handlePermissionsChange = (next: Set<number>) => {
    setSelectedPermissionIds(next);
    setIsCustom(true);
  };

  const selectedRoleName = roles.find((r) => r.id === selectedRoleId)?.name ?? "";
  const perfilLabel = isCustom ? "Personalizado" : selectedRoleName;

  const onSubmit = async (values: CreateUserFormValues) => {
    setSubmitting(true);
    try {
      const created = await userService.create(values);
      if (isCustom) {
        await userService.updatePermissions(created.id, {
          permission_ids: [...selectedPermissionIds],
        });
      }
      message.success("Usuário criado com sucesso.");
      navigate("/settings/users");
    } catch {
      message.error("Erro ao criar usuário. Confira os dados e tente novamente.");
    } finally {
      setSubmitting(false);
    }
  };

  // Se o usuário tentar salvar com erro num campo que está numa aba
  // escondida, leva ele até lá em vez de deixar o erro invisível.
  const onInvalid = (formErrors: FieldErrors<CreateUserFormValues>) => {
    if (DADOS_FIELDS.some((field) => formErrors[field])) {
      setActiveTab("dados");
    } else if (formErrors.role_id) {
      setActiveTab("permissoes");
    }
  };

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
                    <FormField label="Usuário" error={errors.username}>
                      <Controller
                        name="username"
                        control={control}
                        render={({ field }) => (
                          <Input {...field} placeholder="ex: carla.souza" />
                        )}
                      />
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
                    <FormField label="Senha temporária" error={errors.password}>
                      <Controller
                        name="password"
                        control={control}
                        render={({ field }) => (
                          <Input.Password {...field} placeholder="••••••••" />
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
                    {loadingPresets ? (
                      <Skeleton active paragraph={{ rows: 3 }} />
                    ) : (
                      <RolePresetPicker
                        roles={roles}
                        selectedRoleId={selectedRoleId}
                        isCustom={isCustom}
                        onSelectRole={(id) => void handleSelectRole(id)}
                      />
                    )}
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
                    {loadingPresets ? (
                      <Skeleton active paragraph={{ rows: 5 }} />
                    ) : (
                      <ModulePermissionsPanel
                        permissionsByModule={permissionsByModule}
                        selectedIds={selectedPermissionIds}
                        onChange={handlePermissionsChange}
                        activeModule={activeModule}
                      />
                    )}
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
