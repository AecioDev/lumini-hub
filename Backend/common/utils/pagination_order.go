package utils

import (
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// safeOrderBy monta o ORDER BY de uma listagem paginada sem confiar no texto recebido.
//
// `sort` e `order` chegam da requisição (order_by_column no corpo do /filter, `sort` e `order`
// na query string), e antes iam direto para o SQL: era uma injeção de SQL na ordenação. Aqui:
//   - a coluna só vale se for uma coluna real do modelo e não estiver escondida do JSON
//     (json:"-", como o hash da senha), senão ordenar por ela vazaria a ordem do dado;
//   - o sentido só vale se for "asc" ou "desc";
//   - o resultado é um clause.OrderByColumn, que o GORM escreve com aspas e com o nome da
//     tabela, nunca texto solto;
//   - qualquer outra coisa cai na ordenação padrão: chave primária, decrescente.
//
// `sort` pode trazer o sentido junto ("first_name desc"), como os services de busca montam.
func safeOrderBy(value interface{}, db *gorm.DB, sort, order string) (clause.OrderByColumn, bool) {
	stmt := &gorm.Statement{DB: db}
	if err := stmt.Parse(value); err != nil || stmt.Schema == nil {
		return clause.OrderByColumn{}, false
	}

	fallback := clause.OrderByColumn{}
	hasFallback := false
	if pk := stmt.Schema.PrioritizedPrimaryField; pk != nil {
		fallback = clause.OrderByColumn{Column: clause.Column{Table: clause.CurrentTable, Name: pk.DBName}, Desc: true}
		hasFallback = true
	}

	parts := strings.Fields(sort)
	if len(parts) == 0 || len(parts) > 2 {
		return fallback, hasFallback
	}

	field, ok := stmt.Schema.FieldsByDBName[strings.ToLower(parts[0])]
	if !ok || field.Tag.Get("json") == "-" {
		return fallback, hasFallback
	}

	direction := strings.ToLower(strings.TrimSpace(order))
	if len(parts) == 2 {
		direction = strings.ToLower(parts[1])
	}
	switch direction {
	case "", "asc", "desc":
	default:
		return fallback, hasFallback
	}

	return clause.OrderByColumn{
		Column: clause.Column{Table: clause.CurrentTable, Name: field.DBName},
		Desc:   direction == "desc",
	}, true
}
