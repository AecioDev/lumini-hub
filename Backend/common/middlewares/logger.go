package middlewares

import (
	"log"
	"net/http"
	"strings"
	"time"

	"lumini-hub/common/database"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Tamanhos das colunas de system_logs (database.SystemLog): o que passa disso é cortado antes do
// banco, para o registro não ser recusado em silêncio.
const (
	auditActionMax    = 150
	auditEntityMax    = 50
	auditEntityIDMax  = 50
	auditUserAgentMax = 255
)

// auditSkippedPrefixes são os caminhos que não entram na auditoria: a sessão (login, refresh,
// favoritos...) é ruído e o webhook vem de fora, sem usuário.
var auditSkippedPrefixes = []string{
	"/api/auth/",
	"/api/integrations/webhooks/",
}

// LoggerMiddleware registra na tabela system_logs as requisições que alteram dados (POST, PUT,
// PATCH, DELETE) e terminaram com sucesso (status abaixo de 400). Leituras (GET), buscas
// (POST .../filter) e a sessão não são registradas.
//
// Além de quem fez o quê, guarda em details o que mudou: o objeto criado (after), o excluído
// (before) ou, numa alteração, só os campos que mudaram com o valor de antes e de depois
// (changes). O "depois" é o que a própria resposta da API devolveu. O "antes" vem de um GET interno
// no mesmo recurso, feito antes de a alteração rodar, com a sessão de quem pediu; sem GET
// equivalente ou sem permissão de leitura, a alteração fica só com o depois. Rota nova de escrita
// já é auditada sozinha, sem mudar o handler. O corpo da requisição nunca é gravado, e campos com
// cara de segredo (senha, token, hash, certificado) saem como "***".
//
// O registro é gravado em segundo plano, sem atrasar a resposta; um erro de gravação só aparece no
// log do servidor. `engine` é o roteador onde o middleware foi registrado (usado no GET do "antes").
func LoggerMiddleware(db *gorm.DB, engine *gin.Engine) gin.HandlerFunc {
	return auditMiddleware(engine, func(entry database.SystemLog) {
		go func() {
			if err := db.Create(&entry).Error; err != nil {
				log.Printf("[auditoria] não foi possível gravar %q: %v", entry.Action, err)
			}
		}()
	})
}

func auditMiddleware(engine *gin.Engine, save func(database.SystemLog)) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		method := c.Request.Method

		var before any
		var recorder *bodyCapture
		if isWriteMethod(method) && auditableRoute(c.FullPath()) && auditCapturesBody(c) {
			if method != http.MethodPost {
				before = fetchBefore(engine, c)
			}
			recorder = &bodyCapture{ResponseWriter: c.Writer}
			c.Writer = recorder
		}

		c.Next()

		var after any
		if recorder != nil {
			after = responseData(recorder.buf.Bytes())
		}
		if entry := buildAuditEntry(c, time.Since(startTime), before, after); entry != nil {
			save(*entry)
		}
	}
}

func isWriteMethod(method string) bool {
	return method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch || method == http.MethodDelete
}

// auditableRoute diz se uma escrita naquela rota (modelo, como /api/users/:id) entra na auditoria.
func auditableRoute(route string) bool {
	if route == "" {
		return false
	}
	// POST .../filter é busca de lista (usa POST só por causa do corpo), não escrita.
	if strings.HasSuffix(route, "/filter") {
		return false
	}
	for _, prefix := range auditSkippedPrefixes {
		if strings.HasPrefix(route, prefix) {
			return false
		}
	}
	return true
}

// buildAuditEntry monta o registro de auditoria da requisição que acabou de ser tratada, ou
// devolve nil quando ela não deve ser auditada. before e after são o estado do recurso antes e
// depois da escrita (nil quando não se sabe).
func buildAuditEntry(c *gin.Context, latency time.Duration, before, after any) *database.SystemLog {
	if !isWriteMethod(c.Request.Method) || c.Writer.Status() >= 400 {
		return nil
	}
	route := c.FullPath() // o modelo da rota (/api/users/:id): nunca o ID nem a query na ação
	if !auditableRoute(route) {
		return nil
	}

	entry := &database.SystemLog{
		Action:     truncate(c.Request.Method+" "+route, auditActionMax),
		EntityType: truncate(entityTypeFromRoute(route), auditEntityMax),
		EntityID:   truncate(c.Param("id"), auditEntityIDMax),
		IPAddress:  c.ClientIP(),
		Details: database.JSONB{
			"status":     c.Writer.Status(),
			"latency_ms": latency.Milliseconds(),
			"user_agent": truncate(c.Request.UserAgent(), auditUserAgentMax),
		},
	}
	if userID, ok := c.Get("userID"); ok {
		if id, ok := userID.(uint); ok {
			entry.UserID = &id
		}
	}
	// O nome fica no registro porque o usuário pode ser excluído depois.
	if username, ok := c.Get("username"); ok {
		if name, ok := username.(string); ok && name != "" {
			entry.Details["username"] = name
		}
	}
	addChanges(entry.Details, c.Request.Method, before, after)
	return entry
}

// entityTypeFromRoute devolve o primeiro segmento depois de /api/ ("/api/users/:id" -> "users").
func entityTypeFromRoute(route string) string {
	rest := strings.TrimPrefix(route, "/api/")
	if i := strings.IndexByte(rest, '/'); i >= 0 {
		rest = rest[:i]
	}
	return rest
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max]
}
