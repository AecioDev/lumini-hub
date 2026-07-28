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

// DeveloperModule é o valor de Permission.Module usado pras permissões que
// só fazem sentido pro DEVELOP (roles.*, permissions.*, admin.create_permissions
// — ver seed em erp_system.permissions). Fica escondido de qualquer
// requisitante que não seja DEVELOP em toda consulta ao catálogo de
// Permissões (lista, agrupamento por módulo, lista de módulos) — do
// contrário um ADMIN poderia enxergar/atribuir essas permissões a um Perfil
// customizado e conceder, sem querer, poderes de DEVELOP pra esse Perfil.
const DeveloperModule = "Develop"

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

// IsDeveloperOnlyPermission indica se uma permissão é reservada ao perfil
// "DEVELOP". O catálogo de Perfis (`roles.*` — criar/editar/excluir/atribuir
// permissões) é do ADMIN: ele decide como agrupar as permissões da empresa
// em Perfis, sem restrição — quem manda é a permissão em si, checada em
// tempo real pelo middleware. O que fica exclusivo do DEVELOP é:
//   - `permissions.create`/`permissions.edit`/`permissions.delete`: mexer no
//     catálogo de Permissões em si (quais capacidades existem no sistema) —
//     criar uma permissão sem um RequirePermission("...") correspondente no
//     código não faz nada, então é ferramenta de dev. `permissions.view`
//     fica de fora do bloqueio (mesmo motivo de `roles.view`): é só leitura,
//     e o próprio módulo DeveloperModule já sai da resposta pra quem não é
//     DEVELOP (ver GetPermissions/excludeModuleFor), então não vaza nada.
//   - `admin.create_permissions`: CRUD de /menu-items (árvore do menu lateral).
//
// Guardas específicas contra abuso do poder amplo sobre `roles.*` (não
// deixar o ADMIN criar/renomear um Perfil pra "DEVELOP", nem excluir o
// Perfil "ADMIN") ficam em RoleValidator, não aqui — este helper só decide
// QUEM pode chamar o endpoint.
func IsDeveloperOnlyPermission(permission string) bool {
	if permission == "permissions.view" {
		return false
	}
	return permission == "admin.create_permissions" ||
		strings.HasPrefix(permission, "permissions.")
}
