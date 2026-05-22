import api from "../common/api";
import {
  CreatePermissionFormData,
  Permission,
  PermissionList,
  UpdatePermissionFormData,
} from "./permission-schema";

const PermissionService = {
  /**
   * Obtém permissões paginadas e filtradas.
   * @param page O número da página.
   * @param limit O tamanho da página.
   * @param filters Filtros opcionais (nome da permissão, módulo).
   * @returns Um objeto com a lista de permissões e os dados de paginação.
   */
  async getPermissions(
    page = 1,
    limit = 10,
    filters?: {
      name?: string;
      module?: string;
      roleId?: string;
      isLinkedToRole?: boolean;
    }
  ): Promise<PermissionList> {
    const params = new URLSearchParams();
    params.append("page", page.toString());
    params.append("limit", limit.toString());

    if (filters?.name) {
      params.append("name", filters.name);
    }
    if (filters?.module) {
      params.append("module", filters.module);
    }
    if (filters?.roleId) {
      params.append("roleId", filters.roleId);
    }
    if (filters?.isLinkedToRole !== undefined) {
      params.append("isLinkedToRole", String(filters.isLinkedToRole));
    }

    let url = `/permissions?${params.toString()}`;
    const response = await api.get(url);

    if (!response.data.success) {
      const errorMessage =
        response.data.error ||
        response.data.message ||
        "Erro desconhecido ao obter permissões.";
      throw new Error(errorMessage);
    }

    const backendWrapper = response.data.data;

    if (!backendWrapper || !backendWrapper.data || !backendWrapper.pagination) {
      console.error("Dados ou paginação ausentes na resposta aninhada da API.");
      throw new Error("Formato de resposta da API inválido.");
    }

    return backendWrapper;
  },

  // Obtém uma permissão pelo ID
  async getPermissionById(id: string): Promise<Permission> {
    const response = await api.get(`/permissions/${id}`);
    if (!response.data.success) {
      throw new Error(
        response.data.error ||
          response.data.message ||
          "Erro ao obter permissão por ID."
      );
    }
    return response.data.data;
  },

  // Obtem as permissões por módulo
  async getPermissionsByModule(): Promise<Record<string, Permission[]>> {
    const response = await api.get("/permissions/by-module");
    if (!response.data.success) {
      throw new Error(
        response.data.error ||
          response.data.message ||
          "Erro ao obter permissões por módulo."
      );
    }
    return response.data.data;
  },

  // Obtem uma lista de modulos
  async getAvailableModules(): Promise<string[]> {
    const response = await api.get("/permissions/modules");
    if (!response.data.success) {
      throw new Error(
        response.data.error ||
          response.data.message ||
          "Erro ao obter módulos disponíveis."
      );
    }
    return response.data.data;
  },

  // Cria uma Permissão
  async createPermission(
    permission: CreatePermissionFormData
  ): Promise<Permission> {
    const response = await api.post("/permissions", permission);
    if (!response.data.success) {
      throw new Error(
        response.data.error ||
          response.data.message ||
          "Erro ao criar permissão."
      );
    }
    return response.data.data;
  },

  // Atualiza uma permissão
  async updatePermission(
    id: string,
    data: UpdatePermissionFormData
  ): Promise<Permission> {
    const response = await api.put(`/permissions/${id}`, data);
    if (!response.data.success) {
      throw new Error(
        response.data.error ||
          response.data.message ||
          "Erro ao atualizar permissão."
      );
    }
    return response.data.data;
  },

  // Deleta uma Permissão
  async deletePermission(id: string): Promise<void> {
    const response = await api.delete(`/permissions/${id}`);
    if (!response.data.success) {
      throw new Error(
        response.data.error ||
          response.data.message ||
          "Erro ao deletar permissão."
      );
    }
  },
};

export default PermissionService;
