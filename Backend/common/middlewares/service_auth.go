package middlewares

import (
	"net/http"
	"strings"

	"lumini-hub/common/config"
	"lumini-hub/common/utils"

	"github.com/gin-gonic/gin"
)

// ServiceAuthMiddleware protege as rotas internas (/api/internal/...), chamadas só por outros microsserviços
// com um token de serviço (utils.GenerateServiceToken). O gateway não repassa essas rotas para fora.
func ServiceAuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		token, found := strings.CutPrefix(header, "Bearer ")
		if !found || token == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Não autorizado", "Token de serviço não fornecido")
			c.Abort()
			return
		}

		caller, err := utils.ValidateServiceToken(token, cfg)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Não autorizado", "Token de serviço inválido")
			c.Abort()
			return
		}

		c.Set("serviceCaller", caller)
		c.Next()
	}
}
