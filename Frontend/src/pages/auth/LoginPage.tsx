import { useEffect, useState } from "react";
import { zodResolver } from "@hookform/resolvers/zod";
import { Controller, useForm } from "react-hook-form";
import { App as AntdApp, Button, Checkbox, Input, Typography } from "antd";
import { Link, useNavigate } from "react-router-dom";
import { useAuth } from "@/contexts/AuthContext";
import { FormField } from "@/components/common/FormField";
import { loginSchema, type LoginFormValues } from "@/schemas/login-schema";

const { Title, Paragraph } = Typography;

export function LoginPage() {
  const { login, isAuthenticated } = useAuth();
  const navigate = useNavigate();
  const { message } = AntdApp.useApp();
  const [submitting, setSubmitting] = useState(false);

  const {
    control,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginFormValues>({
    resolver: zodResolver(loginSchema),
    defaultValues: { username: "", password: "", remember: false },
  });

  useEffect(() => {
    if (isAuthenticated) navigate("/dashboard", { replace: true });
  }, [isAuthenticated, navigate]);

  const onSubmit = async (values: LoginFormValues) => {
    setSubmitting(true);
    try {
      await login({ username: values.username, password: values.password });
      navigate("/dashboard", { replace: true });
    } catch {
      message.error("Usuário ou senha inválidos.");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div>
      <Title level={3} style={{ marginTop: 0, marginBottom: 6 }}>
        Bem-vindo de volta
      </Title>
      <Paragraph type="secondary" style={{ marginBottom: 28 }}>
        Entre com sua conta para continuar
      </Paragraph>

      <form onSubmit={handleSubmit(onSubmit)} noValidate>
        <FormField label="Usuário" error={errors.username}>
          <Controller
            name="username"
            control={control}
            render={({ field }) => (
              <Input {...field} size="large" placeholder="seu.usuario" autoFocus />
            )}
          />
        </FormField>

        <FormField label="Senha" error={errors.password}>
          <Controller
            name="password"
            control={control}
            render={({ field }) => (
              <Input.Password {...field} size="large" placeholder="••••••••" />
            )}
          />
        </FormField>

        <div
          style={{
            display: "flex",
            justifyContent: "space-between",
            alignItems: "center",
            margin: "-4px 0 22px",
          }}
        >
          <Controller
            name="remember"
            control={control}
            render={({ field: { value, onChange } }) => (
              <Checkbox checked={value} onChange={(e) => onChange(e.target.checked)}>
                Lembrar-me
              </Checkbox>
            )}
          />
          <Link to="/forgot-password">Esqueceu a senha?</Link>
        </div>

        <Button type="primary" htmlType="submit" size="large" block loading={submitting}>
          Entrar
        </Button>
      </form>
    </div>
  );
}
