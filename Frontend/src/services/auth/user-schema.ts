import { z } from "zod";
import { Pagination } from "../common/pagination-service";
import { Role } from "./role-schema";
import { Permission } from "./permission-schema";

export interface MenuItem {
  id: number;
  name: string;
  icon: string;
  href: string;
  children?: MenuItem[];
}

export interface User {
  id: number;
  name: string;
  username: string;
  email: string;
  phone?: string;
  role_id: number;
  role: Role;
  is_active: boolean;
  permissions?: Permission[];
  menu_items?: MenuItem[];
  created_at: string;
  updated_at: string;
}

export interface UserList {
  users: User[];
  pagination: Pagination;
}

export interface CreateUserDto {
  username: string;
  password: string;
  name: string;
  email?: string;
  phone?: string;
  role_id: number;
}

export interface UpdateUserDto {
  name?: string;
  email?: string;
  phone?: string;
  role_id?: number;
  is_active?: boolean;
}

export interface ChangePasswordDto {
  current_password: string;
  new_password: string;
  confirm_password: string;
}

export const createUserFormSchema = z.object({
  username: z
    .string()
    .min(3, "O nome de usuário deve ter no mínimo 3 caracteres.")
    .max(50, "O nome de usuário não pode exceder 50 caracteres."),
  password: z.string().min(6, "A senha deve ter no mínimo 6 caracteres."),
  name: z.string().min(1, "O nome completo é obrigatório."),
  email: z.string().min(1, "O email é obrigatório.").email("Email inválido."),
  phone: z.string().optional(),
  role_id: z.coerce
    .number({ invalid_type_error: "Selecione um perfil." })
    .min(1, "Selecione um perfil."),
});

export type CreateUserFormInput = z.infer<typeof createUserFormSchema>;

export const updateUserFormSchema = z.object({
  name: z.string().min(1, "O nome completo é obrigatório."),
  email: z.string().min(1, "O email é obrigatório.").email("Email inválido."),
  phone: z.string().optional(),
  role_id: z.coerce
    .number({ invalid_type_error: "Selecione um perfil." })
    .min(1, "Selecione um perfil."),
});

export type UpdateUserFormInput = z.infer<typeof updateUserFormSchema>;
