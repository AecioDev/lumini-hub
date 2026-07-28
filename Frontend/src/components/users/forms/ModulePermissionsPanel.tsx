import { Checkbox, Typography } from "antd";
import type { ApiPermission, ApiPermissionsByModule } from "@/types/permission";

const { Text } = Typography;

const ACTION_LABELS: Record<string, string> = {
  view: "Ver",
  create: "Criar",
  edit: "Editar",
  delete: "Excluir",
};

const MODULE_LABELS: Record<string, string> = {
  users: "Usuários",
  roles: "Perfis",
  permissions: "Permissões",
  customers: "Clientes",
  suppliers: "Fornecedores",
  admin: "Administração",
};

const CANONICAL_ACTION_ORDER = ["view", "create", "edit", "delete"];

export function moduleLabel(module: string): string {
  return MODULE_LABELS[module] ?? module;
}

// Fallback pra quando o nome do Perfil não bate com nenhum módulo do
// catálogo (ex.: "Gerente"/"ADMIN" não são módulos) — usa o primeiro módulo
// que tiver alguma permissão selecionada por aquele preset.
export function pickDefaultModule(
  permissionsByModule: ApiPermissionsByModule[],
  selectedIds: Set<number>
): string | null {
  const withSelection = permissionsByModule.find((entry) =>
    entry.permissions.some((perm) => selectedIds.has(perm.id))
  );
  return withSelection?.module ?? permissionsByModule[0]?.module ?? null;
}

// O painel de permissões é controlado pelos próprios botões de Perfil (card
// "Perfil", à esquerda) — não existe uma lista de módulos separada. Ao clicar
// num Perfil, mostramos o módulo cujo nome bate com o nome do Perfil (ex.:
// perfil "Vendas" -> módulo "Vendas"); quando não bate (ex.: "Gerente",
// "ADMIN" não são módulos — isso é uma organização a ajustar no catálogo de
// perfis/permissões do banco), cai no fallback acima.
export function pickModuleForRole(
  roleName: string,
  permissionsByModule: ApiPermissionsByModule[],
  selectedIds: Set<number>
): string | null {
  const normalized = roleName.trim().toLowerCase();
  const exactMatch = permissionsByModule.find(
    (entry) => entry.module.toLowerCase() === normalized
  );
  if (exactMatch) return exactMatch.module;
  return pickDefaultModule(permissionsByModule, selectedIds);
}

function extractAction(code: string): string {
  const lastDot = code.lastIndexOf(".");
  return lastDot >= 0 ? code.slice(lastDot + 1) : code;
}

export interface PermissionItem {
  permission: ApiPermission;
  label: string;
}

// Ver/Criar/Editar/Excluir usam o rótulo em PT-BR; qualquer outra ação (ou
// uma que colida — ex.: o módulo "dashboard" tem 5 permissões diferentes
// terminando em ".view", uma por perfil de dashboard) usa a própria
// descrição da permissão. No fim é só uma lista só, sem distinção visual.
export function buildPermissionItems(permissions: ApiPermission[]): PermissionItem[] {
  const grouped = new Map<string, ApiPermission[]>();
  for (const perm of permissions) {
    const action = extractAction(perm.permission);
    const bucket = grouped.get(action);
    if (bucket) bucket.push(perm);
    else grouped.set(action, [perm]);
  }

  const items: PermissionItem[] = [];

  for (const action of CANONICAL_ACTION_ORDER) {
    const bucket = grouped.get(action);
    if (bucket && bucket.length === 1) {
      items.push({ permission: bucket[0], label: ACTION_LABELS[action] });
      grouped.delete(action);
    }
  }

  for (const bucket of grouped.values()) {
    for (const perm of bucket) {
      items.push({ permission: perm, label: perm.description || perm.permission });
    }
  }

  return items;
}

interface ModulePermissionsPanelProps {
  permissionsByModule: ApiPermissionsByModule[];
  selectedIds: Set<number>;
  onChange: (next: Set<number>) => void;
  activeModule: string | null;
}

export function ModulePermissionsPanel({
  permissionsByModule,
  selectedIds,
  onChange,
  activeModule,
}: ModulePermissionsPanelProps) {
  const toggle = (permissionId: number) => {
    const next = new Set(selectedIds);
    if (next.has(permissionId)) next.delete(permissionId);
    else next.add(permissionId);
    onChange(next);
  };

  const current = permissionsByModule.find((entry) => entry.module === activeModule);
  const items = current ? buildPermissionItems(current.permissions) : [];

  if (!current) {
    return <Text type="secondary">Selecione um Perfil ao lado para ver as permissões.</Text>;
  }

  return (
    <div>
      <Text strong style={{ fontSize: 13, display: "block", marginBottom: 14 }}>
        {moduleLabel(current.module)}
      </Text>

      {items.length === 0 ? (
        <Text type="secondary" style={{ fontSize: 13 }}>
          Este módulo não tem permissões cadastradas.
        </Text>
      ) : (
        <div style={{ display: "flex", flexWrap: "wrap", gap: "12px 28px" }}>
          {items.map(({ permission, label }) => (
            <Checkbox
              key={permission.id}
              checked={selectedIds.has(permission.id)}
              onChange={() => toggle(permission.id)}
            >
              {label}
            </Checkbox>
          ))}
        </div>
      )}
    </div>
  );
}
