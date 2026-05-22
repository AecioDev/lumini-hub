import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { ArrowDownIcon, ArrowUpIcon, DollarSign, Wallet } from "lucide-react";
import { formatCurrency } from "@/lib/utils";
import type { FinancialSummary } from "@/services/financial-service";

interface FinancialSummaryCardsProps {
  summary: FinancialSummary;
}

export function FinancialSummaryCards({ summary }: FinancialSummaryCardsProps) {
  return (
    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
      <Card>
        <CardHeader className="flex flex-row items-center justify-between pb-2">
          <CardTitle className="text-sm font-medium">Receitas Totais</CardTitle>
          <ArrowUpIcon className="h-4 w-4 text-green-500" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">
            {formatCurrency(summary.totalRevenue)}
          </div>
          <p className="text-xs text-muted-foreground">
            Receitas no período selecionado
          </p>
        </CardContent>
      </Card>
      <Card>
        <CardHeader className="flex flex-row items-center justify-between pb-2">
          <CardTitle className="text-sm font-medium">Despesas Totais</CardTitle>
          <ArrowDownIcon className="h-4 w-4 text-red-500" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">
            {formatCurrency(summary.totalExpenses)}
          </div>
          <p className="text-xs text-muted-foreground">
            Despesas no período selecionado
          </p>
        </CardContent>
      </Card>
      <Card>
        <CardHeader className="flex flex-row items-center justify-between pb-2">
          <CardTitle className="text-sm font-medium">Saldo</CardTitle>
          <Wallet className="h-4 w-4 text-primary" />
        </CardHeader>
        <CardContent>
          <div
            className={`text-2xl font-bold ${
              summary.balance >= 0 ? "text-green-600" : "text-red-600"
            }`}
          >
            {formatCurrency(summary.balance)}
          </div>
          <p className="text-xs text-muted-foreground">
            Saldo no período selecionado
          </p>
        </CardContent>
      </Card>
      <Card>
        <CardHeader className="flex flex-row items-center justify-between pb-2">
          <CardTitle className="text-sm font-medium">Pendentes</CardTitle>
          <DollarSign className="h-4 w-4 text-amber-500" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">
            {formatCurrency(summary.pendingRevenue - summary.pendingExpenses)}
          </div>
          <p className="text-xs text-muted-foreground">
            Saldo de transações pendentes
          </p>
        </CardContent>
      </Card>
    </div>
  );
}
