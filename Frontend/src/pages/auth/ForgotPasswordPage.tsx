import { useState } from "react";
import { zodResolver } from "@hookform/resolvers/zod";
import { Controller, useForm } from "react-hook-form";
import { App as AntdApp, Button, Input, Typography } from "antd";
import { Link } from "react-router-dom";
import { FormField } from "@/components/common/FormField";
import {
  forgotPasswordSchema,
  type ForgotPasswordFormValues,
} from "@/schemas/forgot-password-schema";

const { Title, Paragraph } = Typography;

export function ForgotPasswordPage() {
  const { message } = AntdApp.useApp();
  const [submitting, setSubmitting] = useState(false);

  const {
    control,
    handleSubmit,
    formState: { errors },
  } = useForm<ForgotPasswordFormValues>({
    resolver: zodResolver(forgotPasswordSchema),
    defaultValues: { email: "" },
  });

  const onSubmit = async () => {
    setSubmitting(true);
    try {
      // O backend ainda não expõe um endpoint de recuperação de senha
      // (só /auth/login|refresh-token|logout|me). A tela fica pronta para
      // plugar em POST /auth/forgot-password (ou equivalente) assim que existir.
      await new Promise((resolve) => setTimeout(resolve, 500));
      message.info(
        "Recuperação de senha ainda não está disponível: falta o endpoint de envio de e-mail no backend."
      );
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div>
      <Title level={3} style={{ marginTop: 0, marginBottom: 6 }}>
        Recuperar senha
      </Title>
      <Paragraph type="secondary" style={{ marginBottom: 28 }}>
        Enviaremos um link de redefinição para seu e-mail
      </Paragraph>

      <form onSubmit={handleSubmit(onSubmit)} noValidate>
        <FormField label="E-mail" error={errors.email}>
          <Controller
            name="email"
            control={control}
            render={({ field }) => (
              <Input {...field} size="large" placeholder="voce@empresa.com" autoFocus />
            )}
          />
        </FormField>

        <Button
          type="primary"
          htmlType="submit"
          size="large"
          block
          loading={submitting}
          style={{ marginTop: 8 }}
        >
          Enviar link de recuperação
        </Button>
      </form>

      <Paragraph style={{ textAlign: "center", marginTop: 20, marginBottom: 0 }}>
        <Link to="/login">Voltar ao login</Link>
      </Paragraph>
    </div>
  );
}
