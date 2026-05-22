import { Customer } from "../customer/customer-service";
import { Supplier } from "../supplier/supplier-service";

export interface Contact {
  id: number;
  contact: string;
  type: string;
  name: string;
  customer_id: number;
  customer: Customer;
  supplier_id: number;
  supplier: Supplier;
}
