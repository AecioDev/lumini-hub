# Menu de Dashboards Dinâmicos na Sidebar - Plano de Implementação

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Mostrar os dashboards disponíveis para o usuário como submenus abaixo do item "Dashboard" na Sidebar caso ele tenha acesso a mais de um, e remover o dropdown de troca no cabeçalho.

**Architecture:** Modificar o `SystemSidebarMenu` para usar `useAuth()` e obter as permissões do usuário logado, filtrando a lista `dashboardOptions` e gerando dinamicamente `dropdownItems` no item "Dashboard". Modificar o `SystemHeader` para remover o `DashboardSwitcher`.

**Tech Stack:** React, Next.js, Jotai, TailwindCSS, Rizzui.

---

### Task 1: Remover o DashboardSwitcher do SystemHeader

**Files:**
- Modify: `src/components/layout/system-header.tsx:45-49`

- [ ] **Step 1: Remover o import do DashboardSwitcher**
  Excluir a linha:
  ```typescript
  import { DashboardSwitcher } from './dashboard-switcher';
  ```

- [ ] **Step 2: Remover a renderização do DashboardSwitcher**
  Remover a renderização do switcher do corpo do componente:
  ```typescript
  {/* Switcher de Dashboard (reaproveitado do layout original) */}
  <div className="hidden sm:block">
    <DashboardSwitcher />
  </div>
  ```

---

### Task 2: Implementar Submenus de Dashboards Dinâmicos no SystemSidebarMenu

**Files:**
- Modify: `src/components/layout/system-sidebar-menu.tsx`

- [ ] **Step 1: Adicionar imports de useAuth, useMemo e dashboardOptions**
  Importar os hooks e as opções de dashboards no arquivo `system-sidebar-menu.tsx`:
  ```typescript
  import { useMemo } from 'react';
  import { useAuth } from '@/hooks/use-auth';
  import { dashboardOptions } from '@/config/navigation';
  ```

- [ ] **Step 2: Calcular os dashboards disponíveis e construir itens de menu dinâmicos**
  Inserir a lógica de cálculo das permissões e geração de submenus para "Dashboard" dentro do componente `SystemSidebarMenu`:
  ```typescript
  const { user, isLoading: isLoadingUser } = useAuth();

  const availableDashboards = useMemo(() => {
    if (isLoadingUser || !user) {
      return [];
    }
    let userPermissions: string[] = [];
    if (user.role?.permissions && Array.isArray(user.role.permissions)) {
      userPermissions = user.role.permissions.map((p: any) => p.permission);
    }
    if (Array.isArray(user.permissions)) {
      user.permissions.forEach((p: any) => {
        const permName = typeof p === "string" ? p : p.permission;
        if (permName && !userPermissions.includes(permName)) {
          userPermissions.push(permName);
        }
      });
    }
    return dashboardOptions.filter((option) =>
      option.requiredPermission ? userPermissions.includes(option.requiredPermission) : true
    );
  }, [user, isLoadingUser]);

  const dynamicMenuItems = useMemo(() => {
    return menuItems.map((item) => {
      if (item.name === 'Dashboard') {
        if (availableDashboards.length > 1) {
          return {
            ...item,
            href: undefined, // remove href para abrir o collapse
            dropdownItems: availableDashboards.map((dash) => ({
              name: dash.title,
              href: dash.href,
            })),
          };
        } else if (availableDashboards.length === 1) {
          return {
            ...item,
            href: availableDashboards[0].href,
            dropdownItems: undefined,
          };
        }
      }
      return item;
    });
  }, [availableDashboards]);
  ```

- [ ] **Step 3: Atualizar o mapeamento para renderizar os itens dinâmicos**
  Alterar a renderização do menu para utilizar `dynamicMenuItems.map` em vez de `menuItems.map`:
  ```typescript
  return (
    <div className="mt-4 pb-3 3xl:mt-6">
      {dynamicMenuItems.map((item, index) => {
  ```
