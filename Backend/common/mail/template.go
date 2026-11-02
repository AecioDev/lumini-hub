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
const TemplatePasswordReset = "password_reset"

// templateSubjects guarda o assunto de cada template.
var templateSubjects = map[string]string{
	TemplatePasswordReset: "Redefinição de senha",
}

// TemplateData é o que os templates enxergam. Campos que um template não usa ficam vazios.
type TemplateData struct {
	Subject     string // preenchido pelo Render
	CompanyName string // nome no cabeçalho do e-mail
	Name        string // nome de quem recebe

	// Redefinição de senha
	Link           string
	ExpiresMinutes int
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
