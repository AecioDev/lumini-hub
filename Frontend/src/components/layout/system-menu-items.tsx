import { routes } from '@/config/routes';
import {
  PiCurrencyCircleDollarFill,
  PiShoppingCart,
  PiPackage,
  PiChartBar,
  PiSquaresFour,
  PiGridFour,
  PiUserGear,
  PiUser,
  PiBriefcase,
} from 'react-icons/pi';

export interface MenuItem {
  name: string;
  href?: string;
  icon?: React.ReactNode;
  badge?: string;
  dropdownItems?: MenuItem[];
}

export const menuItems: MenuItem[] = [
  // Visão Geral
  {
    name: 'Visão Geral',
  },
  {
    name: 'Dashboard',
    href: routes.dashboards.root,
    icon: <PiSquaresFour size={20} />,
  },
  
  // Cadastros Principais
  {
    name: 'Cadastros',
  },
  {
    name: 'Clientes',
    href: routes.customers.root,
    icon: <PiUser size={20} />,
  },
  {
    name: 'Fornecedores',
    href: routes.suppliers.root,
    icon: <PiBriefcase size={20} />,
  },
  {
    name: 'Produtos',
    href: routes.products.root,
    icon: <PiPackage size={20} />,
  },

  // Módulos ERP
  {
    name: 'Operações',
  },
  {
    name: 'Vendas',
    href: routes.sales.root,
    icon: <PiChartBar size={20} />,
  },
  {
    name: 'Compras',
    href: routes.purchases.root,
    icon: <PiShoppingCart size={20} />,
  },
  {
    name: 'Estoque',
    href: routes.inventory.root,
    icon: <PiGridFour size={20} />,
  },
  {
    name: 'Financeiro',
    href: routes.financial.root,
    icon: <PiCurrencyCircleDollarFill size={20} />,
    dropdownItems: [
      {
        name: 'Contas a Pagar',
        href: routes.financial.accountsPayable,
      },
      {
        name: 'Contas a Receber',
        href: routes.financial.accountsReceivable,
      },
    ],
  },

  // Configurações
  {
    name: 'Administração',
  },
  {
    name: 'Configurações',
    href: '#',
    icon: <PiUserGear size={20} />,
    dropdownItems: [
      {
        name: 'Usuários e Perfis',
        href: routes.settings.users.root,
      },
    ],
  },
];
