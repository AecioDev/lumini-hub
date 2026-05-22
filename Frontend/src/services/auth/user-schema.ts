import { Pagination } from "../common/pagination-service";
import { Role } from "./role-schema";

export interface User {
  id: number;
  name: string;
  username: string;
  email: string;
  phone?: string;
  role_id: number;
  role: Role;
  is_active: boolean;
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
