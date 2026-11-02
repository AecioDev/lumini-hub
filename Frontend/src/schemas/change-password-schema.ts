import { z } from "zod";

// Troca de senha do próprio usuário (PUT /users/:id/password): senha atual + nova, 6 a 72 caracteres.
export const changePasswordSchema = z
  .object({
    current_password: z.string().min(1, "Informe a senha atual"),
    new_password: z.string().min(6, "Mínimo 6 caracteres").max(72, "Máximo 72 caracteres"),
    confirm_password: z.string().min(1, "Confirme a nova senha"),
  })
  .refine((values) => values.new_password === values.confirm_password, {
    path: ["confirm_password"],
    message: "As senhas não conferem",
  })
  .refine((values) => values.new_password !== values.current_password, {
    path: ["new_password"],
    message: "A nova senha precisa ser diferente da atual",
  });

export type ChangePasswordFormValues = z.infer<typeof changePasswordSchema>;
