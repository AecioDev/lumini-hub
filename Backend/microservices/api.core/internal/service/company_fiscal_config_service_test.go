package service

import (
	"errors"
	"os"
	"testing"
	"time"

	"lumini-hub/api.core/internal/domain"
	"lumini-hub/common/utils"
)

// Configuração fiscal da empresa e certificado A1. Precisam de Postgres: sem TEST_DATABASE_URL são
// ignorados. O certificado de testdata/ é autoassinado e só existe para os testes (CN
// "EMPRESA TESTE LTDA:12345678000199", senha abaixo, validade de 100 anos).

const (
	chaveDeTeste     = "000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f" // AES-256 em hex
	senhaCertificado = "senha-de-teste"
)

func novoServicoFiscal(t *testing.T) (*CompanyFiscalConfigService, *CompanyService) {
	t.Helper()
	uow := newTestUow(t)
	return NewCompanyFiscalConfigService(uow, chaveDeTeste), NewCompanyService(uow)
}

func TestFiscal_CriacaoExigeEmpresaComCNPJ(t *testing.T) {
	fiscal, empresas := novoServicoFiscal(t)
	matriz := criarMatriz(t, empresas, "Matriz", "11111111000111")
	semCNPJ, err := empresas.CreateCompany(domain.CreateCompanyRequest{ParentID: &matriz, LegalName: "Filial Sem CNPJ"}, 1)
	if err != nil {
		t.Fatalf("criar filial: %v", err)
	}

	_, err = fiscal.CreateCompanyFiscalConfig(domain.CreateCompanyFiscalConfigRequest{CompanyID: 9999, TaxRegime: "SIMPLES"})
	exigeErroNoCampo(t, err, "company_id")

	_, err = fiscal.CreateCompanyFiscalConfig(domain.CreateCompanyFiscalConfigRequest{CompanyID: semCNPJ.ID, TaxRegime: "SIMPLES"})
	exigeErroNoCampo(t, err, "company_id")

	_, err = fiscal.CreateCompanyFiscalConfig(domain.CreateCompanyFiscalConfigRequest{CompanyID: matriz, TaxRegime: "MEI"})
	exigeErroNoCampo(t, err, "tax_regime")
}

func TestFiscal_UmaConfiguracaoPorEmpresa(t *testing.T) {
	fiscal, empresas := novoServicoFiscal(t)
	matriz := criarMatriz(t, empresas, "Matriz", "11111111000111")

	cfg, err := fiscal.CreateCompanyFiscalConfig(domain.CreateCompanyFiscalConfigRequest{
		CompanyID: matriz, TaxRegime: "PRESUMIDO", AccountantName: "Contabilidade Silva",
	})
	if err != nil {
		t.Fatalf("criar: %v", err)
	}
	if cfg.HasCertificate {
		t.Error("configuração nova não tem certificado")
	}

	_, err = fiscal.CreateCompanyFiscalConfig(domain.CreateCompanyFiscalConfigRequest{CompanyID: matriz, TaxRegime: "REAL"})
	if msg := exigeErroNoCampo(t, err, "company_id"); msg != "esta empresa já possui configuração fiscal cadastrada" {
		t.Errorf("mensagem = %q", msg)
	}

	porEmpresa, err := fiscal.GetByCompanyID(matriz)
	if err != nil || porEmpresa.TaxRegime != "PRESUMIDO" || porEmpresa.AccountantName != "Contabilidade Silva" {
		t.Errorf("GetByCompanyID = %+v, %v", porEmpresa, err)
	}
	if _, err := fiscal.GetByCompanyID(9999); !errors.Is(err, utils.ErrNotFound) {
		t.Errorf("empresa sem configuração: erro = %v, quer ErrNotFound", err)
	}
}

func TestFiscal_Atualizacao(t *testing.T) {
	fiscal, empresas := novoServicoFiscal(t)
	matriz := criarMatriz(t, empresas, "Matriz", "11111111000111")
	cfg, err := fiscal.CreateCompanyFiscalConfig(domain.CreateCompanyFiscalConfigRequest{CompanyID: matriz, TaxRegime: "SIMPLES"})
	if err != nil {
		t.Fatalf("criar: %v", err)
	}

	_, err = fiscal.UpdateCompanyFiscalConfig(cfg.ID, domain.UpdateCompanyFiscalConfigRequest{TaxRegime: "OUTRO"})
	exigeErroNoCampo(t, err, "tax_regime")

	_, err = fiscal.UpdateCompanyFiscalConfig(9999, domain.UpdateCompanyFiscalConfigRequest{TaxRegime: "REAL"})
	exigeErroNoCampo(t, err, "id")

	atualizada, err := fiscal.UpdateCompanyFiscalConfig(cfg.ID, domain.UpdateCompanyFiscalConfigRequest{
		TaxRegime: "REAL", AccountantName: "Nova Contabilidade", AccountantDocument: "12345678000190", AccountantContact: "contato@x.com",
	})
	if err != nil {
		t.Fatalf("atualizar: %v", err)
	}
	if atualizada.TaxRegime != "REAL" || atualizada.AccountantName != "Nova Contabilidade" {
		t.Errorf("atualizada = %+v", atualizada)
	}
}

func TestFiscal_CertificadoA1(t *testing.T) {
	uow := newTestUow(t)
	fiscal := NewCompanyFiscalConfigService(uow, chaveDeTeste)
	empresas := NewCompanyService(uow)
	matriz := criarMatriz(t, empresas, "Matriz", "11111111000111")
	cfg, err := fiscal.CreateCompanyFiscalConfig(domain.CreateCompanyFiscalConfigRequest{CompanyID: matriz, TaxRegime: "SIMPLES"})
	if err != nil {
		t.Fatalf("criar: %v", err)
	}

	arquivo, err := os.ReadFile("testdata/certificado_teste.pfx")
	if err != nil {
		t.Fatalf("lendo o certificado de teste: %v", err)
	}

	// senha errada e arquivo que não é um certificado: erro no campo "password", sem detalhe interno
	_, err = fiscal.SetCertificate(cfg.ID, arquivo, "senha-errada")
	exigeErroNoCampo(t, err, "password")
	_, err = fiscal.SetCertificate(cfg.ID, []byte("isto não é um pfx"), senhaCertificado)
	exigeErroNoCampo(t, err, "password")

	// o certificado não pode ter sido gravado por causa das tentativas que falharam
	if semCert, _ := fiscal.GetByCompanyID(matriz); semCert.HasCertificate {
		t.Fatal("tentativas que falharam não podem gravar certificado")
	}

	_, err = fiscal.SetCertificate(9999, arquivo, senhaCertificado)
	if !errors.Is(err, utils.ErrNotFound) {
		t.Errorf("configuração inexistente: erro = %v, quer ErrNotFound", err)
	}

	dto, err := fiscal.SetCertificate(cfg.ID, arquivo, senhaCertificado)
	if err != nil {
		t.Fatalf("SetCertificate: %v", err)
	}
	if !dto.HasCertificate {
		t.Error("HasCertificate deveria ser true")
	}
	if dto.CertificateSubjectName != "EMPRESA TESTE LTDA" || dto.CertificateSubjectDocument != "12345678000199" {
		t.Errorf("titular = %q / %q, quer EMPRESA TESTE LTDA / 12345678000199", dto.CertificateSubjectName, dto.CertificateSubjectDocument)
	}
	if dto.CertificateExpiry == nil || !dto.CertificateExpiry.After(time.Now().AddDate(50, 0, 0)) {
		t.Errorf("vencimento = %v, quer daqui a mais de 50 anos (vem do próprio arquivo)", dto.CertificateExpiry)
	}

	// a senha nunca fica em texto puro no banco, e a chave certa a recupera
	gravada, err := uow.CompanyFiscalConfigs().FindByID(cfg.ID)
	if err != nil || gravada == nil {
		t.Fatalf("FindByID: %v", err)
	}
	if gravada.CertificatePassword == "" || gravada.CertificatePassword == senhaCertificado {
		t.Fatalf("senha gravada = %q: deveria estar criptografada", gravada.CertificatePassword)
	}
	aberta, err := utils.DecryptAES(gravada.CertificatePassword, chaveDeTeste)
	if err != nil || aberta != senhaCertificado {
		t.Errorf("DecryptAES = %q, %v, quer a senha original", aberta, err)
	}
	if len(gravada.CertificateFile) != len(arquivo) {
		t.Errorf("arquivo gravado com %d bytes, quer %d", len(gravada.CertificateFile), len(arquivo))
	}
}
