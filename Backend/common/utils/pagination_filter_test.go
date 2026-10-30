package utils

import "testing"

func TestPaginationFromFilter(t *testing.T) {
	two, fifty := 2, 50
	f, tr := false, true

	casos := []struct {
		nome             string
		pageNo, pageSize *int
		col              string
		isAsc            *bool
		page, limit      int
		sort, order      string
	}{
		{"padrões", nil, nil, "", nil, 1, 10, "id", "desc"},
		{"página e tamanho", &two, &fifty, "", nil, 2, 50, "id", "desc"},
		{"valores inválidos caem no padrão", new(int), new(int), "", nil, 1, 10, "id", "desc"},
		{"coluna sem is_asc é crescente", nil, nil, "name", nil, 1, 10, "name", "asc"},
		{"coluna decrescente", nil, nil, "name", &f, 1, 10, "name", "desc"},
		{"coluna crescente", nil, nil, "name", &tr, 1, 10, "name", "asc"},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			got := PaginationFromFilter(c.pageNo, c.pageSize, c.col, c.isAsc, "id")
			if got.Page != c.page || got.Limit != c.limit || got.Sort != c.sort || got.Order != c.order {
				t.Fatalf("got %+v", got)
			}
		})
	}
}
