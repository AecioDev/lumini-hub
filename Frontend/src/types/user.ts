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
}

export interface UpdateUserRequest {
  name?: string;
  email?: string;
  phone?: string;
  role_id?: number;
  is_active?: boolean;
}

export interface UpdateUserPermissionsRequest {
  permission_ids: number[];
}

export interface ApiUserListPaginated {
  users: ApiUser[];
  pagination?: ApiPagination;
}
