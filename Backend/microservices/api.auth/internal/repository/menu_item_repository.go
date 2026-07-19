package repository

import (
	"errors"

	"lumini-hub/api.auth/internal/domain"

	"gorm.io/gorm"
)

// MenuItemRepository define as operações de acesso a dados para itens de menu
type MenuItemRepository interface {
	Repository
	FindAllFlat() ([]domain.MenuItem, error)
	FindByID(id uint) (*domain.MenuItem, error)
	Create(menuItem *domain.MenuItem) error
	Update(menuItem *domain.MenuItem) error
	Delete(id uint) error
	HasChildren(id uint) (bool, error)
}

// GormMenuItemRepository implementa MenuItemRepository usando GORM
type GormMenuItemRepository struct {
	*BaseRepository
}

// NewMenuItemRepository cria um novo repository de itens de menu
func NewMenuItemRepository(db *gorm.DB) MenuItemRepository {
	return &GormMenuItemRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// FindAllFlat retorna todos os itens de menu (sem paginação), com a permissão pré-carregada,
// ordenados por posição — usado para montar a árvore em memória.
func (r *GormMenuItemRepository) FindAllFlat() ([]domain.MenuItem, error) {
	var items []domain.MenuItem
	if err := r.GetDB().Preload("Permission").Order("position").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// FindByID busca um item de menu pelo ID
func (r *GormMenuItemRepository) FindByID(id uint) (*domain.MenuItem, error) {
	var item domain.MenuItem
	if err := r.GetDB().Preload("Permission").First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

// Create cria um novo item de menu
func (r *GormMenuItemRepository) Create(menuItem *domain.MenuItem) error {
	return r.GetDB().Create(menuItem).Error
}

// Update atualiza um item de menu existente
func (r *GormMenuItemRepository) Update(menuItem *domain.MenuItem) error {
	return r.GetDB().Save(menuItem).Error
}

// Delete exclui um item de menu
func (r *GormMenuItemRepository) Delete(id uint) error {
	return r.GetDB().Delete(&domain.MenuItem{}, id).Error
}

// HasChildren verifica se o item possui filhos diretos (impede exclusão de nó com descendentes)
func (r *GormMenuItemRepository) HasChildren(id uint) (bool, error) {
	var count int64
	err := r.GetDB().Model(&domain.MenuItem{}).Where("parent_id = ?", id).Count(&count).Error
	return count > 0, err
}
