"use client";

import { useId, useMemo, useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Checkbox } from "@/components/ui/checkbox";
import { Loader2, Search, X } from "lucide-react";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Permission } from "@/services/auth/permission-schema";

type FlatPermission = Permission & { module: string };

function toId(id: Permission["id"]): number {
  return typeof id === "string" ? parseInt(id, 10) : id;
}

interface PermissionsPickerProps {
  permissionsByModule: Record<string, Permission[]>;
  isLoading?: boolean;
  selectedIds: Set<number>;
  onToggle: (id: number) => void;
  onBulkToggle: (ids: number[], checked: boolean) => void;
  selectedLabel?: string;
  /** Ex: um Select de "Filtrar por Perfil" nos diálogos de usuário. Renderizado na barra de filtros. */
  extraFilterControl?: React.ReactNode;
  extraFilterPredicate?: (permission: FlatPermission) => boolean;
  isExtraFilterActive?: boolean;
  onClearExtraFilter?: () => void;
}

export function PermissionsPicker({
  permissionsByModule,
  isLoading = false,
  selectedIds,
  onToggle,
  onBulkToggle,
  selectedLabel = "selecionada(s)",
  extraFilterControl,
  extraFilterPredicate,
  isExtraFilterActive = false,
  onClearExtraFilter,
}: PermissionsPickerProps) {
  const idPrefix = useId();
  const [searchQuery, setSearchQuery] = useState("");
  const [filterModule, setFilterModule] = useState("all");

  const modulesList = useMemo(
    () => Object.keys(permissionsByModule).sort(),
    [permissionsByModule]
  );

  const groupedFiltered = useMemo(() => {
    const q = searchQuery.trim().toLowerCase();

    const matches = (p: FlatPermission) => {
      if (q) {
        const nameMatch = p.permission.toLowerCase().includes(q);
        const descMatch = p.description?.toLowerCase().includes(q) || false;
        const modMatch = p.module.toLowerCase().includes(q);
        if (!nameMatch && !descMatch && !modMatch) return false;
      }
      if (filterModule !== "all" && p.module !== filterModule) return false;
      if (extraFilterPredicate && !extraFilterPredicate(p)) return false;
      return true;
    };

    return modulesList
      .map((moduleName) => ({
        module: moduleName,
        permissions: (permissionsByModule[moduleName] || [])
          .map((p): FlatPermission => ({ ...p, module: moduleName }))
          .filter(matches)
          .sort((a, b) => a.permission.localeCompare(b.permission)),
      }))
      .filter((group) => group.permissions.length > 0);
  }, [modulesList, permissionsByModule, searchQuery, filterModule, extraFilterPredicate]);

  const filteredFlat = useMemo(
    () => groupedFiltered.flatMap((g) => g.permissions),
    [groupedFiltered]
  );

  const isAllFilteredChecked = useMemo(() => {
    if (filteredFlat.length === 0) return false;
    return filteredFlat.every((p) => selectedIds.has(toId(p.id)));
  }, [filteredFlat, selectedIds]);

  const handleClearFilters = () => {
    setSearchQuery("");
    setFilterModule("all");
    onClearExtraFilter?.();
  };

  const hasActiveFilters =
    !!searchQuery || filterModule !== "all" || isExtraFilterActive;

  return (
    <div className="space-y-3">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-2">
        <p className="text-xs text-muted-foreground">
          Marque as permissões desejadas. Elas podem ser de qualquer módulo,
          independente do que já vem por padrão.
        </p>
        <span className="text-xs text-muted-foreground bg-accent px-2.5 py-1 rounded-full font-medium shrink-0 self-start sm:self-center">
          {selectedIds.size} {selectedLabel}
        </span>
      </div>

      {/* Barra de Filtros */}
      <div className="grid grid-cols-1 sm:grid-cols-12 gap-2 bg-muted/30 p-2.5 rounded-lg border">
        <div
          className={`relative ${
            extraFilterControl ? "sm:col-span-5" : "sm:col-span-8"
          }`}
        >
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
        {extraFilterControl && (
          <div className="sm:col-span-3">{extraFilterControl}</div>
        )}
        <div className="sm:col-span-1 flex justify-center items-center">
          <Button
            variant="ghost"
            size="icon"
            onClick={handleClearFilters}
            disabled={!hasActiveFilters}
            className="h-9 w-9 shrink-0"
            title="Limpar Filtros"
          >
            <X className="h-4 w-4" />
          </Button>
        </div>
      </div>

      {isLoading ? (
        <div className="flex flex-col items-center justify-center py-12 border rounded-lg bg-muted/10">
          <Loader2 className="h-7 w-7 animate-spin text-primary mb-2.5" />
          <p className="text-xs text-muted-foreground font-medium">
            Carregando permissões...
          </p>
        </div>
      ) : filteredFlat.length === 0 ? (
        <div className="flex flex-col items-center justify-center py-12 border rounded-lg bg-muted/10 text-center px-4">
          <p className="text-xs text-muted-foreground font-medium">
            Nenhuma permissão encontrada com os filtros aplicados.
          </p>
          {hasActiveFilters && (
            <Button
              variant="link"
              size="sm"
              onClick={handleClearFilters}
              className="text-xs mt-1 h-auto p-0"
            >
              Limpar Filtros
            </Button>
          )}
        </div>
      ) : (
        <div className="flex flex-col border rounded-lg overflow-hidden bg-card">
          <div className="relative flex items-center gap-2 px-4 py-2 border-b bg-muted/50 text-[10px] font-bold uppercase tracking-wider text-muted-foreground">
            <Checkbox
              id={`${idPrefix}-select-all-filtered`}
              checked={isAllFilteredChecked}
              onCheckedChange={(checked) =>
                onBulkToggle(
                  filteredFlat.map((p) => toId(p.id)),
                  !!checked
                )
              }
              className="h-4 w-4"
            />
            <Label
              htmlFor={`${idPrefix}-select-all-filtered`}
              className="cursor-pointer"
            >
              Selecionar todas as filtradas ({filteredFlat.length})
            </Label>
          </div>

          <div className="max-h-[480px] overflow-y-auto divide-y">
            {groupedFiltered.map((group) => (
              <div key={group.module}>
                <div className="sticky top-0 z-10 px-4 py-1.5 bg-muted/95 backdrop-blur-sm border-b border-t text-[10px] font-bold uppercase tracking-wider text-muted-foreground capitalize">
                  {group.module} ({group.permissions.length})
                </div>
                <div className="divide-y">
                  {group.permissions.map((permission) => {
                    const idNum = toId(permission.id);
                    const isChecked = selectedIds.has(idNum);

                    return (
                      <div
                        key={permission.id}
                        className="relative flex items-start gap-3 px-4 py-2 text-xs hover:bg-muted/10 transition-colors"
                      >
                        <Checkbox
                          id={`${idPrefix}-perm-${permission.id}`}
                          checked={isChecked}
                          onCheckedChange={() => onToggle(idNum)}
                          className="h-4 w-4 mt-0.5"
                        />
                        <Label
                          htmlFor={`${idPrefix}-perm-${permission.id}`}
                          className="font-semibold text-foreground cursor-pointer block truncate flex-1"
                          title={permission.description || permission.permission}
                        >
                          {permission.description || permission.permission}
                          {permission.description && (
                            <span className="block text-[10px] font-normal text-muted-foreground font-mono truncate">
                              {permission.permission}
                            </span>
                          )}
                        </Label>
                      </div>
                    );
                  })}
                </div>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
