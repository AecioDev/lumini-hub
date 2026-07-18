"use client";

import { useState } from "react";
import { ColumnDef, RowSelectionState } from "@tanstack/react-table";
import { DataTable } from "@/components/common/table/data-table";
import { Loader2 } from "lucide-react";
import { useServerPagination } from "@/hooks/use-server-pagination";
import SyncLogService from "@/services/integrations/sync-log-service";
import { SyncLog } from "@/services/integrations/sync-log-schema";

const columns: ColumnDef<SyncLog>[] = [
  {
    accessorKey: "direction",
    header: "Direção",
    cell: ({ row }) =>
      row.original.direction === "li_to_sql"
        ? "Loja Integrada → SQL Server"
        : row.original.direction === "sql_to_li"
        ? "SQL Server → Loja Integrada"
        : row.original.direction,
  },
  { accessorKey: "entity_type", header: "Entidade" },
  { accessorKey: "reference_id", header: "Referência" },
  { accessorKey: "status", header: "Status" },
  { accessorKey: "message", header: "Mensagem" },
  {
    accessorKey: "created_at",
    header: "Data",
    cell: ({ row }) => new Date(row.original.created_at).toLocaleString("pt-BR"),
  },
];

export function SyncLogsTable() {
  const [rowSelection, setRowSelection] = useState<RowSelectionState>({});

  const {
    data: syncLogs,
    isLoading,
    currentPage,
    totalPages,
    totalItems,
    goToNextPage,
    goToPreviousPage,
    goToFirstPage,
    goToLastPage,
  } = useServerPagination<SyncLog>({
    fetchDataService: SyncLogService.list,
    pageSize: 10,
  });

  if (isLoading) {
    return (
      <div className="flex flex-col items-center justify-center h-48 border rounded-md">
        <Loader2 className="h-8 w-8 animate-spin text-primary mb-3" />
        <p className="text-muted-foreground">Carregando logs de sincronização...</p>
      </div>
    );
  }

  return (
    <DataTable
      columns={columns}
      data={syncLogs}
      currentPage={currentPage}
      totalPages={totalPages}
      totalItems={totalItems}
      onPreviousPage={goToPreviousPage}
      onNextPage={goToNextPage}
      onGoToFirstPage={goToFirstPage}
      onGoToLastPage={goToLastPage}
      rowSelection={rowSelection}
      onRowSelectionChange={setRowSelection}
    />
  );
}
