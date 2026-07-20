// Espelha o envelope padrão de resposta do backend (Backend/common/utils/response.go)
export interface ApiResponse<T> {
  success: boolean;
  statusCode: number;
  message?: string;
  validationErrors?: string[];
  data?: T;
  error?: string;
  meta?: unknown;
}

// Espelha utils.ApiPagination (Backend/common/utils/pagination.go)
export interface ApiPagination {
  page: number;
  limit: number;
  sort: string;
  order: string;
  totalRows: number;
  totalPages: number;
}
