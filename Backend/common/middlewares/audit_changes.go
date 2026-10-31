package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
)

// Limites do que a auditoria guarda do "antes e depois" de uma escrita.
const (
	auditMaxBody    = 256 << 10 // maior resposta que ainda é lida (acima disso, só a ação é gravada)
	auditMaxDetails = 32 << 10  // maior "antes/depois" gravado em details (acima disso, só a lista de campos)
)

// auditNoBodySuffixes/Prefixes são as rotas cujo conteúdo não vai para a auditoria (a ação em si
// continua registrada): senha, certificado e logo, e as chaves de integração.
var (
	auditNoBodySuffixes = []string{"/password", "/certificate", "/logo"}
	auditNoBodyPrefixes = []string{"/api/integrations/settings"}
)

// Campos que mudam em toda gravação e só fazem barulho no "antes e depois".
var auditIgnoredFields = map[string]bool{
	"created_at": true, "updated_at": true, "deleted_at": true,
	"created_by": true, "updated_by": true, "created_by_id": true, "updated_by_id": true,
}

// Pedaços de nome de campo cujo valor nunca vai para a auditoria, mesmo que a API devolva.
var auditSecretFields = []string{"password", "senha", "token", "secret", "hash", "apikey", "api_key", "certificate_file"}

// bodyCapture copia o corpo da resposta enquanto ele segue para o cliente.
type bodyCapture struct {
	gin.ResponseWriter
	buf bytes.Buffer
}

func (w *bodyCapture) Write(b []byte) (int, error) {
	if w.buf.Len()+len(b) <= auditMaxBody {
		w.buf.Write(b)
	}
	return w.ResponseWriter.Write(b)
}

func (w *bodyCapture) WriteString(s string) (int, error) {
	if w.buf.Len()+len(s) <= auditMaxBody {
		w.buf.WriteString(s)
	}
	return w.ResponseWriter.WriteString(s)
}

// auditCapturesBody diz se o antes e o depois dessa rota podem ser guardados.
func auditCapturesBody(c *gin.Context) bool {
	route := c.FullPath()
	for _, suffix := range auditNoBodySuffixes {
		if strings.HasSuffix(route, suffix) {
			return false
		}
	}
	for _, prefix := range auditNoBodyPrefixes {
		if strings.HasPrefix(route, prefix) {
			return false
		}
	}
	return true
}

// fetchBefore devolve como o registro estava antes de uma alteração (PUT/PATCH/DELETE): faz um GET
// interno no mesmo caminho, com a sessão de quem pediu (então vale a mesma permissão de leitura),
// e, para sub-recursos como /users/:id/permissions, no registro pai. Sem GET equivalente, ou sem
// permissão de leitura, devolve nil e a auditoria fica só com o depois.
func fetchBefore(engine *gin.Engine, c *gin.Context) any {
	if engine == nil {
		return nil
	}
	candidates := []string{c.Request.URL.Path}
	routeSegs := strings.Split(c.FullPath(), "/")
	if n := len(routeSegs); n >= 3 && strings.HasPrefix(routeSegs[n-2], ":") {
		pathSegs := strings.Split(c.Request.URL.Path, "/")
		candidates = append(candidates, strings.Join(pathSegs[:len(pathSegs)-1], "/"))
	}
	for _, path := range candidates {
		req, err := http.NewRequestWithContext(c.Request.Context(), http.MethodGet, path, nil)
		if err != nil {
			continue
		}
		req.Header = c.Request.Header.Clone()
		req.Header.Del("Content-Length")
		req.RemoteAddr = c.Request.RemoteAddr
		rec := httptest.NewRecorder()
		engine.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			continue
		}
		if data := responseData(rec.Body.Bytes()); data != nil {
			return data
		}
	}
	return nil
}

// responseData extrai o campo "data" do envelope utils.Response.
func responseData(body []byte) any {
	var envelope struct {
		Data any `json:"data"`
	}
	if len(body) == 0 || json.Unmarshal(body, &envelope) != nil {
		return nil
	}
	return envelope.Data
}

// addChanges coloca em details o que mudou: o objeto criado (POST), o objeto excluído (DELETE) ou,
// numa alteração, só os campos que mudaram, cada um com o valor de antes ("de") e de depois ("para").
func addChanges(details map[string]any, method string, before, after any) {
	switch method {
	case http.MethodPost:
		if after != nil {
			details["after"] = scrub(after)
		}
	case http.MethodDelete:
		if before != nil {
			details["before"] = scrub(before)
		}
	default:
		beforeMap, okBefore := before.(map[string]any)
		afterMap, okAfter := after.(map[string]any)
		switch {
		case okBefore && okAfter:
			if changes := diffMaps(scrub(beforeMap).(map[string]any), scrub(afterMap).(map[string]any)); len(changes) > 0 {
				details["changes"] = changes
			}
		case after != nil:
			details["after"] = scrub(after)
		}
	}

	// Valor grande demais não cabe em details: fica só a lista do que mudou.
	if raw, err := json.Marshal(details); err == nil && len(raw) > auditMaxDetails {
		for _, key := range []string{"changes", "before", "after"} {
			if value, ok := details[key].(map[string]any); ok {
				names := make([]string, 0, len(value))
				for name := range value {
					names = append(names, name)
				}
				delete(details, key)
				details["truncated"] = names
			}
		}
	}
}

// diffMaps compara dois objetos JSON campo a campo.
func diffMaps(before, after map[string]any) map[string]any {
	changes := map[string]any{}
	names := map[string]bool{}
	for name := range before {
		names[name] = true
	}
	for name := range after {
		names[name] = true
	}
	for name := range names {
		if auditIgnoredFields[name] || reflect.DeepEqual(before[name], after[name]) {
			continue
		}
		if listChange, ok := diffLists(before[name], after[name]); ok {
			changes[name] = listChange
			continue
		}
		changes[name] = map[string]any{"de": before[name], "para": after[name]}
	}
	return changes
}

// diffLists resume a mudança numa lista de objetos com "id" (as permissões de um usuário, por
// exemplo) como o que entrou e o que saiu, em vez de repetir a lista inteira.
func diffLists(before, after any) (map[string]any, bool) {
	oldItems, okOld := labelsByID(before)
	newItems, okNew := labelsByID(after)
	if !okOld || !okNew {
		return nil, false
	}
	added, removed := []string{}, []string{}
	for id, label := range newItems {
		if _, found := oldItems[id]; !found {
			added = append(added, label)
		}
	}
	for id, label := range oldItems {
		if _, found := newItems[id]; !found {
			removed = append(removed, label)
		}
	}
	if len(added) == 0 && len(removed) == 0 {
		return nil, false
	}
	sort.Strings(added)
	sort.Strings(removed)
	return map[string]any{"adicionados": added, "removidos": removed}, true
}

// labelsByID indexa uma lista de objetos pelo "id", com um rótulo legível de cada um. Vazio e nil
// contam como lista vazia; qualquer outra forma não é uma lista de objetos com id.
func labelsByID(value any) (map[string]string, bool) {
	if value == nil {
		return map[string]string{}, true
	}
	list, ok := value.([]any)
	if !ok {
		return nil, false
	}
	items := make(map[string]string, len(list))
	for _, element := range list {
		object, ok := element.(map[string]any)
		if !ok {
			return nil, false
		}
		id, ok := object["id"]
		if !ok {
			return nil, false
		}
		label := fmt.Sprint(id)
		for _, field := range []string{"permission", "name", "description"} {
			if text, ok := object[field].(string); ok && text != "" {
				label = text
				break
			}
		}
		items[fmt.Sprint(id)] = label
	}
	return items, true
}

// scrub devolve uma cópia do valor com "***" no lugar dos campos que parecem segredo.
func scrub(value any) any {
	switch v := value.(type) {
	case map[string]any:
		out := make(map[string]any, len(v))
		for name, inner := range v {
			if isSecretField(name) {
				out[name] = "***"
				continue
			}
			out[name] = scrub(inner)
		}
		return out
	case []any:
		out := make([]any, len(v))
		for i, inner := range v {
			out[i] = scrub(inner)
		}
		return out
	default:
		return value
	}
}

func isSecretField(name string) bool {
	lower := strings.ToLower(name)
	for _, part := range auditSecretFields {
		if strings.Contains(lower, part) {
			return true
		}
	}
	return false
}
