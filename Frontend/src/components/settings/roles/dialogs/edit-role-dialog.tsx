"use client";

import { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Checkbox } from "@/components/ui/checkbox";
import { useToast } from "@/components/ui/use-toast";
import { Loader2, ShieldAlert } from "lucide-react";
import { Permission } from "@/services/auth/permission-schema";
import { Role } from "@/services/auth/role-schema";
import PermissionService from "@/services/auth/permission-service";
import RoleService from "@/services/auth/role-service";

interface EditRoleDialogProps {
  isOpen: boolean;
  onOpenChange: (open: boolean) => void;
  role: Role | null;
  onRoleUpdated: () => void;
}

export function EditRoleDialog({
  isOpen,
  onOpenChange,
  role,
  onRoleUpdated,
}: EditRoleDialogProps) {
  const { toast } = useToast();
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [permissionsByModule, setPermissionsByModule] = useState<Record<string, Permission[]>>({});
  const [selectedPermissionIds, setSelectedPermissionIds] = useState<Set<number>>(new Set());
  const [isLoadingPermissions, setIsLoadingPermissions] = useState(false);
  const [isSaving, setIsSaving] = useState(false);

  // Carregar as permissões e dados da Role ao abrir o diálogo
  useEffect(() => {
    if (isOpen && role) {
      setName(role.name);
      setDescription(role.description || "");

      const loadRoleDataAndPermissions = async () => {
        setIsLoadingPermissions(true);
        try {
          // Carregar todas as permissões do sistema
          const allPermissions = await PermissionService.getPermissionsByModule();
          setPermissionsByModule(allPermissions);

          // Carregar os dados detalhados da Role (incluindo as permissões vinculadas)
          const detailedRole = await RoleService.getRoleById(role.id);
          const activePermissionIds = new Set<number>();
          if (detailedRole.permissions) {
            detailedRole.permissions.forEach((p) => {
              const idNum = typeof p.id === "string" ? parseInt(p.id, 10) : p.id;
              activePermissionIds.add(idNum);
            });
          }
          setSelectedPermissionIds(activePermissionIds);
        } catch (error: any) {
          console.error("Erro ao carregar dados do perfil e permissões:", error);
          toast({
            variant: "destructive",
            title: "Erro",
            description: "Não foi possível carregar as permissões do perfil.",
          });
        } finally {
          setIsLoadingPermissions(false);
        }
      };

      loadRoleDataAndPermissions();
    }
  }, [isOpen, role, toast]);

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

  const handleToggleModuleAll = (modulePermissions: Permission[], isChecked: boolean) => {
    setSelectedPermissionIds((prev) => {
      const next = new Set(prev);
      modulePermissions.forEach((p) => {
        const idNum = typeof p.id === "string" ? parseInt(p.id, 10) : p.id;
        if (isChecked) {
          next.add(idNum);
        } else {
          next.delete(idNum);
        }
      });
      return next;
    });
  };

  const handleSave = async () => {
    if (!role) return;

    if (!name.trim()) {
      toast({
        variant: "destructive",
        title: "Erro",
        description: "O nome do perfil é obrigatório.",
      });
      return;
    }

    setIsSaving(true);
    try {
      await RoleService.updateRole(role.id, {
        name,
        description,
        permissionIds: Array.from(selectedPermissionIds),
      });

      toast({
        title: "Perfil atualizado",
        description: `O perfil "${name}" foi atualizado com sucesso.`,
      });
      onRoleUpdated();
      onOpenChange(false);
    } catch (error: any) {
      console.error("Erro ao atualizar perfil:", error);
      toast({
        variant: "destructive",
        title: "Erro ao atualizar perfil",
        description: error.message || "Ocorreu um erro ao salvar as alterações do perfil.",
      });
    } finally {
      setIsSaving(false);
    }
  };

  if (!role) return null;

  return (
    <Dialog open={isOpen} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-3xl max-h-[90vh] flex flex-col p-6 overflow-hidden">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2 text-xl font-bold">
            <ShieldAlert className="h-5 w-5 text-primary" />
            Editar Perfil de Acesso: {role.name}
          </DialogTitle>
          <DialogDescription>
            Atualize o nome, a descrição e altere quais permissões estão vinculadas a este perfil.
          </DialogDescription>
        </DialogHeader>

        <div className="flex-1 overflow-y-auto pr-1 py-4 space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="grid gap-2">
              <Label htmlFor="edit-role-name" className="text-sm font-medium">Nome do Perfil *</Label>
              <Input
                id="edit-role-name"
                placeholder="Ex: Vendedor, Gerente, Suporte"
                value={name}
                onChange={(e) => setName(e.target.value)}
              />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="edit-role-description" className="text-sm font-medium">Descrição</Label>
              <Input
                id="edit-role-description"
                placeholder="Explique o propósito deste perfil"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
              />
            </div>
          </div>

          <div className="space-y-4">
            <div className="border-b pb-2 flex items-center justify-between">
              <Label className="text-base font-semibold text-foreground">Permissões do Sistema</Label>
              <span className="text-xs text-muted-foreground bg-accent px-2.5 py-1 rounded-full font-medium">
                {selectedPermissionIds.size} selecionada(s)
              </span>
            </div>

            {isLoadingPermissions ? (
              <div className="flex flex-col items-center justify-center py-12">
                <Loader2 className="h-8 w-8 animate-spin text-primary mb-2" />
                <p className="text-sm text-muted-foreground">Carregando permissões do perfil...</p>
              </div>
            ) : Object.keys(permissionsByModule).length === 0 ? (
              <p className="text-sm text-muted-foreground text-center py-8">Nenhuma permissão disponível.</p>
            ) : (
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {Object.entries(permissionsByModule).map(([moduleName, items]) => {
                  const modulePermissionsIds = items.map((p) => typeof p.id === "string" ? parseInt(p.id, 10) : p.id);
                  const isAllModuleSelected = modulePermissionsIds.every((id) => selectedPermissionIds.has(id));
                  const isSomeModuleSelected = modulePermissionsIds.some((id) => selectedPermissionIds.has(id)) && !isAllModuleSelected;

                  return (
                    <div
                      key={moduleName}
                      className="border rounded-lg p-4 bg-card hover:border-primary/20 transition-all shadow-sm flex flex-col justify-between"
                    >
                      <div>
                        <div className="flex items-center justify-between border-b pb-2 mb-3">
                          <span className="font-semibold text-sm capitalize text-foreground">
                            {moduleName}
                          </span>
                          <div className="flex items-center gap-1.5">
                            <Checkbox
                              id={`edit-module-all-${moduleName}`}
                              checked={isAllModuleSelected ? true : isSomeModuleSelected ? "indeterminate" : false}
                              onCheckedChange={(checked) => handleToggleModuleAll(items, !!checked)}
                            />
                            <Label
                              htmlFor={`edit-module-all-${moduleName}`}
                              className="text-xs text-muted-foreground cursor-pointer hover:text-foreground transition-colors"
                            >
                              Selecionar Tudo
                            </Label>
                          </div>
                        </div>

                        <div className="space-y-2.5">
                          {items.map((permission) => {
                            const idNum = typeof permission.id === "string" ? parseInt(permission.id, 10) : permission.id;
                            const isChecked = selectedPermissionIds.has(idNum);

                            return (
                              <div key={permission.id} className="flex items-start gap-2.5">
                                <Checkbox
                                  id={`edit-perm-${permission.id}`}
                                  checked={isChecked}
                                  onCheckedChange={() => handleTogglePermission(idNum)}
                                  className="mt-0.5"
                                />
                                <div className="grid gap-0.5 leading-none">
                                  <Label
                                    htmlFor={`edit-perm-${permission.id}`}
                                    className="text-sm font-medium text-foreground cursor-pointer hover:text-primary transition-colors"
                                  >
                                    {permission.permission.split(".")[1] || permission.permission}
                                  </Label>
                                  {permission.description && (
                                    <p className="text-xs text-muted-foreground">
                                      {permission.description}
                                    </p>
                                  )}
                                </div>
                              </div>
                            );
                          })}
                        </div>
                      </div>
                    </div>
                  );
                })}
              </div>
            )}
          </div>
        </div>

        <DialogFooter className="border-t pt-4 mt-2">
          <Button
            variant="outline"
            onClick={() => onOpenChange(false)}
            disabled={isSaving}
          >
            Cancelar
          </Button>
          <Button onClick={handleSave} disabled={isSaving}>
            {isSaving ? (
              <>
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                Salvando...
              </>
            ) : (
              "Salvar Alterações"
            )}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
