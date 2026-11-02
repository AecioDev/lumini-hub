package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"log"
	"strings"
	"time"

	"lumini-hub/api.auth/internal/domain"
	"lumini-hub/common/internalmail"
	"lumini-hub/common/utils"

	"gorm.io/gorm"
)

const (
	// PasswordResetTTL é por quanto tempo o link de redefinição vale.
	PasswordResetTTL = 30 * time.Minute
	// InviteTTL é por quanto tempo o link do convite de um novo usuário vale (mais que o da recuperação).
	InviteTTL = 7 * 24 * time.Hour
	// maxResetRequestsPerHour limita os pedidos por usuário numa hora (os pedidos além disso são ignorados em silêncio).
	maxResetRequestsPerHour = 3
)

// ErrInvalidResetToken: o link não existe, já foi usado ou expirou. A mensagem é a mesma nos três casos.
var ErrInvalidResetToken = errors.New("link inválido ou expirado: peça uma nova redefinição de senha")

// PasswordResetService gera os links de redefinição de senha e troca a senha com eles.
type PasswordResetService struct {
	db     *gorm.DB
	mailer internalmail.Sender
	// dispatch executa o envio do e-mail fora da requisição, para que o tempo de resposta não revele se o
	// e-mail existe. Os testes trocam por uma chamada direta.
	dispatch func(func())
}

// NewPasswordResetService cria o serviço. mailer pede o envio do e-mail ao api.core.
func NewPasswordResetService(db *gorm.DB, mailer internalmail.Sender) *PasswordResetService {
	return &PasswordResetService{db: db, mailer: mailer, dispatch: func(f func()) { go f() }}
}

// ForgotPassword cria um token e dispara o e-mail para a conta ativa com esse e-mail. Não devolve erro quando
// o e-mail não existe (nem quando o limite de pedidos foi atingido ou o envio falha): quem chama responde
// sempre a mesma coisa, para não revelar quais e-mails estão cadastrados. Só erro de banco volta.
func (s *PasswordResetService) ForgotPassword(email, requestedIP string) error {
	var user domain.User
	err := s.db.Where("LOWER(email) = LOWER(?) AND is_active = ?", strings.TrimSpace(email), true).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	if err != nil {
		return err
	}

	var recent int64
	if err := s.db.Model(&domain.PasswordResetToken{}).
		Where("user_id = ? AND purpose = ? AND created_at > ?", user.ID, domain.TokenPurposeReset, time.Now().Add(-time.Hour)).
		Count(&recent).Error; err != nil {
		return err
	}
	if recent >= maxResetRequestsPerHour {
		log.Printf("Recuperação de senha: limite de pedidos atingido para o usuário %d", user.ID)
		return nil
	}

	token, err := s.createToken(user.ID, domain.TokenPurposeReset, PasswordResetTTL, requestedIP)
	if err != nil {
		return err
	}

	request := internalmail.PasswordResetRequest{
		To:             user.Email,
		Name:           user.Name,
		Token:          token,
		CompanyID:      user.CompanyID,
		ExpiresMinutes: int(PasswordResetTTL.Minutes()),
	}
	userID := user.ID
	s.dispatch(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 35*time.Second)
		defer cancel()
		if err := s.mailer.SendPasswordReset(ctx, request); err != nil {
			log.Printf("Recuperação de senha: e-mail não enviado ao usuário %d: %v", userID, err)
		}
	})
	return nil
}

// createToken gera um token novo para o usuário e invalida os anteriores ainda abertos (mas os mantém, para a
// contagem do limite de pedidos). Devolve o valor puro, que só existe para ir no link do e-mail.
func (s *PasswordResetService) createToken(userID uint, purpose string, ttl time.Duration, requestedIP string) (string, error) {
	token, hash, err := newResetToken()
	if err != nil {
		return "", err
	}
	now := time.Now()
	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&domain.PasswordResetToken{}).
			Where("user_id = ? AND used_at IS NULL", userID).
			Update("used_at", now).Error; err != nil {
			return err
		}
		return tx.Create(&domain.PasswordResetToken{
			UserID:      userID,
			Purpose:     purpose,
			TokenHash:   hash,
			ExpiresAt:   now.Add(ttl),
			RequestedIP: requestedIP,
		}).Error
	})
	if err != nil {
		return "", err
	}
	return token, nil
}

// IssueInvite cria o link de convite do usuário (válido por InviteTTL, invalidando os links anteriores dele) e
// pede o e-mail ao api.core na hora, sem enviar em segundo plano: quem cadastra precisa saber se o e-mail saiu.
// Quando o link foi criado mas o e-mail falhou, devolve mailErr e err nil (o convite pode ser reenviado).
func (s *PasswordResetService) IssueInvite(user domain.User, invitedBy string) (mailErr error, err error) {
	token, err := s.createToken(user.ID, domain.TokenPurposeInvite, InviteTTL, "")
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 35*time.Second)
	defer cancel()
	mailErr = s.mailer.SendUserInvite(ctx, internalmail.UserInviteRequest{
		To:          user.Email,
		Name:        user.Name,
		Username:    user.Username,
		Token:       token,
		CompanyID:   user.CompanyID,
		ExpiresDays: int(InviteTTL.Hours() / 24),
		InvitedBy:   invitedBy,
	})
	if mailErr != nil {
		log.Printf("Convite de usuário: e-mail não enviado ao usuário %d: %v", user.ID, mailErr)
	}
	return mailErr, nil
}

// SendWelcome pede ao api.core o e-mail com os dados de acesso (usuário e senha provisória) de um usuário cadastrado
// com senha inicial. Na hora, como o convite: quem cadastra precisa saber se o e-mail saiu. A senha só passa por aqui,
// em memória, a caminho do e-mail: nunca é gravada nem vai para o log.
func (s *PasswordResetService) SendWelcome(user domain.User, password string, mustChangePassword bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 35*time.Second)
	defer cancel()
	err := s.mailer.SendUserWelcome(ctx, internalmail.UserWelcomeRequest{
		To:                 user.Email,
		Name:               user.Name,
		Username:           user.Username,
		Password:           password,
		CompanyID:          user.CompanyID,
		MustChangePassword: mustChangePassword,
	})
	if err != nil {
		log.Printf("Cadastro de usuário: e-mail de acesso não enviado ao usuário %d: %v", user.ID, err)
	}
	return err
}

// ResetPassword troca a senha da conta dona do token e o inutiliza (junto com qualquer outro pedido aberto dela).
func (s *PasswordResetService) ResetPassword(token, newPassword string) error {
	hash := hashResetToken(token)
	now := time.Now()

	var record domain.PasswordResetToken
	err := s.db.Where("token_hash = ? AND used_at IS NULL AND expires_at > ?", hash, now).First(&record).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrInvalidResetToken
	}
	if err != nil {
		return err
	}

	passwordHash, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		// O UPDATE condicional garante uso único mesmo com duas requisições simultâneas com o mesmo token.
		claimed := tx.Model(&domain.PasswordResetToken{}).
			Where("id = ? AND used_at IS NULL", record.ID).
			Update("used_at", now)
		if claimed.Error != nil {
			return claimed.Error
		}
		if claimed.RowsAffected != 1 {
			return ErrInvalidResetToken
		}

		// Quem define a própria senha pelo link não precisa mais trocá-la no primeiro acesso.
		updated := tx.Model(&domain.User{}).
			Where("id = ? AND is_active = ?", record.UserID, true).
			Updates(map[string]any{"password_hash": passwordHash, "must_change_password": false})
		if updated.Error != nil {
			return updated.Error
		}
		if updated.RowsAffected != 1 {
			return ErrInvalidResetToken
		}

		return tx.Model(&domain.PasswordResetToken{}).
			Where("user_id = ? AND used_at IS NULL", record.UserID).
			Update("used_at", now).Error
	})
}

// newResetToken gera 32 bytes aleatórios (o token que vai no link) e o hash que fica no banco.
func newResetToken() (token, hash string, err error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", "", err
	}
	token = base64.RawURLEncoding.EncodeToString(raw)
	return token, hashResetToken(token), nil
}

func hashResetToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
