package repository

import (
	"errors"

	"gorm.io/gorm"
)

// Repository define a interface genérica para persistência de dados no GORM
type Repository[T any] interface {
	GetDB() *gorm.DB
	FindAll() ([]T, error)
	FindByID(id uint) (*T, error)
	Create(entity *T) error
	Update(entity *T) error
	Delete(id uint) error
}

// GormRepository implementa a interface Repository usando GORM
type GormRepository[T any] struct {
	db *gorm.DB
}

// NewGormRepository cria um novo repositório genérico do GORM
func NewGormRepository[T any](db *gorm.DB) *GormRepository[T] {
	return &GormRepository[T]{db: db}
}

// GetDB retorna a instância ativa do banco de dados GORM
func (r *GormRepository[T]) GetDB() *gorm.DB {
	return r.db
}

// FindAll retorna todos os registros da entidade
func (r *GormRepository[T]) FindAll() ([]T, error) {
	var entities []T
	err := r.db.Find(&entities).Error
	return entities, err
}

// FindByID busca um registro pelo ID correspondente
func (r *GormRepository[T]) FindByID(id uint) (*T, error) {
	var entity T
	err := r.db.First(&entity, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}

// Create insere um novo registro no banco de dados
func (r *GormRepository[T]) Create(entity *T) error {
	return r.db.Create(entity).Error
}

// Update salva todas as alterações do registro no banco
func (r *GormRepository[T]) Update(entity *T) error {
	return r.db.Save(entity).Error
}

// Delete remove o registro correspondente pelo ID (Soft Delete do GORM se suportado pela entidade)
func (r *GormRepository[T]) Delete(id uint) error {
	var entity T
	return r.db.Delete(&entity, id).Error
}
