import api from "@/services/common/api";
import { SyncLogList } from "./sync-log-schema";

const SyncLogService = {
  async list(
    page = 1,
    limit = 10,
    filters?: { direction?: string; entity_type?: string; status?: string }
  ): Promise<SyncLogList> {
    const response = await api.post("/integrations/sync-logs/filter", {
      page_no: page,
      page_size: limit,
      direction: filters?.direction,
      entity_type: filters?.entity_type,
      status: filters?.status,
    });

    if (!response.data.success) {
      throw new Error(
        response.data.error ||
          response.data.message ||
          "Erro ao obter logs de sincronização."
      );
    }

    return response.data.data;
  },
};

export default SyncLogService;
