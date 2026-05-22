import api from "../common/api";
import { Address } from "../common/address-service";
import { Contact } from "../common/contacts";
import { Pagination } from "../common/pagination-service";
import { z } from "zod";
import { User } from "../auth/user-schema";

export interface Customer {
  id: number;
  first_name: string;
  last_name: string;
  person_type: "fisica" | "juridica" | string;
  document_number: string;
  company_name: string;
  is_active: boolean;

  addresses?: Address[] | null;
  contacts?: Contact[] | null;

  created_at: string; // backend envia time.Time, mas no frontend geralmente string ISO
  updated_at: string;
}

export interface CustomerDetail extends Customer {
  created_by?: User | null;
  updated_by?: User | null;

  //sales?: Sale[] | null;
  //transactions?: Transaction[] | null;
}

export interface CustomerList {
  data: Customer[];
  pagination: Pagination; // defina Pagination conforme seu DTO PaginationDTO do backend
}

export const customerFormSchema = z
  .object({
    first_name: z
      .string()
      .min(1, "Nome é obrigatório")
      .transform((str) => str.trim()),
    last_name: z
      .string()
      .optional()
      .transform((str) => str?.trim() || ""),
    person_type: z.enum(["fisica", "juridica"]),
    document_number: z
      .string()
      .min(1, "Número do documento é obrigatório")
      .transform((str) => str.trim()),
    company_name: z
      .string()
      .optional()
      .transform((str) => str?.trim()),
    is_active: z.boolean().default(true),
    notes: z
      .string()
      .optional()
      .transform((str) => str?.trim() || ""),
  })
  .refine(
    (data) => (data.person_type === "juridica" ? !!data.company_name : true),
    {
      message: "Razão social é obrigatória para pessoa jurídica",
      path: ["company_name"],
    }
  );

export type CreateCustomerInput = z.infer<typeof customerFormSchema>;

export interface UpdateCustomerDto extends CreateCustomerInput {}

const CustomerService = {
  async getCustomers(
    page = 1,
    limit = 10,
    filters?: {
      name?: string;
      document?: string;
      is_active?: boolean;
    }
  ): Promise<CustomerList> {
    let url = `/customers?page=${page}&limit=${limit}`;

    if (filters) {
      if (filters.name) url += `&name=${filters.name}`;
      if (filters.document) url += `&document=${filters.document}`;
      if (filters.is_active !== undefined)
        url += `&is_active=${filters.is_active}`;
    }

    const response = await api.get(url);
    return response.data.data;
  },

  async getCustomerById(id: number): Promise<Customer> {
    const response = await api.get(`/customers/${id}`);
    return response.data;
  },

  async createCustomer(customer: CreateCustomerInput): Promise<Customer> {
    const response = await api.post("/customers", customer);
    return response.data;
  },

  async updateCustomer(
    id: number,
    customer: CreateCustomerInput
  ): Promise<Customer> {
    const response = await api.put(`/customers/${id}`, customer);
    return response.data;
  },

  async deleteCustomer(id: number): Promise<void> {
    await api.delete(`/customers/${id}`);
  },
};

export default CustomerService;
