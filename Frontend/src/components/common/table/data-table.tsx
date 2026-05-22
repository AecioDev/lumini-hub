// src/components/common/data-table/data-table.tsx
"use client";

import * as React from "react";
import {
  ColumnDef,
  flexRender,
  getCoreRowModel,
  useReactTable,
  getSortedRowModel,
  SortingState,
  // getFilteredRowModel, // Pode ser útil para filtros globais ou locais
  RowSelectionState,
} from "@tanstack/react-table";

import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { Icon } from "@iconify/react/dist/iconify.js";

interface DataTableProps<TData, TValue> {
  columns: ColumnDef<TData, TValue>[];
  data: TData[];
  currentPage: number;
  totalPages: number;
  totalItems: number;
  onPreviousPage: () => void;
  onNextPage: () => void;
  onGoToFirstPage: () => void;
  onGoToLastPage: () => void;
  pageSize?: number;
  rowSelection: RowSelectionState;
  onRowSelectionChange: (
    updater: React.SetStateAction<RowSelectionState>
  ) => void;
}

export function DataTable<TData, TValue>({
  columns,
  data,
  currentPage,
  totalPages,
  totalItems,
  onPreviousPage,
  onNextPage,
  onGoToFirstPage,
  onGoToLastPage,
  rowSelection,
  onRowSelectionChange,
}: DataTableProps<TData, TValue>) {
  const [sorting, setSorting] = React.useState<SortingState>([]);

  const table = useReactTable({
    data,
    columns,
    getCoreRowModel: getCoreRowModel(),
    onSortingChange: setSorting,
    getSortedRowModel: getSortedRowModel(),
    manualPagination: true,
    pageCount: totalPages,
    enableRowSelection: true,
    onRowSelectionChange: onRowSelectionChange,
    state: {
      sorting,
      pagination: {
        pageIndex: currentPage - 1,
        pageSize: data.length > 0 ? data.length : 10,
      },
      rowSelection,
    },
    // Se você for usar filtros globais na DataTable, adicione:
    // getFilteredRowModel: getFilteredRowModel(),
  });

  return (
    <div className="space-y-4">
      <div className="rounded-md border bg-card text-card-foreground shadow-sm overflow-hidden">
        <Table>
          <TableHeader className="bg-muted text-muted-foreground">
            {table.getHeaderGroups().map((headerGroup) => (
              <TableRow key={headerGroup.id}>
                {headerGroup.headers.map((header) => {
                  return (
                    <TableHead
                      key={header.id}
                      className="h-12 px-4 text-left align-middle font-semibold text-muted-foreground [&:has([role=checkbox])]:pr-0"
                    >
                      {header.isPlaceholder ? null : (
                        <div
                          {...{
                            className: header.column.getCanSort()
                              ? "cursor-pointer select-none flex items-center gap-1"
                              : "",
                            onClick: header.column.getToggleSortingHandler(),
                          }}
                        >
                          {flexRender(
                            header.column.columnDef.header,
                            header.getContext()
                          )}
                          {{
                            asc: <span className="ml-1">▲</span>,
                            desc: <span className="ml-1">▼</span>,
                          }[header.column.getIsSorted() as string] ?? null}
                        </div>
                      )}
                    </TableHead>
                  );
                })}
              </TableRow>
            ))}
          </TableHeader>
          <TableBody>
            {table.getRowModel().rows?.length ? (
              table.getRowModel().rows.map((row) => (
                <TableRow
                  key={row.id}
                  data-state={row.getIsSelected() && "selected"}
                  className="border-b transition-colors hover:bg-muted/50 data-[state=selected]:bg-muted"
                >
                  {row.getVisibleCells().map((cell) => (
                    <TableCell
                      key={cell.id}
                      className="p-4 align-middle [&:has([role=checkbox])]:pr-0"
                    >
                      {flexRender(
                        cell.column.columnDef.cell,
                        cell.getContext()
                      )}
                    </TableCell>
                  ))}
                </TableRow>
              ))
            ) : (
              <TableRow>
                <TableCell
                  colSpan={columns.length}
                  className="h-24 text-center"
                >
                  Nenhum resultado.
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>

      {/* Componentes de Paginação */}
      <div className="flex items-center justify-end space-x-3 py-4 text-sm text-muted-foreground">
        <span className="text-muted-foreground mr-auto">
          {totalItems} registros encontrados
        </span>

        <Button
          variant="outline"
          size="icon"
          onClick={onGoToFirstPage}
          disabled={currentPage === 1}
        >
          <Icon icon="mdi:page-first" className="h-4 w-4" />
        </Button>

        <Button
          variant="outline"
          size="icon"
          onClick={onPreviousPage}
          disabled={currentPage === 1}
        >
          <Icon icon="mdi:chevron-left" className="h-4 w-4" />
        </Button>

        <span className="text-muted-foreground">
          Página {currentPage} de {totalPages}
        </span>

        <Button
          variant="outline"
          size="icon"
          onClick={onNextPage}
          disabled={currentPage === totalPages}
        >
          <Icon icon="mdi:chevron-right" className="h-4 w-4" />
        </Button>

        <Button
          variant="outline"
          size="icon"
          onClick={onGoToLastPage}
          disabled={currentPage === totalPages}
        >
          <Icon icon="mdi:page-last" className="h-4 w-4" />
        </Button>
      </div>
    </div>
  );
}
