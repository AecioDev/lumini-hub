package dto

import "simple-erp-service/internal/data-structure/models"

// Dados básicos de um Papel
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

// ToDTO converte um modelo Role para RoleDTO
func ApiRoleFromModel(r models.Role) ApiRole {
	return ApiRole{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
	}
}

// ToDetailDTO converte um modelo Role para RoleDetailDTO
func ApiRoleDetailFromModel(r models.Role) ApiRoleDetail {
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
