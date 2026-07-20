import { z } from "zod";

// Login real é por `username`, não e-mail (Backend/microservices/api.auth/internal/handlers/auth.go LoginRequest).
export const loginSchema = z.object({
  username: z.string().min(1, "Informe o usuário"),
  password: z.string().min(1, "Informe a senha"),
  remember: z.boolean().optional(),
});

export type LoginFormValues = z.infer<typeof loginSchema>;
