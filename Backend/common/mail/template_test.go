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
