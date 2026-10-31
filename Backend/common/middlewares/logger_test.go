package middlewares

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"lumini-hub/common/database"

	"github.com/gin-gonic/gin"
)

// auditar monta um roteador com a auditoria, deixa o teste registrar as rotas e devolve os registros
// que seriam gravados por uma requisição.
func auditar(t *testing.T, method, path string, rotas func(r *gin.Engine)) []database.SystemLog {
	t.Helper()
	gin.SetMode(gin.TestMode)
	var gravados []database.SystemLog
	router := gin.New()
	router.Use(auditMiddleware(router, func(entry database.SystemLog) { gravados = append(gravados, entry) }))
	rotas(router)
	router.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(method, path, strings.NewReader("{}")))
	return gravados
}

func responder(status int, data any) gin.HandlerFunc {
	return func(c *gin.Context) {
		if data == nil {
			c.JSON(status, gin.H{"success": true})
			return
		}
		c.JSON(status, gin.H{"success": true, "data": data})
	}
}

// usuarioAntes e usuarioDepois simulam o GET e o PUT de /api/users/:id.
func usuario(nome, email string, permissoes ...string) map[string]any {
	list := []any{}
	for i, p := range permissoes {
		list = append(list, map[string]any{"id": float64(i + 1), "permission": p})
	}
	return map[string]any{"id": 5, "name": nome, "email": email, "updated_at": nome, "permissions": list}
}

func TestAuditoria_AlteracaoGuardaSoOQueMudou(t *testing.T) {
	atual := usuario("Ana", "ana@x.com", "users.view")
	got := auditar(t, http.MethodPut, "/api/users/5", func(r *gin.Engine) {
		r.GET("/api/users/:id", func(c *gin.Context) { c.JSON(200, gin.H{"data": atual}) })
		r.PUT("/api/users/:id", func(c *gin.Context) {
			c.Set("userID", uint(7))
			c.Set("username", "admin")
			atual = usuario("Ana Souza", "ana@x.com", "users.view", "users.edit")
			c.JSON(200, gin.H{"data": atual})
		})
	})
	if len(got) != 1 {
		t.Fatalf("registros = %d, quer 1", len(got))
	}
	entry := got[0]
	if entry.Action != "PUT /api/users/:id" || entry.EntityType != "users" || entry.EntityID != "5" {
		t.Errorf("registro = %q %q %q", entry.Action, entry.EntityType, entry.EntityID)
	}
	if entry.UserID == nil || *entry.UserID != 7 || entry.Details["username"] != "admin" {
		t.Errorf("usuário = %v / %v", entry.UserID, entry.Details["username"])
	}
	changes, ok := entry.Details["changes"].(map[string]any)
	if !ok {
		t.Fatalf("sem changes: %v", entry.Details)
	}
	name, _ := changes["name"].(map[string]any)
	if name["de"] != "Ana" || name["para"] != "Ana Souza" {
		t.Errorf("name = %v", changes["name"])
	}
	if _, found := changes["email"]; found {
		t.Error("email não mudou e não deveria aparecer")
	}
	if _, found := changes["updated_at"]; found {
		t.Error("updated_at é ruído e não deveria aparecer")
	}
	perms, _ := changes["permissions"].(map[string]any)
	added, _ := perms["adicionados"].([]string)
	if len(added) != 1 || added[0] != "users.edit" {
		t.Errorf("permissions = %v, quer adicionados [users.edit]", changes["permissions"])
	}
}

func TestAuditoria_SubRecursoUsaOPaiComoAntes(t *testing.T) {
	atual := usuario("Ana", "ana@x.com", "users.view", "users.edit")
	got := auditar(t, http.MethodPut, "/api/users/5/permissions", func(r *gin.Engine) {
		r.GET("/api/users/:id", func(c *gin.Context) { c.JSON(200, gin.H{"data": atual}) })
		r.PUT("/api/users/:id/permissions", func(c *gin.Context) {
			c.JSON(200, gin.H{"data": usuario("Ana", "ana@x.com", "users.view")})
		})
	})
	if len(got) != 1 {
		t.Fatalf("registros = %d, quer 1", len(got))
	}
	perms, _ := got[0].Details["changes"].(map[string]any)["permissions"].(map[string]any)
	removed, _ := perms["removidos"].([]string)
	if len(removed) != 1 || removed[0] != "users.edit" {
		t.Errorf("permissions = %v, quer removidos [users.edit]", perms)
	}
}

func TestAuditoria_CriacaoGuardaODepoisSemSegredos(t *testing.T) {
	got := auditar(t, http.MethodPost, "/api/users", func(r *gin.Engine) {
		r.POST("/api/users", responder(201, map[string]any{"id": 9, "name": "Bia", "password_hash": "$2a$abc", "token": "t"}))
	})
	if len(got) != 1 {
		t.Fatalf("registros = %d, quer 1", len(got))
	}
	after, _ := got[0].Details["after"].(map[string]any)
	if after["name"] != "Bia" {
		t.Errorf("after = %v", after)
	}
	if after["password_hash"] != "***" || after["token"] != "***" {
		t.Errorf("segredos deviam sair como ***: %v", after)
	}
}

func TestAuditoria_ExclusaoGuardaOAntes(t *testing.T) {
	got := auditar(t, http.MethodDelete, "/api/users/5", func(r *gin.Engine) {
		r.GET("/api/users/:id", responder(200, usuario("Ana", "ana@x.com")))
		r.DELETE("/api/users/:id", responder(200, nil))
	})
	if len(got) != 1 {
		t.Fatalf("registros = %d, quer 1", len(got))
	}
	before, _ := got[0].Details["before"].(map[string]any)
	if before["name"] != "Ana" {
		t.Errorf("before = %v", got[0].Details)
	}
}

func TestAuditoria_SemPermissaoDeLeituraFicaSoODepois(t *testing.T) {
	got := auditar(t, http.MethodPut, "/api/users/5", func(r *gin.Engine) {
		r.GET("/api/users/:id", func(c *gin.Context) { c.AbortWithStatus(http.StatusForbidden) })
		r.PUT("/api/users/:id", responder(200, usuario("Ana", "ana@x.com")))
	})
	if len(got) != 1 {
		t.Fatalf("registros = %d, quer 1", len(got))
	}
	if _, found := got[0].Details["changes"]; found {
		t.Errorf("sem o antes não há changes: %v", got[0].Details)
	}
	if _, found := got[0].Details["after"]; !found {
		t.Errorf("deveria guardar o depois: %v", got[0].Details)
	}
}

func TestAuditoria_RotasSemConteudo(t *testing.T) {
	got := auditar(t, http.MethodPut, "/api/users/5/password", func(r *gin.Engine) {
		r.GET("/api/users/:id", responder(200, usuario("Ana", "ana@x.com")))
		r.PUT("/api/users/:id/password", responder(200, usuario("Ana", "ana@x.com")))
	})
	if len(got) != 1 {
		t.Fatalf("registros = %d, quer 1 (a troca de senha é auditada, sem conteúdo)", len(got))
	}
	for _, key := range []string{"changes", "before", "after"} {
		if _, found := got[0].Details[key]; found {
			t.Errorf("troca de senha não pode guardar %q", key)
		}
	}
}

func TestAuditoria_NaoAudita(t *testing.T) {
	casos := []struct {
		nome, method, route, path string
		status                    int
	}{
		{"leitura", http.MethodGet, "/api/users", "/api/users", http.StatusOK},
		{"falha", http.MethodPost, "/api/users", "/api/users", http.StatusBadRequest},
		{"erro interno", http.MethodDelete, "/api/users/:id", "/api/users/1", http.StatusInternalServerError},
		{"busca de lista", http.MethodPost, "/api/users/filter", "/api/users/filter", http.StatusOK},
		{"busca na auditoria", http.MethodPost, "/api/audit-logs/filter", "/api/audit-logs/filter", http.StatusOK},
		{"login", http.MethodPost, "/api/auth/login", "/api/auth/login", http.StatusOK},
		{"favoritos", http.MethodPut, "/api/auth/favorites", "/api/auth/favorites", http.StatusOK},
		{"webhook", http.MethodPost, "/api/integrations/webhooks/loja-integrada", "/api/integrations/webhooks/loja-integrada", http.StatusOK},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			got := auditar(t, c.method, c.path, func(r *gin.Engine) {
				r.Handle(c.method, c.route, responder(c.status, nil))
			})
			if len(got) != 0 {
				t.Errorf("não deveria auditar, veio %+v", got)
			}
		})
	}
}

func TestAuditoria_CortaNoTamanhoDaColuna(t *testing.T) {
	longo := strings.Repeat("a", 300)
	got := auditar(t, http.MethodPost, "/api/"+longo, func(r *gin.Engine) {
		r.POST("/api/"+longo, responder(201, nil))
	})
	if len(got) != 1 {
		t.Fatalf("registros = %d, quer 1", len(got))
	}
	if len(got[0].Action) != auditActionMax || len(got[0].EntityType) != auditEntityMax {
		t.Errorf("tamanhos = %d/%d, quer %d/%d", len(got[0].Action), len(got[0].EntityType), auditActionMax, auditEntityMax)
	}
}

func TestAuditoria_ConteudoGrandeFicaSoComAListaDeCampos(t *testing.T) {
	grande := map[string]any{"id": 1, "nome": strings.Repeat("x", auditMaxDetails)}
	got := auditar(t, http.MethodPost, "/api/things", func(r *gin.Engine) {
		r.POST("/api/things", responder(201, grande))
	})
	if len(got) != 1 {
		t.Fatalf("registros = %d, quer 1", len(got))
	}
	if _, found := got[0].Details["after"]; found {
		t.Error("conteúdo grande não deveria ser guardado")
	}
	if _, found := got[0].Details["truncated"]; !found {
		t.Errorf("deveria listar os campos em truncated: %v", got[0].Details)
	}
}
