// src/components/settings/permissions-table.tsx
"use client";

import { useRouter } from "next/navigation";
import { ColumnDef, RowSelectionState } from "@tanstack/react-table";
import { Permission } from "@/services/auth/permission-schema";
import { routes } from "@/config/routes";
import { DataTable } from "@/components/common/table/data-table";
import { Icon } from "@iconify/react";
import { Button } from "@/components/ui/button";
import {
  DataTableRowActions,
  RowAction,
} from "@/components/common/table/data-table-row-action";

// IMPORTANTE: Importar o EditPermissionModal e o useState
import { useState } from "react";
import { EditPermissionModal } from "../dialogs/edit-permission-modal";
import { ViewPermissionModal } from "../dialogs/view-permission-modal";

interface PermissionsTableProps {
  permissions: Permission[];
  onConfirmDelete: (permission: Permission) => void;
  currentPage: number;
  totalPages: number;
  totalItems: number;
  onPreviousPage: () => void;
  onNextPage: () => void;
  onGoToFirstPage: () => void;
  onGoToLastPage: () => void;
  // Nova prop: Função para recarregar os dados na lista pai após uma edição
  onPermissionUpdated: () => void;
}

export function PermissionsTable({
  permissions,
  onConfirmDelete,
  currentPage,
  totalPages,
  totalItems,
  onPreviousPage,
  onNextPage,
  onGoToFirstPage,
  onGoToLastPage,
  onPermissionUpdated, // Nova prop
}: PermissionsTableProps) {
  const router = useRouter();

  // Estados para o modal de edição
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);
  const [editingPermissionId, setEditingPermissionId] = useState<string | null>(
    null
  );
  const [isViewModalOpen, setIsViewModalOpen] = useState(false);
  const [viewingPermissionId, setViewingPermissionId] = useState<string | null>(
    null
  );
  // Estado para gerenciar a seleção de linha do TanStack Table
  const [rowSelection, setRowSelection] = useState<RowSelectionState>({});

  const viewPermissionDetails = (id: string) => {
    setViewingPermissionId(id);
    setIsViewModalOpen(true);
  };

  // Alteração aqui: Abrir o modal em vez de navegar para outra página
  const handleEditPermission = (id: string) => {
    setEditingPermissionId(id);
    setIsEditModalOpen(true);
  };

  // 1. Defina as colunas para a sua tabela de permissões
  const columns: ColumnDef<Permission>[] = [
    {
      accessorKey: "module",
      header: ({ column }) => (
        <Button
          variant="ghost"
          onClick={column.getToggleSortingHandler()}
          className="p-0 h-auto"
        >
          Módulo
          <Icon
            icon={
              column.getIsSorted() === "asc" ? "mdi:arrow-up" : "mdi:arrow-down"
            }
            className={
              column.getIsSorted() ? "ml-1 h-3 w-3" : "ml-1 h-3 w-3 opacity-0"
            }
          />
        </Button>
      ),
      cell: ({ row }) => row.original.module,
    },
    {
      accessorKey: "permission",
      header: ({ column }) => (
        <Button
          variant="ghost"
          onClick={column.getToggleSortingHandler()}
          className="p-0 h-auto"
        >
          Permissão
          <Icon
            icon={
              column.getIsSorted() === "asc" ? "mdi:arrow-up" : "mdi:arrow-down"
            }
            className={
              column.getIsSorted() ? "ml-1 h-3 w-3" : "ml-1 h-3 w-3 opacity-0"
            }
          />
        </Button>
      ),
      cell: ({ row }) => row.original.permission,
    },
    {
      accessorKey: "description",
      header: ({ column }) => (
        <Button
          variant="ghost"
          onClick={column.getToggleSortingHandler()}
          className="p-0 h-auto"
        >
          Descrição
          <Icon
            icon={
              column.getIsSorted() === "asc" ? "mdi:arrow-up" : "mdi:arrow-down"
            }
            className={
              column.getIsSorted() ? "ml-1 h-3 w-3" : "ml-1 h-3 w-3 opacity-0"
            }
          />
        </Button>
      ),
      cell: ({ row }) => row.original.description,
    },
    {
      id: "actions",
      header: "Ações",
      enableHiding: false,
      cell: ({ row }) => {
        const perm = row.original;

        const actions: RowAction[] = [
          {
            label: "Visualizar",
            onClick: () => viewPermissionDetails(perm.id),
            icon: "mdi:eye",
            permission: "permissions.view",
          },
          {
            label: "Editar",
            onClick: () => handleEditPermission(perm.id),
            icon: "mdi:pencil",
            permission: "permissions.edit",
          },
          {
            label: "Excluir",
            onClick: () => onConfirmDelete(perm),
            icon: "mdi:trash-can",
            variant: "destructive",
            permission: "permissions.delete",
          },
        ];

        return <DataTableRowActions actions={actions} />;
      },
    },
  ];

  return (
    <>
      <DataTable
        columns={columns}
        data={permissions}
        currentPage={currentPage}
        totalPages={totalPages}
        totalItems={totalItems}
        onPreviousPage={onPreviousPage}
        onNextPage={onNextPage}
        onGoToFirstPage={onGoToFirstPage}
        onGoToLastPage={onGoToLastPage}
        rowSelection={rowSelection}
        onRowSelectionChange={setRowSelection}
      />

      {isEditModalOpen && editingPermissionId && (
        <EditPermissionModal
          open={isEditModalOpen}
          onOpenChange={(open) => {
            setIsEditModalOpen(open);
            if (!open) {
              setEditingPermissionId(null); // Limpa o ID quando o modal fecha
            }
          }}
          permissionId={editingPermissionId}
          onSuccess={onPermissionUpdated} // <--- CHAMA O CALLBACK DA LISTA PAI PARA RECARREGAR
        />
      )}

      {isViewModalOpen && viewingPermissionId && (
        <ViewPermissionModal
          open={isViewModalOpen}
          onOpenChange={(open) => {
            setIsViewModalOpen(open);
            if (!open) {
              setViewingPermissionId(null); // Limpa o ID quando o modal fecha
            }
          }}
          permissionId={viewingPermissionId}
        />
      )}
    </>
  );
}
