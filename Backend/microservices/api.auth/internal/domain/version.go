package domain

// ApiVersionInfo é a resposta de GET /auth/version: em que versão o sistema e o banco estão.
type ApiVersionInfo struct {
	SystemVersion           string `json:"system_version"`
	DatabaseVersion         int    `json:"database_version"`
	RequiredDatabaseVersion int    `json:"required_database_version"`
	Environment             string `json:"environment"`
}
