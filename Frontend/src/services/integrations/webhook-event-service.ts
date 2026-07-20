import { api } from "../common/api";
import type { ApiResponse } from "@/types/common";
import type {
  ApiWebhookEventListPaginated,
  WebhookEventFilterRequest,
} from "@/types/integration";

export const webhookEventService = {
  async filter(
    payload: WebhookEventFilterRequest = {}
  ): Promise<ApiWebhookEventListPaginated> {
    const { data } = await api.post<ApiResponse<ApiWebhookEventListPaginated>>(
      "/integrations/webhook-events/filter",
      payload
    );
    return data.data!;
  },
};
