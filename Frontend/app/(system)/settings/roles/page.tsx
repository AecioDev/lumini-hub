// src/app/(system)/settings/roles/page.tsx
"use client";

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { PagePermissionGuard } from "@/components/layout/PagePermissionGuard";
import { PermissionedLinkButton } from "@/components/common/PermissionedLinkButton";
import { useState } from "react";
import { RoleList } from "@/components/settings/roles/tables/roles-list";

export default function RolesPage() {
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);

  // Esta função será passada para o modal de criação
  // e será responsável por, provavelmente, recarregar a lista de perfis
  // e/ou atualizar o perfil selecionado após a criação.
  const handleRoleListRefresh = () => {
    // Aqui você chamaria uma função que atualize a lista de roles,
    // por exemplo, buscando novamente todos os roles ou
    // atualizando o estado do RoleList se ele tiver um mecanismo para isso.
    // Por enquanto, podemos deixar vazio, e ajustar quando tivermos o RoleList.
  };

  return (
    <PagePermissionGuard
      requiredPermissions={["roles.view"]}
      accessDeniedMessage="Você não tem permissão para visualizar a lista de Perfis."
    >
      <div className="space-y-6 px-4 py-6 md:px-8 md:py-8 bg-background text-foreground min-h-screen">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold tracking-tight lg:text-4xl">
              Perfis
            </h1>
            <p className="text-lg text-muted-foreground mt-1">
              Gerencie os perfis de usuário
            </p>
          </div>
          <div>
            <PermissionedLinkButton
              href="#" // O href pode ser "#" já que vai abrir um modal
              onClick={() => setIsCreateModalOpen(true)}
              permission="roles.create"
              tooltipMessage="Você não pode criar novos perfis."
              className="w-full @lg:w-auto"
              iconName="mdi:plus"
            >
              Novo Perfil
            </PermissionedLinkButton>
          </div>
        </div>

        <Card className="shadow-lg border-primary/20">
          <CardHeader>
            <CardTitle className="text-xl font-semibold">
              Lista de Perfis
            </CardTitle>
            <CardDescription>
              Visualize e gerencie os perfis de usuário cadastrados no sistema.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <RoleList onRoleCreated={handleRoleListRefresh} />
          </CardContent>
        </Card>
      </div>

      {/* Modal de Criação de Perfil 
      <CreateRoleModal
        open={isCreateModalOpen}
        onOpenChange={setIsCreateModalOpen}
        onSuccess={handleRoleListRefresh} // Chama a função para atualizar a lista após criar um perfil
      />*/}
    </PagePermissionGuard>
  );
}
