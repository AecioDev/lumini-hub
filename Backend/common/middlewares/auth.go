package middlewares

import (
	"log"
	"net/http"
	"strings"

	"lumini-hub/common/config"
	"lumini-hub/common/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware verifica se o usuário está autenticado
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// 1. Tentar obter o token do cookie (nova forma)
		accessTokenFromCookie, err := c.Cookie("access_token")
		if err == nil && accessTokenFromCookie != "" {
			tokenString = accessTokenFromCookie
			log.Printf("Middleware: Token de acesso encontrado no cookie.")
		} else {
			log.Printf("Middleware: Token de acesso NÃO encontrado no cookie (erro: %v).", err)
			// 2. Se não encontrou no cookie, tentar obter do cabeçalho Authorization (fallback)
			authHeader := c.GetHeader("Authorization")
			if authHeader != "" {
				log.Printf("Middleware: Token encontrado no cabeçalho Authorization.")
				parts := strings.Split(authHeader, " ")
				if len(parts) == 2 && parts[0] == "Bearer" {
					tokenString = parts[1]
				}
			}
		}

		if tokenString == "" {
			log.Printf("Middleware: Nenhum token de autenticação encontrado.")
			utils.ErrorResponse(c, http.StatusUnauthorized, "Não autorizado", "Token não fornecido")
			c.Abort()
			return
		}

		// Validar token
		claims, err := utils.ValidateToken(tokenString, cfg)
		if err != nil {
			log.Printf("Middleware: Erro ao validar token: %v", err)
			utils.ErrorResponse(c, http.StatusUnauthorized, "Não autorizado", err.Error())
			c.Abort()
			return
		}

		// Armazenar claims no contexto
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("roleID", claims.RoleID)
		c.Set("role", claims.Role)
		c.Set("permissions", claims.Permissions)

		log.Printf("Middleware: Token validado com sucesso para UserID: %d", claims.UserID)
		c.Next()
	}
}
