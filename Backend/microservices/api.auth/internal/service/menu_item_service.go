package service

import (
	"lumini-hub/api.auth/internal/domain"
	"lumini-hub/api.auth/internal/repository"
	"lumini-hub/api.auth/internal/validator"
	"lumini-hub/common/utils"
)

// MenuItemService gerencia operações relacionadas a itens do menu lateral
type MenuItemService struct {
	menuItemRepo repository.MenuItemRepository
	validator    *validator.MenuItemValidator
}

// NewMenuItemService cria um novo serviço de itens de menu
func NewMenuItemService(menuItemRepo repository.MenuItemRepository, permRepo repository.PermissionRepository) *MenuItemService {
	return &MenuItemService{
		menuItemRepo: menuItemRepo,
		validator:    validator.NewMenuItemValidator(menuItemRepo, permRepo),
	}
}

// GetMenuTree retorna a árvore completa de itens de menu (sem filtro de permissão) — uso administrativo
func (s *MenuItemService) GetMenuTree() ([]domain.ApiMenuItem, error) {
	flat, err := s.menuItemRepo.FindAllFlat()
	if err != nil {
		return nil, err
	}
	return domain.BuildMenuTree(flat), nil
}

// GetMenuItemByID busca um item de menu pelo ID
func (s *MenuItemService) GetMenuItemByID(id uint) (*domain.ApiMenuItem, error) {
	item, err := s.menuItemRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, utils.ErrNotFound
	}
	dto := domain.ApiMenuItemFromModel(*item)
	return &dto, nil
}

// CreateMenuItem cria um novo item de menu
func (s *MenuItemService) CreateMenuItem(req domain.CreateMenuItemRequest) (*domain.ApiMenuItem, error) {
	if err := s.validator.ValidateForCreation(req); err != nil {
		return nil, err
	}

	item := domain.MenuItem{
		Name:         req.Name,
		Icon:         req.Icon,
		Href:         req.Href,
		ParentID:     req.ParentID,
		PermissionID: req.PermissionID,
		Position:     req.Position,
		IsActive:     true,
	}

	if err := s.menuItemRepo.Create(&item); err != nil {
		return nil, err
	}

	dto := domain.ApiMenuItemFromModel(item)
	return &dto, nil
}

// UpdateMenuItem atualiza um item de menu existente
func (s *MenuItemService) UpdateMenuItem(id uint, req domain.UpdateMenuItemRequest) (*domain.ApiMenuItem, error) {
	if err := s.validator.ValidateForUpdate(id, req); err != nil {
		return nil, err
	}

	item, err := s.menuItemRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, utils.ErrNotFound
	}

	item.Name = req.Name
	item.Icon = req.Icon
	item.Href = req.Href
	item.ParentID = req.ParentID
	item.PermissionID = req.PermissionID
	item.Position = req.Position
	if req.IsActive != nil {
		item.IsActive = *req.IsActive
	}

	if err := s.menuItemRepo.Update(item); err != nil {
		return nil, err
	}

	updated, err := s.menuItemRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	dto := domain.ApiMenuItemFromModel(*updated)
	return &dto, nil
}

// DeleteMenuItem exclui um item de menu
func (s *MenuItemService) DeleteMenuItem(id uint) error {
	if err := s.validator.ValidateForDeletion(id); err != nil {
		return err
	}
	return s.menuItemRepo.Delete(id)
}
