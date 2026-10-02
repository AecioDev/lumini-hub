package config

import (
	"fmt"
	"strings"
	"testing"
)

const secretPassword = "SENHA-SECRETA-123"

// A senha do admin inicial não pode vazar se alguém imprimir o Config (log de debug,
// erro formatado com %+v etc.).
func TestBootstrapAdminConfig_NaoVazaSenhaAoImprimir(t *testing.T) {
	cfg := &Config{Bootstrap: BootstrapAdminConfig{Username: "admin", Password: secretPassword, Email: "a@b.c"}}

	for name, value := range map[string]any{
		"*Config": cfg,
		"Config":  *cfg,
		"campo":   cfg.Bootstrap,
		"campo*":  &cfg.Bootstrap,
	} {
		for _, verb := range []string{"%v", "%+v", "%#v", "%s"} {
			if out := fmt.Sprintf(verb, value); strings.Contains(out, secretPassword) {
				t.Errorf("%s com %s vazou a senha: %s", name, verb, out)
			}
		}
	}
}

func TestLoad_BootstrapSemSenhaPadrao(t *testing.T) {
	t.Setenv("BOOTSTRAP_ADMIN_PASSWORD", "")
	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Bootstrap.Password != "" {
		t.Error("não pode haver senha default para o admin inicial")
	}
	if cfg.Bootstrap.Username != "admin" {
		t.Errorf("username default = %q, esperado admin", cfg.Bootstrap.Username)
	}
}

func TestLoad_AutoMigrate(t *testing.T) {
	cases := []struct {
		env     string
		want    bool
		wantErr bool
	}{
		{"", true, false},
		{"true", true, false},
		{"false", false, false},
		{"flase", false, true},
	}
	for _, tc := range cases {
		t.Run("DB_AUTO_MIGRATE="+tc.env, func(t *testing.T) {
			t.Setenv("DB_AUTO_MIGRATE", tc.env)
			cfg, err := Load()
			if (err != nil) != tc.wantErr {
				t.Fatalf("erro = %v, esperava erro: %v", err, tc.wantErr)
			}
			if err == nil && cfg.Database.AutoMigrate != tc.want {
				t.Errorf("AutoMigrate = %v, esperado %v", cfg.Database.AutoMigrate, tc.want)
			}
		})
	}
}
