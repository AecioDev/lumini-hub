import { App as AntdApp, ConfigProvider } from "antd";
import ptBR from "antd/locale/pt_BR";
import { RouterProvider } from "react-router-dom";
import { AuthProvider } from "@/contexts/AuthContext";
import { ThemeProvider, useThemeMode } from "@/contexts/ThemeContext";
import { buildAntdTheme } from "@/theme/antd-theme";
import { router } from "@/routes/router";

function ThemedApp() {
  const { mode } = useThemeMode();

  return (
    <ConfigProvider theme={buildAntdTheme(mode)} locale={ptBR}>
      <AntdApp notification={{ placement: "topRight" }}>
        <AuthProvider>
          <RouterProvider router={router} future={{ v7_startTransition: true }} />
        </AuthProvider>
      </AntdApp>
    </ConfigProvider>
  );
}

function App() {
  return (
    <ThemeProvider>
      <ThemedApp />
    </ThemeProvider>
  );
}

export default App;
