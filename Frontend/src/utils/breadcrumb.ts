import type { ApiUserMenuItem } from "@/types/menu";

export interface BreadcrumbEntry {
  title: string;
  /** Rota para onde o item leva; ausente em grupos do menu (sem tela) e no item atual. */
  href?: string;
}

// Telas de cadastro e de edição não estão no menu: o caminho delas é o da listagem de origem + o nome
// da própria tela. `parent` é a rota (do menu) da listagem à qual a tela pertence.
interface FormPage {
  base: string;
  parent: string;
  create: string;
  edit: string;
}

const FORM_PAGES: FormPage[] = [
  { base: "/settings/users", parent: "/settings/users", create: "Novo Usuário", edit: "Editar Usuário" },
  { base: "/settings/roles", parent: "/settings/roles", create: "Novo Perfil", edit: "Editar Perfil" },
  { base: "/settings/companies", parent: "/settings/companies", create: "Nova Empresa", edit: "Editar Empresa" },
  // A lista de permissões vive dentro da tela "Perfis e Permissões".
  { base: "/settings/permissions", parent: "/settings/roles", create: "Nova Permissão", edit: "Editar Permissão" },
];

// Caminho do menu (do grupo raiz até o item) cuja rota é `pathname`. Se a mesma rota aparece em mais de
// um lugar da árvore (ex.: Clientes), vale a primeira.
function findMenuPath(items: ApiUserMenuItem[], pathname: string): ApiUserMenuItem[] | null {
  for (const item of items) {
    if (item.href === pathname) return [item];
    if (item.children) {
      const inside = findMenuPath(item.children, pathname);
      if (inside) return [item, ...inside];
    }
  }
  return null;
}

function toEntries(path: ApiUserMenuItem[], linkLast: boolean): BreadcrumbEntry[] {
  return path.map((item, index) => ({
    title: item.name,
    href: item.href && (linkLast || index < path.length - 1) ? item.href : undefined,
  }));
}

// Monta o caminho exibido acima do conteúdo, ex.: "Configurações / Usuários / Editar Usuário".
// Devolve lista vazia quando a rota não é conhecida (o chamador não mostra a faixa).
export function buildBreadcrumb(menuItems: ApiUserMenuItem[], pathname: string): BreadcrumbEntry[] {
  const own = findMenuPath(menuItems, pathname);
  if (own) return toEntries(own, false);

  for (const page of FORM_PAGES) {
    const isCreate = pathname === `${page.base}/create`;
    const isEdit = new RegExp(`^${page.base}/\\d+/edit$`).test(pathname);
    if (!isCreate && !isEdit) continue;

    const parent = findMenuPath(menuItems, page.parent);
    if (!parent) return [{ title: isCreate ? page.create : page.edit }];
    return [...toEntries(parent, true), { title: isCreate ? page.create : page.edit }];
  }

  return [];
}
