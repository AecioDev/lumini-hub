package database

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"lumini-hub/common/testutil"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const secretHash = "$2a$10$SEGREDO-NAO-PODE-VAZAR"

func traceSQL(t *testing.T, appEnv string) string {
	t.Helper()
	var buf bytes.Buffer
	l := newLogger(appEnv, &buf)
	l.Trace(context.Background(), time.Now(), func() (string, int64) {
		return `UPDATE "users" SET "password_hash"='` + secretHash + `' WHERE id = 1`, 1
	}, nil)
	return buf.String()
}

// Fora de development nenhuma consulta bem-sucedida vai pro log.
func TestNewLogger_ForaDeDevelopmentNaoLogaConsultas(t *testing.T) {
	for _, env := range []string{"production", "staging", "test", "", "  ", "Production", "desconhecido"} {
		if out := traceSQL(t, env); out != "" {
			t.Errorf("APP_ENV=%q logou uma consulta bem-sucedida: %s", env, out)
		}
	}
}

func TestNewLogger_DevelopmentContinuaLogandoTudo(t *testing.T) {
	for _, env := range []string{"development", "Development", " development "} {
		out := traceSQL(t, env)
		if !strings.Contains(out, "UPDATE") {
			t.Errorf("APP_ENV=%q deveria logar as consultas (nível Info), saiu: %q", env, out)
		}
	}
}

func TestNewLogger_ParametrosSoVaoParaOLogEmDevelopment(t *testing.T) {
	type paramsFilter interface {
		ParamsFilter(ctx context.Context, sql string, params ...interface{}) (string, []interface{})
	}
	for env, wantParams := range map[string]bool{"development": true, "staging": false, "production": false, "": false} {
		pf, ok := newLogger(env, &bytes.Buffer{}).(paramsFilter)
		if !ok {
			t.Fatalf("APP_ENV=%q: o logger não implementa ParamsFilter, o GORM não filtraria os valores", env)
		}
		_, params := pf.ParamsFilter(context.Background(), "UPDATE users SET password_hash = ?", secretHash)
		if (len(params) > 0) != wantParams {
			t.Errorf("APP_ENV=%q: parâmetros presentes=%v, esperado %v", env, len(params) > 0, wantParams)
		}
	}
}

// Em development os valores continuam visíveis, inclusive no Scan: o filtro global do gravador
// (que fora de development esconde os valores) é restaurado ao padrão.
func TestNewLogger_DevelopmentMostraOsValoresMesmoNoScan(t *testing.T) {
	db := testutil.OpenTestDB(t)
	t.Cleanup(func() { logger.RecorderParamsFilter = defaultRecorderParamsFilter })

	// primeiro um ambiente "fechado", pra provar que a volta a development restaura o filtro
	_ = newLogger("production", &bytes.Buffer{})

	var buf bytes.Buffer
	dev := db.Session(&gorm.Session{Logger: newLogger("development", &buf)})
	var out struct{ Value string }
	if err := dev.Raw("SELECT ? AS value", "valor-visivel-em-dev").Scan(&out).Error; err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "valor-visivel-em-dev") {
		t.Errorf("em development o Scan deveria logar o valor do parâmetro, saiu: %q", buf.String())
	}
}

// Prova ponta a ponta com banco real: nem a consulta que dá ERRO (que o nível Warn loga) pode
// carregar o valor do parâmetro. Precisa de TEST_DATABASE_URL; pula sem ela.
func TestNewLogger_ErroEmProducaoNaoVazaValorDoParametro(t *testing.T) {
	db := testutil.OpenTestDB(t)

	var buf bytes.Buffer
	quiet := db.Session(&gorm.Session{Logger: newLogger("production", &buf)})
	if err := quiet.Exec("INSERT INTO tabela_que_nao_existe (password_hash) VALUES (?)", secretHash).Error; err == nil {
		t.Fatal("a consulta deveria falhar (tabela inexistente)")
	}
	if buf.Len() == 0 {
		t.Error("o erro da consulta deveria ser logado (nível Warn)")
	}
	if strings.Contains(buf.String(), secretHash) || strings.Contains(buf.String(), "SEGREDO") {
		t.Errorf("o log de erro vazou o valor do parâmetro: %s", buf.String())
	}

	// consulta LENTA que dá certo: o nível Warn a loga (acima de 200 ms), e também sem o valor.
	// Cobre o outro caminho em que um SQL chega ao log em produção (o do erro está acima).
	buf.Reset()
	var slow struct{ Value string }
	if err := quiet.Raw("SELECT pg_sleep(0.3), ? AS value", secretHash).Scan(&slow).Error; err != nil {
		t.Fatalf("a consulta lenta deveria funcionar: %v", err)
	}
	if !strings.Contains(buf.String(), "SLOW SQL") {
		t.Errorf("a consulta lenta deveria ser logada como SLOW SQL (nível Warn), saiu: %q", buf.String())
	}
	if strings.Contains(buf.String(), secretHash) || strings.Contains(buf.String(), "SEGREDO") {
		t.Errorf("o log de consulta lenta vazou o valor do parâmetro: %s", buf.String())
	}

	// Scan que dá ERRO: o Scan grava o SQL por um gravador interno do GORM que ignora o filtro do
	// logger (por isso a consulta lenta acima vazava antes de newLogger sobrescrever o filtro global).
	buf.Reset()
	var scanned struct{ Value string }
	if err := quiet.Raw("SELECT ? AS value FROM tabela_que_nao_existe", secretHash).Scan(&scanned).Error; err == nil {
		t.Fatal("o Scan deveria falhar (tabela inexistente)")
	}
	if strings.Contains(buf.String(), secretHash) || strings.Contains(buf.String(), "SEGREDO") {
		t.Errorf("o log de erro de um Scan vazou o valor do parâmetro: %s", buf.String())
	}

	// e o registro-não-encontrado, que o código usa como fluxo normal, não vira linha de erro
	buf.Reset()
	var row struct{ ID int }
	_ = quiet.Raw("SELECT 1 AS id WHERE false").First(&row).Error
	if buf.Len() != 0 {
		t.Errorf("record not found não deveria ser logado em produção: %s", buf.String())
	}
}
