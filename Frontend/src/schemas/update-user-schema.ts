import { z } from "zod";

// Espelha UpdateUserRequest (Backend/microservices/api.auth/internal/domain/user.go).
// Sem senha: troca de senha é um fluxo separado (PUT /users/:id/password), fora do escopo desta tela.
export const updateUserSchema = z.object({
  name: z.string().min(1, "Informe o nome"),
  email: z.string().min(1, "Informe o e-mail").email("E-mail inválido"),
  phone: z.string().optional(),
  role_id: z
    .number({ invalid_type_error: "Selecione um perfil" })
    .int()
    .positive("Selecione um perfil"),
  is_active: z.boolean(),
});

export type UpdateUserFormValues = z.infer<typeof updateUserSchema>;
