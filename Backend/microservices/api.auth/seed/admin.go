package seed

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"lumini-hub/api.auth/internal/domain"
	"lumini-hub/common/utils"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	// AdminName é o nome de exibição do usuário admin inicial.
	AdminName = "Administrador"
	// MinAdminPasswordLength é o tamanho mínimo da senha do admin inicial. Maior que o
	// mínimo da API de usuários (6) de propósito: este usuário nasce com perfil ADMIN,
	// num servidor que pode estar exposto na internet.
	MinAdminPasswordLength = 8
	// defaultEmailDomain completa o e-mail do admin quando BOOTSTRAP_ADMIN_EMAIL não vem:
	// users.email é unique, então "" duplicaria assim que outro usuário sem e-mail existisse.
	defaultEmailDomain = "@lumini.local"
)

// AdminBootstrap reúne os dados do primeiro usuário ADMIN. Transitório: no multi-tenant
// quem provisiona o primeiro admin é o app comercial.
type AdminBootstrap struct {
	Username string
	Password string
	Email    string
}

// String e GoString omitem a senha, pra ela não vazar se alguém imprimir o valor (%v, %+v, %#v).
func (a AdminBootstrap) String() string {
	return fmt.Sprintf("{Username:%s Email:%s Password:<omitida>}", a.Username, a.Email)
}

func (a AdminBootstrap) GoString() string { return a.String() }

// BootstrapAdmin cria o primeiro usuário ADMIN e devolve se criou. Idempotente e conservador:
//
//   - Só cria se a tabela users estiver VAZIA (contando os excluídos logicamente): reiniciar
//     o serviço não duplica nem altera a senha, e apagar o único usuário não o recria.
//   - Sem senha (BOOTSTRAP_ADMIN_PASSWORD) não cria nada e apenas avisa. Não há senha default.
//   - Senha com menos de MinAdminPasswordLength caracteres é erro de configuração.
//   - A senha nunca é logada nem entra em mensagem de erro.
//   - O usuário nasce com perfil ADMIN, ativo e sem company_id (master: vê todas as empresas).
//
// Precisa rodar depois de SyncCatalog, que cria a role ADMIN no primeiro boot.
func BootstrapAdmin(db *gorm.DB, in AdminBootstrap) (bool, error) {
	var created bool

	err := db.Transaction(func(tx *gorm.DB) error {
		var total int64
		if err := tx.Unscoped().Model(&domain.User{}).Count(&total).Error; err != nil {
			return fmt.Errorf("contando usuários: %w", err)
		}
		if total > 0 {
			return nil
		}

		if in.Password == "" {
			log.Printf("[api.auth] BOOTSTRAP_ADMIN_PASSWORD não definida: o admin inicial não foi criado")
			return nil
		}
		if len(in.Password) < MinAdminPasswordLength {
			return fmt.Errorf("BOOTSTRAP_ADMIN_PASSWORD curta demais: mínimo de %d caracteres", MinAdminPasswordLength)
		}
		username := strings.TrimSpace(in.Username)
		if username == "" {
			return errors.New("BOOTSTRAP_ADMIN_USERNAME vazio")
		}
		email := strings.TrimSpace(in.Email)
		if email == "" {
			email = username + defaultEmailDomain
		}

		var admin domain.Role
		if err := tx.Where("name = ?", AdminRoleName).First(&admin).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("role %s não encontrada: o admin inicial não pode ser criado", AdminRoleName)
			}
			return fmt.Errorf("buscando a role %s: %w", AdminRoleName, err)
		}

		passwordHash, err := utils.HashPassword(in.Password)
		if err != nil {
			return fmt.Errorf("gerando o hash da senha: %w", err)
		}

		user := domain.User{
			Username:     username,
			PasswordHash: passwordHash,
			Name:         AdminName,
			Email:        email,
			IsActive:     true,
			RoleID:       admin.ID,
		}
		// DoNothing: se outra instância criou o mesmo usuário no boot simultâneo, o unique
		// de username não derruba a transação.
		result := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&user)
		if result.Error != nil {
			return fmt.Errorf("criando o admin inicial: %w", result.Error)
		}
		if result.RowsAffected == 0 {
			return nil
		}
		created = true
		log.Printf("[api.auth] Admin inicial %q criado com o perfil %s", username, AdminRoleName)
		return nil
	})
	if err != nil {
		return false, err
	}
	return created, nil
}
