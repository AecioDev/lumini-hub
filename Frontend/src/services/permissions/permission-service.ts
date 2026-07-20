import { api } from "../common/api";
import type { ApiResponse } from "@/types/common";
import type {
  ApiPermissionDetail,
  ApiPermissionListPaginated,
  ApiPermissionsByModule,
  CreatePermissionRequest,
  GetPermissionsFilters,
  UpdatePermissionRequest,
} from "@/types/permission";

export const permissionService = {
  async list(filters: GetPermissionsFilters = {}): Promise<ApiPermissionListPaginated> {
    const { data } = await api.get<ApiResponse<ApiPermissionListPaginated>>(
      "/permissions",
      { params: filters }
    );
    return data.data!;
  },

  async getById(id: number): Promise<ApiPermissionDetail> {
    const { data } = await api.get<ApiResponse<ApiPermissionDetail>>(
      `/permissions/${id}`
    );
    return data.data!;
  },

  async byModule(): Promise<ApiPermissionsByModule[]> {
    const { data } = await api.get<ApiResponse<ApiPermissionsByModule[]>>(
      "/permissions/by-module"
    );
    return data.data ?? [];
  },

  async modules(): Promise<string[]> {
    const { data } = await api.get<ApiResponse<string[]>>("/permissions/modules");
    return data.data ?? [];
  },

  async create(payload: CreatePermissionRequest) {
    const { data } = await api.post<ApiResponse<ApiPermissionDetail>>(
      "/permissions",
      payload
    );
    return data.data!;
  },

  async update(id: number, payload: UpdatePermissionRequest) {
    const { data } = await api.put<ApiResponse<ApiPermissionDetail>>(
      `/permissions/${id}`,
      payload
    );
    return data.data!;
  },

  async remove(id: number): Promise<void> {
    await api.delete(`/permissions/${id}`);
  },
};
