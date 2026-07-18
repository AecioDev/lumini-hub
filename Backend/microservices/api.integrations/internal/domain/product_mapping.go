package domain

import (
	"time"

	"gorm.io/gorm"
)

// ProductMapping espelha, no Postgres, o vínculo entre o produto do legado (CADPRO.codint)
// e o produto correspondente na Loja Integrada, junto com metadados de sincronização.
type ProductMapping struct {
	gorm.Model

	CodPro       string     `gorm:"size:30;unique;not null" json:"cod_pro"`
	LiProductID  string     `gorm:"size:50" json:"li_product_id"`
	LiSku        string     `gorm:"size:50" json:"li_sku"`
	LastSyncedAt *time.Time `json:"last_synced_at"`
	SyncStatus   string     `gorm:"size:20" json:"sync_status"`
}

// TableName especifica o nome da tabela
func (ProductMapping) TableName() string {
	return "product_mappings"
}
