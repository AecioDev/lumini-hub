import { Badge } from "antd";
import { moduleLabel } from "@/components/users/forms/ModulePermissionsPanel";
import type { ApiPermissionsByModule } from "@/types/permission";

interface ModuleListProps {
  permissionsByModule: ApiPermissionsByModule[];
  selectedIds: Set<number>;
  activeModule: string | null;
  onSelectModule: (module: string) => void;
}

// Diferente do formulário de Usuário (onde os botões de Perfil já servem de
// navegador de módulo), aqui não existe atalho nenhum — é o próprio Perfil
// sendo editado, então precisa de uma lista de módulos dedicada.
export function ModuleList({
  permissionsByModule,
  selectedIds,
  activeModule,
  onSelectModule,
}: ModuleListProps) {
  return (
    <div style={{ display: "flex", flexDirection: "column", gap: 6 }}>
      {permissionsByModule.map((entry) => {
        const selectedCount = entry.permissions.filter((perm) =>
          selectedIds.has(perm.id)
        ).length;
        const isActive = entry.module === activeModule;

        return (
          <button
            key={entry.module}
            type="button"
            onClick={() => onSelectModule(entry.module)}
            style={{
              display: "flex",
              alignItems: "center",
              justifyContent: "space-between",
              textAlign: "left",
              padding: "8px 12px",
              borderRadius: 6,
              border: "1px solid",
              borderColor: isActive ? "#1677ff" : "transparent",
              background: isActive ? "rgba(22,119,255,0.12)" : "transparent",
              color: "inherit",
              cursor: "pointer",
              font: "inherit",
            }}
          >
            <span style={{ fontSize: 13, fontWeight: isActive ? 600 : 400 }}>
              {moduleLabel(entry.module)}
            </span>
            {selectedCount > 0 && <Badge count={selectedCount} size="small" />}
          </button>
        );
      })}
    </div>
  );
}
