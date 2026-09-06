import { api } from "../common/api";
import { isNotFoundError } from "@/utils/api-error";
import type { ApiResponse } from "@/types/common";
import type {
  ApiCompanyFiscalConfigDetail,
  CreateCompanyFiscalConfigRequest,
  UpdateCompanyFiscalConfigRequest,
} from "@/types/company-fiscal-config";

export const companyFiscalConfigService = {
  // Devolve null (não lança) quando a empresa ainda não tem Configuração
  // Fiscal cadastrada — a tela usa isso pra decidir entre criar ou editar.
  async getByCompanyId(companyId: number): Promise<ApiCompanyFiscalConfigDetail | null> {
    try {
      const { data } = await api.get<ApiResponse<ApiCompanyFiscalConfigDetail>>(
        `/company-fiscal-configs/by-company/${companyId}`
      );
      return data.data ?? null;
    } catch (error) {
      if (isNotFoundError(error)) return null;
      throw error;
    }
  },

  async create(
    payload: CreateCompanyFiscalConfigRequest
  ): Promise<ApiCompanyFiscalConfigDetail> {
    const { data } = await api.post<ApiResponse<ApiCompanyFiscalConfigDetail>>(
      "/company-fiscal-configs",
      payload
    );
    return data.data!;
  },

  async update(
    id: number,
    payload: UpdateCompanyFiscalConfigRequest
  ): Promise<ApiCompanyFiscalConfigDetail> {
    const { data } = await api.put<ApiResponse<ApiCompanyFiscalConfigDetail>>(
      `/company-fiscal-configs/${id}`,
      payload
    );
    return data.data!;
  },

  // Sem `remove()` de propósito — config 1:1 da empresa, ver comentário em
  // CompanyFiscalConfigService.go (backend nem expõe mais o endpoint).

  // Upload multipart separado do PUT de configuração — a senha nunca volta
  // em nenhuma resposta (backend só expõe has_certificate).
  async uploadCertificate(
    id: number,
    file: File,
    password: string
  ): Promise<ApiCompanyFiscalConfigDetail> {
    const formData = new FormData();
    formData.append("certificate", file);
    formData.append("password", password);

    const { data } = await api.post<ApiResponse<ApiCompanyFiscalConfigDetail>>(
      `/company-fiscal-configs/${id}/certificate`,
      formData,
      { headers: { "Content-Type": "multipart/form-data" } }
    );
    return data.data!;
  },
};
