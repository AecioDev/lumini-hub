import type { ApiPermission } from "./permission";

// Espelha Backend/microservices/api.auth/internal/domain/role.go

export interface ApiRole {
  id: number;
  name: string;
  description: string;
}

export interface ApiRoleDetail extends ApiRole {
  permissions: ApiPermission[];
  created_at: string;
  updated_at: string;
}

export interface CreateRoleRequest {
  name: string;
  description?: string;
}

export interface UpdateRoleRequest {
  name?: string;
  description?: string;
}

export interface UpdateRolePermissionsRequest {
  permission_ids: number[];
}
