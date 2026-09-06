import { api } from "../common/api";
import type { ApiResponse } from "@/types/common";
import type {
  ApiCompanyColorPalette,
  CreateCompanyColorPaletteRequest,
} from "@/types/company-color-palette";

export const companyColorPaletteService = {
  async listByCompanyId(companyId: number): Promise<ApiCompanyColorPalette[]> {
    const { data } = await api.get<ApiResponse<ApiCompanyColorPalette[]>>(
      `/company-color-palettes/by-company/${companyId}`
    );
    return data.data ?? [];
  },

  async create(
    payload: CreateCompanyColorPaletteRequest
  ): Promise<ApiCompanyColorPalette> {
    const { data } = await api.post<ApiResponse<ApiCompanyColorPalette>>(
      "/company-color-palettes",
      payload
    );
    return data.data!;
  },

  async remove(id: number): Promise<void> {
    await api.delete(`/company-color-palettes/${id}`);
  },
};
