import api from "./common/api";

export interface FinancialTransaction {
  id: number;
  transaction_type: "receita" | "despesa";
  amount: number;
  description: string;
  reference_id?: number;
  reference_type?: string;
  transaction_date: string;
  payment_method_id: number;
  payment_method?: string;
  status: "pendente" | "pago" | "cancelado";
  due_date?: string;
  payment_date?: string;
  created_by: number;
  created_at: string;
  updated_at: string;
}

export interface CreateTransactionDto {
  transaction_type: "receita" | "despesa";
  amount: number;
  description: string;
  reference_id?: number;
  reference_type?: string;
  transaction_date: string;
  payment_method_id: number;
  status: "pendente" | "pago" | "cancelado";
  due_date?: string;
  payment_date?: string;
}

export interface UpdateTransactionDto {
  amount?: number;
  description?: string;
  transaction_date?: string;
  payment_method_id?: number;
  status?: "pendente" | "pago" | "cancelado";
  due_date?: string;
  payment_date?: string;
}

export interface FinancialSummary {
  totalRevenue: number;
  totalExpenses: number;
  balance: number;
  pendingRevenue: number;
  pendingExpenses: number;
}

export interface CashFlowItem {
  date: string;
  revenue: number;
  expenses: number;
  balance: number;
}

const FinancialService = {
  async getTransactions(
    page = 1,
    limit = 10,
    filters?: {
      type?: string;
      startDate?: string;
      endDate?: string;
      status?: string;
    }
  ): Promise<{ data: FinancialTransaction[]; total: number }> {
    let url = `/financial/transactions?page=${page}&limit=${limit}`;

    if (filters) {
      if (filters.type) url += `&type=${filters.type}`;
      if (filters.startDate) url += `&startDate=${filters.startDate}`;
      if (filters.endDate) url += `&endDate=${filters.endDate}`;
      if (filters.status) url += `&status=${filters.status}`;
    }

    const response = await api.get(url);
    return response.data;
  },

  async getTransactionById(id: number): Promise<FinancialTransaction> {
    const response = await api.get(`/financial/transactions/${id}`);
    return response.data;
  },

  async createTransaction(
    transaction: CreateTransactionDto
  ): Promise<FinancialTransaction> {
    const response = await api.post("/financial/transactions", transaction);
    return response.data;
  },

  async updateTransaction(
    id: number,
    transaction: UpdateTransactionDto
  ): Promise<FinancialTransaction> {
    const response = await api.put(
      `/financial/transactions/${id}`,
      transaction
    );
    return response.data;
  },

  async deleteTransaction(id: number): Promise<void> {
    await api.delete(`/financial/transactions/${id}`);
  },

  async getSummary(
    startDate?: string,
    endDate?: string
  ): Promise<FinancialSummary> {
    let url = "/financial/summary";
    if (startDate && endDate) {
      url += `?startDate=${startDate}&endDate=${endDate}`;
    }

    const response = await api.get(url);
    return response.data;
  },

  async getCashFlow(
    startDate?: string,
    endDate?: string
  ): Promise<CashFlowItem[]> {
    let url = "/financial/cash-flow";
    if (startDate && endDate) {
      url += `?startDate=${startDate}&endDate=${endDate}`;
    }

    const response = await api.get(url);
    return response.data;
  },
};

export default FinancialService;
