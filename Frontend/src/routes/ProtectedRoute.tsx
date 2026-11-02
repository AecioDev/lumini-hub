import { Spin } from "antd";
import { Navigate, Outlet, useLocation } from "react-router-dom";
import { useAuth } from "@/contexts/AuthContext";
import { CompanySelectionGate } from "@/components/companies/CompanySelectionGate";

// Tela de troca de senha obrigatória (ChangePasswordPage).
const CHANGE_PASSWORD_PATH = "/change-password";

export function ProtectedRoute() {
  const { isAuthenticated, isBootstrapping, user } = useAuth();
  const location = useLocation();

  if (isBootstrapping) {
    return (
      <div
        style={{
          minHeight: "100vh",
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
        }}
      >
        <Spin size="large" />
      </div>
    );
  }

  if (!isAuthenticated) {
    return <Navigate to="/login" replace state={{ from: location }} />;
  }

  // Troca de senha obrigatória (cadastro sem convite): nenhuma outra tela abre antes dela.
  const onChangePasswordPage = location.pathname === CHANGE_PASSWORD_PATH;
  if (user?.must_change_password && !onChangePasswordPage) {
    return <Navigate to={CHANGE_PASSWORD_PATH} replace />;
  }
  if (!user?.must_change_password && onChangePasswordPage) {
    return <Navigate to="/dashboard" replace />;
  }

  if (user?.requires_company_selection) {
    return <CompanySelectionGate />;
  }

  return <Outlet />;
}
