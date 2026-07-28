package middlewares

import (
	"net/http"

	"lumini-hub/common/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// permissionDB é o banco usado por RequirePermission para consultar as
// permissões atuais do usuário a cada requisição, em vez de confiar numa
// lista embutida no JWT. Configurado uma vez por microsserviço via
// InitPermissionChecker — mesmo espírito de AuthMiddleware receber cfg.
var permissionDB *gorm.DB

// InitPermissionChecker configura o banco usado por RequirePermission. Deve
// ser chamado uma vez no main.go de cada microsserviço, logo após
// database.InitDB.
func InitPermissionChecker(db *gorm.DB) {
	permissionDB = db
}

// RequirePermission verifica se o usuário tem a permissão necessária.
//
// Bypass hierárquico: "DEVELOP" (utils.RoleDeveloper) passa em qualquer
// checagem. "ADMIN" passa em quase tudo, EXCETO nas permissões do catálogo
// de Perfis e Permissões (utils.IsDeveloperOnlyPermission) — esse cadastro
// fica reservado ao DEVELOP. Fora desses bypasses, consulta user_permissions
// diretamente no banco (join indexado por user_id + permission).
func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDValue, exists := c.Get("userID")
		if !exists {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Não autorizado", "Usuário não autenticado")
			c.Abort()
			return
		}
		userID, ok := userIDValue.(uint)
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

		hasPermission, err := userHasPermission(userID, permission)
		if err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Erro interno", "Erro ao verificar permissões")
			c.Abort()
			return
		}

		if !hasPermission {
			utils.ErrorResponse(c, http.StatusForbidden, "Acesso negado", "Permissão insuficiente")
			c.Abort()
			return
		}

		c.Next()
	}
}

// userHasPermission consulta user_permissions + permissions direto (sem
// depender dos modelos Go de domínio de nenhum microsserviço específico,
// já que common/middlewares é compartilhado por todos eles).
func userHasPermission(userID uint, permission string) (bool, error) {
	var count int64
	err := permissionDB.Table("user_permissions AS up").
		Joins("JOIN permissions AS p ON p.id = up.permission_id").
		Where("up.user_id = ? AND p.permission = ?", userID, permission).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
