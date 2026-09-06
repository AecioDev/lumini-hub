import { z } from "zod";

// Espelha as tags `binding` de CreateUserRequest
// (Backend/microservices/api.auth/internal/domain/user.go).
export const createUserSchema = z.object({
  username: z
    .string()
    .min(3, "Mínimo 3 caracteres")
    .max(50, "Máximo 50 caracteres"),
  password: z.string().min(6, "Mínimo 6 caracteres"),
  name: z.string().min(1, "Informe o nome"),
  email: z.string().min(1, "Informe o e-mail").email("E-mail inválido"),
  phone: z.string().optional(),
  role_id: z
    .number({ invalid_type_error: "Selecione um perfil" })
    .int()
    .positive("Selecione um perfil"),
  // null = usuário "master" (sem empresa vinculada, vê todas a partir da Matriz).
  company_id: z.number().int().positive().nullable(),
});

export type CreateUserFormValues = z.infer<typeof createUserSchema>;
