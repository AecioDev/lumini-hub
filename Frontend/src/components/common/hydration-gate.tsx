// components/hydration-gate.tsx
"use client"; // ESSENCIAL para usar hooks e acessar localStorage

import { useEffect, useState, type ReactNode } from "react";
import { useAtomValue } from "jotai"; // Assumindo que userAtom é do Jotai
import { userAtom } from "@/atoms/userAtom"; // Ajuste o caminho para seu userAtom
import { ContainerLoader } from "@/components/ContainerLoader"; // Seu componente de loader

interface HydrationGateProps {
  children: ReactNode;
}

export function HydrationGate({ children }: HydrationGateProps) {
  const user = useAtomValue(userAtom); // Seu átomo do usuário
  const [hasHydrated, setHasHydrated] = useState(false);

  useEffect(() => {
    // A "gambiarra infernal" (eu te entendo, rsrs) para garantir a hidratação no cliente.
    // Isso evita o erro de divergência entre o render do servidor e o primeiro render do cliente.
    setTimeout(() => setHasHydrated(true), 10);
  }, []);

  if (!hasHydrated) {
    // Enquanto não hidratou, exibe o loader
    return (
      <div className="flex h-screen items-center justify-center">
        <ContainerLoader>
          <h1 className="text-xl">carregando...</h1>
        </ContainerLoader>
      </div>
    );
  }

  // Depois de hidratado, renderiza o conteúdo normal da aplicação
  return <>{children}</>;
}
