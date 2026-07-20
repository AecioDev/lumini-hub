import { z } from "zod";

// Espelha InCreatePermission/InUpdatePermission
// (Backend/microservices/api.auth/internal/domain/permission.go, tags `binding`).
export const permissionSchema = z.object({
  permission: z
    .string()
    .min(3, "Mínimo 3 caracteres")
    .max(100, "Máximo 100 caracteres"),
  description: z
    .string()
    .min(3, "Mínimo 3 caracteres")
    .max(500, "Máximo 500 caracteres"),
  module: z.string().min(3, "Mínimo 3 caracteres").max(50, "Máximo 50 caracteres"),
});

export type PermissionFormValues = z.infer<typeof permissionSchema>;
