"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "./use-auth";
import { useToast } from "@/components/ui/use-toast";

export interface RequirePermissionOptions {
  requiredPermissions?: string[];
  requiredRoles?: string[];
  redirectPath?: string;
  accessDeniedMessage?: string;
}

export function useRequirePermission(options: RequirePermissionOptions) {
  const { user, isLoading } = useAuth();
  const router = useRouter();
  const { toast } = useToast();

  const [hasPermissionCheckCompleted, setHasPermissionCheckCompleted] =
    useState(false);

  useEffect(() => {
    if (isLoading) {
      return;
    }

    // Caso 1: Usuário não está logado
    if (!user) {
      toast({
        variant: "destructive",
        title: "Acesso Negado",
        description: "Você precisa estar logado para acessar esta página.",
      });
      const timer = setTimeout(() => {
        router.replace("/login"); // Sempre para /login se não há usuário
      }, 50); // O timeout é pequeno, mas a função logout() deveria ser mais rápida
      return () => clearTimeout(timer); // Importante: previne setHasPermissionCheckCompleted(true)
    }

    // Caso 2: Usuário está logado, verificar permissões/roles
    let hasRequiredPermissions = true;
    if (options.requiredPermissions && options.requiredPermissions.length > 0) {
      const userPermissions =
        user.role?.permissions?.map((p) => p.permission) || [];
      hasRequiredPermissions = options.requiredPermissions.some((perm) =>
        userPermissions.includes(perm)
      );
    }

    let hasRequiredRoles = true;
    if (options.requiredRoles && options.requiredRoles.length > 0) {
      const userRoleName = user.role?.name;
      hasRequiredRoles = options.requiredRoles.includes(userRoleName || "");
    }

    // Se chegou aqui e alguma verificação (permissão OU role) falhou
    if (!hasRequiredPermissions || !hasRequiredRoles) {
      toast({
        variant: "destructive",
        title: "Acesso Negado",
        description:
          options.accessDeniedMessage ||
          "Você não tem permissão para acessar esta funcionalidade.",
      });
      const timer = setTimeout(() => {
        router.replace(options.redirectPath || "/dashboard"); // Redireciona para o path definido ou /dashboard
      }, 50);
      return () => clearTimeout(timer); // Importante: previne setHasPermissionCheckCompleted(true)
    }

    // Se passou por todas as verificações e o usuário está logado com acesso
    setHasPermissionCheckCompleted(true);
  }, [
    isLoading,
    user,
    router,
    toast,
    // Passar o objeto options inteiro para a dependência
    // Ou as propriedades específicas que *realmente* podem mudar
    options.requiredPermissions,
    options.requiredRoles,
    options.redirectPath, // Agora estas são dependências diretas de `options`
    options.accessDeniedMessage, // O que é mais estável
  ]);

  return { hasPermissionCheckCompleted };
}
