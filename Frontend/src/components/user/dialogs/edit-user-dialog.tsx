"use client";

import { useEffect, useState, useMemo } from "react";
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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useToast } from "@/components/ui/use-toast";
import { Edit2, Loader2, ShieldCheck, Search, X, Filter } from "lucide-react";
import { Role } from "@/services/auth/role-schema";
import { User } from "@/services/auth/user-schema";
import { Permission } from "@/services/auth/permission-schema";
import UserService from "@/services/auth/user-service";
import RoleService from "@/services/auth/role-service";
import PermissionService from "@/services/auth/permission-service";

interface EditUserDialogProps {
  isOpen: boolean;
  onOpenChange: (open: boolean) => void;
  user: User | null;
  roles: Role[];
  onUserUpdated: (user: User) => void;
  onUserChange: (user: User) => void;
}

export function EditUserDialog({
  isOpen,
  onOpenChange,
  user,
  roles,
  onUserUpdated,
  onUserChange,
}: EditUserDialogProps) {
  // Permissões do sistema
  const [allPermissionsByModule, setAllPermissionsByModule] = useState<Record<string, Permission[]>>({});
  const [isLoadingAllPermissions, setIsLoadingAllPermissions] = useState(false);

  // Controle de permissões selecionadas e herdadas da Role
  const [defaultRolePermissionIds, setDefaultRolePermissionIds] = useState<Set<number>>(new Set());
  const [selectedPermissionIds, setSelectedPermissionIds] = useState<Set<number>>(new Set());
  const [isLoadingUserPermissions, setIsLoadingUserPermissions] = useState(false);
  
  const [isSaving, setIsSaving] = useState(false);

  // Filtros locais para personalização de permissões
  const [searchQuery, setSearchQuery] = useState("");
  const [filterModule, setFilterModule] = useState("all");
  const [filterRoleId, setFilterRoleId] = useState("all");
  const [filterRolePermIds, setFilterRolePermIds] = useState<Set<number>>(new Set());
  const [isLoadingFilterRole, setIsLoadingFilterRole] = useState(false);

  // Carrega permissões do perfil filtrado
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
        if (roleDetails.permissions) {
          roleDetails.permissions.forEach((p) => {
            const idNum = typeof p.id === "string" ? parseInt(p.id, 10) : p.id;
            ids.add(idNum);
          });
        }
        setFilterRolePermIds(ids);
      } catch (error) {
        console.error("Erro ao carregar permissões do perfil filtrado:", error);
      } finally {
        setIsLoadingFilterRole(false);
      }
    };
    loadFilterRolePerms();
  }, [filterRoleId]);

  // Lista plana de todas as permissões para filtragem fácil
  const allPermissionsFlat = useMemo(() => {
    const flatList: (Permission & { module: string })[] = [];
    Object.entries(allPermissionsByModule).forEach(([moduleName, items]) => {
      items.forEach((p) => {
        flatList.push({ ...p, module: moduleName });
      });
    });
    return flatList.sort((a, b) => {
      const modCompare = a.module.localeCompare(b.module);
      if (modCompare !== 0) return modCompare;
      return a.permission.localeCompare(b.permission);
    });
  }, [allPermissionsByModule]);

  // Lista de módulos extraída das permissões carregadas
  const modulesList = useMemo(() => {
    return Object.keys(allPermissionsByModule).sort();
  }, [allPermissionsByModule]);

  // Permissões filtradas localmente
  const filteredPermissions = useMemo(() => {
    return allPermissionsFlat.filter((p) => {
      const idNum = typeof p.id === "string" ? parseInt(p.id, 10) : p.id;

      // 1. Filtro de busca de texto
      if (searchQuery.trim()) {
        const q = searchQuery.toLowerCase();
        const nameMatch = p.permission.toLowerCase().includes(q);
        const descMatch = p.description?.toLowerCase().includes(q) || false;
        const modMatch = p.module.toLowerCase().includes(q);
        if (!nameMatch && !descMatch && !modMatch) return false;
      }

      // 2. Filtro de módulo
      if (filterModule !== "all" && p.module !== filterModule) {
        return false;
      }

      // 3. Filtro de perfil
      if (filterRoleId !== "all" && !filterRolePermIds.has(idNum)) {
        return false;
      }

      return true;
    });
  }, [allPermissionsFlat, searchQuery, filterModule, filterRoleId, filterRolePermIds]);

  const isAllFilteredChecked = useMemo(() => {
    if (filteredPermissions.length === 0) return false;
    return filteredPermissions.every((p) => {
      const idNum = typeof p.id === "string" ? parseInt(p.id, 10) : p.id;
      return selectedPermissionIds.has(idNum);
    });
  }, [filteredPermissions, selectedPermissionIds]);

  const handleSelectAllFiltered = (checked: boolean) => {
    setSelectedPermissionIds((prev) => {
      const next = new Set(prev);
      filteredPermissions.forEach((p) => {
        const idNum = typeof p.id === "string" ? parseInt(p.id, 10) : p.id;
        if (checked) {
          next.add(idNum);
        } else {
          next.delete(idNum);
        }
      });
      return next;
    });
  };

  const handleClearFilters = () => {
    setSearchQuery("");
    setFilterModule("all");
    setFilterRoleId("all");
  };

  const { toast } = useToast();
  const validRoles = roles.filter((role) => role.id && role.name);

  // 1. Carrega todas as permissões do sistema quando o modal abre
  useEffect(() => {
    if (isOpen) {
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
    }
  }, [isOpen]);

  // 2. Carrega as permissões da Role do usuário e suas permissões ativas atuais
  useEffect(() => {
    const loadUserPermissions = async () => {
      if (!user || !user.role_id) return;
      setIsLoadingUserPermissions(true);
      try {
        // Pega as permissões padrão da Role
        const roleDetails = await RoleService.getRoleById(user.role_id);
        const roleIds = new Set<number>();
        if (roleDetails.permissions) {
          roleDetails.permissions.forEach((p) => {
            const idNum = typeof p.id === "string" ? parseInt(p.id, 10) : p.id;
            roleIds.add(idNum);
          });
        }
        setDefaultRolePermissionIds(roleIds);

        // Pega os detalhes completos do usuário no banco (que inclui as permissões em user_permissions)
        const detailedUser = await UserService.getUserById(user.id);
        console.log("[EditUserDialog] detailedUser obtido:", detailedUser);
        console.log("[EditUserDialog] detailedUser.permissions:", detailedUser?.permissions);
        const activeIds = new Set<number>();
        if (detailedUser.permissions) {
          detailedUser.permissions.forEach((p) => {
            const idNum = typeof p.id === "string" ? parseInt(p.id, 10) : p.id;
            activeIds.add(idNum);
          });
        }
        console.log("[EditUserDialog] activeIds identificados:", Array.from(activeIds));
        setSelectedPermissionIds(activeIds);
      } catch (error) {
        console.error("Erro ao carregar permissões detalhadas do usuário:", error);
      } finally {
        setIsLoadingUserPermissions(false);
      }
    };

    if (isOpen && user) {
      loadUserPermissions();
    }
  }, [isOpen, user?.id]);

  // 3. Ao mudar a seleção de perfil (RoleID) no select, atualiza os templates
  const handleRoleSelectChange = async (roleIdStr: string) => {
    const newRoleId = Number.parseInt(roleIdStr);
    if (!user) return;

    // Atualiza a Role no estado do usuário
    onUserChange({ ...user, role_id: newRoleId });

    // Busca as novas permissões do perfil selecionado e atualiza
    setIsLoadingUserPermissions(true);
    try {
      const roleDetails = await RoleService.getRoleById(newRoleId);
      const roleIds = new Set<number>();
      if (roleDetails.permissions) {
        roleDetails.permissions.forEach((p) => {
          const idNum = typeof p.id === "string" ? parseInt(p.id, 10) : p.id;
          roleIds.add(idNum);
        });
      }
      setDefaultRolePermissionIds(roleIds);
      setSelectedPermissionIds(new Set(roleIds)); // Quando muda a role, resetamos para o padrão dela
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

  const handleUpdateUser = async () => {
    if (!user) return;

    setIsSaving(true);
    try {
      // 1. Atualizar informações básicas (e perfil)
      const updatedUser = await UserService.updateUser(user.id, {
        name: user.name,
        email: user.email,
        role_id: user.role_id,
        is_active: user.is_active,
      });

      // 2. Atualizar permissões customizadas diretas do usuário
      await UserService.updateUserPermissions(user.id, Array.from(selectedPermissionIds));

      onUserUpdated(updatedUser);
      onOpenChange(false);

      toast({
        title: "Usuário atualizado",
        description: `O usuário ${updatedUser.name} foi atualizado com acessos personalizados.`,
      });
    } catch (error: any) {
      toast({
        variant: "destructive",
        title: "Erro ao atualizar usuário",
        description:
          error.response?.data?.message ||
          "Ocorreu um erro ao atualizar o usuário",
      });
    } finally {
      setIsSaving(false);
    }
  };

  if (!user) return null;

  return (
    <Dialog open={isOpen} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-2xl max-h-[90vh] flex flex-col p-6 overflow-hidden">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2 text-xl font-bold">
            <Edit2 className="h-5 w-5 text-primary" />
            Editar Usuário
          </DialogTitle>
          <DialogDescription>
            Atualize os dados e personalize as permissões de acesso (extras e overrides) do usuário.
          </DialogDescription>
        </DialogHeader>

        <div className="flex-1 overflow-y-auto pr-1 py-4 space-y-6">
          {/* Informações Cadastrais */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="grid gap-1.5 col-span-full">
              <Label htmlFor="edit-name" className="text-sm font-medium">Nome Completo</Label>
              <Input
                id="edit-name"
                value={user.name}
                onChange={(e) => onUserChange({ ...user, name: e.target.value })}
              />
            </div>

            <div className="grid gap-1.5">
              <Label htmlFor="edit-email" className="text-sm font-medium">Email</Label>
              <Input
                id="edit-email"
                type="email"
                placeholder="exemplo@lumini.com.br"
                value={user.email || ""}
                onChange={(e) => onUserChange({ ...user, email: e.target.value })}
              />
            </div>

            <div className="grid gap-1.5">
              <Label htmlFor="edit-role" className="text-sm font-medium">Perfil de Acesso (Modelo)</Label>
              <Select
                value={user.role_id ? user.role_id.toString() : ""}
                onValueChange={handleRoleSelectChange}
              >
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
            </div>
          </div>

          {/* Configuração de Permissões (Híbrido: herdado + overrides) */}
          {user.role_id > 0 && (
            <div className="space-y-4 pt-4 border-t flex flex-col min-h-0 flex-1">
              <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-2">
                <div>
                  <Label className="text-sm font-bold flex items-center gap-1">
                    <ShieldCheck className="h-4 w-4 text-primary" />
                    Personalização de Permissões de Acesso
                  </Label>
                  <p className="text-xs text-muted-foreground mt-0.5">
                    As permissões padrão do perfil estão marcadas. Adicione permissões extras ou revogue acessos.
                  </p>
                </div>
                <span className="text-xs text-muted-foreground bg-accent px-2.5 py-1 rounded-full font-medium shrink-0 self-start sm:self-center">
                  {selectedPermissionIds.size} ativas
                </span>
              </div>

              {/* Barra de Filtros */}
              <div className="grid grid-cols-1 sm:grid-cols-12 gap-2 bg-muted/30 p-2.5 rounded-lg border">
                <div className="sm:col-span-5 relative">
                  <Search className="absolute left-2.5 top-2.5 h-3.5 w-3.5 text-muted-foreground" />
                  <Input
                    placeholder="Buscar permissão..."
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    className="pl-8 h-9 text-xs"
                  />
                </div>
                <div className="sm:col-span-3">
                  <Select value={filterModule} onValueChange={setFilterModule}>
                    <SelectTrigger className="h-9 text-xs">
                      <SelectValue placeholder="Módulo" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="all">Todos Módulos</SelectItem>
                      {modulesList.map((mod) => (
                        <SelectItem key={mod} value={mod} className="capitalize text-xs">
                          {mod}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>
                <div className="sm:col-span-3">
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
                </div>
                <div className="sm:col-span-1 flex justify-center items-center">
                  <Button
                    variant="ghost"
                    size="icon"
                    onClick={handleClearFilters}
                    disabled={!searchQuery && filterModule === "all" && filterRoleId === "all"}
                    className="h-9 w-9 shrink-0"
                    title="Limpar Filtros"
                  >
                    <X className="h-4 w-4" />
                  </Button>
                </div>
              </div>

              {isLoadingUserPermissions || isLoadingAllPermissions || isLoadingFilterRole ? (
                <div className="flex flex-col items-center justify-center py-12 border rounded-lg bg-muted/10">
                  <Loader2 className="h-7 w-7 animate-spin text-primary mb-2.5" />
                  <p className="text-xs text-muted-foreground font-medium">Carregando permissões...</p>
                </div>
              ) : filteredPermissions.length === 0 ? (
                <div className="flex flex-col items-center justify-center py-12 border rounded-lg bg-muted/10 text-center px-4">
                  <p className="text-xs text-muted-foreground font-medium">Nenhuma permissão encontrada com os filtros aplicados.</p>
                  {(searchQuery || filterModule !== "all" || filterRoleId !== "all") && (
                    <Button variant="link" size="sm" onClick={handleClearFilters} className="text-xs mt-1 h-auto p-0">
                      Limpar Filtros
                    </Button>
                  )}
                </div>
              ) : (
                <div className="flex flex-col border rounded-lg overflow-hidden bg-card flex-1 min-h-[220px]">
                  {/* Tabela Header */}
                  <div className="grid grid-cols-12 gap-4 px-4 py-2 border-b bg-muted/50 text-[10px] font-bold uppercase tracking-wider text-muted-foreground items-center">
                    <div className="col-span-1 flex items-center">
                      <Checkbox
                        id="edit-user-select-all-filtered"
                        checked={isAllFilteredChecked}
                        onCheckedChange={(checked) => handleSelectAllFiltered(!!checked)}
                        className="h-4 w-4"
                      />
                    </div>
                    <div className="col-span-8">Permissão</div>
                    <div className="col-span-3">Módulo</div>
                  </div>

                  {/* Tabela Body */}
                  <div className="overflow-y-auto max-h-64 flex-1 divide-y">
                    {filteredPermissions.map((permission) => {
                      const idNum = typeof permission.id === "string" ? parseInt(permission.id, 10) : permission.id;
                      const isChecked = selectedPermissionIds.has(idNum);

                      return (
                        <div
                          key={permission.id}
                          className="grid grid-cols-12 gap-4 px-4 py-2 text-xs items-center hover:bg-muted/10 transition-colors"
                        >
                          <div className="col-span-1 flex items-center">
                            <Checkbox
                              id={`edit-user-perm-${permission.id}`}
                              checked={isChecked}
                              onCheckedChange={() => handleTogglePermission(idNum)}
                              className="h-4 w-4"
                            />
                          </div>
                          <div className="col-span-8 pr-2">
                            <Label
                              htmlFor={`edit-user-perm-${permission.id}`}
                              className="font-semibold text-foreground cursor-pointer block truncate text-xs"
                              title={permission.description || permission.permission}
                            >
                              {permission.description || permission.permission}
                            </Label>
                            {permission.description && (
                              <p className="text-[10px] text-muted-foreground font-mono truncate" title={permission.permission}>
                                {permission.permission}
                              </p>
                            )}
                          </div>
                          <div className="col-span-3">
                            <span className="inline-block text-[10px] bg-slate-100 dark:bg-slate-800 text-slate-700 dark:text-slate-300 px-2 py-0.5 rounded capitalize font-medium">
                              {permission.module}
                            </span>
                          </div>
                        </div>
                      );
                    })}
                  </div>
                </div>
              )}
            </div>
          )}
        </div>

        <DialogFooter className="border-t pt-4">
          <Button variant="outline" onClick={() => onOpenChange(false)} disabled={isSaving}>
            Cancelar
          </Button>
          <Button onClick={handleUpdateUser} disabled={isSaving}>
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
