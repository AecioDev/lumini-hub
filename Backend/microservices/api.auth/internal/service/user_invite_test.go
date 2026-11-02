package service

import (
	"errors"
	"testing"
	"time"

	"lumini-hub/api.auth/internal/domain"
	"lumini-hub/api.auth/internal/repository"
	"lumini-hub/common/testutil"
	"lumini-hub/common/utils"

	"gorm.io/gorm"
)

// Convite de novo usuário e troca obrigatória de senha (precisa de Postgres: sem TEST_DATABASE_URL o teste é ignorado).

func novoServicoDeUsuario(t *testing.T) (*UserService, *PasswordResetService, *mailerFalso, *gorm.DB, domain.Role) {
	t.Helper()
	db := testutil.OpenTestDB(t)
	if err := db.AutoMigrate(&domain.Permission{}, &domain.Role{}, &domain.User{}, &domain.PasswordResetToken{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	role := domain.Role{Name: "VENDEDOR"}
	if err := db.Create(&role).Error; err != nil {
		t.Fatal(err)
	}
	mailer := &mailerFalso{}
	reset := NewPasswordResetService(db, mailer)
	reset.dispatch = func(f func()) { f() }
	svc := NewUserService(repository.NewUserRepository(db), repository.NewRoleRepository(db), repository.NewPermissionRepository(db), reset)
	return svc, reset, mailer, db, role
}

func pedidoDeUsuario(role domain.Role) domain.CreateUserRequest {
	return domain.CreateUserRequest{Username: "novo.usuario", Name: "Novo", Email: "novo@exemplo.com", RoleID: role.ID}
}

func TestConvite_UsuarioNasceSemSenhaUtilizavelEOEmailSai(t *testing.T) {
	svc, _, mailer, db, role := novoServicoDeUsuario(t)
	req := pedidoDeUsuario(role)
	req.SendInvite = true
	req.MustChangePassword = true // sem efeito com convite

	created, err := svc.CreateUser(req, "ADMIN", "admin")
	if err != nil {
		t.Fatalf("criar: %v", err)
	}
	if !created.InviteSent || created.InviteError != "" || !created.PendingInvite {
		t.Errorf("resposta inesperada: %+v", created)
	}

	if len(mailer.convites) != 1 {
		t.Fatalf("deveria pedir 1 e-mail de convite, pediu %d", len(mailer.convites))
	}
	c := mailer.convites[0]
	if c.To != "novo@exemplo.com" || c.Username != "novo.usuario" || c.ExpiresDays != 7 || c.InvitedBy != "admin" || len(c.Token) < 40 {
		t.Errorf("convite inesperado: %+v", c)
	}

	var user domain.User
	db.First(&user, created.ID)
	if user.MustChangePassword {
		t.Error("com convite não há troca obrigatória: o usuário define a senha pelo link")
	}
	if utils.CheckPasswordHash("", user.PasswordHash) || utils.CheckPasswordHash(c.Token, user.PasswordHash) {
		t.Error("a senha inicial não pode ser utilizável")
	}

	var token domain.PasswordResetToken
	db.Where("user_id = ?", user.ID).First(&token)
	if token.Purpose != domain.TokenPurposeInvite || token.TokenHash != hashResetToken(c.Token) || time.Until(token.ExpiresAt) < 6*24*time.Hour {
		t.Errorf("token inesperado: %+v", token)
	}
}

func TestConvite_AceitarDefineASenhaELimpaATrocaObrigatoria(t *testing.T) {
	svc, reset, mailer, db, role := novoServicoDeUsuario(t)
	req := pedidoDeUsuario(role)
	req.SendInvite = true
	created, err := svc.CreateUser(req, "ADMIN", "admin")
	if err != nil {
		t.Fatal(err)
	}
	db.Model(&domain.User{}).Where("id = ?", created.ID).Update("must_change_password", true) // cenário defensivo

	if err := reset.ResetPassword(mailer.convites[0].Token, "minha-senha-1"); err != nil {
		t.Fatalf("aceitar convite: %v", err)
	}
	var user domain.User
	db.First(&user, created.ID)
	if !utils.CheckPasswordHash("minha-senha-1", user.PasswordHash) || user.MustChangePassword {
		t.Errorf("senha definida e troca obrigatória desligada esperadas: %+v", user)
	}
	if err := reset.ResetPassword(mailer.convites[0].Token, "outra-senha-2"); !errors.Is(err, ErrInvalidResetToken) {
		t.Errorf("o link do convite vale uma vez: %v", err)
	}
}

func TestConvite_FalhaNoEmailNaoDesfazOCadastroEPermiteReenviar(t *testing.T) {
	svc, _, mailer, db, role := novoServicoDeUsuario(t)
	req := pedidoDeUsuario(role)
	req.SendInvite = true
	mailer.err = errors.New("api.core recusou o envio (HTTP 422): endereço do sistema não configurado")

	created, err := svc.CreateUser(req, "ADMIN", "admin")
	if err != nil {
		t.Fatalf("o cadastro não pode falhar por causa do e-mail: %v", err)
	}
	if created.InviteSent || created.InviteError == "" {
		t.Errorf("deveria informar a falha do envio: %+v", created)
	}
	var count int64
	db.Model(&domain.User{}).Where("email = ?", "novo@exemplo.com").Count(&count)
	if count != 1 {
		t.Fatalf("o usuário deveria existir (%d)", count)
	}

	// Reenviar com o e-mail já funcionando: invalida o link anterior e manda o novo.
	mailer.err = nil
	if err := svc.ResendInvite(created.ID, "ADMIN", "admin"); err != nil {
		t.Fatalf("reenviar: %v", err)
	}
	if len(mailer.convites) != 1 {
		t.Fatalf("deveria ter enviado 1 convite, enviou %d", len(mailer.convites))
	}
	var open int64
	db.Model(&domain.PasswordResetToken{}).Where("user_id = ? AND used_at IS NULL", created.ID).Count(&open)
	if open != 1 {
		t.Errorf("só o link mais novo deveria ficar aberto (%d)", open)
	}
}

func TestConvite_ReenvioSoParaQuemNuncaEntrouEEstaAtivo(t *testing.T) {
	svc, _, _, db, role := novoServicoDeUsuario(t)
	req := pedidoDeUsuario(role)
	req.SendInvite = true
	created, _ := svc.CreateUser(req, "ADMIN", "admin")

	db.Model(&domain.User{}).Where("id = ?", created.ID).Update("last_login", time.Now())
	if err := svc.ResendInvite(created.ID, "ADMIN", "admin"); !errors.Is(err, ErrInviteNotAllowed) {
		t.Errorf("quem já entrou não recebe convite: %v", err)
	}

	db.Model(&domain.User{}).Where("id = ?", created.ID).Updates(map[string]any{"last_login": nil, "is_active": false})
	if err := svc.ResendInvite(created.ID, "ADMIN", "admin"); !errors.Is(err, ErrInviteNotAllowed) {
		t.Errorf("usuário inativo não recebe convite: %v", err)
	}
	if err := svc.ResendInvite(9999, "ADMIN", "admin"); !errors.Is(err, utils.ErrNotFound) {
		t.Errorf("usuário inexistente: %v", err)
	}
}

func TestCadastro_SemConviteExigeSenhaETrocaObrigatoriaValeSoSemConvite(t *testing.T) {
	svc, _, mailer, db, role := novoServicoDeUsuario(t)

	req := pedidoDeUsuario(role)
	if _, err := svc.CreateUser(req, "ADMIN", "admin"); err == nil {
		t.Fatal("sem convite e sem senha deveria falhar")
	}

	req.Password = "senha-inicial"
	req.MustChangePassword = true
	created, err := svc.CreateUser(req, "ADMIN", "admin")
	if err != nil {
		t.Fatalf("criar: %v", err)
	}
	if created.InviteSent || len(mailer.convites) != 0 {
		t.Error("sem convite nenhum e-mail sai")
	}
	var user domain.User
	db.First(&user, created.ID)
	if !user.MustChangePassword || !utils.CheckPasswordHash("senha-inicial", user.PasswordHash) {
		t.Errorf("troca obrigatória e senha digitada esperadas: %+v", user)
	}
}

func TestTrocaDeSenha_OProprioUsuarioDesligaATrocaObrigatoria(t *testing.T) {
	svc, _, _, db, role := novoServicoDeUsuario(t)
	req := pedidoDeUsuario(role)
	req.Password = "senha-inicial"
	req.MustChangePassword = true
	created, err := svc.CreateUser(req, "ADMIN", "admin")
	if err != nil {
		t.Fatal(err)
	}

	// Um administrador redefinindo a senha de outra pessoa não cumpre a troca obrigatória dela.
	if err := svc.ChangePassword(created.ID, "", "senha-do-admin", true); err != nil {
		t.Fatal(err)
	}
	var user domain.User
	db.First(&user, created.ID)
	if !user.MustChangePassword {
		t.Error("troca feita por administrador não desliga a marca")
	}

	if err := svc.ChangePassword(created.ID, "senha-do-admin", "senha-propria-9", false); err != nil {
		t.Fatalf("o próprio usuário trocando: %v", err)
	}
	db.First(&user, created.ID)
	if user.MustChangePassword || !utils.CheckPasswordHash("senha-propria-9", user.PasswordHash) {
		t.Errorf("a troca pelo próprio usuário deveria desligar a marca: %+v", user)
	}
}
