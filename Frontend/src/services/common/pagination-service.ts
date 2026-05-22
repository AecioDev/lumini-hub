export interface Pagination {
  page: number;
  limit: number;
  sort: string;
  order: string;
  totalRows: number;
  totalPages: number;
}
