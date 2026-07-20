import { Result } from "antd";
import { Outlet } from "react-router-dom";
import { useAuth } from "@/contexts/AuthContext";

interface RequirePermissionProps {
  permission: string;
}

// Gate de UI — a autorização real é sempre imposta pelo backend
// (middlewares.RequirePermission), isto só evita mostrar telas que a API
// vai recusar de qualquer forma.
export function RequirePermission({ permission }: RequirePermissionProps) {
  const { hasPermission } = useAuth();

  if (!hasPermission(permission)) {
    return (
      <Result
        status="403"
        title="Sem permissão"
        subTitle="Você não tem permissão para acessar esta página."
      />
    );
  }

  return <Outlet />;
}
