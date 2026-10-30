package domain

import (
	"strings"
	"testing"

	"lumini-hub/common/testutil"

	"github.com/gin-gonic/gin/binding"
)

// Todo campo com gorm "size:N" que aparece numa request de criar/editar deve ter binding "max=N": sem
// ele o backend aceita o excesso e só o banco recusa. Ao criar uma entidade nova no api.auth,
// acrescente-a aqui.
func TestRequests_MaxAlinhadoAoTamanhoDaColuna(t *testing.T) {
	testutil.AssertRequestMaxMatchesColumnSize(t, User{}, CreateUserRequest{}, UpdateUserRequest{})
	testutil.AssertRequestMaxMatchesColumnSize(t, Role{}, CreateRoleRequest{}, UpdateRoleRequest{})
	testutil.AssertRequestMaxMatchesColumnSize(t, MenuItem{}, CreateMenuItemRequest{}, UpdateMenuItemRequest{})
}

// A senha vira bcrypt, que recusa mais de 72 bytes: o limite é checado antes, na request.
func TestRequests_SenhaAcimaDoLimiteDoBcryptEhRecusada(t *testing.T) {
	user := CreateUserRequest{Username: "ana", Password: strings.Repeat("a", 72), Name: "Ana", Email: "ana@x.com", RoleID: 1}
	if err := binding.Validator.ValidateStruct(user); err != nil {
		t.Fatalf("72 caracteres deveria passar: %v", err)
	}
	user.Password = strings.Repeat("a", 73)
	if err := binding.Validator.ValidateStruct(user); err == nil {
		t.Error("73 caracteres deveria ser recusado")
	}

	change := ChangePasswordRequest{CurrentPassword: "x", NewPassword: strings.Repeat("a", 73)}
	if err := binding.Validator.ValidateStruct(change); err == nil {
		t.Error("nova senha com 73 caracteres deveria ser recusada")
	}
}

func TestRequests_ExcessoDeTamanhoEhRecusadoAntesDoBanco(t *testing.T) {
	text := func(n int) string { return strings.Repeat("ã", n) }
	cases := []struct {
		name string
		make func(n int) any
		max  int
	}{
		{"usuário: nome", func(n int) any {
			return CreateUserRequest{Username: "ana", Password: "123456", Name: text(n), Email: "a@x.com", RoleID: 1}
		}, 100},
		{"usuário: telefone", func(n int) any {
			return CreateUserRequest{Username: "ana", Password: "123456", Name: "Ana", Email: "a@x.com", Phone: text(n), RoleID: 1}
		}, 20},
		{"perfil: nome", func(n int) any { return CreateRoleRequest{Name: text(n)} }, 50},
		{"item de menu: href", func(n int) any { return CreateMenuItemRequest{Name: "Menu", Href: text(n)} }, 255},
	}
	for _, tc := range cases {
		if err := binding.Validator.ValidateStruct(tc.make(tc.max)); err != nil {
			t.Errorf("%s com %d caracteres deveria passar: %v", tc.name, tc.max, err)
		}
		if err := binding.Validator.ValidateStruct(tc.make(tc.max + 1)); err == nil {
			t.Errorf("%s com %d caracteres deveria ser recusado", tc.name, tc.max+1)
		}
	}
}
