// src/hooks/useCustomersPagination.ts
import { useState, useEffect, useCallback } from "react";
import { useToast } from "@/components/ui/use-toast";
import CustomerService, {
  type Customer,
} from "@/services/customer/customer-service";

interface UseCustomersPaginationProps {
  initialPage?: number;
  pageSize?: number;
}

export function useCustomersPagination({
  initialPage = 1,
  pageSize = 10,
}: UseCustomersPaginationProps = {}) {
  const [customers, setCustomers] = useState<Customer[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [currentPage, setCurrentPage] = useState(initialPage);
  const [totalPages, setTotalPages] = useState(1);
  const [searchName, setSearchName] = useState("");
  const [searchDocument, setSearchDocument] = useState("");

  const { toast } = useToast();

  const fetchCustomers = useCallback(async () => {
    try {
      setIsLoading(true);
      const filters: { name?: string; document?: string } = {};

      if (searchName) filters.name = searchName;
      if (searchDocument) filters.document = searchDocument;

      const response = await CustomerService.getCustomers(
        currentPage,
        pageSize,
        filters
      );
      setCustomers(response.data);
      setTotalPages(Math.ceil(response.pagination.totalPages / pageSize));
    } catch (error) {
      console.error("Erro ao carregar clientes:", error);
      toast({
        variant: "destructive",
        title: "Erro ao carregar dados",
        description: "Não foi possível carregar a lista de clientes.",
      });
    } finally {
      setIsLoading(false);
    }
  }, [currentPage, pageSize, searchName, searchDocument, toast]);

  useEffect(() => {
    fetchCustomers();
  }, [fetchCustomers]);

  const handleSearch = () => setCurrentPage(1);
  const goToNextPage = () =>
    setCurrentPage((prev) => Math.min(prev + 1, totalPages));
  const goToPreviousPage = () =>
    setCurrentPage((prev) => Math.max(prev - 1, 1));

  return {
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
    fetchCustomers, // Expondo para que possa ser chamado após uma exclusão, por exemplo
  };
}
