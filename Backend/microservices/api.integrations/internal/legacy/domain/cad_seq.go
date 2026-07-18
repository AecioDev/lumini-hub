package domain

// CadSeq mapeia a tabela de sequências do ERP legado (FOCCO_ERP), usada para gerar
// números sequenciais de documentos (ex.: "NOT" para notas) via o padrão GetSequencia.
type CadSeq struct {
	CodSeq string `gorm:"column:codseq" json:"cod_seq"`
	DesSeq string `gorm:"column:desseq" json:"des_seq"`
	NumSeq int64  `gorm:"column:numseq" json:"num_seq"`
}

// TableName especifica o nome da tabela no SQL Server
func (CadSeq) TableName() string {
	return "CADSEQ"
}
