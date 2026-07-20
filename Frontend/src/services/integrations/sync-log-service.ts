import { api } from "../common/api";
import type { ApiResponse } from "@/types/common";
import type { ApiSyncLogListPaginated, SyncLogFilterRequest } from "@/types/integration";

export const syncLogService = {
  async filter(payload: SyncLogFilterRequest = {}): Promise<ApiSyncLogListPaginated> {
    const { data } = await api.post<ApiResponse<ApiSyncLogListPaginated>>(
      "/integrations/sync-logs/filter",
      payload
    );
    return data.data!;
  },
};
