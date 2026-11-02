package service

import (
	"errors"
	"testing"

	"lumini-hub/api.auth/internal/domain"
	"lumini-hub/common/config"
	"lumini-hub/common/testutil"

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
