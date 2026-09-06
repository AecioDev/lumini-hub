import { api } from "../common/api";
import { isNotFoundError } from "@/utils/api-error";
import type { ApiResponse } from "@/types/common";
import type {
  ApiCompanyVisualConfigDetail,
  CreateCompanyVisualConfigRequest,
  UpdateCompanyVisualConfigRequest,
} from "@/types/company-visual-config";

// URL direta pro binário do logo (fora do envelope ApiResponse) — usada
// como src de <img>, o navegador manda o cookie de sessão sozinho (mesmo
// domínio efetivo do resto da API, ver services/common/api.ts).
export function companyVisualConfigLogoUrl(configId: number): string {
  return `${import.meta.env.VITE_API_URL}/company-visual-configs/${configId}/logo`;
}

export const companyVisualConfigService = {
  // Devolve null (não lança) quando a empresa ainda não tem Configuração
  // Visual cadastrada — a tela usa isso pra decidir entre criar ou editar.
  async getByCompanyId(companyId: number): Promise<ApiCompanyVisualConfigDetail | null> {
    try {
      const { data } = await api.get<ApiResponse<ApiCompanyVisualConfigDetail>>(
        `/company-visual-configs/by-company/${companyId}`
      );
      return data.data ?? null;
    } catch (error) {
      if (isNotFoundError(error)) return null;
      throw error;
    }
  },

  async create(
    payload: CreateCompanyVisualConfigRequest
  ): Promise<ApiCompanyVisualConfigDetail> {
    const { data } = await api.post<ApiResponse<ApiCompanyVisualConfigDetail>>(
      "/company-visual-configs",
      payload
    );
    return data.data!;
  },

  async update(
    id: number,
    payload: UpdateCompanyVisualConfigRequest
  ): Promise<ApiCompanyVisualConfigDetail> {
    const { data } = await api.put<ApiResponse<ApiCompanyVisualConfigDetail>>(
      `/company-visual-configs/${id}`,
      payload
    );
    return data.data!;
  },

  // Sem `remove()` de propósito — config 1:1 da empresa, mesmo princípio de
  // company-fiscal-config-service.ts.

  // Upload multipart separado do PUT de cores.
  async uploadLogo(id: number, file: File): Promise<ApiCompanyVisualConfigDetail> {
    const formData = new FormData();
    formData.append("logo", file);

    const { data } = await api.post<ApiResponse<ApiCompanyVisualConfigDetail>>(
      `/company-visual-configs/${id}/logo`,
      formData,
      { headers: { "Content-Type": "multipart/form-data" } }
    );
    return data.data!;
  },
};
