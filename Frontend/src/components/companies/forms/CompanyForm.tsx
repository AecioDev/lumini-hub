import { useEffect, useState } from "react";
import { zodResolver } from "@hookform/resolvers/zod";
import { Controller, useForm } from "react-hook-form";
import { Button, Card, Input, Select, Switch } from "antd";
import { FormField } from "@/components/common/FormField";
import { companyService } from "@/services/companies/company-service";
import { companySchema, type CompanyFormValues } from "@/schemas/company-schema";
import type { ApiCompany } from "@/types/company";

interface CompanyFormProps {
  initialValues?: CompanyFormValues;
  submitting: boolean;
  submitLabel: string;
  /** Modo edição mostra o campo Ativo e exclui a própria empresa das opções de vínculo. */
  editingId?: number;
  onSubmit: (values: CompanyFormValues) => void;
  onCancel: () => void;
}

export function CompanyForm({
  initialValues,
  submitting,
  submitLabel,
  editingId,
  onSubmit,
  onCancel,
}: CompanyFormProps) {
  const [companies, setCompanies] = useState<ApiCompany[]>([]);

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
      tax_id: "",
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

  return (
    <form onSubmit={handleSubmit(onSubmit)} noValidate>
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
            error={errors.tax_id}
          >
            <Controller
              name="tax_id"
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
        </Card>

        <div style={{ display: "flex", justifyContent: "flex-end", gap: 10, marginTop: 16 }}>
          <Button onClick={onCancel}>Cancelar</Button>
          <Button type="primary" htmlType="submit" loading={submitting}>
            {submitLabel}
          </Button>
        </div>
      </div>
    </form>
  );
}
