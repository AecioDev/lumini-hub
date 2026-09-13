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
// (usada no full-bleed das telas de auth e no Sider/Header) — cinza neutro
// quase preto no dark (recalibrado 2026-09-13, tom extraído de um print real
// do Claude Desktop — antes era um grafite com tingimento azulado, decisão
// explícita do usuário de trocar por cinza puro), branco no light.
export const CHROME_BG: Record<ThemeMode, string> = {
  dark: "#111111",
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
      colorBgLayout: isDark ? "#151515" : "#F5F7FA",
      colorBgContainer: isDark ? "#1C1C1C" : "#ffffff",
      colorBorder: isDark ? "#2A2A2A" : "#E5E9F0",
      colorBorderSecondary: isDark ? "#2A2A2A" : "#E5E9F0",
    },
    components: {
      Card: { borderRadiusLG: 12 },
      Layout: {
        siderBg: CHROME_BG[mode],
        headerBg: CHROME_BG[mode],
      },
      Menu: {
        // Hover usa um azul translúcido; item selecionado usa o Azul Royal
        // sólido com texto branco — mesmo tratamento nos dois temas (pedido
        // explícito do usuário em 2026-09-13: o claro usava fundo pálido
        // (#EFF4FF) + texto azul, bem mais discreto que o escuro, que já
        // ficava sólido "de graça" só porque o antd ignora itemSelectedBg/
        // itemSelectedColor no preset dark e cai no próprio default sólido
        // dele — ver darkItemSelectedBg/darkItemSelectedColor abaixo, agora
        // explícitos em vez de acidentais).
        itemHoverBg: isDark ? "rgba(37,99,235,0.14)" : "rgba(37,99,235,0.06)",
        itemSelectedBg: BRAND.royal,
        itemSelectedColor: "#fff",
        // Sem isso, o preset dark do Menu (SidebarMenu.tsx usa theme="dark")
        // aplica fundos navy hardcoded — passavam despercebidos contra o
        // CHROME_BG antigo (também navy), mas destoam contra o cinza neutro
        // novo (CFG-6.3.1). Dois tokens dark-específicos diferentes, achados
        // via DevTools + confirmados lendo node_modules/antd/lib/menu/style:
        // darkItemBg (#001529) pinta o <ul> raiz do menu (e, por seletor
        // descendente, vaza pros <ul> de submenu que não tiverem override
        // mais específico); darkSubMenuItemBg (#000c17) é essa própria
        // sobreposição mais específica pro <ul> de um submenu aberto.
        darkItemBg: "transparent",
        darkSubMenuItemBg: "transparent",
        // Mesma pegadinha dos dois de cima: itemHoverBg (acima) nunca
        // funcionou no modo escuro — o preset dark ignora essa chave e usa
        // darkItemHoverBg (default do antd: 'transparent', por isso o hover
        // só mudava a cor do texto, não o fundo do item inteiro). Reaplica o
        // mesmo tom translúcido do Azul Royal já usado no modo claro.
        darkItemHoverBg: "rgba(37,99,235,0.14)",
        darkItemSelectedBg: BRAND.royal,
        darkItemSelectedColor: "#fff",
      },
    },
  };
}
