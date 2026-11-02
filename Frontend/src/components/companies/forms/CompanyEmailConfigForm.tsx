import { useEffect, useState } from "react";
import { Alert, Button, Popconfirm, Skeleton } from "antd";
import { SmtpConfigForm } from "@/components/settings/SmtpConfigForm";
import { useAuth } from "@/contexts/AuthContext";
import { useFeedback } from "@/hooks/useFeedback";
import { companyEmailConfigService } from "@/services/companies/company-email-config-service";
import { getApiErrorMessage } from "@/utils/api-error";
import type { ApiCompanyEmailConfig } from "@/types/email-config";

interface CompanyEmailConfigFormProps {
  companyId: number;
}

// Texto de "de onde sai o e-mail desta empresa hoje", conforme a herança empresa -> Matriz -> sistema.
function effectiveNotice(config: ApiCompanyEmailConfig) {
  switch (config.effective_source) {
    case "company":
      return { type: "success" as const, text: "Os e-mails desta empresa saem pelo servidor configurado abaixo." };
    case "parent":
      return { type: "info" as const, text: `Esta empresa não usa servidor próprio: herda o da Matriz (${config.effective_name}).` };
    case "system":
      return { type: "info" as const, text: "Esta empresa não usa servidor próprio: usa o e-mail padrão do sistema." };
    default:
      return {
        type: "warning" as const,
        text: "Nenhum servidor de e-mail ativo para esta empresa, nem na Matriz, nem no sistema. Os e-mails não serão enviados.",
      };
  }
}

// Aba "Configuração de E-mail" da Empresa. Sem servidor próprio a empresa herda da Matriz e, por fim, do
// sistema; configurar aqui é só para quem quer outro remetente.
export function CompanyEmailConfigForm({ companyId }: CompanyEmailConfigFormProps) {
  const { hasPermission } = useAuth();
  const feedback = useFeedback();
  const [config, setConfig] = useState<ApiCompanyEmailConfig | null>(null);

  useEffect(() => {
    let cancelled = false;
    companyEmailConfigService
      .getByCompanyId(companyId)
      .then((result) => {
        if (!cancelled) setConfig(result);
      })
      .catch(() => feedback.error("Erro ao carregar a configuração de e-mail."));
    return () => {
      cancelled = true;
    };
  }, [companyId, feedback]);

  if (!config) return <Skeleton active paragraph={{ rows: 8 }} />;

  const notice = effectiveNotice(config);

  const handleRemove = async () => {
    try {
      await companyEmailConfigService.remove(companyId);
      setConfig(await companyEmailConfigService.getByCompanyId(companyId));
      feedback.success("Servidor próprio removido: a empresa voltou a herdar.");
    } catch (error) {
      feedback.error(getApiErrorMessage(error, "Erro ao remover a configuração de e-mail."));
    }
  };

  return (
    <SmtpConfigForm
      title="Servidor SMTP da empresa"
      settings={config}
      canEdit={hasPermission("companies.email_config.edit")}
      notice={<Alert type={notice.type} showIcon message={notice.text} style={{ marginBottom: 16 }} />}
      extraActions={
        config.configured && (
          <Popconfirm
            title="Remover o servidor próprio?"
            description="A empresa volta a usar o da Matriz ou o do sistema."
            okText="Remover"
            cancelText="Cancelar"
            onConfirm={() => void handleRemove()}
          >
            <Button danger>Remover servidor próprio</Button>
          </Popconfirm>
        )
      }
      onSave={async (values) => {
        // O endereço do sistema é do e-mail do sistema, não da empresa.
        const smtp = { ...values };
        delete smtp.app_base_url;
        setConfig(await companyEmailConfigService.save(companyId, smtp));
      }}
      onSendTest={(to) => companyEmailConfigService.sendTest(companyId, to)}
    />
  );
}
