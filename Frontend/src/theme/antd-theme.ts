import { theme as antdTheme, type ThemeConfig } from "antd";
import type { ThemeMode } from "@/contexts/ThemeContext";

// Tokens replicados de Documentos/exemplos/design_handoff_erp_frontend/README.md
// ("Design Tokens"). colorPrimary/success/warning/error já batem com o default
// do antd — mantidos explícitos aqui só para documentar a intenção.
const SHARED_TOKENS: ThemeConfig["token"] = {
  colorPrimary: "#1677ff",
  colorSuccess: "#52c41a",
  colorWarning: "#faad14",
  colorError: "#ff4d4f",
  borderRadius: 6,
  fontFamily:
    "-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif",
};

// Fora do sistema de tokens do antd: cor de fundo "chrome" da aplicação
// (usada no full-bleed das telas de auth e no Sider/Header), igual ao
// sidebarBg/bgLayout do mock — mais escura que o colorBgLayout do conteúdo.
export const CHROME_BG: Record<ThemeMode, string> = {
  dark: "#000000",
  light: "#ffffff",
};

export function buildAntdTheme(mode: ThemeMode): ThemeConfig {
  const isDark = mode === "dark";

  return {
    algorithm: isDark ? antdTheme.darkAlgorithm : antdTheme.defaultAlgorithm,
    token: {
      ...SHARED_TOKENS,
      // colorBgLayout = fundo da área de conteúdo ("contentBg" no mock)
      // colorBgContainer = fundo dos Cards/Inputs ("cardBg" no mock)
      colorBgLayout: isDark ? "#141414" : "#f5f5f5",
      colorBgContainer: isDark ? "#1f1f1f" : "#ffffff",
      colorBorder: isDark ? "#303030" : "#f0f0f0",
      colorBorderSecondary: isDark ? "#303030" : "#f0f0f0",
    },
    components: {
      Card: { borderRadiusLG: 8 },
      Layout: {
        siderBg: CHROME_BG[mode],
        headerBg: CHROME_BG[mode],
      },
      Menu: {
        // Hover neutro (não confundir com o item ativo, que usa o azul
        // primário) — mesmos tons de destaque do mock (ERP.dc.html).
        itemHoverBg: isDark ? "rgba(255,255,255,0.08)" : "rgba(0,0,0,0.04)",
        itemSelectedBg: isDark ? "rgba(22,119,255,0.15)" : "#e6f4ff",
        itemSelectedColor: "#1677ff",
      },
    },
  };
}
