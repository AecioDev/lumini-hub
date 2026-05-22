import api from "./common/api";

export interface Product {
  id: number;
  sku: string;
  barcode?: string;
  name: string;
  description?: string;
  category_id: number;
  category_name?: string;
  unit_id: number;
  unit_name?: string;
  cost_price: number;
  selling_price: number;
  min_stock: number;
  max_stock?: number;
  current_stock: number;
  is_active: boolean;
  created_by: number;
  created_at: string;
  updated_at: string;
}

export interface CreateProductDto {
  sku: string;
  barcode?: string;
  name: string;
  description?: string;
  category_id: number;
  unit_id: number;
  cost_price: number;
  selling_price: number;
  min_stock: number;
  max_stock?: number;
  current_stock: number;
}

export interface UpdateProductDto {
  sku?: string;
  barcode?: string;
  name?: string;
  description?: string;
  category_id?: number;
  unit_id?: number;
  cost_price?: number;
  selling_price?: number;
  min_stock?: number;
  max_stock?: number;
  is_active?: boolean;
}

const ProductService = {
  async getProducts(
    page = 1,
    limit = 10,
    filters?: {
      name?: string;
      category_id?: number;
      is_active?: boolean;
    }
  ): Promise<{ data: Product[]; total: number }> {
    let url = `/products?page=${page}&limit=${limit}`;

    if (filters) {
      if (filters.name) url += `&name=${filters.name}`;
      if (filters.category_id) url += `&category_id=${filters.category_id}`;
      if (filters.is_active !== undefined)
        url += `&is_active=${filters.is_active}`;
    }

    const response = await api.get(url);
    return response.data;
  },

  async getProductById(id: number): Promise<Product> {
    const response = await api.get(`/products/${id}`);
    return response.data;
  },

  async createProduct(product: CreateProductDto): Promise<Product> {
    const response = await api.post("/products", product);
    return response.data;
  },

  async updateProduct(id: number, product: UpdateProductDto): Promise<Product> {
    const response = await api.put(`/products/${id}`, product);
    return response.data;
  },

  async deleteProduct(id: number): Promise<void> {
    await api.delete(`/products/${id}`);
  },

  async getLowStockProducts(): Promise<Product[]> {
    const response = await api.get("/products/low-stock");
    return response.data;
  },

  async getProductsByCategory(categoryId: number): Promise<Product[]> {
    const response = await api.get(`/products/by-category/${categoryId}`);
    return response.data;
  },
};

export default ProductService;
