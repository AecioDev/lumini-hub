package utils

import (
	"errors"
	"slices"

	"gorm.io/gorm"
)

// ResolveVisibleCompanyIDs resolve o conjunto de Company IDs que o usuário
// enxerga, conforme a regra de visibilidade fechada em 2026-07-21
// (Documentos/Planejamento/Modulo_0_Configuracao/plano_empresa.md):
//
//   - Usuário sem Company vinculada ("master") vê todas as empresas, a
//     partir da Matriz — sinalizado pelo retorno unrestricted=true, pra não
//     obrigar o chamador a enumerar toda a tabela só pra "não filtrar nada".
//   - Usuário com Company vinculada vê só a própria, a menos que tenha a
//     permission companies.hierarchy.view — nesse caso vê também as
//     subsidiárias abaixo dela na hierarquia self-referencing.
//
// Lê as tabelas `users`/`companies` direto (mesmo banco físico
// compartilhado nesta fase de migração, mesmo padrão de UserHasPermission)
// — qualquer microsserviço com uma *gorm.DB pra esse banco pode chamar,
// sem round-trip HTTP entre serviços. Pensada pra ser reaproveitada por
// qualquer entidade futura com escopo "hard-scoped" por Company (ver Step 0
// da skill lumini_hub_entity_creation): o chamador aplica
// `WHERE company_id IN (...)` na própria tabela usando o retorno daqui.
func ResolveVisibleCompanyIDs(db *gorm.DB, userID uint) (companyIDs []uint, unrestricted bool, err error) {
	var row struct {
		CompanyID *uint `gorm:"column:company_id"`
	}
	if err := db.Table("users").Select("company_id").Where("id = ?", userID).Take(&row).Error; err != nil {
		return nil, false, err
	}

	if row.CompanyID == nil {
		return nil, true, nil
	}

	hasHierarchyView, err := UserHasPermission(db, userID, "companies.hierarchy.view")
	if err != nil {
		return nil, false, err
	}
	if !hasHierarchyView {
		return []uint{*row.CompanyID}, false, nil
	}

	var ids []uint
	err = db.Raw(`
		WITH RECURSIVE subtree AS (
			SELECT id FROM companies WHERE id = ? AND deleted_at IS NULL
			UNION ALL
			SELECT c.id FROM companies c
			JOIN subtree s ON c.parent_id = s.id
			WHERE c.deleted_at IS NULL
		)
		SELECT id FROM subtree
	`, *row.CompanyID).Scan(&ids).Error
	if err != nil {
		return nil, false, err
	}
	return ids, false, nil
}

// companyExistsAndActive checa se uma Company existe, não está soft-deleted
// e está ativa — usado só pelo mecanismo de "empresa ativa" abaixo (operar
// "como" uma empresa desativada não faz sentido, diferente de só listá-la).
func companyExistsAndActive(db *gorm.DB, companyID uint) (bool, error) {
	var count int64
	err := db.Table("companies").
		Where("id = ? AND deleted_at IS NULL AND is_active = true", companyID).
		Count(&count).Error
	return count > 0, err
}

// ResolveActiveCompany decide qual Company está "ativa" na sessão do
// usuário e se ele precisa escolher uma antes de prosseguir — mecanismo de
// "empresa ativa" do usuário master, fechado em plano_empresa.md (Opção B:
// guardado no backend, fora do JWT; hoje na própria linha do usuário,
// users.active_company_id, revalidado a cada leitura, nunca confiando
// cegamente no valor guardado).
//
// Se o usuário só enxerga UMA empresa (não-master sem companies.hierarchy.view),
// essa é sempre a ativa — não há escolha real, então nunca requiresSelection.
// Se enxerga mais de uma (master, ou tem hierarchy.view), usa o valor
// gravado em storedActiveCompanyID SE ainda for uma opção válida; senão,
// requiresSelection=true (o chamador deve pedir pro usuário escolher, ver
// SetActiveCompany).
func ResolveActiveCompany(db *gorm.DB, userID uint, storedActiveCompanyID *uint) (activeCompanyID *uint, requiresSelection bool, err error) {
	visibleIDs, unrestricted, err := ResolveVisibleCompanyIDs(db, userID)
	if err != nil {
		return nil, false, err
	}

	if !unrestricted && len(visibleIDs) == 1 {
		id := visibleIDs[0]
		return &id, false, nil
	}

	// Master (unrestricted) só tem escolha de verdade se existir mais de uma
	// Company ativa no tenant inteiro — com só uma, não faz sentido pedir
	// pra "escolher" entre uma opção só.
	if unrestricted {
		var allIDs []uint
		if err := db.Table("companies").Where("deleted_at IS NULL AND is_active = true").
			Pluck("id", &allIDs).Error; err != nil {
			return nil, false, err
		}
		if len(allIDs) == 1 {
			return &allIDs[0], false, nil
		}
	}

	if storedActiveCompanyID != nil {
		if unrestricted {
			ok, err := companyExistsAndActive(db, *storedActiveCompanyID)
			if err != nil {
				return nil, false, err
			}
			if ok {
				return storedActiveCompanyID, false, nil
			}
		} else if slices.Contains(visibleIDs, *storedActiveCompanyID) {
			return storedActiveCompanyID, false, nil
		}
	}

	return nil, true, nil
}

// SetActiveCompany grava a empresa ativa do usuário, validando antes que ela
// está entre as que ele pode enxergar (ResolveVisibleCompanyIDs) e que
// existe/está ativa — nunca aceita cegamente o ID vindo do cliente.
// ErrForbidden se a empresa existe mas não é visível pro usuário;
// ErrNotFound se a empresa não existe (ou está inativa/excluída).
func SetActiveCompany(db *gorm.DB, userID uint, companyID uint) error {
	visibleIDs, unrestricted, err := ResolveVisibleCompanyIDs(db, userID)
	if err != nil {
		return err
	}

	if unrestricted {
		ok, err := companyExistsAndActive(db, companyID)
		if err != nil {
			return err
		}
		if !ok {
			return ErrNotFound
		}
	} else if !slices.Contains(visibleIDs, companyID) {
		return ErrForbidden
	}

	return db.Table("users").Where("id = ?", userID).Update("active_company_id", companyID).Error
}

// CompanyOption é uma versão mínima e "sem permission" de uma Company — só
// o suficiente pra rotular um switcher/dropdown (id + nome de exibição),
// nunca o registro completo que `GET /companies` expõe (esse sim é gated
// por companies.view, porque é o cadastro administrativo de Empresas).
// Saber "em qual empresa eu estou" e trocar entre as que já enxergo é uma
// questão de identidade da própria sessão, não de administrar o cadastro —
// por isso não pode depender de companies.view (a maioria dos perfis
// operacionais, ex. Financeiro/Vendas, nunca vai ter essa permission).
type CompanyOption struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// ResolveCompanyOptions resolve {id, nome} das Companies em ids (ou, se
// unrestricted, de todas as Companies ativas do tenant) — usado por
// AuthService pra embutir a lista do switcher direto em ApiUserDetail
// (login/refresh/me), sem round-trip pelo endpoint gated de api.core.
func ResolveCompanyOptions(db *gorm.DB, ids []uint, unrestricted bool) ([]CompanyOption, error) {
	if !unrestricted && len(ids) == 0 {
		return []CompanyOption{}, nil
	}

	query := db.Table("companies").
		Select("id, COALESCE(NULLIF(trade_name, ''), legal_name) AS name").
		Where("deleted_at IS NULL AND is_active = true").
		Order("name")
	if !unrestricted {
		query = query.Where("id IN ?", ids)
	}

	options := []CompanyOption{}
	if err := query.Scan(&options).Error; err != nil {
		return nil, err
	}
	return options, nil
}

// CompanyVisualConfigOption é a identidade visual (logo + paleta) da
// empresa ativa, resolvida sem depender de companies.visual_config.view —
// aplicar o tema pra quem já está "dentro" daquela empresa é identidade de
// sessão (mesmo raciocínio de CompanyOption acima), não administração da
// configuração em si (isso continua exigindo a permission, via
// GET /company-visual-configs/by-company/:id, tela de Configurações).
type CompanyVisualConfigOption struct {
	ID             uint    `json:"id"`
	HasLogo        bool    `json:"has_logo"`
	PrimaryColor   string  `json:"primary_color"`
	SecondaryColor *string `json:"secondary_color"`
	AccentColor    *string `json:"accent_color"`
	TextColor      *string `json:"text_color"`
}

// ResolveActiveCompanyVisualConfig busca a identidade visual (se existir)
// da Company informada, pra embutir em ApiUserDetail junto da empresa
// ativa. nil, nil quando a empresa não tem CompanyVisualConfig cadastrada
// — fallback pra paleta padrão da Lumini Hub é responsabilidade do
// consumidor/frontend (ver plano_empresa.md § CompanyVisualConfig). Não
// seleciona logo_file (bytea) — só se tem logo ou não, o binário em si
// continua servido por GET /company-visual-configs/:id/logo.
func ResolveActiveCompanyVisualConfig(db *gorm.DB, companyID uint) (*CompanyVisualConfigOption, error) {
	var row struct {
		ID             uint
		HasLogo        bool `gorm:"column:has_logo"`
		PrimaryColor   string
		SecondaryColor *string
		AccentColor    *string
		TextColor      *string
	}
	err := db.Table("company_visual_configs").
		Select("id, (length(logo_file) > 0) AS has_logo, primary_color, secondary_color, accent_color, text_color").
		Where("company_id = ? AND deleted_at IS NULL", companyID).
		Take(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &CompanyVisualConfigOption{
		ID:             row.ID,
		HasLogo:        row.HasLogo,
		PrimaryColor:   row.PrimaryColor,
		SecondaryColor: row.SecondaryColor,
		AccentColor:    row.AccentColor,
		TextColor:      row.TextColor,
	}, nil
}
