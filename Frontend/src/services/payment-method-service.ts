import api from "./common/api";

export interface PaymentMethod {
  id: number;
  name: string;
  description?: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

const PaymentMethodService = {
  async getPaymentMethods(): Promise<PaymentMethod[]> {
    const response = await api.get("/payment-methods");
    return response.data;
  },

  async getPaymentMethodById(id: number): Promise<PaymentMethod> {
    const response = await api.get(`/payment-methods/${id}`);
    return response.data;
  },

  async createPaymentMethod(paymentMethod: {
    name: string;
    description?: string;
  }): Promise<PaymentMethod> {
    const response = await api.post("/payment-methods", paymentMethod);
    return response.data;
  },

  async updatePaymentMethod(
    id: number,
    paymentMethod: { name?: string; description?: string; is_active?: boolean }
  ): Promise<PaymentMethod> {
    const response = await api.put(`/payment-methods/${id}`, paymentMethod);
    return response.data;
  },

  async deletePaymentMethod(id: number): Promise<void> {
    await api.delete(`/payment-methods/${id}`);
  },
};

export default PaymentMethodService;
