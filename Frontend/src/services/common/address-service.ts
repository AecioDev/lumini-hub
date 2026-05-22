import api from "./api";
import { Customer } from "../customer/customer-service";
import { Supplier } from "../supplier/supplier-service";
import { City } from "./city-service";

export interface Address {
  id: number;
  street: string;
  number: string;
  neighborhood: string;
  zip_code: string;
  city_id: number;
  city: City;
  customer_id: number;
  customer: Customer;
  supplier_id: number;
  supplier: Supplier;
}

export const AddressService = {
  async getAddressesByCustomer(
    customerId: number,
    page = 1,
    limit = 10
  ): Promise<Address[]> {
    const response = await api.get(
      `/customers/${customerId}/addresses?page=${page}&limit=${limit}`
    );
    return response.data.data;
  },

  async getAddressById(id: number): Promise<Address> {
    const response = await api.get(`/addresses/${id}`);
    return response.data;
  },

  async createAddress(data: Partial<Address>): Promise<Address> {
    const response = await api.post("/addresses", data);
    return response.data;
  },

  async updateAddress(id: number, data: Partial<Address>): Promise<Address> {
    const response = await api.put(`/addresses/${id}`, data);
    return response.data;
  },

  async deleteAddress(id: number): Promise<void> {
    await api.delete(`/addresses/${id}`);
  },
};
