import type { ReactNode } from "react";
import { ConfigProvider } from "antd";
import { useAuth } from "@/contexts/AuthContext";
import { useThemeMode } from "@/contexts/ThemeContext";
import { buildCompanyThemeOverride } from "./antd-theme";

// Aplica a identidade visual (logo + paleta) da empresa ativa por cima do
// tema padrão da Lumini Hub, via um <ConfigProvider> aninhado — o antd
// mescla automaticamente com o tema do provider pai (algoritmo dark/light
// e os demais tokens de buildAntdTheme continuam valendo). Precisa estar
// dentro de <AuthProvider> (usa useAuth()). Sem CompanyVisualConfig
// cadastrada pra empresa ativa, não monta ConfigProvider nenhum — o
// fallback é simplesmente não sobrepor nada (paleta padrão intacta).
export function CompanyThemeProvider({ children }: { children: ReactNode }) {
  const { user } = useAuth();
  const { mode } = useThemeMode();

  const override = buildCompanyThemeOverride(user?.active_company_visual_config, mode);
  if (!override) return <>{children}</>;

  return <ConfigProvider theme={override}>{children}</ConfigProvider>;
}
