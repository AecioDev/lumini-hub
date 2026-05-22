// UI para listar clientes (rota: routes.vendas.cadastros.clientes.root)
"use client";

import EmConstrucao from "@/components/em-construcao";
import { PagePermissionGuard } from "@/components/layout/PagePermissionGuard";

export default function CustomersPage() {
  return (
    <PagePermissionGuard
      requiredPermissions={["inventory.view"]}
      accessDeniedMessage="Você não tem permissão para visualizar os locais de estoque."
    >
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold">Locais de Estoque</h1>
            <p className="text-muted-foreground">
              Gerencie os estoques do seu negócio
            </p>
          </div>
        </div>
        <EmConstrucao />
      </div>
    </PagePermissionGuard>
  );
}
