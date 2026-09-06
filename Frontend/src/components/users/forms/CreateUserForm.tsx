import { useEffect, useMemo, useState } from "react";
import { zodResolver } from "@hookform/resolvers/zod";
import { Controller, useForm, type FieldErrors } from "react-hook-form";
import { Button, Card, Col, Input, Row, Select, Skeleton, Tabs } from "antd";
import { useNavigate } from "react-router-dom";
import { FormField } from "@/components/common/FormField";
import { AccessDeniedResult } from "@/components/common/AccessDeniedResult";
import {
  ModulePermissionsPanel,
  pickModuleForRole,
  getOrphanModules,
  OTHERS_FILTER_ID,
  OTHERS_MODULE_KEY,
} from "./ModulePermissionsPanel";
import { RolePresetPicker } from "./RolePresetPicker";
import { UserPermissionsSummary } from "./UserPermissionsSummary";
import { createUserSchema, type CreateUserFormValues } from "@/schemas/create-user-schema";
import { roleService } from "@/services/roles/role-service";
import { permissionService } from "@/services/permissions/permission-service";
import { userService } from "@/services/users/user-service";
import { companyService } from "@/services/companies/company-service";
import { useFeedback } from "@/hooks/useFeedback";
import { isForbiddenError } from "@/utils/api-error";
import type { ApiRole } from "@/types/role";
import type { ApiPermissionsByModule } from "@/types/permission";
import type { ApiCompany } from "@/types/company";

// Opção explícita no Select "Empresa" pra virar "master" (sem empresa
// vinculada) — 0 nunca é um ID de Company de verdade (seriais do banco
// começam em 1). Antes disso existir só dava pra "limpar" via allowClear,
// que era pouco descoberto e também expunha um bug real: o backend
// ignorava `company_id: null` no update (corrigido em UserService.
// UpdateUser, achado testando).
const MASTER_OPTION_VALUE = 0;

const DADOS_FIELDS: (keyof CreateUserFormValues)[] = [
  "username",
  "password",
  "name",
  "email",
  "phone",
  "role_id",
  "company_id",
];

export function CreateUserForm() {
  const navigate = useNavigate();
  const feedback = useFeedback();

  const [roles, setRoles] = useState<ApiRole[]>([]);
  const [permissionsByModule, setPermissionsByModule] = useState<
    ApiPermissionsByModule[]
  >([]);
  const [companies, setCompanies] = useState<ApiCompany[]>([]);
  const [loadingPresets, setLoadingPresets] = useState(true);
  const [forbidden, setForbidden] = useState(false);

  // Perfil atribuído ao usuário (campo "Perfil" na aba Dados do Usuário).
  const [roleId, setRoleId] = useState<number | null>(null);
  // Perfil usado só como filtro para pré-carregar a matriz de permissões
  // (aba Permissões) — independente do Perfil atribuído acima.
  const [filterRoleId, setFilterRoleId] = useState<number | null>(null);
  const [selectedPermissionIds, setSelectedPermissionIds] = useState<Set<number>>(
    new Set()
  );
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
      company_id: null,
    },
  });

  useEffect(() => {
    let cancelled = false;

    Promise.all([roleService.list(), permissionService.byModule(), companyService.list()])
      .then(async ([rolesResult, permsResult, companiesResult]) => {
        if (cancelled) return;
        setRoles(rolesResult);
        setPermissionsByModule(permsResult);
        setCompanies(companiesResult);

        if (rolesResult.length > 0) {
          const first = rolesResult[0];
          setRoleId(first.id);
          setFilterRoleId(first.id);
          setValue("role_id", first.id, { shouldValidate: true });
          const detail = await roleService.getById(first.id);
          if (!cancelled) {
            const ids = new Set(detail.permissions.map((p) => p.id));
            setSelectedPermissionIds(ids);
            setActiveModule(pickModuleForRole(detail.name, permsResult, ids));
          }
        }
      })
      .catch((error) => {
        if (cancelled) return;
        if (isForbiddenError(error)) {
          setForbidden(true);
        } else {
          feedback.error("Erro ao carregar perfis/permissões.");
        }
      })
      .finally(() => {
        if (!cancelled) setLoadingPresets(false);
      });

    return () => {
      cancelled = true;
    };
  }, [setValue, feedback]);

  // Select "Perfil" na aba Dados do Usuário — muda o Perfil atribuído ao
  // usuário e pré-preenche a matriz de permissões com o padrão desse perfil.
  const handleChangeUserRole = async (newRoleId: number) => {
    setRoleId(newRoleId);
    setValue("role_id", newRoleId, { shouldValidate: true });
    try {
      const detail = await roleService.getById(newRoleId);
      const ids = new Set(detail.permissions.map((p) => p.id));
      setSelectedPermissionIds(ids);
      setFilterRoleId(newRoleId);
      setActiveModule(pickModuleForRole(detail.name, permissionsByModule, ids));
    } catch {
      feedback.error("Erro ao carregar as permissões do perfil selecionado.");
    }
  };

  // Botões de Perfil na aba Permissões — só navegam até o módulo de
  // permissões daquele perfil para visualização/edição. Não mexem no Perfil
  // atribuído ao usuário (role_id) nem substituem a seleção de permissões já
  // marcada (que pode acumular permissões de vários módulos diferentes).
  const handleFilterRole = (newFilterRoleId: number) => {
    setFilterRoleId(newFilterRoleId);
    if (newFilterRoleId === OTHERS_FILTER_ID) {
      setActiveModule(OTHERS_MODULE_KEY);
      return;
    }
    const role = roles.find((r) => r.id === newFilterRoleId);
    if (role) {
      setActiveModule(pickModuleForRole(role.name, permissionsByModule, selectedPermissionIds));
    }
  };

  const handlePermissionsChange = (next: Set<number>) => {
    setSelectedPermissionIds(next);
  };

  const filterRoleName =
    filterRoleId === OTHERS_FILTER_ID
      ? "Outros"
      : roles.find((r) => r.id === filterRoleId)?.name ?? "";

  // Módulos do catálogo sem Perfil correspondente (ex.: "Empresas") —
  // inalcançáveis pelo filtro normal de Perfil, ver getOrphanModules.
  const orphanModules = useMemo(
    () => getOrphanModules(roles, permissionsByModule),
    [roles, permissionsByModule]
  );
  const othersHasPermissions = useMemo(
    () =>
      orphanModules.some((entry) =>
        entry.permissions.some((p) => selectedPermissionIds.has(p.id))
      ),
    [orphanModules, selectedPermissionIds]
  );

  // Perfis cujo módulo tem alguma permissão marcada — usado pra destacar,
  // no filtro de perfil, quais perfis "contribuíram" com alguma permissão
  // além do que veio do Perfil atribuído ao usuário.
  const rolesWithPermissions = useMemo(() => {
    const ids = new Set<number>();
    for (const role of roles) {
      const entry = permissionsByModule.find(
        (m) => m.module.trim().toLowerCase() === role.name.trim().toLowerCase()
      );
      if (entry?.permissions.some((p) => selectedPermissionIds.has(p.id))) {
        ids.add(role.id);
      }
    }
    return ids;
  }, [roles, permissionsByModule, selectedPermissionIds]);

  const onSubmit = async (values: CreateUserFormValues) => {
    setSubmitting(true);
    try {
      const created = await userService.create(values);
      await userService.updatePermissions(created.id, {
        permission_ids: [...selectedPermissionIds],
      });
      feedback.success("Usuário criado com sucesso.");
      navigate("/settings/users");
    } catch {
      feedback.error("Erro ao criar usuário. Confira os dados e tente novamente.");
    } finally {
      setSubmitting(false);
    }
  };

  // Se o usuário tentar salvar com erro num campo que está numa aba
  // escondida, leva ele até lá em vez de deixar o erro invisível.
  const onInvalid = (formErrors: FieldErrors<CreateUserFormValues>) => {
    if (DADOS_FIELDS.some((field) => formErrors[field])) {
      setActiveTab("dados");
    }
  };

  if (forbidden) {
    return (
      <AccessDeniedResult subTitle="Você não tem permissão para visualizar perfis/permissões, necessários para cadastrar um usuário." />
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
                  <Col xs={24} sm={12}>
                    <FormField label="Perfil" error={errors.role_id}>
                      {loadingPresets ? (
                        <Skeleton.Input active block style={{ width: "100%" }} />
                      ) : (
                        <Select
                          value={roleId ?? undefined}
                          placeholder="Selecione o perfil do usuário"
                          options={roles.map((role) => ({
                            value: role.id,
                            label: role.name,
                          }))}
                          onChange={(id) => void handleChangeUserRole(id)}
                        />
                      )}
                    </FormField>
                  </Col>
                  <Col xs={24} sm={12}>
                    <FormField label="Empresa" error={errors.company_id}>
                      {loadingPresets ? (
                        <Skeleton.Input active block style={{ width: "100%" }} />
                      ) : (
                        <Controller
                          name="company_id"
                          control={control}
                          render={({ field }) => (
                            <Select
                              value={field.value ?? MASTER_OPTION_VALUE}
                              options={[
                                {
                                  value: MASTER_OPTION_VALUE,
                                  label: "Nenhuma (usuário master, vê todas as empresas)",
                                },
                                ...companies.map((company) => ({
                                  value: company.id,
                                  label: company.trade_name || company.legal_name,
                                })),
                              ]}
                              onChange={(id) =>
                                field.onChange(id === MASTER_OPTION_VALUE ? null : id)
                              }
                            />
                          )}
                        />
                      )}
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
              <>
                {!loadingPresets && (
                  <UserPermissionsSummary
                    permissionsByModule={permissionsByModule}
                    selectedIds={selectedPermissionIds}
                  />
                )}
                <Row gutter={[16, 16]}>
                  <Col xs={24} sm={8} lg={6}>
                    <Card title="Filtrar perfil">
                      {loadingPresets ? (
                        <Skeleton active paragraph={{ rows: 3 }} />
                      ) : (
                        <RolePresetPicker
                          roles={roles}
                          selectedRoleId={filterRoleId}
                          onSelectRole={handleFilterRole}
                          rolesWithPermissions={rolesWithPermissions}
                          showOthers={orphanModules.length > 0}
                          othersHasPermissions={othersHasPermissions}
                        />
                      )}
                    </Card>
                  </Col>

                  <Col xs={24} sm={16} lg={18}>
                    <Card title={`Permissões — ${filterRoleName}`}>
                      {loadingPresets ? (
                        <Skeleton active paragraph={{ rows: 5 }} />
                      ) : (
                        <ModulePermissionsPanel
                          permissionsByModule={permissionsByModule}
                          selectedIds={selectedPermissionIds}
                          onChange={handlePermissionsChange}
                          activeModule={activeModule}
                          orphanModules={orphanModules}
                        />
                      )}
                    </Card>
                  </Col>
                </Row>
              </>
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
