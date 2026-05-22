"use client";

import { useEffect, useState } from "react";
import {
  Bar,
  BarChart,
  CartesianGrid,
  Legend,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";
import type { CashFlowItem } from "@/services/financial-service";
import { formatCurrency } from "@/lib/utils";

interface CashFlowChartProps {
  data: CashFlowItem[];
}

export function CashFlowChart({ data }: CashFlowChartProps) {
  const [chartData, setChartData] = useState<any[]>([]);

  useEffect(() => {
    // Formatar dados para o gráfico
    const formattedData = data.map((item) => ({
      date: new Date(item.date).toLocaleDateString("pt-BR"),
      receitas: item.revenue,
      despesas: item.expenses,
      saldo: item.balance,
    }));

    setChartData(formattedData);
  }, [data]);

  return (
    <div className="h-[400px] w-full">
      <ResponsiveContainer width="100%" height="100%">
        <BarChart
          data={chartData}
          margin={{
            top: 20,
            right: 30,
            left: 20,
            bottom: 5,
          }}
        >
          <CartesianGrid strokeDasharray="3 3" />
          <XAxis dataKey="date" />
          <YAxis
            tickFormatter={(value) => formatCurrency(value).replace("R$", "")}
          />
          <Tooltip
            formatter={(value) => formatCurrency(Number(value))}
            labelFormatter={(label) => `Data: ${label}`}
          />
          <Legend />
          <Bar dataKey="receitas" name="Receitas" fill="#4ade80" />
          <Bar dataKey="despesas" name="Despesas" fill="#f87171" />
        </BarChart>
      </ResponsiveContainer>
    </div>
  );
}
