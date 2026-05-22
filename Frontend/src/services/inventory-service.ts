import api from "./common/api";

export interface InventoryMovement {
  id: number;
  product_id: number;
  product_name?: string;
  quantity: number;
  previous_stock: number;
  new_stock: number;
  movement_type: "entrada" | "saida" | "ajuste";
  reference_id?: number;
  reference_type?: string;
  notes?: string;
  created_by: number;
  created_at: string;
}

export interface CreateMovementDto {
  product_id: number;
  quantity: number;
  movement_type: "entrada" | "saida" | "ajuste";
  reference_id?: number;
  reference_type?: string;
  notes?: string;
}

const InventoryService = {
  async getMovements(
    page = 1,
    limit = 10,
    filters?: {
      product_id?: number;
      movement_type?: string;
      startDate?: string;
      endDate?: string;
    }
  ): Promise<{ data: InventoryMovement[]; total: number }> {
    let url = `/inventory/movements?page=${page}&limit=${limit}`;

    if (filters) {
      if (filters.product_id) url += `&product_id=${filters.product_id}`;
      if (filters.movement_type)
        url += `&movement_type=${filters.movement_type}`;
      if (filters.startDate) url += `&startDate=${filters.startDate}`;
      if (filters.endDate) url += `&endDate=${filters.endDate}`;
    }

    const response = await api.get(url);
    return response.data;
  },

  async getMovementById(id: number): Promise<InventoryMovement> {
    const response = await api.get(`/inventory/movements/${id}`);
    return response.data;
  },

  async createMovement(
    movement: CreateMovementDto
  ): Promise<InventoryMovement> {
    const response = await api.post("/inventory/movements", movement);
    return response.data;
  },

  async getProductMovements(productId: number): Promise<InventoryMovement[]> {
    const response = await api.get(`/inventory/movements/product/${productId}`);
    return response.data;
  },

  async getInventorySummary(): Promise<{
    totalProducts: number;
    totalStock: number;
    lowStockCount: number;
    outOfStockCount: number;
    totalValue: number;
  }> {
    const response = await api.get("/inventory/summary");
    return response.data;
  },
};

export default InventoryService;
