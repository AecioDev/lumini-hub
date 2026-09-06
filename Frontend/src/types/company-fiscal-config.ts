// Espelha Backend/microservices/api.core/internal/domain/company_fiscal_config.go

export type TaxRegime = "SIMPLES" | "PRESUMIDO" | "REAL";

export interface ApiCompanyFiscalConfig {
  id: number;
  company_id: number;
  has_certificate: boolean;
  // certificate_expiry e os campos certificate_subject_* são extraídos do
  // próprio arquivo no upload (ver SetCertificate no backend) — nunca
  // digitados pelo usuário.
  certificate_expiry: string | null;
  certificate_subject_name: string;
  certificate_subject_document: string;
  tax_regime: TaxRegime;
  accountant_name: string;
  accountant_document: string;
  accountant_contact: string;
  created_at: string;
  updated_at: string;
}

export type ApiCompanyFiscalConfigDetail = ApiCompanyFiscalConfig;

export interface CreateCompanyFiscalConfigRequest {
  company_id: number;
  tax_regime: TaxRegime;
  accountant_name?: string;
  accountant_document?: string;
  accountant_contact?: string;
}

export interface UpdateCompanyFiscalConfigRequest {
  tax_regime: TaxRegime;
  accountant_name?: string;
  accountant_document?: string;
  accountant_contact?: string;
}
