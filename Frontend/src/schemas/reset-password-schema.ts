import { z } from "zod";

// Mesma regra de senha do cadastro de usuário (6 a 72 caracteres; 72 é o limite do bcrypt).
export const resetPasswordSchema = z
  .object({
    new_password: z.string().min(6, "Mínimo 6 caracteres").max(72, "Máximo 72 caracteres"),
    confirm_password: z.string().min(1, "Confirme a nova senha"),
  })
  .refine((values) => values.new_password === values.confirm_password, {
    path: ["confirm_password"],
    message: "As senhas não conferem",
  });

export type ResetPasswordFormValues = z.infer<typeof resetPasswordSchema>;
