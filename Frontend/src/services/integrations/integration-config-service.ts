import { api } from "../common/api";
import type { ApiResponse } from "@/types/common";
import type {
  ApiIntegrationSettings,
  UpdateIntegrationSettingsRequest,
} from "@/types/integration";

export const integrationConfigService = {
  async getSettings(): Promise<ApiIntegrationSettings> {
    const { data } = await api.get<ApiResponse<ApiIntegrationSettings>>(
      "/integrations/settings"
    );
    return data.data!;
  },

  async updateSettings(
    payload: UpdateIntegrationSettingsRequest
  ): Promise<ApiIntegrationSettings> {
    const { data } = await api.put<ApiResponse<ApiIntegrationSettings>>(
      "/integrations/settings",
      payload
    );
    return data.data!;
  },
};
