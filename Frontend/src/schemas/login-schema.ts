import { z } from "zod";

// O campo vai no JSON como `username` (Backend/microservices/api.auth/internal/handlers/auth.go LoginRequest),
// mas o backend aceita tanto o nome de usuário quanto o e-mail da conta.
export const loginSchema = z.object({
  username: z.string().min(1, "Informe o usuário ou o e-mail"),
  password: z.string().min(1, "Informe a senha"),
});

export type LoginFormValues = z.infer<typeof loginSchema>;
