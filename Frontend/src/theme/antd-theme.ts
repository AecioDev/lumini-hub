import { theme as antdTheme, type ThemeConfig } from "antd";
import type { ThemeMode } from "@/contexts/ThemeContext";
import type { ApiCompanyVisualConfigOption } from "@/types/auth";

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

function hexToRgb(hex: string): [number, number, number] {
  const clean = hex.replace("#", "");
  return [
    parseInt(clean.slice(0, 2), 16),
    parseInt(clean.slice(2, 4), 16),
    parseInt(clean.slice(4, 6), 16),
  ];
}

// Sobrepõe a identidade visual (logo + paleta) da empresa ativa por cima do
// tema padrão da Lumini Hub — usado num <ConfigProvider> aninhado (ver
// CompanyThemeProvider.tsx), que o antd mescla automaticamente com o tema
// do provider pai (algoritmo dark/light e demais tokens continuam vindo
// de buildAntdTheme). null quando a empresa não tem CompanyVisualConfig
// cadastrada — nesse caso o chamador nem monta o ConfigProvider extra,
// mantendo a paleta padrão sem diferença visual (ver plano_empresa.md §
// CompanyVisualConfig, "Fallback").
//
// accent_color e text_color não são mapeados aqui de propósito. accent_color
// não tem token antd correspondente em lugar nenhum do projeto hoje, só
// aparece no preview isolado de CompanyVisualConfigForm.tsx. text_color
// CHEGOU a ser mapeado pra colorTextLightSolid, mas foi revertido: esse é um
// token de alias global do antd consumido por Menu (cor do item NÃO
// selecionado no tema dark — ver darkItemColor em antd/es/menu/style/index.js),
// além de Button, Tag, Switch, Tooltip, Badge, Avatar etc. — não é um token
// "texto sobre botão primário" isolado como o nome sugere. Como o valor
// cadastrado é pensado pra contrastar com primary_color (ex.: quase-preto
// pra combinar com um dourado), sobrepor esse token globalmente deixou o
// texto do menu lateral quase ilegível (texto escuro sobre o navy escuro do
// Sider). Sem um token antd que isole "texto sobre a cor primária" sem
// vazar pro resto do app, text_color fica de fora do override global por
// enquanto — só aparece no preview isolado do form, igual accent_color.
export function buildCompanyThemeOverride(
  visualConfig: ApiCompanyVisualConfigOption | null | undefined,
  mode: ThemeMode
): ThemeConfig | null {
  if (!visualConfig) return null;

  const { primary_color, secondary_color } = visualConfig;
  const [r, g, b] = hexToRgb(primary_color);
  const isDark = mode === "dark";

  return {
    token: {
      colorPrimary: primary_color,
      colorLink: primary_color,
      ...(secondary_color ? { colorInfo: secondary_color } : {}),
    },
    components: {
      // itemSelectedColor/itemHoverBg/itemSelectedBg de buildAntdTheme são
      // fixos em BRAND.royal — sem essa sobreposição, o menu do sidebar
      // continuaria azul mesmo com uma cor primária diferente aplicada em
      // botões/links.
      Menu: {
        itemHoverBg: `rgba(${r}, ${g}, ${b}, ${isDark ? 0.14 : 0.06})`,
        itemSelectedBg: `rgba(${r}, ${g}, ${b}, ${isDark ? 0.2 : 0.1})`,
        itemSelectedColor: primary_color,
      },
    },
  };
}
