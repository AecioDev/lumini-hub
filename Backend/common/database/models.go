package database

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

// SystemLog representa um log de auditoria do sistema
type SystemLog struct {
	gorm.Model
	UserID     *uint  `json:"user_id"`
	Action     string `gorm:"size:100;not null" json:"action"`
	EntityType string `gorm:"size:50" json:"entity_type"`
	EntityID   string `json:"entity_id"`
	Details    JSONB  `gorm:"type:jsonb" json:"details"`
	IPAddress  string `gorm:"size:45" json:"ip_address"`
}

// TableName especifica o nome da tabela
func (SystemLog) TableName() string {
	return "system_logs"
}

// JSONB define um tipo mapeado para colunas jsonb do Postgres no GORM
type JSONB map[string]interface{}

// Value converte o JSONB para escrita no banco
func (j JSONB) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan converte o valor lido do banco para a estrutura JSONB
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = make(JSONB)
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("falha ao converter valor do banco para []byte")
	}
	return json.Unmarshal(bytes, j)
}
