package database

import (
	"context"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"gorm.io/gorm/logger"
)

// defaultRecorderParamsFilter guarda o filtro original do GORM (no-op), pra restaurá-lo em development.
var defaultRecorderParamsFilter = logger.RecorderParamsFilter

// newLogger devolve o logger do GORM conforme o ambiente (APP_ENV).
//
// Só em "development" ele loga TODA consulta (nível Info), com os valores dos parâmetros: útil
// pra depurar localmente, mas o SQL de um login ou de um cadastro carrega password_hash e dados
// pessoais, então isso não pode ir pra servidor. Em qualquer outro ambiente (staging, production
// e também APP_ENV vazio ou desconhecido) o nível é Warn: só consultas lentas e erros, sem os
// valores dos parâmetros e sem poluir o log com "record not found", que o código usa como fluxo
// normal (ex.: os seeds checam se o registro já existe).
//
// Os valores são retirados do SQL por dois mecanismos, e os dois são necessários:
//   - ParameterizedQueries, no logger: cobre Exec, First, Find, Create etc.;
//   - logger.RecorderParamsFilter, uma variável GLOBAL do GORM: o Scan() grava o SQL por um
//     "gravador" interno que ignora o filtro do logger, e sem isto o valor vazava no log de
//     consulta lenta ou de erro de um Scan (ver log_test.go). Por ser global, vale pro processo
//     inteiro: cada serviço abre a conexão uma vez, na subida.
//
// `out` nil escreve em os.Stdout, como o logger padrão do GORM.
func newLogger(appEnv string, out io.Writer) logger.Interface {
	if out == nil {
		out = os.Stdout
	}

	if strings.EqualFold(strings.TrimSpace(appEnv), "development") {
		logger.RecorderParamsFilter = defaultRecorderParamsFilter
		return logger.New(log.New(out, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold: 200 * time.Millisecond,
			LogLevel:      logger.Info,
			Colorful:      true,
		})
	}

	logger.RecorderParamsFilter = func(ctx context.Context, sql string, params ...interface{}) (string, []interface{}) {
		return sql, nil
	}
	return logger.New(log.New(out, "\r\n", log.LstdFlags), logger.Config{
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  logger.Warn,
		IgnoreRecordNotFoundError: true,
		ParameterizedQueries:      true,
		Colorful:                  false,
	})
}
