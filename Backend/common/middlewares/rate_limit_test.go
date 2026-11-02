package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestRateLimitByIP(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/x", RateLimitByIP(2, 50*time.Millisecond), func(c *gin.Context) { c.Status(http.StatusOK) })

	call := func(ip string) int {
		req := httptest.NewRequest(http.MethodPost, "/x", nil)
		req.RemoteAddr = ip + ":1234"
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		return w.Code
	}

	if call("10.0.0.1") != 200 || call("10.0.0.1") != 200 {
		t.Fatal("as duas primeiras deveriam passar")
	}
	if got := call("10.0.0.1"); got != http.StatusTooManyRequests {
		t.Errorf("a terceira deveria ser limitada, veio %d", got)
	}
	if call("10.0.0.2") != 200 {
		t.Error("outro IP não é afetado")
	}
	time.Sleep(70 * time.Millisecond)
	if call("10.0.0.1") != 200 {
		t.Error("depois da janela volta a passar")
	}
}
