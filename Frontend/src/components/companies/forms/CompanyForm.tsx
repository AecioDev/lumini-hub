import { useEffect, useMemo, useState } from "react";
import { zodResolver } from "@hookform/resolvers/zod";
import { Controller, useForm } from "react-hook-form";
import { Button, Card, Input, Select, Switch, Typography, Upload } from "antd";
import type { UploadProps } from "antd";
import { DeleteOutlined, UploadOutlined } from "@ant-design/icons";
import { FormField } from "@/components/common/FormField";
import { companyLogoUrl, companyService } from "@/services/companies/company-service";
import { getApiErrorMessage } from "@/utils/api-error";
import { useFeedback } from "@/hooks/useFeedback";
import { companySchema, type CompanyFormValues } from "@/schemas/company-schema";
import type { ApiCompany } from "@/types/company";

const { Text } = Typography;

const ALLOWED_LOGO_TYPES = ["image/png", "image/jpeg", "image/svg+xml"];
const MAX_LOGO_SIZE = 2 * 1024 * 1024; // 2MB, mesmo limite do backend (CFG-7.1.2)

interface CompanyFormProps {
  initialValues?: CompanyFormValues;
  submitting: boolean;
  submitLabel: string;
  /** Modo edição mostra o campo Ativo, a seção de Logo e exclui a própria empresa das opções de vínculo. */
  editingId?: number;
  /** Empresa já tem logo cadastrado — só relevante em modo edição. */
  hasLogo?: boolean;
  /** Chamado depois de upload/remoção de logo bem-sucedidos, com o novo valor de has_logo. */
  onLogoChange?: (hasLogo: boolean) => void;
  /**
   * Salva os dados cadastrais da empresa. Precisa devolver/propagar rejeição
   * em caso de erro — CompanyForm usa isso pra decidir se segue ou não pro
   * passo de persistir o logo pendente (ver comentário em handleFormSubmit).
   */
  onSubmit: (values: CompanyFormValues) => void | Promise<void>;
  onCancel: () => void;
}

export function CompanyForm({
  initialValues,
  submitting,
  submitLabel,
  editingId,
  hasLogo,
  onLogoChange,
  onSubmit,
  onCancel,
}: CompanyFormProps) {
  const feedback = useFeedback();
  const [companies, setCompanies] = useState<ApiCompany[]>([]);

  // Logo é "rascunho" — nada é enviado ao servidor até o clique em "Salvar
  // Empresa" (CFG-7.4.1, redesenhado em 2026-09-07: a versão "ação isolada e
  // imediata" original tinha um bug real onde salvar os dados da empresa
  // limpava o logo recém-enviado; a causa exata não foi isolada apesar de
  // várias tentativas de reprodução, então a decisão foi adotar o mesmo
  // padrão de rascunho que a antiga CompanyVisualConfigForm.tsx já usava).
  // logoFile = arquivo novo escolhido, ainda não enviado (preview via blob
  // URL). logoMarkedForRemoval = usuário pediu pra remover o logo já salvo,
  // efetivado só no próximo Salvar.
  const [logoFile, setLogoFile] = useState<File | null>(null);
  const [logoMarkedForRemoval, setLogoMarkedForRemoval] = useState(false);
  const [savingLogo, setSavingLogo] = useState(false);
  // Cache-busting: a URL do logo salvo não muda quando o arquivo é trocado/
  // removido (mesmo ID), então o navegador serviria a imagem antiga do cache
  // sem isso.
  const [logoVersion, setLogoVersion] = useState(0);

  const stagedLogoPreviewUrl = useMemo(
    () => (logoFile ? URL.createObjectURL(logoFile) : null),
    [logoFile]
  );
  useEffect(() => {
    return () => {
      if (stagedLogoPreviewUrl) URL.revokeObjectURL(stagedLogoPreviewUrl);
    };
  }, [stagedLogoPreviewUrl]);

  useEffect(() => {
    companyService.list().then(setCompanies).catch(() => setCompanies([]));
  }, []);

  const {
    control,
    handleSubmit,
    watch,
    formState: { errors },
  } = useForm<CompanyFormValues>({
    resolver: zodResolver(companySchema),
    defaultValues: initialValues ?? {
      parent_id: null,
      legal_name: "",
      trade_name: "",
      cnpj: "",
      is_active: true,
    },
  });

  // CNPJ só é obrigatório pra Matriz (sem vínculo) — reflete isso no label
  // em tempo real conforme o usuário mexe no Select "Vinculada a".
  const isSubsidiary = watch("parent_id") !== null;

  // Só pode vincular a uma empresa que já não seja ela mesma — o backend
  // também valida (e detecta ciclo pra vínculos indiretos), isso aqui é só
  // pra não oferecer a opção óbvia de auto-vínculo na UI.
  const parentOptions = companies
    .filter((c) => c.id !== editingId)
    .map((c) => ({ value: c.id, label: c.trade_name || c.legal_name }));

  // beforeUpload só guarda o arquivo localmente (preview via blob URL) e
  // devolve Upload.LIST_IGNORE pra impedir o antd de tentar subir sozinho —
  // o envio de verdade só acontece em handleFormSubmit, junto do resto dos
  // dados da empresa.
  const handleLogoSelected: UploadProps["beforeUpload"] = (file) => {
    if (!ALLOWED_LOGO_TYPES.includes(file.type)) {
      feedback.error("Formato não suportado — envie PNG, JPEG ou SVG.");
      return Upload.LIST_IGNORE;
    }
    if (file.size > MAX_LOGO_SIZE) {
      feedback.error("O logo não pode ultrapassar 2MB.");
      return Upload.LIST_IGNORE;
    }

    setLogoFile(file);
    setLogoMarkedForRemoval(false);
    return Upload.LIST_IGNORE;
  };

  // Remover também é só local: se havia um arquivo recém-selecionado ainda
  // não salvo, só descarta a seleção; se era o logo já salvo (hasLogo), marca
  // pra remoção no próximo Salvar.
  const handleRemoveLogoClick = () => {
    if (logoFile) {
      setLogoFile(null);
      return;
    }
    if (hasLogo) {
      setLogoMarkedForRemoval(true);
    }
  };

  // Salva os dados da empresa primeiro (delegado ao onSubmit do pai) e só
  // depois, se o passo anterior deu certo, persiste o logo pendente — nessa
  // ordem, nunca antes (CFG-7.4.1, redesenhado em 2026-09-07, ver comentário
  // no topo do componente). Em modo criação (sem editingId) não há logo
  // pendente possível — a página de criação cuida do próprio feedback.
  const handleFormSubmit = async (values: CompanyFormValues) => {
    try {
      await onSubmit(values);
    } catch {
      // onSubmit já mostrou o próprio erro — não mexe no logo pendente,
      // usuário pode tentar salvar de novo sem perder a seleção.
      return;
    }

    if (!editingId || (!logoFile && !logoMarkedForRemoval)) {
      if (editingId) feedback.success("Empresa atualizada com sucesso.");
      return;
    }

    setSavingLogo(true);
    try {
      const updated = logoFile
        ? await companyService.uploadLogo(editingId, logoFile)
        : await companyService.removeLogo(editingId);
      setLogoFile(null);
      setLogoMarkedForRemoval(false);
      setLogoVersion((v) => v + 1);
      onLogoChange?.(updated.has_logo);
      feedback.success("Empresa atualizada com sucesso.");
    } catch (error) {
      feedback.error(
        getApiErrorMessage(
          error,
          "Dados da empresa salvos, mas houve um erro ao salvar o logo. Tente novamente."
        )
      );
    } finally {
      setSavingLogo(false);
    }
  };

  const logoPreviewUrl = logoMarkedForRemoval
    ? null
    : stagedLogoPreviewUrl ?? (hasLogo && editingId ? `${companyLogoUrl(editingId)}?v=${logoVersion}` : null);

  return (
    <form onSubmit={handleSubmit(handleFormSubmit)} noValidate>
      <div style={{ maxWidth: 560, margin: "0 auto" }}>
        <Card>
          <FormField label="Razão Social" error={errors.legal_name}>
            <Controller
              name="legal_name"
              control={control}
              render={({ field }) => <Input {...field} placeholder="Razão Social LTDA" />}
            />
          </FormField>

          <FormField label="Nome Fantasia" error={errors.trade_name}>
            <Controller
              name="trade_name"
              control={control}
              render={({ field }) => <Input {...field} placeholder="Nome Fantasia" />}
            />
          </FormField>

          <FormField
            label={isSubsidiary ? "CNPJ (opcional)" : "CNPJ"}
            error={errors.cnpj}
          >
            <Controller
              name="cnpj"
              control={control}
              render={({ field }) => (
                <Input
                  {...field}
                  placeholder={
                    isSubsidiary
                      ? "Deixe em branco se a parte fiscal ficar com a Matriz"
                      : "00.000.000/0000-00"
                  }
                />
              )}
            />
          </FormField>

          <FormField label="Vinculada a" error={errors.parent_id}>
            <Controller
              name="parent_id"
              control={control}
              render={({ field }) => (
                <Select
                  {...field}
                  allowClear
                  placeholder="Nenhuma — esta é a empresa Matriz"
                  options={parentOptions}
                  onChange={(value) => field.onChange(value ?? null)}
                />
              )}
            />
          </FormField>

          {editingId && (
            <FormField label="Ativa">
              <Controller
                name="is_active"
                control={control}
                render={({ field: { value, onChange } }) => (
                  <Switch checked={value} onChange={onChange} />
                )}
              />
            </FormField>
          )}

          {editingId && (
            <FormField label="Logo" stacked>
              <div style={{ display: "flex", alignItems: "center", gap: 16 }}>
                {logoPreviewUrl ? (
                  <img
                    src={logoPreviewUrl}
                    alt="Logo da empresa"
                    style={{
                      height: 64,
                      maxWidth: 160,
                      objectFit: "contain",
                      border: "1px solid #d9d9d9",
                      borderRadius: 8,
                      padding: 4,
                    }}
                  />
                ) : (
                  <div
                    style={{
                      height: 64,
                      width: 160,
                      display: "flex",
                      alignItems: "center",
                      justifyContent: "center",
                      border: "1px dashed #d9d9d9",
                      borderRadius: 8,
                      color: "rgba(0, 0, 0, 0.45)",
                      fontSize: 12,
                    }}
                  >
                    Sem logo
                  </div>
                )}

                <Upload
                  accept={ALLOWED_LOGO_TYPES.join(",")}
                  showUploadList={false}
                  beforeUpload={handleLogoSelected}
                >
                  <Button icon={<UploadOutlined />}>
                    {logoPreviewUrl ? "Trocar Logo" : "Enviar Logo"}
                  </Button>
                </Upload>

                {logoPreviewUrl && (
                  <Button danger type="text" icon={<DeleteOutlined />} onClick={handleRemoveLogoClick}>
                    Remover
                  </Button>
                )}
              </div>

              {logoFile && (
                <Text type="secondary" style={{ fontSize: 12, display: "block", marginTop: 6 }}>
                  Novo logo selecionado — será enviado ao clicar em "{submitLabel}".
                </Text>
              )}
              {logoMarkedForRemoval && (
                <Text type="secondary" style={{ fontSize: 12, display: "block", marginTop: 6 }}>
                  O logo será removido ao clicar em "{submitLabel}".
                </Text>
              )}
            </FormField>
          )}
        </Card>

        <div style={{ display: "flex", justifyContent: "flex-end", gap: 10, marginTop: 16 }}>
          <Button onClick={onCancel}>Cancelar</Button>
          <Button type="primary" htmlType="submit" loading={submitting || savingLogo}>
            {submitLabel}
          </Button>
        </div>
      </div>
    </form>
  );
}
