import type { ApiPermission } from "./permission";
import type { ApiUserMenuItem } from "./menu";

// Espelha Backend/microservices/api.auth/internal/domain/user.go

export interface ApiUser {
  id: number;
  username: string;
  name: string;
  email?: string;
  role_id: number;
  role?: string;
  is_active: boolean;
  company_id: number | null;
}

export interface ApiUserRole {
  id: number;
  name: string;
  description: string;
  created_at: string;
  updated_at: string;
}

export interface ApiUserDetail {
  id: number;
  username: string;
  name: string;
  email?: string;
  phone: string;
  role_id: number;
  role: ApiUserRole;
  is_active: boolean;
  last_login?: string;
  company_id: number | null;
  // active_company_id/requires_company_selection só vêm revalidados de
  // verdade nas respostas de login/refresh/me — ver AuthService no backend.
  active_company_id: number | null;
  requires_company_selection?: boolean;
  permissions?: ApiPermission[];
  menu_items?: ApiUserMenuItem[];
  created_at: string;
  updated_at: string;
}

export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginSuccessResponse {
  user: ApiUserDetail;
}
