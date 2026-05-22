// src/config/navigation.ts
import { routes } from "./routes";

export interface NavItem {
  title: string;
  href?: string;
  icon: string; // O tipo do ícone é uma string para o nome do ícone do Iconify
  requiredRoles?: string[];
  requiredPermission?: string;
  children?: NavItem[];
}

export const systemNavItems: NavItem[] = [
  {
    title: "Home",
    href: routes.dashboards.root,
    icon: "mdi:home", // Ícone MDI para Home
    requiredPermission: "dashboard.view_default", // Adicionando permissão para o Dashboard padrão
  },
  {
    title: "Vendas",
    icon: "mdi:cart", // Ícone MDI para carrinho de compras
    requiredRoles: ["ADMIN", "VENDAS", "GERENTE"],
    children: [
      {
        title: "Pedidos",
        href: routes.sales.orders.root,
        icon: "mdi:clipboard-list", // MDI para lista/clipboard
        requiredPermission: "orders.view", // Descomentar e adicionar permissão
      },
      {
        title: "Cadastros",
        icon: "mdi:cog", // MDI para engrenagem/settings
        children: [
          {
            title: "Clientes",
            href: routes.customers.root,
            icon: "mdi:account-group", // MDI para grupo de usuários
            requiredPermission: "customers.view", // Descomentar e adicionar permissão
          },
        ],
      },
    ],
  },
  {
    title: "Compras",
    icon: "mdi:briefcase", // MDI para pasta/maleta
    requiredRoles: ["ADMIN", "COMPRAS", "GERENTE"],
    children: [
      {
        title: "Pedidos de Compras",
        href: routes.purchases.orders.root,
        icon: "mdi:clipboard-list",
        requiredPermission: "purchase_orders.view", // Descomentar e adicionar permissão
      },
      {
        title: "Cadastros",
        icon: "mdi:cog",
        children: [
          {
            title: "Fornecedores",
            href: routes.suppliers.root,
            icon: "mdi:truck",
            requiredPermission: "suppliers.view", // Descomentar e adicionar permissão
          },
        ],
      },
    ],
  },
  {
    title: "Estoque",
    icon: "mdi:package-variant", // MDI para pacote
    requiredRoles: ["ADMIN", "ESTOQUE", "GERENTE"],
    children: [
      {
        title: "Produtos",
        href: routes.products.root,
        icon: "mdi:archive", // MDI para arquivo/caixa
        requiredPermission: "products.view", // Descomentar e adicionar permissão
      },
      {
        title: "Movimentação",
        icon: "mdi:swap-horizontal", // MDI para troca/movimentação
        children: [
          {
            title: "Notas de Entrada",
            href: routes.inventory.movements.inbound,
            icon: "mdi:arrow-down-box", // MDI para seta para baixo em caixa
            requiredPermission: "inventory.movements.inbound.view",
          },
          {
            title: "Notas de Saídas",
            href: routes.inventory.movements.outbound,
            icon: "mdi:arrow-up-box", // MDI para seta para cima em caixa
            requiredPermission: "inventory.movements.outbound.view",
          },
          {
            title: "Ajuste de Estoque",
            href: routes.inventory.movements.adjustments,
            icon: "mdi:tune-variant", // MDI para ajuste
            requiredPermission: "inventory.movements.adjustments.view",
          },
        ],
      },
      {
        title: "Níveis de Estoque",
        href: routes.inventory.stockLevels,
        icon: "mdi:package-variant-closed", // MDI para pacote fechado
        requiredPermission: "inventory.stock_levels.view",
      },
      {
        title: "Cadastro",
        icon: "mdi:cog",
        children: [
          {
            title: "Locations",
            href: routes.inventory.setup.locations,
            icon: "mdi:map-marker", // MDI para localização
            requiredPermission: "inventory.setup.locations.view", // Permissão específica
          },
        ],
      },
    ],
  },
  {
    title: "Financeiro",
    icon: "mdi:wallet", // MDI para carteira
    requiredRoles: ["ADMIN", "FINANCEIRO", "GERENTE"],
    children: [
      {
        title: "Dashboard Financeiro",
        href: routes.dashboards.financial,
        icon: "mdi:view-dashboard", // MDI para dashboard
        requiredPermission: "dashboard.financial.view",
      },
      {
        title: "Contas a Receber",
        href: routes.financial.accountsReceivable,
        icon: "mdi:receipt-text", // MDI para recibo/conta
        requiredPermission: "financial.accounts_receivable.view", // Permissão para o módulo
        children: [
          {
            title: "Clientes",
            href: routes.customers.root,
            icon: "mdi:account-group",
            requiredPermission: "customers.view",
          },
        ],
      },
      {
        title: "Contas a Pagar",
        href: routes.financial.accountsPayable,
        icon: "mdi:bank", // MDI para banco/instituição financeira
        requiredPermission: "financial.accounts_payable.view", // Permissão para o módulo
        children: [
          {
            title: "Fornecedores",
            href: routes.suppliers.root,
            icon: "mdi:truck",
            requiredPermission: "suppliers.view",
          },
        ],
      },
    ],
  },
  {
    title: "Configurações",
    icon: "mdi:cog-box", // MDI para caixa de configurações
    requiredRoles: ["ADMIN"],
    children: [
      {
        title: "Usuários",
        href: routes.settings.users.root,
        icon: "mdi:account-circle", // MDI para perfil de usuário
        requiredPermission: "users.view",
      },
      {
        title: "Permissões",
        href: routes.settings.permissions.root,
        icon: "mdi:lock", // MDI para cadeado
        requiredPermission: "admin.create_permissions",
      },
      {
        title: "Perfis de Acesso",
        href: routes.settings.roles,
        icon: "mdi:shield-account", // MDI para perfil de segurança
        requiredPermission: "roles.view",
      },
    ],
  },
];

// --- OPÇÕES PARA O DASHBOARD SWITCHER (Header) ---
export interface DashboardOption {
  title: string;
  href: string;
  requiredPermission?: string; // Mantido como opcional
}

export const dashboardOptions: DashboardOption[] = [
  {
    title: "Meu Dashboard",
    href: routes.dashboards.general,
    requiredPermission: "dashboard.view_default",
  },
  {
    title: "Dashboard Admin",
    href: routes.dashboards.admin,
    requiredPermission: "dashboard.admin.view",
  },
  {
    title: "Dashboard Compras",
    href: routes.dashboards.purchase,
    requiredPermission: "dashboard.purchase.view",
  },
  {
    title: "Dashboard Estoque",
    href: routes.dashboards.inventory,
    requiredPermission: "dashboard.inventory.view",
  },
  {
    title: "Dashboard Vendas",
    href: routes.dashboards.sales,
    requiredPermission: "dashboard.sales.view",
  },
  {
    title: "Dashboard Financeiro",
    href: routes.dashboards.financial,
    requiredPermission: "dashboard.financial.view",
  },
  {
    title: "Dashboard Gerencial",
    href: routes.dashboards.manager,
    requiredPermission: "dashboard.manager.view",
  },
];
