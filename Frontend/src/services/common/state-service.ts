import api from "./api";
import { Country } from "./country-service";

export interface State {
  id: number;
  name: string;
  uf: string;
  ibge_code: string;
  country_id: number;
  country: Country;
}

export const StateService = {
  async getStates(
    page = 1,
    limit = 10,
    filters?: { name?: string; uf?: string; country_id?: number }
  ): Promise<State[]> {
    let url = `/states?page=${page}&limit=${limit}`;
    if (filters?.name) url += `&name=${filters.name}`;
    if (filters?.uf) url += `&uf=${filters.uf}`;
    if (filters?.country_id) url += `&country_id=${filters.country_id}`;

    const response = await api.get(url);
    return response.data.data;
  },

  async getStateById(id: number): Promise<State> {
    const response = await api.get(`/states/${id}`);
    return response.data;
  },

  async createState(data: Partial<State>): Promise<State> {
    const response = await api.post("/states", data);
    return response.data;
  },

  async updateState(id: number, data: Partial<State>): Promise<State> {
    const response = await api.put(`/states/${id}`, data);
    return response.data;
  },

  async deleteState(id: number): Promise<void> {
    await api.delete(`/states/${id}`);
  },
};
