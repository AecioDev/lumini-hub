import type { ApiPagination } from "./common";
import type { ApiUser } from "./auth";

// Espelha Backend/microservices/api.auth/internal/domain/user.go (DTOs de entrada) e handlers/users.go

export interface CreateUserRequest {
  username: string;
  // Obrigatória sem convite; com `send_invite` o usuário define a própria senha pelo link do e-mail.
  password?: string;
  name: string;
  email: string;
  phone?: string;
  role_id: number;
  company_id?: number | null;
  send_invite?: boolean;
  // Só sem convite: exige a troca da senha no primeiro acesso.
  must_change_password?: boolean;
}

// Resposta do cadastro: o usuário e, quando houve convite, se o e-mail saiu (falha no e-mail não desfaz o cadastro).
export interface ApiUserCreated extends ApiUser {
  email_sent: boolean;
  email_error?: string;
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
