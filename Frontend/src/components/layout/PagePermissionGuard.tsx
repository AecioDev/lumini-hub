// src/components/layout/PagePermissionGuard.tsx
"use client";

import React, { useState, useEffect, ReactNode } from "react";
// Importe o hook e a interface de opções dele
import {
  useRequirePermission,
  RequirePermissionOptions,
} from "@/hooks/use-require-permission";

import { LoadingSpinner } from "@/components/ui/loading-spinner"; // Certifique-se do caminho

// Props do guardião agora usam a interface importada
interface PagePermissionGuardProps extends RequirePermissionOptions {
  children: ReactNode;
  loadingMessage?: string;
}

export function PagePermissionGuard({
  // Props de RequirePermissionOptions
  requiredPermissions,
  requiredRoles,
  redirectPath, // Usado pelo hook se o acesso for negado
  accessDeniedMessage, // Usado pelo hook se o acesso for negado
  // Props específicas do Guard
  children,
  loadingMessage = "Carregando, aguarde...",
}: PagePermissionGuardProps) {
  // Estado de carregamento interno do guardião
  const [isCheckingPermission, setIsCheckingPermission] = useState(true);

  const { hasPermissionCheckCompleted } = useRequirePermission({
    requiredPermissions,
    requiredRoles,
    redirectPath,
    accessDeniedMessage,
  });

  useEffect(() => {
    if (hasPermissionCheckCompleted) {
      // Se a verificação foi completada, significa que o acesso foi concedido
      // (pois o hook redirecionaria em caso de negação antes de completar)
      setIsCheckingPermission(false);
    }
    // Se hasPermissionCheckCompleted for false e o hook estiver no processo de redirecionamento,
    // este componente provavelmente será desmontado, e isCheckingPermission não importará mais.
  }, [hasPermissionCheckCompleted]);

  if (isCheckingPermission) {
    // Mostra o spinner enquanto hasPermissionCheckCompleted ainda é false
    // ou enquanto o estado interno isCheckingPermission não foi atualizado.
    return (
      <div className="flex h-screen w-full items-center justify-center">
        <LoadingSpinner message={loadingMessage} />
      </div>
    );
  }

  // Se isCheckingPermission é false, significa que hasPermissionCheckCompleted foi true,
  // o que implica que o acesso foi concedido e o hook não redirecionou.
  // Então, podemos renderizar o conteúdo da página.
  return <>{children}</>;
}
