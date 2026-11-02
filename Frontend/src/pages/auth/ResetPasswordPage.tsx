import { useState } from "react";
import { zodResolver } from "@hookform/resolvers/zod";
import { Controller, useForm } from "react-hook-form";
import { Alert, Button, Input, Typography } from "antd";
import { Link, useNavigate, useSearchParams } from "react-router-dom";
import { FormField } from "@/components/common/FormField";
import { useFeedback } from "@/hooks/useFeedback";
import { authService } from "@/services/auth/auth-service";
import { getApiErrorMessage } from "@/utils/api-error";
import { resetPasswordSchema, type ResetPasswordFormValues } from "@/schemas/reset-password-schema";

const { Title, Paragraph } = Typography;

// Destino do link do e-mail de recuperação: /reset-password?token=...
export function ResetPasswordPage() {
  const [searchParams] = useSearchParams();
  const token = searchParams.get("token") ?? "";
  const navigate = useNavigate();
  const feedback = useFeedback();
  const [submitting, setSubmitting] = useState(false);
  const [linkError, setLinkError] = useState<string | null>(null);

  const {
    control,
    handleSubmit,
    formState: { errors },
  } = useForm<ResetPasswordFormValues>({
    resolver: zodResolver(resetPasswordSchema),
    defaultValues: { new_password: "", confirm_password: "" },
  });

  const onSubmit = async (values: ResetPasswordFormValues) => {
    setSubmitting(true);
    setLinkError(null);
    try {
      await authService.resetPassword(token, values.new_password);
      feedback.success("Senha redefinida com sucesso. Entre com a nova senha.");
      navigate("/login", { replace: true });
    } catch (error) {
      setLinkError(getApiErrorMessage(error, "Não foi possível redefinir a senha. Tente novamente."));
    } finally {
      setSubmitting(false);
    }
  };

  if (!token) {
    return (
      <div>
        <Title level={3} style={{ marginTop: 0, marginBottom: 6 }}>
          Link inválido
        </Title>
        <Paragraph type="secondary">
          O link de redefinição está incompleto. Peça um novo para receber outro e-mail.
        </Paragraph>
        <Paragraph style={{ marginBottom: 0 }}>
          <Link to="/forgot-password">Pedir novo link</Link>
        </Paragraph>
      </div>
    );
  }

  return (
    <div>
      <Title level={3} style={{ marginTop: 0, marginBottom: 6 }}>
        Nova senha
      </Title>
      <Paragraph type="secondary" style={{ marginBottom: 28 }}>
        Escolha a nova senha da sua conta
      </Paragraph>

      {linkError && (
        <Alert
          type="error"
          showIcon
          style={{ marginBottom: 20 }}
          message={linkError}
          description={<Link to="/forgot-password">Pedir novo link</Link>}
        />
      )}

      <form onSubmit={handleSubmit(onSubmit)} noValidate>
        <FormField label="Nova senha" error={errors.new_password}>
          <Controller
            name="new_password"
            control={control}
            render={({ field }) => (
              <Input.Password {...field} size="large" autoComplete="new-password" autoFocus />
            )}
          />
        </FormField>

        <FormField label="Confirmar senha" error={errors.confirm_password}>
          <Controller
            name="confirm_password"
            control={control}
            render={({ field }) => <Input.Password {...field} size="large" autoComplete="new-password" />}
          />
        </FormField>

        <Button type="primary" htmlType="submit" size="large" block loading={submitting} style={{ marginTop: 8 }}>
          Redefinir senha
        </Button>
      </form>

      <Paragraph style={{ textAlign: "center", marginTop: 20, marginBottom: 0 }}>
        <Link to="/login">Voltar ao login</Link>
      </Paragraph>
    </div>
  );
}
