// Package mail envia e-mail por SMTP. Não conhece empresa, banco nem template: recebe a configuração
// do servidor já descriptografada e a mensagem pronta (quem resolve qual configuração usar é o api.core).
package mail

import (
	"crypto/tls"
	"errors"
	"fmt"
	"mime"
	"net"
	"net/mail"
	"net/smtp"
	"strconv"
	"strings"
	"time"
)

// Segurança da conexão com o servidor SMTP.
const (
	SecurityStartTLS = "starttls" // conecta em texto e sobe para TLS (porta 587, o padrão)
	SecuritySSL      = "ssl"      // TLS desde o início (porta 465)
	SecurityNone     = "none"     // sem criptografia (só para servidor local de teste, como o Mailpit)
)

// ValidSecurity diz se o valor é um dos modos aceitos.
func ValidSecurity(security string) bool {
	switch security {
	case SecurityStartTLS, SecuritySSL, SecurityNone:
		return true
	}
	return false
}

// ServerConfig é a configuração de um servidor SMTP, com a senha em texto puro (só em memória).
type ServerConfig struct {
	Host      string
	Port      int
	Security  string
	Username  string
	Password  string
	FromName  string
	FromEmail string
}

// Message é um e-mail pronto para enviar. Text é a versão simples (opcional); HTML também é opcional,
// mas pelo menos um dos dois precisa existir.
type Message struct {
	To      string
	Subject string
	Text    string
	HTML    string
}

// Sender envia mensagens. Existe para os testes trocarem o envio real por um falso.
type Sender interface {
	Send(cfg ServerConfig, msg Message) error
}

// SMTPSender é o Sender de verdade, por net/smtp.
type SMTPSender struct {
	Timeout time.Duration
}

// NewSMTPSender cria o Sender padrão, com 20 segundos de limite para conectar e enviar.
func NewSMTPSender() *SMTPSender { return &SMTPSender{Timeout: 20 * time.Second} }

// Send entrega a mensagem. Os erros de rede/protocolo voltam como estão: o chamador decide o que mostrar.
func (s *SMTPSender) Send(cfg ServerConfig, msg Message) error {
	body, from, to, err := buildMessage(cfg, msg)
	if err != nil {
		return err
	}

	timeout := s.Timeout
	if timeout <= 0 {
		timeout = 20 * time.Second
	}
	addr := net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))
	tlsConfig := &tls.Config{ServerName: cfg.Host, MinVersion: tls.VersionTLS12}
	dialer := &net.Dialer{Timeout: timeout}

	var conn net.Conn
	if cfg.Security == SecuritySSL {
		conn, err = tls.DialWithDialer(dialer, "tcp", addr, tlsConfig)
	} else {
		conn, err = dialer.Dial("tcp", addr)
	}
	if err != nil {
		return fmt.Errorf("não foi possível conectar a %s: %w", addr, err)
	}
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(timeout))

	client, err := smtp.NewClient(conn, cfg.Host)
	if err != nil {
		return fmt.Errorf("o servidor %s não respondeu como SMTP: %w", addr, err)
	}
	defer client.Close()

	if cfg.Security == SecurityStartTLS {
		if ok, _ := client.Extension("STARTTLS"); !ok {
			return errors.New("o servidor não oferece STARTTLS: troque a segurança da conexão")
		}
		if err := client.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("falha ao iniciar TLS: %w", err)
		}
	}

	if cfg.Username != "" {
		if ok, _ := client.Extension("AUTH"); !ok {
			return errors.New("o servidor não oferece autenticação, mas há usuário configurado")
		}
		if err := client.Auth(plainAuth{username: cfg.Username, password: cfg.Password}); err != nil {
			return fmt.Errorf("autenticação recusada pelo servidor: %w", err)
		}
	}

	if err := client.Mail(from); err != nil {
		return fmt.Errorf("remetente recusado: %w", err)
	}
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("destinatário recusado: %w", err)
	}
	w, err := client.Data()
	if err != nil {
		return err
	}
	if _, err := w.Write(body); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	return client.Quit()
}

// plainAuth é o AUTH PLAIN sem a trava do net/smtp.PlainAuth (que recusa conexão sem TLS): aqui a
// decisão de usar ou não criptografia é da configuração, e "none" existe para servidores de teste.
type plainAuth struct{ username, password string }

func (a plainAuth) Start(*smtp.ServerInfo) (string, []byte, error) {
	return "PLAIN", []byte("\x00" + a.username + "\x00" + a.password), nil
}

func (a plainAuth) Next(_ []byte, more bool) ([]byte, error) {
	if more {
		return nil, errors.New("resposta inesperada do servidor durante a autenticação")
	}
	return nil, nil
}

// buildMessage monta o texto RFC 5322 da mensagem e devolve também o remetente e o destinatário puros
// (só o endereço) para o envelope SMTP. Recusa quebra de linha em qualquer cabeçalho (injeção de cabeçalho).
func buildMessage(cfg ServerConfig, msg Message) (body []byte, fromAddr, toAddr string, err error) {
	if strings.TrimSpace(msg.Text) == "" && strings.TrimSpace(msg.HTML) == "" {
		return nil, "", "", errors.New("a mensagem está vazia")
	}
	for _, header := range []string{cfg.FromName, cfg.FromEmail, msg.To, msg.Subject} {
		if strings.ContainsAny(header, "\r\n") {
			return nil, "", "", errors.New("cabeçalho do e-mail com quebra de linha")
		}
	}

	from, err := mail.ParseAddress(cfg.FromEmail)
	if err != nil {
		return nil, "", "", fmt.Errorf("e-mail do remetente inválido: %w", err)
	}
	to, err := mail.ParseAddress(msg.To)
	if err != nil {
		return nil, "", "", fmt.Errorf("e-mail do destinatário inválido: %w", err)
	}
	from.Name = cfg.FromName

	var b strings.Builder
	write := func(format string, args ...any) { fmt.Fprintf(&b, format, args...) }
	write("From: %s\r\n", from.String())
	write("To: %s\r\n", to.Address)
	write("Subject: %s\r\n", mime.QEncoding.Encode("utf-8", msg.Subject))
	write("Date: %s\r\n", time.Now().Format(time.RFC1123Z))
	write("MIME-Version: 1.0\r\n")

	const boundary = "lumini-hub-alt"
	switch {
	case msg.Text != "" && msg.HTML != "":
		write("Content-Type: multipart/alternative; boundary=%q\r\n\r\n", boundary)
		write("--%s\r\nContent-Type: text/plain; charset=utf-8\r\nContent-Transfer-Encoding: 8bit\r\n\r\n%s\r\n", boundary, crlf(msg.Text))
		write("--%s\r\nContent-Type: text/html; charset=utf-8\r\nContent-Transfer-Encoding: 8bit\r\n\r\n%s\r\n", boundary, crlf(msg.HTML))
		write("--%s--\r\n", boundary)
	case msg.HTML != "":
		write("Content-Type: text/html; charset=utf-8\r\nContent-Transfer-Encoding: 8bit\r\n\r\n%s\r\n", crlf(msg.HTML))
	default:
		write("Content-Type: text/plain; charset=utf-8\r\nContent-Transfer-Encoding: 8bit\r\n\r\n%s\r\n", crlf(msg.Text))
	}
	return []byte(b.String()), from.Address, to.Address, nil
}

// crlf normaliza as quebras de linha para CRLF, como o SMTP exige.
func crlf(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, "\r\n", "\n"), "\n", "\r\n")
}
