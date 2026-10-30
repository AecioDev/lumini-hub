import type { ApiPagination } from "./common";

// Espelha Backend/microservices/api.auth/internal/domain/permission.go

export interface ApiPermission {
  id: number;
  permission: string; // código no formato "modulo.acao", ex: "users.view"
  description: string;
  module: string;
}

export interface ApiPermissionDetail extends ApiPermission {
  created_at: string;
  updated_at: string;
}

export interface ApiPermissionsByModule {
  module: string;
  permissions: ApiPermission[];
}

export interface ApiPermissionListPaginated {
  data: ApiPermission[];
  pagination: ApiPagination;
}

export interface CreatePermissionRequest {
  permission: string;
  description: string;
  module: string;
}

export interface UpdatePermissionRequest {
  permission?: string;
  description?: string;
  module?: string;
}

// Espelha InGetPermissionsFilters — vai como query string (GET), não POST /filter
// como Customers/Suppliers.
export interface GetPermissionsFilters {
  page?: number;
  limit?: number;
  sort?: string;
  name?: string;
  module?: string;
  roleId?: number;
  isLinkedToRole?: boolean;
}

// Espelha PermissionFilterRequest (POST /permissions/filter). `search` procura no código e na descrição.
export interface PermissionFilterRequest {
  search?: string;
  module?: string;
  role_id?: number;
  is_linked_to_role?: boolean;
  page_no?: number;
  page_size?: number;
  order_by_column?: string;
  is_asc?: boolean;
}
