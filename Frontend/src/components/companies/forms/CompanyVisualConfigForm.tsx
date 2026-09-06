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
  Tag,
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

interface PresetPalette {
  name: string;
  primary: string;
  secondary: string;
  accent: string;
}

// Paletas prontas pra preencher os 3 campos de uma vez, sem precisar
// garimpar hex manualmente — pedido do usuário depois de ver o formulário
// só com os color pickers "vazios". A primeira é a própria identidade da
// Lumini Hub (src/theme/antd-theme.ts), usada como padrão quando a empresa
// ainda não configurou nada.
const PRESET_PALETTES: PresetPalette[] = [
  { name: "Azul Royal (padrão Lumini Hub)", primary: "#2563EB", secondary: "#06B6D4", accent: "#7C3AED" },
  { name: "Verde Esmeralda", primary: "#059669", secondary: "#10B981", accent: "#F59E0B" },
  { name: "Roxo Índigo", primary: "#4F46E5", secondary: "#818CF8", accent: "#EC4899" },
  { name: "Vermelho Coral", primary: "#DC2626", secondary: "#F87171", accent: "#FBBF24" },
  { name: "Grafite Neutro", primary: "#334155", secondary: "#64748B", accent: "#0EA5E9" },
  { name: "Laranja Vibrante", primary: "#EA580C", secondary: "#FB923C", accent: "#2563EB" },
];

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

  const applyPalette = (palette: PresetPalette) => {
    setValue("primary_color", palette.primary, { shouldValidate: true, shouldDirty: true });
    setValue("secondary_color", palette.secondary, { shouldValidate: true, shouldDirty: true });
    setValue("accent_color", palette.accent, { shouldValidate: true, shouldDirty: true });
  };

  if (loading) {
    return <Skeleton active paragraph={{ rows: 6 }} />;
  }

  const primaryColor = watch("primary_color");
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

          <Divider orientation="left" plain>
            Paletas Prontas
          </Divider>
          <Space wrap>
            {PRESET_PALETTES.map((palette) => (
              <Button key={palette.name} onClick={() => applyPalette(palette)}>
                <Space size={4}>
                  <span style={{ display: "flex" }}>
                    {[palette.primary, palette.secondary, palette.accent].map((color, index) => (
                      <span
                        key={index}
                        style={{
                          width: 14,
                          height: 14,
                          borderRadius: "50%",
                          background: color,
                          border: "1px solid rgba(0,0,0,0.15)",
                          marginLeft: index === 0 ? 0 : -5,
                        }}
                      />
                    ))}
                  </span>
                  {palette.name}
                </Space>
              </Button>
            ))}
          </Space>

          <Divider orientation="left" plain>
            Pré-visualização
          </Divider>
          <div style={{ border: "1px solid rgba(128,128,128,0.3)", borderRadius: 8, overflow: "hidden" }}>
            <div style={{ background: primaryColor, color: "#fff", padding: "10px 16px", fontWeight: 600 }}>
              Nome da Empresa
            </div>
            <div style={{ padding: 16, display: "flex", gap: 8, flexWrap: "wrap", alignItems: "center" }}>
              <Button type="primary" style={{ background: primaryColor, borderColor: primaryColor }}>
                Ação Primária
              </Button>
              {secondaryColor && <Tag color={secondaryColor}>Detalhe Secundário</Tag>}
              {accentColor && <Tag color={accentColor}>Destaque</Tag>}
            </div>
          </div>
          <Text type="secondary" style={{ fontSize: 12, display: "block", marginTop: 4 }}>
            Prévia isolada desta tela — o sistema ainda não aplica essas cores em outro lugar.
          </Text>

          <div style={{ display: "flex", justifyContent: "flex-end", marginTop: 16 }}>
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
