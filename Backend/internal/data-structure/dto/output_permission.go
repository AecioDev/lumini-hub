package dto

import "simple-erp-service/internal/data-structure/models"

// ApiPermission representa os dados de uma permissão que são seguros para enviar ao frontend
type ApiPermission struct {
	ID          uint   `json:"id"`
	Permission  string `json:"permission"`
	Description string `json:"description"`
	Module      string `json:"module"`
}

// ApiPermission representa os dados de uma permissão que são seguros para enviar ao frontend
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

// CustomerListDTO representa uma lista paginada de clientes
type ApiPermissionListPaginated struct {
	Permission []ApiPermission `json:"data"`
	Pagination ApiPagination   `json:"pagination"`
}

// ApiPermissionFromModel converte um modelo Permission para ApiPermission
func ApiPermissionFromModel(p models.Permission) ApiPermission {
	return ApiPermission{
		ID:          p.ID,
		Permission:  p.Permission,
		Description: p.Description,
		Module:      p.Module,
	}
}

// ToDTO converte um modelo PermissionsByModule para PermissionsByModuleDTO
func ApiPermissionByModuleFromModel(p models.PermissionsByModule) ApiPermissionsByModule {
	permissionDTOs := make([]ApiPermission, 0, len(p.Permissions))
	for _, perm := range p.Permissions {
		permissionDTOs = append(permissionDTOs, ApiPermissionFromModel(perm))
	}

	return ApiPermissionsByModule{
		Module:      p.Module,
		Permissions: permissionDTOs,
	}
}

// ToDetailDTO converte um modelo Role para RoleDetailDTO
func ApiPermissionDetailFromModel(r models.Permission) ApiPermissionDetail {
	return ApiPermissionDetail{
		ID:          r.ID,
		Permission:  r.Permission,
		Description: r.Description,
		Module:      r.Module,
		CreatedAt:   r.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   r.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
