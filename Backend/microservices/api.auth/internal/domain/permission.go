package domain

import (
	"lumini-hub/common/utils"

	"gorm.io/gorm"
)

// Permission representa uma permissão no sistema
type Permission struct {
	gorm.Model

	Permission  string `gorm:"size:100;not null;unique" json:"permission"`
	Description string `json:"description"`
	Module      string `gorm:"size:50;not null" json:"module"`

	// Relacionamento many-to-many com Roles
	Roles []Role `gorm:"many2many:role_permissions;" json:"-"`
}

// TableName especifica o nome da tabela
func (Permission) TableName() string {
	return "permissions"
}

// PermissionsByModule agrupa permissões por módulo
type PermissionsByModule struct {
	Module      string       `json:"module"`
	Permissions []Permission `json:"permissions"`
}

// DTOs de Entrada

// InCreatePermission representa os dados para criar uma permissão
type InCreatePermission struct {
	Description string `json:"description" binding:"required,min=3,max=500"`
	Permission  string `json:"permission" binding:"required,min=3,max=100"`
	Module      string `json:"module" binding:"required,min=3,max=50"`
}

// InUpdatePermission representa os dados para atualizar uma permissão
type InUpdatePermission struct {
	Description *string `json:"description" binding:"required,min=3,max=500"`
	Permission  *string `json:"permission" binding:"required,min=3,max=100"`
	Module      *string `json:"module" binding:"required,min=3,max=50"`
}

// InGetPermissionsFilters representa os parâmetros de filtro para buscar permissões
type InGetPermissionsFilters struct {
	Name           string `form:"name"`
	Module         string `form:"module"`
	RoleId         uint   `form:"roleId"`
	IsLinkedToRole *bool  `form:"isLinkedToRole"`
}

// DTOs de Saída

// ApiPermission representa os dados de uma permissão que são seguros para enviar ao frontend
type ApiPermission struct {
	ID          uint   `json:"id"`
	Permission  string `json:"permission"`
	Description string `json:"description"`
	Module      string `json:"module"`
}

// ApiPermissionDetail representa os dados detalhados de uma permissão
type ApiPermissionDetail struct {
	ID          uint   `json:"id"`
	Permission  string `json:"permission"`
	Description string `json:"description"`
	Module      string `json:"module"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// ApiPermissionsByModule representa permissões agrupadas por módulo
type ApiPermissionsByModule struct {
	Module      string          `json:"module"`
	Permissions []ApiPermission `json:"permissions"`
}

// ApiPermissionListPaginated representa uma lista paginada de permissões
type ApiPermissionListPaginated struct {
	Permission []ApiPermission     `json:"data"`
	Pagination utils.ApiPagination `json:"pagination"`
}

// Conversores

// ApiPermissionFromModel converte um modelo Permission para ApiPermission
func ApiPermissionFromModel(p Permission) ApiPermission {
	return ApiPermission{
		ID:          p.ID,
		Permission:  p.Permission,
		Description: p.Description,
		Module:      p.Module,
	}
}

// ApiPermissionByModuleFromModel converte um modelo PermissionsByModule para ApiPermissionsByModule
func ApiPermissionByModuleFromModel(p PermissionsByModule) ApiPermissionsByModule {
	permissionDTOs := make([]ApiPermission, 0, len(p.Permissions))
	for _, perm := range p.Permissions {
		permissionDTOs = append(permissionDTOs, ApiPermissionFromModel(perm))
	}

	return ApiPermissionsByModule{
		Module:      p.Module,
		Permissions: permissionDTOs,
	}
}

// ApiPermissionDetailFromModel converte um modelo Permission para ApiPermissionDetail
func ApiPermissionDetailFromModel(p Permission) ApiPermissionDetail {
	return ApiPermissionDetail{
		ID:          p.ID,
		Permission:  p.Permission,
		Description: p.Description,
		Module:      p.Module,
		CreatedAt:   p.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   p.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
