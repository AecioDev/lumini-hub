package domain

// Cencus mapeia a tabela de empresas do ERP legado (FOCCO_ERP). Usada para configurar
// a empresa vinculada à Loja Integrada (codemp).
type Cencus struct {
	CodCencus int    `gorm:"column:codcencus" json:"cod_cencus"`
	DesCencus string `gorm:"column:descencus" json:"des_cencus"`
	CgcCencus string `gorm:"column:cgccencus" json:"cgc_cencus"`
}

// TableName especifica o nome da tabela no SQL Server
func (Cencus) TableName() string {
	return "CENCUS"
}
