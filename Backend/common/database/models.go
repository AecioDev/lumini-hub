package database

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// SystemLog é um registro de auditoria: quem (UserID, e o username em Details) fez qual escrita
// (Action = método + modelo da rota) em qual entidade. Gravado pelo middlewares.LoggerMiddleware
// em todos os serviços e lido pelo api.auth (GET /audit-logs). Sem soft delete: auditoria não se
// apaga. O UserID não tem chave estrangeira de propósito (o registro sobrevive à exclusão do usuário).
type SystemLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	CreatedAt  time.Time `gorm:"index" json:"created_at"`
	UserID     *uint     `gorm:"index" json:"user_id"`
	Action     string    `gorm:"size:150;not null" json:"action"`
	EntityType string    `gorm:"size:50;index" json:"entity_type"`
	EntityID   string    `gorm:"size:50" json:"entity_id"`
	Details    JSONB     `gorm:"type:jsonb" json:"details"`
	IPAddress  string    `gorm:"size:45" json:"ip_address"`
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
