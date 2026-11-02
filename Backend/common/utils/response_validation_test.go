package utils

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type validationTestRequest struct {
	FirstName string   `json:"first_name" binding:"required,max=5"`
	Email     string   `json:"email" binding:"omitempty,email"`
	Age       int      `json:"age" binding:"min=18"`
	Tags      []string `json:"tags" binding:"max=2"`
	Kind      string   `json:"kind" binding:"omitempty,oneof=a b"`
	Hidden    string   `json:"-" binding:"required"`
}

// bindJSON faz o mesmo que os handlers: ShouldBindJSON numa requisição com o corpo dado.
func bindJSON(t *testing.T, body string) error {
	t.Helper()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodPost, "/x", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	var req validationTestRequest
	return c.ShouldBindJSON(&req)
}

func validationMessages(t *testing.T, err error) (messages []string, body string) {
	t.Helper()
	capturaLog(t)
	c, rec := newTestContext(http.MethodPost, "/api/x")
	ValidationErrorResponse(c, "Dados inválidos", err)
	r := decode(t, rec)
	if r.StatusCode != http.StatusBadRequest || r.Message != "Dados inválidos" {
		t.Fatalf("resposta = %+v", r)
	}
	return r.ValidationErrors, rec.Body.String()
}

func assertNoInternals(t *testing.T, body string) {
	t.Helper()
	for _, internal := range []string{"Key:", "Field validation", "validationTestRequest", "FirstName", "failed on the", "Go struct field", "json: cannot", "invalid character"} {
		if strings.Contains(body, internal) {
			t.Errorf("a resposta mostra detalhe interno %q: %s", internal, body)
		}
	}
}

func TestValidationErrorResponse_UsaONomeDoCampoNoJSONEFrasesEmPortugues(t *testing.T) {
	err := bindJSON(t, `{"first_name":"123456","email":"nao-e-email","age":3,"tags":["a","b","c"],"kind":"z"}`)
	messages, body := validationMessages(t, err)

	want := []string{
		"first_name: no máximo 5 caracteres",
		"email: e-mail inválido",
		"age: valor mínimo 18",
		"tags: no máximo 2 itens",
		"kind: valor não permitido",
	}
	got := strings.Join(messages, " | ")
	for _, w := range want {
		if !strings.Contains(got, w) {
			t.Errorf("faltou %q nas mensagens: %s", w, got)
		}
	}
	assertNoInternals(t, body)
}

func TestValidationErrorResponse_CampoObrigatorioEMinimoDeTexto(t *testing.T) {
	messages, body := validationMessages(t, bindJSON(t, `{"age":30}`))
	if len(messages) == 0 || !strings.Contains(strings.Join(messages, "|"), "first_name: campo obrigatório") {
		t.Errorf("mensagens = %v", messages)
	}
	assertNoInternals(t, body)
}

func TestValidationErrorResponse_ErrosDeJSON(t *testing.T) {
	cases := map[string]string{
		`{"first_name":"ok","age":"abc"}`: "age: tipo de valor inválido",
		`{bad json}`:                      "o corpo da requisição não é um JSON válido",
		`{"first_name":`:                  "o corpo da requisição não é um JSON válido",
		``:                                "o corpo da requisição está vazio",
	}
	for body, want := range cases {
		messages, resp := validationMessages(t, bindJSON(t, body))
		if len(messages) != 1 || !strings.Contains(messages[0], want) {
			t.Errorf("corpo %q: mensagens = %v, esperado %q", body, messages, want)
		}
		assertNoInternals(t, resp)
	}
}

// Erro de regra de negócio (escrito por nós) continua indo ao usuário como está.
func TestValidationErrorResponse_RegraDeNegocioContinuaComOTextoDela(t *testing.T) {
	messages, _ := validationMessages(t, errors.New("já existe um cliente com este documento"))
	if len(messages) != 1 || messages[0] != "já existe um cliente com este documento" {
		t.Errorf("mensagens = %v", messages)
	}
}

// Mesmo recebendo o erro em si (e não o texto), erro de banco continua escondido.
func TestValidationErrorResponse_ErroDeBancoComoErroContinuaEscondido(t *testing.T) {
	messages, body := validationMessages(t, errors.New(pgErro))
	if len(messages) != 1 || messages[0] != GenericOperationError {
		t.Errorf("mensagens = %v", messages)
	}
	if strings.Contains(body, "SQLSTATE") {
		t.Errorf("a resposta vazou o SQLSTATE: %s", body)
	}
}

// O nome do campo vem do JSON só no validador do Gin: o campo sem tag json ou com "-" usa o nome em Go.
func TestValidationErrorResponse_CampoComJSONHifenUsaONomeEmGo(t *testing.T) {
	err := binding.Validator.ValidateStruct(&validationTestRequest{FirstName: "ok", Age: 30})
	messages, _ := validationMessages(t, err)
	if len(messages) != 1 || !strings.HasPrefix(messages[0], "Hidden: ") {
		t.Errorf("mensagens = %v", messages)
	}
}

// regraDeNegocioFalsa imita os erros dos validators dos serviços (lista de campo + mensagem).
type regraDeNegocioFalsa map[string]string

func (r regraDeNegocioFalsa) Error() string { return "regra de negócio" }
func (r regraDeNegocioFalsa) FieldMessages() map[string]string {
	return r
}

func TestValidationErrorResponse_FieldErrorsApontamOCampo(t *testing.T) {
	capturaLog(t)

	// Erro de bind: o campo sai pelo nome do JSON.
	c, rec := newTestContext(http.MethodPost, "/api/x")
	ValidationErrorResponse(c, "Dados inválidos", bindJSON(t, `{"first_name":"abcdefgh","age":20,"hidden":"x"}`))
	got := decode(t, rec)
	if got.FieldErrors["first_name"] != "no máximo 5 caracteres" {
		t.Errorf("fieldErrors do bind = %v", got.FieldErrors)
	}

	// Regra de negócio que implementa FieldMessager.
	c, rec = newTestContext(http.MethodPost, "/api/x")
	ValidationErrorResponse(c, "Dados inválidos", regraDeNegocioFalsa{"email": "email já está em uso"})
	got = decode(t, rec)
	if got.FieldErrors["email"] != "email já está em uso" || len(got.FieldErrors) != 1 {
		t.Errorf("fieldErrors da regra = %v", got.FieldErrors)
	}

	// Erro sem campo: a resposta não traz o mapa.
	c, rec = newTestContext(http.MethodPost, "/api/x")
	ValidationErrorResponse(c, "Dados inválidos", errors.New("algo deu errado"))
	if got = decode(t, rec); got.FieldErrors != nil {
		t.Errorf("erro sem campo não deveria ter fieldErrors: %v", got.FieldErrors)
	}
}
