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
import { Permission } from "@/services/auth/role-service"; // Sua tipagem de Permission

export function DashboardSwitcher() {
  const { user, isLoading: isLoadingUser } = useAuth();
  const router = useRouter();
  const pathname = usePathname(); // Para saber qual dashboard está ativo

  // Filtra os dashboards que o usuário tem permissão para ver
  const availableUserDashboards = useMemo(() => {
    if (isLoadingUser || !user || !user.role || !user.role.permissions) {
      return [];
    }
    const userPermissions = user.role.permissions.map(
      (p: Permission) => p.name
    );
    return dashboardOptions.filter((option) =>
      userPermissions.includes(option.requiredPermission)
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
