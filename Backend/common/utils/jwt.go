package utils

import (
	"errors"
	"time"

	"lumini-hub/common/config"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims representa os claims do token JWT.
//
// Não carrega a lista de permissões do usuário — ela crescia sem limite
// conforme o usuário acumulava permissões (chegando a inviabilizar
// cookies/URIs grandes, um problema real para o futuro uso de WebSockets).
// RequirePermission (common/middlewares/permission.go) consulta a tabela
// user_permissions direto no banco a cada requisição em vez de confiar numa
// cópia embutida no token.
type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	RoleID   uint   `json:"role_id"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateAccessToken gera um novo token JWT de acesso
func GenerateAccessToken(userID uint, username string, roleID uint, role string, cfg *config.Config) (string, error) {
	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		RoleID:   roleID,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.JWT.AccessTokenExp)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "lumini-hub",
			Subject:   username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWT.Secret))
}

// GenerateRefreshToken gera um novo token JWT de refresh
func GenerateRefreshToken(userID uint, username string, cfg *config.Config) (string, error) {
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.JWT.RefreshTokenExp)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Issuer:    "lumini-hub",
		Subject:   username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWT.Secret))
}

// ValidateToken valida um token JWT
func ValidateToken(tokenString string, cfg *config.Config) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("token inválido")
}
