package service

import (
	"errors"
	"fmt"
	"strings"

	"lumini-hub/api.core/internal/domain"
	"lumini-hub/api.core/internal/repository"
	"lumini-hub/api.core/internal/validator"
	"lumini-hub/common/mail"
	"lumini-hub/common/utils"
)

// Origem da configuração SMTP efetivamente usada num envio.
const (
	MailSourceCompany = "company" // a própria empresa
	MailSourceParent  = "parent"  // uma empresa acima, na cadeia até a Matriz
	MailSourceSystem  = "system"  // o padrão do sistema
)

// maxCompanyChain limita a subida pelos ParentID (o validator de empresa já impede ciclos; é só uma trava).
const maxCompanyChain = 20

// ErrMailNotConfigured: nenhuma configuração SMTP ativa foi encontrada na cadeia da empresa nem no sistema.
var ErrMailNotConfigured = errors.New("e-mail não configurado")

// MailNotConfiguredError carrega a mensagem pronta para o usuário (qual empresa, qual Matriz).
type MailNotConfiguredError struct{ Message string }

func (e *MailNotConfiguredError) Error() string { return e.Message }

// Is faz errors.Is(err, ErrMailNotConfigured) valer para este erro.
func (e *MailNotConfiguredError) Is(target error) bool { return target == ErrMailNotConfigured }

// MailSendError é a falha do servidor SMTP ao entregar (conexão, autenticação, destinatário recusado...).
type MailSendError struct{ Err error }

func (e *MailSendError) Error() string { return e.Err.Error() }
func (e *MailSendError) Unwrap() error { return e.Err }

// ResolvedMail é a configuração SMTP escolhida para um envio, com a senha já descriptografada.
type ResolvedMail struct {
	Server mail.ServerConfig
	Source string // MailSource*
	Name   string // quem é a dona: nome da empresa, ou "Sistema"
}

// MailConfigService guarda e resolve a configuração de e-mail: a global do sistema e a de cada empresa.
//
// Resolução (decidida com o usuário em 2026-11-01): a empresa usa a própria configuração ativa; sem ela,
// sobe pelos ParentID até a Matriz do grupo; sem nenhuma, usa a do sistema. Nunca olha uma empresa fora
// da cadeia. Sem nenhuma configuração ativa, para com MailNotConfiguredError.
type MailConfigService struct {
	uow           repository.UnitOfWork
	sender        mail.Sender
	encryptionKey string
}

// NewMailConfigService cria o serviço. encryptionKey é a chave AES-256 em hex usada na senha do SMTP.
func NewMailConfigService(uow repository.UnitOfWork, sender mail.Sender, encryptionKey string) *MailConfigService {
	return &MailConfigService{uow: uow, sender: sender, encryptionKey: encryptionKey}
}

// GetSystemEmail devolve o e-mail padrão do sistema
func (s *MailConfigService) GetSystemEmail() (*domain.ApiSystemEmailSettings, error) {
	settings, err := s.uow.SystemSettings().Get()
	if err != nil {
		return nil, err
	}
	return &domain.ApiSystemEmailSettings{
		ApiSMTPSettings: domain.ApiSMTPSettingsFromModel(settings.SMTP),
		UpdatedAt:       settings.UpdatedAt,
	}, nil
}

// SaveSystemEmail grava o e-mail padrão do sistema
func (s *MailConfigService) SaveSystemEmail(req domain.SaveSMTPRequest, userID uint) (*domain.ApiSystemEmailSettings, error) {
	settings, err := s.uow.SystemSettings().Get()
	if err != nil {
		return nil, err
	}
	if err := s.applyRequest(&settings.SMTP, req); err != nil {
		return nil, err
	}
	settings.UpdatedByID = &userID

	if err := s.uow.Execute(func(uow repository.UnitOfWork) error {
		return uow.SystemSettings().Save(settings)
	}); err != nil {
		return nil, err
	}
	return s.GetSystemEmail()
}

// GetCompanyEmail devolve a configuração de e-mail de uma empresa (padrões se ela herda) e de onde o envio sairia hoje
func (s *MailConfigService) GetCompanyEmail(companyID uint) (*domain.ApiCompanyEmailConfig, error) {
	if err := s.requireCompany(companyID); err != nil {
		return nil, err
	}
	dto := domain.ApiCompanyEmailConfig{
		ApiSMTPSettings: domain.ApiSMTPSettingsFromModel(domain.SMTPSettings{Port: 587, Security: mail.SecurityStartTLS}),
		CompanyID:       companyID,
	}

	config, err := s.uow.CompanyEmailConfigs().FindByCompanyID(companyID)
	if err != nil {
		return nil, err
	}
	if config != nil {
		dto.ApiSMTPSettings = domain.ApiSMTPSettingsFromModel(config.SMTP)
		dto.Configured = true
		dto.UpdatedAt = config.UpdatedAt
	}

	if resolved, err := s.Resolve(&companyID); err == nil {
		dto.EffectiveSource = resolved.Source
		dto.EffectiveName = resolved.Name
	} else if !errors.Is(err, ErrMailNotConfigured) {
		return nil, err
	}
	return &dto, nil
}

// SaveCompanyEmail cria ou atualiza a configuração de e-mail própria de uma empresa
func (s *MailConfigService) SaveCompanyEmail(companyID uint, req domain.SaveSMTPRequest) (*domain.ApiCompanyEmailConfig, error) {
	if err := s.requireCompany(companyID); err != nil {
		return nil, err
	}
	config, err := s.uow.CompanyEmailConfigs().FindByCompanyID(companyID)
	if err != nil {
		return nil, err
	}
	isNew := config == nil
	if isNew {
		config = &domain.CompanyEmailConfig{CompanyID: companyID}
	}
	if err := s.applyRequest(&config.SMTP, req); err != nil {
		return nil, err
	}

	if err := s.uow.Execute(func(uow repository.UnitOfWork) error {
		if isNew {
			return uow.CompanyEmailConfigs().Create(config)
		}
		return uow.CompanyEmailConfigs().Update(config)
	}); err != nil {
		return nil, err
	}
	return s.GetCompanyEmail(companyID)
}

// DeleteCompanyEmail remove a configuração própria: a empresa volta a herdar da Matriz ou do sistema
func (s *MailConfigService) DeleteCompanyEmail(companyID uint) error {
	if err := s.requireCompany(companyID); err != nil {
		return err
	}
	config, err := s.uow.CompanyEmailConfigs().FindByCompanyID(companyID)
	if err != nil {
		return err
	}
	if config == nil {
		return utils.ErrNotFound
	}
	return s.uow.Execute(func(uow repository.UnitOfWork) error {
		return uow.CompanyEmailConfigs().Delete(config.ID)
	})
}

// Resolve escolhe a configuração SMTP de um envio. companyID nil (sem empresa) vai direto para a do sistema.
func (s *MailConfigService) Resolve(companyID *uint) (*ResolvedMail, error) {
	var chain []string // nomes das empresas consultadas, para a mensagem de erro

	current := companyID
	for depth := 0; current != nil && depth < maxCompanyChain; depth++ {
		company, err := s.uow.Companies().FindByID(*current)
		if err != nil {
			return nil, err
		}
		if company == nil {
			return nil, utils.ErrNotFound
		}
		name := companyDisplayName(company.TradeName, company.LegalName)
		chain = append(chain, name)

		config, err := s.uow.CompanyEmailConfigs().FindByCompanyID(company.ID)
		if err != nil {
			return nil, err
		}
		if config != nil && usable(config.SMTP) {
			source := MailSourceCompany
			if depth > 0 {
				source = MailSourceParent
			}
			return s.resolved(config.SMTP, source, name)
		}
		current = company.ParentID
	}

	settings, err := s.uow.SystemSettings().Get()
	if err != nil {
		return nil, err
	}
	if usable(settings.SMTP) {
		return s.resolved(settings.SMTP, MailSourceSystem, "Sistema")
	}

	if len(chain) == 0 {
		return nil, &MailNotConfiguredError{Message: "O e-mail do sistema não está configurado."}
	}
	return nil, &MailNotConfiguredError{Message: fmt.Sprintf(
		"Nenhum e-mail está configurado para %s, e o e-mail padrão do sistema também não está ativo.",
		strings.Join(chain, " nem para "))}
}

// Send entrega uma mensagem usando a configuração resolvida para a empresa. Devolve de onde saiu.
func (s *MailConfigService) Send(companyID *uint, msg mail.Message) (*ResolvedMail, error) {
	resolved, err := s.Resolve(companyID)
	if err != nil {
		return nil, err
	}
	if err := s.sender.Send(resolved.Server, msg); err != nil {
		return nil, &MailSendError{Err: err}
	}
	return resolved, nil
}

// SendTest manda um e-mail de teste pela configuração que valeria para a empresa (ou a do sistema, sem empresa).
func (s *MailConfigService) SendTest(companyID *uint, to string) (*domain.ApiTestEmailResult, error) {
	resolved, err := s.Send(companyID, mail.Message{
		To:      to,
		Subject: "Teste de e-mail - Lumini Hub",
		Text:    "Este é um e-mail de teste do Lumini Hub.\n\nSe você o recebeu, o envio está funcionando.",
	})
	if err != nil {
		return nil, err
	}
	return &domain.ApiTestEmailResult{Source: resolved.Source, Name: resolved.Name}, nil
}

// applyRequest aplica os campos da request sobre a configuração, criptografando a senha nova (vazia mantém
// a atual; sem usuário não há senha) e validando o resultado antes de qualquer gravação.
func (s *MailConfigService) applyRequest(target *domain.SMTPSettings, req domain.SaveSMTPRequest) error {
	updated := *target
	updated.Host = strings.TrimSpace(req.Host)
	updated.Port = req.Port
	if updated.Port == 0 {
		updated.Port = 587
	}
	updated.Security = strings.TrimSpace(req.Security)
	if updated.Security == "" {
		updated.Security = mail.SecurityStartTLS
	}
	updated.Username = strings.TrimSpace(req.Username)
	updated.FromName = strings.TrimSpace(req.FromName)
	updated.FromEmail = strings.TrimSpace(req.FromEmail)
	updated.Enabled = req.Enabled

	if updated.Username == "" {
		updated.Password = ""
	} else if req.Password != "" {
		encrypted, err := utils.EncryptAES(req.Password, s.encryptionKey)
		if err != nil {
			return err
		}
		updated.Password = encrypted
	}

	if err := validator.ValidateSMTPSettings(updated, updated.Password != ""); err != nil {
		return err
	}
	*target = updated
	return nil
}

func (s *MailConfigService) requireCompany(companyID uint) error {
	company, err := s.uow.Companies().FindByID(companyID)
	if err != nil {
		return err
	}
	if company == nil {
		return utils.ErrNotFound
	}
	return nil
}

func (s *MailConfigService) resolved(smtp domain.SMTPSettings, source, name string) (*ResolvedMail, error) {
	password := ""
	if smtp.Password != "" {
		var err error
		if password, err = utils.DecryptAES(smtp.Password, s.encryptionKey); err != nil {
			return nil, fmt.Errorf("não foi possível ler a senha do e-mail de %s: %w", name, err)
		}
	}
	return &ResolvedMail{
		Server: mail.ServerConfig{
			Host: smtp.Host, Port: smtp.Port, Security: smtp.Security,
			Username: smtp.Username, Password: password,
			FromName: smtp.FromName, FromEmail: smtp.FromEmail,
		},
		Source: source,
		Name:   name,
	}, nil
}

// usable diz se a configuração está ligada e completa o bastante para enviar.
func usable(s domain.SMTPSettings) bool {
	return s.Enabled && s.Host != "" && s.Port > 0 && s.FromEmail != ""
}

func companyDisplayName(tradeName, legalName string) string {
	if tradeName != "" {
		return tradeName
	}
	return legalName
}
