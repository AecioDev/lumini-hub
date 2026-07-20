// Espelha Backend/microservices/api.core/internal/domain/{address,contact}.go

export interface City {
  id: number;
  name: string;
  ibge_code: string;
  state_id: number;
}

export interface Address {
  id: number;
  street: string;
  number: string;
  neighborhood: string;
  zip_code: string;
  city_id: number;
  city: City;
  customer_id?: number;
  supplier_id?: number;
}

export interface Contact {
  id: number;
  customer_id?: number;
  supplier_id?: number;
  type: string; // email | Telefone | Celular
  contact: string;
  name?: string;
}
