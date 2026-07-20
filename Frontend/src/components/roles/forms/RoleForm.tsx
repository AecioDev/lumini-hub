import { zodResolver } from "@hookform/resolvers/zod";
import { Controller, useForm } from "react-hook-form";
import { Button, Card, Input } from "antd";
import { FormField } from "@/components/common/FormField";
import { roleSchema, type RoleFormValues } from "@/schemas/role-schema";

interface RoleFormProps {
  initialValues?: RoleFormValues;
  submitting: boolean;
  submitLabel: string;
  onSubmit: (values: RoleFormValues) => void;
  onCancel: () => void;
}

export function RoleForm({
  initialValues,
  submitting,
  submitLabel,
  onSubmit,
  onCancel,
}: RoleFormProps) {
  const {
    control,
    handleSubmit,
    formState: { errors },
  } = useForm<RoleFormValues>({
    resolver: zodResolver(roleSchema),
    defaultValues: initialValues ?? { name: "", description: "" },
  });

  return (
    <form onSubmit={handleSubmit(onSubmit)} noValidate>
      <Card style={{ maxWidth: 480 }}>
        <FormField label="Nome do perfil" error={errors.name}>
          <Controller
            name="name"
            control={control}
            render={({ field }) => <Input {...field} placeholder="Ex: Vendedor" />}
          />
        </FormField>
        <FormField label="Descrição" error={errors.description}>
          <Controller
            name="description"
            control={control}
            render={({ field }) => (
              <Input.TextArea {...field} rows={3} placeholder="Descrição do perfil" />
            )}
          />
        </FormField>
      </Card>

      <div style={{ display: "flex", justifyContent: "flex-end", gap: 10, marginTop: 16 }}>
        <Button onClick={onCancel}>Cancelar</Button>
        <Button type="primary" htmlType="submit" loading={submitting}>
          {submitLabel}
        </Button>
      </div>
    </form>
  );
}
