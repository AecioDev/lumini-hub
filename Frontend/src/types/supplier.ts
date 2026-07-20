import type { ApiPagination } from "./common";
import type { Address, Contact } from "./address-contact";
import type { PersonType } from "./customer";

// Espelha Backend/microservices/api.core/internal/domain/supplier.go

export interface ApiSupplier {
  id: number;
  first_name: string;
  last_name: string;
  person_type: PersonType;
  document_number: string;
  company_name: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
  addresses?: Address[];
  contacts?: Contact[];
}

export interface CreateSupplierRequest {
  first_name: string;
  last_name?: string;
  person_type: PersonType;
  document_number: string;
  company_name?: string;
  is_active: boolean;
  notes?: string;
}

export interface UpdateSupplierRequest {
  first_name: string;
  last_name?: string;
  document_number: string;
  company_name?: string;
  is_active: boolean;
  notes?: string;
}

export interface SupplierFilterRequest {
  id?: number;
  first_name?: string;
  last_name?: string;
  person_type?: string;
  document_number?: string;
  company_name?: string;
  is_active?: boolean;
  page_no?: number;
  page_size?: number;
  order_by_column?: string;
  is_asc?: boolean;
}

export interface ApiSupplierListPaginated {
  data: ApiSupplier[];
  pagination: ApiPagination;
}
