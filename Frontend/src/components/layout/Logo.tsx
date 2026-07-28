import { BRAND } from "@/theme/antd-theme";

interface LogoProps {
  size?: number;
  showWordmark?: boolean;
}

// Losango com gradiente em CSS puro, sem asset externo — cores da identidade
// visual (Documentos/Imagens/Base da identidade visual.jpeg). Se um dia
// existir um export vetorial do ícone real (barras + anel), trocar aqui.
export function Logo({ size = 26, showWordmark = true }: LogoProps) {
  return (
    <div style={{ display: "flex", alignItems: "center", gap: 10 }}>
      <div
        style={{
          width: size,
          height: size,
          flexShrink: 0,
          background: BRAND.gradient,
          transform: "rotate(45deg)",
          borderRadius: 6,
        }}
      />
      {showWordmark && (
        <span style={{ fontSize: 19, fontWeight: 600, whiteSpace: "nowrap" }}>
          Lumini{" "}
          <span
            style={{
              backgroundImage: BRAND.gradient,
              backgroundClip: "text",
              WebkitBackgroundClip: "text",
              color: "transparent",
            }}
          >
            hub
          </span>
        </span>
      )}
    </div>
  );
}
