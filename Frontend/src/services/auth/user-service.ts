import api from "../common/api";
import {
  ChangePasswordDto,
  CreateUserDto,
  UpdateUserDto,
  User,
  UserList,
} from "./user-schema";

const UserService = {
  async getUsers(page = 1, limit = 10): Promise<UserList> {
    const response = await api.get(`/users?page=${page}&limit=${limit}`);

    return response.data.data;
  },

  async getUserById(id: number): Promise<User> {
    const response = await api.get(`/users/${id}`);
    return response.data.data;
  },

  async createUser(user: CreateUserDto): Promise<User> {
    const response = await api.post("/users", user);
    return response.data.data;
  },

  async updateUser(id: number, user: UpdateUserDto): Promise<User> {
    const response = await api.put(`/users/${id}`, user);
    return response.data.data;
  },

  async deleteUser(id: number): Promise<void> {
    await api.delete(`/users/${id}`);
  },

  async changePassword(id: number, data: ChangePasswordDto): Promise<void> {
    await api.put(`/users/${id}/password`, data);
  },

  async updateUserPermissions(id: number, permissionIds: number[]): Promise<User> {
    const response = await api.put(`/users/${id}/permissions`, {
      permission_ids: permissionIds,
    });
    return response.data.data;
  },
};

export default UserService;
