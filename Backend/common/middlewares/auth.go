package middlewares

import (
	"log"
	"net/http"
	"strings"

	"lumini-hub/common/config"
	"lumini-hub/common/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// userStillActive confere no banco se a conta do token continua existindo e ativa. É o que faz desativar
// ou excluir um usuário valer na hora: o token de acesso (JWT) é válido por minutos e, sozinho, não sabe
// que a conta mudou. É uma busca pela chave primária, a mesma ordem de custo da consulta de permissões que
// o RequirePermission já faz a cada requisição. Sem banco configurado (permissionDB nulo, só em testes que
// não chamam InitPermissionChecker), não há o que conferir.
func userStillActive(userID uint) (bool, error) {
	if permissionDB == nil {
		return true, nil
	}
	var row struct{ IsActive bool }
	err := permissionDB.Table("users").Select("is_active").
		Where("id = ? AND deleted_at IS NULL", userID).Take(&row).Error
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return row.IsActive, nil
}

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

		active, err := userStillActive(claims.UserID)
		if err != nil {
			log.Printf("Middleware: erro ao conferir a conta do usuário %d: %v", claims.UserID, err)
			utils.ErrorResponse(c, http.StatusInternalServerError, "Erro interno", "Erro ao verificar a sessão")
			c.Abort()
			return
		}
		if !active {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Não autorizado", "Sessão encerrada")
			c.Abort()
			return
		}

		// Armazenar claims no contexto
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("roleID", claims.RoleID)
		c.Set("role", claims.Role)

		log.Printf("Middleware: Token validado com sucesso para UserID: %d", claims.UserID)
		c.Next()
	}
}
