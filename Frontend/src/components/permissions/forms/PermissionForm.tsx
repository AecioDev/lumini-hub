import { useEffect, useState } from "react";
import { zodResolver } from "@hookform/resolvers/zod";
import { Controller, useForm } from "react-hook-form";
import { AutoComplete, Button, Card, Input } from "antd";
import { FormField } from "@/components/common/FormField";
import { permissionService } from "@/services/permissions/permission-service";
import { useFeedback } from "@/hooks/useFeedback";
import {
  permissionSchema,
  type PermissionFormValues,
} from "@/schemas/permission-schema";

interface PermissionFormProps {
  initialValues?: PermissionFormValues;
  submitting: boolean;
  submitLabel: string;
  onSubmit: (values: PermissionFormValues) => void;
  onCancel: () => void;
}

export function PermissionForm({
  initialValues,
  submitting,
  submitLabel,
  onSubmit,
  onCancel,
}: PermissionFormProps) {
  const feedback = useFeedback();
  const [modules, setModules] = useState<string[]>([]);

  useEffect(() => {
    permissionService
      .modules()
      .then(setModules)
      .catch(() => feedback.error("Erro ao carregar módulos existentes."));
  }, [feedback]);

  const {
    control,
    handleSubmit,
    formState: { errors },
  } = useForm<PermissionFormValues>({
    resolver: zodResolver(permissionSchema),
    defaultValues: initialValues ?? { permission: "", description: "", module: "" },
  });

  return (
    <form onSubmit={handleSubmit(onSubmit)} noValidate>
      <Card style={{ maxWidth: 480 }}>
        <FormField label="Código da permissão" error={errors.permission}>
          <Controller
            name="permission"
            control={control}
            render={({ field }) => (
              <Input {...field} placeholder="Ex: reports.view (modulo.acao)" />
            )}
          />
        </FormField>
        <FormField label="Descrição" error={errors.description}>
          <Controller
            name="description"
            control={control}
            render={({ field }) => (
              <Input.TextArea {...field} rows={3} placeholder="O que essa permissão libera" />
            )}
          />
        </FormField>
        <FormField label="Módulo" error={errors.module}>
          <Controller
            name="module"
            control={control}
            render={({ field }) => (
              <AutoComplete
                {...field}
                options={modules.map((module) => ({ value: module }))}
                filterOption={(inputValue, option) =>
                  (option?.value ?? "")
                    .toLowerCase()
                    .includes(inputValue.toLowerCase())
                }
                placeholder="Selecione um módulo existente ou digite um novo"
              />
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
