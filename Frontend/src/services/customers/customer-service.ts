import { api } from "../common/api";
import type { ApiResponse } from "@/types/common";
import type {
  ApiCustomer,
  ApiCustomerListPaginated,
  CreateCustomerRequest,
  CustomerFilterRequest,
  UpdateCustomerRequest,
} from "@/types/customer";

export interface ListCustomersParams {
  page?: number;
  limit?: number;
  sort?: string;
}

export const customerService = {
  async list(params: ListCustomersParams = {}): Promise<ApiCustomerListPaginated> {
    const { data } = await api.get<ApiResponse<ApiCustomerListPaginated>>(
      "/customers",
      { params }
    );
    return data.data!;
  },

  async filter(payload: CustomerFilterRequest): Promise<ApiCustomerListPaginated> {
    const { data } = await api.post<ApiResponse<ApiCustomerListPaginated>>(
      "/customers/filter",
      payload
    );
    return data.data!;
  },

  async getById(id: number): Promise<ApiCustomer> {
    const { data } = await api.get<ApiResponse<ApiCustomer>>(`/customers/${id}`);
    return data.data!;
  },

  async create(payload: CreateCustomerRequest): Promise<ApiCustomer> {
    const { data } = await api.post<ApiResponse<ApiCustomer>>(
      "/customers",
      payload
    );
    return data.data!;
  },

  async update(id: number, payload: UpdateCustomerRequest): Promise<ApiCustomer> {
    const { data } = await api.put<ApiResponse<ApiCustomer>>(
      `/customers/${id}`,
      payload
    );
    return data.data!;
  },

  async remove(id: number): Promise<void> {
    await api.delete(`/customers/${id}`);
  },
};
