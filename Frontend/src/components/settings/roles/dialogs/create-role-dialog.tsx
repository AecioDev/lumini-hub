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
import { useToast } from "@/components/ui/use-toast";
import { Loader2, Shield, Search, X } from "lucide-react";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Permission } from "@/services/auth/permission-schema";
import PermissionService from "@/services/auth/permission-service";
import RoleService from "@/services/auth/role-service";

interface CreateRoleDialogProps {
  isOpen: boolean;
  onOpenChange: (open: boolean) => void;
  onRoleCreated: () => void;
}

export function CreateRoleDialog({
  isOpen,
  onOpenChange,
  onRoleCreated,
}: CreateRoleDialogProps) {
  const { toast } = useToast();
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [permissionsByModule, setPermissionsByModule] = useState<Record<string, Permission[]>>({});
  const [selectedPermissionIds, setSelectedPermissionIds] = useState<Set<number>>(new Set());
  const [isLoadingPermissions, setIsLoadingPermissions] = useState(false);
  const [isSaving, setIsSaving] = useState(false);

  // Filtros locais para personalização de permissões
  const [searchQuery, setSearchQuery] = useState("");
  const [filterModule, setFilterModule] = useState("all");

  // Lista plana de todas as permissões para filtragem
  const allPermissionsFlat = useMemo(() => {
    const flatList: (Permission & { module: string })[] = [];
    Object.entries(permissionsByModule).forEach(([moduleName, items]) => {
      items.forEach((p) => {
        flatList.push({ ...p, module: moduleName });
      });
    });
    return flatList.sort((a, b) => {
      const modCompare = a.module.localeCompare(b.module);
      if (modCompare !== 0) return modCompare;
      return a.permission.localeCompare(b.permission);
    });
  }, [permissionsByModule]);

  // Lista de módulos para o Select
  const modulesList = useMemo(() => {
    return Object.keys(permissionsByModule).sort();
  }, [permissionsByModule]);

  // Permissões filtradas
  const filteredPermissions = useMemo(() => {
    return allPermissionsFlat.filter((p) => {
      if (searchQuery.trim()) {
        const q = searchQuery.toLowerCase();
        const nameMatch = p.permission.toLowerCase().includes(q);
        const descMatch = p.description?.toLowerCase().includes(q) || false;
        const modMatch = p.module.toLowerCase().includes(q);
        if (!nameMatch && !descMatch && !modMatch) return false;
      }
      if (filterModule !== "all" && p.module !== filterModule) {
        return false;
      }
      return true;
    });
  }, [allPermissionsFlat, searchQuery, filterModule]);

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
  };

  // Carregar as permissões agrupadas por módulo ao abrir o diálogo
  useEffect(() => {
    if (isOpen) {
      const loadPermissions = async () => {
        setIsLoadingPermissions(true);
        try {
          const data = await PermissionService.getPermissionsByModule();
          setPermissionsByModule(data);
        } catch (error: any) {
          console.error("Erro ao carregar permissões por módulo:", error);
          toast({
            variant: "destructive",
            title: "Erro",
            description: "Não foi possível carregar as permissões disponíveis.",
          });
        } finally {
          setIsLoadingPermissions(false);
        }
      };

      loadPermissions();
      // Resetar form
      setName("");
      setDescription("");
      setSelectedPermissionIds(new Set());
    }
  }, [isOpen, toast]);

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
      await RoleService.createRole({
        name,
        description,
        permissionIds: Array.from(selectedPermissionIds),
      });

      toast({
        title: "Perfil criado",
        description: `O perfil "${name}" foi criado com sucesso.`,
      });
      onRoleCreated();
      onOpenChange(false);
    } catch (error: any) {
      console.error("Erro ao criar perfil:", error);
      toast({
        variant: "destructive",
        title: "Erro ao criar perfil",
        description: error.message || "Ocorreu um erro ao salvar o perfil.",
      });
    } finally {
      setIsSaving(false);
    }
  };

  return (
    <Dialog open={isOpen} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-3xl max-h-[90vh] flex flex-col p-6 overflow-hidden">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2 text-xl font-bold">
            <Shield className="h-5 w-5 text-primary" />
            Novo Perfil de Acesso
          </DialogTitle>
          <DialogDescription>
            Defina o nome, a descrição e configure quais permissões de sistema este perfil terá.
          </DialogDescription>
        </DialogHeader>

        <div className="flex-1 overflow-y-auto pr-1 py-4 space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="grid gap-2">
              <Label htmlFor="role-name" className="text-sm font-medium">Nome do Perfil *</Label>
              <Input
                id="role-name"
                placeholder="Ex: Vendedor, Gerente, Suporte"
                value={name}
                onChange={(e) => setName(e.target.value)}
              />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="role-description" className="text-sm font-medium">Descrição</Label>
              <Input
                id="role-description"
                placeholder="Explique o propósito deste perfil"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
              />
            </div>
          </div>

          <div className="space-y-4 pt-4 border-t flex flex-col min-h-0 flex-1">
            <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-2">
              <div>
                <Label className="text-sm font-bold">Permissões do Sistema</Label>
                <p className="text-xs text-muted-foreground mt-0.5">
                  Selecione as permissões que pertencerão a este novo perfil de acesso.
                </p>
              </div>
              <span className="text-xs text-muted-foreground bg-accent px-2.5 py-1 rounded-full font-medium shrink-0 self-start sm:self-center">
                {selectedPermissionIds.size} selecionada(s)
              </span>
            </div>

            {/* Barra de Filtros */}
            <div className="grid grid-cols-1 sm:grid-cols-12 gap-2 bg-muted/30 p-2.5 rounded-lg border">
              <div className="sm:col-span-8 relative">
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
              <div className="sm:col-span-1 flex justify-center items-center">
                <Button
                  variant="ghost"
                  size="icon"
                  onClick={handleClearFilters}
                  disabled={!searchQuery && filterModule === "all"}
                  className="h-9 w-9 shrink-0"
                  title="Limpar Filtros"
                >
                  <X className="h-4 w-4" />
                </Button>
              </div>
            </div>

            {isLoadingPermissions ? (
              <div className="flex flex-col items-center justify-center py-12 border rounded-lg bg-muted/10">
                <Loader2 className="h-7 w-7 animate-spin text-primary mb-2.5" />
                <p className="text-xs text-muted-foreground font-medium">Carregando permissões...</p>
              </div>
            ) : filteredPermissions.length === 0 ? (
              <div className="flex flex-col items-center justify-center py-12 border rounded-lg bg-muted/10 text-center px-4">
                <p className="text-xs text-muted-foreground font-medium">Nenhuma permissão encontrada com os filtros aplicados.</p>
                {(searchQuery || filterModule !== "all") && (
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
                      id="create-role-select-all-filtered"
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
                            id={`create-role-perm-${permission.id}`}
                            checked={isChecked}
                            onCheckedChange={() => handleTogglePermission(idNum)}
                            className="h-4 w-4"
                          />
                        </div>
                        <div className="col-span-8 pr-2">
                          <Label
                            htmlFor={`create-role-perm-${permission.id}`}
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
                Criando Perfil...
              </>
            ) : (
              "Criar Perfil"
            )}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
