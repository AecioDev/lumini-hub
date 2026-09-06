// Espelha Backend/microservices/api.core/internal/domain/company.go

export interface ApiCompany {
  id: number;
  parent_id: number | null;
  legal_name: string;
  trade_name: string;
  // null = sem CNPJ próprio — só permitido pra empresa vinculada (parent_id
  // != null) cuja parte fiscal fica a cargo da Matriz (decidido 2026-09-05).
  tax_id: string | null;
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
  // String vazia = sem CNPJ (só válido se parent_id != null) — backend
  // normaliza pra null internamente, ver CompanyValidator.
  tax_id: string;
}

export interface UpdateCompanyRequest extends CreateCompanyRequest {
  is_active: boolean;
}
