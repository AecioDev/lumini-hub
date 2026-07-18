"use client";

import { useState } from "react";
import { ColumnDef, RowSelectionState } from "@tanstack/react-table";
import { DataTable } from "@/components/common/table/data-table";
import { Loader2 } from "lucide-react";
import { useServerPagination } from "@/hooks/use-server-pagination";
import WebhookEventService from "@/services/integrations/webhook-event-service";
import { WebhookEvent } from "@/services/integrations/webhook-event-schema";

const columns: ColumnDef<WebhookEvent>[] = [
  { accessorKey: "source", header: "Origem" },
  { accessorKey: "event_type", header: "Tipo de Evento" },
  { accessorKey: "status", header: "Status" },
  {
    accessorKey: "received_at",
    header: "Recebido em",
    cell: ({ row }) => new Date(row.original.received_at).toLocaleString("pt-BR"),
  },
  {
    accessorKey: "processed_at",
    header: "Processado em",
    cell: ({ row }) =>
      row.original.processed_at
        ? new Date(row.original.processed_at).toLocaleString("pt-BR")
        : "—",
  },
  { accessorKey: "error_message", header: "Erro" },
];

export function WebhookEventsTable() {
  const [rowSelection, setRowSelection] = useState<RowSelectionState>({});

  const {
    data: webhookEvents,
    isLoading,
    currentPage,
    totalPages,
    totalItems,
    goToNextPage,
    goToPreviousPage,
    goToFirstPage,
    goToLastPage,
  } = useServerPagination<WebhookEvent>({
    fetchDataService: WebhookEventService.list,
    pageSize: 10,
  });

  if (isLoading) {
    return (
      <div className="flex flex-col items-center justify-center h-48 border rounded-md">
        <Loader2 className="h-8 w-8 animate-spin text-primary mb-3" />
        <p className="text-muted-foreground">Carregando eventos de webhook...</p>
      </div>
    );
  }

  return (
    <DataTable
      columns={columns}
      data={webhookEvents}
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
