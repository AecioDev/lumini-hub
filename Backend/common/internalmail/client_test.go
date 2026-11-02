package internalmail

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"lumini-hub/common/config"
	"lumini-hub/common/utils"
)

func TestClient_EnviaComTokenDeServicoEPropagaOErro(t *testing.T) {
	var gotPath, gotCaller string
	var gotBody PasswordResetRequest
	cfg := &config.Config{JWT: config.JWTConfig{Secret: "segredo-de-teste"}}

	status := http.StatusOK
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		token, _ := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
		gotCaller, _ = utils.ValidateServiceToken(token, cfg)
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		w.WriteHeader(status)
		_, _ = w.Write([]byte(`{"success":false,"message":"O e-mail do sistema não está configurado."}`))
	}))
	defer server.Close()
	cfg.Services.CoreURL = server.URL

	client := NewClient("api.auth", cfg)
	req := PasswordResetRequest{To: "ana@exemplo.com", Name: "Ana", Token: "tok", ExpiresMinutes: 30}

	if err := client.SendPasswordReset(context.Background(), req); err != nil {
		t.Fatalf("envio: %v", err)
	}
	if gotPath != "/api/internal/mail/password-reset" || gotCaller != "api.auth" || gotBody.To != "ana@exemplo.com" || gotBody.Token != "tok" {
		t.Errorf("chamada inesperada: path=%q caller=%q body=%+v", gotPath, gotCaller, gotBody)
	}

	status = http.StatusUnprocessableEntity
	err := client.SendPasswordReset(context.Background(), req)
	if err == nil || !strings.Contains(err.Error(), "422") || !strings.Contains(err.Error(), "não está configurado") {
		t.Errorf("erro deveria trazer o status e a mensagem do api.core: %v", err)
	}
}
