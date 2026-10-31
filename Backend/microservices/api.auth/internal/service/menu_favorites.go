package service

import (
	"errors"

	"lumini-hub/api.auth/internal/domain"

	"gorm.io/gorm"
)

// ErrInvalidFavorites: a lista de favoritos enviada tem item repetido, acima do limite ou que o
// usuário não enxerga no menu.
var ErrInvalidFavorites = errors.New("lista de favoritos inválida")

// destinationIDs devolve os IDs dos itens de menu que levam a uma tela (têm href) na árvore visível.
func destinationIDs(items []domain.ApiUserMenuItem, into map[uint]bool) map[uint]bool {
	for _, item := range items {
		if item.Href != "" {
			into[item.ID] = true
		}
		destinationIDs(item.Children, into)
	}
	return into
}

// FillFavorites preenche os favoritos do usuário na resposta de login/refresh/me, na ordem em que ele
// os arrumou. Linhas de itens que sumiram do menu ou que ele deixou de enxergar são ignoradas. Em caso
// de erro de leitura devolve a lista vazia: favoritos não podem impedir a sessão de abrir.
func (s *AuthService) FillFavorites(userDetail *domain.ApiUserDetail) {
	userDetail.FavoriteMenuItemIDs = []uint{}

	var rows []domain.UserMenuFavorite
	if err := s.db.Where("user_id = ?", userDetail.ID).Order("position asc, id asc").Find(&rows).Error; err != nil {
		return
	}
	visible := destinationIDs(userDetail.MenuItems, map[uint]bool{})
	for _, row := range rows {
		if visible[row.MenuItemID] {
			userDetail.FavoriteMenuItemIDs = append(userDetail.FavoriteMenuItemIDs, row.MenuItemID)
		}
	}
}

// SetMenuFavorites troca a lista de favoritos do usuário pela enviada (na ordem). Só aceita telas que
// ele realmente enxerga no menu, sem repetição e até domain.MaxMenuFavorites. A troca é atômica.
func (s *AuthService) SetMenuFavorites(userID uint, menuItemIDs []uint) ([]uint, error) {
	if len(menuItemIDs) > domain.MaxMenuFavorites {
		return nil, ErrInvalidFavorites
	}
	user, err := s.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	menu, err := s.GetMenuItemsForUser(*user)
	if err != nil {
		return nil, err
	}
	visible := destinationIDs(menu, map[uint]bool{})

	seen := make(map[uint]bool, len(menuItemIDs))
	rows := make([]domain.UserMenuFavorite, 0, len(menuItemIDs))
	for position, id := range menuItemIDs {
		if seen[id] || !visible[id] {
			return nil, ErrInvalidFavorites
		}
		seen[id] = true
		rows = append(rows, domain.UserMenuFavorite{UserID: userID, MenuItemID: id, Position: position})
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userID).Delete(&domain.UserMenuFavorite{}).Error; err != nil {
			return err
		}
		if len(rows) == 0 {
			return nil
		}
		return tx.Create(&rows).Error
	})
	if err != nil {
		return nil, err
	}
	return append([]uint{}, menuItemIDs...), nil
}
