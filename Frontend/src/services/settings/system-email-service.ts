import { api } from "../common/api";
import type { ApiResponse } from "@/types/common";
import type {
  ApiSystemEmailSettings,
  ApiTestEmailResult,
  SaveSmtpRequest,
} from "@/types/email-config";

export const systemEmailService = {
  async get(): Promise<ApiSystemEmailSettings> {
    const { data } = await api.get<ApiResponse<ApiSystemEmailSettings>>("/system-settings/email");
    return data.data!;
  },

  async save(payload: SaveSmtpRequest): Promise<ApiSystemEmailSettings> {
    const { data } = await api.put<ApiResponse<ApiSystemEmailSettings>>(
      "/system-settings/email",
      payload
    );
    return data.data!;
  },

  // Usa a configuração já salva, não os campos em edição.
  async sendTest(to: string): Promise<ApiTestEmailResult> {
    const { data } = await api.post<ApiResponse<ApiTestEmailResult>>(
      "/system-settings/email/test",
      { to }
    );
    return data.data!;
  },
};
