import { Button, Typography } from "antd";
import type { ApiRole } from "@/types/role";

const { Paragraph } = Typography;

interface RolePresetPickerProps {
  roles: ApiRole[];
  selectedRoleId: number | null;
  isCustom: boolean;
  onSelectRole: (roleId: number) => void;
}

// Perfil = preset/template: escolher um botão pré-preenche a matriz de
// permissões ao lado, mas cada checkbox continua editável individualmente
// depois (ver Documentos/exemplos/design_handoff_erp_frontend/README.md,
// seção "Permission / Role Model"). "Personalizado" não é uma Role real no
// banco — é só um rótulo de UI que aparece quando a matriz diverge do preset.
export function RolePresetPicker({
  roles,
  selectedRoleId,
  isCustom,
  onSelectRole,
}: RolePresetPickerProps) {
  return (
    <div>
      <Paragraph type="secondary" style={{ fontSize: 12, marginBottom: 14 }}>
        Um modelo que pré-preenche permissões — cada uma ainda pode ser ajustada
        individualmente ao lado.
      </Paragraph>
      <div style={{ display: "flex", flexDirection: "column", gap: 8 }}>
        {roles.map((role) => (
          <Button
            key={role.id}
            type={selectedRoleId === role.id && !isCustom ? "primary" : "default"}
            onClick={() => onSelectRole(role.id)}
            style={{ textAlign: "left" }}
          >
            {role.name}
          </Button>
        ))}
        {isCustom && (
          <Button type="primary" disabled style={{ textAlign: "left" }}>
            Personalizado
          </Button>
        )}
      </div>
    </div>
  );
}
