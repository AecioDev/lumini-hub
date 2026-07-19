"use client";

import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useToast } from "@/components/ui/use-toast";
import { Loader2, Save, ShieldCheck } from "lucide-react";
import { FormField } from "@/components/common/form-field";
import {
  UpdateUserFormInput,
  updateUserFormSchema,
  User,
} from "@/services/auth/user-schema";
import { Role } from "@/services/auth/role-schema";
import { Permission } from "@/services/auth/permission-schema";
import UserService from "@/services/auth/user-service";
import RoleService from "@/services/auth/role-service";
import PermissionService from "@/services/auth/permission-service";
import { PermissionsPicker } from "@/components/settings/permissions/permissions-picker";
import { FormFooter } from "@/components/common/form-footer";
import { routes } from "@/config/routes";

function toId(id: Permission["id"]): number {
  return typeof id === "string" ? parseInt(id, 10) : id;
}

interface EditUserFormProps {
  user: User;
  roles: Role[];
}

export function EditUserForm({ user, roles }: EditUserFormProps) {
  const router = useRouter();
  const { toast } = useToast();

  const {
    register,
    handleSubmit,
    watch,
    setValue,
    formState: { errors, isSubmitting },
  } = useForm<UpdateUserFormInput>({
    resolver: zodResolver(updateUserFormSchema),
    defaultValues: {
      name: user.name,
      email: user.email || "",
      phone: user.phone || "",
      role_id: user.role_id,
    },
  });

  const roleId = watch("role_id");
  const validRoles = roles.filter((role) => role.id && role.name);

  const [allPermissionsByModule, setAllPermissionsByModule] = useState<Record<string, Permission[]>>({});
  const [isLoadingAllPermissions, setIsLoadingAllPermissions] = useState(false);
  const [selectedPermissionIds, setSelectedPermissionIds] = useState<Set<number>>(new Set());
  const [isLoadingUserPermissions, setIsLoadingUserPermissions] = useState(false);

  const [filterRoleId, setFilterRoleId] = useState("all");
  const [filterRolePermIds, setFilterRolePermIds] = useState<Set<number>>(new Set());
  const [isLoadingFilterRole, setIsLoadingFilterRole] = useState(false);

  useEffect(() => {
    const loadAllPermissions = async () => {
      setIsLoadingAllPermissions(true);
      try {
        const data = await PermissionService.getPermissionsByModule();
        setAllPermissionsByModule(data);
      } catch (error) {
        console.error("Erro ao carregar permissões do sistema:", error);
      } finally {
        setIsLoadingAllPermissions(false);
      }
    };
    loadAllPermissions();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  // Carrega as permissões ativas atuais do usuário (tabela user_permissions)
  useEffect(() => {
    const loadUserPermissions = async () => {
      setIsLoadingUserPermissions(true);
      try {
        const detailedUser = await UserService.getUserById(user.id);
        const activeIds = new Set<number>();
        (detailedUser.permissions || []).forEach((p) => activeIds.add(toId(p.id)));
        setSelectedPermissionIds(activeIds);
      } catch (error) {
        console.error("Erro ao carregar permissões detalhadas do usuário:", error);
      } finally {
        setIsLoadingUserPermissions(false);
      }
    };
    loadUserPermissions();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [user.id]);

  useEffect(() => {
    if (filterRoleId === "all") {
      setFilterRolePermIds(new Set());
      return;
    }
    const loadFilterRolePerms = async () => {
      setIsLoadingFilterRole(true);
      try {
        const roleDetails = await RoleService.getRoleById(parseInt(filterRoleId, 10));
        const ids = new Set<number>();
        (roleDetails.permissions || []).forEach((p) => ids.add(toId(p.id)));
        setFilterRolePermIds(ids);
      } catch (error) {
        console.error("Erro ao carregar permissões do perfil filtrado:", error);
      } finally {
        setIsLoadingFilterRole(false);
      }
    };
    loadFilterRolePerms();
  }, [filterRoleId]);

  // Ao mudar a seleção de perfil, reseta as permissões para o padrão do novo perfil
  const handleRoleSelectChange = async (value: string) => {
    const newRoleId = Number.parseInt(value);
    setValue("role_id", newRoleId, { shouldValidate: true });

    setIsLoadingUserPermissions(true);
    try {
      const roleDetails = await RoleService.getRoleById(newRoleId);
      const roleIds = new Set<number>();
      (roleDetails.permissions || []).forEach((p) => roleIds.add(toId(p.id)));
      setSelectedPermissionIds(roleIds);
    } catch (error) {
      console.error("Erro ao mudar perfil e obter permissões:", error);
    } finally {
      setIsLoadingUserPermissions(false);
    }
  };

  const handleTogglePermission = (id: number) => {
    setSelectedPermissionIds((prev) => {
      const next = new Set(prev);
      if (next.has(id)) {
        next.delete(id);
      } else {
        next.add(id);
      }
      return next;
    });
  };

  const handleBulkTogglePermissions = (ids: number[], checked: boolean) => {
    setSelectedPermissionIds((prev) => {
      const next = new Set(prev);
      ids.forEach((id) => (checked ? next.add(id) : next.delete(id)));
      return next;
    });
  };

  const onSubmit = async (data: UpdateUserFormInput) => {
    try {
      const updatedUser = await UserService.updateUser(user.id, {
        ...data,
        is_active: user.is_active,
      });
      await UserService.updateUserPermissions(user.id, Array.from(selectedPermissionIds));

      toast({
        title: "Usuário atualizado",
        description: `O usuário ${updatedUser.name} foi atualizado com acessos personalizados.`,
      });
      router.push(routes.settings.users.root);
    } catch (error: any) {
      toast({
        variant: "destructive",
        title: "Erro ao atualizar usuário",
        description:
          error.response?.data?.message ||
          "Ocorreu um erro ao atualizar o usuário",
      });
    }
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <FormField label="Nome Completo *" id="edit-name" className="col-span-full" errorMessage={errors.name?.message}>
          <Input id="edit-name" {...register("name")} />
        </FormField>

        <FormField label="Nome de Usuário" id="edit-username">
          <Input id="edit-username" value={user.username} disabled readOnly />
        </FormField>

        <FormField label="Email *" id="edit-email" errorMessage={errors.email?.message}>
          <Input id="edit-email" type="email" placeholder="exemplo@lumini.com.br" {...register("email")} />
        </FormField>

        <FormField label="Perfil de Acesso (Modelo)" id="edit-role" className="col-span-full" errorMessage={errors.role_id?.message}>
          <Select value={roleId ? roleId.toString() : ""} onValueChange={handleRoleSelectChange}>
            <SelectTrigger id="edit-role" className="w-full">
              <SelectValue placeholder="Selecione um perfil" />
            </SelectTrigger>
            <SelectContent>
              {validRoles.length > 0 ? (
                validRoles.map((role) => (
                  <SelectItem key={role.id} value={role.id.toString()}>
                    {role.name}
                  </SelectItem>
                ))
              ) : (
                <div className="p-2 text-xs text-muted-foreground text-center">
                  Nenhum perfil disponível.
                </div>
              )}
            </SelectContent>
          </Select>
        </FormField>
      </div>

      {!!roleId && (
        <div className="space-y-1 pt-4 border-t">
          <label className="text-sm font-bold flex items-center gap-1">
            <ShieldCheck className="h-4 w-4 text-primary" />
            Personalização de Permissões de Acesso
          </label>
          <p className="text-xs text-muted-foreground mb-3">
            As permissões padrão do perfil estão marcadas. Adicione permissões extras ou revogue acessos.
          </p>
          <PermissionsPicker
            permissionsByModule={allPermissionsByModule}
            isLoading={isLoadingUserPermissions || isLoadingAllPermissions || isLoadingFilterRole}
            selectedIds={selectedPermissionIds}
            onToggle={handleTogglePermission}
            onBulkToggle={handleBulkTogglePermissions}
            selectedLabel="ativas"
            isExtraFilterActive={filterRoleId !== "all"}
            onClearExtraFilter={() => setFilterRoleId("all")}
            extraFilterPredicate={
              filterRoleId === "all"
                ? undefined
                : (p) => filterRolePermIds.has(toId(p.id))
            }
            extraFilterControl={
              <Select value={filterRoleId} onValueChange={setFilterRoleId}>
                <SelectTrigger className="h-9 text-xs">
                  <SelectValue placeholder="Filtrar por Perfil" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">Todos Perfis</SelectItem>
                  {validRoles.map((role) => (
                    <SelectItem key={role.id} value={role.id.toString()} className="text-xs">
                      {role.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            }
          />
        </div>
      )}

      <FormFooter>
        <Button
          type="button"
          variant="outline"
          onClick={() => router.push(routes.settings.users.root)}
          disabled={isSubmitting}
        >
          Cancelar
        </Button>
        <Button type="submit" disabled={isSubmitting}>
          {isSubmitting ? (
            <>
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              Salvando...
            </>
          ) : (
            <>
              <Save className="mr-2 h-4 w-4" />
              Salvar Alterações
            </>
          )}
        </Button>
      </FormFooter>
    </form>
  );
}
