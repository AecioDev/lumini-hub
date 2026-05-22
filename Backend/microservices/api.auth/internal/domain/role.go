package domain

import (
	"gorm.io/gorm"
)

// Role representa um perfil de usuário
type Role struct {
	gorm.Model

	Name        string       `gorm:"size:50;not null;unique" json:"name"`
	Description string       `json:"description"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
	Users       []User       `gorm:"foreignKey:RoleID" json:"-"`
}

func (Role) TableName() string {
	return "roles"
}

type RolePermissions struct {
	role_id       uint
	permission_id uint
}

func (RolePermissions) TableName() string {
	return "role_permissions"
}

// CreateRoleRequest representa os dados para criar um novo perfil
type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// UpdateRoleRequest representa os dados para atualizar um perfil
type UpdateRoleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// UpdateRolePermissionsRequest representa os dados para atualizar permissões de um perfil
type UpdateRolePermissionsRequest struct {
	PermissionIDs []uint `json:"permission_ids" binding:"required"`
}

// ApiRole representa os dados básicos de um Papel
type ApiRole struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ApiRoleDetail representa os dados detalhados de um papel, incluindo suas permissões
type ApiRoleDetail struct {
	ID          uint            `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Permissions []ApiPermission `json:"permissions"`
	CreatedAt   string          `json:"created_at"`
	UpdatedAt   string          `json:"updated_at"`
}

// ApiRoleFromModel converte um modelo Role para ApiRole
func ApiRoleFromModel(r Role) ApiRole {
	return ApiRole{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
	}
}

// ApiRoleDetailFromModel converte um modelo Role para ApiRoleDetail
func ApiRoleDetailFromModel(r Role) ApiRoleDetail {
	permissionDTOs := make([]ApiPermission, 0, len(r.Permissions))
	for _, perm := range r.Permissions {
		permissionDTOs = append(permissionDTOs, ApiPermissionFromModel(perm))
	}

	return ApiRoleDetail{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Permissions: permissionDTOs,
		CreatedAt:   r.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   r.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
