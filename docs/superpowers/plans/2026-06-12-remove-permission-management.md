# Remover Gerenciamento de Permissões - Plano de Implementação

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Remover as opções e telas do menu e código frontend para gerenciar permissões, mantendo o serviço de listagem que é consumido pelos modais de Perfil e Usuários.

**Architecture:** Modificar `system-menu-items.tsx` para tirar a opção "Permissões". Excluir a rota `app/(system)/settings/permissions` e a pasta de componentes `src/components/settings/permission`.

**Tech Stack:** React, Next.js.

---

### Task 1: Remover Opção de Menu e Excluir Arquivos de Permissão

**Files:**
- Modify: `src/components/layout/system-menu-items.tsx:104-108`
- Delete: `app/(system)/settings/permissions/page.tsx`
- Delete: `src/components/settings/permission/`

- [ ] **Step 1: Remover o item "Permissões" do menu lateral**
  Excluir do array `menuItems` em `system-menu-items.tsx`:
  ```typescript
  {
    name: 'Permissões',
    href: routes.settings.permissions.root,
  },
  ```

- [ ] **Step 2: Excluir a rota física de permissões**
  Remover o arquivo `app/(system)/settings/permissions/page.tsx`.

- [ ] **Step 3: Excluir a pasta de componentes de permissão**
  Remover a pasta `src/components/settings/permission/` de forma recursiva.
