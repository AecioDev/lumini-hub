// src/hooks/permissions/use-permissions-pagination.ts (versão revisada novamente)
import { useState, useEffect, useCallback } from "react";
import { useToast } from "@/components/ui/use-toast";
import { Permission } from "@/services/auth/permission-schema";
import permissionService from "@/services/auth/permission-service";

interface UsePermissionsPaginationProps {
  initialPage?: number;
  pageSize?: number;
}

interface Filters {
  name?: string;
  module?: string;
  roleId?: string;
}

interface PermissionsPaginationResult {
  permissions: Permission[];
  isLoading: boolean;
  currentPage: number;
  totalPages: number;
  totalItems: number;
  searchName: string;
  setSearchName: (value: string) => void;
  searchModule: string;
  setSearchModule: (value: string) => void;
  searchRoleId: string;
  setSearchRoleId: (value: string) => void;
  applyFiltersAndFetch: () => void;
  goToNextPage: () => void;
  goToPreviousPage: () => void;
  goToFirstPage: () => void;
  goToLastPage: () => void;
  goToPage: (page: number) => void;
}

export function usePermissionsPagination({
  initialPage = 1,
  pageSize = 10,
}: UsePermissionsPaginationProps = {}): PermissionsPaginationResult {
  const [permissions, setPermissions] = useState<Permission[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [currentPage, setCurrentPage] = useState(initialPage);
  const [totalPages, setTotalPages] = useState(1);
  const [totalItems, setTotalItems] = useState(0);

  // Estados para os inputs dos filtros
  const [searchName, setSearchName] = useState("");
  const [searchModule, setSearchModule] = useState("");
  const [searchRoleId, setSearchRoleId] = useState("");

  // NOVO: Estado que REALMENTE dispara a busca de filtros
  // Inicializamos com os valores dos estados de input para a primeira carga.
  // IMPORTANTE: Para a primeira carga, se os inputs iniciarem vazios, o appliedFilters também deve.
  // Se quiser que ele busque com filtros iniciais definidos, mude o initialPage/pageSize OU
  // initialize os `search*` com valores padrão.
  const [appliedFilters, setAppliedFilters] = useState<Filters>({
    name: "", // Inicia vazio
    module: "",
    roleId: "",
  });

  const { toast } = useToast();

  // Função para buscar permissões do backend.
  // Ela pega os filtros do estado 'appliedFilters'.
  const fetchPermissionsData = useCallback(async () => {
    setIsLoading(true);
    console.log(
      `Fetching permissions for page ${currentPage} with filters: `,
      appliedFilters // Usa o estado `appliedFilters` que está nas dependências
    );
    try {
      const result = await permissionService.getPermissions(
        currentPage,
        pageSize,
        appliedFilters
      );

      setPermissions(result.data);
      setTotalPages(result.pagination.totalPages);
      setTotalItems(result.pagination.totalRows);
    } catch (error: any) {
      console.error("Erro ao carregar permissões:", error);
      toast({
        title: "Erro ao carregar dados",
        description:
          error.message || "Não foi possível carregar a lista de permissões.",
        variant: "destructive",
      });
      setPermissions([]);
      setTotalPages(1);
      setTotalItems(0);
    } finally {
      setIsLoading(false);
    }
  }, [currentPage, pageSize, appliedFilters, toast]); // Dependências: currentPage, pageSize, appliedFilters, toast

  // useEffect principal que dispara a busca.
  useEffect(() => {
    console.log("useEffect: Disparando busca de permissões...");
    fetchPermissionsData();
  }, [currentPage, appliedFilters, fetchPermissionsData]);

  // Função que o componente da tela chamará para aplicar os filtros
  // Esta função agora GARANTE que `appliedFilters` recebe os valores mais recentes
  const applyFiltersAndFetch = useCallback(() => {
    console.log("applyFiltersAndFetch chamado. Definindo filtros aplicados...");
    // Captura os valores mais recentes dos estados de input NO MOMENTO desta chamada
    const currentSearchName = searchName;
    const currentSearchModule = searchModule;
    const currentSearchRoleId = searchRoleId;

    setAppliedFilters({
      name: currentSearchName,
      module: currentSearchModule,
      roleId: currentSearchRoleId,
    });

    // Se a página for diferente de 1, setCurrentPage(1) causará um disparo do useEffect.
    // Se a página JÁ FOR 1, a mudança em `appliedFilters` ainda DISPARARÁ o useEffect.
    setCurrentPage(1); // Redefine para a primeira página ao aplicar novos filtros
  }, [searchName, searchModule, searchRoleId]); // Depende dos estados dos inputs

  // Funções de navegação de página
  const goToNextPage = () => {
    setCurrentPage((prev) => Math.min(prev + 1, totalPages));
  };
  const goToPreviousPage = () => {
    setCurrentPage((prev) => Math.max(prev - 1, 1));
  };
  const goToFirstPage = () => {
    setCurrentPage(1);
  };
  const goToLastPage = () => {
    setCurrentPage(totalPages);
  };
  const goToPage = (page: number) => {
    setCurrentPage(page);
  };

  return {
    permissions,
    isLoading,
    currentPage,
    totalPages,
    totalItems,
    searchName,
    setSearchName,
    searchModule,
    setSearchModule,
    searchRoleId,
    setSearchRoleId,
    applyFiltersAndFetch,
    goToNextPage,
    goToPreviousPage,
    goToFirstPage,
    goToLastPage,
    goToPage,
  };
}
