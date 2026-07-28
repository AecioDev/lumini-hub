// Espelha Backend/microservices/api.core/internal/domain/company.go

export interface ApiCompany {
  id: number;
  parent_id: number | null;
  legal_name: string;
  trade_name: string;
  tax_id: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface ApiCompanyDetail extends ApiCompany {
  created_by?: { id: number; username: string; name: string };
  updated_by?: { id: number; username: string; name: string };
}

export interface CreateCompanyRequest {
  parent_id: number | null;
  legal_name: string;
  trade_name?: string;
  tax_id: string;
}

export interface UpdateCompanyRequest extends CreateCompanyRequest {
  is_active: boolean;
}
