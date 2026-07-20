package utils

import "strings"

// Papéis com tratamento especial no RBAC hierárquico.
const (
	RoleAdmin     = "ADMIN"
	RoleDeveloper = "DEVELOP"
)

// IsDeveloperOnlyPermission indica se uma permissão faz parte do catálogo de
// RBAC (cadastro/manutenção de Perfis e Permissões) — reservado ao perfil
// "DEVELOP". O bypass do ADMIN não vale para essas: criar um Perfil ou uma
// Permissão sem existir um RequirePermission(...) correspondente já no
// código não tem efeito nenhum, então esse catálogo é ferramenta de
// desenvolvedor, não de administrador do negócio.
func IsDeveloperOnlyPermission(permission string) bool {
	return permission == "admin.create_permissions" ||
		strings.HasPrefix(permission, "roles.") ||
		strings.HasPrefix(permission, "permissions.")
}
