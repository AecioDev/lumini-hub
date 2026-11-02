package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"lumini-hub/api.auth/internal/domain"
	"lumini-hub/common/internalmail"
	"lumini-hub/common/testutil"
	"lumini-hub/common/utils"

	"gorm.io/gorm"
)

// Recuperação de senha por e-mail (precisa de Postgres: sem TEST_DATABASE_URL o teste é ignorado).

// mailerFalso guarda os pedidos de envio no lugar de chamar o api.core.
type mailerFalso struct {
	pedidos []internalmail.PasswordResetRequest
	err     error
}

func (m *mailerFalso) SendPasswordReset(_ context.Context, req internalmail.PasswordResetRequest) error {
	if m.err != nil {
		return m.err
	}
	m.pedidos = append(m.pedidos, req)
	return nil
}

func novoServicoDeSenha(t *testing.T) (*PasswordResetService, *mailerFalso, *gorm.DB, domain.User) {
	t.Helper()
	db := testutil.OpenTestDB(t)
	if err := db.AutoMigrate(&domain.Permission{}, &domain.Role{}, &domain.User{}, &domain.PasswordResetToken{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	role := domain.Role{Name: "ADMIN"}
	if err := db.Create(&role).Error; err != nil {
		t.Fatal(err)
	}
	oldHash, _ := utils.HashPassword("senha-antiga")
	companyID := uint(7)
	user := domain.User{Username: "ana", Name: "Ana", Email: "Ana@Exemplo.com", PasswordHash: oldHash, RoleID: role.ID, IsActive: true, CompanyID: &companyID}
	if err := db.Create(&user).Error; err != nil {
		t.Fatal(err)
	}

	mailer := &mailerFalso{}
	svc := NewPasswordResetService(db, mailer)
	svc.dispatch = func(f func()) { f() } // envio síncrono nos testes
	return svc, mailer, db, user
}

func TestSenha_PedidoGeraTokenSoComHashEEnviaEmail(t *testing.T) {
	svc, mailer, db, user := novoServicoDeSenha(t)

	if err := svc.ForgotPassword("  ana@exemplo.COM ", "10.0.0.1"); err != nil {
		t.Fatalf("forgot: %v", err)
	}
	if len(mailer.pedidos) != 1 {
		t.Fatalf("deveria ter pedido 1 e-mail, pediu %d", len(mailer.pedidos))
	}
	p := mailer.pedidos[0]
	if p.To != "Ana@Exemplo.com" || p.Name != "Ana" || p.ExpiresMinutes != 30 || p.CompanyID == nil || *p.CompanyID != 7 || len(p.Token) < 40 {
		t.Errorf("pedido inesperado: %+v", p)
	}

	var saved domain.PasswordResetToken
	if err := db.Where("user_id = ?", user.ID).First(&saved).Error; err != nil {
		t.Fatal(err)
	}
	if saved.TokenHash == p.Token || saved.TokenHash != hashResetToken(p.Token) {
		t.Error("o banco deve guardar só o hash do token")
	}
	if saved.UsedAt != nil || saved.RequestedIP != "10.0.0.1" || time.Until(saved.ExpiresAt) > 31*time.Minute || time.Until(saved.ExpiresAt) < 29*time.Minute {
		t.Errorf("registro inesperado: %+v", saved)
	}
}

func TestSenha_EmailDesconhecidoOuInativoNaoGeraNadaENaoDaErro(t *testing.T) {
	svc, mailer, db, user := novoServicoDeSenha(t)

	if err := svc.ForgotPassword("ninguem@exemplo.com", "1.1.1.1"); err != nil {
		t.Fatalf("e-mail desconhecido não pode dar erro: %v", err)
	}
	db.Model(&domain.User{}).Where("id = ?", user.ID).Update("is_active", false)
	if err := svc.ForgotPassword(user.Email, "1.1.1.1"); err != nil {
		t.Fatalf("conta inativa não pode dar erro: %v", err)
	}

	var count int64
	db.Model(&domain.PasswordResetToken{}).Count(&count)
	if count != 0 || len(mailer.pedidos) != 0 {
		t.Errorf("nada deveria ser gerado: %d tokens, %d e-mails", count, len(mailer.pedidos))
	}
}

func TestSenha_FalhaNoEnvioNaoViraErroParaQuemPediu(t *testing.T) {
	svc, mailer, _, user := novoServicoDeSenha(t)
	mailer.err = errors.New("api.core indisponível")

	if err := svc.ForgotPassword(user.Email, "1.1.1.1"); err != nil {
		t.Errorf("falha de envio só vai para o log: %v", err)
	}
}

func TestSenha_LimiteDeTresPedidosPorHoraEPedidoNovoInvalidaOAnterior(t *testing.T) {
	svc, mailer, _, user := novoServicoDeSenha(t)

	for i := 0; i < 5; i++ {
		if err := svc.ForgotPassword(user.Email, "1.1.1.1"); err != nil {
			t.Fatal(err)
		}
	}
	if len(mailer.pedidos) != 3 {
		t.Fatalf("só 3 pedidos por hora deveriam enviar e-mail, enviou %d", len(mailer.pedidos))
	}

	// O primeiro link foi invalidado pelo segundo; só o último ainda vale.
	if err := svc.ResetPassword(mailer.pedidos[0].Token, "nova-senha-1"); !errors.Is(err, ErrInvalidResetToken) {
		t.Errorf("link antigo deveria estar invalidado: %v", err)
	}
	if err := svc.ResetPassword(mailer.pedidos[2].Token, "nova-senha-1"); err != nil {
		t.Errorf("o último link deveria valer: %v", err)
	}
}

func TestSenha_RedefinirTrocaASenhaEOTokenValeUmaVez(t *testing.T) {
	svc, mailer, db, user := novoServicoDeSenha(t)
	if err := svc.ForgotPassword(user.Email, "1.1.1.1"); err != nil {
		t.Fatal(err)
	}
	token := mailer.pedidos[0].Token

	if err := svc.ResetPassword(token, "senha-nova-123"); err != nil {
		t.Fatalf("reset: %v", err)
	}
	var updated domain.User
	db.First(&updated, user.ID)
	if !utils.CheckPasswordHash("senha-nova-123", updated.PasswordHash) || utils.CheckPasswordHash("senha-antiga", updated.PasswordHash) {
		t.Error("a senha deveria ter sido trocada")
	}

	if err := svc.ResetPassword(token, "outra-senha-456"); !errors.Is(err, ErrInvalidResetToken) {
		t.Errorf("segundo uso do mesmo token deveria falhar: %v", err)
	}
	db.First(&updated, user.ID)
	if !utils.CheckPasswordHash("senha-nova-123", updated.PasswordHash) {
		t.Error("o segundo uso não pode alterar a senha")
	}
}

func TestSenha_TokenInvalidoExpiradoOuDeContaInativaEhRecusado(t *testing.T) {
	svc, mailer, db, user := novoServicoDeSenha(t)

	if err := svc.ResetPassword("token-que-nao-existe", "senha-nova-123"); !errors.Is(err, ErrInvalidResetToken) {
		t.Errorf("token inexistente: %v", err)
	}

	if err := svc.ForgotPassword(user.Email, "1.1.1.1"); err != nil {
		t.Fatal(err)
	}
	token := mailer.pedidos[0].Token

	// Expirado
	db.Model(&domain.PasswordResetToken{}).Where("user_id = ?", user.ID).Update("expires_at", time.Now().Add(-time.Minute))
	if err := svc.ResetPassword(token, "senha-nova-123"); !errors.Is(err, ErrInvalidResetToken) {
		t.Errorf("token expirado: %v", err)
	}

	// Conta desativada depois do pedido
	db.Model(&domain.PasswordResetToken{}).Where("user_id = ?", user.ID).Update("expires_at", time.Now().Add(10*time.Minute))
	db.Model(&domain.User{}).Where("id = ?", user.ID).Update("is_active", false)
	if err := svc.ResetPassword(token, "senha-nova-123"); !errors.Is(err, ErrInvalidResetToken) {
		t.Errorf("conta inativa: %v", err)
	}
	var open int64
	db.Model(&domain.PasswordResetToken{}).Where("user_id = ? AND used_at IS NULL", user.ID).Count(&open)
	if open != 1 {
		t.Errorf("a falha deveria desfazer a marcação de uso do token (abertos = %d)", open)
	}
}
