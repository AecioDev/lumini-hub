package domain

import (
	"time"

	"lumini-hub/common/utils"

	"gorm.io/gorm"
)

// User representa um usuário do sistema
type User struct {
	gorm.Model

	Username     string       `gorm:"size:50;not null;unique" json:"username"`
	PasswordHash string       `gorm:"size:255;not null;column:password_hash" json:"-"`
	Name         string       `gorm:"size:100;not null" json:"name"`
	Email        string       `gorm:"size:100;unique" json:"email"`
	Phone        string       `gorm:"size:20" json:"phone"`
	IsActive     bool         `gorm:"default:true" json:"is_active"`
	LastLogin    *time.Time   `json:"last_login"`
	RoleID       uint         `json:"role_id"`
	Role         *Role        `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	Permissions  []Permission `gorm:"many2many:user_permissions;" json:"permissions,omitempty"`

	// CompanyID referencia Company (api.core) por ID puro, sem GORM relation
	// cross-service (mesmo padrão de Customer/User descrito no CLAUDE.md).
	// Nulo = usuário "master": vê todas as Empresas a partir de qualquer
	// Matriz (ver regra de visibilidade em plano_empresa.md). Não nulo =
	// usuário só vê a própria Empresa, exceto com companies.hierarchy.view.
	CompanyID *uint `json:"company_id"`

	// ActiveCompanyID é a Company que o usuário está "operando como" no
	// momento — só é uma escolha real pra quem enxerga mais de uma Company
	// (master, ou tem companies.hierarchy.view); pra todo mundo mais é
	// sempre igual a CompanyID. Revalidado a cada leitura via
	// utils.ResolveActiveCompany — nunca confiar cegamente neste valor.
	ActiveCompanyID *uint `json:"active_company_id"`
}

// TableName especifica o nome da tabela
func (User) TableName() string {
	return "users"
}

// CreateUserRequest representa os dados para criar um novo usuário
type CreateUserRequest struct {
	Username  string `json:"username" binding:"required,min=3,max=50"`
	Password  string `json:"password" binding:"required,min=6"`
	Name      string `json:"name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Phone     string `json:"phone"`
	RoleID    uint   `json:"role_id" binding:"required"`
	CompanyID *uint  `json:"company_id"`
}

// UpdateUserRequest representa os dados para atualizar um usuário
type UpdateUserRequest struct {
	Name      string `json:"name"`
	Email     string `json:"email" binding:"omitempty,email"`
	Phone     string `json:"phone"`
	RoleID    uint   `json:"role_id"`
	IsActive  *bool  `json:"is_active"`
	CompanyID *uint  `json:"company_id"`
}

// ChangePasswordRequest representa os dados para alterar a senha
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
}

// ApiUser representa os dados de um usuário que são seguros para enviar ao frontend
type ApiUser struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Name      string `json:"name"`
	Email     string `json:"email,omitempty"`
	RoleID    uint   `json:"role_id"`
	Role      string `json:"role,omitempty"`
	IsActive  bool   `json:"is_active"`
	CompanyID *uint  `json:"company_id"`
}

// ApiUserRole representa os dados básicos do papel do usuário, sem incluir as permissões aninhadas
type ApiUserRole struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// ApiUserDetail representa os dados detalhados de um usuário
type ApiUserDetail struct {
	ID        uint        `json:"id"`
	Username  string      `json:"username"`
	Name      string      `json:"name"`
	Email     string      `json:"email,omitempty"`
	Phone     string      `json:"phone"`
	RoleID    uint        `json:"role_id"`
	Role      ApiUserRole `json:"role"`
	IsActive  bool        `json:"is_active"`
	LastLogin string      `json:"last_login,omitempty"`
	CompanyID *uint       `json:"company_id"`
	// ActiveCompanyID/ActiveCompanyName/RequiresCompanySelection/
	// VisibleCompanies só vêm com valor "de verdade" (revalidado) na
	// resposta de login/refresh/me — ver AuthService. Em outras respostas de
	// ApiUserDetail (ex.: admin olhando outro usuário) ActiveCompanyID é só
	// o valor bruto salvo, sem repetir a validação, e os outros três ficam
	// vazios. VisibleCompanies é proposital e propositalmente "sem
	// permission" (utils.ResolveCompanyOptions) — mostrar em qual empresa o
	// próprio usuário está e trocar entre as que ele já enxerga é identidade
	// de sessão, não administração do cadastro (isso continua exigindo
	// companies.view, em GET /companies).
	ActiveCompanyID          *uint                 `json:"active_company_id"`
	ActiveCompanyName        string                `json:"active_company_name,omitempty"`
	RequiresCompanySelection bool                  `json:"requires_company_selection,omitempty"`
	VisibleCompanies         []utils.CompanyOption `json:"visible_companies,omitempty"`
	// ActiveCompanyVisualConfig é a identidade visual (logo + paleta) da
	// empresa ativa — mesmo raciocínio de VisibleCompanies (sem depender de
	// companies.visual_config.view), pra CFG-6.1.2 aplicar no ConfigProvider
	// do frontend. nil quando a empresa não tem CompanyVisualConfig
	// cadastrada (fallback pra paleta padrão é responsabilidade do frontend).
	ActiveCompanyVisualConfig *utils.CompanyVisualConfigOption `json:"active_company_visual_config,omitempty"`
	Permissions              []ApiPermission       `json:"permissions,omitempty"`
	MenuItems                []ApiUserMenuItem     `json:"menu_items,omitempty"`
	CreatedAt                string                `json:"created_at"`
	UpdatedAt                string                `json:"updated_at"`
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
		ID:        u.ID,
		Username:  u.Username,
		Name:      u.Name,
		Email:     u.Email,
		RoleID:    u.RoleID,
		IsActive:  u.IsActive,
		CompanyID: u.CompanyID,
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
		ID:              u.ID,
		Username:        u.Username,
		Name:            u.Name,
		Email:           u.Email,
		Phone:           u.Phone,
		RoleID:          u.RoleID,
		IsActive:        u.IsActive,
		CompanyID:       u.CompanyID,
		ActiveCompanyID: u.ActiveCompanyID,
		CreatedAt:       u.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:       u.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	// Adicionar o último login se existir
	if u.LastLogin != nil {
		dto.LastLogin = u.LastLogin.Format("2006-01-02 15:04:05")
	}

	// Adicionar o perfil se estiver carregado
	if u.Role != nil && u.Role.ID != 0 {
		dto.Role = ApiUserRole{
			ID:          u.Role.ID,
			Name:        u.Role.Name,
			Description: u.Role.Description,
			CreatedAt:   u.Role.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   u.Role.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	// Adicionar permissões do usuário (lidas da relação direta user_permissions)
	if len(u.Permissions) > 0 {
		permissionDTOs := make([]ApiPermission, 0, len(u.Permissions))
		for _, perm := range u.Permissions {
			permissionDTOs = append(permissionDTOs, ApiPermissionFromModel(perm))
		}
		dto.Permissions = permissionDTOs
	}

	return dto
}

// UpdateUserPermissionsRequest representa os dados para atualizar permissões de um usuário
type UpdateUserPermissionsRequest struct {
	PermissionIDs []uint `json:"permission_ids" binding:"required"`
}
