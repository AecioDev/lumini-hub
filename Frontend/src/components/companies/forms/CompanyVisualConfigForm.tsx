import { useEffect, useState } from "react";
import { zodResolver } from "@hookform/resolvers/zod";
import { Controller, useForm } from "react-hook-form";
import {
  Alert,
  Button,
  Card,
  Col,
  ColorPicker,
  Divider,
  Row,
  Skeleton,
  Space,
  Typography,
  Upload,
} from "antd";
import { UploadOutlined } from "@ant-design/icons";
import { FormField } from "@/components/common/FormField";
import {
  companyVisualConfigService,
  companyVisualConfigLogoUrl,
} from "@/services/companies/company-visual-config-service";
import {
  companyVisualConfigSchema,
  type CompanyVisualConfigFormValues,
} from "@/schemas/company-visual-config-schema";
import { getApiErrorMessage } from "@/utils/api-error";
import { useFeedback } from "@/hooks/useFeedback";
import type { ApiCompanyVisualConfigDetail } from "@/types/company-visual-config";

const { Text } = Typography;

const DEFAULT_PRIMARY_COLOR = "#2563EB";

interface CompanyVisualConfigFormProps {
  companyId: number;
}

// Upload de logo é state irmão, fora do react-hook-form — só existe depois
// que a Configuração Visual já foi criada, mesmo padrão do certificado em
// CompanyFiscalConfigForm.tsx.
export function CompanyVisualConfigForm({ companyId }: CompanyVisualConfigFormProps) {
  const feedback = useFeedback();

  const [config, setConfig] = useState<ApiCompanyVisualConfigDetail | null>(null);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);

  const [logoFile, setLogoFile] = useState<File | null>(null);
  const [uploadingLogo, setUploadingLogo] = useState(false);
  const [logoVersion, setLogoVersion] = useState(0);

  const {
    control,
    handleSubmit,
    reset,
    setValue,
    watch,
    formState: { errors },
  } = useForm<CompanyVisualConfigFormValues>({
    resolver: zodResolver(companyVisualConfigSchema),
    defaultValues: {
      primary_color: DEFAULT_PRIMARY_COLOR,
      secondary_color: "",
      accent_color: "",
    },
  });

  useEffect(() => {
    let cancelled = false;
    companyVisualConfigService
      .getByCompanyId(companyId)
      .then((result) => {
        if (cancelled) return;
        setConfig(result);
        if (result) {
          reset({
            primary_color: result.primary_color,
            secondary_color: result.secondary_color ?? "",
            accent_color: result.accent_color ?? "",
          });
        }
      })
      .catch(() => feedback.error("Erro ao carregar a identidade visual."))
      .finally(() => {
        if (!cancelled) setLoading(false);
      });
    return () => {
      cancelled = true;
    };
  }, [companyId, feedback, reset]);

  const onSubmit = async (values: CompanyVisualConfigFormValues) => {
    setSubmitting(true);
    try {
      const payload = {
        primary_color: values.primary_color,
        secondary_color: values.secondary_color || undefined,
        accent_color: values.accent_color || undefined,
      };
      if (config) {
        const updated = await companyVisualConfigService.update(config.id, payload);
        setConfig(updated);
        feedback.success("Identidade visual atualizada com sucesso.");
      } else {
        const created = await companyVisualConfigService.create({
          company_id: companyId,
          ...payload,
        });
        setConfig(created);
        feedback.success("Identidade visual criada com sucesso.");
      }
    } catch (error) {
      feedback.error(getApiErrorMessage(error, "Erro ao salvar identidade visual."));
    } finally {
      setSubmitting(false);
    }
  };

  const handleUploadLogo = async () => {
    if (!config || !logoFile) return;
    setUploadingLogo(true);
    try {
      const updated = await companyVisualConfigService.uploadLogo(config.id, logoFile);
      setConfig(updated);
      setLogoFile(null);
      setLogoVersion((v) => v + 1);
      feedback.success("Logo enviado com sucesso.");
    } catch (error) {
      feedback.error(getApiErrorMessage(error, "Erro ao enviar logo."));
    } finally {
      setUploadingLogo(false);
    }
  };

  if (loading) {
    return <Skeleton active paragraph={{ rows: 6 }} />;
  }

  const secondaryColor = watch("secondary_color");
  const accentColor = watch("accent_color");

  return (
    <div style={{ maxWidth: 720, margin: "0 auto" }}>
      <form onSubmit={handleSubmit(onSubmit)} noValidate>
        <Card title="Identidade Visual">
          <Row gutter={16}>
            <Col xs={24} sm={8}>
              <FormField label="Cor Primária" error={errors.primary_color} required>
                <Controller
                  name="primary_color"
                  control={control}
                  render={({ field }) => (
                    <ColorPicker
                      value={field.value}
                      onChange={(color) => field.onChange(color.toHexString().toUpperCase())}
                      format="hex"
                      disabledAlpha
                      showText
                    />
                  )}
                />
              </FormField>
            </Col>

            <Col xs={24} sm={8}>
              <FormField label="Cor Secundária" error={errors.secondary_color}>
                <Space>
                  <Controller
                    name="secondary_color"
                    control={control}
                    render={({ field }) => (
                      <ColorPicker
                        value={field.value || undefined}
                        onChange={(color) => field.onChange(color.toHexString().toUpperCase())}
                        format="hex"
                        disabledAlpha
                        showText
                      />
                    )}
                  />
                  {secondaryColor && (
                    <Button type="link" size="small" onClick={() => setValue("secondary_color", "")}>
                      Limpar
                    </Button>
                  )}
                </Space>
              </FormField>
            </Col>

            <Col xs={24} sm={8}>
              <FormField label="Cor de Destaque" error={errors.accent_color}>
                <Space>
                  <Controller
                    name="accent_color"
                    control={control}
                    render={({ field }) => (
                      <ColorPicker
                        value={field.value || undefined}
                        onChange={(color) => field.onChange(color.toHexString().toUpperCase())}
                        format="hex"
                        disabledAlpha
                        showText
                      />
                    )}
                  />
                  {accentColor && (
                    <Button type="link" size="small" onClick={() => setValue("accent_color", "")}>
                      Limpar
                    </Button>
                  )}
                </Space>
              </FormField>
            </Col>
          </Row>

          <div style={{ display: "flex", justifyContent: "flex-end", marginTop: 8 }}>
            <Button type="primary" htmlType="submit" loading={submitting}>
              {config ? "Salvar Identidade Visual" : "Criar Identidade Visual"}
            </Button>
          </div>

          <Divider>Logo</Divider>

          {config ? (
            <Row gutter={16} align="middle">
              <Col xs={24} sm={12}>
                <Space direction="vertical" style={{ width: "100%" }} size="small">
                  {config.has_logo ? (
                    <img
                      key={logoVersion}
                      src={`${companyVisualConfigLogoUrl(config.id)}?v=${logoVersion}`}
                      alt="Logo da empresa"
                      style={{ maxHeight: 64, maxWidth: "100%", objectFit: "contain" }}
                    />
                  ) : (
                    <Text type="secondary">Nenhum logo enviado ainda.</Text>
                  )}
                </Space>
              </Col>

              <Col xs={24} sm={12}>
                <Space direction="vertical" style={{ width: "100%" }} size="small">
                  <Upload
                    beforeUpload={(file) => {
                      setLogoFile(file);
                      return false;
                    }}
                    showUploadList={false}
                    accept=".png,.jpg,.jpeg,.svg"
                  >
                    <Button icon={<UploadOutlined />}>Selecionar arquivo</Button>
                  </Upload>
                  {logoFile && <Text type="secondary">{logoFile.name}</Text>}
                  <Button
                    type="primary"
                    onClick={() => void handleUploadLogo()}
                    loading={uploadingLogo}
                    disabled={!logoFile}
                  >
                    Enviar Logo
                  </Button>
                </Space>
              </Col>
            </Row>
          ) : (
            <Alert
              type="info"
              showIcon
              message="Salve a identidade visual primeiro pra poder enviar o logo."
            />
          )}
        </Card>
      </form>
    </div>
  );
}
