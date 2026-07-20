import { api } from "../common/api";
import type { ApiResponse } from "@/types/common";
import type { ApiUserDetail } from "@/types/auth";
import type {
  ApiUserListPaginated,
  CreateUserRequest,
  UpdateUserPermissionsRequest,
  UpdateUserRequest,
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

  async getById(id: number): Promise<ApiUserDetail> {
    const { data } = await api.get<ApiResponse<ApiUserDetail>>(`/users/${id}`);
    return data.data!;
  },

  async create(payload: CreateUserRequest): Promise<ApiUserDetail> {
    const { data } = await api.post<ApiResponse<ApiUserDetail>>("/users", payload);
    return data.data!;
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
