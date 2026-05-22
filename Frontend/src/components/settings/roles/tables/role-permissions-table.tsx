// src/app/(system)/settings/roles/components/tables/role-permissions-table.tsx
"use client";

import { ColumnDef, RowSelectionState } from "@tanstack/react-table";
import { Permission } from "@/services/auth/permission-schema";
import { DataTable } from "@/components/common/table/data-table";
import { Icon } from "@iconify/react";
import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import { useToast } from "@/components/ui/use-toast";
import { ConfirmDeleteDialog } from "@/components/common/confirm-delete-dialog";
import { useEffect, useState } from "react";
import { useRequirePermission } from "@/hooks/use-require-permission";
import { Loader2 } from "lucide-react";
import RoleService from "@/services/auth/role-service";
import { useHasPermission } from "@/hooks/use-has-permission";

interface RolePermissionsTableProps {
  permissions: Permission[];
  roleId: string; // ID do perfil ao qual as permissões estão sendo associadas
  showLinkedPermissions: boolean; // Indica se estamos mostrando permissões vinculadas ou não
  onPermissionsChanged: () => void; // Callback para recarregar a lista de permissões na RoleList

  // Props de paginação do useServerPagination
  currentPage: number;
  totalPages: number;
  totalItems: number;
  onPreviousPage: () => void;
  onNextPage: () => void;
  onGoToFirstPage: () => void;
  onGoToLastPage: () => void;
}

export function RolePermissionsTable({
  permissions,
  roleId,
  showLinkedPermissions,
  onPermissionsChanged,
  currentPage,
  totalPages,
  totalItems,
  onPreviousPage,
  onNextPage,
  onGoToFirstPage,
  onGoToLastPage,
}: RolePermissionsTableProps) {
  const { toast } = useToast();
  const hasEditPermission = useHasPermission("roles.edit");

  // Estado para armazenar o ID da permissão a ser desvinculada (para o modal de confirmação)
  const [permissionToUnlink, setPermissionToUnlink] =
    useState<Permission | null>(null);
  const [isConfirmUnlinkOpen, setIsConfirmUnlinkOpen] = useState(false);
  const [isLinkingUnlinking, setIsLinkingUnlinking] = useState(false);

  // Estado para gerenciar as permissões selecionadas para VINCULAR
  // Guarda os IDs das permissões que estão DESVINCULADAS mas que o usuário MARCOU para vincular
  const [selectedPermissionsToLink, setSelectedPermissionsToLink] = useState<
    Set<string>
  >(new Set());

  // Estado para gerenciar a seleção de linha do TanStack Table
  // IMPORTANTE: Esse estado é necessário para o checkbox mestre e para o .getIsSelected() funcionar.
  const [rowSelection, setRowSelection] = useState<RowSelectionState>({});

  // Resetar seleções e inicializar `rowSelection`
  useEffect(() => {
    setSelectedPermissionsToLink(new Set()); // Sempre limpa as seleções para "vincular"

    // NOVO: Inicializa o `rowSelection` do TanStack Table.
    // Se estiver mostrando vinculadas, todas as linhas devem vir "selecionadas" (para poder desmarcar/desvincular).
    // Se estiver mostrando não vinculadas, começam desmarcadas.
    const initialTanStackSelection: RowSelectionState = {};
    if (showLinkedPermissions) {
      permissions.forEach((p) => {
        initialTanStackSelection[p.id] = true;
      });
    }
    setRowSelection(initialTanStackSelection);
  }, [roleId, showLinkedPermissions, permissions]); // `permissions` é importante aqui

  // Função para lidar com a mudança do checkbox individual
  const handleCheckboxChange = async (
    permission: Permission,
    isChecked: boolean
  ) => {
    if (showLinkedPermissions) {
      // Estamos no modo "Mostrar Permissões Vinculadas"
      // Se desmarcar, significa DESVINCULAR
      if (!isChecked) {
        setPermissionToUnlink(permission); // Guarda a permissão para a confirmação
        setIsConfirmUnlinkOpen(true); // Abre o modal de confirmação
      }
      // Se marcar no modo "vinculadas", não faz nada (elas já estão vinculadas)
      // Ou seja, o checkbox de uma permissão vinculada é só para desvincular
    } else {
      // Estamos no modo "Mostrar Permissões NÃO Vinculadas"
      // Se marcar, significa que o usuário quer VINCULAR essa permissão
      if (isChecked) {
        setSelectedPermissionsToLink(
          (prev) => new Set(prev.add(permission.id))
        );
      } else {
        // Se desmarcar, significa que o usuário mudou de ideia e não quer mais vincular
        setSelectedPermissionsToLink((prev) => {
          const newSet = new Set(prev);
          newSet.delete(permission.id);
          return newSet;
        });
      }
    }
  };

  // Função para confirmar e executar a desvinculação
  const confirmUnlinkPermission = async () => {
    if (!permissionToUnlink || !roleId) return;

    setIsLinkingUnlinking(true);
    try {
      await RoleService.unlinkPermissionFromRole(roleId, permissionToUnlink.id);
      toast({
        title: "Sucesso!",
        description: `Permissão "${permissionToUnlink.permission}" desvinculada do perfil.`,
        variant: "default",
      });
      onPermissionsChanged(); // Recarrega a tabela de permissões na RoleList
    } catch (error: any) {
      console.error("Erro ao desvincular permissão:", error);
      toast({
        title: "Erro",
        description:
          error.message || "Ocorreu um erro ao desvincular a permissão.",
        variant: "destructive",
      });
    } finally {
      setIsConfirmUnlinkOpen(false);
      setPermissionToUnlink(null);
      setIsLinkingUnlinking(false);
    }
  };

  // Função para vincular as permissões selecionadas
  const handleLinkSelectedPermissions = async () => {
    if (selectedPermissionsToLink.size === 0 || !roleId) return;

    setIsLinkingUnlinking(true);
    try {
      const permissionIdsArray = Array.from(selectedPermissionsToLink);
      await RoleService.linkPermissionsToRole(roleId, permissionIdsArray);
      toast({
        title: "Sucesso!",
        description: `${permissionIdsArray.length} permissões vinculadas ao perfil.`,
        variant: "default",
      });
      setSelectedPermissionsToLink(new Set()); // Limpa as seleções
      onPermissionsChanged(); // Recarrega a tabela
    } catch (error: any) {
      console.error("Erro ao vincular permissões:", error);
      toast({
        title: "Erro",
        description:
          error.message || "Ocorreu um erro ao vincular as permissões.",
        variant: "destructive",
      });
    } finally {
      setIsLinkingUnlinking(false);
    }
  };

  // Definição das colunas da tabela
  const columns: ColumnDef<Permission>[] = [
    {
      id: "select", // ID para a coluna de seleção (checkbox)
      header: ({ table }) => (
        <Checkbox
          // O checked do header precisa refletir o estado de seleção apropriado
          checked={
            showLinkedPermissions
              ? table.getIsAllPageRowsSelected() ||
                (table.getIsSomePageRowsSelected() && "indeterminate")
              : selectedPermissionsToLink.size > 0 &&
                selectedPermissionsToLink.size ===
                  table.getRowModel().rows.length
              ? true
              : selectedPermissionsToLink.size > 0 &&
                selectedPermissionsToLink.size < table.getRowModel().rows.length
              ? "indeterminate"
              : false
          }
          onCheckedChange={(value) => {
            // No modo "vinculadas", usa o toggleAllPageRowsSelected do TanStack Table para desmarcar tudo.
            if (showLinkedPermissions) {
              table.toggleAllPageRowsSelected(!!value);
              // Dispara a desvinculação em massa se desmarcar (opcional, mas consistente com a UX)
              // Você pode implementar uma lógica aqui para um modal de "desvincular todos desta página"
              // ou deixar que o usuário desvincule um por um. Por enquanto, só desmarca visualmente.
            } else {
              // No modo "NÃO vinculadas", gerencia nosso `selectedPermissionsToLink`
              const newSelections = new Set(selectedPermissionsToLink);
              table.getRowModel().rows.forEach((row) => {
                if (!!value) {
                  // Se o checkbox mestre foi MARCADO
                  newSelections.add(row.original.id);
                } else {
                  // Se o checkbox mestre foi DESMARCADO
                  newSelections.delete(row.original.id);
                }
              });
              setSelectedPermissionsToLink(newSelections);
              // Também atualiza a seleção da linha do TanStack Table para visual
              table.toggleAllPageRowsSelected(!!value);
            }
          }}
          aria-label="Selecionar tudo"
          // Desabilita o mestre no modo vinculado se não tiver permissão ou estiver processando,
          // ou se não houver linhas para selecionar/desmarcar
          disabled={
            !hasEditPermission ||
            isLinkingUnlinking ||
            (showLinkedPermissions && table.getRowModel().rows.length === 0)
          }
        />
      ),
      cell: ({ row }) => {
        // Verifica se a permissão já está vinculada
        // Sua lógica simplificada permanece aqui, perfeita!
        const isChecked = showLinkedPermissions
          ? showLinkedPermissions
          : selectedPermissionsToLink.has(row.original.id);
        const isDisabled = !hasEditPermission || isLinkingUnlinking;

        return (
          <Checkbox
            checked={isChecked}
            onCheckedChange={(value) => {
              // Sempre que o checkbox é clicado, atualiza a seleção da linha no TanStack Table.
              // Isso é fundamental para o checkbox mestre e para a API da tabela funcionarem.
              row.toggleSelected(!!value);
              handleCheckboxChange(row.original, !!value);
            }}
            aria-label="Selecionar linha"
            disabled={isDisabled}
          />
        );
      },
      enableSorting: false,
      enableHiding: false,
    },
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
  ];

  // Verifica se o botão "Vincular Permissões Selecionadas" deve ser visível
  const shouldShowLinkButton =
    selectedPermissionsToLink.size > 0 &&
    hasEditPermission &&
    !isLinkingUnlinking &&
    !showLinkedPermissions;

  return (
    <>
      {shouldShowLinkButton && (
        <div className="flex justify-end mb-4">
          <Button
            onClick={handleLinkSelectedPermissions}
            disabled={isLinkingUnlinking}
          >
            {isLinkingUnlinking ? (
              <>
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                Vinculando...
              </>
            ) : (
              `Vincular ${selectedPermissionsToLink.size} Permissão(ões) Selecionada(s)`
            )}
          </Button>
        </div>
      )}

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

      {/* Diálogo de confirmação para desvincular */}
      {isConfirmUnlinkOpen && permissionToUnlink && (
        <ConfirmDeleteDialog
          open={isConfirmUnlinkOpen}
          onOpenChange={setIsConfirmUnlinkOpen}
          itemName={
            permissionToUnlink.permission
              ? `a permissão "${permissionToUnlink.permission}"`
              : "o item selecionado"
          }
          onConfirm={confirmUnlinkPermission}
          title="Desvincular Permissão?"
          description={`Tem certeza que deseja desvincular a permissão "${permissionToUnlink.permission}" do perfil selecionado? Esta ação não pode ser desfeita facilmente.`}
          confirmButtonText="Desvincular"
        />
      )}
    </>
  );
}
