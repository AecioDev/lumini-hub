package service

import (
	"errors"
	"testing"

	"lumini-hub/api.auth/internal/domain"
	"lumini-hub/common/config"
	"lumini-hub/common/testutil"
)

// Cobre os favoritos do menu (precisa de Postgres: sem TEST_DATABASE_URL o teste é ignorado).
func TestAuthService_FavoritosDoMenu(t *testing.T) {
	db := testutil.OpenTestDB(t)
	if err := db.AutoMigrate(&domain.Permission{}, &domain.Role{}, &domain.User{}, &domain.MenuItem{}, &domain.UserMenuFavorite{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	admin := domain.Role{Name: "ADMIN"}
	if err := db.Create(&admin).Error; err != nil {
		t.Fatal(err)
	}
	user := domain.User{Username: "ana", Name: "Ana", Email: "ana@x.com", RoleID: admin.ID, IsActive: true}
	other := domain.User{Username: "bia", Name: "Bia", Email: "bia@x.com", RoleID: admin.ID, IsActive: true}
	for _, u := range []*domain.User{&user, &other} {
		if err := db.Create(u).Error; err != nil {
			t.Fatal(err)
		}
	}
	group := domain.MenuItem{Name: "Cadastros", IsActive: true}
	if err := db.Create(&group).Error; err != nil {
		t.Fatal(err)
	}
	var screens []domain.MenuItem
	for _, name := range []string{"A", "B", "C", "D", "E", "F"} {
		m := domain.MenuItem{Name: name, Href: "/" + name, ParentID: &group.ID, IsActive: true}
		if err := db.Create(&m).Error; err != nil {
			t.Fatal(err)
		}
		screens = append(screens, m)
	}
	ids := func(from, to int) []uint {
		var out []uint
		for _, m := range screens[from:to] {
			out = append(out, m.ID)
		}
		return out
	}

	svc := NewAuthService(db, &config.Config{})
	read := func(userID uint) []uint {
		detail := domain.ApiUserDetail{ID: userID}
		loaded, err := svc.GetUserByID(userID)
		if err != nil {
			t.Fatal(err)
		}
		detail.MenuItems, _ = svc.GetMenuItemsForUser(*loaded)
		svc.FillFavorites(&detail)
		return detail.FavoriteMenuItemIDs
	}

	if got := read(user.ID); got == nil || len(got) != 0 {
		t.Fatalf("sem favoritos devia ser lista vazia, veio %v", got)
	}

	// grava na ordem enviada (a ordem não é a dos IDs)
	want := []uint{screens[2].ID, screens[0].ID, screens[1].ID}
	if _, err := svc.SetMenuFavorites(user.ID, want); err != nil {
		t.Fatalf("SetMenuFavorites: %v", err)
	}
	if got := read(user.ID); len(got) != 3 || got[0] != want[0] || got[1] != want[1] || got[2] != want[2] {
		t.Fatalf("favoritos = %v, quer %v", got, want)
	}
	if got := read(other.ID); len(got) != 0 {
		t.Fatalf("favoritos de outro usuário vazaram: %v", got)
	}

	// trocar a lista substitui, e esvaziar funciona
	if _, err := svc.SetMenuFavorites(user.ID, ids(3, 5)); err != nil {
		t.Fatal(err)
	}
	if got := read(user.ID); len(got) != 2 || got[0] != screens[3].ID {
		t.Fatalf("depois da troca = %v", got)
	}
	if _, err := svc.SetMenuFavorites(user.ID, nil); err != nil {
		t.Fatal(err)
	}
	if got := read(user.ID); len(got) != 0 {
		t.Fatalf("depois de esvaziar = %v", got)
	}

	// recusas: acima do limite, repetido, grupo sem tela, item inexistente
	invalid := map[string][]uint{
		"seis favoritos": ids(0, 6),
		"repetido":       {screens[0].ID, screens[0].ID},
		"grupo":          {group.ID},
		"inexistente":    {9999},
	}
	for name, list := range invalid {
		if _, err := svc.SetMenuFavorites(user.ID, list); !errors.Is(err, ErrInvalidFavorites) {
			t.Errorf("%s: err = %v, quer ErrInvalidFavorites", name, err)
		}
	}
	if got := read(user.ID); len(got) != 0 {
		t.Fatalf("pedido inválido alterou os favoritos: %v", got)
	}
}
