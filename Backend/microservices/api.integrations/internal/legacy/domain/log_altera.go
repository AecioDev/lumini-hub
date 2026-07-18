package domain

import "time"

// Tipos de alteração monitorados no LogAltera para o polling de sincronismo
const (
	TipoAltProduto    = "P" // Produto
	TipoAltPrecoVenda = "V" // Preço de venda
	TipoAltCusto      = "U" // Custo
)

// LogAltera mapeia a tabela de log de alterações do ERP legado (FOCCO_ERP), usada
// como cursor de polling para identificar produtos/preços/custos alterados desde a última sincronização.
type LogAltera struct {
	TipoAlt   string    `gorm:"column:TipoAlt" json:"tipo_alt"`
	DatFimAlt time.Time `gorm:"column:DatFimAlt" json:"dat_fim_alt"`
	CodRegAlt string    `gorm:"column:codRegAlt" json:"cod_reg_alt"`
}

// TableName especifica o nome da tabela no SQL Server
func (LogAltera) TableName() string {
	return "LogAltera"
}
