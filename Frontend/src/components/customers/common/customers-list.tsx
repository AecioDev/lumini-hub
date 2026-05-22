// src/components/customers/customers-list.tsx
"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Loader2 } from "lucide-react";

import CustomerService, {
  type Customer,
} from "@/services/customer/customer-service";
import { DeleteCustomerDialog } from "@/components/customers/common/delete-customer-dialog"; // Certifique-se do caminho correto
import { useCustomersPagination } from "@/hooks/customers/use-customers-pagination"; // Importa o novo hook
import { CustomersTable } from "@/components/customers/tables/customers-table"; // Importa o novo componente de tabela
import { useToast } from "@/components/ui/use-toast";

export function CustomersList() {
  const {
    customers,
    isLoading,
    currentPage,
    totalPages,
    searchName,
    setSearchName,
    searchDocument,
    setSearchDocument,
    handleSearch,
    goToNextPage,
    goToPreviousPage,
    fetchCustomers, // Usado para recarregar após exclusão
  } = useCustomersPagination();

  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
  const [selectedCustomer, setSelectedCustomer] = useState<Customer | null>(
    null
  );
  const { toast } = useToast();

  const confirmDelete = (customer: Customer) => {
    setSelectedCustomer(customer);
    setIsDeleteDialogOpen(true);
  };

  const handleDeleteCustomer = async () => {
    if (!selectedCustomer) return;

    try {
      await CustomerService.deleteCustomer(selectedCustomer.id);
      toast({
        title: "Cliente excluído",
        description: `O cliente foi excluído com sucesso`,
      });
      fetchCustomers(); // Recarrega a lista após a exclusão
    } catch (error: any) {
      toast({
        variant: "destructive",
        title: "Erro ao excluir cliente",
        description:
          error.response?.data?.message || "Erro ao excluir cliente.",
      });
    } finally {
      setIsDeleteDialogOpen(false);
      setSelectedCustomer(null);
    }
  };

  if (isLoading) {
    return (
      <div className="flex h-40 items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    );
  }

  return (
    <>
      {/* Filtros */}
      <div className="mb-4 grid grid-cols-1 gap-4 md:grid-cols-3">
        <div className="grid gap-2">
          <Label htmlFor="search-name">Buscar por nome</Label>
          <Input
            id="search-name"
            value={searchName}
            onChange={(e) => setSearchName(e.target.value)}
            placeholder="Nome do cliente"
          />
        </div>
        <div className="grid gap-2">
          <Label htmlFor="search-document">Buscar por documento</Label>
          <Input
            id="search-document"
            value={searchDocument}
            onChange={(e) => setSearchDocument(e.target.value)}
            placeholder="CPF ou CNPJ"
          />
        </div>
        <div className="flex items-end">
          <Button onClick={handleSearch} className="mb-0.5">
            Buscar
          </Button>
        </div>
      </div>

      {/* Tabela de Clientes */}
      <CustomersTable customers={customers} onConfirmDelete={confirmDelete} />

      {/* Paginação */}
      {totalPages > 1 && (
        <div className="mt-4 flex items-center justify-end gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={goToPreviousPage}
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
            onClick={goToNextPage}
            disabled={currentPage === totalPages}
          >
            Próxima
          </Button>
        </div>
      )}

      {/* Diálogo de exclusão */}
      <DeleteCustomerDialog
        open={isDeleteDialogOpen}
        onOpenChange={setIsDeleteDialogOpen}
        customerName={selectedCustomer?.first_name || ""}
        onConfirm={handleDeleteCustomer}
      />
    </>
  );
}
