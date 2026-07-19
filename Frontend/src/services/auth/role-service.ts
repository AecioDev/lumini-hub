// src/services/auth/role-service.ts (Exemplo)
import api from "@/services/common/api";
import { Role, CreateRoleFormData, UpdateRoleFormData } from "./role-schema";
import { Permission } from "./permission-schema";

const RoleService = {
  // 1. Listar todos os perfis (para o Select na página principal)
  async getRoles(): Promise<Role[]> {
    const response = await api.get("/roles");
    return response.data.data;
  },

  // 2. Obter um perfil por ID (para o modal de edição)
  async getRoleById(id: number): Promise<Role> {
    const response = await api.get(`/roles/${id}`);
    return response.data.data;
  },

  // 3. Criar um novo perfil
  async createRole(data: CreateRoleFormData): Promise<Role> {
    // Certifique-se de que os dados enviados estão em snake_case se o backend espera assim
    const payload = {
      name: data.name,
      description: data.description,
      permission_ids: data.permissionIds, // Backend espera snake_case
    };
    const response = await api.post("/roles", payload);
    return response.data.data;
  },

  // 4. Atualizar um perfil existente
  async updateRole(id: number, data: UpdateRoleFormData): Promise<Role> {
    const payload = {
      name: data.name,
      description: data.description,
      permission_ids: data.permissionIds, // Backend espera snake_case
    };
    const response = await api.put(`/roles/${id}`, payload);
    return response.data.data;
  },

  // 5. Deletar um perfil
  async deleteRole(id: number): Promise<void> {
    await api.delete(`/roles/${id}`);
  },
};

export default RoleService;
