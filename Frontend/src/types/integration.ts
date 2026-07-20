import type { ApiPagination } from "./common";

// Espelha Backend/microservices/api.integrations/internal/domain/integration_config.go

export interface ApiIntegrationSettings {
  li_api_key: string;
  li_app_key: string;
  li_webhook_secret: string;
  codtipnot: string;
  codlocarm_oficial: string;
  codlocarm_reserva: string;
  codemp: string;
}

export interface UpdateIntegrationSettingsRequest {
  li_api_key?: string;
  li_app_key?: string;
  li_webhook_secret?: string;
  codtipnot?: string;
  codlocarm_oficial?: string;
  codlocarm_reserva?: string;
  codemp?: string;
}

// Espelha .../internal/domain/sync_log.go

export type SyncDirection = "li_to_sql" | "sql_to_li";
export type SyncStatus = "success" | "error" | "pending";

export interface ApiSyncLog {
  id: number;
  direction: SyncDirection;
  entity_type: string;
  reference_id: string;
  status: SyncStatus;
  message: string;
  details: unknown;
  created_at: string;
}

export interface SyncLogFilterRequest {
  direction?: string;
  entity_type?: string;
  status?: string;
  page_no?: number;
  page_size?: number;
  order_by_column?: string;
  is_asc?: boolean;
}

export interface ApiSyncLogListPaginated {
  data: ApiSyncLog[];
  pagination: ApiPagination;
}

// Espelha .../internal/domain/webhook_event.go

export type WebhookStatus = "received" | "processed" | "error";

export interface ApiWebhookEvent {
  id: number;
  source: string;
  event_type: string;
  payload: unknown;
  status: WebhookStatus;
  received_at: string;
  processed_at: string | null;
  error_message: string;
}

export interface WebhookEventFilterRequest {
  source?: string;
  event_type?: string;
  status?: string;
  page_no?: number;
  page_size?: number;
  order_by_column?: string;
  is_asc?: boolean;
}

export interface ApiWebhookEventListPaginated {
  data: ApiWebhookEvent[];
  pagination: ApiPagination;
}

// Espelha .../internal/legacy/domain/{loc_arm_model,cencus}.go — consultas de apoio ao
// ERP legado (SQL Server), usadas só pros selects da tela de configuração. Podem não estar
// disponíveis (503) se a conexão com o SQL Server legado não estiver configurada.

export interface LegacyLocation {
  cod_loc_am: number;
  des_loc_am: string;
  v_mais_local_empresa: number;
  tip_estoque: string;
  tipe_rc: string;
}

export interface LegacyCompany {
  cod_cencus: number;
  des_cencus: string;
  cgc_cencus: string;
}
