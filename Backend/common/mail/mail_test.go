package mail

import (
	"bufio"
	"net"
	"strconv"
	"strings"
	"testing"
)

func TestBuildMessage_MontaCabecalhosEEnvelope(t *testing.T) {
	cfg := ServerConfig{FromName: "Lumini Hub", FromEmail: "noreply@example.com"}
	body, from, to, err := buildMessage(cfg, Message{To: "ana@example.com", Subject: "Recuperação de senha", Text: "olá", HTML: "<p>olá</p>"})
	if err != nil {
		t.Fatalf("buildMessage: %v", err)
	}
	if from != "noreply@example.com" || to != "ana@example.com" {
		t.Errorf("envelope = %q -> %q", from, to)
	}
	text := string(body)
	for _, want := range []string{"From: \"Lumini Hub\" <noreply@example.com>", "To: ana@example.com", "multipart/alternative", "=?utf-8?q?Recupera", "text/plain", "text/html"} {
		if !strings.Contains(text, want) {
			t.Errorf("mensagem sem %q:\n%s", want, text)
		}
	}
}

func TestBuildMessage_RecusaInjecaoDeCabecalhoEMensagemVazia(t *testing.T) {
	cfg := ServerConfig{FromEmail: "noreply@example.com"}
	cases := map[string]Message{
		"assunto com quebra":      {To: "ana@example.com", Subject: "oi\r\nBcc: x@y.com", Text: "a"},
		"destinatário com quebra": {To: "ana@example.com\r\nBcc: x@y.com", Subject: "oi", Text: "a"},
		"mensagem vazia":          {To: "ana@example.com", Subject: "oi"},
		"destinatário inválido":   {To: "isto-nao-e-email", Subject: "oi", Text: "a"},
	}
	for name, msg := range cases {
		if _, _, _, err := buildMessage(cfg, msg); err == nil {
			t.Errorf("%s: deveria ser recusado", name)
		}
	}
	if _, _, _, err := buildMessage(ServerConfig{FromEmail: "ruim"}, Message{To: "ana@example.com", Text: "a"}); err == nil {
		t.Error("remetente inválido deveria ser recusado")
	}
}

func TestValidSecurity(t *testing.T) {
	for _, ok := range []string{"starttls", "ssl", "none"} {
		if !ValidSecurity(ok) {
			t.Errorf("%q deveria ser válido", ok)
		}
	}
	if ValidSecurity("") || ValidSecurity("tls") {
		t.Error("valores fora da lista deveriam ser inválidos")
	}
}

// fakeSMTP sobe um servidor SMTP mínimo (sem TLS) e devolve o endereço e o que recebeu no DATA.
func fakeSMTP(t *testing.T, requireAuth bool) (host string, port int, received chan string) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { ln.Close() })
	received = make(chan string, 1)

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		r := bufio.NewReader(conn)
		say := func(s string) { conn.Write([]byte(s + "\r\n")) }
		say("220 fake ESMTP")
		var data strings.Builder
		inData := false
		for {
			line, err := r.ReadString('\n')
			if err != nil {
				return
			}
			line = strings.TrimRight(line, "\r\n")
			if inData {
				if line == "." {
					inData = false
					received <- data.String()
					say("250 ok")
					continue
				}
				data.WriteString(line + "\n")
				continue
			}
			switch cmd := strings.ToUpper(line); {
			case strings.HasPrefix(cmd, "EHLO"):
				if requireAuth {
					say("250-fake\r\n250 AUTH PLAIN")
				} else {
					say("250 fake")
				}
			case strings.HasPrefix(cmd, "AUTH"):
				say("235 ok")
			case strings.HasPrefix(cmd, "MAIL"), strings.HasPrefix(cmd, "RCPT"):
				say("250 ok")
			case cmd == "DATA":
				inData = true
				say("354 go")
			case cmd == "QUIT":
				say("221 bye")
				return
			default:
				say("250 ok")
			}
		}
	}()

	h, p, _ := net.SplitHostPort(ln.Addr().String())
	port, _ = strconv.Atoi(p)
	return h, port, received
}

func TestSMTPSender_EntregaSemCriptografia(t *testing.T) {
	host, port, received := fakeSMTP(t, true)
	cfg := ServerConfig{Host: host, Port: port, Security: SecurityNone, Username: "u", Password: "p", FromEmail: "noreply@example.com"}

	if err := NewSMTPSender().Send(cfg, Message{To: "ana@example.com", Subject: "teste", Text: "corpo"}); err != nil {
		t.Fatalf("Send: %v", err)
	}
	if got := <-received; !strings.Contains(got, "corpo") || !strings.Contains(got, "To: ana@example.com") {
		t.Errorf("servidor recebeu: %q", got)
	}
}

func TestSMTPSender_ErroClaroQuandoNaoConecta(t *testing.T) {
	ln, _ := net.Listen("tcp", "127.0.0.1:0")
	_, p, _ := net.SplitHostPort(ln.Addr().String())
	ln.Close() // porta fechada
	port, _ := strconv.Atoi(p)

	err := NewSMTPSender().Send(ServerConfig{Host: "127.0.0.1", Port: port, Security: SecurityNone, FromEmail: "a@b.com"}, Message{To: "c@d.com", Text: "x"})
	if err == nil || !strings.Contains(err.Error(), "não foi possível conectar") {
		t.Errorf("erro = %v", err)
	}
}

func TestSMTPSender_StartTLSExigeSuporteDoServidor(t *testing.T) {
	host, port, _ := fakeSMTP(t, false)
	err := NewSMTPSender().Send(ServerConfig{Host: host, Port: port, Security: SecurityStartTLS, FromEmail: "a@b.com"}, Message{To: "c@d.com", Text: "x"})
	if err == nil || !strings.Contains(err.Error(), "STARTTLS") {
		t.Errorf("erro = %v", err)
	}
}
