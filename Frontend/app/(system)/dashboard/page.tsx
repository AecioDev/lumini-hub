"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/hooks/use-auth"; // Seu hook de autenticação
import { LoadingSpinner } from "@/components/ui/loading-spinner";

export default function DashboardIndexPage() {
  const { user, isLoading } = useAuth();
  const router = useRouter();

  useEffect(() => {
    // Só tentamos redirecionar se o carregamento do usuário estiver completo
    // e se tivermos um objeto de usuário.
    if (!isLoading && user) {
      // Usar router.replace para que o usuário não possa voltar para esta página de redirecionamento
      // com o botão 'voltar' do navegador.
      switch (user.role.name) {
        case "ADMIN":
          router.replace("/dashboard/admin");
          break;
        case "FINANCEIRO":
          router.replace("/dashboard/finance");
          break;
        case "ESTOQUE":
          router.replace("/dashboard/inventory");
          break;
        case "GERENTE":
          router.replace("/dashboard/manager");
          break;
        case "COMPRAS":
          router.replace("/dashboard/purchases");
          break;
        case "VENDAS":
          router.replace("/dashboard/sales");
          break;
        default:
          // Caso a role não seja reconhecida ou não tenha um dashboard específico
          router.replace("/dashboard/general"); // Ou uma página de erro/dashboard padrão
          break;
      }
    }
  }, [isLoading, user, router]); // Dependências do useEffect

  // Exibir um spinner enquanto o redirecionamento acontece
  // Ou um componente de carregamento simples.
  if (isLoading) {
    <div className="flex h-screen items-center justify-center">
      <LoadingSpinner message="Carregando..." />
    </div>;
  }

  // Se o usuário está presente mas ainda não redirecionou (por algum atraso),
  // ou se não há role correspondente, podemos exibir algo temporário ou null.
  return (
    <div className="flex h-screen items-center justify-center">
      <LoadingSpinner message="Redirecionando para o seu dashboard..." />
    </div>
  );
}
