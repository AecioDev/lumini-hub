// UI para listar clientes (rota: routes.vendas.cadastros.clientes.root)
"use client";

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { routes } from "@/config/routes";
import { CustomersList } from "@/components/customers/common/customers-list";
import { PagePermissionGuard } from "@/components/layout/PagePermissionGuard";
import { PermissionedLinkButton } from "@/components/common/PermissionedLinkButton";

export default function CustomersPage() {
  return (
    <PagePermissionGuard
      requiredPermissions={["customers.view"]}
      accessDeniedMessage="Você não tem permissão para visualizar clientes."
    >
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold">Clientes</h1>
            <p className="text-muted-foreground">
              Gerencie os clientes do seu negócio
            </p>
          </div>
          <div>
            <PermissionedLinkButton
              href={routes.customers.create}
              permission="customers.create"
              tooltipMessage="Você não tem permissão para criar clientes."
              className="w-full @lg:w-auto"
              iconName="mdi:plus" // Certifique-se de que este ícone está disponível
            >
              Novo Cliente
            </PermissionedLinkButton>
          </div>
        </div>

        <Card>
          <CardHeader>
            <CardTitle>Lista de Clientes</CardTitle>
            <CardDescription>
              Gerencie os clientes cadastrados no sistema
            </CardDescription>
          </CardHeader>
          <CardContent>
            {/* O componente CustomersList agora deve ser responsável por
                carregar e exibir a lista de clientes, e gerenciar sua própria
                lógica de edição, caso ela ainda seja necessária. */}
            <CustomersList />
          </CardContent>
        </Card>
      </div>
    </PagePermissionGuard>
  );
}
