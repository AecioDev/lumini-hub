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
//   - Por segurança o nome do banco precisa terminar em "_test": o helper recusa
//     (t.Fatalf) qualquer outro, pra nunca rodar contra o banco de desenvolvimento.
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
	dbName := strings.TrimPrefix(parsed.Path, "/")
	if !strings.HasSuffix(dbName, "_test") {
		t.Fatalf("TEST_DATABASE_URL recusada: o banco %q não termina em \"_test\" (proteção contra rodar teste no banco de desenvolvimento)", dbName)
	}

	cfg := &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)}

	admin, err := gorm.Open(postgres.Open(raw), cfg)
	if err != nil {
		t.Fatalf("conectando ao banco de teste: %v", err)
	}

	suffix := make([]byte, 6)
	if _, err := rand.Read(suffix); err != nil {
		t.Fatalf("gerando nome do schema: %v", err)
	}
	schema := "t_" + hex.EncodeToString(suffix)
	if err := admin.Exec(fmt.Sprintf(`CREATE SCHEMA "%s"`, schema)).Error; err != nil {
		t.Fatalf("criando schema %s: %v", schema, err)
	}

	query := parsed.Query()
	query.Set("search_path", schema)
	parsed.RawQuery = query.Encode()
	db, err := gorm.Open(postgres.Open(parsed.String()), cfg)
	if err != nil {
		t.Fatalf("conectando ao schema %s: %v", schema, err)
	}

	t.Cleanup(func() {
		if sqlDB, err := db.DB(); err == nil {
			_ = sqlDB.Close()
		}
		_ = admin.Exec(fmt.Sprintf(`DROP SCHEMA "%s" CASCADE`, schema)).Error
		if sqlDB, err := admin.DB(); err == nil {
			_ = sqlDB.Close()
		}
	})
	return db
}
