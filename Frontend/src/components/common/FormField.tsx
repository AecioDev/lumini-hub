import { Form } from "antd";
import type { ReactNode } from "react";
import type { FieldError } from "react-hook-form";

interface FormFieldProps {
  label?: string;
  error?: FieldError;
  required?: boolean;
  children: ReactNode;
}

// Wrapper fino em cima de Form.Item para exibir label + mensagem de erro do
// react-hook-form de forma consistente em todos os formulários da aplicação.
export function FormField({ label, error, required, children }: FormFieldProps) {
  return (
    <Form.Item
      label={label}
      required={required}
      validateStatus={error ? "error" : ""}
      help={error?.message}
      style={{ marginBottom: 16 }}
    >
      {children}
    </Form.Item>
  );
}
