// src/config/routes.ts
export const routes = {
  public: {
    home: "/",
  },
  auth: {
    login: "/login",
    signIn: "/signin", // Adicionado com base na sua pasta (auth)/signin
  },
  dashboards: {
    root: "/dashboard", // O redirecionador principal
    admin: "/dashboard/admin",
    financial: "/dashboard/financial",
    general: "/dashboard/general",
    inventory: "/dashboard/inventory",
    manager: "/dashboard/manager",
    purchase: "/dashboard/purchase",
    sales: "/dashboard/sales",
  },
  // --- MÓDULOS/ENTIDADES CENTRAIS ---
  addresses: {
    root: "/addresses",
    create: "/addresses/create",
    view: (id: string | number) => `/addresses/${id}`,
    edit: (id: string | number) => `/addresses/${id}/edit`,
  },
  contacts: {
    root: "/contacts",
    create: "/contacts/create",
    view: (id: string | number) => `/contacts/${id}`,
    edit: (id: string | number) => `/contacts/${id}/edit`,
  },
  customers: {
    root: "/customers",
    create: "/customers/create",
    view: (id: string | number) => `/customers/${id}`,
    edit: (id: string | number) => `/customers/${id}/edit`,
  },
  documents: {
    root: "/documents",
    create: "/documents/create",
    view: (id: string | number) => `/documents/${id}`,
    edit: (id: string | number) => `/documents/${id}/edit`,
  },
  products: {
    root: "/products",
    create: "/products/create",
    view: (id: string | number) => `/products/${id}`,
    edit: (id: string | number) => `/products/${id}/edit`,
  },
  suppliers: {
    root: "/suppliers",
    create: "/suppliers/create",
    view: (id: string | number) => `/suppliers/${id}`,
    edit: (id: string | number) => `/suppliers/${id}/edit`,
  },

  // --- MÓDULOS FUNCIONAIS (com suas sub-rotas específicas) ---
  financial: {
    root: "/financial",
    accountsPayable: "/financial/accounts-payable",
    accountsReceivable: "/financial/accounts-receivable",
    // Se houver sub-rotas como /financial/accounts-receivable/boletos,
    // você adicionaria pastas e page.tsx correspondentes e as definiria aqui. Ex:
    // accountsReceivable: {
    //   root: "/financial/accounts-receivable",
    //   boletos: "/financial/accounts-receivable/boletos",
    // }
  },
  inventory: {
    // Anteriormente "estoque"
    root: "/inventory", // Página principal do módulo de inventário
    setup: {
      // Seus cadastros de inventário (produtos são centrais agora)
      root: "/inventory/setup",
      locations: "/inventory/setup/locations",
    },
    movements: {
      // Movimentações de Estoque
      root: "/inventory/movements", // Opcional: uma página geral de movimentações ou um log combinado
      inbound: "/inventory/movements/inbound", // Para Entradas de Notas/Mercadorias
      outbound: "/inventory/movements/outbound", // Para Saídas de Notas/Mercadorias
      adjustments: "/inventory/movements/adjustments", // Para Ajustes Manuais de Estoque
      // Se cada tipo de movimento gerar um "documento" que pode ser visualizado/editado:
      // viewInboundDoc: (id: string | number) => `/inventory/movements/inbound/${id}`,
      // createInbound: "/inventory/movements/inbound/create",
    },
    stockLevels: "/inventory/stock-levels", // Exemplo para consulta de níveis de estoque
    // Outras funcionalidades de inventário, como "Transferências", "Contagem Cíclica", etc.
  },
  purchases: {
    // Para o módulo de Compras e suas funcionalidades
    root: "/purchases", // Página principal do módulo de compras
    orders: {
      root: "/purchases/orders",
      create: "/purchases/orders/create",
      view: (id: string | number) => `/purchases/orders/${id}`,
      edit: (id: string | number) => `/purchases/orders/${id}/edit`,
    },
    // Outras funcionalidades de compras, como "XML Imports" se voltarem para cá,
    // ou "Purchase Reports", etc.
  },
  sales: {
    // Para o módulo de Vendas e suas funcionalidades
    root: "/sales", // Página principal do módulo de vendas
    orders: {
      root: "/sales/orders",
      create: "/sales/orders/create",
      view: (id: string | number) => `/sales/orders/${id}`,
      edit: (id: string | number) => `/sales/orders/${id}/edit`,
    },
    // Outras funcionalidades de vendas, como "Payment Plans" se voltarem para cá,
    // ou "Sales Reports", etc.
  },
  settings: {
    root: "/settings", // Uma página de configurações geral, se tiver
    permissions: {
      root: "/settings/permissions",
      create: "/settings/permissions/create",
      view: (id: string | number) => `/settings/permissions/${id}`,
      edit: (id: string | number) => `/settings/permissions/${id}/edit`,
    },
    roles: "/settings/roles", // Para a atribuição de permissões a roles (admin)
    users: {
      root: "/settings/users",
      create: "/settings/users/create",
      view: (id: string | number) => `/settings/users/${id}`,
      edit: (id: string | number) => `/settings/users/${id}/edit`,
    },
  },
};
