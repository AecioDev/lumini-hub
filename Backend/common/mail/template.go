package mail

import (
	"bytes"
	"embed"
	"fmt"
	htmltemplate "html/template"
	texttemplate "text/template"
)

//go:embed templates/*
var templateFiles embed.FS

// Nomes dos templates de e-mail do sistema. Cada um tem um .html (dentro do layout.html) e um .txt em templates/.
const (
	TemplatePasswordReset = "password_reset"
	TemplateUserInvite    = "user_invite"
	TemplateUserWelcome   = "user_welcome"
)

// templateSubjects guarda o assunto de cada template.
var templateSubjects = map[string]string{
	TemplatePasswordReset: "Redefinição de senha",
	TemplateUserInvite:    "Convite para acessar o Lumini Hub",
	TemplateUserWelcome:   "Seus dados de acesso ao Lumini Hub",
}

// TemplateData é o que os templates enxergam. Campos que um template não usa ficam vazios.
type TemplateData struct {
	Subject     string // preenchido pelo Render
	CompanyName string // nome no cabeçalho do e-mail
	Name        string // nome de quem recebe

	// Link para definir a senha (redefinição de senha e convite)
	Link string

	// Redefinição de senha
	ExpiresMinutes int

	// Convite de novo usuário
	Username    string // nome de usuário de acesso
	InvitedBy   string // quem cadastrou (opcional)
	ExpiresDays int

	// Dados de acesso de um usuário cadastrado com senha inicial
	Password           string // senha provisória, em texto
	MustChangePassword bool   // true: a troca no primeiro acesso é obrigatória; false: só sugerida
}

// Render monta a mensagem (assunto, texto e HTML) de um template para o destinatário.
func Render(name, to string, data TemplateData) (Message, error) {
	subject, ok := templateSubjects[name]
	if !ok {
		return Message{}, fmt.Errorf("template de e-mail desconhecido: %q", name)
	}
	data.Subject = subject

	htmlTpl, err := htmltemplate.ParseFS(templateFiles, "templates/layout.html", "templates/"+name+".html")
	if err != nil {
		return Message{}, fmt.Errorf("lendo o template HTML %q: %w", name, err)
	}
	var htmlBody bytes.Buffer
	if err := htmlTpl.ExecuteTemplate(&htmlBody, "layout", data); err != nil {
		return Message{}, fmt.Errorf("montando o HTML %q: %w", name, err)
	}

	textTpl, err := texttemplate.ParseFS(templateFiles, "templates/"+name+".txt")
	if err != nil {
		return Message{}, fmt.Errorf("lendo o template de texto %q: %w", name, err)
	}
	var textBody bytes.Buffer
	if err := textTpl.ExecuteTemplate(&textBody, name+".txt", data); err != nil {
		return Message{}, fmt.Errorf("montando o texto %q: %w", name, err)
	}

	return Message{To: to, Subject: subject, Text: textBody.String(), HTML: htmlBody.String()}, nil
}
