import { api } from "../common/api";
import type { ApiResponse } from "@/types/common";
import type {
  ApiSupplier,
  ApiSupplierListPaginated,
  CreateSupplierRequest,
  SupplierFilterRequest,
  UpdateSupplierRequest,
} from "@/types/supplier";

export interface ListSuppliersParams {
  page?: number;
  limit?: number;
  sort?: string;
}

export const supplierService = {
  async list(params: ListSuppliersParams = {}): Promise<ApiSupplierListPaginated> {
    const { data } = await api.get<ApiResponse<ApiSupplierListPaginated>>(
      "/suppliers",
      { params }
    );
    return data.data!;
  },

  async filter(payload: SupplierFilterRequest): Promise<ApiSupplierListPaginated> {
    const { data } = await api.post<ApiResponse<ApiSupplierListPaginated>>(
      "/suppliers/filter",
      payload
    );
    return data.data!;
  },

  async getById(id: number): Promise<ApiSupplier> {
    const { data } = await api.get<ApiResponse<ApiSupplier>>(`/suppliers/${id}`);
    return data.data!;
  },

  async create(payload: CreateSupplierRequest): Promise<ApiSupplier> {
    const { data } = await api.post<ApiResponse<ApiSupplier>>(
      "/suppliers",
      payload
    );
    return data.data!;
  },

  async update(id: number, payload: UpdateSupplierRequest): Promise<ApiSupplier> {
    const { data } = await api.put<ApiResponse<ApiSupplier>>(
      `/suppliers/${id}`,
      payload
    );
    return data.data!;
  },

  async remove(id: number): Promise<void> {
    await api.delete(`/suppliers/${id}`);
  },
};
