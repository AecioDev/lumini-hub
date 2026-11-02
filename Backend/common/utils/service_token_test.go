package utils

import (
	"testing"
	"time"

	"lumini-hub/common/config"

	"github.com/golang-jwt/jwt/v5"
)

func TestServiceToken_IdaEVolta(t *testing.T) {
	cfg := &config.Config{JWT: config.JWTConfig{Secret: "segredo-de-teste"}}

	token, err := GenerateServiceToken("api.auth", cfg)
	if err != nil {
		t.Fatalf("gerar: %v", err)
	}
	caller, err := ValidateServiceToken(token, cfg)
	if err != nil || caller != "api.auth" {
		t.Fatalf("validar = %q, %v", caller, err)
	}
}

func TestServiceToken_RecusaTokenDeUsuarioSegredoErradoEVencido(t *testing.T) {
	cfg := &config.Config{JWT: config.JWTConfig{Secret: "segredo-de-teste", AccessTokenExp: time.Minute}}

	// Um token de acesso de usuário (sem a audiência de serviço) não vale como token de serviço.
	userToken, _ := GenerateAccessToken(1, "ana", 1, "ADMIN", cfg)
	if _, err := ValidateServiceToken(userToken, cfg); err == nil {
		t.Error("token de usuário não pode valer como token de serviço")
	}

	// Segredo diferente.
	token, _ := GenerateServiceToken("api.auth", cfg)
	if _, err := ValidateServiceToken(token, &config.Config{JWT: config.JWTConfig{Secret: "outro"}}); err == nil {
		t.Error("segredo errado deveria falhar")
	}

	// Vencido.
	claims := jwt.RegisteredClaims{Subject: "api.auth", Audience: jwt.ClaimStrings{serviceTokenAudience}, ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Minute))}
	expired, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(cfg.JWT.Secret))
	if _, err := ValidateServiceToken(expired, cfg); err == nil {
		t.Error("token vencido deveria falhar")
	}

	// Sem validade.
	forever := jwt.RegisteredClaims{Subject: "api.auth", Audience: jwt.ClaimStrings{serviceTokenAudience}}
	noExp, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, forever).SignedString([]byte(cfg.JWT.Secret))
	if _, err := ValidateServiceToken(noExp, cfg); err == nil {
		t.Error("token sem validade deveria falhar")
	}
}
