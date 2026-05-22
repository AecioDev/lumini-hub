// src/services/auth/permission-schema.ts
import * as z from "zod";
import { Pagination } from "../common/pagination-service";

// Schema para a criação de uma nova permissão (Input do formulário)
export const createPermissionSchema = z.object({
  permission: z
    .string()
    .min(3, "O nome da permissão deve ter no mínimo 3 caracteres.")
    .max(100, "O nome da permissão não pode exceder 100 caracteres.")
    .regex(
      /^[a-z0-9_.-]+$/,
      "Formato inválido. Use apenas letras minúsculas, números, '.', '_' ou '-'. Ex: 'module.action'"
    ),
  description: z
    .string()
    .min(5, "A descrição deve ter no mínimo 5 caracteres.")
    .max(255, "A descrição não pode exceder 255 caracteres."),
  module: z
    .string()
    .min(2, "O módulo deve ter no mínimo 2 caracteres.")
    .max(50, "O módulo não pode exceder 50 caracteres."),
});

// Tipagem para os dados de entrada
export type CreatePermissionFormData = z.infer<typeof createPermissionSchema>;

// Schema para a edição de uma nova permissão
export const updatePermissionSchema = z.object({
  permission: z
    .string()
    .min(3, "O nome da permissão deve ter no mínimo 3 caracteres.")
    .max(100, "O nome da permissão não pode exceder 100 caracteres.")
    .regex(
      /^[a-z0-9_.-]+$/,
      "Formato inválido. Use apenas letras minúsculas, números, '.', '_' ou '-'. Ex: 'module.action'"
    ),
  description: z
    .string()
    .min(5, "A descrição deve ter no mínimo 5 caracteres.")
    .max(255, "A descrição não pode exceder 255 caracteres."),
  module: z
    .string()
    .min(2, "O módulo deve ter no mínimo 2 caracteres.")
    .max(50, "O módulo não pode exceder 50 caracteres."),
});

export type UpdatePermissionFormData = z.infer<typeof updatePermissionSchema>;

// Schema para o modelo de Permissão que vem do backend (Output da API)
export const permissionSchema = z.object({
  id: z.string(),
  createdAt: z.string(),
  updatedAt: z.string(),
  deletedAt: z.string().nullable(),
  permission: z.string(),
  description: z.string(),
  module: z.string(),
});

// Tipagem para a permissão vinda do backend
export type Permission = z.infer<typeof permissionSchema>;

// Schema para a resposta padrão da sua API (utils.Response no Go)
export const apiResponseSchema = z.object({
  success: z.boolean(),
  message: z.string().optional(),
  data: z.any().optional(), // Pode ser tipado mais especificamente dependendo do endpoint
  error: z.string().optional(),
  meta: z.any().optional(), // Pode ser tipado mais especificamente
});

export type ApiResponse<T> = z.infer<typeof apiResponseSchema> & {
  data?: T; // Overwrite para tipar 'data' de forma mais específica
};

export interface PermissionList {
  data: Permission[];
  pagination: Pagination; // defina Pagination conforme seu DTO PaginationDTO do backend
}
