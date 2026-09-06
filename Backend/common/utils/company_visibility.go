package utils

import (
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
