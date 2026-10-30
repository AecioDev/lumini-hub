import { z } from "zod";

export const roleSchema = z.object({
  name: z
    .string()
    .min(1, "Informe o nome do perfil")
    .max(50, "Máximo 50 caracteres"),
  description: z.string().optional(),
});

export type RoleFormValues = z.infer<typeof roleSchema>;
