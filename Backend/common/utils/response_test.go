package utils

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

const pgErro = `ERROR: duplicate key value violates unique constraint "uni_users_email" (SQLSTATE 23505) valor-secreto@x.com`

func newTestContext(method, path string) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(method, path+"?token=segredo", nil)
	return c, rec
}

func capturaLog(t *testing.T) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	log.SetOutput(&buf)
	t.Cleanup(func() { log.SetOutput(os.Stderr) })
	return &buf
}

func decode(t *testing.T, rec *httptest.ResponseRecorder) Response {
	t.Helper()
	var r Response
	if err := json.Unmarshal(rec.Body.Bytes(), &r); err != nil {
		t.Fatalf("resposta não é JSON: %v: %s", err, rec.Body.String())
	}
	return r
}

// No 500 o erro completo vai só para o log do servidor; o cliente recebe a message.
func TestErrorResponse_500NaoDevolveOErroCruAoCliente(t *testing.T) {
	logs := capturaLog(t)
	c, rec := newTestContext(http.MethodPost, "/api/users")

	ErrorResponse(c, http.StatusInternalServerError, "Erro ao criar usuário", pgErro)

	body := rec.Body.String()
	if strings.Contains(body, "uni_users_email") || strings.Contains(body, "valor-secreto") {
		t.Fatalf("o erro cru vazou na resposta: %s", body)
	}
	r := decode(t, rec)
	if r.Success || r.StatusCode != 500 || r.Message != "Erro ao criar usuário" || r.Error != "" {
		t.Errorf("resposta = %+v", r)
	}

	out := logs.String()
	for _, want := range []string{"POST", "/api/users", "Erro ao criar usuário", "uni_users_email"} {
		if !strings.Contains(out, want) {
			t.Errorf("o log do servidor deveria conter %q: %s", want, out)
		}
	}
	if strings.Contains(out, "token=segredo") {
		t.Errorf("a query string não deve ir para o log: %s", out)
	}
}

// Os demais status são mensagens de negócio e seguem indo ao cliente.
func TestErrorResponse_OutrosStatusMantemOErro(t *testing.T) {
	for _, status := range []int{http.StatusBadRequest, http.StatusNotFound, http.StatusConflict, http.StatusServiceUnavailable} {
		logs := capturaLog(t)
		c, rec := newTestContext(http.MethodGet, "/api/x")

		ErrorResponse(c, status, "msg", "texto para o usuário")

		if r := decode(t, rec); r.Error != "texto para o usuário" {
			t.Errorf("status %d: error = %q, esperado o texto original", status, r.Error)
		}
		if logs.Len() != 0 {
			t.Errorf("status %d não deveria gravar log de erro interno: %s", status, logs.String())
		}
	}
}

func TestErrorResponse_500SemContextoNaoQuebra(t *testing.T) {
	capturaLog(t)
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec) // sem c.Request

	ErrorResponse(c, http.StatusInternalServerError, "msg", "")

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d", rec.Code)
	}
}
