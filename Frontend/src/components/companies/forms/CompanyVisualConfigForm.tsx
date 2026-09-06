import { useEffect, useState } from "react";
import { zodResolver } from "@hookform/resolvers/zod";
import { Controller, useForm } from "react-hook-form";
import {
  Button,
  Card,
  Col,
  ColorPicker,
  Divider,
  Input,
  Popconfirm,
  Row,
  Skeleton,
  Space,
  Tag,
  Typography,
  Upload,
  type UploadFile,
} from "antd";
import { CloseOutlined, InboxOutlined } from "@ant-design/icons";
import { FormField } from "@/components/common/FormField";
import {
  companyVisualConfigService,
  companyVisualConfigLogoUrl,
} from "@/services/companies/company-visual-config-service";
import { companyColorPaletteService } from "@/services/companies/company-color-palette-service";
import {
  companyVisualConfigSchema,
  type CompanyVisualConfigFormValues,
} from "@/schemas/company-visual-config-schema";
import { companyColorPaletteSchema } from "@/schemas/company-color-palette-schema";
import { getApiErrorMessage } from "@/utils/api-error";
import { useFeedback } from "@/hooks/useFeedback";
import type { ApiCompanyVisualConfigDetail } from "@/types/company-visual-config";
import type { ApiCompanyColorPalette } from "@/types/company-color-palette";

const { Text } = Typography;

const DEFAULT_PRIMARY_COLOR = "#2563EB";

interface ColorTriple {
  primary: string;
  secondary: string;
  accent: string;
}

interface PresetPalette extends ColorTriple {
  name: string;
}

// Paletas prontas pra preencher os 3 campos de uma vez, sem precisar
// garimpar hex manualmente — pedido do usuário depois de ver o formulário
// só com os color pickers "vazios". A primeira é a própria identidade da
// Lumini Hub (src/theme/antd-theme.ts), usada como padrão quando a empresa
// ainda não configurou nada.
const PRESET_PALETTES: PresetPalette[] = [
  {
    name: "Azul Royal (padrão Lumini Hub)",
    primary: "#2563EB",
    secondary: "#06B6D4",
    accent: "#7C3AED",
  },
  {
    name: "Verde Esmeralda",
    primary: "#059669",
    secondary: "#10B981",
    accent: "#F59E0B",
  },
  {
    name: "Roxo Índigo",
    primary: "#4F46E5",
    secondary: "#818CF8",
    accent: "#EC4899",
  },
  {
    name: "Vermelho Coral",
    primary: "#DC2626",
    secondary: "#F87171",
    accent: "#FBBF24",
  },
  {
    name: "Grafite Neutro",
    primary: "#334155",
    secondary: "#64748B",
    accent: "#0EA5E9",
  },
  {
    name: "Laranja Vibrante",
    primary: "#EA580C",
    secondary: "#FB923C",
    accent: "#2563EB",
  },
];

interface CompanyVisualConfigFormProps {
  companyId: number;
}

// Largura fixa pros botões de paleta (prontas + personalizadas) — sem isso
// cada botão fica do tamanho do próprio texto, deixando a fileira toda
// desalinhada (achado testando: "Roxo Índigo" bem mais estreito que "Azul
// Royal (padrão Lumini Hub)"). Nome mais longo que couber é truncado com
// reticências (o hex/nome completo continua no title/tooltip nativo).
const PALETTE_BUTTON_WIDTH = 210;

// Altura compartilhada entre a dropzone e o preview do logo (itemRender) —
// pra imagem ficar grande e proporcional ao lado do uploader, não uma
// thumbnail pequena perdida ao lado de um texto de nome de arquivo.
const LOGO_HEIGHT = 160;

// 3 bolinhas de cor sobrepostas, usadas tanto nas paletas prontas quanto nas
// personalizadas — cores vazias (paleta personalizada sem secundária/
// destaque) são simplesmente omitidas.
function PaletteSwatch({ colors }: { colors: (string | null | undefined)[] }) {
  return (
    <span style={{ display: "flex" }}>
      {colors
        .filter((color): color is string => !!color)
        .map((color, index) => (
          <span
            key={color}
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
  );
}

// Upload de logo é state irmão, fora do react-hook-form. Selecionar um
// arquivo (ou marcar o atual pra remover) só atualiza esse state local —
// nada é persistido até o clique em "Salvar Identidade Visual", junto com
// as cores. Decidido depois de um review do usuário: uma versão anterior
// chamava o endpoint de upload/remoção imediatamente ao selecionar/excluir,
// o que quebrava a semântica de "cancelar" (sair da tela sem salvar não
// desfazia nada) e desalinhava do fluxo de cores (só grava no Salvar).
export function CompanyVisualConfigForm({
  companyId,
}: CompanyVisualConfigFormProps) {
  const feedback = useFeedback();

  const [config, setConfig] = useState<ApiCompanyVisualConfigDetail | null>(
    null,
  );
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);

  // logoFile = arquivo novo escolhido, ainda não enviado (preview local via
  // blob URL). logoMarkedForRemoval = usuário clicou pra remover o logo já
  // salvo, mas isso só efetiva no backend quando o form for submetido.
  // logoFileList é só a representação visual dos dois (pro Upload.Dragger).
  const [logoFile, setLogoFile] = useState<File | null>(null);
  const [logoMarkedForRemoval, setLogoMarkedForRemoval] = useState(false);
  const [logoFileList, setLogoFileList] = useState<UploadFile[]>([]);

  const [customPalettes, setCustomPalettes] = useState<
    ApiCompanyColorPalette[]
  >([]);
  const [paletteName, setPaletteName] = useState("");
  const [savingPalette, setSavingPalette] = useState(false);

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
      text_color: "",
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
            text_color: result.text_color ?? "",
          });
          setLogoFileList(
            result.has_logo
              ? [
                  {
                    uid: "current-logo",
                    name: "Logo atual",
                    status: "done",
                    percent: 100,
                    thumbUrl: `${companyVisualConfigLogoUrl(result.id)}?v=0`,
                  },
                ]
              : [],
          );
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

  // Paletas personalizadas são hard-scoped por CompanyID direto (não pelo
  // id da CompanyVisualConfig) — carregam independente da config principal
  // já existir ou não.
  useEffect(() => {
    let cancelled = false;
    companyColorPaletteService
      .listByCompanyId(companyId)
      .then((result) => {
        if (!cancelled) setCustomPalettes(result);
      })
      .catch(() => feedback.error("Erro ao carregar paletas personalizadas."));
    return () => {
      cancelled = true;
    };
  }, [companyId, feedback]);

  // Salva cores e logo juntos, nessa ordem — o logo só é enviado/removido
  // de verdade aqui, nunca antes (ver comentário no topo do componente).
  const onSubmit = async (values: CompanyVisualConfigFormValues) => {
    setSubmitting(true);
    try {
      const payload = {
        primary_color: values.primary_color,
        secondary_color: values.secondary_color || undefined,
        accent_color: values.accent_color || undefined,
        text_color: values.text_color || undefined,
      };
      const wasCreate = !config;
      let result = config
        ? await companyVisualConfigService.update(config.id, payload)
        : await companyVisualConfigService.create({
            company_id: companyId,
            ...payload,
          });

      if (logoFile) {
        result = await companyVisualConfigService.uploadLogo(
          result.id,
          logoFile,
        );
        setLogoFile(null);
        setLogoFileList([
          {
            uid: "current-logo",
            name: "Logo atual",
            status: "done",
            percent: 100,
            thumbUrl: `${companyVisualConfigLogoUrl(result.id)}?v=${Date.now()}`,
          },
        ]);
      } else if (logoMarkedForRemoval) {
        result = await companyVisualConfigService.removeLogo(result.id);
        setLogoMarkedForRemoval(false);
        setLogoFileList([]);
      }

      setConfig(result);
      feedback.success(
        wasCreate
          ? "Identidade visual criada com sucesso."
          : "Identidade visual atualizada com sucesso.",
      );
    } catch (error) {
      feedback.error(
        getApiErrorMessage(error, "Erro ao salvar identidade visual."),
      );
    } finally {
      setSubmitting(false);
    }
  };

  // beforeUpload só guarda o arquivo localmente (preview via blob URL) e
  // devolve false pra impedir o antd de tentar subir sozinho — o envio de
  // verdade só acontece no onSubmit.
  const handleLogoSelected = (file: File) => {
    setLogoFile(file);
    setLogoMarkedForRemoval(false);
    setLogoFileList([
      {
        uid: "pending-logo",
        name: file.name,
        status: "done",
        percent: 100,
        thumbUrl: URL.createObjectURL(file),
      },
    ]);
    return false;
  };

  // Remover também é só local: se o item era o logo já salvo, marca pra
  // remoção no próximo Salvar; se era um arquivo recém-escolhido (ainda não
  // enviado), só descarta a seleção.
  const handleLogoRemoveClick = (file: UploadFile) => {
    if (file.uid === "current-logo") {
      setLogoMarkedForRemoval(true);
    }
    setLogoFile(null);
    setLogoFileList([]);
  };

  const applyPalette = (colors: ColorTriple) => {
    setValue("primary_color", colors.primary, {
      shouldValidate: true,
      shouldDirty: true,
    });
    setValue("secondary_color", colors.secondary, {
      shouldValidate: true,
      shouldDirty: true,
    });
    setValue("accent_color", colors.accent, {
      shouldValidate: true,
      shouldDirty: true,
    });
  };

  const handleSaveCustomPalette = async () => {
    const parsed = companyColorPaletteSchema.safeParse({
      name: paletteName.trim(),
      primary_color: watch("primary_color"),
      secondary_color: watch("secondary_color"),
      accent_color: watch("accent_color"),
    });
    if (!parsed.success) {
      feedback.error(
        parsed.error.issues[0]?.message ?? "Não foi possível salvar a paleta.",
      );
      return;
    }
    setSavingPalette(true);
    try {
      const created = await companyColorPaletteService.create({
        company_id: companyId,
        name: parsed.data.name,
        primary_color: parsed.data.primary_color,
        secondary_color: parsed.data.secondary_color || undefined,
        accent_color: parsed.data.accent_color || undefined,
      });
      setCustomPalettes((prev) => [created, ...prev]);
      setPaletteName("");
      feedback.success("Paleta personalizada salva.");
    } catch (error) {
      feedback.error(
        getApiErrorMessage(error, "Erro ao salvar paleta personalizada."),
      );
    } finally {
      setSavingPalette(false);
    }
  };

  const handleDeleteCustomPalette = async (id: number) => {
    try {
      await companyColorPaletteService.remove(id);
      setCustomPalettes((prev) => prev.filter((p) => p.id !== id));
      feedback.success("Paleta removida.");
    } catch (error) {
      feedback.error(getApiErrorMessage(error, "Erro ao remover paleta."));
    }
  };

  if (loading) {
    return <Skeleton active paragraph={{ rows: 6 }} />;
  }

  const primaryColor = watch("primary_color");
  const secondaryColor = watch("secondary_color");
  const accentColor = watch("accent_color");
  const textColor = watch("text_color");

  return (
    <div style={{ maxWidth: 720, margin: "0 auto" }}>
      <form onSubmit={handleSubmit(onSubmit)} noValidate>
        <Card title="Identidade Visual">
          <Row gutter={16}>
            <Col xs={24} sm={6}>
              <FormField
                label="Cor Primária"
                error={errors.primary_color}
                required
                stacked
              >
                <Controller
                  name="primary_color"
                  control={control}
                  render={({ field }) => (
                    <ColorPicker
                      value={field.value}
                      onChange={(color) =>
                        field.onChange(color.toHexString().toUpperCase())
                      }
                      format="hex"
                      disabledAlpha
                      showText
                    />
                  )}
                />
              </FormField>
            </Col>

            <Col xs={24} sm={6}>
              <FormField
                label="Cor Secundária"
                error={errors.secondary_color}
                stacked
              >
                <Space>
                  <Controller
                    name="secondary_color"
                    control={control}
                    render={({ field }) => (
                      <ColorPicker
                        value={field.value || undefined}
                        onChange={(color) =>
                          field.onChange(color.toHexString().toUpperCase())
                        }
                        format="hex"
                        disabledAlpha
                        showText
                      />
                    )}
                  />
                  {secondaryColor && (
                    <Button
                      type="link"
                      size="small"
                      onClick={() => setValue("secondary_color", "")}
                    >
                      Limpar
                    </Button>
                  )}
                </Space>
              </FormField>
            </Col>

            <Col xs={24} sm={6}>
              <FormField
                label="Cor de Destaque"
                error={errors.accent_color}
                stacked
              >
                <Space>
                  <Controller
                    name="accent_color"
                    control={control}
                    render={({ field }) => (
                      <ColorPicker
                        value={field.value || undefined}
                        onChange={(color) =>
                          field.onChange(color.toHexString().toUpperCase())
                        }
                        format="hex"
                        disabledAlpha
                        showText
                      />
                    )}
                  />
                  {accentColor && (
                    <Button
                      type="link"
                      size="small"
                      onClick={() => setValue("accent_color", "")}
                    >
                      Limpar
                    </Button>
                  )}
                </Space>
              </FormField>
            </Col>

            <Col xs={24} sm={6}>
              <FormField label="Cor do Texto" error={errors.text_color} stacked>
                <Space>
                  <Controller
                    name="text_color"
                    control={control}
                    render={({ field }) => (
                      <ColorPicker
                        value={field.value || undefined}
                        onChange={(color) =>
                          field.onChange(color.toHexString().toUpperCase())
                        }
                        format="hex"
                        disabledAlpha
                        showText
                      />
                    )}
                  />
                  {textColor && (
                    <Button
                      type="link"
                      size="small"
                      onClick={() => setValue("text_color", "")}
                    >
                      Limpar
                    </Button>
                  )}
                </Space>
              </FormField>
            </Col>
          </Row>
          <Text type="secondary" style={{ fontSize: 12 }}>
            Cor do texto sobre elementos na cor primária (cabeçalho, botões) —
            em branco se não escolher nada.
          </Text>

          <div
            style={{
              display: "flex",
              gap: 8,
              alignItems: "center",
              flexWrap: "wrap",
              marginTop: 8,
            }}
          >
            <Input
              placeholder="Nome da paleta personalizada"
              value={paletteName}
              onChange={(e) => setPaletteName(e.target.value)}
              style={{ maxWidth: 280 }}
            />
            <Button
              onClick={() => void handleSaveCustomPalette()}
              loading={savingPalette}
              disabled={!paletteName.trim()}
            >
              Salvar como paleta personalizada
            </Button>
          </div>

          <Divider orientation="left" plain>
            Paletas Prontas
          </Divider>
          <Space wrap>
            {PRESET_PALETTES.map((palette) => (
              <Button
                key={palette.name}
                onClick={() => applyPalette(palette)}
                style={{
                  width: PALETTE_BUTTON_WIDTH,
                  justifyContent: "flex-start",
                }}
                title={palette.name}
              >
                <Space size={4} style={{ overflow: "hidden", width: "100%" }}>
                  <PaletteSwatch
                    colors={[
                      palette.primary,
                      palette.secondary,
                      palette.accent,
                    ]}
                  />
                  <span
                    style={{
                      overflow: "hidden",
                      textOverflow: "ellipsis",
                      whiteSpace: "nowrap",
                    }}
                  >
                    {palette.name}
                  </span>
                </Space>
              </Button>
            ))}
          </Space>

          {customPalettes.length > 0 && (
            <>
              <Divider orientation="left" plain>
                Minhas Paletas
              </Divider>
              <Space wrap align="center">
                {customPalettes.map((palette) => (
                  <div
                    key={palette.id}
                    style={{
                      display: "flex",
                      alignItems: "stretch",
                      border: "1px solid rgba(128,128,128,0.35)",
                      borderRadius: 6,
                      overflow: "hidden",
                      width: PALETTE_BUTTON_WIDTH,
                    }}
                  >
                    <Button
                      type="text"
                      onClick={() =>
                        applyPalette({
                          primary: palette.primary_color,
                          secondary: palette.secondary_color ?? "",
                          accent: palette.accent_color ?? "",
                        })
                      }
                      style={{
                        flex: 1,
                        minWidth: 0,
                        justifyContent: "flex-start",
                        borderRadius: 0,
                      }}
                      title={palette.name}
                    >
                      <Space
                        size={4}
                        style={{ overflow: "hidden", width: "100%" }}
                      >
                        <PaletteSwatch
                          colors={[
                            palette.primary_color,
                            palette.secondary_color,
                            palette.accent_color,
                          ]}
                        />
                        <span
                          style={{
                            overflow: "hidden",
                            textOverflow: "ellipsis",
                            whiteSpace: "nowrap",
                          }}
                        >
                          {palette.name}
                        </span>
                      </Space>
                    </Button>
                    <Popconfirm
                      title="Remover esta paleta?"
                      okText="Remover"
                      cancelText="Cancelar"
                      onConfirm={() =>
                        void handleDeleteCustomPalette(palette.id)
                      }
                    >
                      <Button
                        type="text"
                        danger
                        icon={<CloseOutlined />}
                        style={{
                          borderRadius: 0,
                          borderLeft: "1px solid rgba(128,128,128,0.35)",
                          flexShrink: 0,
                        }}
                      />
                    </Popconfirm>
                  </div>
                ))}
              </Space>
            </>
          )}

          <Divider orientation="left" plain>
            Pré-visualização
          </Divider>
          <div
            style={{
              border: "1px solid rgba(128,128,128,0.3)",
              borderRadius: 8,
              overflow: "hidden",
            }}
          >
            <div
              style={{
                background: primaryColor,
                color: textColor || "#fff",
                padding: "10px 16px",
                fontWeight: 600,
              }}
            >
              Nome da Empresa
            </div>
            <div
              style={{
                padding: 16,
                display: "flex",
                gap: 8,
                flexWrap: "wrap",
                alignItems: "center",
              }}
            >
              <Button
                type="primary"
                style={{
                  background: primaryColor,
                  borderColor: primaryColor,
                  color: textColor || undefined,
                }}
              >
                Ação Primária
              </Button>
              {secondaryColor && (
                <Tag color={secondaryColor}>Detalhe Secundário</Tag>
              )}
              {accentColor && <Tag color={accentColor}>Destaque</Tag>}
            </div>
          </div>
          <Text
            type="secondary"
            style={{ fontSize: 12, display: "block", marginTop: 4 }}
          >
            Prévia isolada desta tela — o sistema ainda não aplica essas cores
            em outro lugar.
          </Text>

          <Divider>Logo</Divider>

          {/* .ant-upload-wrapper agrupa a dropzone (.ant-upload-drag) e a
              lista de arquivos (.ant-upload-list) como filhos diretos — por
              padrão o antd empilha um embaixo do outro; o usuário pediu o
              logo ao lado do uploader, então viram itens de um flex-row via
              essas classes públicas e estáveis do antd (não é
              reestruturação de DOM interno, só CSS).
              Não depende mais de `config` existir — selecionar um logo antes
              mesmo de ter salvo cor nenhuma funciona, o upload acontece
              junto na primeira vez que "Salvar Identidade Visual" for
              clicado (ver onSubmit). */}
          <div
            className="logo-uploader"
            style={{ maxWidth: logoFileList.length > 0 ? 620 : 320 }}
          >
            <style>{`
              .logo-uploader .ant-upload-wrapper { display: flex; gap: 16px; align-items: flex-start; }
              .logo-uploader .ant-upload-list { flex: 1; min-width: 160px; margin-top: 0; }
              .logo-uploader .ant-upload-list-item { padding: 0; border: none; height: auto; }
            `}</style>
            <Upload.Dragger
              listType="picture"
              fileList={logoFileList}
              beforeUpload={handleLogoSelected}
              onRemove={handleLogoRemoveClick}
              accept=".png,.jpg,.jpeg,.svg"
              maxCount={1}
              style={{ width: 320, height: LOGO_HEIGHT, flexShrink: 0 }}
              // Item da lista customizado: só a imagem grande + um botão de
              // remover ao lado — o item padrão do antd (thumbnail pequena +
              // nome do arquivo em linha) ficava desproporcional ao lado de
              // uma dropzone desse tamanho.
              itemRender={(_originNode, file, _list, actions) => (
                <div style={{ display: "flex", alignItems: "center", gap: 8 }}>
                  <img
                    src={file.thumbUrl || file.url}
                    alt={file.name}
                    style={{
                      height: LOGO_HEIGHT,
                      maxWidth: 160,
                      objectFit: "contain",
                      borderRadius: 8,
                      border: "1px solid rgba(128,128,128,0.35)",
                    }}
                  />
                  <Button
                    type="text"
                    danger
                    icon={<CloseOutlined />}
                    onClick={() => actions.remove()}
                  />
                </div>
              )}
            >
              <p className="ant-upload-drag-icon" style={{ margin: "4px 0" }}>
                <InboxOutlined />
              </p>
              <p className="ant-upload-text">
                Clique ou arraste o logo pra esta área
              </p>
              <p className="ant-upload-hint">PNG, JPEG ou SVG</p>
            </Upload.Dragger>
          </div>

          <div
            style={{
              display: "flex",
              justifyContent: "flex-end",
              marginTop: 24,
            }}
          >
            <Button type="primary" htmlType="submit" loading={submitting}>
              {config ? "Salvar Identidade Visual" : "Criar Identidade Visual"}
            </Button>
          </div>
        </Card>
      </form>
    </div>
  );
}
