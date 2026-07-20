// Espelha Backend/microservices/api.auth/internal/domain/menu_item.go
// Árvore já filtrada pelas permissões do usuário, embutida no login/refresh/me.
export interface ApiUserMenuItem {
  id: number;
  name: string;
  icon: string; // string Iconify, ex: "ph:squares-four"
  href: string;
  children?: ApiUserMenuItem[];
}
