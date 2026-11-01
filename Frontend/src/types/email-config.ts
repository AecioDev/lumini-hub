// Espelha Backend/microservices/api.core/internal/domain/email_config.go

export type SmtpSecurity = "starttls" | "ssl" | "none";

// A senha nunca volta da API: só se ela existe.
export interface ApiSmtpSettings {
  host: string;
  port: number;
  security: SmtpSecurity;
  username: string;
  has_password: boolean;
  from_name: string;
  from_email: string;
  enabled: boolean;
}

export interface ApiSystemEmailSettings extends ApiSmtpSettings {
  updated_at: string;
}

// "company" = a própria empresa, "parent" = herdada de uma empresa acima (Matriz), "system" = padrão do sistema.
export type EmailSource = "company" | "parent" | "system";

export interface ApiCompanyEmailConfig extends ApiSmtpSettings {
  company_id: number;
  // false = a empresa não tem servidor próprio e herda.
  configured: boolean;
  // De onde o envio sairia hoje; vazio quando não há nenhuma configuração ativa.
  effective_source: EmailSource | "";
  effective_name: string;
  updated_at: string;
}

// Senha vazia mantém a que já está gravada.
export interface SaveSmtpRequest {
  host: string;
  port: number;
  security: SmtpSecurity;
  username: string;
  password: string;
  from_name: string;
  from_email: string;
  enabled: boolean;
}

export interface ApiTestEmailResult {
  source: EmailSource;
  name: string;
}
