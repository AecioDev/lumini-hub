import { useEffect, useState } from "react";
import { Descriptions, Drawer, Skeleton, Space, Tag } from "antd";
import { userService } from "@/services/users/user-service";
import { companyService } from "@/services/companies/company-service";
import { useFeedback } from "@/hooks/useFeedback";
import type { ApiUserDetail } from "@/types/auth";
import { getTagColor } from "@/utils/avatar";

interface UserDetailDrawerProps {
  userId: number | null;
  onClose: () => void;
}

export function UserDetailDrawer({ userId, onClose }: UserDetailDrawerProps) {
  const [user, setUser] = useState<ApiUserDetail | null>(null);
  const [companyName, setCompanyName] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const feedback = useFeedback();

  useEffect(() => {
    if (userId === null) {
      setUser(null);
      setCompanyName(null);
      return;
    }
    setLoading(true);
    userService
      .getById(userId)
      .then(async (result) => {
        setUser(result);
        if (result.company_id) {
          const company = await companyService.getById(result.company_id);
          setCompanyName(company.trade_name || company.legal_name);
        } else {
          setCompanyName(null);
        }
      })
      .catch(() => feedback.error("Erro ao carregar detalhes do usuário."))
      .finally(() => setLoading(false));
  }, [userId, feedback]);

  const permissionsByModule = (user?.permissions ?? []).reduce<Record<string, string[]>>(
    (acc, perm) => {
      (acc[perm.module] ??= []).push(perm.permission);
      return acc;
    },
    {}
  );

  return (
    <Drawer title="Detalhes do Usuário" open={userId !== null} onClose={onClose} width={420}>
      {loading || !user ? (
        <Skeleton active paragraph={{ rows: 6 }} />
      ) : (
        <>
          <Descriptions column={1} size="small" bordered>
            <Descriptions.Item label="Usuário">{user.username}</Descriptions.Item>
            <Descriptions.Item label="Nome">{user.name}</Descriptions.Item>
            <Descriptions.Item label="E-mail">{user.email || "—"}</Descriptions.Item>
            <Descriptions.Item label="Telefone">{user.phone || "—"}</Descriptions.Item>
            <Descriptions.Item label="Perfil">{user.role?.name}</Descriptions.Item>
            <Descriptions.Item label="Empresa">
              {companyName ?? "Nenhuma (usuário master)"}
            </Descriptions.Item>
            <Descriptions.Item label="Status">
              <Tag color={user.is_active ? "success" : "default"}>
                {user.is_active ? "Ativo" : "Inativo"}
              </Tag>
            </Descriptions.Item>
            <Descriptions.Item label="Último login">
              {user.last_login || "Nunca acessou"}
            </Descriptions.Item>
            <Descriptions.Item label="Criado em">{user.created_at}</Descriptions.Item>
          </Descriptions>

          <div style={{ marginTop: 20 }}>
            <div style={{ fontSize: 13, fontWeight: 600, marginBottom: 10 }}>Permissões</div>
            {Object.keys(permissionsByModule).length === 0 && (
              <span style={{ fontSize: 13, opacity: 0.6 }}>Nenhuma permissão atribuída.</span>
            )}
            {Object.entries(permissionsByModule).map(([module, permissions]) => (
              <div key={module} style={{ marginBottom: 10 }}>
                <div style={{ fontSize: 12, opacity: 0.6, marginBottom: 4 }}>{module}</div>
                <Space size={[4, 4]} wrap>
                  {permissions.map((permission) => (
                    <Tag key={permission} color={getTagColor(module)}>
                      {permission}
                    </Tag>
                  ))}
                </Space>
              </div>
            ))}
          </div>
        </>
      )}
    </Drawer>
  );
}
