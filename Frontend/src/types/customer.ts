import type { ApiPagination } from "./common";
import type { Address, Contact } from "./address-contact";

// Espelha Backend/microservices/api.core/internal/domain/customer.go

export type PersonType = "F" | "J";

export interface ApiCustomer {
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

export interface CreateCustomerRequest {
  first_name: string;
  last_name?: string;
  person_type: PersonType;
  document_number: string;
  company_name?: string;
  is_active: boolean;
  notes?: string;
}

export interface UpdateCustomerRequest {
  first_name: string;
  last_name?: string;
  document_number: string;
  company_name?: string;
  is_active: boolean;
  notes?: string;
}

export interface CustomerFilterRequest {
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

export interface ApiCustomerListPaginated {
  data: ApiCustomer[];
  pagination: ApiPagination;
}
