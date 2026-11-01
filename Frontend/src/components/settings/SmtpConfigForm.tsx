import { useEffect, useState, type ReactNode } from "react";
import { zodResolver } from "@hookform/resolvers/zod";
import { Controller, useForm } from "react-hook-form";
import { Alert, Button, Card, Col, Divider, Input, InputNumber, Row, Select, Space, Switch, Typography } from "antd";
import { FormField } from "@/components/common/FormField";
import { emailConfigSchema, testEmailSchema, type EmailConfigFormValues } from "@/schemas/email-config-schema";
import { getApiErrorMessage } from "@/utils/api-error";
import { useFeedback } from "@/hooks/useFeedback";
import type { ApiSmtpSettings, ApiTestEmailResult, SaveSmtpRequest, SmtpSecurity } from "@/types/email-config";

const { Text } = Typography;

const SECURITY_OPTIONS = [
  { value: "starttls", label: "STARTTLS (porta 587, o mais comum)" },
  { value: "ssl", label: "SSL/TLS (porta 465)" },
  { value: "none", label: "Sem criptografia (só servidor local de teste)" },
];

// Porta típica de cada modo: ao trocar a segurança, a porta acompanha se ainda for uma das típicas.
const DEFAULT_PORTS: Record<SmtpSecurity, number> = { starttls: 587, ssl: 465, none: 25 };
const TYPICAL_PORTS = Object.values(DEFAULT_PORTS);

const SOURCE_LABEL: Record<ApiTestEmailResult["source"], string> = {
  company: "da própria empresa",
  parent: "herdada da Matriz",
  system: "padrão do sistema",
};

function toFormValues(settings: ApiSmtpSettings): EmailConfigFormValues {
  return {
    host: settings.host,
    port: settings.port,
    security: settings.security,
    username: settings.username,
    password: "",
    from_name: settings.from_name,
    from_email: settings.from_email,
    enabled: settings.enabled,
  };
}

interface SmtpConfigFormProps {
  title: string;
  settings: ApiSmtpSettings;
  canEdit: boolean;
  // Quem usa grava e atualiza `settings` com a resposta (o formulário se refaz a partir dela).
  onSave: (values: SaveSmtpRequest) => Promise<void>;
  onSendTest: (to: string) => Promise<ApiTestEmailResult>;
  // Botões extras ao lado de "Salvar" (ex.: remover a configuração própria da empresa).
  extraActions?: ReactNode;
  // Aviso acima do formulário (ex.: de onde a empresa herda hoje).
  notice?: ReactNode;
}

// Formulário de servidor SMTP compartilhado pelo "E-mail do sistema" e pela aba "Configuração de E-mail"
// da Empresa. O e-mail de teste usa a configuração JÁ SALVA, por isso fica bloqueado com alterações pendentes.
export function SmtpConfigForm({ title, settings, canEdit, onSave, onSendTest, extraActions, notice }: SmtpConfigFormProps) {
  const feedback = useFeedback();
  const [saving, setSaving] = useState(false);
  const [testTo, setTestTo] = useState("");
  const [testing, setTesting] = useState(false);

  const {
    control,
    handleSubmit,
    reset,
    setValue,
    getValues,
    formState: { errors, isDirty },
  } = useForm<EmailConfigFormValues>({
    resolver: zodResolver(emailConfigSchema),
    defaultValues: toFormValues(settings),
  });

  useEffect(() => {
    reset(toFormValues(settings));
  }, [settings, reset]);

  const onSubmit = async (values: EmailConfigFormValues) => {
    setSaving(true);
    try {
      await onSave(values);
      feedback.success("Configuração de e-mail salva com sucesso.");
    } catch (error) {
      feedback.error(getApiErrorMessage(error, "Erro ao salvar a configuração de e-mail."));
    } finally {
      setSaving(false);
    }
  };

  const handleTest = async () => {
    const parsed = testEmailSchema.safeParse(testTo.trim());
    if (!parsed.success) {
      feedback.error(parsed.error.issues[0].message);
      return;
    }
    setTesting(true);
    try {
      const result = await onSendTest(parsed.data);
      feedback.success(`E-mail enviado a ${parsed.data} pela configuração ${SOURCE_LABEL[result.source]} (${result.name}).`);
    } catch (error) {
      feedback.error(getApiErrorMessage(error, "Erro ao enviar o e-mail de teste."), "Falha no envio");
    } finally {
      setTesting(false);
    }
  };

  return (
    <div style={{ maxWidth: 720, margin: "0 auto" }}>
      {notice}
      <form onSubmit={handleSubmit(onSubmit)} noValidate>
        <Card title={title}>
          <Row gutter={16}>
            <Col xs={24}>
              <FormField label="Envio ativo" error={errors.enabled} stacked>
                <Controller
                  name="enabled"
                  control={control}
                  render={({ field }) => (
                    <Switch checked={field.value} onChange={field.onChange} disabled={!canEdit} />
                  )}
                />
              </FormField>
            </Col>

            <Col xs={24} sm={16}>
              <FormField label="Servidor SMTP" error={errors.host} stacked>
                <Controller
                  name="host"
                  control={control}
                  render={({ field }) => <Input {...field} placeholder="smtp.gmail.com" disabled={!canEdit} />}
                />
              </FormField>
            </Col>

            <Col xs={24} sm={8}>
              <FormField label="Porta" error={errors.port} stacked>
                <Controller
                  name="port"
                  control={control}
                  render={({ field }) => (
                    <InputNumber
                      value={field.value}
                      onChange={(value) => field.onChange(value ?? 0)}
                      min={1}
                      max={65535}
                      style={{ width: "100%" }}
                      disabled={!canEdit}
                    />
                  )}
                />
              </FormField>
            </Col>

            <Col xs={24}>
              <FormField label="Segurança" error={errors.security} stacked>
                <Controller
                  name="security"
                  control={control}
                  render={({ field }) => (
                    <Select
                      value={field.value}
                      options={SECURITY_OPTIONS}
                      disabled={!canEdit}
                      onChange={(value: SmtpSecurity) => {
                        field.onChange(value);
                        if (TYPICAL_PORTS.includes(getValues("port"))) {
                          setValue("port", DEFAULT_PORTS[value], { shouldDirty: true });
                        }
                      }}
                    />
                  )}
                />
              </FormField>
            </Col>

            <Col xs={24} sm={12}>
              <FormField label="Usuário" error={errors.username} stacked>
                <Controller
                  name="username"
                  control={control}
                  render={({ field }) => (
                    <Input {...field} placeholder="Geralmente o e-mail" autoComplete="off" disabled={!canEdit} />
                  )}
                />
              </FormField>
            </Col>

            <Col xs={24} sm={12}>
              <FormField label="Senha" error={errors.password} stacked>
                <Controller
                  name="password"
                  control={control}
                  render={({ field }) => (
                    <Input.Password
                      {...field}
                      autoComplete="new-password"
                      disabled={!canEdit}
                      placeholder={settings.has_password ? "Gravada (deixe em branco para manter)" : "Senha ou senha de app"}
                    />
                  )}
                />
              </FormField>
            </Col>
          </Row>

          <Divider>Remetente</Divider>

          <Row gutter={16}>
            <Col xs={24} sm={12}>
              <FormField label="Nome" error={errors.from_name} stacked>
                <Controller
                  name="from_name"
                  control={control}
                  render={({ field }) => <Input {...field} placeholder="Ex: Lumini Hub" disabled={!canEdit} />}
                />
              </FormField>
            </Col>

            <Col xs={24} sm={12}>
              <FormField label="E-mail" error={errors.from_email} stacked>
                <Controller
                  name="from_email"
                  control={control}
                  render={({ field }) => <Input {...field} placeholder="noreply@empresa.com.br" disabled={!canEdit} />}
                />
              </FormField>
            </Col>
          </Row>

          <Text type="secondary">
            Gmail: servidor smtp.gmail.com, porta 587, STARTTLS, e uma senha de app (a senha normal da conta não funciona).
          </Text>

          {canEdit && (
            <div style={{ display: "flex", justifyContent: "flex-end", gap: 8, flexWrap: "wrap", marginTop: 24 }}>
              {extraActions}
              <Button type="primary" htmlType="submit" loading={saving}>
                Salvar
              </Button>
            </div>
          )}
        </Card>
      </form>

      {canEdit && (
        <Card title="E-mail de teste" style={{ marginTop: 16 }}>
          <Space direction="vertical" style={{ width: "100%" }}>
            {isDirty && (
              <Alert type="info" showIcon message="Salve as alterações antes de testar: o teste usa a configuração gravada." />
            )}
            <Space.Compact style={{ width: "100%" }}>
              <Input
                value={testTo}
                onChange={(e) => setTestTo(e.target.value)}
                placeholder="Enviar para (seu e-mail)"
                onPressEnter={() => void handleTest()}
              />
              <Button onClick={() => void handleTest()} loading={testing} disabled={isDirty || !testTo.trim()}>
                Enviar teste
              </Button>
            </Space.Compact>
          </Space>
        </Card>
      )}
    </div>
  );
}
