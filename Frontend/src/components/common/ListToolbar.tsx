import type { ReactNode } from "react";
import { Input } from "antd";

interface ListToolbarProps {
  searchPlaceholder: string;
  onSearch: (text: string) => void;
  /** Botão que abre a gaveta de filtros (FilterButton), quando a lista tem filtros. */
  filter?: ReactNode;
  /** Botões da direita (ex.: "Novo Usuário"). */
  actions?: ReactNode;
}

// Barra no topo do cartão de toda lista: busca por texto e botão de filtros à esquerda, ações de criar à
// direita. Quebra de linha em tela estreita. O total de registros não fica aqui, e sim no pé da lista
// (ResponsiveTable).
export function ListToolbar({ searchPlaceholder, onSearch, filter, actions }: ListToolbarProps) {
  return (
    <div style={{ display: "flex", flexWrap: "wrap", alignItems: "center", gap: 8, marginBottom: 16 }}>
      <div style={{ display: "flex", gap: 8, flex: "1 1 260px", maxWidth: 480 }}>
        <Input.Search
          placeholder={searchPlaceholder}
          allowClear
          style={{ flex: 1, minWidth: 0 }}
          onChange={(e) => onSearch(e.target.value)}
        />
        {filter}
      </div>
      {actions && <div style={{ marginLeft: "auto" }}>{actions}</div>}
    </div>
  );
}
