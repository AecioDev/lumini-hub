import type { ApiPagination } from "@/types/common";

// Espelha domain.ApiAuditLog (Backend/microservices/api.auth/internal/domain/audit_log.go)
export interface ApiAuditLog {
  id: number;
  created_at: string;
  user_id: number | null;
  username: string;
  /** Método + modelo da rota, ex.: "PUT /api/users/:id". */
  action: string;
  entity_type: string;
  entity_id: string;
  ip_address: string;
  status: number;
  /** O que mudou: `changes` (campo -> de/para, ou adicionados/removidos numa lista), `after` (objeto criado), `before` (objeto excluído) ou `truncated` (conteúdo grande demais). */
  details?: {
    changes?: Record<string, unknown>;
    after?: Record<string, unknown>;
    before?: Record<string, unknown>;
    truncated?: string[];
  };
}

export interface ApiAuditLogListPaginated {
  data: ApiAuditLog[];
  pagination?: ApiPagination;
}

// Espelha AuditLogFilterRequest (POST /audit-logs/filter). `search` procura na ação e no nome do
// usuário; `date_from` (inclusivo) e `date_to` (exclusivo) são instantes ISO 8601.
export interface AuditLogFilterRequest {
  search?: string;
  entity_type?: string;
  method?: "POST" | "PUT" | "PATCH" | "DELETE";
  user_id?: number;
  date_from?: string;
  date_to?: string;
  page_no?: number;
  page_size?: number;
  order_by_column?: string;
  is_asc?: boolean;
}
