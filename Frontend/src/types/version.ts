// Espelha domain.ApiVersionInfo (api.auth), resposta de GET /auth/version.
export interface ApiVersionInfo {
  system_version: string;
  database_version: number;
  required_database_version: number;
  environment: string;
}
