// Package testutil reúne helpers só pra testes que precisam de um Postgres real.
package testutil

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"strings"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// OpenTestDB devolve um *gorm.DB isolado num schema próprio e descartável.
//
//   - Sem TEST_DATABASE_URL definida o teste é PULADO (t.Skip), pra `go test ./...`
//     não quebrar em quem não tem Postgres de teste.
//   - A URL precisa estar no formato postgres://... e o banco REAL da conexão
//     (current_database(), não o texto da URL) precisa terminar em "_test": o helper
//     recusa (t.Fatalf) qualquer outro, pra nunca rodar contra o banco de
//     desenvolvimento. Não há checagem de host: quem aponta pra um servidor remoto
//     cujo banco termina em _test assume o risco (o helper só cria e apaga schemas t_*).
//   - Cada chamada cria um schema novo (t_<hex>) e o apaga no t.Cleanup; a conexão
//     devolvida usa esse schema no search_path. Assim pacotes de teste rodando em
//     paralelo (go test ./...) não se atropelam e nada fora do schema é tocado.
//
// Exemplo: TEST_DATABASE_URL=postgres://postgres:senha@localhost:5432/lumini_test?sslmode=disable
func OpenTestDB(t testing.TB) *gorm.DB {
	t.Helper()

	raw := os.Getenv("TEST_DATABASE_URL")
	if raw == "" {
		t.Skip("TEST_DATABASE_URL não definida: teste que precisa de Postgres foi ignorado")
	}

	parsed, err := url.Parse(raw)
	if err != nil {
		t.Fatalf("TEST_DATABASE_URL inválida: %v", err)
	}
	if parsed.Scheme != "postgres" && parsed.Scheme != "postgresql" {
		t.Fatalf("TEST_DATABASE_URL recusada: use o formato URL (postgres://usuario:senha@host:porta/banco_test)")
	}
	if name := strings.TrimPrefix(parsed.Path, "/"); !strings.HasSuffix(name, "_test") {
		t.Fatalf("TEST_DATABASE_URL recusada: o banco %q não termina em \"_test\" (proteção contra rodar teste no banco de desenvolvimento)", name)
	}

	cfg := &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)}

	// Cada recurso ganha o seu cleanup assim que existe (LIFO: fecha a conexão do
	// teste, apaga o schema, fecha a conexão admin), pra um Fatalf no meio não vazar nada.
	admin, err := gorm.Open(postgres.Open(raw), cfg)
	if err != nil {
		t.Fatalf("conectando ao banco de teste: %v", err)
	}
	t.Cleanup(func() { closeDB(admin) })

	// Confere o banco real da conexão: pega parâmetros na query/variáveis de ambiente
	// que sobrescrevam o nome validado acima.
	var current string
	if err := admin.Raw("SELECT current_database()").Scan(&current).Error; err != nil {
		t.Fatalf("lendo current_database(): %v", err)
	}
	if !strings.HasSuffix(current, "_test") {
		t.Fatalf("TEST_DATABASE_URL recusada: a conexão caiu no banco %q, que não termina em \"_test\"", current)
	}

	suffix := make([]byte, 6)
	if _, err := rand.Read(suffix); err != nil {
		t.Fatalf("gerando nome do schema: %v", err)
	}
	schema := "t_" + hex.EncodeToString(suffix)
	if err := admin.Exec(fmt.Sprintf(`CREATE SCHEMA "%s"`, schema)).Error; err != nil {
		t.Fatalf("criando schema %s: %v", schema, err)
	}
	t.Cleanup(func() {
		if err := admin.Exec(fmt.Sprintf(`DROP SCHEMA "%s" CASCADE`, schema)).Error; err != nil {
			t.Logf("ATENÇÃO: não consegui apagar o schema %s do banco de teste: %v", schema, err)
		}
	})

	query := parsed.Query()
	query.Set("search_path", schema)
	parsed.RawQuery = query.Encode()
	db, err := gorm.Open(postgres.Open(parsed.String()), cfg)
	if err != nil {
		t.Fatalf("conectando ao schema %s: %v", schema, err)
	}
	t.Cleanup(func() { closeDB(db) })
	return db
}

func closeDB(db *gorm.DB) {
	if sqlDB, err := db.DB(); err == nil {
		_ = sqlDB.Close()
	}
}
