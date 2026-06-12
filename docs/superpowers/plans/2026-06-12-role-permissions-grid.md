# Grade de Permissões de Perfis - Plano de Implementação

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Alterar o nome "Sistema ERP" para "Lumini Hub" no login/layout e unificar a visualização de permissões nos modais de perfil (criação e edição) para usarem a grade unificada com busca e filtro por módulo que existe em usuários.

**Architecture:** Modificar `login/page.tsx` e `layout.tsx` para mudar o texto/title. Modificar `edit-role-dialog.tsx` e `create-role-dialog.tsx` para implementar lógica de filtro e busca de permissões, adicionando a tabela de grid e ocultando a lógica de cards.

**Tech Stack:** React, Next.js, TailwindCSS, Lucide Icons, Shadcn components.

---

### Task 1: Ajustar Nome do Projeto no Login e Layout

**Files:**
- Modify: `app/(auth)/login/page.tsx:12`
- Modify: `app/layout.tsx:12-13`

- [ ] **Step 1: Modificar h1 em login/page.tsx**
  Substituir "Sistema ERP" por "Lumini Hub":
  ```tsx
  <h1 className="text-2xl font-bold">Lumini Hub</h1>
  ```

- [ ] **Step 2: Modificar metadados de title no layout.tsx**
  Substituir "Sistema ERP" por "Lumini Hub" no objeto de metadados:
  ```typescript
  export const metadata = {
    title: "Lumini Hub - Gestão Empresarial",
    description: "Sistema integrado de gestão empresarial",
  ```

---

### Task 2: Implementar a Grade de Permissões no EditRoleDialog

**Files:**
- Modify: `src/components/settings/roles/dialogs/edit-role-dialog.tsx`

- [ ] **Step 1: Adicionar imports necessários para Filtros e Ícones**
  Adicionar imports de Search, X, Checkbox e Input:
  ```typescript
  import { Search, X } from "lucide-react";
  import { Input } from "@/components/ui/input";
  import { useMemo } from "react";
  import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
  } from "@/components/ui/select";
  ```

- [ ] **Step 2: Adicionar estados e filtros locais**
  No corpo do componente, instanciar os estados de busca e filtros:
  ```typescript
  const [searchQuery, setSearchQuery] = useState("");
  const [filterModule, setFilterModule] = useState("all");

  // Lista plana para filtros fáceis
  const allPermissionsFlat = useMemo(() => {
    const flatList: (Permission & { module: string })[] = [];
    Object.entries(permissionsByModule).forEach(([moduleName, items]) => {
      items.forEach((p) => {
        flatList.push({ ...p, module: moduleName });
      });
    });
    return flatList.sort((a, b) => {
      const modCompare = a.module.localeCompare(b.module);
      if (modCompare !== 0) return modCompare;
      return a.permission.localeCompare(b.permission);
    });
  }, [permissionsByModule]);

  // Lista de módulos para o Select
  const modulesList = useMemo(() => {
    return Object.keys(permissionsByModule).sort();
  }, [permissionsByModule]);

  // Filtra as permissões
  const filteredPermissions = useMemo(() => {
    return allPermissionsFlat.filter((p) => {
      const idNum = typeof p.id === "string" ? parseInt(p.id, 10) : p.id;
      if (searchQuery.trim()) {
        const q = searchQuery.toLowerCase();
        const nameMatch = p.permission.toLowerCase().includes(q);
        const descMatch = p.description?.toLowerCase().includes(q) || false;
        const modMatch = p.module.toLowerCase().includes(q);
        if (!nameMatch && !descMatch && !modMatch) return false;
      }
      if (filterModule !== "all" && p.module !== filterModule) {
        return false;
      }
      return true;
    });
  }, [allPermissionsFlat, searchQuery, filterModule]);

  const isAllFilteredChecked = useMemo(() => {
    if (filteredPermissions.length === 0) return false;
    return filteredPermissions.every((p) => {
      const idNum = typeof p.id === "string" ? parseInt(p.id, 10) : p.id;
      return selectedPermissionIds.has(idNum);
    });
  }, [filteredPermissions, selectedPermissionIds]);

  const handleSelectAllFiltered = (checked: boolean) => {
    setSelectedPermissionIds((prev) => {
      const next = new Set(prev);
      filteredPermissions.forEach((p) => {
        const idNum = typeof p.id === "string" ? parseInt(p.id, 10) : p.id;
        if (checked) {
          next.add(idNum);
        } else {
          next.delete(idNum);
        }
      });
      return next;
    });
  };

  const handleClearFilters = () => {
    setSearchQuery("");
    setFilterModule("all");
  };
  ```

- [ ] **Step 3: Substituir renderização do bloco de permissões**
  Substituir o bloco de exibição dos cards pela grade de tabela estruturada com barra de filtros, exatamente como em `EditUserDialog`.

---

### Task 3: Implementar a Grade de Permissões no CreateRoleDialog

**Files:**
- Modify: `src/components/settings/roles/dialogs/create-role-dialog.tsx`

- [ ] **Step 1: Adicionar imports e estados de filtros locais**
  Adicionar os mesmos imports e lógica de estados/filtros locais descritos na Task 2 no arquivo `create-role-dialog.tsx`.

- [ ] **Step 2: Substituir renderização do bloco de permissões**
  Trocar os cards da criação do perfil pela nova grade de tabela estruturada com busca e filtro de módulos.
