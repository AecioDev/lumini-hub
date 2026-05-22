package utils

import (
	"github.com/gin-gonic/gin"
)

// GetUserIDFromContext retorna o ID do usuário autenticado do contexto
func GetUserIDFromContext(c *gin.Context) (uint, bool) {
	userIDRaw, exists := c.Get("userID")
	if !exists {
		return 0, false
	}

	userID, ok := userIDRaw.(uint)
	if !ok {
		return 0, false
	}

	return userID, true
}
