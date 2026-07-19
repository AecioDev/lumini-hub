// c:\Projetos\lumini-hub\Frontend\app\(system)\settings\roles\page.tsx
"use client";

import { useCallback, useEffect, useState } from "react";
import { useSearchParams } from "next/navigation";
import Link from "next/link";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useToast } from "@/components/ui/use-toast";
import { PagePermissionGuard } from "@/components/layout/PagePermissionGuard";
import { PageHeader } from "@/components/layout/page-header";
import { Role } from "@/services/auth/role-schema";
import { Permission } from "@/services/auth/permission-schema";
import RoleService from "@/services/auth/role-service";
import PermissionService from "@/services/auth/permission-service";
import { PermissionsPicker } from "@/components/settings/permissions/permissions-picker";
import { ConfirmDeleteDialog } from "@/components/common/confirm-delete-dialog";
import { FormFooter } from "@/components/common/form-footer";
import { routes } from "@/config/routes";
import { Loader2, Pencil, Plus, Trash2 } from "lucide-react";

function toId(id: Permission["id"]): number {
  return typeof id === "string" ? parseInt(id, 10) : id;
}

export default function RolesAndPermissionsPage() {
  const { toast } = useToast();
  const searchParams = useSearchParams();
  const roleIdFromUrl = searchParams.get("role");

  const [roles, setRoles] = useState<Role[]>([]);
  const [isLoadingRoles, setIsLoadingRoles] = useState(true);
  const [selectedRoleId, setSelectedRoleId] = useState<number | null>(null);

  const [allPermissionsByModule, setAllPermissionsByModule] = useState<Record<string, Permission[]>>({});
  const [isLoadingAllPermissions, setIsLoadingAllPermissions] = useState(true);

  const [selectedPermissionIds, setSelectedPermissionIds] = useState<Set<number>>(new Set());
  const [isLoadingRolePermissions, setIsLoadingRolePermissions] = useState(false);
  const [isDirty, setIsDirty] = useState(false);
  const [isSaving, setIsSaving] = useState(false);

  const [isConfirmDeleteOpen, setIsConfirmDeleteOpen] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);

  const selectedRole = roles.find((r) => r.id === selectedRoleId) || null;

  const loadRoles = useCallback(
    async (selectId?: number) => {
      setIsLoadingRoles(true);
      try {
        const data = [...(await RoleService.getRoles())].sort((a, b) =>
          a.name.localeCompare(b.name)
        );
        setRoles(data);
        if (selectId) {
          setSelectedRoleId(selectId);
        } else if (!data.some((r) => r.id === selectedRoleId)) {
          setSelectedRoleId(data.length > 0 ? data[0].id : null);
        }
      } catch (error) {
        console.error("Erro ao carregar perfis:", error);
        toast({
          variant: "destructive",
          title: "Erro",
          description: "Não foi possível carregar a lista de perfis.",
        });
      } finally {
        setIsLoadingRoles(false);
      }
    },
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [toast]
  );

  useEffect(() => {
    loadRoles(roleIdFromUrl ? Number(roleIdFromUrl) : undefined);
    const loadPermissions = async () => {
      setIsLoadingAllPermissions(true);
      try {
        const data = await PermissionService.getPermissionsByModule();
        setAllPermissionsByModule(data);
      } catch (error) {
        console.error("Erro ao carregar permissões do sistema:", error);
        toast({
          variant: "destructive",
          title: "Erro",
          description: "Não foi possível carregar o catálogo de permissões.",
        });
      } finally {
        setIsLoadingAllPermissions(false);
      }
    };
    loadPermissions();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  useEffect(() => {
    if (!selectedRoleId) {
      setSelectedPermissionIds(new Set());
      return;
    }
    const loadRolePermissions = async () => {
      setIsLoadingRolePermissions(true);
      try {
        const roleDetails = await RoleService.getRoleById(selectedRoleId);
        const ids = new Set<number>();
        (roleDetails.permissions || []).forEach((p) => ids.add(toId(p.id)));
        setSelectedPermissionIds(ids);
        setIsDirty(false);
      } catch (error) {
        console.error("Erro ao carregar permissões do perfil:", error);
      } finally {
        setIsLoadingRolePermissions(false);
      }
    };
    loadRolePermissions();
  }, [selectedRoleId]);

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
    setIsDirty(true);
  };

  const handleBulkTogglePermissions = (ids: number[], checked: boolean) => {
    setSelectedPermissionIds((prev) => {
      const next = new Set(prev);
      ids.forEach((id) => (checked ? next.add(id) : next.delete(id)));
      return next;
    });
    setIsDirty(true);
  };

  const handleSave = async () => {
    if (!selectedRole) return;
    setIsSaving(true);
    try {
      await RoleService.updateRole(selectedRole.id, {
        name: selectedRole.name,
        description: selectedRole.description || "",
        permissionIds: Array.from(selectedPermissionIds),
      });
      toast({
        title: "Permissões salvas",
        description: `As permissões do perfil "${selectedRole.name}" foram atualizadas.`,
      });
      setIsDirty(false);
    } catch (error: any) {
      console.error("Erro ao salvar permissões do perfil:", error);
      toast({
        variant: "destructive",
        title: "Erro ao salvar",
        description: error.response?.data?.message || "Não foi possível salvar as permissões do perfil.",
      });
    } finally {
      setIsSaving(false);
    }
  };

  const handleConfirmDelete = async () => {
    if (!selectedRole) return;
    setIsDeleting(true);
    try {
      await RoleService.deleteRole(selectedRole.id);
      toast({
        title: "Perfil excluído",
        description: `O perfil "${selectedRole.name}" foi excluído com sucesso.`,
      });
      setIsConfirmDeleteOpen(false);
      setSelectedRoleId(null);
      loadRoles();
    } catch (error: any) {
      console.error("Erro ao excluir perfil:", error);
      toast({
        variant: "destructive",
        title: "Erro ao excluir perfil",
        description:
          error.response?.data?.message ||
          "Não foi possível excluir o perfil (certifique-se de que ele não está associado a nenhum usuário).",
      });
    } finally {
      setIsDeleting(false);
    }
  };

  return (
    <PagePermissionGuard
      requiredPermissions={["admin.create_permissions"]}
      accessDeniedMessage="Esta tela é restrita à equipe de desenvolvimento."
    >
      <div className="space-y-6">
        <PageHeader
          title="Perfis e Permissões"
          description="Ferramenta de desenvolvedor: cadastro de perfis e o catálogo de permissões que cada perfil concede por padrão."
        />

        <Card>
          <CardHeader className="pb-4">
            <div className="flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
              <div className="flex flex-col sm:flex-row sm:items-end gap-2 flex-1">
                <div className="grid gap-1.5 w-full sm:max-w-xs">
                  <CardTitle className="text-sm font-medium text-muted-foreground">Perfil</CardTitle>
                  <Select
                    value={selectedRoleId ? selectedRoleId.toString() : ""}
                    onValueChange={(value) => setSelectedRoleId(Number.parseInt(value))}
                    disabled={isLoadingRoles || roles.length === 0}
                  >
                    <SelectTrigger>
                      <SelectValue placeholder="Selecione um perfil" />
                    </SelectTrigger>
                    <SelectContent>
                      {roles.map((role) => (
                        <SelectItem key={role.id} value={role.id.toString()}>
                          {role.name}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>

                {selectedRole && (
                  <div className="flex gap-1.5">
                    <Button variant="outline" size="icon" title="Editar nome/descrição" asChild>
                      <Link href={routes.settings.roles.edit(selectedRole.id)}>
                        <Pencil className="h-4 w-4" />
                      </Link>
                    </Button>
                    <Button
                      variant="outline"
                      size="icon"
                      title="Excluir perfil"
                      className="text-rose-600 hover:text-rose-600 hover:bg-rose-50 dark:hover:bg-rose-950/20"
                      onClick={() => setIsConfirmDeleteOpen(true)}
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </div>
                )}
              </div>

              <Button asChild>
                <Link href={routes.settings.roles.create}>
                  <Plus className="mr-1.5 h-4 w-4" />
                  Novo Perfil
                </Link>
              </Button>
            </div>
            {selectedRole?.description && (
              <CardDescription className="pt-1">{selectedRole.description}</CardDescription>
            )}
          </CardHeader>

          <CardContent>
            {isLoadingRoles ? (
              <div className="flex flex-col items-center justify-center py-16">
                <Loader2 className="h-8 w-8 animate-spin text-primary mb-3" />
                <p className="text-sm text-muted-foreground">Carregando perfis...</p>
              </div>
            ) : roles.length === 0 ? (
              <div className="flex flex-col items-center justify-center py-16 text-center">
                <p className="text-sm text-muted-foreground mb-3">
                  Nenhum perfil cadastrado ainda.
                </p>
                <Button variant="outline" asChild>
                  <Link href={routes.settings.roles.create}>
                    <Plus className="mr-1.5 h-4 w-4" />
                    Criar o primeiro perfil
                  </Link>
                </Button>
              </div>
            ) : (
              <div className="space-y-4 pt-2">
                <PermissionsPicker
                  permissionsByModule={allPermissionsByModule}
                  isLoading={isLoadingAllPermissions || isLoadingRolePermissions}
                  selectedIds={selectedPermissionIds}
                  onToggle={handleTogglePermission}
                  onBulkToggle={handleBulkTogglePermissions}
                />
                <FormFooter>
                  <Button onClick={handleSave} disabled={!isDirty || isSaving}>
                    {isSaving ? (
                      <>
                        <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                        Salvando...
                      </>
                    ) : (
                      "Salvar Alterações"
                    )}
                  </Button>
                </FormFooter>
              </div>
            )}
          </CardContent>
        </Card>

        {selectedRole && (
          <ConfirmDeleteDialog
            open={isConfirmDeleteOpen}
            onOpenChange={setIsConfirmDeleteOpen}
            itemName={`o perfil "${selectedRole.name}"`}
            onConfirm={handleConfirmDelete}
            title="Excluir Perfil?"
            description={`Tem certeza que deseja excluir permanentemente o perfil "${selectedRole.name}"? Esta ação não pode ser desfeita e irá falhar caso haja usuários vinculados a ele.`}
            confirmButtonText={isDeleting ? "Excluindo..." : "Confirmar Exclusão"}
          />
        )}
      </div>
    </PagePermissionGuard>
  );
}
