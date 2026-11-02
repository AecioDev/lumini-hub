import { z } from "zod";

// Espelha as tags `binding` de CreateUserRequest
// (Backend/microservices/api.auth/internal/domain/user.go).
export const createUserSchema = z
  .object({
    username: z
      .string()
      .min(3, "Mínimo 3 caracteres")
      .max(50, "Máximo 50 caracteres"),
    // Obrigatória só sem convite (ver superRefine); com convite o usuário define a própria senha pelo link.
    password: z
      .string()
      // O bcrypt recusa senhas acima de 72 bytes.
      .max(72, "Máximo 72 caracteres"),
    name: z.string().min(1, "Informe o nome").max(100, "Máximo 100 caracteres"),
    email: z
      .string()
      .min(1, "Informe o e-mail")
      .email("E-mail inválido")
      .max(100, "Máximo 100 caracteres"),
    phone: z.string().max(20, "Máximo 20 caracteres").optional(),
    role_id: z
      .number({ invalid_type_error: "Selecione um perfil" })
      .int()
      .positive("Selecione um perfil"),
    // null = usuário "master" (sem empresa vinculada, vê todas a partir da Matriz).
    company_id: z.number().int().positive().nullable(),
    // Convite por e-mail: o usuário nasce sem senha e define a sua pelo link recebido.
    send_invite: z.boolean(),
    // Só sem convite: o usuário precisa trocar a senha inicial no primeiro acesso.
    must_change_password: z.boolean(),
  })
  .superRefine((values, ctx) => {
    if (!values.send_invite && values.password.length < 6) {
      ctx.addIssue({
        code: "custom",
        path: ["password"],
        message: values.password.length === 0 ? "Informe a senha inicial" : "Mínimo 6 caracteres",
      });
    }
  });

export type CreateUserFormValues = z.infer<typeof createUserSchema>;
