import { api } from "../common/api";
import type { ApiResponse } from "@/types/common";
import type { ApiAuditLogListPaginated, AuditLogFilterRequest } from "@/types/audit";

export const auditService = {
  async filter(payload: AuditLogFilterRequest): Promise<ApiAuditLogListPaginated> {
    const { data } = await api.post<ApiResponse<ApiAuditLogListPaginated>>("/audit-logs/filter", payload);
    return data.data!;
  },

  async entityTypes(): Promise<string[]> {
    const { data } = await api.get<ApiResponse<string[]>>("/audit-logs/entity-types");
    return data.data ?? [];
  },
};
