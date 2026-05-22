import api from "./common/api";

export interface Unit {
  id: number;
  name: string;
  abbreviation: string;
  created_at: string;
  updated_at: string;
}

const UnitService = {
  async getUnits(): Promise<Unit[]> {
    const response = await api.get("/measurement-units");
    return response.data;
  },

  async getUnitById(id: number): Promise<Unit> {
    const response = await api.get(`/measurement-units/${id}`);
    return response.data;
  },

  async createUnit(unit: {
    name: string;
    abbreviation: string;
  }): Promise<Unit> {
    const response = await api.post("/measurement-units", unit);
    return response.data;
  },

  async updateUnit(
    id: number,
    unit: { name?: string; abbreviation?: string }
  ): Promise<Unit> {
    const response = await api.put(`/measurement-units/${id}`, unit);
    return response.data;
  },

  async deleteUnit(id: number): Promise<void> {
    await api.delete(`/measurement-units/${id}`);
  },
};

export default UnitService;
