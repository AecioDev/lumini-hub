import { Card, Tag, Typography } from "antd";
import { buildPermissionItems, moduleLabel } from "./ModulePermissionsPanel";
import type { ApiPermissionsByModule } from "@/types/permission";

const { Text } = Typography;

interface UserPermissionsSummaryProps {
  permissionsByModule: ApiPermissionsByModule[];
  selectedIds: Set<number>;
}

// Resumo de todas as permissões marcadas para o usuário, agrupadas por
// módulo — dá pra ver de uma vez só o que foi acumulado ao navegar entre os
// filtros de perfil abaixo (que só trocam o módulo em exibição, sem
// substituir a seleção).
export function UserPermissionsSummary({
  permissionsByModule,
  selectedIds,
}: UserPermissionsSummaryProps) {
  const groups = permissionsByModule
    .map((entry) => ({
      module: entry.module,
      items: buildPermissionItems(entry.permissions).filter((item) =>
        selectedIds.has(item.permission.id)
      ),
    }))
    .filter((group) => group.items.length > 0);

  return (
    <Card title="Permissões do usuário" size="small" style={{ marginBottom: 16 }}>
      {groups.length === 0 ? (
        <Text type="secondary" style={{ fontSize: 13 }}>
          Nenhuma permissão selecionada.
        </Text>
      ) : (
        <div style={{ display: "flex", flexDirection: "column", gap: 6 }}>
          {groups.map((group) => (
            <div key={group.module}>
              <Text strong style={{ fontSize: 12, marginRight: 8 }}>
                {moduleLabel(group.module)}:
              </Text>
              {group.items.map((item) => (
                <Tag key={item.permission.id} style={{ marginBottom: 4 }}>
                  {item.label}
                </Tag>
              ))}
            </div>
          ))}
        </div>
      )}
    </Card>
  );
}
