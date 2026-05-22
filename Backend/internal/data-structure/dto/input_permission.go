package dto

// Estrutura para receber os dados do front e fazer o Insert no banco
type InCreatePermission struct {
	Description string `json:"description"  binding:"required,min=3,max=500"`
	Permission  string `json:"permission"  binding:"required,min=3,max=100"`
	Module      string `json:"module"  binding:"required,min=3,max=50"`
}

// Estrutura para receber os dados do front e fazer o Update no banco
type InUpdatePermission struct {
	Description *string `json:"description"  binding:"required,min=3,max=500"`
	Permission  *string `json:"permission"  binding:"required,min=3,max=100"`
	Module      *string `json:"module"  binding:"required,min=3,max=50"`
}

// InGetPermissionsFilters representa os parâmetros de filtro para buscar permissões
type InGetPermissionsFilters struct {
	Name           string `form:"name"`           // Usará para o filtro por nome da permissão (ex: users.view)
	Module         string `form:"module"`         // Usará para o filtro por módulo (ex: users)
	RoleId         uint   `form:"roleId"`         // Opcional: para o filtro por Role ID (se for implementar)
	IsLinkedToRole *bool  `form:"isLinkedToRole"` // Opcional: para o filtro se está vinculado a role ou não.
}
