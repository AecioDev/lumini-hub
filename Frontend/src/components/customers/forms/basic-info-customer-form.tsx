"use client";

import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  customerFormSchema,
  CreateCustomerInput,
} from "@/services/customer/customer-service";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import { Switch } from "@/components/ui/switch";
import { useToast } from "@/components/ui/use-toast";
import CustomerService from "@/services/customer/customer-service";
import { useRouter } from "next/navigation";
import React from "react";

// --- Componente Reutilizável para Campo de Formulário ---
interface FormFieldProps {
  label: string;
  id: string;
  children: React.ReactNode;
  errorMessage?: string;
  className?: string;
}

const FormField: React.FC<FormFieldProps> = ({
  label,
  id,
  children,
  errorMessage,
  className,
}) => (
  <div className={className}>
    <Label htmlFor={id} className="text-foreground">
      {label}
    </Label>
    {children}
    {errorMessage && (
      <p className="text-sm text-destructive mt-1">{errorMessage}</p>
    )}
  </div>
);

// --- Componente Principal do Formulário ---
export default function BasicInfoCustomerForm() {
  const {
    register,
    handleSubmit,
    watch,
    setValue,
    formState: { errors, isSubmitting },
  } = useForm<CreateCustomerInput>({
    resolver: zodResolver(customerFormSchema),
    defaultValues: {
      is_active: true,
      person_type: "fisica",
    },
  });

  const { toast } = useToast();
  const router = useRouter();

  const personType = watch("person_type");
  const isActive = watch("is_active");

  const onSubmit = async (data: CreateCustomerInput) => {
    try {
      await CustomerService.createCustomer(data);
      toast({
        title: "Cliente cadastrado com sucesso!",
        description: "Você será redirecionado em breve.",
      });
      router.push("/vendas/cadastros/clientes");
    } catch (error) {
      console.error("Erro ao cadastrar cliente:", error);
      toast({
        variant: "destructive",
        title: "Erro ao cadastrar cliente",
        description:
          "Verifique os dados e tente novamente. Detalhes: " +
          (error instanceof Error ? error.message : "Erro desconhecido"),
      });
    }
  };

  return (
    <form
      onSubmit={handleSubmit(onSubmit)}
      className="bg-form-background space-y-6 p-6 rounded-lg shadow-md"
    >
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <FormField
          label="Nome"
          id="first_name"
          errorMessage={errors.first_name?.message}
        >
          <Input
            id="first_name"
            {...register("first_name")}
            className="border-input focus-visible:ring-ring"
          />
        </FormField>

        <FormField
          label="Sobrenome"
          id="last_name"
          errorMessage={errors.last_name?.message}
        >
          <Input
            id="last_name"
            {...register("last_name")}
            className="border-input focus-visible:ring-ring"
          />
        </FormField>

        <FormField label="Tipo de Pessoa" id="person_type">
          <select
            id="person_type"
            {...register("person_type")}
            className="w-full p-2 border border-input rounded-md focus:outline-none focus:ring-2 focus:ring-ring bg-input text-foreground"
          >
            <option value="fisica">Pessoa Física</option>
            <option value="juridica">Pessoa Jurídica</option>
          </select>
        </FormField>

        <FormField
          label="Número do Documento"
          id="document_number"
          errorMessage={errors.document_number?.message}
        >
          <Input
            id="document_number"
            {...register("document_number")}
            placeholder={personType === "fisica" ? "CPF" : "CNPJ"}
            className="border-input focus-visible:ring-ring"
          />
        </FormField>

        {personType === "juridica" && (
          <FormField
            label="Razão Social"
            id="company_name"
            errorMessage={errors.company_name?.message}
            className="md:col-span-2"
          >
            <Input
              id="company_name"
              {...register("company_name")}
              className="border-input focus-visible:ring-ring"
            />
          </FormField>
        )}

        <FormField label="Observações" id="notes" className="md:col-span-2">
          <Input
            id="notes"
            {...register("notes")}
            className="border-input focus-visible:ring-ring"
          />
        </FormField>

        <div className="flex items-center space-x-2 md:col-span-2">
          <Switch
            id="is_active"
            checked={isActive}
            onCheckedChange={(checked) => setValue("is_active", checked)}
            {...register("is_active")}
          />
          <Label htmlFor="is_active" className="text-foreground">
            Cliente Ativo
          </Label>
        </div>
      </div>

      <Button
        type="submit"
        disabled={isSubmitting}
        className="w-full md:w-auto mt-6 px-8 py-2"
      >
        {isSubmitting ? "Salvando..." : "Salvar Cliente"}
      </Button>
    </form>
  );
}
