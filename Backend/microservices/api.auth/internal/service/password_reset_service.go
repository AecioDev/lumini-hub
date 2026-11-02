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
		Where("user_id = ? AND created_at > ?", user.ID, time.Now().Add(-time.Hour)).
		Count(&recent).Error; err != nil {
		return err
	}
	if recent >= maxResetRequestsPerHour {
		log.Printf("Recuperação de senha: limite de pedidos atingido para o usuário %d", user.ID)
		return nil
	}

	token, hash, err := newResetToken()
	if err != nil {
		return err
	}
	now := time.Now()
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// Um pedido novo invalida os anteriores ainda abertos (mas os mantém, para a contagem do limite).
		if err := tx.Model(&domain.PasswordResetToken{}).
			Where("user_id = ? AND used_at IS NULL", user.ID).
			Update("used_at", now).Error; err != nil {
			return err
		}
		return tx.Create(&domain.PasswordResetToken{
			UserID:      user.ID,
			TokenHash:   hash,
			ExpiresAt:   now.Add(PasswordResetTTL),
			RequestedIP: requestedIP,
		}).Error
	})
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

		updated := tx.Model(&domain.User{}).
			Where("id = ? AND is_active = ?", record.UserID, true).
			Update("password_hash", passwordHash)
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
