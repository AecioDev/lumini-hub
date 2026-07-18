import api from "@/services/common/api";
import { LegacyCompany, LegacyLocation } from "./legacy-lookup-schema";

const LegacyLookupService = {
  async getLocations(): Promise<LegacyLocation[]> {
    const response = await api.get("/integrations/legacy/locations");
    if (!response.data.success) {
      throw new Error(
        response.data.error ||
          response.data.message ||
          "Erro ao obter locais de armazenamento do ERP legado."
      );
    }
    return response.data.data;
  },

  async getCompanies(): Promise<LegacyCompany[]> {
    const response = await api.get("/integrations/legacy/companies");
    if (!response.data.success) {
      throw new Error(
        response.data.error ||
          response.data.message ||
          "Erro ao obter empresas do ERP legado."
      );
    }
    return response.data.data;
  },
};

export default LegacyLookupService;
