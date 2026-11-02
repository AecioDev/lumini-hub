import { useState } from "react";
import { zodResolver } from "@hookform/resolvers/zod";
import { Controller, useForm } from "react-hook-form";
import { Alert, Button, Input, Typography } from "antd";
import { Link } from "react-router-dom";
import { FormField } from "@/components/common/FormField";
import { useFeedback } from "@/hooks/useFeedback";
import { authService } from "@/services/auth/auth-service";
import { getApiErrorMessage } from "@/utils/api-error";
import {
  forgotPasswordSchema,
  type ForgotPasswordFormValues,
} from "@/schemas/forgot-password-schema";

const { Title, Paragraph } = Typography;

export function ForgotPasswordPage() {
  const feedback = useFeedback();
  const [submitting, setSubmitting] = useState(false);
  const [sent, setSent] = useState(false);

  const {
    control,
    handleSubmit,
    formState: { errors },
  } = useForm<ForgotPasswordFormValues>({
    resolver: zodResolver(forgotPasswordSchema),
    defaultValues: { email: "" },
  });

  const onSubmit = async (values: ForgotPasswordFormValues) => {
    setSubmitting(true);
    try {
      await authService.forgotPassword(values.email);
      setSent(true);
    } catch (error) {
      feedback.error(getApiErrorMessage(error, "Não foi possível enviar o link agora. Tente novamente."));
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

      {sent && (
        <Alert
          type="success"
          showIcon
          style={{ marginBottom: 20 }}
          message="Verifique seu e-mail"
          description="Se o e-mail estiver cadastrado, enviamos um link para redefinir a senha. Ele vale por 30 minutos."
        />
      )}

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
          {sent ? "Enviar novamente" : "Enviar link de recuperação"}
        </Button>
      </form>

      <Paragraph style={{ textAlign: "center", marginTop: 20, marginBottom: 0 }}>
        <Link to="/login">Voltar ao login</Link>
      </Paragraph>
    </div>
  );
}
