// components/layout/sidebar.tsx
"use client";

import type React from "react";

import { useEffect, useState, useMemo } from "react";
import { useAuth } from "@/hooks/use-auth";
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import { Logo } from "@/components/ui/logo";
import { ChevronLeft, ChevronRight } from "lucide-react";
import { systemNavItems, NavItem } from "@/config/navigation";
import { SidebarItem } from "./sidebarItem";
import { Permission } from "@/services/auth/permission-schema";

export function Sidebar() {
  const [collapsed, setCollapsed] = useState(false);
  const { user } = useAuth();

  useEffect(() => {
    const handleResize = () => {
      if (window.innerWidth < 768) {
        setCollapsed(true);
      } else {
        setCollapsed(false);
      }
    };

    handleResize();
    window.addEventListener("resize", handleResize);
    return () => window.removeEventListener("resize", handleResize);
  }, []);

  const filteredNavItems = useMemo(() => {
    if (!user || !user.role || !user.role.name || !user.role.permissions) {
      return [];
    }

    const userRoleName = user.role.name;
    const userPermissions = user.role.permissions.map(
      (p: Permission) => p.permission
    );

    const hasPermission = (permissionName: string): boolean => {
      return userPermissions.includes(permissionName);
    };

    const hasRequiredRole = (requiredRoles: string[]): boolean => {
      return requiredRoles.includes(userRoleName);
    };

    const filterItems = (items: NavItem[]): NavItem[] => {
      return items.filter((item) => {
        let isItemVisibleByRules = true; // Visibilidade inicial do item por suas regras

        // Prioridade 1: requiredPermission
        if (item.requiredPermission) {
          isItemVisibleByRules = hasPermission(item.requiredPermission);
        }
        // Prioridade 2: requiredRoles (se não houver requiredPermission no item)
        else if (item.requiredRoles) {
          isItemVisibleByRules = hasRequiredRole(item.requiredRoles);
        }
        // Se não tiver nenhum, isItemVisibleByRules permanece true (visível por padrão)

        // Se o item tem filhos, processa os filhos recursivamente
        if (item.children && item.children.length > 0) {
          const filteredChildren = filterItems(item.children); // Filtra os filhos

          // MUDANÇA CRUCIAL AQUI: Nova lógica para visibilidade do item pai
          if (filteredChildren.length > 0) {
            item.children = filteredChildren; // Atualiza os filhos para apenas os visíveis
            // O item pai é visível SE:
            // 1. Ele mesmo é visível pelas suas regras (role/permission) E tem filhos visíveis
            // OU
            // 2. Ele TEM um href próprio (é um link clicável) E tem filhos visíveis
            // (Esta segunda parte é uma salvaguarda caso a regra de role/permission não se aplique perfeitamente ao pai, mas ele ainda é um link principal)
            return isItemVisibleByRules || !!item.href;
          } else {
            // Se não tem filhos visíveis APÓS a filtragem:
            // Ele só é visível se tiver um href próprio (ou seja, é um item clicável por si só, mesmo sem filhos visíveis)
            return !!item.href;
          }
        }

        // Se o item não tem filhos:
        // Sua visibilidade é determinada apenas por suas próprias regras (isItemVisibleByRules)
        return isItemVisibleByRules;
      });
    };

    return filterItems(systemNavItems);
  }, [user]); // Dependência: A lista de itens filtrada recalcula apenas quando o usuário muda

  return (
    <aside
      className={cn(
        "flex h-screen flex-col border-r bg-card transition-all duration-300",
        collapsed ? "w-16" : "w-64"
      )}
    >
      <div className="flex h-16 items-center justify-between border-b px-4">
        <div
          className={cn(
            "flex items-center gap-2",
            collapsed && "justify-center w-full"
          )}
        >
          <Logo className={cn("h-8 w-8", collapsed ? "mx-auto" : "")} />
          {!collapsed && (
            <span className="text-lg font-semibold">ERP System</span>
          )}
        </div>
        <Button
          variant="ghost"
          size="icon"
          onClick={() => setCollapsed(!collapsed)}
          className={cn("h-8 w-8", collapsed && "right-0 mr-4")}
        >
          {collapsed ? (
            <ChevronRight className="h-4 w-4" />
          ) : (
            <ChevronLeft className="h-4 w-4" />
          )}
        </Button>
      </div>
      <nav className="flex-1 overflow-y-auto py-4">
        <ul className="space-y-1 px-2">
          {filteredNavItems.map((item) => (
            <SidebarItem key={item.title} item={item} collapsed={collapsed} />
          ))}
        </ul>
      </nav>
    </aside>
  );
}
