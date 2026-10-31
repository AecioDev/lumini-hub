package domain

import "time"

// MaxMenuFavorites é o limite de itens do menu que um usuário pode marcar como favoritos
// (aparecem no cabeçalho). O frontend repete este número.
const MaxMenuFavorites = 5

// UserMenuFavorite é um item do menu marcado como favorito por um usuário. Position guarda a ordem
// de exibição no cabeçalho. Sem soft delete: desmarcar apaga a linha. Não há FK para menu_items de
// propósito: um item que sumiu ou que o usuário deixou de enxergar é ignorado na leitura
// (AuthService.FavoriteMenuItemIDs), então a linha órfã é inofensiva.
type UserMenuFavorite struct {
	ID         uint      `gorm:"primaryKey"`
	UserID     uint      `gorm:"not null;uniqueIndex:idx_user_menu_favorite"`
	MenuItemID uint      `gorm:"not null;uniqueIndex:idx_user_menu_favorite"`
	Position   int       `gorm:"not null;default:0"`
	CreatedAt  time.Time
}

// TableName especifica o nome da tabela
func (UserMenuFavorite) TableName() string {
	return "user_menu_favorites"
}

// SetMenuFavoritesRequest substitui a lista inteira de favoritos do usuário logado, na ordem enviada.
type SetMenuFavoritesRequest struct {
	MenuItemIDs []uint `json:"menu_item_ids" binding:"max=5,dive,gt=0"`
}

// SetMenuFavoritesResponse devolve a lista gravada, na ordem.
type SetMenuFavoritesResponse struct {
	FavoriteMenuItemIDs []uint `json:"favorite_menu_item_ids"`
}
