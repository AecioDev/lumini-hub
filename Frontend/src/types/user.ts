import type { ApiPagination } from "./common";
import type { ApiUser } from "./auth";

// Espelha Backend/microservices/api.auth/internal/domain/user.go (DTOs de entrada) e handlers/users.go

export interface CreateUserRequest {
  username: string;
  password: string;
  name: string;
  email: string;
  phone?: string;
  role_id: number;
  company_id?: number | null;
}

export interface UpdateUserRequest {
  name?: string;
  email?: string;
  phone?: string;
  role_id?: number;
  is_active?: boolean;
  company_id?: number | null;
}

export interface UpdateUserPermissionsRequest {
  permission_ids: number[];
}

export interface ApiUserListPaginated {
  users: ApiUser[];
  pagination?: ApiPagination;
}

// Espelha UserFilterRequest (POST /users/filter). `search` procura no nome, no usuário e no e-mail.
export interface UserFilterRequest {
  search?: string;
  role_id?: number;
  is_active?: boolean;
  page_no?: number;
  page_size?: number;
  order_by_column?: string;
  is_asc?: boolean;
}
