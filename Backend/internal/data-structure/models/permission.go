package models

import (
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
