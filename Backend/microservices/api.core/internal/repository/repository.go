package repository

import (
	"gorm.io/gorm"
)

// UnitOfWork gerencia a persistência atômica e transações no microsserviço Core
type UnitOfWork interface {
	Customers() CustomerRepository
	Suppliers() SupplierRepository
	Execute(fn func(uow UnitOfWork) error) error
	GetDB() *gorm.DB
}

// GormUnitOfWork implementa a interface UnitOfWork usando o GORM
type GormUnitOfWork struct {
	db *gorm.DB
}

// NewUnitOfWork cria uma nova instância de Unit of Work
func NewUnitOfWork(db *gorm.DB) UnitOfWork {
	return &GormUnitOfWork{db: db}
}

// Customers fornece acesso ao repositório de clientes sob o escopo atual (banco ou transação ativa)
func (u *GormUnitOfWork) Customers() CustomerRepository {
	return NewCustomerRepository(u.db)
}

// Suppliers fornece acesso ao repositório de fornecedores sob o escopo atual (banco ou transação ativa)
func (u *GormUnitOfWork) Suppliers() SupplierRepository {
	return NewSupplierRepository(u.db)
}

// GetDB expõe a conexão ativa com o banco de dados GORM
func (u *GormUnitOfWork) GetDB() *gorm.DB {
	return u.db
}

// Execute encapsula a execução de operações em bloco sob uma transação física no banco.
// Realiza Rollback se a função retornar erro e Commit caso ocorra sucesso.
func (u *GormUnitOfWork) Execute(fn func(uow UnitOfWork) error) error {
	return u.db.Transaction(func(tx *gorm.DB) error {
		txUow := &GormUnitOfWork{db: tx}
		return fn(txUow)
	})
}
