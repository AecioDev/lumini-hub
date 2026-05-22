"use client";

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  BarChart3,
  ShoppingCart,
  Package,
  Wallet,
  TrendingUp,
} from "lucide-react";
import { PagePermissionGuard } from "@/components/layout/PagePermissionGuard";

export default function DashboardPage() {
  return (
    <PagePermissionGuard
      requiredPermissions={["dashboard.manager.view"]}
      accessDeniedMessage="Você não tem permissão para visualizar este Dashboard."
    >
      <div className="space-y-6">
        <div>
          <h1 className="text-3xl font-bold">Dashboard Gerencial</h1>
          <p className="text-muted-foreground">Visão gerencial do sistema</p>
        </div>

        <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between pb-2">
              <CardTitle className="text-sm font-medium">
                Vendas Totais
              </CardTitle>
              <ShoppingCart className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">R$ 45.231,89</div>
              <p className="text-xs text-muted-foreground">
                +20.1% em relação ao mês anterior
              </p>
            </CardContent>
          </Card>
          <Card>
            <CardHeader className="flex flex-row items-center justify-between pb-2">
              <CardTitle className="text-sm font-medium">
                Produtos em Estoque
              </CardTitle>
              <Package className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">1.234</div>
              <p className="text-xs text-muted-foreground">
                +12 novos produtos este mês
              </p>
            </CardContent>
          </Card>
          <Card>
            <CardHeader className="flex flex-row items-center justify-between pb-2">
              <CardTitle className="text-sm font-medium">
                Receita Mensal
              </CardTitle>
              <Wallet className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">R$ 12.543,00</div>
              <p className="text-xs text-muted-foreground">
                +2.5% em relação ao mês anterior
              </p>
            </CardContent>
          </Card>
          <Card>
            <CardHeader className="flex flex-row items-center justify-between pb-2">
              <CardTitle className="text-sm font-medium">Crescimento</CardTitle>
              <TrendingUp className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">+18.2%</div>
              <p className="text-xs text-muted-foreground">Crescimento anual</p>
            </CardContent>
          </Card>
        </div>

        <div className="grid gap-6 md:grid-cols-2">
          <Card className="col-span-1">
            <CardHeader>
              <CardTitle>Vendas Recentes</CardTitle>
              <CardDescription>
                Últimas 5 vendas realizadas no sistema
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {[1, 2, 3, 4, 5].map((i) => (
                  <div
                    key={i}
                    className="flex items-center justify-between border-b pb-2"
                  >
                    <div>
                      <p className="font-medium">Pedido #{1000 + i}</p>
                      <p className="text-sm text-muted-foreground">
                        Cliente {i}
                      </p>
                    </div>
                    <div className="text-right">
                      <p className="font-medium">
                        R$ {(Math.random() * 1000).toFixed(2)}
                      </p>
                      <p className="text-sm text-muted-foreground">Hoje</p>
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
          <Card className="col-span-1">
            <CardHeader>
              <CardTitle>Atividades Recentes</CardTitle>
              <CardDescription>Últimas atividades no sistema</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {[
                  "Usuário adicionou um novo produto",
                  "Venda #1005 foi finalizada",
                  "Estoque do produto X foi atualizado",
                  "Novo usuário foi cadastrado",
                  "Relatório mensal foi gerado",
                ].map((activity, i) => (
                  <div
                    key={i}
                    className="flex items-center gap-4 border-b pb-2"
                  >
                    <div className="rounded-full bg-primary/10 p-2">
                      <BarChart3 className="h-4 w-4 text-primary" />
                    </div>
                    <div>
                      <p className="font-medium">{activity}</p>
                      <p className="text-sm text-muted-foreground">
                        Há {i + 1} horas
                      </p>
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </PagePermissionGuard>
  );
}
