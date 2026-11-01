package validator

import (
	netmail "net/mail"
	"strings"

	"lumini-hub/api.core/internal/domain"
	"lumini-hub/common/mail"
)

// ValidateSMTPSettings valida o que vai ser gravado como configuração SMTP (já com os padrões aplicados:
// porta, segurança e remetente). Desligada, aceita campos incompletos (dá para ir preenchendo); ligada,
// exige servidor, porta e remetente válidos, e senha quando há usuário. hasPassword diz se a senha já
// existe (a nova ou a que estava gravada).
func ValidateSMTPSettings(s domain.SMTPSettings, hasPassword bool) error {
	var errs ValidationErrors

	if s.Security != "" && !mail.ValidSecurity(s.Security) {
		errs.AddError("security", "segurança inválida: use starttls, ssl ou none")
	}
	if s.FromEmail != "" {
		if _, err := netmail.ParseAddress(s.FromEmail); err != nil || strings.ContainsAny(s.FromEmail, "\r\n") {
			errs.AddError("from_email", "e-mail do remetente inválido")
		}
	}
	if strings.ContainsAny(s.FromName, "\r\n") {
		errs.AddError("from_name", "nome do remetente inválido")
	}

	if s.Enabled {
		if s.Host == "" {
			errs.AddError("host", "informe o servidor SMTP para ativar o envio")
		}
		if s.Port == 0 {
			errs.AddError("port", "informe a porta do servidor SMTP para ativar o envio")
		}
		if s.FromEmail == "" {
			errs.AddError("from_email", "informe o e-mail do remetente para ativar o envio")
		}
		if s.Username != "" && !hasPassword {
			errs.AddError("password", "informe a senha do usuário SMTP")
		}
	}

	if errs.HasErrors() {
		return errs
	}
	return nil
}
