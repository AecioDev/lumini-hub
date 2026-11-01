package domain

import (
	"time"

	"gorm.io/gorm"
)

// SMTPSettings são as colunas de um servidor SMTP, compartilhadas pela configuração do sistema
// (SystemSettings) e pela de cada empresa (CompanyEmailConfig). A senha fica criptografada (AES, a mesma
// chave do certificado) e nunca sai em JSON: a API só informa HasPassword.
type SMTPSettings struct {
	Host      string `gorm:"column:smtp_host;size:150" json:"host"`
	Port      int    `gorm:"column:smtp_port;default:587" json:"port"`
	Security  string `gorm:"column:smtp_security;size:10;default:starttls" json:"security"`
	Username  string `gorm:"column:smtp_username;size:150" json:"username"`
	Password  string `gorm:"column:smtp_password;type:text" json:"-"`
	FromName  string `gorm:"column:smtp_from_name;size:100" json:"from_name"`
	FromEmail string `gorm:"column:smtp_from_email;size:150" json:"from_email"`
	Enabled   bool   `gorm:"column:smtp_enabled" json:"enabled"`
}

// SystemSettingsID é o ID da única linha de system_settings.
const SystemSettingsID uint = 1

// SystemSettings são os parâmetros globais do sistema (valem para o banco inteiro, qualquer empresa).
// É sempre UMA linha, de ID fixo 1 (a constraint abaixo impede qualquer outra): novos parâmetros globais
// entram como colunas aqui, com o prefixo do assunto (smtp_...), sem tabela nova.
type SystemSettings struct {
	ID uint `gorm:"primaryKey;check:chk_system_settings_singleton,id = 1" json:"id"`

	SMTP SMTPSettings `gorm:"embedded" json:"smtp"`

	UpdatedByID *uint     `gorm:"column:updated_by" json:"updated_by_id"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName especifica o nome da tabela
func (SystemSettings) TableName() string {
	return "system_settings"
}

// CompanyEmailConfig é o servidor SMTP próprio de uma empresa (1:1). Sem ela, a empresa herda a da
// Matriz e, por fim, a do sistema (ver MailConfigService.Resolve). Apagar a linha devolve a empresa à herança.
type CompanyEmailConfig struct {
	gorm.Model

	// Unique parcial pelo mesmo motivo da configuração fiscal: o soft delete não pode prender o índice.
	CompanyID uint `gorm:"not null;uniqueIndex:idx_company_email_configs_company_id,where:deleted_at IS NULL" json:"company_id"`

	SMTP SMTPSettings `gorm:"embedded"`
}

// TableName especifica o nome da tabela
func (CompanyEmailConfig) TableName() string {
	return "company_email_configs"
}

// SaveSMTPRequest é o corpo do PUT de uma configuração SMTP (a do sistema ou a de uma empresa). Salvar
// desligada aceita campos incompletos; ligada, o validator exige servidor, porta e remetente. Password
// vazia mantém a senha já gravada.
type SaveSMTPRequest struct {
	Host      string `json:"host" binding:"omitempty,max=150"`
	Port      int    `json:"port" binding:"omitempty,min=1,max=65535"`
	Security  string `json:"security" binding:"omitempty,max=10"`
	Username  string `json:"username" binding:"omitempty,max=150"`
	Password  string `json:"password" binding:"omitempty,max=255"`
	FromName  string `json:"from_name" binding:"omitempty,max=100"`
	FromEmail string `json:"from_email" binding:"omitempty,max=150"`
	Enabled   bool   `json:"enabled"`
}

// SendTestEmailRequest pede um e-mail de teste para o endereço informado.
type SendTestEmailRequest struct {
	To string `json:"to" binding:"required,email,max=150"`
}

// ApiSMTPSettings é a configuração SMTP para exibição: sem a senha, só se ela existe.
type ApiSMTPSettings struct {
	Host        string `json:"host"`
	Port        int    `json:"port"`
	Security    string `json:"security"`
	Username    string `json:"username"`
	HasPassword bool   `json:"has_password"`
	FromName    string `json:"from_name"`
	FromEmail   string `json:"from_email"`
	Enabled     bool   `json:"enabled"`
}

// ApiSystemEmailSettings é a resposta do e-mail do sistema.
type ApiSystemEmailSettings struct {
	ApiSMTPSettings
	UpdatedAt time.Time `json:"updated_at"`
}

// ApiCompanyEmailConfig é a resposta do e-mail de uma empresa. Configured é false quando a empresa não
// tem configuração própria; EffectiveSource diz de onde o envio sairia hoje ("company", "parent",
// "system" ou vazio se não há nenhuma), para a tela mostrar a herança.
type ApiCompanyEmailConfig struct {
	ApiSMTPSettings
	CompanyID       uint      `json:"company_id"`
	Configured      bool      `json:"configured"`
	EffectiveSource string    `json:"effective_source"`
	EffectiveName   string    `json:"effective_name"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// ApiTestEmailResult diz por qual configuração o e-mail de teste saiu.
type ApiTestEmailResult struct {
	Source string `json:"source"`
	Name   string `json:"name"`
}

// ApiSMTPSettingsFromModel converte as colunas SMTP para exibição
func ApiSMTPSettingsFromModel(s SMTPSettings) ApiSMTPSettings {
	return ApiSMTPSettings{
		Host:        s.Host,
		Port:        s.Port,
		Security:    s.Security,
		Username:    s.Username,
		HasPassword: s.Password != "",
		FromName:    s.FromName,
		FromEmail:   s.FromEmail,
		Enabled:     s.Enabled,
	}
}
