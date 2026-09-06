import { api } from "../common/api";
import type { ApiResponse } from "@/types/common";
import type { ApiUserDetail, LoginRequest, LoginSuccessResponse } from "@/types/auth";

export const authService = {
  async login(payload: LoginRequest): Promise<ApiUserDetail> {
    const { data } = await api.post<ApiResponse<LoginSuccessResponse>>(
      "/auth/login",
      payload
    );
    return data.data!.user;
  },

  async me(): Promise<ApiUserDetail> {
    const { data } = await api.get<ApiResponse<LoginSuccessResponse>>("/auth/me");
    return data.data!.user;
  },

  async refreshToken(): Promise<ApiUserDetail> {
    const { data } = await api.post<ApiResponse<LoginSuccessResponse>>(
      "/auth/refresh-token"
    );
    return data.data!.user;
  },

  async logout(): Promise<void> {
    await api.post("/auth/logout");
  },

  async setActiveCompany(companyId: number): Promise<ApiUserDetail> {
    const { data } = await api.put<ApiResponse<LoginSuccessResponse>>(
      "/auth/active-company",
      { company_id: companyId }
    );
    return data.data!.user;
  },
};
