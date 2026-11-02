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
	for _, want := range []string{"Admin criou", "novo.usuario", "7 dias", "Criar minha senha"} {
		if !strings.Contains(msg.HTML, want) {
			t.Errorf("HTML sem %q", want)
		}
	}
	if !strings.Contains(msg.Text, "https://app.exemplo.com/reset-password?token=abc&invite=1") || !strings.Contains(msg.Text, "Admin criou") {
		t.Errorf("texto inesperado:\n%s", msg.Text)
	}

	// Sem quem convidou, o texto continua natural.
	msg, err = Render(TemplateUserInvite, "novo@example.com", TemplateData{CompanyName: "Focco", Name: "Novo", Username: "n", Link: "https://x/y", ExpiresDays: 7})
	if err != nil || !strings.Contains(msg.Text, "Foi criado um acesso") {
		t.Errorf("sem InvitedBy: %v\n%s", err, msg.Text)
	}
}
