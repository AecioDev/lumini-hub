package service

import (
	"errors"
	"strings"
	"testing"

	"lumini-hub/api.core/internal/domain"
	"lumini-hub/common/mail"
)

// Configuração de e-mail: sistema, empresa e a regra de resolução empresa -> Matriz -> sistema. Precisam de
// Postgres: sem TEST_DATABASE_URL são ignorados.

// remetenteFalso guarda o que seria enviado, no lugar de abrir conexão SMTP.
type remetenteFalso struct {
	enviados []mail.ServerConfig
	err      error
}

func (r *remetenteFalso) Send(cfg mail.ServerConfig, _ mail.Message) error {
	if r.err != nil {
		return r.err
	}
	r.enviados = append(r.enviados, cfg)
	return nil
}

func novoServicoEmail(t *testing.T) (*MailConfigService, *CompanyService, *remetenteFalso) {
	t.Helper()
	uow := newTestUow(t)
	fake := &remetenteFalso{}
	return NewMailConfigService(uow, fake, chaveDeTeste), NewCompanyService(uow), fake
}

func smtpLigado(host string) domain.SaveSMTPRequest {
	return domain.SaveSMTPRequest{Host: host, Port: 587, Security: "starttls", Username: "user", Password: "segredo-" + host, FromEmail: "noreply@" + host, Enabled: true}
}

func filial(t *testing.T, empresas *CompanyService, pai uint, nome string) uint {
	t.Helper()
	c, err := empresas.CreateCompany(domain.CreateCompanyRequest{ParentID: &pai, LegalName: nome}, 1)
	if err != nil {
		t.Fatalf("criar filial %q: %v", nome, err)
	}
	return c.ID
}

func TestEmail_SistemaVazioNaoEnviaENaoVazaSenha(t *testing.T) {
	svc, _, _ := novoServicoEmail(t)

	got, err := svc.GetSystemEmail()
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Enabled || got.HasPassword || got.Port != 587 || got.Security != "starttls" {
		t.Errorf("padrão inesperado: %+v", got)
	}

	_, err = svc.Resolve(nil)
	if !errors.Is(err, ErrMailNotConfigured) {
		t.Fatalf("sem configuração deveria ser ErrMailNotConfigured, veio %v", err)
	}
}

func TestEmail_SalvarCriptografaSenhaEVazioMantem(t *testing.T) {
	svc, _, _ := novoServicoEmail(t)

	if _, err := svc.SaveSystemEmail(smtpLigado("smtp.exemplo.com"), 1); err != nil {
		t.Fatalf("salvar: %v", err)
	}
	stored, _ := svc.uow.SystemSettings().Get()
	if stored.SMTP.Password == "" || strings.Contains(stored.SMTP.Password, "segredo") {
		t.Fatalf("a senha deveria estar criptografada no banco: %q", stored.SMTP.Password)
	}

	// Salvar de novo sem senha mantém a gravada; a resolução devolve a senha em texto para o envio.
	req := smtpLigado("smtp.exemplo.com")
	req.Password = ""
	api, err := svc.SaveSystemEmail(req, 1)
	if err != nil {
		t.Fatalf("salvar sem senha: %v", err)
	}
	if !api.HasPassword {
		t.Error("a senha deveria continuar gravada")
	}
	resolved, err := svc.Resolve(nil)
	if err != nil || resolved.Server.Password != "segredo-smtp.exemplo.com" || resolved.Source != MailSourceSystem {
		t.Errorf("resolve = %+v, %v", resolved, err)
	}

	// Sem usuário não há autenticação: a senha é descartada.
	req.Username = ""
	api, err = svc.SaveSystemEmail(req, 1)
	if err != nil || api.HasPassword {
		t.Errorf("sem usuário a senha deveria sumir: %+v, %v", api, err)
	}
}

func TestEmail_LigadoExigeServidorRemetenteESenha(t *testing.T) {
	svc, _, _ := novoServicoEmail(t)

	_, err := svc.SaveSystemEmail(domain.SaveSMTPRequest{Enabled: true}, 1)
	exigeErroNoCampo(t, err, "host")

	_, err = svc.SaveSystemEmail(domain.SaveSMTPRequest{Host: "h", Enabled: true, FromEmail: "x"}, 1)
	exigeErroNoCampo(t, err, "from_email")

	_, err = svc.SaveSystemEmail(domain.SaveSMTPRequest{Host: "h", Enabled: true, FromEmail: "a@b.com", Username: "u"}, 1)
	exigeErroNoCampo(t, err, "password")

	_, err = svc.SaveSystemEmail(domain.SaveSMTPRequest{Host: "h", Security: "tls", FromEmail: "a@b.com"}, 1)
	exigeErroNoCampo(t, err, "security")

	// Desligado aceita campos pela metade.
	if _, err := svc.SaveSystemEmail(domain.SaveSMTPRequest{Host: "so-o-servidor"}, 1); err != nil {
		t.Errorf("desligado deveria aceitar: %v", err)
	}
}

func TestEmail_EmpresaUsaAPropriaDepoisAMatrizDepoisOSistema(t *testing.T) {
	svc, empresas, _ := novoServicoEmail(t)
	matriz := criarMatriz(t, empresas, "Matriz", "11111111000111")
	filha := filial(t, empresas, matriz, "Filial")

	if _, err := svc.SaveSystemEmail(smtpLigado("sistema.com"), 1); err != nil {
		t.Fatal(err)
	}
	// Só o sistema configurado: a filial cai nele.
	r, err := svc.Resolve(&filha)
	if err != nil || r.Source != MailSourceSystem {
		t.Fatalf("esperava sistema, veio %+v, %v", r, err)
	}

	// A Matriz configura: a filial herda dela.
	if _, err := svc.SaveCompanyEmail(matriz, smtpLigado("matriz.com")); err != nil {
		t.Fatal(err)
	}
	r, err = svc.Resolve(&filha)
	if err != nil || r.Source != MailSourceParent || r.Server.Host != "matriz.com" || r.Name != "Matriz" {
		t.Fatalf("esperava herdar da Matriz, veio %+v, %v", r, err)
	}

	// A filial configura a própria.
	if _, err := svc.SaveCompanyEmail(filha, smtpLigado("filial.com")); err != nil {
		t.Fatal(err)
	}
	r, _ = svc.Resolve(&filha)
	if r.Source != MailSourceCompany || r.Server.Host != "filial.com" {
		t.Fatalf("esperava a própria, veio %+v", r)
	}

	// Desligar a da filial devolve à Matriz; apagar também.
	off := smtpLigado("filial.com")
	off.Enabled = false
	if _, err := svc.SaveCompanyEmail(filha, off); err != nil {
		t.Fatal(err)
	}
	if r, _ = svc.Resolve(&filha); r.Server.Host != "matriz.com" {
		t.Errorf("filial desligada deveria herdar, veio %+v", r)
	}
	if err := svc.DeleteCompanyEmail(filha); err != nil {
		t.Fatal(err)
	}
	api, _ := svc.GetCompanyEmail(filha)
	if api.Configured || api.EffectiveSource != MailSourceParent || api.EffectiveName != "Matriz" {
		t.Errorf("depois de apagar: %+v", api)
	}
}

func TestEmail_NuncaCaiEmEmpresaForaDoGrupo(t *testing.T) {
	svc, empresas, _ := novoServicoEmail(t)
	grupoA := criarMatriz(t, empresas, "Grupo A", "11111111000111")
	grupoB := criarMatriz(t, empresas, "Grupo B", "22222222000122")
	filhaA := filial(t, empresas, grupoA, "Filial A")

	// Só o grupo B tem e-mail; o sistema está vazio.
	if _, err := svc.SaveCompanyEmail(grupoB, smtpLigado("grupob.com")); err != nil {
		t.Fatal(err)
	}

	_, err := svc.Resolve(&filhaA)
	if !errors.Is(err, ErrMailNotConfigured) {
		t.Fatalf("a filial do grupo A não pode usar o grupo B: %v", err)
	}
	if msg := err.Error(); !strings.Contains(msg, "Filial A") || !strings.Contains(msg, "Grupo A") || strings.Contains(msg, "Grupo B") {
		t.Errorf("mensagem deveria citar a cadeia da empresa: %q", msg)
	}

	// Sem empresa, só vale o sistema (e o grupo B não entra).
	if _, err := svc.Resolve(nil); !errors.Is(err, ErrMailNotConfigured) {
		t.Errorf("sem empresa e sem sistema deveria falhar: %v", err)
	}
}

func TestEmail_EnvioDeTesteUsaConfiguracaoEfetivaEPropagaFalha(t *testing.T) {
	svc, empresas, fake := novoServicoEmail(t)
	matriz := criarMatriz(t, empresas, "Matriz", "11111111000111")
	if _, err := svc.SaveCompanyEmail(matriz, smtpLigado("matriz.com")); err != nil {
		t.Fatal(err)
	}
	filha := filial(t, empresas, matriz, "Filial")

	res, err := svc.SendTest(&filha, "ana@exemplo.com")
	if err != nil || res.Source != MailSourceParent {
		t.Fatalf("teste: %+v, %v", res, err)
	}
	if len(fake.enviados) != 1 || fake.enviados[0].Host != "matriz.com" || fake.enviados[0].Password != "segredo-matriz.com" {
		t.Errorf("enviado = %+v", fake.enviados)
	}

	fake.err = errors.New("535 autenticação recusada")
	_, err = svc.SendTest(&filha, "ana@exemplo.com")
	var sendErr *MailSendError
	if !errors.As(err, &sendErr) {
		t.Errorf("falha do servidor deveria ser MailSendError: %v", err)
	}

	if _, err := svc.SendTest(uintPtr(9999), "ana@exemplo.com"); err == nil {
		t.Error("empresa inexistente deveria falhar")
	}
}
