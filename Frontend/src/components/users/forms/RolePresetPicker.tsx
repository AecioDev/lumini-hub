import { CheckCircleFilled } from "@ant-design/icons";
import { Button, Typography } from "antd";
import type { ApiRole } from "@/types/role";

const { Paragraph } = Typography;

interface RolePresetPickerProps {
  roles: ApiRole[];
  selectedRoleId: number | null;
  onSelectRole: (roleId: number) => void;
  rolesWithPermissions?: Set<number>;
}

// Filtro de navegação: escolher um botão só troca o módulo de permissões
// exibido na matriz ao lado (ex.: ver as permissões do módulo Estoque), sem
// alterar a seleção já marcada — permissões de outros módulos continuam
// marcadas. NÃO altera o Perfil atribuído ao usuário (isso é feito
// exclusivamente pelo campo "Perfil" na aba "Dados do Usuário"). Serve para
// dar ao usuário acesso pontual a permissões de outros módulos além das do
// seu próprio Perfil, sem perder o que já foi marcado.
//
// `rolesWithPermissions` marca com um check os perfis cujo módulo tem
// alguma permissão atualmente selecionada no usuário — ajuda a ver de
// relance quais outros perfis (além do atribuído) contribuíram permissões.
export function RolePresetPicker({
  roles,
  selectedRoleId,
  onSelectRole,
  rolesWithPermissions,
}: RolePresetPickerProps) {
  return (
    <div>
      <Paragraph type="secondary" style={{ fontSize: 12, marginBottom: 14 }}>
        Filtro para navegar até as permissões de um módulo. As permissões já
        marcadas em outros módulos são mantidas. Não altera o Perfil do
        usuário — isso é feito na aba "Dados do Usuário".
      </Paragraph>
      <div style={{ display: "flex", flexDirection: "column", gap: 8 }}>
        {roles.map((role) => {
          const isActive = selectedRoleId === role.id;
          const hasPermissions = rolesWithPermissions?.has(role.id) ?? false;
          return (
            <Button
              key={role.id}
              type={isActive ? "primary" : "default"}
              onClick={() => onSelectRole(role.id)}
              style={{
                textAlign: "left",
                display: "flex",
                alignItems: "center",
                justifyContent: "space-between",
              }}
            >
              <span>{role.name}</span>
              {hasPermissions && (
                <CheckCircleFilled
                  style={{ color: isActive ? "#fff" : "#52c41a" }}
                />
              )}
            </Button>
          );
        })}
      </div>
    </div>
  );
}
