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
import { Loader2, ShieldCheck, UserPlus } from "lucide-react";
import { FormField } from "@/components/common/form-field";
import {
  CreateUserFormInput,
  createUserFormSchema,
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

interface CreateUserFormProps {
  roles: Role[];
}

export function CreateUserForm({ roles }: CreateUserFormProps) {
  const router = useRouter();
  const { toast } = useToast();

  const {
    register,
    handleSubmit,
    watch,
    setValue,
    formState: { errors, isSubmitting },
  } = useForm<CreateUserFormInput>({
    resolver: zodResolver(createUserFormSchema),
  });

  const roleId = watch("role_id");
  const validRoles = roles.filter((role) => role.id && role.name);

  const [allPermissionsByModule, setAllPermissionsByModule] = useState<Record<string, Permission[]>>({});
  const [isLoadingAllPermissions, setIsLoadingAllPermissions] = useState(false);
  const [selectedPermissionIds, setSelectedPermissionIds] = useState<Set<number>>(new Set());
  const [isLoadingRolePermissions, setIsLoadingRolePermissions] = useState(false);

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
        toast({
          variant: "destructive",
          title: "Erro",
          description: "Não foi possível carregar as permissões do sistema.",
        });
      } finally {
        setIsLoadingAllPermissions(false);
      }
    };
    loadAllPermissions();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  useEffect(() => {
    const loadRolePermissions = async () => {
      if (!roleId) {
        setSelectedPermissionIds(new Set());
        return;
      }
      setIsLoadingRolePermissions(true);
      try {
        const roleDetails = await RoleService.getRoleById(roleId);
        const ids = new Set<number>();
        (roleDetails.permissions || []).forEach((p) => ids.add(toId(p.id)));
        setSelectedPermissionIds(ids);
      } catch (error) {
        console.error("Erro ao carregar permissões do perfil:", error);
      } finally {
        setIsLoadingRolePermissions(false);
      }
    };
    loadRolePermissions();
  }, [roleId]);

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

  const onSubmit = async (data: CreateUserFormInput) => {
    try {
      const createdUser = await UserService.createUser(data);
      await UserService.updateUserPermissions(createdUser.id, Array.from(selectedPermissionIds));

      toast({
        title: "Usuário adicionado",
        description: `O usuário ${createdUser.name} foi criado com acessos personalizados.`,
      });
      router.push(routes.settings.users.root);
    } catch (error: any) {
      toast({
        variant: "destructive",
        title: "Erro ao adicionar usuário",
        description:
          error.response?.data?.message ||
          "Ocorreu um erro ao adicionar o usuário",
      });
    }
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <FormField label="Nome Completo *" id="name" className="col-span-full" errorMessage={errors.name?.message}>
          <Input id="name" placeholder="Ex: João Silva" {...register("name")} />
        </FormField>

        <FormField label="Nome de Usuário *" id="username" errorMessage={errors.username?.message}>
          <Input id="username" placeholder="joao_silva" {...register("username")} />
        </FormField>

        <FormField label="Senha *" id="password" errorMessage={errors.password?.message}>
          <Input id="password" type="password" placeholder="Mínimo 6 dígitos" {...register("password")} />
        </FormField>

        <FormField label="Email *" id="email" errorMessage={errors.email?.message}>
          <Input id="email" type="email" placeholder="exemplo@lumini.com.br" {...register("email")} />
        </FormField>

        <FormField label="Perfil de Acesso (Modelo) *" id="role" errorMessage={errors.role_id?.message}>
          <Select
            value={roleId ? roleId.toString() : ""}
            onValueChange={(value) => setValue("role_id", Number.parseInt(value), { shouldValidate: true })}
          >
            <SelectTrigger id="role" className="w-full">
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
            isLoading={isLoadingRolePermissions || isLoadingAllPermissions || isLoadingFilterRole}
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
              Criando...
            </>
          ) : (
            <>
              <UserPlus className="mr-2 h-4 w-4" />
              Criar Usuário
            </>
          )}
        </Button>
      </FormFooter>
    </form>
  );
}
