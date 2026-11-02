import { api } from "../common/api";
import type { ApiResponse } from "@/types/common";
import type { ApiUserDetail } from "@/types/auth";
import type {
  ApiUserCreated,
  ApiUserListPaginated,
  CreateUserRequest,
  UpdateUserPermissionsRequest,
  UpdateUserRequest,
  UserFilterRequest,
} from "@/types/user";

export interface ListUsersParams {
  page?: number;
  limit?: number;
  sort?: string; // ex: "username asc"
}

export const userService = {
  async list(params: ListUsersParams = {}): Promise<ApiUserListPaginated> {
    const { data } = await api.get<ApiResponse<ApiUserListPaginated>>("/users", {
      params,
    });
    return data.data!;
  },

  async filter(payload: UserFilterRequest): Promise<ApiUserListPaginated> {
    const { data } = await api.post<ApiResponse<ApiUserListPaginated>>("/users/filter", payload);
    return data.data!;
  },

  async getById(id: number): Promise<ApiUserDetail> {
    const { data } = await api.get<ApiResponse<ApiUserDetail>>(`/users/${id}`);
    return data.data!;
  },

  async create(payload: CreateUserRequest): Promise<ApiUserCreated> {
    const { data } = await api.post<ApiResponse<ApiUserCreated>>("/users", payload);
    return data.data!;
  },

  // Gera um novo link de convite (o anterior deixa de valer) e reenvia o e-mail; só para quem nunca entrou.
  async resendInvite(id: number): Promise<void> {
    await api.post(`/users/${id}/invite`);
  },

  // O próprio usuário troca a senha informando a atual.
  async changePassword(id: number, currentPassword: string, newPassword: string): Promise<void> {
    await api.put(`/users/${id}/password`, {
      current_password: currentPassword,
      new_password: newPassword,
    });
  },

  async update(id: number, payload: UpdateUserRequest): Promise<ApiUserDetail> {
    const { data } = await api.put<ApiResponse<ApiUserDetail>>(
      `/users/${id}`,
      payload
    );
    return data.data!;
  },

  async remove(id: number): Promise<void> {
    await api.delete(`/users/${id}`);
  },

  async updatePermissions(
    id: number,
    payload: UpdateUserPermissionsRequest
  ): Promise<void> {
    await api.put(`/users/${id}/permissions`, payload);
  },
};
