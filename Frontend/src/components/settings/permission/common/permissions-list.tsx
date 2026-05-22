// src/components/settings/permissions-list.tsx
"use client";

import { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Loader2 } from "lucide-react";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Permission } from "@/services/auth/permission-schema";
import { ConfirmDeleteDialog } from "@/components/common/confirm-delete-dialog";
import { PermissionsTable } from "../tables/permissions-table";
import { Role } from "@/services/auth/role-schema";
import RoleService from "@/services/auth/role-service";
import PermissionService from "@/services/auth/permission-service";
import { useServerPagination } from "@/hooks/use-server-pagination";

export function PermissionsList() {
  const [searchName, setSearchName] = useState("");
  const [searchModule, setSearchModule] = useState("");
  const [searchRoleId, setSearchRoleId] = useState("");

  const {
    data: permissions, // Renomeie 'data' para 'permissions' para clareza na tela
    isLoading,
    currentPage,
    totalPages,
    totalItems,
    fetchDataWithFilters, // Função para buscar dados com filtros
    goToNextPage,
    goToPreviousPage,
    goToFirstPage,
    goToLastPage,
  } = useServerPagination<Permission>({
    // Especifique o tipo da entidade aqui
    fetchDataService: PermissionService.getPermissions, // <--- Passando a função específica
    pageSize: 5,
  });

  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
  const [availableModules, setAvailableModules] = useState<string[]>([]);
  const [availableRoles, setAvailableRoles] = useState<Role[]>([]);
  const [selectedPermission, setSelectedPermission] =
    useState<Permission | null>(null);

  async function getModules() {
    try {
      const data = await PermissionService.getAvailableModules();
      setAvailableModules(data);
    } catch (error) {
      console.log(error);
    }
  }

  async function getRoles() {
    try {
      const data = await RoleService.getRoles();
      setAvailableRoles(data);
    } catch (error) {
      console.log(error);
    }
  }

  useEffect(() => {
    getModules();
    getRoles();
    // A primeira busca é disparada pelo useEffect interno do useServerPagination
    // que usa a fetchDataService injetada com os filtros iniciais (vazios neste caso).
  }, []);

  const handleSearch = () => {
    fetchDataWithFilters({
      name: searchName,
      module: searchModule,
      roleId: searchRoleId,
    });
  };

  const handleModuleChange = (value: string) => {
    setSearchModule(value === "all" ? "" : value);
    fetchDataWithFilters({
      name: searchName,
      module: value === "all" ? "" : value,
      roleId: searchRoleId,
    });
  };

  const handleRoleChange = (value: string) => {
    setSearchRoleId(value === "all" ? "" : value);
    fetchDataWithFilters({
      name: searchName,
      module: searchModule,
      roleId: value === "all" ? "" : value,
    });
  };

  const handleClearFilters = () => {
    setSearchName("");
    setSearchModule("all");
    setSearchRoleId("all");
    fetchDataWithFilters({ name: "", module: "", roleId: "" });
  };

  const confirmDelete = (permission: Permission) => {
    setSelectedPermission(permission);
    setIsDeleteDialogOpen(true);
  };

  const handleDeletePermission = async () => {
    if (!selectedPermission) return;

    try {
      await PermissionService.deletePermission(selectedPermission.id);
      // Recarrega a lista com os filtros atuais
      fetchDataWithFilters(
        {
          name: searchName,
          module: searchModule,
          roleId: searchRoleId,
        },
        false
      ); // Passa false para não resetar a página após exclusão (opcional)
    } catch (error: any) {
      // O toast já é exibido pelo service ou pelo hook de paginação
    } finally {
      setIsDeleteDialogOpen(false);
      setSelectedPermission(null);
    }
  };

  // Função que será passada para a tabela para recarregar os dados após uma edição
  const handlePermissionUpdated = () => {
    fetchDataWithFilters(
      {
        name: searchName,
        module: searchModule === "all" ? "" : searchModule,
        roleId: searchRoleId === "all" ? "" : searchRoleId,
      },
      false // Não reseta a página ao atualizar
    );
  };

  const pageSize = 5;
  const minHeightPerItem = 48;
  const calculatedMinHeight = pageSize * minHeightPerItem;

  return (
    <>
      <div className="mb-6 p-4 rounded-lg border bg-card shadow-sm grid grid-cols-1 gap-6 md:grid-cols-4">
        <div className="grid gap-2">
          <Label
            htmlFor="search-permission"
            className="text-sm font-medium text-foreground"
          >
            Buscar por Permissão
          </Label>
          <Input
            id="search-permission"
            value={searchName}
            onChange={(e) => setSearchName(e.target.value)}
            placeholder="Nome da permissão (ex: users.view)"
            className="border-input focus:ring-primary focus:border-primary"
            onKeyDown={(e) => {
              if (e.key === "Enter") {
                handleSearch();
              }
            }}
          />
        </div>
        <div className="grid gap-2">
          <Label
            htmlFor="search-module"
            className="text-sm font-medium text-foreground"
          >
            Filtrar por Módulo
          </Label>
          <Select value={searchModule} onValueChange={handleModuleChange}>
            <SelectTrigger id="search-module">
              <SelectValue placeholder="Selecione um módulo" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">Todos os Módulos</SelectItem>
              {availableModules.map((module) => (
                <SelectItem key={module} value={module}>
                  {module}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
        <div className="grid gap-2">
          <Label
            htmlFor="search-role"
            className="text-sm font-medium text-foreground"
          >
            Filtrar por Perfil
          </Label>
          <Select value={searchRoleId} onValueChange={handleRoleChange}>
            <SelectTrigger id="search-role">
              <SelectValue placeholder="Selecione um perfil" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">Todos os Perfis</SelectItem>
              {availableRoles.map((role) => (
                <SelectItem key={role.id} value={role.id.toString()}>
                  {role.name}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
        <div className="flex items-end gap-2">
          <Button onClick={handleSearch} className="mb-0.5">
            Buscar
          </Button>
          <Button
            onClick={handleClearFilters}
            variant="outline"
            className="mb-0.5"
          >
            Limpar Filtros
          </Button>
        </div>
      </div>

      {isLoading ? (
        <div
          className="flex flex-col items-center justify-center text-muted-foreground border rounded-md"
          style={{ minHeight: `${calculatedMinHeight + 120}px` }}
        >
          <Loader2 className="h-8 w-8 animate-spin text-primary mb-3" />
          <p>Consultando servidor...</p>
        </div>
      ) : (
        <PermissionsTable
          permissions={permissions} // Renomeado de 'data' para 'permissions'
          onConfirmDelete={confirmDelete}
          currentPage={currentPage}
          totalPages={totalPages}
          totalItems={totalItems}
          onPreviousPage={goToPreviousPage}
          onNextPage={goToNextPage}
          onGoToFirstPage={goToFirstPage}
          onGoToLastPage={goToLastPage}
          onPermissionUpdated={handlePermissionUpdated}
        />
      )}

      {isDeleteDialogOpen && (
        <ConfirmDeleteDialog
          open={isDeleteDialogOpen}
          onOpenChange={setIsDeleteDialogOpen}
          itemName={
            selectedPermission?.permission
              ? `a permissão "${selectedPermission.permission}"`
              : "o item selecionado"
          }
          onConfirm={handleDeletePermission}
          title="Confirma a exclusão da Permissão?"
          description={`Tem certeza que deseja excluir a permissão "${selectedPermission?.permission}"? Esta ação não pode ser desfeita.`}
        />
      )}
    </>
  );
}
