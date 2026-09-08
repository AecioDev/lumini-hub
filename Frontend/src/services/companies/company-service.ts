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

// URL direta pro binário do logo (fora do envelope ApiResponse) — usada
// como src de <img>, o navegador manda o cookie de sessão sozinho (mesmo
// domínio efetivo do resto da API, ver services/common/api.ts). Mesmo
// padrão de companyVisualConfigLogoUrl em company-visual-config-service.ts
// (arquivo antigo, removido em CFG-7.3.1).
export function companyLogoUrl(companyId: number): string {
  return `${import.meta.env.VITE_API_URL}/companies/${companyId}/logo`;
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

  // Chamadas de baixo nível pro logo — quem decide QUANDO chamá-las é
  // CompanyForm.tsx, que trata upload/remoção como "rascunho" local (nada é
  // enviado até o clique em "Salvar Empresa", junto com o resto dos dados
  // da empresa; ver handleFormSubmit lá — redesenhado em 2026-09-07 depois
  // de um bug real na versão anterior "ação isolada e imediata").
  async uploadLogo(id: number, file: File): Promise<ApiCompany> {
    const formData = new FormData();
    formData.append("logo", file);

    const { data } = await api.post<ApiResponse<ApiCompany>>(`/companies/${id}/logo`, formData, {
      headers: { "Content-Type": "multipart/form-data" },
    });
    return data.data!;
  },

  async removeLogo(id: number): Promise<ApiCompany> {
    const { data } = await api.delete<ApiResponse<ApiCompany>>(`/companies/${id}/logo`);
    return data.data!;
  },
};
