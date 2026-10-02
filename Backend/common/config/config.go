package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config armazena todas as configurações da aplicação
type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	SQLServer SQLServerConfig
	JWT       JWTConfig
	Security  SecurityConfig
	Bootstrap BootstrapAdminConfig
	App       AppConfig
}

// AppConfig armazena configurações gerais da aplicação
type AppConfig struct {
	Env string // Ex: "development", "production", "test"
}

// ServerConfig armazena configurações do servidor HTTP
type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// DatabaseConfig armazena configurações do banco de dados
type DatabaseConfig struct {
	Host         string
	Port         string
	User         string
	Password     string
	DBName       string
	SSLMode      string
	DatabaseLink string
	// AutoMigrate liga/desliga o AutoMigrate dos serviços ao subir (DB_AUTO_MIGRATE).
	// Ligado por padrão pra um banco vazio subir sozinho; em produção com dados
	// reais o operador desliga explicitamente (DB_AUTO_MIGRATE=false).
	AutoMigrate bool
}

// SQLServerConfig armazena configurações do banco de dados SQL Server
type SQLServerConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// JWTConfig armazena configurações do JWT
type JWTConfig struct {
	Secret          string
	AccessTokenExp  time.Duration
	RefreshTokenExp time.Duration
}

// SecurityConfig armazena chaves usadas pra criptografar dados sensíveis
// reversíveis (diferente de hash irreversível, ex. senha de usuário/bcrypt).
type SecurityConfig struct {
	// CertificateEncryptionKey é uma chave AES-256 em hex (64 caracteres),
	// usada pra criptografar/decriptografar a senha do certificado digital
	// A1 de uma empresa (ver utils.EncryptAES/DecryptAES). Gerar com
	// `openssl rand -hex 32`.
	CertificateEncryptionKey string
}

// BootstrapAdminConfig armazena as credenciais do primeiro usuário ADMIN, criado
// só no primeiro boot (banco sem usuários). Transitório: no multi-tenant quem
// provisiona o primeiro admin é o app comercial. Não há senha default de
// propósito — sem BOOTSTRAP_ADMIN_PASSWORD o admin simplesmente não é criado.
type BootstrapAdminConfig struct {
	Username string
	Password string
	Email    string
}

// Load carrega as configurações do ambiente
func Load() (*Config, error) {
	// Carregar variáveis de ambiente do arquivo .env se existir
	// Tenta buscar no diretório atual, um nível acima ou dois níveis acima (útil para microsserviços)
	_ = godotenv.Load()
	_ = godotenv.Load("../.env")
	_ = godotenv.Load("../../.env")
	_ = godotenv.Load("../../../.env")

	// Configurações do servidor
	port := getEnv("SERVER_PORT", "4000")
	readTimeout, _ := strconv.Atoi(getEnv("SERVER_READ_TIMEOUT", "10"))
	writeTimeout, _ := strconv.Atoi(getEnv("SERVER_WRITE_TIMEOUT", "10"))
	idleTimeout, _ := strconv.Atoi(getEnv("SERVER_IDLE_TIMEOUT", "60"))

	// Configurações do banco de dados
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "erp_system")
	dbSSLMode := getEnv("DB_SSLMODE", "disable")
	dbLink := getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/erp_system?sslmode=disable")
	dbAutoMigrate, err := getEnvBool("DB_AUTO_MIGRATE", true)
	if err != nil {
		return nil, err
	}

	// Configurações do SQL Server
	mssqlHost := getEnv("MSSQL_HOST", "localhost")
	mssqlPort := getEnv("MSSQL_PORT", "1433")
	mssqlUser := getEnv("MSSQL_USER", "sa")
	mssqlPassword := getEnv("MSSQL_PASSWORD", "")
	mssqlName := getEnv("MSSQL_NAME", "FOCCO_ERP")

	// Configurações do JWT
	jwtSecret := getEnv("JWT_SECRET", "your-secret-key")
	jwtAccessExp, _ := strconv.Atoi(getEnv("JWT_ACCESS_EXP", "15"))      // 15 minutos
	jwtRefreshExp, _ := strconv.Atoi(getEnv("JWT_REFRESH_EXP", "10080")) // 7 dias

	// Chave de criptografia de dados sensíveis reversíveis
	certificateEncryptionKey := getEnv("CERTIFICATE_ENCRYPTION_KEY", "")

	// Primeiro usuário ADMIN (sem default de senha)
	bootstrapAdmin := BootstrapAdminConfig{
		Username: getEnv("BOOTSTRAP_ADMIN_USERNAME", "admin"),
		Password: getEnv("BOOTSTRAP_ADMIN_PASSWORD", ""),
		Email:    getEnv("BOOTSTRAP_ADMIN_EMAIL", ""),
	}

	// Configurações gerais da aplicação
	appEnv := getEnv("APP_ENV", "development")

	return &Config{
		Server: ServerConfig{
			Port:         port,
			ReadTimeout:  time.Duration(readTimeout) * time.Second,
			WriteTimeout: time.Duration(writeTimeout) * time.Second,
			IdleTimeout:  time.Duration(idleTimeout) * time.Second,
		},
		Database: DatabaseConfig{
			Host:         dbHost,
			Port:         dbPort,
			User:         dbUser,
			Password:     dbPassword,
			DBName:       dbName,
			SSLMode:      dbSSLMode,
			DatabaseLink: dbLink,
			AutoMigrate:  dbAutoMigrate,
		},
		SQLServer: SQLServerConfig{
			Host:     mssqlHost,
			Port:     mssqlPort,
			User:     mssqlUser,
			Password: mssqlPassword,
			DBName:   mssqlName,
		},
		JWT: JWTConfig{
			Secret:          jwtSecret,
			AccessTokenExp:  time.Duration(jwtAccessExp) * time.Minute,
			RefreshTokenExp: time.Duration(jwtRefreshExp) * time.Minute,
		},
		Security: SecurityConfig{
			CertificateEncryptionKey: certificateEncryptionKey,
		},
		Bootstrap: bootstrapAdmin,
		App: AppConfig{
			Env: appEnv,
		},
	}, nil
}

// DSN retorna a string de conexão com o banco de dados
func (c *DatabaseConfig) DSN() string {
	if c.DatabaseLink != "" {
		return c.DatabaseLink
	}

	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

// DSN retorna a string de conexão para SQL Server
func (c *SQLServerConfig) DSN() string {
	query := url.Values{}
	query.Add("database", c.DBName)

	u := &url.URL{
		Scheme:   "sqlserver",
		User:     url.UserPassword(c.User, c.Password),
		Host:     fmt.Sprintf("%s:%s", c.Host, c.Port),
		RawQuery: query.Encode(),
	}
	return u.String()
}

// getEnv retorna o valor da variável de ambiente ou o valor padrão
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvBool lê uma variável booleana. Valor inválido é erro (em vez de cair
// no default em silêncio): um "DB_AUTO_MIGRATE=flase" digitado errado não pode
// deixar o AutoMigrate ligado em produção sem ninguém perceber.
func getEnvBool(key string, defaultValue bool) (bool, error) {
	value, exists := os.LookupEnv(key)
	if !exists || strings.TrimSpace(value) == "" {
		return defaultValue, nil
	}
	parsed, err := strconv.ParseBool(strings.TrimSpace(value))
	if err != nil {
		return false, fmt.Errorf("%s inválido (%q): use true ou false", key, value)
	}
	return parsed, nil
}
