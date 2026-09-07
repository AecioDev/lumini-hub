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

export interface ApiCompanyOption {
  id: number;
  name: string;
}

// Espelha utils.CompanyVisualConfigOption (Backend/common/utils/company_visibility.go)
// — identidade visual da empresa ativa, sem depender de companies.visual_config.view.
export interface ApiCompanyVisualConfigOption {
  id: number;
  has_logo: boolean;
  primary_color: string;
  secondary_color: string | null;
  accent_color: string | null;
  text_color: string | null;
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
  // active_company_id/active_company_name/requires_company_selection/
  // visible_companies só vêm revalidados de verdade nas respostas de
  // login/refresh/me — ver AuthService no backend. visible_companies não
  // depende de companies.view (é identidade de sessão, não administração
  // do cadastro) — sempre reflete o que o próprio usuário pode ver.
  active_company_id: number | null;
  active_company_name?: string;
  requires_company_selection?: boolean;
  visible_companies?: ApiCompanyOption[];
  active_company_visual_config?: ApiCompanyVisualConfigOption;
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
