// Espelha Backend/microservices/api.core/internal/domain/company_color_palette.go

export interface ApiCompanyColorPalette {
  id: number;
  company_id: number;
  name: string;
  primary_color: string;
  secondary_color: string | null;
  accent_color: string | null;
  created_at: string;
}

export interface CreateCompanyColorPaletteRequest {
  company_id: number;
  name: string;
  primary_color: string;
  secondary_color?: string;
  accent_color?: string;
}
