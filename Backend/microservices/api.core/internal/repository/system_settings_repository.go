package repository

import (
	"lumini-hub/api.core/internal/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// SystemSettingsRepository dá acesso à linha única de parâmetros globais do sistema (ID fixo 1).
type SystemSettingsRepository interface {
	// Get devolve os parâmetros globais; se a linha ainda não existe, cria uma vazia (nenhum parâmetro ligado).
	Get() (*domain.SystemSettings, error)
	Save(settings *domain.SystemSettings) error
}

// GormSystemSettingsRepository implementa SystemSettingsRepository usando GORM
type GormSystemSettingsRepository struct {
	db *gorm.DB
}

// NewSystemSettingsRepository cria um novo repository dos parâmetros globais
func NewSystemSettingsRepository(db *gorm.DB) SystemSettingsRepository {
	return &GormSystemSettingsRepository{db: db}
}

// Get busca a linha única, criando-a vazia na primeira leitura (ON CONFLICT DO NOTHING cobre dois pedidos simultâneos)
func (r *GormSystemSettingsRepository) Get() (*domain.SystemSettings, error) {
	empty := domain.SystemSettings{ID: domain.SystemSettingsID, SMTP: domain.SMTPSettings{Port: 587, Security: "starttls"}}
	if err := r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&empty).Error; err != nil {
		return nil, err
	}
	var settings domain.SystemSettings
	if err := r.db.First(&settings, domain.SystemSettingsID).Error; err != nil {
		return nil, err
	}
	return &settings, nil
}

// Save grava a linha única inteira
func (r *GormSystemSettingsRepository) Save(settings *domain.SystemSettings) error {
	settings.ID = domain.SystemSettingsID
	return r.db.Save(settings).Error
}
