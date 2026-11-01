import { api } from "../common/api";
import type { ApiResponse } from "@/types/common";
import type {
  ApiCompanyEmailConfig,
  ApiTestEmailResult,
  SaveSmtpRequest,
} from "@/types/email-config";

export const companyEmailConfigService = {
  async getByCompanyId(companyId: number): Promise<ApiCompanyEmailConfig> {
    const { data } = await api.get<ApiResponse<ApiCompanyEmailConfig>>(
      `/company-email-configs/by-company/${companyId}`
    );
    return data.data!;
  },

  // Cria ou atualiza (PUT único por empresa).
  async save(companyId: number, payload: SaveSmtpRequest): Promise<ApiCompanyEmailConfig> {
    const { data } = await api.put<ApiResponse<ApiCompanyEmailConfig>>(
      `/company-email-configs/by-company/${companyId}`,
      payload
    );
    return data.data!;
  },

  // Apagar a configuração própria devolve a empresa à herança (Matriz, depois sistema).
  async remove(companyId: number): Promise<void> {
    await api.delete(`/company-email-configs/by-company/${companyId}`);
  },

  async sendTest(companyId: number, to: string): Promise<ApiTestEmailResult> {
    const { data } = await api.post<ApiResponse<ApiTestEmailResult>>(
      `/company-email-configs/by-company/${companyId}/test`,
      { to }
    );
    return data.data!;
  },
};
