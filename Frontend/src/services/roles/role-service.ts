import { api } from "../common/api";
import type { ApiResponse } from "@/types/common";
import type {
  ApiRole,
  ApiRoleDetail,
  CreateRoleRequest,
  UpdateRoleRequest,
  UpdateRolePermissionsRequest,
} from "@/types/role";

export const roleService = {
  async list(): Promise<ApiRole[]> {
    const { data } = await api.get<ApiResponse<ApiRole[]>>("/roles");
    // GET /roles não garante ordem estável (sem ORDER BY no backend) — sem
    // isso, telas que assumem roles[0] como padrão (ex.: preset inicial no
    // formulário de usuário) mudariam de perfil a cada carregamento.
    return [...(data.data ?? [])].sort((a, b) => a.name.localeCompare(b.name));
  },

  async getById(id: number): Promise<ApiRoleDetail> {
    const { data } = await api.get<ApiResponse<ApiRoleDetail>>(`/roles/${id}`);
    return data.data!;
  },

  async create(payload: CreateRoleRequest): Promise<ApiRoleDetail> {
    const { data } = await api.post<ApiResponse<ApiRoleDetail>>("/roles", payload);
    return data.data!;
  },

  async update(id: number, payload: UpdateRoleRequest): Promise<ApiRoleDetail> {
    const { data } = await api.put<ApiResponse<ApiRoleDetail>>(
      `/roles/${id}`,
      payload
    );
    return data.data!;
  },

  async updatePermissions(
    id: number,
    payload: UpdateRolePermissionsRequest
  ): Promise<void> {
    await api.put(`/roles/${id}/permissions`, payload);
  },

  async remove(id: number): Promise<void> {
    await api.delete(`/roles/${id}`);
  },
};
