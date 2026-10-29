package utils

import (
	"strings"
	"testing"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type orderTestItem struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	Secret    string `json:"-"`
	CreatedAt time.Time
}

// dryRunDB devolve um *gorm.DB que só monta o SQL, sem conectar em nenhum banco.
func dryRunDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(postgres.Open("host=localhost user=x dbname=x sslmode=disable"), &gorm.Config{
		DryRun:               true,
		DisableAutomaticPing: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	return db
}

// orderSQL roda o Paginate e devolve o SQL da listagem, para conferir o ORDER BY gerado.
func orderSQL(t *testing.T, sort, order string) (string, *Pagination) {
	t.Helper()
	db := dryRunDB(t)
	pagination := &Pagination{Page: 1, Limit: 10, Sort: sort, Order: order}
	query, err := Paginate(&orderTestItem{}, pagination, db.Model(&orderTestItem{}))
	if err != nil {
		t.Fatal(err)
	}
	// Fora do DryRun o GORM limpa o SQL do Count antes do Find; aqui é preciso fazer o mesmo.
	query.Statement.SQL.Reset()
	query.Statement.Vars = nil
	var items []orderTestItem
	sql := query.Find(&items).Statement.SQL.String()
	return sql, pagination
}

const defaultOrder = `ORDER BY "order_test_items"."id" DESC`

func TestPaginate_OrdenaPorColunaValida(t *testing.T) {
	cases := []struct {
		sort, order, want string
	}{
		{"name", "asc", `ORDER BY "order_test_items"."name"`},
		{"name", "desc", `ORDER BY "order_test_items"."name" DESC`},
		{"name desc", "", `ORDER BY "order_test_items"."name" DESC`},
		{"Name ASC", "", `ORDER BY "order_test_items"."name"`},
		{"created_at", "desc", `ORDER BY "order_test_items"."created_at" DESC`},
	}
	for _, tc := range cases {
		sql, _ := orderSQL(t, tc.sort, tc.order)
		if !strings.Contains(sql, tc.want) {
			t.Errorf("sort=%q order=%q: SQL = %s, esperado conter %s", tc.sort, tc.order, sql, tc.want)
		}
	}
}

// Nada que não seja uma coluna do modelo pode chegar ao SQL: cai na ordenação padrão.
func TestPaginate_TentativasDeInjecaoCaemNaOrdenacaoPadrao(t *testing.T) {
	attacks := []string{
		"id; DROP TABLE users",
		"name desc; --",
		"(CASE WHEN 1=1 THEN id ELSE name END)",
		"name, (SELECT password_hash FROM users LIMIT 1)",
		`name"`,
		"name'--",
		"1",
		"id desc, name",
		"name asc nulls last",
		"",
		"   ",
		"coluna_que_nao_existe",
		"order_test_items.name",
	}
	for _, sort := range attacks {
		sql, pagination := orderSQL(t, sort, "")
		if !strings.Contains(sql, defaultOrder) {
			t.Errorf("sort=%q: SQL = %s, esperado a ordenação padrão", sort, sql)
		}
		for _, bad := range []string{"DROP", "SELECT password_hash", "CASE WHEN", "--"} {
			if strings.Contains(sql, bad) {
				t.Errorf("sort=%q: o SQL contém %q: %s", sort, bad, sql)
			}
		}
		if pagination.Sort != "id" || pagination.Order != "desc" {
			t.Errorf("sort=%q: pagination ficou com %q/%q, esperado id/desc", sort, pagination.Sort, pagination.Order)
		}
	}
}

// Campo escondido do JSON (senha, hash, segredo) não pode ser usado para ordenar.
func TestPaginate_CampoEscondidoDoJSONNaoOrdena(t *testing.T) {
	sql, _ := orderSQL(t, "secret", "asc")
	if strings.Contains(sql, "secret") || !strings.Contains(sql, defaultOrder) {
		t.Errorf("ordenou por campo escondido do JSON: %s", sql)
	}
}

func TestPaginate_SentidoInvalidoCaiNaOrdenacaoPadrao(t *testing.T) {
	for _, order := range []string{"asc; DROP TABLE x", "up", "desc,id"} {
		sql, _ := orderSQL(t, "name", order)
		if !strings.Contains(sql, defaultOrder) {
			t.Errorf("order=%q: SQL = %s, esperado a ordenação padrão", order, sql)
		}
	}
}

func TestPaginate_CorrigeSortEOrderDaPaginacao(t *testing.T) {
	_, pagination := orderSQL(t, "name desc", "")
	if pagination.Sort != "name" || pagination.Order != "desc" {
		t.Errorf("pagination ficou com %q/%q, esperado name/desc", pagination.Sort, pagination.Order)
	}
}
