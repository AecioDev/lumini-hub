import { api } from "../common/api";
import type { ApiResponse } from "@/types/common";
import type { LegacyCompany, LegacyLocation } from "@/types/integration";

// Endpoints retornam 503 quando a conexão com o SQL Server legado não está
// configurada — quem chama deve tratar essa falha separadamente do restante
// da tela (ver IntegrationsPage.tsx).
export const legacyLookupService = {
  async getLocations(): Promise<LegacyLocation[]> {
    const { data } = await api.get<ApiResponse<LegacyLocation[]>>(
      "/integrations/legacy/locations"
    );
    return data.data ?? [];
  },

  async getCompanies(): Promise<LegacyCompany[]> {
    const { data } = await api.get<ApiResponse<LegacyCompany[]>>(
      "/integrations/legacy/companies"
    );
    return data.data ?? [];
  },
};
