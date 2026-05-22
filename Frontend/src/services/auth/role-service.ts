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
    return {
      id: response.data.id,
      name: response.data.name,
      description: response.data.description,
      createdAt: response.data.created_at,
      updatedAt: response.data.updated_at,
    };
  },

  // 4. Atualizar um perfil existente
  async updateRole(id: number, data: UpdateRoleFormData): Promise<Role> {
    const payload = {
      name: data.name,
      description: data.description,
      permission_ids: data.permissionIds, // Backend espera snake_case
    };
    const response = await api.put(`/roles/${id}`, payload);
    return {
      id: response.data.id,
      name: response.data.name,
      description: response.data.description,
      createdAt: response.data.created_at,
      updatedAt: response.data.updated_at,
    };
  },

  // 5. Deletar um perfil
  async deleteRole(id: number): Promise<void> {
    await api.delete(`/roles/${id}`);
  },

  // NOVO: Função para vincular (linkar) múltiplas permissões a um perfil
  async linkPermissionsToRole(
    roleId: string,
    permissionIds: string[]
  ): Promise<void> {
    console.log(`Linking permissions ${permissionIds} to role ${roleId}`);
    const response = await api.post(`/roles/${roleId}/permissions/link`, {
      permission_ids: permissionIds, // Backend espera snake_case para o array de IDs
    });
    if (!response.data.success) {
      throw new Error(
        response.data.error ||
          response.data.message ||
          "Erro ao vincular permissões ao perfil."
      );
    }
  },

  // NOVO: Função para desvincular (unlinkar) uma única permissão de um perfil
  async unlinkPermissionFromRole(
    roleId: string,
    permissionId: string
  ): Promise<void> {
    console.log(`Unlinking permission ${permissionId} from role ${roleId}`);
    const response = await api.delete(
      `/roles/${roleId}/permissions/unlink/${permissionId}`
    );
    if (!response.data.success) {
      throw new Error(
        response.data.error ||
          response.data.message ||
          "Erro ao desvincular permissão do perfil."
      );
    }
  },
};

export default RoleService;
