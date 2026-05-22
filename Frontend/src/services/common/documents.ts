import { Customer } from "../customer/customer-service";
import { Supplier } from "../supplier/supplier-service";
import { State } from "./address";

export interface Document {
  id: number;
  type: string;
  number: string;
  validate: string;
  emission_date: string;
  department: string;
  state_id: number;
  state: State;
  customer_id: number;
  customer: Customer;
  supplier_id: number;
  supplier: Supplier;
}
