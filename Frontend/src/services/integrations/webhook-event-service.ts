import api from "@/services/common/api";
import { WebhookEventList } from "./webhook-event-schema";

const WebhookEventService = {
  async list(
    page = 1,
    limit = 10,
    filters?: { source?: string; event_type?: string; status?: string }
  ): Promise<WebhookEventList> {
    const response = await api.post("/integrations/webhook-events/filter", {
      page_no: page,
      page_size: limit,
      source: filters?.source,
      event_type: filters?.event_type,
      status: filters?.status,
    });

    if (!response.data.success) {
      throw new Error(
        response.data.error ||
          response.data.message ||
          "Erro ao obter eventos de webhook."
      );
    }

    return response.data.data;
  },
};

export default WebhookEventService;
