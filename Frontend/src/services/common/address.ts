import { Customer } from "../customer/customer-service";
import { Supplier } from "../supplier/supplier-service";

export interface ddAddress {
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
