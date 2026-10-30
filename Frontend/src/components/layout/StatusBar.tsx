import { useEffect, useState } from "react";
import { Tag, theme } from "antd";
import { Logo } from "@/components/layout/Logo";
import { useAuth } from "@/contexts/AuthContext";
import { useThemeMode } from "@/contexts/ThemeContext";
import { authService } from "@/services/auth/auth-service";
import { CHROME_BG } from "@/theme/antd-theme";
import type { ApiVersionInfo } from "@/types/version";

// Portal onde ficam as notas de cada versão (uma página por versão, ex.: .../releases/2026.2.2.0/).
// Pode ser trocado por VITE_RELEASES_URL no .env.local.
const RELEASES_URL = (
  import.meta.env.VITE_RELEASES_URL ?? "https://lumini-release.vercel.app/releases"
).replace(/\/+$/, "");

// Nome e cor do ambiente (APP_ENV). Vazio ou desconhecido não mostra etiqueta.
const ENVIRONMENTS: Record<string, { label: string; color: string }> = {
  development: { label: "Desenvolvimento", color: "blue" },
  staging: { label: "Teste", color: "orange" },
  production: { label: "Produção", color: "green" },
};

// Barra de status fixa no rodapé do app: versão do sistema, versão do banco, empresa ativa (só leitura:
// a troca é no pé do menu lateral) e ambiente. Outras informações do sistema entram aqui com o tempo. Se a consulta falhar, a barra
// continua no lugar, vazia, para o layout não pular.
export function StatusBar() {
  const { token } = theme.useToken();
  const { mode } = useThemeMode();
  const { user } = useAuth();
  const companyName = user?.active_company_name;
  const [info, setInfo] = useState<ApiVersionInfo | null>(null);

  useEffect(() => {
    let active = true;
    authService
      .versionInfo()
      .then((data) => {
        if (active) setInfo(data);
      })
      .catch(() => {});
    return () => {
      active = false;
    };
  }, []);

  const environment = info ? ENVIRONMENTS[info.environment] : undefined;
  const databaseBehind = info !== null && info.database_version < info.required_database_version;

  return (
    <div
      style={{
        flexShrink: 0,
        height: 28,
        padding: "0 16px",
        display: "flex",
        alignItems: "center",
        justifyContent: "space-between",
        gap: 12,
        fontSize: 12,
        // Uma linha só, em qualquer largura: o que não couber é cortado com "…" (a empresa primeiro).
        whiteSpace: "nowrap",
        color: token.colorTextSecondary,
        // Mesma cor do menu lateral e do cabeçalho.
        background: CHROME_BG[mode],
        borderTop: `1px solid ${token.colorBorderSecondary}`,
      }}
    >
      <div style={{ display: "flex", alignItems: "center", gap: 12, minWidth: 0 }}>
        <Logo size={9} showWordmark={false} />
        {info && (
          <span style={{ flexShrink: 0 }}>
            <a
              href={`${RELEASES_URL}/${encodeURIComponent(info.system_version)}/`}
              target="_blank"
              rel="noopener noreferrer"
              title="Abrir as notas desta versão"
              style={{
                color: "inherit",
                textDecoration: "underline dotted",
                textUnderlineOffset: 3,
              }}
            >
              Versão {info.system_version}
            </a>{" "}
            - BD v{info.database_version}
          </span>
        )}
        {databaseBehind && (
          <Tag color="warning" style={{ margin: 0, flexShrink: 0 }}>
            Banco desatualizado (exige v{info.required_database_version})
          </Tag>
        )}
        {companyName && (
          <>
            {info && (
              <span aria-hidden="true" style={{ flexShrink: 0 }}>
                |
              </span>
            )}
            <span style={{ minWidth: 0, overflow: "hidden", textOverflow: "ellipsis" }}>
              Empresa: {companyName}
            </span>
          </>
        )}
      </div>

      {environment && (
        <Tag color={environment.color} style={{ margin: 0, flexShrink: 0 }}>
          {environment.label}
        </Tag>
      )}
    </div>
  );
}
