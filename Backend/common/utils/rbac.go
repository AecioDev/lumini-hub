package utils

import (
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

// UserHasPermission consulta user_permissions + permissions diretamente
// (sem depender dos modelos Go de domínio de nenhum microsserviço
// específico — só nomes de tabela, mesmo truque que RequirePermission já
// fazia antes desta função ser extraída pra cá). Não aplica o bypass
// hierárquico DEVELOP/ADMIN (isso é decisão de gate de rota, não faz
// sentido pra checagens de negócio como visibilidade de dados — ver
// ResolveVisibleCompanyIDs). `db` é a conexão do microsserviço chamador
// (mesmo banco físico compartilhado nesta fase de migração).
func UserHasPermission(db *gorm.DB, userID uint, permission string) (bool, error) {
	var count int64
	err := db.Table("user_permissions AS up").
		Joins("JOIN permissions AS p ON p.id = up.permission_id").
		Where("up.user_id = ? AND p.permission = ?", userID, permission).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
