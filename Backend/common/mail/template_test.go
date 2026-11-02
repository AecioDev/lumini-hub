package mail

import (
	"strings"
	"testing"
)

func TestRender_RedefinicaoDeSenha(t *testing.T) {
	msg, err := Render(TemplatePasswordReset, "ana@example.com", TemplateData{
		CompanyName: "Empresa & Cia", Name: "Ana <script>", Link: "https://app.exemplo.com/reset-password?token=abc_-123", ExpiresMinutes: 30,
	})
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if msg.To != "ana@example.com" || msg.Subject != "Redefinição de senha" {
		t.Errorf("cabeçalho = %q / %q", msg.To, msg.Subject)
	}
	for _, want := range []string{"https://app.exemplo.com/reset-password?token=abc_-123", "30 minutos", "Empresa &amp; Cia"} {
		if !strings.Contains(msg.HTML, want) {
			t.Errorf("HTML sem %q", want)
		}
	}
	if strings.Contains(msg.HTML, "<script>") {
		t.Error("o nome precisa sair escapado no HTML")
	}
	if !strings.Contains(msg.Text, "https://app.exemplo.com/reset-password?token=abc_-123") {
		t.Errorf("texto inesperado:\n%s", msg.Text)
	}
}

func TestRender_TemplateDesconhecido(t *testing.T) {
	if _, err := Render("nao_existe", "a@b.com", TemplateData{}); err == nil {
		t.Error("template desconhecido deveria falhar")
	}
}

func TestRender_ConviteDeUsuario(t *testing.T) {
	msg, err := Render(TemplateUserInvite, "novo@example.com", TemplateData{
		CompanyName: "Focco", Name: "Novo", Username: "novo.usuario", InvitedBy: "Admin",
		Link: "https://app.exemplo.com/reset-password?token=abc&invite=1", ExpiresDays: 7,
	})
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if msg.Subject != "Convite para acessar o Lumini Hub" {
		t.Errorf("assunto = %q", msg.Subject)
	}
	for _, want := range []string{"Olá, Novo!", "Seu cadastro no Lumini Hub da empresa Focco está pronto.", "novo.usuario", "também pode entrar com este e-mail", "7 dias", "Criar minha senha"} {
		if !strings.Contains(msg.HTML, want) {
			t.Errorf("HTML sem %q", want)
		}
	}
	if !strings.Contains(msg.Text, "https://app.exemplo.com/reset-password?token=abc&invite=1") || !strings.Contains(msg.Text, "Seu cadastro no Lumini Hub da empresa Focco está pronto.") {
		t.Errorf("texto inesperado:\n%s", msg.Text)
	}

	// Sem empresa (o nome cai no padrão "Lumini Hub"), a frase não repete o nome.
	msg, err = Render(TemplateUserInvite, "novo@example.com", TemplateData{CompanyName: "Lumini Hub", Name: "Novo", Username: "n", Link: "https://x/y", ExpiresDays: 7})
	if err != nil || !strings.Contains(msg.Text, "Seu cadastro no Lumini Hub está pronto.") {
		t.Errorf("sem InvitedBy: %v\n%s", err, msg.Text)
	}
}

func TestRender_DadosDeAcessoDoUsuario(t *testing.T) {
	dados := TemplateData{
		CompanyName: "Focco Iluminação", Name: "Joel", Username: "joel.silva", Password: "Prov<&>123",
		Link: "https://app.exemplo.com/login", MustChangePassword: true,
	}
	msg, err := Render(TemplateUserWelcome, "joel@example.com", dados)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if msg.Subject != "Seus dados de acesso ao Lumini Hub" {
		t.Errorf("assunto = %q", msg.Subject)
	}
	for _, want := range []string{"Olá, Joel!", "Seu cadastro no Lumini Hub da empresa Focco Iluminação está pronto.", "joel.silva", "Prov&lt;&amp;&gt;123", "Acessar o sistema", "será necessário trocar a senha no primeiro acesso"} {
		if !strings.Contains(msg.HTML, want) {
			t.Errorf("HTML sem %q", want)
		}
	}
	if strings.Contains(msg.HTML, "sugerimos") {
		t.Error("com troca obrigatória não há sugestão")
	}
	for _, want := range []string{"Senha provisória: Prov<&>123", "https://app.exemplo.com/login", "ATENÇÃO: será necessário trocar a senha no primeiro acesso."} {
		if !strings.Contains(msg.Text, want) {
			t.Errorf("texto sem %q", want)
		}
	}

	// Sem troca obrigatória: só a sugestão.
	dados.MustChangePassword = false
	msg, err = Render(TemplateUserWelcome, "joel@example.com", dados)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(msg.HTML, "será necessário") || strings.Contains(msg.Text, "será necessário") ||
		!strings.Contains(msg.HTML, "sugerimos que você troque essa senha provisória") || !strings.Contains(msg.Text, "sugerimos que você troque") {
		t.Errorf("deveria só sugerir a troca:\n%s", msg.Text)
	}
}
