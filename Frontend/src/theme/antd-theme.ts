import { theme as antdTheme, type ThemeConfig } from "antd";
import type { ThemeMode } from "@/contexts/ThemeContext";

// Identidade visual Lumini Hub (Documentos/Imagens/Base da identidade visual.jpeg):
// Azul Royal (confiança/tecnologia), Azul Ciano (inovação) e Roxo (inteligência/
// transformação digital) sobre uma base grafite/branco.
export const BRAND = {
  royal: "#2563EB",
  cyan: "#06B6D4",
  purple: "#7C3AED",
  gradient: "linear-gradient(135deg, #2563EB 0%, #06B6D4 55%, #7C3AED 100%)",
} as const;

const SHARED_TOKENS: ThemeConfig["token"] = {
  colorPrimary: BRAND.royal,
  colorInfo: BRAND.cyan,
  colorLink: BRAND.royal,
  colorSuccess: "#52c41a",
  colorWarning: "#faad14",
  colorError: "#ff4d4f",
  borderRadius: 8,
  fontFamily:
    "-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif",
};

// Fora do sistema de tokens do antd: cor de fundo "chrome" da aplicação
// (usada no full-bleed das telas de auth e no Sider/Header) — grafite bem
// escuro (quase navy) no dark, igual ao fundo do brand board, branco no light.
export const CHROME_BG: Record<ThemeMode, string> = {
  dark: "#0A0E1A",
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
      colorBgLayout: isDark ? "#0F172A" : "#F5F7FA",
      colorBgContainer: isDark ? "#1E293B" : "#ffffff",
      colorBorder: isDark ? "#263449" : "#E5E9F0",
      colorBorderSecondary: isDark ? "#263449" : "#E5E9F0",
    },
    components: {
      Card: { borderRadiusLG: 12 },
      Layout: {
        siderBg: CHROME_BG[mode],
        headerBg: CHROME_BG[mode],
      },
      Menu: {
        // Hover/seleção usam tons do Azul Royal (marca), não mais cinza neutro.
        itemHoverBg: isDark ? "rgba(37,99,235,0.14)" : "rgba(37,99,235,0.06)",
        itemSelectedBg: isDark ? "rgba(37,99,235,0.20)" : "#EFF4FF",
        itemSelectedColor: BRAND.royal,
      },
    },
  };
}
