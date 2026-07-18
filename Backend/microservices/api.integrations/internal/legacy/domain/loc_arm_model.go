package domain

// LocArm mapeia a tabela de locais de armazenamento do ERP legado (FOCCO_ERP).
// Usada para configurar o local oficial e o local de reserva provisória de estoque
// das vendas originadas na Loja Integrada.
type LocArm struct {
	CodLocAm          int    `gorm:"column:CodLocAm" json:"cod_loc_am"`
	DesLocAm          string `gorm:"column:DesLocAm" json:"des_loc_am"`
	VMaisLocalEmpresa int    `gorm:"column:VMais_Local_Empresa" json:"v_mais_local_empresa"`
	TipEstoque        string `gorm:"column:tipestoque" json:"tip_estoque"`
	TipeRc            string `gorm:"column:tiperc" json:"tipe_rc"`
}

// TableName especifica o nome da tabela no SQL Server
func (LocArm) TableName() string {
	return "LOCARM"
}
