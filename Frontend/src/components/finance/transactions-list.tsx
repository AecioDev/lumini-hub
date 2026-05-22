"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useToast } from "@/components/ui/use-toast";
import { Edit, Eye, Loader2, MoreHorizontal, Trash2 } from "lucide-react";
import { format } from "date-fns";
import { formatCurrency } from "@/lib/utils";
import FinancialService, {
  type FinancialTransaction,
} from "@/services/financial-service";
import PaymentMethodService, {
  type PaymentMethod,
} from "@/services/payment-method-service";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";

interface TransactionsListProps {
  limit?: number;
  showFilters?: boolean;
  showPagination?: boolean;
}

export function FinancialTransactionsList({
  limit = 10,
  showFilters = false,
  showPagination = true,
}: TransactionsListProps) {
  const [transactions, setTransactions] = useState<FinancialTransaction[]>([]);
  const [paymentMethods, setPaymentMethods] = useState<PaymentMethod[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [filters, setFilters] = useState({
    type: "",
    status: "",
    startDate: "",
    endDate: "",
  });
  const [selectedTransaction, setSelectedTransaction] =
    useState<FinancialTransaction | null>(null);
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);

  const { toast } = useToast();
  const router = useRouter();

  // Carregar transações e métodos de pagamento
  useEffect(() => {
    const loadData = async () => {
      try {
        setIsLoading(true);

        const [transactionsResponse, methodsResponse] = await Promise.all([
          FinancialService.getTransactions(currentPage, limit, filters),
          PaymentMethodService.getPaymentMethods(),
        ]);

        setTransactions(transactionsResponse.data);
        setTotalPages(Math.ceil(transactionsResponse.total / limit));
        setPaymentMethods(methodsResponse);
      } catch (error) {
        console.error("Erro ao carregar transações:", error);
        toast({
          variant: "destructive",
          title: "Erro ao carregar dados",
          description: "Não foi possível carregar as transações financeiras.",
        });
      } finally {
        setIsLoading(false);
      }
    };

    loadData();
  }, [currentPage, limit, filters, toast]);

  const handleFilterChange = (key: string, value: string) => {
    setFilters((prev) => ({ ...prev, [key]: value }));
    setCurrentPage(1); // Resetar para a primeira página ao filtrar
  };

  const handleDeleteTransaction = async () => {
    if (!selectedTransaction) return;

    try {
      setIsDeleting(true);
      await FinancialService.deleteTransaction(selectedTransaction.id);

      setTransactions(
        transactions.filter((t) => t.id !== selectedTransaction.id)
      );
      setIsDeleteDialogOpen(false);
      setSelectedTransaction(null);

      toast({
        title: "Transação excluída",
        description: "A transação foi excluída com sucesso.",
      });
    } catch (error: any) {
      toast({
        variant: "destructive",
        title: "Erro ao excluir transação",
        description:
          error.response?.data?.message ||
          "Ocorreu um erro ao excluir a transação.",
      });
    } finally {
      setIsDeleting(false);
    }
  };

  const getPaymentMethodName = (id: number) => {
    const method = paymentMethods.find((m) => m.id === id);
    return method ? method.name : "-";
  };

  if (isLoading) {
    return (
      <div className="flex h-40 items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {showFilters && (
        <div className="grid grid-cols-1 gap-4 md:grid-cols-4">
          <div>
            <Select
              value={filters.type}
              onValueChange={(value) => handleFilterChange("type", value)}
            >
              <SelectTrigger>
                <SelectValue placeholder="Tipo de Transação" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">Todos os Tipos</SelectItem>
                <SelectItem value="receita">Receitas</SelectItem>
                <SelectItem value="despesa">Despesas</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div>
            <Select
              value={filters.status}
              onValueChange={(value) => handleFilterChange("status", value)}
            >
              <SelectTrigger>
                <SelectValue placeholder="Status" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">Todos os Status</SelectItem>
                <SelectItem value="pendente">Pendente</SelectItem>
                <SelectItem value="pago">Pago</SelectItem>
                <SelectItem value="cancelado">Cancelado</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div>
            <Input
              type="date"
              placeholder="Data Inicial"
              value={filters.startDate}
              onChange={(e) => handleFilterChange("startDate", e.target.value)}
            />
          </div>

          <div>
            <Input
              type="date"
              placeholder="Data Final"
              value={filters.endDate}
              onChange={(e) => handleFilterChange("endDate", e.target.value)}
            />
          </div>
        </div>
      )}

      <div className="rounded-md border">
        <table className="w-full">
          <thead>
            <tr className="border-b bg-muted/50 text-left text-sm font-medium">
              <th className="px-4 py-3">Descrição</th>
              <th className="px-4 py-3">Tipo</th>
              <th className="px-4 py-3">Valor</th>
              <th className="px-4 py-3">Data</th>
              <th className="px-4 py-3">Método</th>
              <th className="px-4 py-3">Status</th>
              <th className="px-4 py-3 text-right">Ações</th>
            </tr>
          </thead>
          <tbody>
            {transactions.length === 0 ? (
              <tr>
                <td
                  colSpan={7}
                  className="px-4 py-8 text-center text-muted-foreground"
                >
                  Nenhuma transação encontrada
                </td>
              </tr>
            ) : (
              transactions.map((transaction) => (
                <tr key={transaction.id} className="border-b">
                  <td className="px-4 py-3">{transaction.description}</td>
                  <td className="px-4 py-3">
                    <span
                      className={`rounded-full px-2 py-1 text-xs font-medium ${
                        transaction.transaction_type === "receita"
                          ? "bg-green-100 text-green-700"
                          : "bg-red-100 text-red-700"
                      }`}
                    >
                      {transaction.transaction_type === "receita"
                        ? "Receita"
                        : "Despesa"}
                    </span>
                  </td>
                  <td className="px-4 py-3 font-medium">
                    {formatCurrency(transaction.amount)}
                  </td>
                  <td className="px-4 py-3">
                    {format(
                      new Date(transaction.transaction_date),
                      "dd/MM/yyyy"
                    )}
                  </td>
                  <td className="px-4 py-3">
                    {getPaymentMethodName(transaction.payment_method_id)}
                  </td>
                  <td className="px-4 py-3">
                    <span
                      className={`rounded-full px-2 py-1 text-xs font-medium ${
                        transaction.status === "pago"
                          ? "bg-green-100 text-green-700"
                          : transaction.status === "pendente"
                          ? "bg-amber-100 text-amber-700"
                          : "bg-red-100 text-red-700"
                      }`}
                    >
                      {transaction.status === "pago"
                        ? "Pago"
                        : transaction.status === "pendente"
                        ? "Pendente"
                        : "Cancelado"}
                    </span>
                  </td>
                  <td className="px-4 py-3 text-right">
                    <DropdownMenu>
                      <DropdownMenuTrigger asChild>
                        <Button variant="ghost" size="icon">
                          <MoreHorizontal className="h-4 w-4" />
                        </Button>
                      </DropdownMenuTrigger>
                      <DropdownMenuContent align="end">
                        <DropdownMenuLabel>Ações</DropdownMenuLabel>
                        <DropdownMenuSeparator />
                        <DropdownMenuItem
                          onClick={() =>
                            router.push(
                              `/dashboard/finance/transactions/${transaction.id}`
                            )
                          }
                        >
                          <Eye className="mr-2 h-4 w-4" />
                          Visualizar
                        </DropdownMenuItem>
                        <DropdownMenuItem
                          onClick={() =>
                            router.push(
                              `/dashboard/finance/transactions/${transaction.id}/edit`
                            )
                          }
                        >
                          <Edit className="mr-2 h-4 w-4" />
                          Editar
                        </DropdownMenuItem>
                        <DropdownMenuItem
                          onClick={() => {
                            setSelectedTransaction(transaction);
                            setIsDeleteDialogOpen(true);
                          }}
                          className="text-destructive"
                        >
                          <Trash2 className="mr-2 h-4 w-4" />
                          Excluir
                        </DropdownMenuItem>
                      </DropdownMenuContent>
                    </DropdownMenu>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>

      {showPagination && totalPages > 1 && (
        <div className="flex items-center justify-end space-x-2">
          <Button
            variant="outline"
            size="sm"
            onClick={() => setCurrentPage((prev) => Math.max(prev - 1, 1))}
            disabled={currentPage === 1}
          >
            Anterior
          </Button>
          <span className="text-sm text-muted-foreground">
            Página {currentPage} de {totalPages}
          </span>
          <Button
            variant="outline"
            size="sm"
            onClick={() =>
              setCurrentPage((prev) => Math.min(prev + 1, totalPages))
            }
            disabled={currentPage === totalPages}
          >
            Próxima
          </Button>
        </div>
      )}

      {/* Dialog de confirmação de exclusão */}
      <Dialog open={isDeleteDialogOpen} onOpenChange={setIsDeleteDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Confirmar Exclusão</DialogTitle>
            <DialogDescription>
              Tem certeza que deseja excluir esta transação? Esta ação não pode
              ser desfeita.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => setIsDeleteDialogOpen(false)}
              disabled={isDeleting}
            >
              Cancelar
            </Button>
            <Button
              variant="destructive"
              onClick={handleDeleteTransaction}
              disabled={isDeleting}
            >
              {isDeleting ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  Excluindo...
                </>
              ) : (
                "Excluir"
              )}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
