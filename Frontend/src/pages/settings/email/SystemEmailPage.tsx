import { useEffect, useState } from "react";
import { Alert, Skeleton } from "antd";
import { PageHeader } from "@/components/common/PageHeader";
import { SmtpConfigForm } from "@/components/settings/SmtpConfigForm";
import { useAuth } from "@/contexts/AuthContext";
import { useFeedback } from "@/hooks/useFeedback";
import { systemEmailService } from "@/services/settings/system-email-service";
import type { ApiSystemEmailSettings } from "@/types/email-config";

export function SystemEmailPage() {
  const { hasPermission } = useAuth();
  const feedback = useFeedback();
  const [settings, setSettings] = useState<ApiSystemEmailSettings | null>(null);

  useEffect(() => {
    let cancelled = false;
    systemEmailService
      .get()
      .then((result) => {
        if (!cancelled) setSettings(result);
      })
      .catch(() => feedback.error("Erro ao carregar a configuração de e-mail do sistema."));
    return () => {
      cancelled = true;
    };
  }, [feedback]);

  return (
    <div>
      <PageHeader
        title="E-mail do sistema"
        subtitle="Servidor de e-mail padrão: usado por qualquer empresa que não tenha um próprio"
      />

      {!settings ? (
        <Skeleton active paragraph={{ rows: 8 }} />
      ) : (
        <SmtpConfigForm
          title="Servidor SMTP padrão"
          settings={settings}
          canEdit={hasPermission("system_settings.edit")}
          notice={
            !settings.enabled && (
              <Alert
                type="warning"
                showIcon
                style={{ marginBottom: 16 }}
                message="O e-mail padrão do sistema está desativado. Recuperação de senha e demais avisos só funcionam em empresas com servidor próprio."
              />
            )
          }
          onSave={async (values) => setSettings(await systemEmailService.save(values))}
          onSendTest={(to) => systemEmailService.sendTest(to)}
        />
      )}
    </div>
  );
}
