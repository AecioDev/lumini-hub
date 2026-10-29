package utils

import (
	"log"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

// GenericOperationError é o texto que o usuário vê quando a operação falha por um erro do banco.
// É de propósito genérico: tamanho, duplicidade e formato são validados no frontend e no backend
// antes de chegar ao banco, então um erro de banco aqui é falha de validação nossa ou tentativa de
// ataque, e quem ataca não deve aprender nada com a resposta (tipo e tamanho de coluna, nome de
// tabela ou de constraint, host e usuário da conexão).
const GenericOperationError = "Não foi possível concluir a operação."

// Todo erro do Postgres vindo do driver traz o código no fim: "ERRO: ... (SQLSTATE 22001)".
var sqlStatePattern = regexp.MustCompile(`\(SQLSTATE [0-9A-Z]{5}\)`)

// Erros que não trazem SQLSTATE mas mostram detalhes do banco ou da conexão (host, porta, usuário,
// nome do banco). "record not found" e erros de regra de negócio não entram aqui.
var databaseErrorMarkers = []string{
	"failed to connect to",
	"dial tcp",
	"connection refused",
	"driver: bad connection",
	"sql: ",
	"pgconn",
	"pgx",
}

// isDatabaseError diz se o texto parece vir do banco de dados e não de uma regra de negócio.
func isDatabaseError(text string) bool {
	if text == "" {
		return false
	}
	if sqlStatePattern.MatchString(text) {
		return true
	}
	lower := strings.ToLower(text)
	for _, marker := range databaseErrorMarkers {
		if strings.Contains(lower, marker) {
			return true
		}
	}
	return false
}

// hideDatabaseError devolve o texto que pode ir ao cliente: o próprio texto, ou a mensagem
// genérica (gravando o original no log do servidor) quando ele é um erro de banco.
func hideDatabaseError(c *gin.Context, statusCode int, message, text string) string {
	if !isDatabaseError(text) {
		return text
	}
	method, path := requestInfo(c)
	log.Printf("[erro de banco %d] %s %s: %s: %s", statusCode, method, path, message, text)
	return GenericOperationError
}

// requestInfo devolve método e caminho (sem a query string) da requisição, para os logs.
func requestInfo(c *gin.Context) (method, path string) {
	method, path = "-", "-"
	if c != nil && c.Request != nil && c.Request.URL != nil {
		method, path = c.Request.Method, c.Request.URL.Path
	}
	return method, path
}
