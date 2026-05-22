// src/hooks/use-server-pagination.ts
import { useState, useEffect, useCallback } from "react";
import { useToast } from "@/components/ui/use-toast";

// Definindo um tipo genérico para os filtros e para os dados
interface Filters {
  [key: string]: any; // Filtros podem ser qualquer chave-valor
}

interface PaginationResult<T> {
  data: T[]; // 'T' é o tipo da entidade (Permission, Product, Customer)
  pagination: {
    totalPages: number;
    totalRows: number;
  };
}

// A interface para a função de busca que será injetada
interface FetchFunction<T> {
  (page: number, limit: number, filters: Filters): Promise<PaginationResult<T>>;
}

interface UseServerPaginationProps<T> {
  fetchDataService: FetchFunction<T>; // <--- A função de busca do serviço específico
  initialPage?: number;
  pageSize?: number;
}

interface UseServerPaginationResult<T> {
  data: T[];
  isLoading: boolean;
  currentPage: number;
  totalPages: number;
  totalItems: number;
  fetchDataWithFilters: (filters: Filters, resetPage?: boolean) => void;
  goToNextPage: () => void;
  goToPreviousPage: () => void;
  goToFirstPage: () => void;
  goToLastPage: () => void;
  goToPage: (page: number) => void;
}

export function useServerPagination<T>({
  fetchDataService,
  initialPage = 1,
  pageSize = 10,
}: UseServerPaginationProps<T>): UseServerPaginationResult<T> {
  const [data, setData] = useState<T[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [currentPage, setCurrentPage] = useState(initialPage);
  const [totalPages, setTotalPages] = useState(1);
  const [totalItems, setTotalItems] = useState(0);

  const [appliedFilters, setAppliedFilters] = useState<Filters>({});

  const { toast } = useToast();

  const fetchData = useCallback(async () => {
    setIsLoading(true);
    // console.log(
    //   `Fetching data for page ${currentPage} with filters: `,
    //   appliedFilters
    // );
    try {
      // Chama a função de busca que foi passada como prop
      const result = await fetchDataService(
        currentPage,
        pageSize,
        appliedFilters
      );

      setData(result.data);
      setTotalPages(result.pagination.totalPages);
      setTotalItems(result.pagination.totalRows);
    } catch (error: any) {
      console.error("Erro ao carregar dados:", error);
      toast({
        title: "Erro ao carregar dados",
        description:
          error.message || "Não foi possível carregar a lista de itens.",
        variant: "destructive",
      });
      setData([]);
      setTotalPages(1);
      setTotalItems(0);
    } finally {
      setIsLoading(false);
    }
  }, [currentPage, pageSize, appliedFilters, fetchDataService, toast]);

  useEffect(() => {
    fetchData();
  }, [currentPage, appliedFilters, fetchData]);

  const fetchDataWithFilters = useCallback(
    (filters: Filters, resetPage: boolean = true) => {
      setAppliedFilters(filters);
      if (resetPage) {
        setCurrentPage(1);
      }
    },
    []
  );

  const goToNextPage = () =>
    setCurrentPage((prev) => Math.min(prev + 1, totalPages));
  const goToPreviousPage = () =>
    setCurrentPage((prev) => Math.max(prev - 1, 1));
  const goToFirstPage = () => setCurrentPage(1);
  const goToLastPage = () => setCurrentPage(totalPages);
  const goToPage = (page: number) => setCurrentPage(page);

  return {
    data,
    isLoading,
    currentPage,
    totalPages,
    totalItems,
    fetchDataWithFilters,
    goToNextPage,
    goToPreviousPage,
    goToFirstPage,
    goToLastPage,
    goToPage,
  };
}
