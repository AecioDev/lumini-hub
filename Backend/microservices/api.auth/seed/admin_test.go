package seed

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"lumini-hub/api.auth/internal/domain"
	"lumini-hub/common/testutil"
	"lumini-hub/common/utils"

	"gorm.io/gorm"
)

// Testes de banco: precisam de TEST_DATABASE_URL (ver testutil.OpenTestDB) e são pulados
// sem ela, exceto o de String(), que não usa banco.

const adminPassword = "SenhaForte-2026"

// newAdminTestDB devolve um banco com o catálogo e a role ADMIN já semeados (o que o boot
// faz antes do BootstrapAdmin).
func newAdminTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db := testutil.OpenTestDB(t)
	if err := db.AutoMigrate(&domain.Permission{}, &domain.Role{}, &domain.User{}); err != nil {
		t.Fatalf("migrando tabelas: %v", err)
	}
	mustSync(t, db)
	return db
}

func capturaLog(t *testing.T) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	log.SetOutput(&buf)
	t.Cleanup(func() { log.SetOutput(os.Stderr) })
	return &buf
}

func countUsers(t *testing.T, db *gorm.DB) int64 {
	t.Helper()
	var n int64
	if err := db.Unscoped().Model(&domain.User{}).Count(&n).Error; err != nil {
		t.Fatal(err)
	}
	return n
}

func input() AdminBootstrap {
	return AdminBootstrap{Username: "admin", Password: adminPassword, Email: "admin@exemplo.com"}
}

func TestAdminBootstrap_StringNaoVazaSenha(t *testing.T) {
	in := input()
	for _, verb := range []string{"%v", "%+v", "%#v", "%s"} {
		if out := fmt.Sprintf(verb, in); strings.Contains(out, adminPassword) {
			t.Errorf("%s vazou a senha: %s", verb, out)
		}
	}
}

func TestBootstrapAdmin_CriaComTabelaVazia(t *testing.T) {
	db := newAdminTestDB(t)
	logs := capturaLog(t)

	created, err := BootstrapAdmin(db, input())
	if err != nil || !created {
		t.Fatalf("created=%v err=%v, esperado criar", created, err)
	}

	var user domain.User
	if err := db.Preload("Role").Where("username = ?", "admin").First(&user).Error; err != nil {
		t.Fatal(err)
	}
	if user.Role == nil || user.Role.Name != AdminRoleName {
		t.Errorf("perfil = %+v, esperado %s", user.Role, AdminRoleName)
	}
	if !user.IsActive || user.CompanyID != nil || user.Email != "admin@exemplo.com" || user.Name != AdminName {
		t.Errorf("usuário inesperado: ativo=%v company=%v email=%q nome=%q", user.IsActive, user.CompanyID, user.Email, user.Name)
	}
	if user.PasswordHash == adminPassword || !utils.CheckPasswordHash(adminPassword, user.PasswordHash) {
		t.Error("a senha não foi guardada como hash bcrypt válido")
	}
	if strings.Contains(logs.String(), adminPassword) {
		t.Errorf("a senha apareceu no log: %s", logs.String())
	}
}

func TestBootstrapAdmin_EmailVazioGanhaPadraoNaoVazio(t *testing.T) {
	db := newAdminTestDB(t)
	in := input()
	in.Email = ""

	if created, err := BootstrapAdmin(db, in); err != nil || !created {
		t.Fatalf("created=%v err=%v", created, err)
	}
	var user domain.User
	db.Where("username = ?", "admin").First(&user)
	if user.Email != "admin"+defaultEmailDomain {
		t.Errorf("e-mail = %q, esperado %q", user.Email, "admin"+defaultEmailDomain)
	}
}

func TestBootstrapAdmin_SegundoBootNaoDuplicaNemMudaSenha(t *testing.T) {
	db := newAdminTestDB(t)
	if _, err := BootstrapAdmin(db, input()); err != nil {
		t.Fatal(err)
	}
	var first domain.User
	db.Where("username = ?", "admin").First(&first)

	other := input()
	other.Password = "OutraSenha-2027"
	created, err := BootstrapAdmin(db, other)
	if err != nil || created {
		t.Fatalf("2º boot: created=%v err=%v, esperado não criar", created, err)
	}
	var after domain.User
	db.Where("username = ?", "admin").First(&after)
	if countUsers(t, db) != 1 || after.PasswordHash != first.PasswordHash {
		t.Error("o 2º boot duplicou o usuário ou alterou a senha")
	}
}

func TestBootstrapAdmin_SemSenhaNaoCriaEApenasAvisa(t *testing.T) {
	db := newAdminTestDB(t)
	logs := capturaLog(t)
	in := input()
	in.Password = ""

	created, err := BootstrapAdmin(db, in)
	if err != nil || created {
		t.Fatalf("created=%v err=%v, esperado não criar e sem erro", created, err)
	}
	if countUsers(t, db) != 0 {
		t.Error("criou usuário sem senha")
	}
	if !strings.Contains(logs.String(), "BOOTSTRAP_ADMIN_PASSWORD") {
		t.Errorf("faltou o aviso no log: %q", logs.String())
	}
}

func TestBootstrapAdmin_JaExistindoUsuarioNaoCria(t *testing.T) {
	newUser := func(db *gorm.DB) domain.User {
		// users.role_id tem FK pra roles: o usuário precisa de um perfil válido (o ADMIN do seed)
		return domain.User{Username: "maria", PasswordHash: "x", Name: "Maria", Email: "m@x.com", RoleID: adminRole(t, db).ID}
	}

	t.Run("usuário ativo", func(t *testing.T) {
		db := newAdminTestDB(t)
		user := newUser(db)
		if err := db.Create(&user).Error; err != nil {
			t.Fatal(err)
		}
		if created, err := BootstrapAdmin(db, input()); err != nil || created {
			t.Errorf("created=%v err=%v, esperado não criar", created, err)
		}
	})

	t.Run("usuário excluído logicamente também conta", func(t *testing.T) {
		db := newAdminTestDB(t)
		user := newUser(db)
		if err := db.Create(&user).Error; err != nil {
			t.Fatal(err)
		}
		if err := db.Delete(&user).Error; err != nil {
			t.Fatal(err)
		}
		if created, err := BootstrapAdmin(db, input()); err != nil || created {
			t.Errorf("created=%v err=%v: apagar o único usuário não pode recriar o admin", created, err)
		}
	})
}

func TestBootstrapAdmin_ErrosDeConfiguracaoNaoVazamSenha(t *testing.T) {
	t.Run("senha curta", func(t *testing.T) {
		db := newAdminTestDB(t)
		in := input()
		in.Password = "xK9#z" // curta (5), e não aparece na mensagem de erro
		created, err := BootstrapAdmin(db, in)
		if err == nil || created {
			t.Fatalf("created=%v err=%v, esperado erro", created, err)
		}
		if strings.Contains(err.Error(), in.Password) {
			t.Errorf("o erro vazou a senha: %v", err)
		}
		if countUsers(t, db) != 0 {
			t.Error("criou usuário com senha curta")
		}
	})

	t.Run("sem a role ADMIN", func(t *testing.T) {
		db := testutil.OpenTestDB(t)
		if err := db.AutoMigrate(&domain.Permission{}, &domain.Role{}, &domain.User{}); err != nil {
			t.Fatal(err)
		}
		created, err := BootstrapAdmin(db, input())
		if err == nil || created {
			t.Fatalf("created=%v err=%v, esperado erro", created, err)
		}
		if strings.Contains(err.Error(), adminPassword) {
			t.Errorf("o erro vazou a senha: %v", err)
		}
	})

	t.Run("username vazio", func(t *testing.T) {
		db := newAdminTestDB(t)
		in := input()
		in.Username = "  "
		if created, err := BootstrapAdmin(db, in); err == nil || created {
			t.Errorf("created=%v err=%v, esperado erro", created, err)
		}
	})
}
