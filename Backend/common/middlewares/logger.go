package middlewares

import (
	"time"

	"lumini-hub/common/database"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// LoggerMiddleware registra informações sobre as requisições que alteram dados (POST, PUT, DELETE, etc.)
func LoggerMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Tempo de início
		startTime := time.Now()

		// Processar requisição
		c.Next()

		// Tempo de término
		endTime := time.Now()
		latency := endTime.Sub(startTime)

		// Obter informações do usuário
		var finalUserID *uint
		userID, exists := c.Get("userID")
		if exists {
			if u, ok := userID.(uint); ok {
				finalUserID = &u
			}
		}

		// Registrar log no banco de dados para operações de escrita de sucesso (CUD)
		if c.Request.Method != "GET" && c.Writer.Status() < 400 {
			sysLog := database.SystemLog{
				UserID:    finalUserID,
				Action:    c.Request.Method + " " + c.Request.URL.Path,
				IPAddress: c.ClientIP(),
				Details: database.JSONB{
					"status":     c.Writer.Status(),
					"latency_ms": latency.Milliseconds(),
					"user_agent": c.Request.UserAgent(),
				},
			}

			// Extrair informações sobre a entidade afetada
			if entityType := c.Param("entity_type"); entityType != "" {
				sysLog.EntityType = entityType
			}

			if entityID := c.Param("id"); entityID != "" {
				sysLog.EntityID = entityID
			}

			// Salvar log de forma assíncrona
			go func(log database.SystemLog) {
				db.Create(&log)
			}(sysLog)
		}
	}
}
