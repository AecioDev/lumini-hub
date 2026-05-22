package dto

import (
	"simple-erp-service/internal/data-structure/models"
)

// ApiUser representa os dados de um usuário que são seguros para enviar ao frontend
type ApiUser struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email,omitempty"`
	RoleID   uint   `json:"role_id"`
	Role     string `json:"role,omitempty"`
	IsActive bool   `json:"is_active"`
}

// ApiUserDetail representa os dados detalhados de um usuário
type ApiUserDetail struct {
	ID        uint          `json:"id"`
	Username  string        `json:"username"`
	Name      string        `json:"name"`
	Email     string        `json:"email,omitempty"`
	Phone     string        `json:"phone"`
	RoleID    uint          `json:"role_id"`
	Role      ApiRoleDetail `json:"role"`
	IsActive  bool          `json:"is_active"`
	LastLogin string        `json:"last_login,omitempty"`
	CreatedAt string        `json:"created_at"`
	UpdatedAt string        `json:"updated_at"`
}

// ApiUserListPaginated representa uma lista paginada de usuários
type ApiUserListPaginated struct {
	Users      []ApiUser      `json:"users"`
	Pagination *ApiPagination `json:"pagination,omitempty"`
}

// ToDTO converte um modelo User para ApiUser
func ApiUserFromModel(u models.User) ApiUser {
	dto := ApiUser{
		ID:       u.ID,
		Username: u.Username,
		Name:     u.Name,
		Email:    u.Email,
		RoleID:   u.RoleID,
		IsActive: u.IsActive,
	}

	// Adicionar o nome do perfil se estiver carregado
	if u.Role.ID != 0 {
		dto.Role = u.Role.Name
	}

	return dto
}

// ToDetailDTO converte um modelo User para UserDetailDTO
func ApiUserDetailFromModel(u models.User) ApiUserDetail {
	dto := ApiUserDetail{
		ID:        u.ID,
		Username:  u.Username,
		Name:      u.Name,
		Email:     u.Email,
		Phone:     u.Phone,
		RoleID:    u.RoleID,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: u.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	// Adicionar o último login se existir
	if u.LastLogin != nil {
		dto.LastLogin = u.LastLogin.Format("2006-01-02 15:04:05")
	}

	// Adicionar o perfil se estiver carregado
	if u.Role.ID != 0 {
		dto.Role = ApiRoleDetailFromModel(*u.Role)
	}

	return dto
}

type LoginSuccessResponse struct {
	User ApiUserDetail `json:"user"`
}

type RefreshTokenSuccessResponse struct {
	User ApiUserDetail `json:"user"`
}
