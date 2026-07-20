package utils

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// Papéis com tratamento especial no RBAC hierárquico.
const (
	RoleAdmin     = "ADMIN"
	RoleDeveloper = "DEVELOP"
)

// RoleFromGinContext extrai o papel do usuário autenticado do contexto do
// Gin, setado pelo AuthMiddleware a partir do JWT. Retorna "" se não houver
// (não deveria acontecer atrás de AuthMiddleware, só evita panic de type
// assertion).
func RoleFromGinContext(c *gin.Context) string {
	value, exists := c.Get("role")
	if !exists {
		return ""
	}
	role, _ := value.(string)
	return role
}

// IsDeveloperOnlyPermission indica se uma permissão faz parte do
// CADASTRO/MANUTENÇÃO do catálogo de RBAC (criar/editar/excluir Perfis e
// Permissões) — reservado ao perfil "DEVELOP". Não inclui `roles.view`/
// `permissions.view`: essas ficam abertas pro ADMIN porque são necessárias
// pra uma operação comum de administrador — atribuir um Perfil e escolher
// permissões ao criar/editar um USUÁRIO (telas de Usuários usam
// GET /roles e GET /permissions/by-module só pra popular esses seletores,
// não pra mexer no catálogo em si).
func IsDeveloperOnlyPermission(permission string) bool {
	if permission == "roles.view" || permission == "permissions.view" {
		return false
	}
	return permission == "admin.create_permissions" ||
		strings.HasPrefix(permission, "roles.") ||
		strings.HasPrefix(permission, "permissions.")
}
