import api from "@/services/common/api";
import {
  IntegrationSettings,
  UpdateIntegrationSettingsData,
} from "./integration-config-schema";

const IntegrationConfigService = {
  async getSettings(): Promise<IntegrationSettings> {
    const response = await api.get("/integrations/settings");
    if (!response.data.success) {
      throw new Error(
        response.data.error ||
          response.data.message ||
          "Erro ao obter configurações."
      );
    }
    return response.data.data;
  },

  async updateSettings(
    data: UpdateIntegrationSettingsData
  ): Promise<IntegrationSettings> {
    const response = await api.put("/integrations/settings", data);
    if (!response.data.success) {
      throw new Error(
        response.data.error ||
          response.data.message ||
          "Erro ao atualizar configurações."
      );
    }
    return response.data.data;
  },
};

export default IntegrationConfigService;
