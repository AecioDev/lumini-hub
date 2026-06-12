// hooks/use-has-permission.ts
"use client";

import { useAuth } from "./use-auth"; // Reutiliza o hook de autenticação existente

export function useHasPermission(permissionName: string): boolean {
  const { user, isLoading } = useAuth();

  // Enquanto estiver carregando ou se não houver usuário, assuma que não tem permissão
  if (isLoading || !user) {
    return false;
  }

  // Verifique se o usuário tem a permissão necessária
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
  return userPermissions.includes(permissionName);
}
