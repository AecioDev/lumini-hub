import api from "./common/api";

export interface Category {
  id: number;
  name: string;
  description?: string;
  parent_id?: number;
  created_at: string;
  updated_at: string;
}

const CategoryService = {
  async getCategories(): Promise<Category[]> {
    const response = await api.get("/product-categories");
    return response.data;
  },

  async getCategoryById(id: number): Promise<Category> {
    const response = await api.get(`/product-categories/${id}`);
    return response.data;
  },

  async createCategory(category: {
    name: string;
    description?: string;
    parent_id?: number;
  }): Promise<Category> {
    const response = await api.post("/product-categories", category);
    return response.data;
  },

  async updateCategory(
    id: number,
    category: { name?: string; description?: string; parent_id?: number }
  ): Promise<Category> {
    const response = await api.put(`/product-categories/${id}`, category);
    return response.data;
  },

  async deleteCategory(id: number): Promise<void> {
    await api.delete(`/product-categories/${id}`);
  },
};

export default CategoryService;
