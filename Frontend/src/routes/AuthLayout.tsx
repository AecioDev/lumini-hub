import { Card, theme } from "antd";
import { Outlet } from "react-router-dom";
import { CHROME_BG } from "@/theme/antd-theme";
import { useThemeMode } from "@/contexts/ThemeContext";
import { Logo } from "@/components/layout/Logo";

export function AuthLayout() {
  const { mode } = useThemeMode();
  const { token } = theme.useToken();

  return (
    <div
      style={{
        minHeight: "100vh",
        width: "100%",
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        background: CHROME_BG[mode],
      }}
    >
      <Card
        style={{
          width: 400,
          borderRadius: 12,
          boxShadow: "0 8px 32px rgba(0,0,0,0.2)",
        }}
        styles={{ body: { padding: 40 } }}
      >
        <div style={{ marginBottom: 32 }}>
          <Logo />
        </div>
        <Outlet context={{ token }} />
      </Card>
    </div>
  );
}
