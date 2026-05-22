import api from "../common/api";
import { User } from "../user-service";

export interface Supplier {
  id: number;
  first_name: string;
  last_name: string;
  person_type: string;
  DocumentNumber: string;
  company_name: string;
  is_active: boolean;
  notes: string;

  created_by: number;
  UserCreate: User;
  updated_by: number;
  UserUpdate: User;

  created_at: string;
  updated_at: string;
}

export interface CreateSupplierDto {
  name: string;
  email?: string;
  phone?: string;
  document?: string;
  document_type?: "cpf" | "cnpj";
  address?: string;
  city?: string;
  state?: string;
  postal_code?: string;
  notes?: string;
}

export interface UpdateSupplierDto {
  name?: string;
  email?: string;
  phone?: string;
  document?: string;
  document_type?: "cpf" | "cnpj";
  address?: string;
  city?: string;
  state?: string;
  postal_code?: string;
  notes?: string;
  is_active?: boolean;
}

const SupplierService = {
  async getSuppliers(
    page = 1,
    limit = 10,
    filters?: {
      name?: string;
      document?: string;
      is_active?: boolean;
    }
  ): Promise<{ data: Supplier[]; total: number }> {
    let url = `/suppliers?page=${page}&limit=${limit}`;

    if (filters) {
      if (filters.name) url += `&name=${filters.name}`;
      if (filters.document) url += `&document=${filters.document}`;
      if (filters.is_active !== undefined)
        url += `&is_active=${filters.is_active}`;
    }

    const response = await api.get(url);
    return response.data;
  },

  async getSupplierById(id: number): Promise<Supplier> {
    const response = await api.get(`/suppliers/${id}`);
    return response.data;
  },

  async createSupplier(supplier: CreateSupplierDto): Promise<Supplier> {
    const response = await api.post("/suppliers", supplier);
    return response.data;
  },

  async updateSupplier(
    id: number,
    supplier: UpdateSupplierDto
  ): Promise<Supplier> {
    const response = await api.put(`/suppliers/${id}`, supplier);
    return response.data;
  },

  async deleteSupplier(id: number): Promise<void> {
    await api.delete(`/suppliers/${id}`);
  },
};

export default SupplierService;
