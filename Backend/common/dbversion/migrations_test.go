package dbversion

import (
	"strings"
	"testing"
	"testing/fstest"

	"lumini-hub/common/testutil"
)

func fakeMigrations(files map[string]string) fstest.MapFS {
	fsys := fstest.MapFS{}
	for name, sql := range files {
		fsys["migrations/"+name] = &fstest.MapFile{Data: []byte(sql)}
	}
	return fsys
}

func TestMigracoes_AplicaUmaVezEmOrdemERegistra(t *testing.T) {
	db := testutil.OpenTestDB(t)
	fsys := fakeMigrations(map[string]string{
		"0002_segunda.sql":  "INSERT INTO ordem (n) VALUES (2);",
		"0001_primeira.sql": "CREATE TABLE ordem (n int); INSERT INTO ordem (n) VALUES (1);",
	})

	pend, err := pendingMigrations(db, fsys)
	if err != nil || len(pend) != 2 {
		t.Fatalf("pendentes antes = %v, %v; esperado 2", pend, err)
	}
	applied, err := applyMigrations(db, fsys)
	if err != nil || strings.Join(applied, ",") != "0001_primeira.sql,0002_segunda.sql" {
		t.Fatalf("aplicadas = %v, %v", applied, err)
	}
	// reiniciar não repete nada
	applied, err = applyMigrations(db, fsys)
	if err != nil || len(applied) != 0 {
		t.Fatalf("segunda rodada aplicou %v, %v", applied, err)
	}
	var n int64
	db.Raw("SELECT count(*) FROM ordem").Scan(&n)
	if n != 2 {
		t.Errorf("linhas em ordem = %d, esperado 2 (cada migração roda só uma vez)", n)
	}
	if pend, _ := pendingMigrations(db, fsys); len(pend) != 0 {
		t.Errorf("pendentes depois = %v", pend)
	}
}

func TestMigracoes_FalhaDesfazTudoENaoRegistra(t *testing.T) {
	db := testutil.OpenTestDB(t)
	fsys := fakeMigrations(map[string]string{
		"0001_boa.sql":  "CREATE TABLE boa (n int);",
		"0002_ruim.sql": "CREATE TABLE meio (n int); SELECT * FROM tabela_que_nao_existe;",
	})
	applied, err := applyMigrations(db, fsys)
	if err == nil || !strings.Contains(err.Error(), "0002_ruim.sql") {
		t.Fatalf("erro = %v, esperado citar 0002_ruim.sql", err)
	}
	if len(applied) != 1 || applied[0] != "0001_boa.sql" {
		t.Errorf("aplicadas = %v, a primeira deveria ter ficado", applied)
	}
	if db.Migrator().HasTable("meio") {
		t.Error("a migração que falhou deixou a tabela 'meio' para trás: deveria ter sido desfeita")
	}
	if pend, _ := pendingMigrations(db, fsys); len(pend) != 1 || pend[0] != "0002_ruim.sql" {
		t.Errorf("pendentes = %v, esperado só a que falhou", pend)
	}
}

// As migrações de verdade precisam rodar num banco sem as tabelas dos outros serviços (o api.auth
// sobe antes do api.core num banco novo) e num banco antigo, com as constraints a derrubar.
func TestMigracoesEmbutidas_BancoNovoEBancoAntigo(t *testing.T) {
	novo := testutil.OpenTestDB(t)
	if _, err := ApplyMigrations(novo); err != nil {
		t.Fatalf("banco sem as tabelas do api.core: %v", err)
	}

	db := testutil.OpenTestDB(t)
	for _, ddl := range []string{
		`CREATE TABLE customers (id serial, document_number varchar(20), deleted_at timestamptz, CONSTRAINT uni_customers_document_number UNIQUE (document_number))`,
		`CREATE TABLE suppliers (id serial, document_number varchar(20), deleted_at timestamptz, CONSTRAINT uni_suppliers_document_number UNIQUE (document_number))`,
		`CREATE TABLE companies (id serial, cnpj varchar(20), deleted_at timestamptz)`,
		`CREATE TABLE permissions (id serial, permission text, module text, description text, deleted_at timestamptz)`,
		`INSERT INTO permissions (permission, module, description) VALUES ('roles.view', 'Develop', 'Visualizar usuários')`,
	} {
		if err := db.Exec(ddl).Error; err != nil {
			t.Fatalf("%s: %v", ddl, err)
		}
	}
	if _, err := ApplyMigrations(db); err != nil {
		t.Fatalf("banco antigo: %v", err)
	}

	// documento de um registro excluído deixa de bloquear o recadastro; entre ativos continua único
	db.Exec(`INSERT INTO customers (document_number, deleted_at) VALUES ('1', now())`)
	if err := db.Exec(`INSERT INTO customers (document_number) VALUES ('1')`).Error; err != nil {
		t.Errorf("recadastro depois da exclusão: %v", err)
	}
	if err := db.Exec(`INSERT INTO customers (document_number) VALUES ('1')`).Error; err == nil {
		t.Error("dois clientes ativos com o mesmo documento deveriam ser recusados")
	}
	db.Exec(`INSERT INTO companies (cnpj) VALUES ('11111111000111')`)
	if err := db.Exec(`INSERT INTO companies (cnpj) VALUES ('11111111000111')`).Error; err == nil {
		t.Error("dois CNPJ ativos iguais deveriam ser recusados")
	}
	var desc string
	db.Raw(`SELECT description FROM permissions WHERE permission = 'roles.view'`).Scan(&desc)
	if desc != "Visualizar perfis" {
		t.Errorf("descrição = %q, esperado o texto corrigido do CFG-041", desc)
	}
}
