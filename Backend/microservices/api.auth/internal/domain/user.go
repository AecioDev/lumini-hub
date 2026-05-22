package domain

import (
	"time"

	"lumini-hub/common/utils"

	"gorm.io/gorm"
)

// User representa um usuário do sistema
type User struct {
	gorm.Model

	Username     string     `gorm:"size:50;not null;unique" json:"username"`
	PasswordHash string     `gorm:"size:255;not null;column:password_hash" json:"-"`
	Name         string     `gorm:"size:100;not null" json:"name"`
	Email        string     `gorm:"size:100;unique" json:"email"`
	Phone        string     `gorm:"size:20" json:"phone"`
	IsActive     bool       `gorm:"default:true" json:"is_active"`
	LastLogin    *time.Time `json:"last_login"`
	RoleID       uint       `json:"role_id"`
	Role         *Role      `gorm:"foreignKey:RoleID" json:"role,omitempty"`
}

// TableName especifica o nome da tabela
func (User) TableName() string {
	return "users"
}

// CreateUserRequest representa os dados para criar um novo usuário
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone"`
	RoleID   uint   `json:"role_id" binding:"required"`
}

// UpdateUserRequest representa os dados para atualizar um usuário
type UpdateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email" binding:"omitempty,email"`
	Phone    string `json:"phone"`
	RoleID   uint   `json:"role_id"`
	IsActive *bool  `json:"is_active"`
}

// ChangePasswordRequest representa os dados para alterar a senha
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
}

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
	Users      []ApiUser            `json:"users"`
	Pagination *utils.ApiPagination `json:"pagination,omitempty"`
}

// LoginSuccessResponse representa a resposta para um login de sucesso
type LoginSuccessResponse struct {
	User ApiUserDetail `json:"user"`
}

// RefreshTokenSuccessResponse representa a resposta para um refresh de token com sucesso
type RefreshTokenSuccessResponse struct {
	User ApiUserDetail `json:"user"`
}

// ApiUserFromModel converte um modelo User para ApiUser
func ApiUserFromModel(u User) ApiUser {
	dto := ApiUser{
		ID:       u.ID,
		Username: u.Username,
		Name:     u.Name,
		Email:    u.Email,
		RoleID:   u.RoleID,
		IsActive: u.IsActive,
	}

	// Adicionar o nome do perfil se estiver carregado
	if u.Role != nil && u.Role.ID != 0 {
		dto.Role = u.Role.Name
	}

	return dto
}

// ApiUserDetailFromModel converte um modelo User para ApiUserDetail
func ApiUserDetailFromModel(u User) ApiUserDetail {
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
	if u.Role != nil && u.Role.ID != 0 {
		dto.Role = ApiRoleDetailFromModel(*u.Role)
	}

	return dto
}
