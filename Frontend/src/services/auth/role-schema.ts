// src/services/auth/role-schema.ts
import { z } from "zod";
import { permissionSchema } from "./permission-schema";

// 1. Schema para o que vem do Backend (Resposta da API - já em camelCase)
// Este tipo representa um Perfil completo, como ele é retornado pelo backend
export const roleSchema = z.object({
  id: z.number().int(),
  name: z.string(),
  description: z.string().nullable().optional(),
  permissions: z.array(permissionSchema).optional(),
  createdAt: z.string(),
  updatedAt: z.string(),
});

export type Role = z.infer<typeof roleSchema>;

// 2. Schema para os dados de Criação de Perfil (Payload para o Backend)
// Usado no formulário de "Novo Perfil"
export const createRoleFormSchema = z.object({
  name: z.string().min(1, "O nome do perfil é obrigatório."),
  description: z.string().optional(),
  // Para criação, enviamos apenas os IDs das permissões selecionadas
  permissionIds: z.array(z.number().int()).optional(), // Array de IDs de permissões
});

export type CreateRoleFormData = z.infer<typeof createRoleFormSchema>;

// 3. Schema para os dados de Atualização de Perfil (Payload para o Backend)
// Usado no formulário de "Editar Perfil"
// Geralmente é o mesmo do de criação, pois os campos editáveis são os mesmos.
export const updateRoleFormSchema = z.object({
  name: z.string().min(1, "O nome do perfil é obrigatório."),
  description: z.string().optional(),
  permissionIds: z.array(z.number().int()).optional(),
});

export type UpdateRoleFormData = z.infer<typeof updateRoleFormSchema>;

// 4. Schema para a associação RolePermission (se precisar usar diretamente)
// Geralmente não é um payload direto, mas pode ser útil para tipagem de dados de tabela de junção
export const rolePermissionSchema = z.object({
  roleId: z.number().int(),
  permissionId: z.number().int(),
});

export type RolePermission = z.infer<typeof rolePermissionSchema>;
