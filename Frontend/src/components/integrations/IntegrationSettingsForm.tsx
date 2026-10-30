import { useState } from "react";
import { zodResolver } from "@hookform/resolvers/zod";
import { Controller, useForm } from "react-hook-form";
import { Alert, Button, Card, Col, Input, Row } from "antd";
import { FormField } from "@/components/common/FormField";
import { integrationConfigService } from "@/services/integrations/integration-config-service";
import { useFeedback } from "@/hooks/useFeedback";
import {
  integrationSettingsSchema,
  type IntegrationSettingsFormValues,
} from "@/schemas/integration-settings-schema";
import type { ApiIntegrationSettings } from "@/types/integration";

interface IntegrationSettingsFormProps {
  settings: ApiIntegrationSettings;
  canEdit: boolean;
  onSaved: (settings: ApiIntegrationSettings) => void;
}

export function IntegrationSettingsForm({
  settings,
  canEdit,
  onSaved,
}: IntegrationSettingsFormProps) {
  const feedback = useFeedback();
  const [submitting, setSubmitting] = useState(false);

  const { control, handleSubmit } = useForm<IntegrationSettingsFormValues>({
    resolver: zodResolver(integrationSettingsSchema),
    defaultValues: {
      li_api_key: settings.li_api_key,
      li_app_key: settings.li_app_key,
      li_webhook_secret: settings.li_webhook_secret,
      codtipnot: settings.codtipnot,
    },
  });

  const onSubmit = async (values: IntegrationSettingsFormValues) => {
    setSubmitting(true);
    try {
      const updated = await integrationConfigService.updateSettings(values);
      onSaved(updated);
      feedback.success("Configurações salvas com sucesso.");
    } catch {
      feedback.error("Erro ao salvar configurações.");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} noValidate>
      {!canEdit && (
        <Alert
          type="info"
          showIcon
          style={{ marginBottom: 16 }}
          message="Você tem acesso só de leitura a estas configurações."
        />
      )}

      <Card title="Loja Integrada">
        <Row gutter={16}>
          <Col xs={24} sm={12}>
            <FormField label="API Key">
              <Controller
                name="li_api_key"
                control={control}
                render={({ field }) => (
                  <Input.Password
                    {...field}
                    disabled={!canEdit}
                    placeholder="Chave de API da Loja Integrada"
                  />
                )}
              />
            </FormField>
          </Col>
          <Col xs={24} sm={12}>
            <FormField label="App Key">
              <Controller
                name="li_app_key"
                control={control}
                render={({ field }) => (
                  <Input.Password
                    {...field}
                    disabled={!canEdit}
                    placeholder="Chave de aplicação da Loja Integrada"
                  />
                )}
              />
            </FormField>
          </Col>
          <Col xs={24} sm={12}>
            <FormField label="Segredo do Webhook">
              <Controller
                name="li_webhook_secret"
                control={control}
                render={({ field }) => (
                  <Input.Password
                    {...field}
                    disabled={!canEdit}
                    placeholder="Usado para validar webhooks recebidos"
                  />
                )}
              />
            </FormField>
          </Col>
          <Col xs={24} sm={12}>
            <FormField label="Código do tipo de nota (codtipnot)">
              <Controller
                name="codtipnot"
                control={control}
                render={({ field }) => (
                  <Input {...field} disabled={!canEdit} placeholder="Ex: 50" />
                )}
              />
            </FormField>
          </Col>
        </Row>
      </Card>

      {canEdit && (
        <div style={{ display: "flex", justifyContent: "flex-end", marginTop: 16 }}>
          <Button type="primary" htmlType="submit" loading={submitting}>
            Salvar Configurações
          </Button>
        </div>
      )}
    </form>
  );
}
