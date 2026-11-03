package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"lumini-hub/common/config"
	"lumini-hub/common/testutil"
	"lumini-hub/common/utils"

	"github.com/gin-gonic/gin"
)

// Desativar ou excluir um usuário derruba o acesso na hora, mesmo com o token de acesso ainda dentro da validade.
// Precisa de Postgres: sem TEST_DATABASE_URL o teste é ignorado.
func TestAuthMiddleware_ContaDesativadaOuExcluidaPerdeOAcesso(t *testing.T) {
	db := testutil.OpenTestDB(t)
	if err := db.Exec(`CREATE TABLE users (id bigserial PRIMARY KEY, is_active boolean NOT NULL DEFAULT true, deleted_at timestamptz)`).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Exec(`INSERT INTO users (id) VALUES (1), (2), (3)`).Error; err != nil {
		t.Fatal(err)
	}
	old := permissionDB
	InitPermissionChecker(db)
	t.Cleanup(func() { permissionDB = old })

	cfg := &config.Config{JWT: config.JWTConfig{Secret: "segredo-de-teste", AccessTokenExp: 15 * time.Minute}}
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/x", AuthMiddleware(cfg), func(c *gin.Context) { c.Status(http.StatusOK) })

	chamar := func(userID uint) int {
		token, err := utils.GenerateAccessToken(userID, "u", 1, "ADMIN", cfg)
		if err != nil {
			t.Fatal(err)
		}
		req := httptest.NewRequest(http.MethodGet, "/x", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		return rec.Code
	}

	if got := chamar(1); got != http.StatusOK {
		t.Fatalf("usuário ativo: status %d, esperava 200", got)
	}
	// o token de antes continua dentro da validade; só a conta mudou
	db.Exec(`UPDATE users SET is_active = false WHERE id = 1`)
	if got := chamar(1); got != http.StatusUnauthorized {
		t.Errorf("usuário desativado: status %d, esperava 401", got)
	}
	db.Exec(`UPDATE users SET deleted_at = now() WHERE id = 2`)
	if got := chamar(2); got != http.StatusUnauthorized {
		t.Errorf("usuário excluído: status %d, esperava 401", got)
	}
	if got := chamar(99); got != http.StatusUnauthorized {
		t.Errorf("usuário inexistente: status %d, esperava 401", got)
	}
	db.Exec(`UPDATE users SET is_active = true WHERE id = 1`)
	if got := chamar(1); got != http.StatusOK {
		t.Errorf("usuário reativado: status %d, esperava 200", got)
	}
	if got := chamar(3); got != http.StatusOK {
		t.Errorf("outro usuário ativo: status %d, esperava 200", got)
	}
}
