// src/components/settings/permissions/permission-form.tsx
"use client";

import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";

import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import { Loader2, HelpCircle } from "lucide-react"; // Importar Loader2 e HelpCircle (ícone de ?)

// Importe o seu componente SimpleTooltip
import { SimpleTooltip } from "@/components/common/SimpleTooltip"; // <--- Importe o seu SimpleTooltip

// Importe o schema e o tipo de dados do formulário
import {
  createPermissionSchema,
  CreatePermissionFormData,
} from "@/services/auth/permission-schema";
import { useEffect } from "react"; // Adicione o useEffect

interface PermissionFormProps {
  onSubmit?: (values: CreatePermissionFormData) => Promise<void>;
  initialValues?: CreatePermissionFormData;
  submitButtonText?: string;
  isSubmitting?: boolean;
  readOnly?: boolean;
}

export function PermissionForm({
  onSubmit,
  initialValues,
  submitButtonText = "Salvar Permissão",
  isSubmitting = false,
  readOnly = false,
}: PermissionFormProps) {
  const form = useForm<CreatePermissionFormData>({
    resolver: zodResolver(createPermissionSchema),
    defaultValues: initialValues || {
      permission: "",
      description: "",
      module: "",
    },
  });

  // useEffect para resetar o formulário com new defaultValues, se mudar (para edição)
  useEffect(() => {
    if (initialValues) {
      form.reset(initialValues);
    }
  }, [initialValues, form]);

  return (
    <Form {...form}>
      <form
        onSubmit={onSubmit ? form.handleSubmit(onSubmit) : undefined}
        className="space-y-6"
      >
        <FormField
          control={form.control}
          name="permission"
          render={({ field }) => (
            <FormItem>
              <div className="flex items-center space-x-1">
                <FormLabel>Nome da Permissão</FormLabel>
                <SimpleTooltip
                  label="Nome único e técnico da permissão (ex: `modulo.acao` ou `modulo.submodulo.acao`)."
                  side="top"
                >
                  <HelpCircle className="h-4 w-4 text-muted-foreground cursor-help" />
                </SimpleTooltip>
              </div>
              <FormControl>
                <Input
                  placeholder="Ex: users.view"
                  {...field}
                  readOnly={readOnly}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="description"
          render={({ field }) => (
            <FormItem>
              <div className="flex items-center space-x-1">
                <FormLabel>Descrição</FormLabel>
                <SimpleTooltip
                  label="Descrição amigável da permissão para fácil entendimento."
                  side="top"
                >
                  <HelpCircle className="h-4 w-4 text-muted-foreground cursor-help" />
                </SimpleTooltip>
              </div>
              <FormControl>
                <Textarea
                  placeholder="Ex: Visualizar lista de usuários"
                  {...field}
                  readOnly={readOnly}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* Campo Módulo */}
        <FormField
          control={form.control}
          name="module"
          render={({ field }) => (
            <FormItem>
              <div className="flex items-center space-x-1">
                <FormLabel>Módulo</FormLabel>
                {/* Use o SimpleTooltip aqui */}
                <SimpleTooltip
                  label="Módulo ao qual esta permissão pertence (ex: `users`, `sales`, `finance`)."
                  side="top"
                >
                  <HelpCircle className="h-4 w-4 text-muted-foreground cursor-help" />
                </SimpleTooltip>
              </div>
              <FormControl>
                <Input placeholder="Ex: users" {...field} readOnly={readOnly} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        {!readOnly && (
          <Button
            type="submit"
            className="w-full"
            disabled={isSubmitting || form.formState.isSubmitting}
          >
            {isSubmitting || form.formState.isSubmitting ? (
              <>
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                Salvando...
              </>
            ) : (
              submitButtonText
            )}
          </Button>
        )}
      </form>
    </Form>
  );
}
