import { useState } from "react";
import { zodResolver } from "@hookform/resolvers/zod";
import { Controller, useForm } from "react-hook-form";
import { Alert, Button, Input, Typography } from "antd";
import { useNavigate } from "react-router-dom";
import { FormField } from "@/components/common/FormField";
import { useAuth } from "@/contexts/AuthContext";
import { useFeedback } from "@/hooks/useFeedback";
import { userService } from "@/services/users/user-service";
import { getApiErrorMessage } from "@/utils/api-error";
import { changePasswordSchema, type ChangePasswordFormValues } from "@/schemas/change-password-schema";

const { Title, Paragraph } = Typography;

// Troca de senha obrigatória: o ProtectedRoute leva para cá o usuário marcado com must_change_password
// (cadastrado sem convite, com a senha inicial definida por quem o cadastrou) antes de qualquer outra tela.
export function ChangePasswordPage() {
  const { user, refreshUser, logout } = useAuth();
  const navigate = useNavigate();
  const feedback = useFeedback();
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const {
    control,
    handleSubmit,
    formState: { errors },
  } = useForm<ChangePasswordFormValues>({
    resolver: zodResolver(changePasswordSchema),
    defaultValues: { current_password: "", new_password: "", confirm_password: "" },
  });

  const onSubmit = async (values: ChangePasswordFormValues) => {
    if (!user) return;
    setSubmitting(true);
    setError(null);
    try {
      await userService.changePassword(user.id, values.current_password, values.new_password);
      // Recarrega a sessão: o servidor já desligou a troca obrigatória e o app libera as demais telas.
      await refreshUser();
      feedback.success("Senha alterada com sucesso.");
      navigate("/dashboard", { replace: true });
    } catch (err) {
      setError(getApiErrorMessage(err, "Não foi possível alterar a senha. Tente novamente."));
    } finally {
      setSubmitting(false);
    }
  };

  const handleLogout = async () => {
    await logout();
    navigate("/login", { replace: true });
  };

  return (
    <div>
      <Title level={3} style={{ marginTop: 0, marginBottom: 6 }}>
        Troque sua senha
      </Title>
      <Paragraph type="secondary" style={{ marginBottom: 24 }}>
        Por segurança, defina uma nova senha antes de continuar. Use a senha que você recebeu como senha atual.
      </Paragraph>

      {error && <Alert type="error" showIcon message={error} style={{ marginBottom: 20 }} />}

      <form onSubmit={handleSubmit(onSubmit)} noValidate>
        <FormField label="Senha atual" error={errors.current_password}>
          <Controller
            name="current_password"
            control={control}
            render={({ field }) => (
              <Input.Password {...field} size="large" autoComplete="current-password" autoFocus />
            )}
          />
        </FormField>

        <FormField label="Nova senha" error={errors.new_password}>
          <Controller
            name="new_password"
            control={control}
            render={({ field }) => <Input.Password {...field} size="large" autoComplete="new-password" />}
          />
        </FormField>

        <FormField label="Confirmar nova senha" error={errors.confirm_password}>
          <Controller
            name="confirm_password"
            control={control}
            render={({ field }) => <Input.Password {...field} size="large" autoComplete="new-password" />}
          />
        </FormField>

        <Button type="primary" htmlType="submit" size="large" block loading={submitting} style={{ marginTop: 8 }}>
          Alterar senha
        </Button>
      </form>

      <Paragraph style={{ textAlign: "center", marginTop: 20, marginBottom: 0 }}>
        <Button type="link" onClick={() => void handleLogout()}>
          Sair
        </Button>
      </Paragraph>
    </div>
  );
}
