import { api } from "../common/api";
import type { ApiResponse } from "@/types/common";
import type {
  ApiCompany,
  ApiCompanyDetail,
  CreateCompanyRequest,
  UpdateCompanyRequest,
} from "@/types/company";

// GetCompanies devolve domain.ApiCompanyList (wrapper com campo `data`), não
// um slice puro como /roles — por isso o duplo `data.data.data` aqui.
interface ApiCompanyListEnvelope {
  data: ApiCompany[];
}

export const companyService = {
  async list(): Promise<ApiCompany[]> {
    const { data } = await api.get<ApiResponse<ApiCompanyListEnvelope>>("/companies");
    return data.data?.data ?? [];
  },

  async getById(id: number): Promise<ApiCompanyDetail> {
    const { data } = await api.get<ApiResponse<ApiCompanyDetail>>(`/companies/${id}`);
    return data.data!;
  },

  async create(payload: CreateCompanyRequest): Promise<ApiCompany> {
    const { data } = await api.post<ApiResponse<ApiCompany>>("/companies", payload);
    return data.data!;
  },

  async update(id: number, payload: UpdateCompanyRequest): Promise<ApiCompany> {
    const { data } = await api.put<ApiResponse<ApiCompany>>(`/companies/${id}`, payload);
    return data.data!;
  },

  async remove(id: number): Promise<void> {
    await api.delete(`/companies/${id}`);
  },
};
