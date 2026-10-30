package repository

import (
	"testing"

	"lumini-hub/api.auth/internal/domain"
	"lumini-hub/common/testutil"
	"lumini-hub/common/utils"
)

// Cobre os filtros do POST /users/filter e do POST /permissions/filter (precisa de Postgres:
// sem TEST_DATABASE_URL o teste é ignorado).
func TestFindByFilterUsuarios(t *testing.T) {
	db := testutil.OpenTestDB(t)
	if err := db.AutoMigrate(&domain.Permission{}, &domain.Role{}, &domain.User{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	admin := domain.Role{Name: "ADMIN"}
	develop := domain.Role{Name: "DEVELOP"}
	vendas := domain.Role{Name: "Vendas"}
	for _, role := range []*domain.Role{&admin, &develop, &vendas} {
		if err := db.Create(role).Error; err != nil {
			t.Fatalf("role: %v", err)
		}
	}
	users := []domain.User{
		{Username: "ana", Name: "Ana Souza", Email: "ana@x.com", RoleID: admin.ID, IsActive: true},
		{Username: "bruno", Name: "Bruno Lima", Email: "bruno@x.com", RoleID: vendas.ID, IsActive: true},
		{Username: "carla", Name: "Carla Dias", Email: "carla@x.com", RoleID: vendas.ID},
		{Username: "dev", Name: "Dev Oculto", Email: "dev@x.com", RoleID: develop.ID, IsActive: true},
	}
	for i := range users {
		if err := db.Create(&users[i]).Error; err != nil {
			t.Fatalf("user: %v", err)
		}
	}
	// IsActive tem default:true no GORM, então o false de Carla precisa ser gravado à parte.
	if err := db.Model(&domain.User{}).Where("username = ?", "carla").Update("is_active", false).Error; err != nil {
		t.Fatalf("inativar: %v", err)
	}

	repo := NewUserRepository(db)
	boolPtr := func(b bool) *bool { return &b }
	uintPtr := func(u uint) *uint { return &u }

	casos := []struct {
		nome   string
		filtro domain.UserFilterRequest
		oculta string
		quer   []string
	}{
		{"sem filtro esconde o DEVELOP", domain.UserFilterRequest{}, "DEVELOP", []string{"carla", "bruno", "ana"}},
		{"DEVELOP vê todos", domain.UserFilterRequest{}, "", []string{"dev", "carla", "bruno", "ana"}},
		{"texto no nome", domain.UserFilterRequest{Search: "souza"}, "DEVELOP", []string{"ana"}},
		{"texto no login", domain.UserFilterRequest{Search: "BRUN"}, "DEVELOP", []string{"bruno"}},
		{"texto no e-mail", domain.UserFilterRequest{Search: "carla@"}, "DEVELOP", []string{"carla"}},
		{"perfil", domain.UserFilterRequest{RoleID: uintPtr(vendas.ID)}, "DEVELOP", []string{"carla", "bruno"}},
		{"só ativos", domain.UserFilterRequest{IsActive: boolPtr(true)}, "DEVELOP", []string{"bruno", "ana"}},
		{"só inativos", domain.UserFilterRequest{IsActive: boolPtr(false)}, "DEVELOP", []string{"carla"}},
		{"perfil e status", domain.UserFilterRequest{RoleID: uintPtr(vendas.ID), IsActive: boolPtr(true)}, "DEVELOP", []string{"bruno"}},
		{"o filtro não revela o DEVELOP", domain.UserFilterRequest{Search: "oculto"}, "DEVELOP", nil},
		{"curingas no texto não quebram", domain.UserFilterRequest{Search: "'; DROP TABLE users;--"}, "DEVELOP", nil},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			pagination := utils.PaginationFromFilter(nil, nil, "", nil, "id")
			got, err := repo.FindByFilter(c.filtro, &pagination, c.oculta)
			if err != nil {
				t.Fatalf("FindByFilter: %v", err)
			}
			var logins []string
			for _, u := range got {
				logins = append(logins, u.Username)
			}
			if len(logins) != len(c.quer) {
				t.Fatalf("logins = %v, quer %v", logins, c.quer)
			}
			for i := range logins {
				if logins[i] != c.quer[i] {
					t.Fatalf("logins = %v, quer %v", logins, c.quer)
				}
			}
			if pagination.TotalRows != int64(len(c.quer)) {
				t.Errorf("TotalRows = %d, quer %d", pagination.TotalRows, len(c.quer))
			}
		})
	}
}

func TestFindAllFilteredPermissoesBuscaCodigoEDescricao(t *testing.T) {
	db := testutil.OpenTestDB(t)
	if err := db.AutoMigrate(&domain.Permission{}, &domain.Role{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	perms := []domain.Permission{
		{Permission: "sales.view", Description: "Visualizar vendas", Module: "Vendas"},
		{Permission: "sales.create", Description: "Criar vendas", Module: "Vendas"},
		{Permission: "users.view", Description: "Visualizar usuários", Module: "Configurações"},
	}
	for i := range perms {
		if err := db.Create(&perms[i]).Error; err != nil {
			t.Fatalf("perm: %v", err)
		}
	}

	repo := NewPermissionRepository(db)
	casos := []struct {
		nome    string
		filtros domain.InGetPermissionsFilters
		quer    int
	}{
		{"pela descrição", domain.InGetPermissionsFilters{Name: "vendas"}, 2},
		{"pelo código", domain.InGetPermissionsFilters{Name: "users"}, 1},
		{"código e módulo", domain.InGetPermissionsFilters{Name: "view", Module: "Vendas"}, 1},
		{"sem texto", domain.InGetPermissionsFilters{}, 3},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			pagination := utils.PaginationFromFilter(nil, nil, "", nil, "id")
			got, err := repo.FindAllFiltered(&pagination, c.filtros, "")
			if err != nil {
				t.Fatalf("FindAllFiltered: %v", err)
			}
			if len(got) != c.quer {
				t.Fatalf("achou %d, quer %d", len(got), c.quer)
			}
		})
	}
}
