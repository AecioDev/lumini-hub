import api from "./api";
import { State } from "./state-service";

export interface City {
  id: number;
  name: string;
  ibge_code: string;
  state_id: number;
  state: State;
}

export const CityService = {
  async getCities(
    page = 1,
    limit = 10,
    filters?: { name?: string; state_id?: number }
  ): Promise<City[]> {
    let url = `/cities?page=${page}&limit=${limit}`;
    if (filters?.name) url += `&name=${filters.name}`;
    if (filters?.state_id) url += `&state_id=${filters.state_id}`;

    const response = await api.get(url);
    return response.data.data;
  },

  async getCityById(id: number): Promise<City> {
    const response = await api.get(`/cities/${id}`);
    return response.data;
  },

  async createCity(data: Partial<City>): Promise<City> {
    const response = await api.post("/cities", data);
    return response.data;
  },

  async updateCity(id: number, data: Partial<City>): Promise<City> {
    const response = await api.put(`/cities/${id}`, data);
    return response.data;
  },

  async deleteCity(id: number): Promise<void> {
    await api.delete(`/cities/${id}`);
  },
};
