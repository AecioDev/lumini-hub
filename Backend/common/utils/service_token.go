package utils

import (
	"errors"
	"time"

	"lumini-hub/common/config"

	"github.com/golang-jwt/jwt/v5"
)

// Token de serviço: a credencial das chamadas internas entre microsserviços (por exemplo, o api.auth
// pedindo ao api.core o envio de um e-mail). Mesmo JWT_SECRET dos tokens de usuário, mas com uma audiência
// própria e validade de segundos: um token de usuário não passa como token de serviço e vice-versa.
const (
	serviceTokenAudience = "lumini-hub-internal"
	serviceTokenTTL      = 30 * time.Second
)

// GenerateServiceToken gera o token de serviço de curta duração. service é o nome de quem chama (ex.: "api.auth").
func GenerateServiceToken(service string, cfg *config.Config) (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Issuer:    "lumini-hub",
		Subject:   service,
		Audience:  jwt.ClaimStrings{serviceTokenAudience},
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now.Add(-5 * time.Second)), // tolera relógios levemente fora de sincronia
		ExpiresAt: jwt.NewNumericDate(now.Add(serviceTokenTTL)),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(cfg.JWT.Secret))
}

// ValidateServiceToken confere o token de serviço e devolve o nome de quem chamou.
func ValidateServiceToken(tokenString string, cfg *config.Config) (string, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		return []byte(cfg.JWT.Secret), nil
	}, jwt.WithValidMethods([]string{"HS256"}), jwt.WithAudience(serviceTokenAudience), jwt.WithExpirationRequired())
	if err != nil {
		return "", err
	}
	if !token.Valid || claims.Subject == "" {
		return "", errors.New("token de serviço inválido")
	}
	return claims.Subject, nil
}
