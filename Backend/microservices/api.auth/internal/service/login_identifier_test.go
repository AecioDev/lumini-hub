package service

import (
	"errors"
	"testing"

	"lumini-hub/api.auth/internal/domain"
	"lumini-hub/common/config"
	"lumini-hub/common/testutil"
	"lumini-hub/common/utils"

	"gorm.io/gorm"
)

// O login aceita o nome de usuário ou o e-mail (precisa de Postgres: sem TEST_DATABASE_URL o teste é ignorado).
func TestLogin_BuscaPorUsernameOuEmail(t *testing.T) {
	db := testutil.OpenTestDB(t)
	if err := db.AutoMigrate(&domain.Permission{}, &domain.Role{}, &domain.User{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	role := domain.Role{Name: "ADMIN"}
	if err := db.Create(&role).Error; err != nil {
		t.Fatal(err)
	}
	ana := domain.User{Username: "ana", Name: "Ana", Email: "Ana@Exemplo.com", RoleID: role.ID, IsActive: true}
	// O username de Bia é igual ao e-mail de Ana: a busca por username tem prioridade.
	bia := domain.User{Username: "ana@exemplo.com", Name: "Bia", Email: "bia@exemplo.com", RoleID: role.ID, IsActive: true}
	for _, u := range []*domain.User{&ana, &bia} {
		if err := db.Create(u).Error; err != nil {
			t.Fatal(err)
		}
	}
	svc := NewAuthService(db, &config.Config{})

	cases := []struct {
		digitado string
		quem     string
	}{
		{"ana", "Ana"},
		{"ANA", "Ana"},
		{"  bia@EXEMPLO.com ", "Bia"},
		{"ana@exemplo.com", "Bia"}, // username de Bia vence o e-mail de Ana
	}
	for _, c := range cases {
		got, err := svc.findUserForLogin(c.digitado)
		if err != nil || got.Name != c.quem {
			t.Errorf("%q: veio %v, %v; esperava %s", c.digitado, got, err, c.quem)
		}
	}

	if _, err := svc.findUserForLogin("ninguem@exemplo.com"); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("desconhecido deveria ser ErrRecordNotFound: %v", err)
	}
	got, err := svc.findUserForLogin("Ana@exemplo.com")
	if err != nil || got.Name != "Bia" {
		t.Errorf("o texto igual ao username de Bia deveria entrar como Bia: %v, %v", got, err)
	}
	if got, err = svc.findUserForLogin("ana@exemplo.com"); err != nil || got.Role == nil || got.Role.Name != "ADMIN" {
		t.Errorf("o perfil deveria vir carregado: %+v, %v", got, err)
	}
}

// CFG-047: usuário inexistente e senha errada recebem a mesma resposta; "inativo" só aparece para quem acertou a senha.
func TestLogin_RespostaUnicaParaCredenciaisInvalidas(t *testing.T) {
	db := testutil.OpenTestDB(t)
	if err := db.AutoMigrate(&domain.Permission{}, &domain.Role{}, &domain.User{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	role := domain.Role{Name: "ADMIN"}
	if err := db.Create(&role).Error; err != nil {
		t.Fatal(err)
	}
	hash, err := utils.HashPassword("senha-correta")
	if err != nil {
		t.Fatal(err)
	}
	ativo := domain.User{Username: "ana", Name: "Ana", Email: "ana@exemplo.com", RoleID: role.ID, IsActive: true, PasswordHash: hash}
	inativo := domain.User{Username: "bia", Name: "Bia", Email: "bia@exemplo.com", RoleID: role.ID, IsActive: true, PasswordHash: hash}
	for _, u := range []*domain.User{&ativo, &inativo} {
		if err := db.Create(u).Error; err != nil {
			t.Fatal(err)
		}
	}
	// IsActive tem default true no modelo: o false precisa ser gravado depois
	if err := db.Model(&inativo).Update("is_active", false).Error; err != nil {
		t.Fatal(err)
	}
	svc := NewAuthService(db, &config.Config{})

	for _, c := range []struct{ usuario, senha string }{
		{"ninguem", "qualquer"},
		{"ana", "senha-errada"},
		{"bia", "senha-errada"}, // conta inativa com senha errada não revela que está inativa
	} {
		if _, err := svc.Login(c.usuario, c.senha); !errors.Is(err, ErrInvalidCredentials) {
			t.Errorf("Login(%q, %q) = %v, esperava ErrInvalidCredentials", c.usuario, c.senha, err)
		}
	}
	if _, err := svc.Login("bia", "senha-correta"); err == nil || errors.Is(err, ErrInvalidCredentials) {
		t.Errorf("conta inativa com a senha certa deveria dizer que está inativa: %v", err)
	}
}
