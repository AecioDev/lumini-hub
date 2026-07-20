package middlewares

import (
	"net/http"

	"lumini-hub/common/utils"

	"github.com/gin-gonic/gin"
)

// RequirePermission verifica se o usuário tem a permissão necessária.
//
// Bypass hierárquico: "DEVELOP" (utils.RoleDeveloper) passa em qualquer
// checagem. "ADMIN" passa em quase tudo, EXCETO nas permissões do catálogo
// de Perfis e Permissões (utils.IsDeveloperOnlyPermission) — esse cadastro
// fica reservado ao DEVELOP.
func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Verificar se o usuário está autenticado
		permissions, exists := c.Get("permissions")
		if !exists {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Não autorizado", "Usuário não autenticado")
			c.Abort()
			return
		}

		// Verificar se o usuário tem a permissão necessária
		userPermissions, ok := permissions.([]string)
		if !ok {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Erro interno", "Erro ao verificar permissões")
			c.Abort()
			return
		}

		if roleValue, exists := c.Get("role"); exists {
			role, _ := roleValue.(string)
			if role == utils.RoleDeveloper {
				c.Next()
				return
			}
			if role == utils.RoleAdmin && !utils.IsDeveloperOnlyPermission(permission) {
				c.Next()
				return
			}
		}

		// Verificar se o usuário tem a permissão específica
		hasPermission := false
		for _, p := range userPermissions {
			if p == permission {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			utils.ErrorResponse(c, http.StatusForbidden, "Acesso negado", "Permissão insuficiente")
			c.Abort()
			return
		}

		c.Next()
	}
}
