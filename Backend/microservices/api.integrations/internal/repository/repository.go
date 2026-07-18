package repository

import "gorm.io/gorm"

// UnitOfWork gerencia a persistência atômica e transações das tabelas próprias do api.integrations
type UnitOfWork interface {
	IntegrationConfigs() IntegrationConfigRepository
	SyncLogs() SyncLogRepository
	ProductMappings() ProductMappingRepository
	WebhookEvents() WebhookEventRepository
	Execute(fn func(uow UnitOfWork) error) error
	GetDB() *gorm.DB
}

// GormUnitOfWork implementa UnitOfWork usando GORM
type GormUnitOfWork struct {
	db *gorm.DB
}

// NewUnitOfWork cria um novo Unit of Work a partir de uma conexão GORM
func NewUnitOfWork(db *gorm.DB) UnitOfWork {
	return &GormUnitOfWork{db: db}
}

func (u *GormUnitOfWork) IntegrationConfigs() IntegrationConfigRepository {
	return NewIntegrationConfigRepository(u.db)
}

func (u *GormUnitOfWork) SyncLogs() SyncLogRepository {
	return NewSyncLogRepository(u.db)
}

func (u *GormUnitOfWork) ProductMappings() ProductMappingRepository {
	return NewProductMappingRepository(u.db)
}

func (u *GormUnitOfWork) WebhookEvents() WebhookEventRepository {
	return NewWebhookEventRepository(u.db)
}

func (u *GormUnitOfWork) GetDB() *gorm.DB {
	return u.db
}

// Execute encapsula o trabalho em uma transação física, com rollback em caso de erro
func (u *GormUnitOfWork) Execute(fn func(uow UnitOfWork) error) error {
	return u.db.Transaction(func(tx *gorm.DB) error {
		txUow := &GormUnitOfWork{db: tx}
		return fn(txUow)
	})
}
