// c:\Projetos\lumini-hub\Frontend\app\(system)\settings\users\page.tsx
"use client";

import { useState, useEffect, useMemo, useCallback } from "react";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useToast } from "@/components/ui/use-toast";
import { PagePermissionGuard } from "@/components/layout/PagePermissionGuard";
import { PageHeader } from "@/components/layout/page-header";
import { PermissionedLinkButton } from "@/components/common/PermissionedLinkButton";
import { User } from "@/services/auth/user-schema";
import { Role } from "@/services/auth/role-schema";
import UserService from "@/services/auth/user-service";
import RoleService from "@/services/auth/role-service";
import { UserTable } from "@/components/user/tables/user-table";
import { routes } from "@/config/routes";
import { Loader2 } from "lucide-react";

export default function UsersPage() {
  const [users, setUsers] = useState<User[]>([]);
  const [roles, setRoles] = useState<Role[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  // Paginação
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);

  // Filtros na listagem de usuários
  const [searchTerm, setSearchTerm] = useState("");
  const [filterRole, setFilterRole] = useState("all");
  const [filterStatus, setFilterStatus] = useState("all");

  const { toast } = useToast();

  const loadData = useCallback(async () => {
    try {
      setIsLoading(true);
      const [usersResponse, rolesResponse] = await Promise.all([
        UserService.getUsers(currentPage, 10),
        RoleService.getRoles(),
      ]);

      setUsers(usersResponse.users);
      setTotalPages(Math.ceil(usersResponse.pagination.totalPages / 10) || 1);
      setRoles(rolesResponse);
    } catch (error) {
      console.error("Erro ao carregar dados:", error);
      toast({
        variant: "destructive",
        title: "Erro ao carregar dados",
        description: "Não foi possível carregar a lista de usuários.",
      });
    } finally {
      setIsLoading(false);
    }
  }, [currentPage, toast]);

  useEffect(() => {
    loadData();
  }, [loadData]);

  const handleUserUpdated = (updatedUser: User) => {
    setUsers((prev) => prev.map((u) => (u.id === updatedUser.id ? updatedUser : u)));
    loadData();
  };

  const handleUserDeleted = (id: number) => {
    setUsers((prev) => prev.filter((u) => u.id !== id));
    loadData();
  };

  // Filtragem dos usuários em memória na página atual (para busca instantânea responsiva)
  const filteredUsers = useMemo(() => {
    return users.filter((user) => {
      const matchesSearch =
        user.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
        user.username.toLowerCase().includes(searchTerm.toLowerCase()) ||
        (user.email && user.email.toLowerCase().includes(searchTerm.toLowerCase()));

      const userRole = roles.find((r) => r.id === user.role_id) || user.role;
      const matchesRole =
        filterRole === "all" ||
        (userRole && userRole.id.toString() === filterRole);

      const matchesStatus =
        filterStatus === "all" ||
        (filterStatus === "active" && user.is_active) ||
        (filterStatus === "inactive" && !user.is_active);

      return matchesSearch && matchesRole && matchesStatus;
    });
  }, [users, roles, searchTerm, filterRole, filterStatus]);

  return (
    <PagePermissionGuard
      requiredPermissions={["users.view"]}
      accessDeniedMessage="Você não tem permissão para visualizar usuários."
    >
      <div className="space-y-6">
        <PageHeader
          title="Usuários"
          description="Gerencie as contas de usuário e suas permissões de acesso."
          actions={
            <PermissionedLinkButton
              href={routes.settings.users.create}
              permission="users.create"
              tooltipMessage="Você não pode criar novos usuários."
              iconName="mdi:account-plus"
            >
              Novo Usuário
            </PermissionedLinkButton>
          }
        />

        {/* Tabela de Usuários */}
        <Card>
          <CardHeader className="pb-4">
            <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
              <div>
                <CardTitle>Todos os Usuários</CardTitle>
                <CardDescription>
                  Visualize, edite e configure as contas de usuário do sistema.
                </CardDescription>
              </div>

              {/* Filtros da Tabela */}
              <div className="flex flex-wrap items-center gap-3">
                <Input
                  placeholder="Buscar usuário..."
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                  className="w-full sm:w-48 h-9 text-xs"
                />

                <Select value={filterRole} onValueChange={setFilterRole}>
                  <SelectTrigger className="w-full sm:w-36 h-9 text-xs">
                    <SelectValue placeholder="Filtrar por Perfil" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">Todos Perfis</SelectItem>
                    {roles.map((role) => (
                      <SelectItem key={role.id} value={role.id.toString()}>
                        {role.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>

                <Select value={filterStatus} onValueChange={setFilterStatus}>
                  <SelectTrigger className="w-full sm:w-32 h-9 text-xs">
                    <SelectValue placeholder="Filtrar Status" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">Todos Status</SelectItem>
                    <SelectItem value="active">Ativo</SelectItem>
                    <SelectItem value="inactive">Inativo</SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </div>
          </CardHeader>

          <CardContent>
            {isLoading ? (
              <div className="flex flex-col items-center justify-center py-20">
                <Loader2 className="h-8 w-8 animate-spin text-primary mb-3" />
                <p className="text-sm text-muted-foreground">Carregando dados dos usuários...</p>
              </div>
            ) : (
              <UserTable
                users={filteredUsers}
                roles={roles}
                onUserUpdated={handleUserUpdated}
                onUserDeleted={handleUserDeleted}
                currentPage={currentPage}
                totalPages={totalPages}
                onPageChange={setCurrentPage}
              />
            )}
          </CardContent>
        </Card>
      </div>
    </PagePermissionGuard>
  );
}
