interface LogoProps {
  size?: number;
  showWordmark?: boolean;
}

// Losango com gradiente em CSS puro, sem asset externo — mesma ideia do mock
// (Documentos/exemplos/design_handoff_erp_frontend/ERP.dc.html).
export function Logo({ size = 26, showWordmark = true }: LogoProps) {
  return (
    <div style={{ display: "flex", alignItems: "center", gap: 10 }}>
      <div
        style={{
          width: size,
          height: size,
          flexShrink: 0,
          background: "linear-gradient(135deg, #1677ff, #69b1ff)",
          transform: "rotate(45deg)",
          borderRadius: 6,
        }}
      />
      {showWordmark && (
        <span style={{ fontSize: 19, fontWeight: 600, whiteSpace: "nowrap" }}>
          Lumini Hub
        </span>
      )}
    </div>
  );
}
