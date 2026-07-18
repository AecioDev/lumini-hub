package repository

import "gorm.io/gorm"

// SequenceRepository replica, no Go, o padrão GetSequencia do ERP legado: consulta o número
// atual da sequência em CADSEQ com lock, incrementa e persiste, retornando o próximo número.
type SequenceRepository interface {
	Next(codSeq string) (int64, error)
}

// GormSequenceRepository implementa SequenceRepository usando GORM sobre o SQL Server
type GormSequenceRepository struct {
	db *gorm.DB
}

// NewSequenceRepository cria um novo repository de sequências
func NewSequenceRepository(db *gorm.DB) SequenceRepository {
	return &GormSequenceRepository{db: db}
}

// Next retorna o próximo número da sequência identificada por codSeq (ex.: "NOT"), criando o
// registro em CADSEQ com valor inicial 1 caso ainda não exista. Toda a operação roda dentro de
// uma transação com WITH (UPDLOCK) para evitar que duas requisições concorrentes obtenham o
// mesmo número, replicando o comportamento do GetSequencia original.
func (r *GormSequenceRepository) Next(codSeq string) (int64, error) {
	var next int64

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var current int64
		result := tx.Raw("SELECT numseq FROM CADSEQ WITH (UPDLOCK) WHERE codseq = UPPER(?)", codSeq).Scan(&current)
		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			// Nenhuma linha encontrada: cria a sequência com valor inicial 1
			next = 1
			return tx.Exec("INSERT INTO CADSEQ (codseq, desseq, numseq) VALUES (UPPER(?), '', 1)", codSeq).Error
		}

		next = current + 1
		return tx.Exec("UPDATE CADSEQ SET numseq = ? WHERE codseq = UPPER(?)", next, codSeq).Error
	})

	if err != nil {
		return 0, err
	}
	return next, nil
}
