# Design Spec: Dashboards dinâmicos na Sidebar

**Data:** 12 de Junho de 2026  
**Status:** Aprovado pelo Usuário  
**Autor:** Antigravity

## Objetivo
Substituir o switcher de dashboard do cabeçalho (dropdown select) por uma lista dinâmica de submenus na sidebar abaixo do item "Dashboard", caso o usuário possua permissão para acessar mais de um dashboard.

## Arquitetura e Fluxo de Dados
1. O componente `SystemSidebarMenu` lerá as permissões do usuário logado através do hook `useAuth()`.
2. Filtraremos a lista de `dashboardOptions` (importada de `@/config/navigation`) com base nessas permissões.
3. Se o usuário tiver acesso a **mais de um** dashboard:
   - Injetamos dinamicamente os dashboards disponíveis como `dropdownItems` no item de menu "Dashboard".
   - Removemos a propriedade `href` do item "Dashboard" raiz para que ele funcione estritamente como um gatilho de recolhimento/expansão (collapse).
4. O switcher de dashboard (`DashboardSwitcher`) será removido do cabeçalho `SystemHeader`.

## Componentes Afetados
- `src/components/layout/system-sidebar-menu.tsx`: Adição da lógica dinâmica de construção do menu de dashboards.
- `src/components/layout/system-header.tsx`: Remoção do componente `<DashboardSwitcher />`.
