"use client";

import React, { useMemo } from "react";
import { useRouter, usePathname } from "next/navigation";
import { useAuth } from "@/hooks/use-auth"; // Para pegar as permissões do usuário
import { dashboardOptions } from "@/config/navigation";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Permission } from "@/services/auth/permission-schema";

export function DashboardSwitcher() {
  const { user, isLoading: isLoadingUser } = useAuth();
  const router = useRouter();
  const pathname = usePathname(); // Para saber qual dashboard está ativo

  // Filtra os dashboards que o usuário tem permissão para ver
  const availableUserDashboards = useMemo(() => {
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

  const handleDashboardChange = (newHref: string) => {
    if (newHref) {
      router.push(newHref);
    }
  };

  // Não renderiza o switcher se estiver carregando,
  // se não houver opções disponíveis, ou se houver apenas uma (ou nenhuma).
  // Ou se não estivermos em uma rota de dashboard (começando com /dashboard/)
  if (
    isLoadingUser ||
    availableUserDashboards.length <= 1 ||
    !pathname.startsWith("/dashboard")
  ) {
    return null;
  }

  return (
    <Select
      value={pathname} // O valor selecionado é a rota atual
      onValueChange={handleDashboardChange}
    >
      <SelectTrigger className="w-[220px]">
        {" "}
        {/* Ajuste o tamanho */}
        <SelectValue placeholder="Mudar Dashboard" />
      </SelectTrigger>
      <SelectContent>
        {availableUserDashboards.map((option) => (
          <SelectItem key={option.href} value={option.href}>
            {option.title}
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  );
}
