// src/app/(system)/settings/roles/components/role-list/role-list.tsx
"use client";

import { useEffect, useState, useCallback, useMemo } from "react";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Label } from "@/components/ui/label";
import { Loader2 } from "lucide-react";
import { useToast } from "@/components/ui/use-toast";
import { Checkbox } from "@/components/ui/checkbox";
import { Role } from "@/services/auth/role-schema";
import { Permission } from "@/services/auth/permission-schema";
import { useServerPagination } from "@/hooks/use-server-pagination";
import { RolePermissionsTable } from "../tables/role-permissions-table";
import RoleService from "@/services/auth/role-service";
import PermissionService from "@/services/auth/permission-service";

interface RoleListProps {
  onRoleCreated: () => void; // Callback para quando um novo perfil for criado
}

export function RoleList({ onRoleCreated }: RoleListProps) {
  const { toast } = useToast();

  const [availableRoles, setAvailableRoles] = useState<Role[]>([]);
  const [selectedRoleId, setSelectedRoleId] = useState<string | null>(null);
  const [isLoadingRoles, setIsLoadingRoles] = useState(true);

  // Removido: permissionSearchName
  const [permissionSearchModule, setPermissionSearchModule] = useState(""); // Mantido para o filtro de módulo

  const [showLinkedPermissions, setShowLinkedPermissions] = useState(true);

  // NOVO: Estado para módulos de permissão disponíveis (para o Select de filtro)
  const [availablePermissionModules, setAvailablePermissionModules] = useState<
    string[]
  >([]);

  // Função para carregar todos os perfis para o Select
  const fetchAllRoles = useCallback(async () => {
    setIsLoadingRoles(true);
    try {
      const roles = await RoleService.getRoles();
      setAvailableRoles(roles);

      if (roles.length > 0) {
        if (
          !selectedRoleId ||
          !roles.some((role) => role.id.toString() === selectedRoleId)
        ) {
          setSelectedRoleId(roles[0].id.toString());
        }
      } else {
        setSelectedRoleId(null);
      }
    } catch (error: any) {
      console.error("Erro ao carregar perfis:", error);
      toast({
        title: "Erro",
        description:
          error.message || "Não foi possível carregar a lista de perfis.",
        variant: "destructive",
      });
      setAvailableRoles([]);
      setSelectedRoleId(null);
    } finally {
      setIsLoadingRoles(false);
    }
  }, [toast, selectedRoleId]);

  // NOVO: Função para carregar os módulos de permissão disponíveis
  const fetchAvailablePermissionModules = useCallback(async () => {
    try {
      const modules = await PermissionService.getAvailableModules();
      setAvailablePermissionModules(modules);
    } catch (error: any) {
      console.error("Erro ao carregar módulos de permissão:", error);
      toast({
        title: "Erro",
        description: "Não foi possível carregar os módulos de permissão.",
        variant: "destructive",
      });
      setAvailablePermissionModules([]);
    }
  }, [toast]);

  const {
    data: permissions,
    isLoading: isLoadingPermissions,
    currentPage,
    totalPages,
    totalItems,
    fetchDataWithFilters,
    goToNextPage,
    goToPreviousPage,
    goToFirstPage,
    goToLastPage,
  } = useServerPagination<Permission>({
    fetchDataService: PermissionService.getPermissions,
    pageSize: 5,
  });

  // Efeito para disparar a busca de permissões sempre que o perfil ou os filtros da tabela mudam
  useEffect(() => {
    if (selectedRoleId) {
      fetchDataWithFilters({
        module: permissionSearchModule === "all" ? "" : permissionSearchModule,
        roleId: selectedRoleId,
        isLinkedToRole: showLinkedPermissions,
      });
    } else {
      fetchDataWithFilters({}); // Busca sem filtro de roleId se "all" estiver selecionado
    }
  }, [
    selectedRoleId,
    permissionSearchModule,
    showLinkedPermissions,
    fetchDataWithFilters,
  ]);

  // Carrega todos os perfis quando o componente monta
  useEffect(() => {
    fetchAllRoles();
    fetchAvailablePermissionModules(); // Carrega os módulos junto
  }, [fetchAllRoles, fetchAvailablePermissionModules]);

  // Handler para quando o valor do Select de perfil muda
  const handleRoleSelectChange = (value: string) => {
    setSelectedRoleId(value === "all" ? null : value);
    // setPermissionSearchName(""); // Removido
    setPermissionSearchModule("all"); // Volta para "Todos os Módulos" ao mudar de perfil
    // showLinkedPermissions mantém o estado
  };

  // Handler para quando o filtro de módulo de permissão muda
  const handlePermissionModuleChange = (value: string) => {
    setPermissionSearchModule(value);
    // O useEffect já vai disparar a busca
  };

  const handlePermissionLinkUnlink = useCallback(() => {
    if (selectedRoleId) {
      fetchDataWithFilters(
        {
          // Removido: name: permissionSearchName,
          module:
            permissionSearchModule === "all" ? "" : permissionSearchModule,
          roleId: selectedRoleId,
          isLinkedToRole: showLinkedPermissions,
        },
        false
      );
    }
  }, [
    selectedRoleId,
    permissionSearchModule,
    showLinkedPermissions,
    fetchDataWithFilters,
  ]);

  const currentSelectedRole = useMemo(() => {
    return availableRoles.find((role) => role.id.toString() === selectedRoleId);
  }, [availableRoles, selectedRoleId]);

  return (
    <div className="space-y-6">
      {/* Seletor de Perfis e Filtros de Permissão na mesma linha */}
      <div className="flex flex-wrap items-end gap-4 mb-6">
        {" "}
        {/* Use flex-wrap para responsividade */}
        {/* Seletor de Perfis */}
        <div className="grid gap-1.5 min-w-[280px]">
          {" "}
          {/* Adiciona um min-width para o Select */}
          <Label htmlFor="select-role" className="text-sm font-medium">
            Selecionar Perfil
          </Label>
          {isLoadingRoles ? (
            <div className="flex items-center space-x-2 h-10 border rounded-md px-3 py-2">
              {" "}
              {/* Estilo para alinhar com o Select */}
              <Loader2 className="h-5 w-5 animate-spin text-primary" />
              <p className="text-muted-foreground">Carregando perfis...</p>
            </div>
          ) : (
            <Select
              value={selectedRoleId || "all"}
              onValueChange={handleRoleSelectChange}
              disabled={availableRoles.length === 0}
            >
              <SelectTrigger id="select-role" className="w-full">
                <SelectValue placeholder="Selecione um perfil" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">Todos os Perfis</SelectItem>{" "}
                {availableRoles.map((role) => (
                  <SelectItem key={role.id} value={role.id.toString()}>
                    {role.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          )}
        </div>
        {/* NOVO: Filtro de Módulo para Permissões (apenas se um perfil estiver selecionado) */}
        {selectedRoleId && (
          <div className="grid gap-1.5 min-w-[200px]">
            <Label
              htmlFor="permission-module-filter"
              className="text-sm font-medium"
            >
              Filtrar por Módulo
            </Label>
            <Select
              value={permissionSearchModule || "all"}
              onValueChange={handlePermissionModuleChange}
            >
              <SelectTrigger id="permission-module-filter" className="w-full">
                <SelectValue placeholder="Todos os Módulos" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">Todos os Módulos</SelectItem>
                {availablePermissionModules.map((module) => (
                  <SelectItem key={module} value={module}>
                    {module}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        )}
        {/* Checkbox "Permissões Vinculadas" (apenas se um perfil estiver selecionado) */}
        {selectedRoleId && (
          <div className="flex items-center space-x-2 self-end mb-1">
            <Checkbox
              id="show-linked-permissions"
              checked={showLinkedPermissions}
              onCheckedChange={(checked) => setShowLinkedPermissions(!!checked)}
            />
            <Label
              htmlFor="show-linked-permissions"
              className="text-sm font-medium"
            >
              Mostrar Permissões Vinculadas
            </Label>
          </div>
        )}
      </div>

      {/* Tabela de Permissões do Perfil Selecionado */}
      {selectedRoleId ? (
        <div className="mt-8">
          <h2 className="text-xl font-semibold mb-4">
            Permissões do Perfil: {currentSelectedRole?.name || "Carregando..."}
          </h2>
          {isLoadingPermissions ? (
            <div className="flex flex-col items-center justify-center h-48 border rounded-md">
              <Loader2 className="h-8 w-8 animate-spin text-primary mb-3" />
              <p className="text-muted-foreground">
                Carregando permissões do perfil...
              </p>
            </div>
          ) : (
            <RolePermissionsTable
              permissions={permissions || []}
              roleId={selectedRoleId}
              onPermissionsChanged={handlePermissionLinkUnlink}
              currentPage={currentPage}
              totalPages={totalPages}
              totalItems={totalItems}
              onPreviousPage={goToPreviousPage}
              onNextPage={goToNextPage}
              onGoToFirstPage={goToFirstPage}
              onGoToLastPage={goToLastPage}
              showLinkedPermissions={showLinkedPermissions}
            />
          )}
        </div>
      ) : (
        <div className="mt-4">
          {!isLoadingRoles && availableRoles.length === 0 && (
            <p className="text-muted-foreground">
              Nenhum perfil disponível. Crie um novo perfil para começar.
            </p>
          )}
          {!isLoadingRoles && availableRoles.length > 0 && (
            <p className="text-muted-foreground">
              Selecione um perfil acima para visualizar e gerenciar suas
              permissões.
            </p>
          )}
        </div>
      )}
    </div>
  );
}
