"use client";

import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useToast } from "@/components/ui/use-toast";
import { Loader2, Save } from "lucide-react";
import { FormField } from "@/components/common/form-field";
import {
  Role,
  UpdateRoleFormData,
  updateRoleFormSchema,
} from "@/services/auth/role-schema";
import RoleService from "@/services/auth/role-service";
import { FormFooter } from "@/components/common/form-footer";
import { routes } from "@/config/routes";

interface EditRoleFormProps {
  role: Role;
}

export function EditRoleForm({ role }: EditRoleFormProps) {
  const router = useRouter();
  const { toast } = useToast();

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<UpdateRoleFormData>({
    resolver: zodResolver(updateRoleFormSchema),
    defaultValues: {
      name: role.name,
      description: role.description || "",
    },
  });

  const onSubmit = async (data: UpdateRoleFormData) => {
    try {
      const currentPermissionIds = (role.permissions || []).map((p) =>
        typeof p.id === "string" ? parseInt(p.id, 10) : p.id
      );
      const updatedRole = await RoleService.updateRole(role.id, {
        ...data,
        permissionIds: currentPermissionIds,
      });

      toast({
        title: "Perfil atualizado",
        description: `O perfil "${updatedRole.name}" foi atualizado com sucesso.`,
      });
      router.push(`${routes.settings.roles.root}?role=${role.id}`);
    } catch (error: any) {
      console.error("Erro ao atualizar perfil:", error);
      toast({
        variant: "destructive",
        title: "Erro ao atualizar perfil",
        description: error.response?.data?.message || "Ocorreu um erro ao salvar as alterações do perfil.",
      });
    }
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <FormField label="Nome do Perfil *" id="edit-role-name-only" errorMessage={errors.name?.message}>
          <Input id="edit-role-name-only" placeholder="Ex: Vendedor, Gerente, Suporte" {...register("name")} />
        </FormField>
        <FormField label="Descrição" id="edit-role-description-only" errorMessage={errors.description?.message}>
          <Input id="edit-role-description-only" placeholder="Explique o propósito deste perfil" {...register("description")} />
        </FormField>
      </div>

      <p className="text-xs text-muted-foreground">
        As permissões deste perfil são gerenciadas na tela principal de Perfis e Permissões.
      </p>

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
              Salvando...
            </>
          ) : (
            <>
              <Save className="mr-2 h-4 w-4" />
              Salvar
            </>
          )}
        </Button>
      </FormFooter>
    </form>
  );
}
