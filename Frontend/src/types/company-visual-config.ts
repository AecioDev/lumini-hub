// Espelha Backend/microservices/api.core/internal/domain/company_visual_config.go

export interface ApiCompanyVisualConfig {
  id: number;
  company_id: number;
  // O logo em si não viaja aqui (é bytea no backend) — ver
  // logoUrl em company-visual-config-service.ts, servido por endpoint
  // dedicado (GET /company-visual-configs/:id/logo).
  has_logo: boolean;
  primary_color: string;
  secondary_color: string | null;
  accent_color: string | null;
  text_color: string | null;
  created_at: string;
  updated_at: string;
}

export type ApiCompanyVisualConfigDetail = ApiCompanyVisualConfig;

export interface CreateCompanyVisualConfigRequest {
  company_id: number;
  primary_color: string;
  secondary_color?: string;
  accent_color?: string;
  text_color?: string;
}

export interface UpdateCompanyVisualConfigRequest {
  primary_color: string;
  secondary_color?: string;
  accent_color?: string;
  text_color?: string;
}
