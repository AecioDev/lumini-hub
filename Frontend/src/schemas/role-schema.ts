import { z } from "zod";

export const roleSchema = z.object({
  name: z.string().min(1, "Informe o nome do perfil"),
  description: z.string().optional(),
});

export type RoleFormValues = z.infer<typeof roleSchema>;
