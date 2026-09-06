package domain

import (
	"time"

	"gorm.io/gorm"
)

// CompanyFiscalConfig representa a configuração fiscal de uma Company (1:1):
// certificado digital A1, tipo de tributação e dados do contador responsável.
// Ver Documentos/Planejamento/Modulo_0_Configuracao/plano_empresa.md pra
// decisões fechadas (certificado em bytea na própria linha, senha
// criptografada e reversível pro uso no ACBrLib).
//
// Numeração/série de nota fiscal NÃO mora aqui — é CompanyDocumentIssuanceConfig,
// tabela separada (EPIC CFG-2), cuja estrutura fina ainda depende do Módulo 6
// (Fiscal) ser desenhado.
type CompanyFiscalConfig struct {
	gorm.Model

	// Unique parcial (só entre linhas não excluídas) — sem o `where`, o soft
	// delete do GORM (`deleted_at`) deixaria uma linha excluída ainda
	// ocupando o índice único, impedindo recriar a configuração fiscal da
	// mesma empresa depois de excluir a antiga (achado testando CFG-1.1.4).
	CompanyID uint `gorm:"not null;uniqueIndex:idx_company_fiscal_configs_company_id,where:deleted_at IS NULL" json:"company_id"`

	// CertificateFile e CertificatePassword nunca são expostos em JSON — ver
	// ApiCompanyFiscalConfig.HasCertificate. Upload é endpoint multipart
	// separado do PUT de configuração (CFG-1.1.4), não passa por aqui como
	// Create/UpdateRequest.
	CertificateFile []byte `gorm:"type:bytea" json:"-"`
	// CertificatePassword precisa ser reversível (uso futuro no ACBrLib), não
	// hash — a criptografia/decriptografia real ainda não está implementada
	// nesta tarefa (CFG-1.1.1, só domain model); fica pra CFG-1.1.3, que
	// também decide onde mora a chave simétrica.
	CertificatePassword string `gorm:"type:text" json:"-"`

	// CertificateExpiry, CertificateSubjectName e CertificateSubjectDocument
	// são extraídos do próprio arquivo do certificado no upload
	// (SetCertificate) — nunca vêm de Create/UpdateRequest, o usuário não
	// digita essas informações.
	CertificateExpiry          *time.Time `json:"certificate_expiry"`
	CertificateSubjectName     string     `gorm:"size:200" json:"certificate_subject_name"`
	CertificateSubjectDocument string     `gorm:"size:20" json:"certificate_subject_document"`

	// TaxRegime aceita: SIMPLES, PRESUMIDO, REAL. Validação do enum de fato
	// é responsabilidade do validator (CFG-1.1.2), não deste struct.
	TaxRegime string `gorm:"size:20;not null" json:"tax_regime"`

	// Dados do contador — campos exatos sujeitos a evoluir conforme o
	// Módulo 6 (Fiscal) for desenhado, ver plano_empresa.md.
	AccountantName     string `gorm:"size:150" json:"accountant_name"`
	AccountantDocument string `gorm:"size:20" json:"accountant_document"`
	AccountantContact  string `gorm:"size:150" json:"accountant_contact"`
}

// TableName especifica o nome da tabela
func (CompanyFiscalConfig) TableName() string {
	return "company_fiscal_configs"
}

// CreateCompanyFiscalConfigRequest representa os dados para criar a configuração fiscal de uma empresa
type CreateCompanyFiscalConfigRequest struct {
	CompanyID          uint   `json:"company_id" binding:"required"`
	TaxRegime          string `json:"tax_regime" binding:"required"`
	AccountantName     string `json:"accountant_name"`
	AccountantDocument string `json:"accountant_document"`
	AccountantContact  string `json:"accountant_contact"`
}

// UpdateCompanyFiscalConfigRequest representa os dados para atualizar a configuração fiscal de uma empresa
type UpdateCompanyFiscalConfigRequest struct {
	TaxRegime          string `json:"tax_regime" binding:"required"`
	AccountantName     string `json:"accountant_name"`
	AccountantDocument string `json:"accountant_document"`
	AccountantContact  string `json:"accountant_contact"`
}

// ApiCompanyFiscalConfig representa os dados de configuração fiscal para exibição
type ApiCompanyFiscalConfig struct {
	ID                         uint       `json:"id"`
	CompanyID                  uint       `json:"company_id"`
	HasCertificate             bool       `json:"has_certificate"`
	CertificateExpiry          *time.Time `json:"certificate_expiry"`
	CertificateSubjectName     string     `json:"certificate_subject_name"`
	CertificateSubjectDocument string     `json:"certificate_subject_document"`
	TaxRegime                  string     `json:"tax_regime"`
	AccountantName             string     `json:"accountant_name"`
	AccountantDocument         string     `json:"accountant_document"`
	AccountantContact          string     `json:"accountant_contact"`
	CreatedAt                  time.Time  `json:"created_at"`
	UpdatedAt                  time.Time  `json:"updated_at"`
}

// ApiCompanyFiscalConfigDetail representa os dados detalhados de configuração fiscal
type ApiCompanyFiscalConfigDetail struct {
	ID                         uint       `json:"id"`
	CompanyID                  uint       `json:"company_id"`
	HasCertificate             bool       `json:"has_certificate"`
	CertificateExpiry          *time.Time `json:"certificate_expiry"`
	CertificateSubjectName     string     `json:"certificate_subject_name"`
	CertificateSubjectDocument string     `json:"certificate_subject_document"`
	TaxRegime                  string     `json:"tax_regime"`
	AccountantName             string     `json:"accountant_name"`
	AccountantDocument         string     `json:"accountant_document"`
	AccountantContact          string     `json:"accountant_contact"`
	CreatedAt                  time.Time  `json:"created_at"`
	UpdatedAt                  time.Time  `json:"updated_at"`
}

// ApiCompanyFiscalConfigFromModel converte um CompanyFiscalConfig para ApiCompanyFiscalConfig
func ApiCompanyFiscalConfigFromModel(c CompanyFiscalConfig) ApiCompanyFiscalConfig {
	return ApiCompanyFiscalConfig{
		ID:                         c.ID,
		CompanyID:                  c.CompanyID,
		HasCertificate:             len(c.CertificateFile) > 0,
		CertificateExpiry:          c.CertificateExpiry,
		CertificateSubjectName:     c.CertificateSubjectName,
		CertificateSubjectDocument: c.CertificateSubjectDocument,
		TaxRegime:                  c.TaxRegime,
		AccountantName:             c.AccountantName,
		AccountantDocument:         c.AccountantDocument,
		AccountantContact:          c.AccountantContact,
		CreatedAt:                  c.CreatedAt,
		UpdatedAt:                  c.UpdatedAt,
	}
}

// ApiCompanyFiscalConfigDetailFromModel converte um CompanyFiscalConfig para ApiCompanyFiscalConfigDetail
func ApiCompanyFiscalConfigDetailFromModel(c CompanyFiscalConfig) ApiCompanyFiscalConfigDetail {
	return ApiCompanyFiscalConfigDetail{
		ID:                         c.ID,
		CompanyID:                  c.CompanyID,
		HasCertificate:             len(c.CertificateFile) > 0,
		CertificateExpiry:          c.CertificateExpiry,
		CertificateSubjectName:     c.CertificateSubjectName,
		CertificateSubjectDocument: c.CertificateSubjectDocument,
		TaxRegime:                  c.TaxRegime,
		AccountantName:             c.AccountantName,
		AccountantDocument:         c.AccountantDocument,
		AccountantContact:          c.AccountantContact,
		CreatedAt:                  c.CreatedAt,
		UpdatedAt:                  c.UpdatedAt,
	}
}
