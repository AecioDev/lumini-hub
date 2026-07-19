"use client";

import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useToast } from "@/components/ui/use-toast";
import { Loader2, Shield } from "lucide-react";
import { FormField } from "@/components/common/form-field";
import {
  CreateRoleFormData,
  createRoleFormSchema,
} from "@/services/auth/role-schema";
import { Permission } from "@/services/auth/permission-schema";
import PermissionService from "@/services/auth/permission-service";
import RoleService from "@/services/auth/role-service";
import { PermissionsPicker } from "@/components/settings/permissions/permissions-picker";
import { FormFooter } from "@/components/common/form-footer";
import { routes } from "@/config/routes";

export function CreateRoleForm() {
  const router = useRouter();
  const { toast } = useToast();

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<CreateRoleFormData>({
    resolver: zodResolver(createRoleFormSchema),
  });

  const [permissionsByModule, setPermissionsByModule] = useState<Record<string, Permission[]>>({});
  const [selectedPermissionIds, setSelectedPermissionIds] = useState<Set<number>>(new Set());
  const [isLoadingPermissions, setIsLoadingPermissions] = useState(true);

  useEffect(() => {
    PermissionService.getPermissionsByModule()
      .then(setPermissionsByModule)
      .catch((error) => {
        console.error("Erro ao carregar permissões por módulo:", error);
        toast({
          variant: "destructive",
          title: "Erro",
          description: "Não foi possível carregar as permissões disponíveis.",
        });
      })
      .finally(() => setIsLoadingPermissions(false));
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

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

  const onSubmit = async (data: CreateRoleFormData) => {
    try {
      const createdRole = await RoleService.createRole({
        ...data,
        permissionIds: Array.from(selectedPermissionIds),
      });

      toast({
        title: "Perfil criado",
        description: `O perfil "${createdRole.name}" foi criado com sucesso.`,
      });
      router.push(`${routes.settings.roles.root}?role=${createdRole.id}`);
    } catch (error: any) {
      console.error("Erro ao criar perfil:", error);
      toast({
        variant: "destructive",
        title: "Erro ao criar perfil",
        description: error.response?.data?.message || "Ocorreu um erro ao salvar o perfil.",
      });
    }
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <FormField label="Nome do Perfil *" id="role-name" errorMessage={errors.name?.message}>
          <Input id="role-name" placeholder="Ex: Vendedor, Gerente, Suporte" {...register("name")} />
        </FormField>
        <FormField label="Descrição" id="role-description" errorMessage={errors.description?.message}>
          <Input id="role-description" placeholder="Explique o propósito deste perfil" {...register("description")} />
        </FormField>
      </div>

      <div className="space-y-1 pt-4 border-t">
        <label className="text-sm font-bold flex items-center gap-1">
          <Shield className="h-4 w-4 text-primary" />
          Permissões do Sistema
        </label>
        <p className="text-xs text-muted-foreground mb-3">
          Selecione as permissões que pertencerão a este novo perfil de acesso.
        </p>
        <PermissionsPicker
          permissionsByModule={permissionsByModule}
          isLoading={isLoadingPermissions}
          selectedIds={selectedPermissionIds}
          onToggle={handleTogglePermission}
          onBulkToggle={handleBulkTogglePermissions}
        />
      </div>

      <FormFooter>
        <Button
          type="button"
          variant="outline"
          onClick={() => router.push(routes.settings.roles.root)}
          disabled={isSubmitting}
        >
          Cancelar
        </Button>
        <Button type="submit" disabled={isSubmitting}>
          {isSubmitting ? (
            <>
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              Criando Perfil...
            </>
          ) : (
            "Criar Perfil"
          )}
        </Button>
      </FormFooter>
    </form>
  );
}
