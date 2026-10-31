package repository

import (
	"fmt"
	"testing"
	"time"

	"lumini-hub/api.auth/internal/domain"
	"lumini-hub/common/database"
	"lumini-hub/common/testutil"
	"lumini-hub/common/utils"
)

// Cobre os filtros do POST /audit-logs/filter (precisa de Postgres: sem TEST_DATABASE_URL o teste
// é ignorado).
func TestAuditLogFindAllFiltered(t *testing.T) {
	db := testutil.OpenTestDB(t)
	if err := db.AutoMigrate(&domain.Permission{}, &domain.Role{}, &domain.User{}, &database.SystemLog{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	admin := domain.Role{Name: "ADMIN"}
	develop := domain.Role{Name: "DEVELOP"}
	for _, role := range []*domain.Role{&admin, &develop} {
		if err := db.Create(role).Error; err != nil {
			t.Fatalf("role: %v", err)
		}
	}
	ana := domain.User{Username: "ana", Name: "Ana", Email: "ana@x.com", RoleID: admin.ID}
	dev := domain.User{Username: "dev", Name: "Dev", Email: "dev@x.com", RoleID: develop.ID}
	for _, u := range []*domain.User{&ana, &dev} {
		if err := db.Create(u).Error; err != nil {
			t.Fatalf("user: %v", err)
		}
	}

	base := time.Date(2026, 10, 30, 12, 0, 0, 0, time.UTC)
	logs := []database.SystemLog{
		{UserID: &ana.ID, Action: "POST /api/users", EntityType: "users", CreatedAt: base, Details: database.JSONB{"username": "ana", "status": 201}},
		{UserID: &ana.ID, Action: "PUT /api/users/:id", EntityType: "users", EntityID: "7", CreatedAt: base.Add(time.Hour), Details: database.JSONB{"username": "ana", "status": 200}},
		{UserID: &ana.ID, Action: "DELETE /api/companies/:id", EntityType: "companies", EntityID: "2", CreatedAt: base.Add(48 * time.Hour), Details: database.JSONB{"username": "ana", "status": 200}},
		{UserID: &dev.ID, Action: "PUT /api/roles/:id", EntityType: "roles", EntityID: "1", CreatedAt: base.Add(2 * time.Hour), Details: database.JSONB{"username": "dev", "status": 200}},
		{Action: "POST /api/customers", EntityType: "customers", CreatedAt: base.Add(3 * time.Hour)},
	}
	for i := range logs {
		if err := db.Create(&logs[i]).Error; err != nil {
			t.Fatalf("log: %v", err)
		}
	}

	repo := NewAuditLogRepository(db)
	from := base.Add(30 * time.Minute)
	to := base.Add(24 * time.Hour)

	casos := []struct {
		nome   string
		filtro domain.InGetAuditLogsFilters
		oculta string
		quer   []string // ações, na ordem padrão (id decrescente = ordem de gravação invertida)
	}{
		{"sem filtro esconde o DEVELOP", domain.InGetAuditLogsFilters{}, "DEVELOP",
			[]string{"POST /api/customers", "DELETE /api/companies/:id", "PUT /api/users/:id", "POST /api/users"}},
		{"DEVELOP vê tudo", domain.InGetAuditLogsFilters{}, "",
			[]string{"POST /api/customers", "PUT /api/roles/:id", "DELETE /api/companies/:id", "PUT /api/users/:id", "POST /api/users"}},
		{"entidade", domain.InGetAuditLogsFilters{EntityType: "users"}, "DEVELOP", []string{"PUT /api/users/:id", "POST /api/users"}},
		{"método", domain.InGetAuditLogsFilters{Method: "DELETE"}, "DEVELOP", []string{"DELETE /api/companies/:id"}},
		{"usuário", domain.InGetAuditLogsFilters{UserID: ana.ID, EntityType: "users"}, "DEVELOP", []string{"PUT /api/users/:id", "POST /api/users"}},
		{"texto na ação", domain.InGetAuditLogsFilters{Search: "companies"}, "DEVELOP", []string{"DELETE /api/companies/:id"}},
		{"texto no nome do usuário", domain.InGetAuditLogsFilters{Search: "ANA", EntityType: "companies"}, "DEVELOP", []string{"DELETE /api/companies/:id"}},
		{"período", domain.InGetAuditLogsFilters{DateFrom: &from, DateTo: &to}, "DEVELOP", []string{"POST /api/customers", "PUT /api/users/:id"}},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			pagination := utils.PaginationFromFilter(nil, nil, "", nil, "id")
			got, err := repo.FindAllFiltered(&pagination, c.filtro, c.oculta)
			if err != nil {
				t.Fatalf("FindAllFiltered: %v", err)
			}
			var actions []string
			for _, entry := range got {
				actions = append(actions, entry.Action)
			}
			if fmt.Sprint(actions) != fmt.Sprint(c.quer) {
				t.Errorf("ações = %v, quer %v", actions, c.quer)
			}
			if int(pagination.TotalRows) != len(c.quer) {
				t.Errorf("total = %d, quer %d", pagination.TotalRows, len(c.quer))
			}
		})
	}

	types, err := repo.FindEntityTypes()
	if err != nil {
		t.Fatalf("FindEntityTypes: %v", err)
	}
	if fmt.Sprint(types) != "[companies customers roles users]" {
		t.Errorf("tipos = %v", types)
	}
}
