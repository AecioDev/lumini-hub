import { Form } from "antd";
import type { ReactNode } from "react";
import type { FieldError } from "react-hook-form";

interface FormFieldProps {
  label?: string;
  error?: FieldError;
  required?: boolean;
  children: ReactNode;
  // Sem <Form> ancestral definindo layout, o Form.Item padrão vira flex-row
  // e o label só quebra pra cima do controle quando o conteúdo é largo o
  // bastante pra forçar isso — funciona por acidente num form de coluna
  // única (Input/Select a 100% da largura), mas fica inconsistente num form
  // de várias colunas lado a lado (ex.: Row/Col de 2-3 campos por linha),
  // onde cada campo tem menos espaço e "decide" sozinho se quebra ou não.
  // Regra combinada com o usuário (2026-09-06): label ao lado quando o
  // campo ocupa a coluna inteira do form; label em cima quando o form tem
  // mais de uma coluna por linha. `stacked` força o segundo caso de forma
  // confiável, com um respiro pequeno (não usa labelCol/wrapperCol do
  // Form.Item, que teria um padding maior que o desejado).
  stacked?: boolean;
}

// Wrapper fino em cima de Form.Item para exibir label + mensagem de erro do
// react-hook-form de forma consistente em todos os formulários da aplicação.
export function FormField({ label, error, required, children, stacked }: FormFieldProps) {
  if (stacked) {
    return (
      <div style={{ marginBottom: 16 }}>
        {label && (
          <div style={{ marginBottom: 4, fontSize: 14 }}>
            {required && <span style={{ color: "#ff4d4f", marginRight: 4 }}>*</span>}
            {label}
          </div>
        )}
        <Form.Item
          validateStatus={error ? "error" : ""}
          help={error?.message}
          style={{ marginBottom: 0 }}
        >
          {children}
        </Form.Item>
      </div>
    );
  }

  return (
    <Form.Item
      label={label}
      required={required}
      validateStatus={error ? "error" : ""}
      help={error?.message}
      style={{ marginBottom: 16 }}
      // labelCol com flex fixo (em vez de deixar o Form.Item quebrar de
      // forma "acidental" conforme a largura do conteúdo) garante duas
      // coisas: o label nunca mais decide sozinho se fica ao lado ou em
      // cima, e — pedido do usuário — os campos de um form de coluna única
      // alinham em coluna (o valor sempre começa na mesma posição
      // horizontal, não importa se o label é "CNPJ" ou "Razão Social").
      labelCol={{ flex: "150px" }}
      labelAlign="left"
      wrapperCol={{ flex: 1 }}
    >
      {children}
    </Form.Item>
  );
}
