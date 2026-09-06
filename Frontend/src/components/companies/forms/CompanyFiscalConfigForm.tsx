import { useEffect, useState } from "react";
import { zodResolver } from "@hookform/resolvers/zod";
import { Controller, useForm } from "react-hook-form";
import {
  Alert,
  Button,
  Card,
  Col,
  Descriptions,
  Divider,
  Input,
  Row,
  Select,
  Skeleton,
  Space,
  Tag,
  Typography,
  Upload,
} from "antd";
import { UploadOutlined } from "@ant-design/icons";
import dayjs from "dayjs";
import { FormField } from "@/components/common/FormField";
import { companyFiscalConfigService } from "@/services/companies/company-fiscal-config-service";
import {
  companyFiscalConfigSchema,
  type CompanyFiscalConfigFormValues,
} from "@/schemas/company-fiscal-config-schema";
import { getApiErrorMessage } from "@/utils/api-error";
import { useFeedback } from "@/hooks/useFeedback";
import { maskDocument, onlyDigits } from "@/utils/masks";
import type { ApiCompanyFiscalConfigDetail } from "@/types/company-fiscal-config";

const { Text } = Typography;

const TAX_REGIME_OPTIONS = [
  { value: "SIMPLES", label: "Simples Nacional" },
  { value: "PRESUMIDO", label: "Lucro Presumido" },
  { value: "REAL", label: "Lucro Real" },
];

interface CompanyFiscalConfigFormProps {
  companyId: number;
  // Usado só pra comparação visual com o documento extraído do certificado
  // (ver certificateDocumentMatch) — nunca bloqueia o upload, só avisa.
  companyCnpj?: string;
}

// Alerta visual de vencimento do certificado — vencido (vermelho), perto de
// vencer em até 30 dias (amarelo), ou válido (verde). null quando não há
// data cadastrada ainda.
function certificateStatus(expiry: string | null): { color: string; text: string } | null {
  if (!expiry) return null;
  const days = dayjs(expiry).startOf("day").diff(dayjs().startOf("day"), "day");
  if (days < 0) return { color: "error", text: "Certificado vencido" };
  if (days <= 30) return { color: "warning", text: `Vence em ${days} dia${days === 1 ? "" : "s"}` };
  return { color: "success", text: "Válido" };
}

// Compara o documento extraído do certificado com o CNPJ da própria
// Company — só pra exibição (o usuário confere visualmente), nunca bloqueia
// o upload em caso de divergência.
function certificateDocumentMatch(
  certificateDocument: string,
  companyCnpj?: string
): { matches: boolean; text: string } | null {
  if (!certificateDocument || !companyCnpj) return null;
  const matches = onlyDigits(certificateDocument) === onlyDigits(companyCnpj);
  return {
    matches,
    text: matches ? "Confere com o CNPJ da empresa" : "Diverge do CNPJ cadastrado na empresa",
  };
}

// Upload de certificado (arquivo + senha) é state irmão, fora do
// react-hook-form — só existe depois que a Configuração Fiscal já foi
// criada (o certificado é anexado a um registro existente).
export function CompanyFiscalConfigForm({ companyId, companyCnpj }: CompanyFiscalConfigFormProps) {
  const feedback = useFeedback();

  const [config, setConfig] = useState<ApiCompanyFiscalConfigDetail | null>(null);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);

  const [certificateFile, setCertificateFile] = useState<File | null>(null);
  const [certificatePassword, setCertificatePassword] = useState("");
  const [uploadingCertificate, setUploadingCertificate] = useState(false);
  // Selecionar arquivo/senha ficam escondidos até o usuário clicar em
  // "Enviar Certificado" — evita mostrar um formulário de upload cada vez
  // que a tela abre, já com um certificado válido cadastrado.
  const [showUploadForm, setShowUploadForm] = useState(false);

  const {
    control,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<CompanyFiscalConfigFormValues>({
    resolver: zodResolver(companyFiscalConfigSchema),
    defaultValues: {
      tax_regime: "SIMPLES",
      accountant_name: "",
      accountant_document: "",
      accountant_contact: "",
    },
  });

  useEffect(() => {
    let cancelled = false;
    companyFiscalConfigService
      .getByCompanyId(companyId)
      .then((result) => {
        if (cancelled) return;
        setConfig(result);
        if (result) {
          reset({
            tax_regime: result.tax_regime,
            accountant_name: result.accountant_name,
            accountant_document: result.accountant_document,
            accountant_contact: result.accountant_contact,
          });
        }
      })
      .catch(() => feedback.error("Erro ao carregar a configuração fiscal."))
      .finally(() => {
        if (!cancelled) setLoading(false);
      });
    return () => {
      cancelled = true;
    };
  }, [companyId, feedback, reset]);

  const onSubmit = async (values: CompanyFiscalConfigFormValues) => {
    setSubmitting(true);
    try {
      if (config) {
        const updated = await companyFiscalConfigService.update(config.id, values);
        setConfig(updated);
        feedback.success("Configuração fiscal atualizada com sucesso.");
      } else {
        const created = await companyFiscalConfigService.create({
          company_id: companyId,
          ...values,
        });
        setConfig(created);
        feedback.success("Configuração fiscal criada com sucesso.");
      }
    } catch (error) {
      feedback.error(getApiErrorMessage(error, "Erro ao salvar configuração fiscal."));
    } finally {
      setSubmitting(false);
    }
  };

  const handleUploadCertificate = async () => {
    if (!config || !certificateFile || !certificatePassword) return;
    setUploadingCertificate(true);
    try {
      const updated = await companyFiscalConfigService.uploadCertificate(
        config.id,
        certificateFile,
        certificatePassword
      );
      setConfig(updated);
      setCertificateFile(null);
      setCertificatePassword("");
      setShowUploadForm(false);
      feedback.success("Certificado enviado com sucesso.");
    } catch (error) {
      feedback.error(getApiErrorMessage(error, "Erro ao enviar certificado."));
    } finally {
      setUploadingCertificate(false);
    }
  };

  if (loading) {
    return <Skeleton active paragraph={{ rows: 6 }} />;
  }

  const status = certificateStatus(config?.certificate_expiry ?? null);
  const documentMatch = certificateDocumentMatch(config?.certificate_subject_document ?? "", companyCnpj);

  return (
    <div style={{ maxWidth: 720, margin: "0 auto" }}>
      <form onSubmit={handleSubmit(onSubmit)} noValidate>
        <Card title="Configuração Fiscal">
          <Row gutter={16}>
            <Col xs={24} sm={12}>
              <FormField label="Tipo de Tributação" error={errors.tax_regime} stacked>
                <Controller
                  name="tax_regime"
                  control={control}
                  render={({ field }) => <Select {...field} options={TAX_REGIME_OPTIONS} />}
                />
              </FormField>
            </Col>

            <Col xs={24} sm={12}>
              <FormField label="Nome do Contador" error={errors.accountant_name} stacked>
                <Controller
                  name="accountant_name"
                  control={control}
                  render={({ field }) => <Input {...field} placeholder="Ex: Fulano de Tal" />}
                />
              </FormField>
            </Col>

            <Col xs={24} sm={12}>
              <FormField label="CPF/CRC do Contador" error={errors.accountant_document} stacked>
                <Controller
                  name="accountant_document"
                  control={control}
                  render={({ field }) => <Input {...field} placeholder="CPF ou CRC" />}
                />
              </FormField>
            </Col>

            <Col xs={24} sm={12}>
              <FormField label="Contato do Contador" error={errors.accountant_contact} stacked>
                <Controller
                  name="accountant_contact"
                  control={control}
                  render={({ field }) => <Input {...field} placeholder="Telefone ou e-mail" />}
                />
              </FormField>
            </Col>
          </Row>

          <Divider>Certificado Digital A1</Divider>

          {config ? (
            <Space direction="vertical" style={{ width: "100%" }} size="middle">
              <div
                style={{
                  display: "flex",
                  justifyContent: "space-between",
                  alignItems: "center",
                  flexWrap: "wrap",
                  gap: 8,
                }}
              >
                <Space>
                  <Text strong>Status:</Text>
                  {config.has_certificate ? (
                    <Tag color="blue">Certificado cadastrado</Tag>
                  ) : (
                    <Tag>Nenhum certificado enviado</Tag>
                  )}
                  {status && <Tag color={status.color}>{status.text}</Tag>}
                </Space>

                {!showUploadForm && (
                  <Button onClick={() => setShowUploadForm(true)}>
                    {config.has_certificate ? "Enviar Novo Certificado" : "Enviar Certificado"}
                  </Button>
                )}
              </div>

              {showUploadForm ? (
                <Space direction="vertical" style={{ width: "100%" }} size="small">
                  <Upload
                    beforeUpload={(file) => {
                      setCertificateFile(file);
                      return false;
                    }}
                    showUploadList={false}
                    accept=".pfx,.p12"
                  >
                    <Button icon={<UploadOutlined />}>Selecionar arquivo (.pfx)</Button>
                  </Upload>
                  {certificateFile && <Text type="secondary">{certificateFile.name}</Text>}

                  <Input.Password
                    placeholder="Senha do certificado"
                    value={certificatePassword}
                    onChange={(e) => setCertificatePassword(e.target.value)}
                  />

                  <Space>
                    <Button
                      type="primary"
                      onClick={() => void handleUploadCertificate()}
                      loading={uploadingCertificate}
                      disabled={!certificateFile || !certificatePassword}
                    >
                      Confirmar Envio
                    </Button>
                    <Button
                      onClick={() => {
                        setShowUploadForm(false);
                        setCertificateFile(null);
                        setCertificatePassword("");
                      }}
                    >
                      Cancelar
                    </Button>
                  </Space>
                </Space>
              ) : (
                config.has_certificate && (
                  <Descriptions column={1} size="small" bordered>
                    <Descriptions.Item label="Titular">
                      {config.certificate_subject_name || "—"}
                    </Descriptions.Item>
                    <Descriptions.Item label="Documento">
                      {config.certificate_subject_document
                        ? maskDocument(config.certificate_subject_document)
                        : "—"}
                    </Descriptions.Item>
                    <Descriptions.Item label="Validade">
                      {config.certificate_expiry
                        ? dayjs(config.certificate_expiry).format("DD/MM/YYYY")
                        : "—"}
                    </Descriptions.Item>
                  </Descriptions>
                )
              )}

              {documentMatch && !documentMatch.matches && (
                <Alert type="warning" showIcon message={documentMatch.text} />
              )}
            </Space>
          ) : (
            <Alert
              type="info"
              showIcon
              message="Salve a configuração fiscal primeiro pra poder enviar o certificado digital."
            />
          )}

          <div style={{ display: "flex", justifyContent: "flex-end", marginTop: 24 }}>
            <Button type="primary" htmlType="submit" loading={submitting}>
              {config ? "Salvar Configuração Fiscal" : "Criar Configuração Fiscal"}
            </Button>
          </div>
        </Card>
      </form>
    </div>
  );
}
