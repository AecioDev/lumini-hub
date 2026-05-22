import api from "./api";

export interface Country {
  id: number;
  name: string;
  bacen_code: string;
  phone_code: string;
}

export const CountryService = {
  async getCountries(
    page = 1,
    limit = 10,
    filters?: { name?: string }
  ): Promise<Country[]> {
    let url = `/countries?page=${page}&limit=${limit}`;
    if (filters?.name) url += `&name=${filters.name}`;

    const response = await api.get(url);
    return response.data.data;
  },

  async getCountryById(id: number): Promise<Country> {
    const response = await api.get(`/countries/${id}`);
    return response.data;
  },

  async createCountry(data: Partial<Country>): Promise<Country> {
    const response = await api.post("/countries", data);
    return response.data;
  },

  async updateCountry(id: number, data: Partial<Country>): Promise<Country> {
    const response = await api.put(`/countries/${id}`, data);
    return response.data;
  },

  async deleteCountry(id: number): Promise<void> {
    await api.delete(`/countries/${id}`);
  },
};
