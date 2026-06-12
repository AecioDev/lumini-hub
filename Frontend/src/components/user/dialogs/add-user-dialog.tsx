"use client";

import { useState, useEffect, useMemo } from "react";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
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
import { UserPlus, Loader2, ShieldCheck, ShieldAlert, Search, X, Filter } from "lucide-react";
import { CreateUserDto, User } from "@/services/auth/user-schema";
import { Role } from "@/services/auth/role-schema";
import { Permission } from "@/services/auth/permission-schema";
import UserService from "@/services/auth/user-service";
import RoleService from "@/services/auth/role-service";
import PermissionService from "@/services/auth/permission-service";

interface AddUserDialogProps {
  roles: Role[];
  onUserAdded: (user: User) => void;
}

export function AddUserDialog({ roles, onUserAdded }: AddUserDialogProps) {
  const [isOpen, setIsOpen] = useState(false);
  const [newUser, setNewUser] = useState<CreateUserDto>({
    username: "",
    password: "",
    name: "",
    email: "",
    role_id: 0,
  });

  // Permissões do sistema
  const [allPermissionsByModule, setAllPermissionsByModule] = useState<Record<string, Permission[]>>({});
  const [isLoadingAllPermissions, setIsLoadingAllPermissions] = useState(false);

  // Controle de permissões selecionadas e herdadas da Role
  const [defaultRolePermissionIds, setDefaultRolePermissionIds] = useState<Set<number>>(new Set());
  const [selectedPermissionIds, setSelectedPermissionIds] = useState<Set<number>>(new Set());
  const [isLoadingRolePermissions, setIsLoadingRolePermissions] = useState(false);
  
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

  // Carregar todas as permissões do sistema quando o modal abrir
  useEffect(() => {
    if (isOpen) {
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
    }
  }, [isOpen, toast]);

  // Escuta a mudança de role_id para carregar as permissões vinculadas por padrão
  useEffect(() => {
    const loadRolePermissions = async () => {
      if (!newUser.role_id) {
        setDefaultRolePermissionIds(new Set());
        setSelectedPermissionIds(new Set());
        return;
      }
      setIsLoadingRolePermissions(true);
      try {
        const roleDetails = await RoleService.getRoleById(newUser.role_id);
        const ids = new Set<number>();
        if (roleDetails.permissions) {
          roleDetails.permissions.forEach((p) => {
            const idNum = typeof p.id === "string" ? parseInt(p.id, 10) : p.id;
            ids.add(idNum);
          });
        }
        setDefaultRolePermissionIds(new Set(ids));
        setSelectedPermissionIds(new Set(ids)); // Inicializa marcando as do perfil
      } catch (error) {
        console.error("Erro ao carregar permissões do perfil:", error);
      } finally {
        setIsLoadingRolePermissions(false);
      }
    };

    loadRolePermissions();
  }, [newUser.role_id]);

  const resetForm = () => {
    setNewUser({
      username: "",
      password: "",
      name: "",
      email: "",
      role_id: 0,
    });
    setDefaultRolePermissionIds(new Set());
    setSelectedPermissionIds(new Set());
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

  const handleAddUser = async () => {
    if (
      !newUser.username ||
      !newUser.password ||
      !newUser.name ||
      !newUser.role_id
    ) {
      toast({
        variant: "destructive",
        title: "Erro ao adicionar usuário",
        description: "Preencha todos os campos obrigatórios (*)",
      });
      return;
    }

    setIsSaving(true);
    try {
      // 1. Criar o usuário básico (e herda as permissões da Role automaticamente no backend)
      const createdUser = await UserService.createUser(newUser);

      // 2. Salvar o conjunto personalizado de permissões em lote (se houver modificações da Role padrão)
      // Para consistência, salvamos a lista de permissões que o administrador configurou no grid.
      await UserService.updateUserPermissions(createdUser.id, Array.from(selectedPermissionIds));

      onUserAdded(createdUser);
      resetForm();
      setIsOpen(false);

      toast({
        title: "Usuário adicionado",
        description: `O usuário ${createdUser.name} foi criado com acessos personalizados.`,
      });
    } catch (error: any) {
      toast({
        variant: "destructive",
        title: "Erro ao adicionar usuário",
        description:
          error.response?.data?.message ||
          "Ocorreu um erro ao adicionar o usuário",
      });
    } finally {
      setIsSaving(false);
    }
  };

  return (
    <Dialog open={isOpen} onOpenChange={(open) => {
      setIsOpen(open);
      if (!open) resetForm();
    }}>
      <DialogTrigger asChild>
        <Button className="shadow-md transition-all duration-200">
          <UserPlus className="mr-2 h-4 w-4" />
          Novo Usuário
        </Button>
      </DialogTrigger>
      <DialogContent className="max-w-2xl max-h-[90vh] flex flex-col p-6 overflow-hidden">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2 text-xl font-bold">
            <UserPlus className="h-5 w-5 text-primary" />
            Adicionar Novo Usuário
          </DialogTitle>
          <DialogDescription>
            Preencha os dados do usuário e selecione as permissões de acesso (com suporte a customização).
          </DialogDescription>
        </DialogHeader>

        <div className="flex-1 overflow-y-auto pr-1 py-4 space-y-6">
          {/* Informações Cadastrais */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="grid gap-1.5 col-span-full">
              <Label htmlFor="name" className="text-sm font-medium">Nome Completo *</Label>
              <Input
                id="name"
                placeholder="Ex: João Silva"
                value={newUser.name}
                onChange={(e) => setNewUser({ ...newUser, name: e.target.value })}
              />
            </div>

            <div className="grid gap-1.5">
              <Label htmlFor="username" className="text-sm font-medium">Nome de Usuário *</Label>
              <Input
                id="username"
                placeholder="joao_silva"
                value={newUser.username}
                onChange={(e) =>
                  setNewUser({ ...newUser, username: e.target.value })
                }
              />
            </div>
            <div className="grid gap-1.5">
              <Label htmlFor="password" className="text-sm font-medium">Senha *</Label>
              <Input
                id="password"
                type="password"
                placeholder="Mínimo 6 dígitos"
                value={newUser.password}
                onChange={(e) =>
                  setNewUser({ ...newUser, password: e.target.value })
                }
              />
            </div>

            <div className="grid gap-1.5">
              <Label htmlFor="email" className="text-sm font-medium">Email</Label>
              <Input
                id="email"
                type="email"
                placeholder="exemplo@lumini.com.br"
                value={newUser.email}
                onChange={(e) =>
                  setNewUser({ ...newUser, email: e.target.value })
                }
              />
            </div>

            <div className="grid gap-1.5">
              <Label htmlFor="role" className="text-sm font-medium">Perfil de Acesso (Modelo) *</Label>
              <Select
                value={newUser.role_id ? newUser.role_id.toString() : ""}
                onValueChange={(value) =>
                  setNewUser({ ...newUser, role_id: Number.parseInt(value) })
                }
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
            </div>
          </div>

          {/* Configuração de Permissões (Híbrido: herdado + overrides) */}
          {newUser.role_id > 0 && (
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

              {isLoadingRolePermissions || isLoadingAllPermissions || isLoadingFilterRole ? (
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
                        id="add-user-select-all-filtered"
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
                              id={`add-user-perm-${permission.id}`}
                              checked={isChecked}
                              onCheckedChange={() => handleTogglePermission(idNum)}
                              className="h-4 w-4"
                            />
                          </div>
                          <div className="col-span-8 pr-2">
                            <Label
                              htmlFor={`add-user-perm-${permission.id}`}
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
          <Button variant="outline" onClick={() => setIsOpen(false)} disabled={isSaving}>
            Cancelar
          </Button>
          <Button onClick={handleAddUser} disabled={isSaving}>
            {isSaving ? (
              <>
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                Criando...
              </>
            ) : (
              "Criar Usuário"
            )}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
