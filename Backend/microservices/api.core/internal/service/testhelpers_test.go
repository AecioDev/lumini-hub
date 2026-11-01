package service

import (
	"errors"
	"testing"

	"lumini-hub/api.core/internal/domain"
	"lumini-hub/api.core/internal/repository"
	"lumini-hub/api.core/internal/validator"
	"lumini-hub/common/testutil"
)

// newTestUow abre um schema descartável do Postgres de teste (sem TEST_DATABASE_URL o teste é
// ignorado), migra as tabelas do api.core e cria as poucas tabelas do api.auth que a visibilidade
// de empresas consulta (users, permissions e user_permissions).
func newTestUow(t *testing.T) repository.UnitOfWork {
	t.Helper()
	db := testutil.OpenTestDB(t)
	if err := db.AutoMigrate(
		&domain.Company{},
		&domain.CompanyFiscalConfig{},
		&domain.CompanyEmailConfig{},
		&domain.SystemSettings{},
		&domain.Country{},
		&domain.State{},
		&domain.City{},
		&domain.Customer{},
		&domain.Supplier{},
		&domain.Address{},
		&domain.Document{},
		&domain.Contact{},
	); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	for _, ddl := range []string{
		`CREATE TABLE users (id bigserial PRIMARY KEY, company_id bigint)`,
		`CREATE TABLE permissions (id bigserial PRIMARY KEY, permission text NOT NULL)`,
		`CREATE TABLE user_permissions (user_id bigint NOT NULL, permission_id bigint NOT NULL)`,
	} {
		if err := db.Exec(ddl).Error; err != nil {
			t.Fatalf("ddl %q: %v", ddl, err)
		}
	}
	return repository.NewUnitOfWork(db)
}

// campoComErro diz se o erro é de validação e aponta o campo informado. Devolve a mensagem para o
// teste conferir o texto que o usuário veria.
func campoComErro(err error, campo string) (string, bool) {
	var ve validator.ValidationErrors
	if !errors.As(err, &ve) {
		return "", false
	}
	for _, e := range ve {
		if e.Field == campo {
			return e.Message, true
		}
	}
	return "", false
}

func exigeErroNoCampo(t *testing.T, err error, campo string) string {
	t.Helper()
	if err == nil {
		t.Fatalf("esperava erro de validação no campo %q, veio nil", campo)
	}
	msg, ok := campoComErro(err, campo)
	if !ok {
		t.Fatalf("esperava erro de validação no campo %q, veio: %v", campo, err)
	}
	return msg
}

func uintPtr(v uint) *uint { return &v }

func hasAnyValidationError(err error) bool { return validator.IsValidationError(err) }
