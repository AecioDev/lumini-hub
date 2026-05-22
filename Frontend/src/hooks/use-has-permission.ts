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
  const userPermissions =
    user.role?.permissions?.map((p) => p.permission) || [];
  return userPermissions.includes(permissionName);
}
