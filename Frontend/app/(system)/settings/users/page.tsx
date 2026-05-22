"use client";

import { useState, useEffect } from "react";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { useToast } from "@/components/ui/use-toast";
import { PagePermissionGuard } from "@/components/layout/PagePermissionGuard";
import { User } from "@/services/auth/user-schema";
import { Role } from "@/services/auth/role-schema";
import UserService from "@/services/auth/user-service";
import RoleService from "@/services/auth/role-service";
import { AddUserDialog } from "@/components/user/dialogs/add-user-dialog";
import { UserTable } from "@/components/user/tables/user-table";
import { EditUserDialog } from "@/components/user/dialogs/edit-user-dialog";

export default function UsersPage() {
  const [users, setUsers] = useState<User[]>([]);
  const [roles, setRoles] = useState<Role[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isEditUserOpen, setIsEditUserOpen] = useState(false);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [selectedUser, setSelectedUser] = useState<User | null>(null);
  const { toast } = useToast();

  // Carregar usuários e perfis
  useEffect(() => {
    const loadData = async () => {
      try {
        setIsLoading(true);
        const [usersResponse, rolesResponse] = await Promise.all([
          UserService.getUsers(currentPage, 10),
          RoleService.getRoles(),
        ]);

        setUsers(usersResponse.users);
        setTotalPages(Math.ceil(usersResponse.pagination.totalPages / 10));
        setRoles(rolesResponse);
      } catch (error) {
        console.error("Erro ao carregar dados:", error);
        toast({
          variant: "destructive",
          title: "Erro ao carregar dados",
          description:
            "Não foi possível carregar a lista de usuários e perfis.",
        });
      } finally {
        setIsLoading(false);
      }
    };

    loadData();
  }, [currentPage, toast]);

  const handleEditUser = (user: User) => {
    setSelectedUser(user);
    setIsEditUserOpen(true);
  };

  const handleUserAdded = (user: User) => {
    setUsers([...users, user]);
  };

  const handleUserUpdated = (updatedUser: User) => {
    setUsers(users.map((u) => (u.id === updatedUser.id ? updatedUser : u)));
  };

  return (
    <PagePermissionGuard
      requiredPermissions={["customers.view"]}
      accessDeniedMessage="Você não tem permissão para visualizar clientes."
    >
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold">Usuários</h1>
            <p className="text-muted-foreground">
              Gerenciamento de usuários do sistema
            </p>
          </div>
          <AddUserDialog roles={roles} onUserAdded={handleUserAdded} />
        </div>

        <Card>
          <CardHeader>
            <CardTitle>Lista de Usuários</CardTitle>
            <CardDescription>
              Gerencie os usuários do sistema e suas permissões
            </CardDescription>
          </CardHeader>
          <CardContent>
            <UserTable
              users={users}
              roles={roles}
              onEditUser={handleEditUser}
              onUserUpdated={handleUserUpdated}
              currentPage={currentPage}
              totalPages={totalPages}
              onPageChange={setCurrentPage}
            />
          </CardContent>
        </Card>

        <EditUserDialog
          isOpen={isEditUserOpen}
          onOpenChange={setIsEditUserOpen}
          user={selectedUser}
          roles={roles}
          onUserUpdated={handleUserUpdated}
          onUserChange={setSelectedUser}
        />
      </div>
    </PagePermissionGuard>
  );
}
