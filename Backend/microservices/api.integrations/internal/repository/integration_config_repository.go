package repository

import (
	"errors"

	"lumini-hub/api.integrations/internal/domain"
	commonrepo "lumini-hub/common/repository"

	"gorm.io/gorm"
)

// IntegrationConfigRepository define as operações de acesso a dados para configurações de integração
type IntegrationConfigRepository interface {
	commonrepo.Repository[domain.IntegrationConfig]
	FindByKey(key string) (*domain.IntegrationConfig, error)
	Upsert(key string, value string) error
}

// GormIntegrationConfigRepository implementa IntegrationConfigRepository usando GORM
type GormIntegrationConfigRepository struct {
	*commonrepo.GormRepository[domain.IntegrationConfig]
}

// NewIntegrationConfigRepository cria um novo repository de configurações de integração
func NewIntegrationConfigRepository(db *gorm.DB) IntegrationConfigRepository {
	return &GormIntegrationConfigRepository{
		GormRepository: commonrepo.NewGormRepository[domain.IntegrationConfig](db),
	}
}

// FindByKey busca uma configuração pela chave
func (r *GormIntegrationConfigRepository) FindByKey(key string) (*domain.IntegrationConfig, error) {
	var config domain.IntegrationConfig
	err := r.GetDB().Where("key = ?", key).First(&config).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &config, nil
}

// Upsert cria ou atualiza o valor de uma configuração pela chave
func (r *GormIntegrationConfigRepository) Upsert(key string, value string) error {
	existing, err := r.FindByKey(key)
	if err != nil {
		return err
	}
	if existing == nil {
		return r.Create(&domain.IntegrationConfig{Key: key, Value: value})
	}
	existing.Value = value
	return r.Update(existing)
}
