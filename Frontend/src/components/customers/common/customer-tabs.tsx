"use client";

import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { PagePermissionGuard } from "@/components/layout/PagePermissionGuard";
import { routes } from "@/config/routes";
import BasicInfoCustomerForm from "../forms/basic-info-customer-form";

export default function CustomerTabs() {
  return (
    <PagePermissionGuard
      requiredPermissions={["customers.create"]}
      accessDeniedMessage="Você não tem permissão para cadastrar clientes."
      redirectPath={routes.customers.root}
    >
      <div className="w-full max-w-4xl mx-auto p-6 bg-card text-card-foreground rounded-lg shadow-lg">
        <h1 className="text-2xl font-bold text-foreground mb-6">
          Cadastro de Clientes
        </h1>

        <Tabs defaultValue="basic" className="w-full">
          <TabsList className="grid w-full grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-2 mb-6 bg-form-background p-1 rounded-lg">
            <TabsTrigger
              value="basic"
              className="text-foreground data-[state=active]:bg-primary data-[state=active]:text-primary-foreground data-[state=active]:shadow-sm rounded-md transition-colors duration-200"
            >
              Dados Básicos
            </TabsTrigger>
            <TabsTrigger
              value="address"
              className="text-foreground data-[state=active]:bg-primary data-[state=active]:text-primary-foreground data-[state=active]:shadow-sm rounded-md transition-colors duration-200"
            >
              Endereço
            </TabsTrigger>
            <TabsTrigger
              value="tax"
              className="text-foreground data-[state=active]:bg-primary data-[state=active]:text-primary-foreground data-[state=active]:shadow-sm rounded-md transition-colors duration-200"
            >
              Informações Tributárias
            </TabsTrigger>
            <TabsTrigger
              value="reports"
              className="text-foreground data-[state=active]:bg-primary data-[state=active]:text-primary-foreground data-[state=active]:shadow-sm rounded-md transition-colors duration-200"
            >
              Relatórios
            </TabsTrigger>
          </TabsList>

          <TabsContent value="basic">
            <BasicInfoCustomerForm />
          </TabsContent>

          <TabsContent value="address">
            <div className="p-6 bg-form-background rounded-lg shadow-sm text-foreground">
              <h2 className="text-xl font-semibold mb-3">Endereço</h2>
              <p className="text-muted-foreground">
                O formulário para dados de endereço será implementado aqui.
                Utilize o `bg-form-background` para manter a consistência
                visual.
              </p>
            </div>
          </TabsContent>

          <TabsContent value="tax">
            <div className="p-6 bg-form-background rounded-lg shadow-sm text-foreground">
              <h2 className="text-xl font-semibold mb-3">
                Informações Tributárias
              </h2>
              <p className="text-muted-foreground">
                Este painel conterá as informações fiscais e tributárias do
                cliente.
              </p>
            </div>
          </TabsContent>

          <TabsContent value="reports">
            <div className="p-6 bg-form-background rounded-lg shadow-sm text-foreground">
              <h2 className="text-xl font-semibold mb-3">Relatórios</h2>
              <p className="text-muted-foreground">
                Aqui você poderá visualizar relatórios e histórico relacionados
                a este cliente.
              </p>
            </div>
          </TabsContent>
        </Tabs>
      </div>
    </PagePermissionGuard>
  );
}
