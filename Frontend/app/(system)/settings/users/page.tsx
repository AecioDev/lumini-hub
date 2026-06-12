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
import { Button } from "@/components/ui/button";
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
import { User } from "@/services/auth/user-schema";
import { Role } from "@/services/auth/role-schema";
import UserService from "@/services/auth/user-service";
import RoleService from "@/services/auth/role-service";
import { AddUserDialog } from "@/components/user/dialogs/add-user-dialog";
import { UserTable } from "@/components/user/tables/user-table";
import { EditUserDialog } from "@/components/user/dialogs/edit-user-dialog";
import { CreateRoleDialog } from "@/components/settings/roles/dialogs/create-role-dialog";
import { EditRoleDialog } from "@/components/settings/roles/dialogs/edit-role-dialog";
import { ConfirmDeleteDialog } from "@/components/common/confirm-delete-dialog";
import {
  Shield,
  ShieldAlert,
  Users,
  Settings,
  Key,
  Lock,
  Plus,
  Loader2,
  Trash2,
  MoreVertical,
} from "lucide-react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

export default function UsersAndRolesPage() {
  const [users, setUsers] = useState<User[]>([]);
  const [roles, setRoles] = useState<Role[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  
  // Paginação
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);

  // Modais de Usuário
  const [isEditUserOpen, setIsEditUserOpen] = useState(false);
  const [selectedUser, setSelectedUser] = useState<User | null>(null);

  // Modais de Perfil (Role)
  const [isCreateRoleOpen, setIsCreateRoleOpen] = useState(false);
  const [isEditRoleOpen, setIsEditRoleOpen] = useState(false);
  const [selectedRole, setSelectedRole] = useState<Role | null>(null);
  
  // Exclusão de Perfil
  const [roleToDelete, setRoleToDelete] = useState<Role | null>(null);
  const [isConfirmDeleteRoleOpen, setIsConfirmDeleteRoleOpen] = useState(false);
  const [isDeletingRole, setIsDeletingRole] = useState(false);

  // Filtros na listagem de usuários
  const [searchTerm, setSearchTerm] = useState("");
  const [filterRole, setFilterRole] = useState("all");
  const [filterStatus, setFilterStatus] = useState("all");

  const { toast } = useToast();

  // Função para carregar usuários e perfis
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
        description: "Não foi possível carregar a lista de usuários e perfis.",
      });
    } finally {
      setIsLoading(false);
    }
  }, [currentPage, toast]);

  useEffect(() => {
    loadData();
  }, [loadData]);

  // Handlers para Usuários
  const handleEditUser = (user: User) => {
    setSelectedUser(user);
    setIsEditUserOpen(true);
  };

  const handleUserAdded = (user: User) => {
    setUsers((prev) => [user, ...prev]);
    loadData(); // Recarrega para manter avatares e contagem dos perfis corretos
  };

  const handleUserUpdated = (updatedUser: User) => {
    setUsers((prev) => prev.map((u) => (u.id === updatedUser.id ? updatedUser : u)));
    loadData(); // Recarrega para manter sincronia
  };

  const handleUserDeleted = (id: number) => {
    setUsers((prev) => prev.filter((u) => u.id !== id));
    loadData(); // Recarrega para manter sincronia
  };

  // Handlers para Perfis (Roles)
  const handleEditRoleClick = (role: Role) => {
    setSelectedRole(role);
    setIsEditRoleOpen(true);
  };

  const handleRoleCreatedOrUpdated = () => {
    loadData();
  };

  const handleDeleteRoleClick = (role: Role) => {
    setRoleToDelete(role);
    setIsConfirmDeleteRoleOpen(true);
  };

  const handleConfirmDeleteRole = async () => {
    if (!roleToDelete) return;
    setIsDeletingRole(true);
    try {
      await RoleService.deleteRole(roleToDelete.id);
      setIsConfirmDeleteRoleOpen(false);
      setRoleToDelete(null);
      toast({
        title: "Perfil excluído",
        description: `O perfil "${roleToDelete.name}" foi excluído com sucesso.`,
      });
      loadData();
    } catch (error: any) {
      console.error("Erro ao excluir perfil:", error);
      toast({
        variant: "destructive",
        title: "Erro ao excluir perfil",
        description: error.response?.data?.message || "Não foi possível excluir o perfil (certifique-se de que ele não está associado a nenhum usuário).",
      });
    } finally {
      setIsDeletingRole(false);
    }
  };

  // Funções Visuais para os Cards de Perfis
  const getRoleColor = (index: number) => {
    const colors = [
      "#3b82f6", // Blue
      "#8b5cf6", // Violet
      "#ec4899", // Pink
      "#f59e0b", // Amber
      "#10b981", // Emerald
      "#06b6d4", // Cyan
    ];
    return colors[index % colors.length];
  };

  const getRoleIcon = (roleName: string) => {
    const name = roleName.toLowerCase();
    if (name.includes("admin")) return <Shield className="h-5 w-5 text-white" />;
    if (name.includes("gerente") || name.includes("manager")) return <ShieldAlert className="h-5 w-5 text-white" />;
    if (name.includes("venda") || name.includes("sales")) return <Users className="h-5 w-5 text-white" />;
    if (name.includes("suporte") || name.includes("support")) return <Settings className="h-5 w-5 text-white" />;
    if (name.includes("desenvolvedor") || name.includes("developer") || name.includes("dev")) return <Key className="h-5 w-5 text-white" />;
    return <Lock className="h-5 w-5 text-white" />;
  };

  const getInitials = (name: string) => {
    return name
      .split(" ")
      .map((n) => n[0])
      .slice(0, 2)
      .join("")
      .toUpperCase();
  };

  const getAvatarColor = (name: string) => {
    const colors = [
      "bg-red-500",
      "bg-orange-500",
      "bg-amber-500",
      "bg-green-500",
      "bg-teal-500",
      "bg-blue-500",
      "bg-indigo-500",
      "bg-purple-500",
      "bg-pink-500",
      "bg-rose-500",
    ];
    let hash = 0;
    for (let i = 0; i < name.length; i++) {
      hash = name.charCodeAt(i) + ((hash << 5) - hash);
    }
    const index = Math.abs(hash) % colors.length;
    return colors[index];
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
      requiredPermissions={["customers.view"]}
      accessDeniedMessage="Você não tem permissão para visualizar usuários e perfis."
    >
      <div className="space-y-8 px-4 py-6 md:px-8 md:py-8 bg-background text-foreground min-h-screen">
        {/* Cabeçalho */}
        <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between border-b pb-6 border-border/40">
          <div>
            <h1 className="text-3xl font-bold tracking-tight lg:text-4xl">
              Usuários e Perfis
            </h1>
            <p className="text-muted-foreground mt-1 text-sm md:text-base">
              Gerencie usuários, perfis de acesso (Roles) e permissões de forma centralizada e visual.
            </p>
          </div>
          <div className="flex gap-2">
            <Button
              variant="outline"
              onClick={() => setIsCreateRoleOpen(true)}
              className="border-primary/20 hover:border-primary/40 hover:bg-primary/5 transition-all text-xs md:text-sm"
            >
              <Plus className="mr-1.5 h-4 w-4" />
              Novo Perfil
            </Button>
            <AddUserDialog roles={roles} onUserAdded={handleUserAdded} />
          </div>
        </div>

        {/* Grid de Cards de Perfis (Roles) */}
        <div className="space-y-4">
          <div className="flex items-center justify-between">
            <h2 className="text-lg font-semibold tracking-tight">Perfis de Acesso</h2>
            <span className="text-xs text-muted-foreground bg-accent px-2 py-0.5 rounded-full font-medium">
              {roles.length} perfis cadastrados
            </span>
          </div>

          {isLoading ? (
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
              {[1, 2, 3, 4].map((i) => (
                <Card key={i} className="animate-pulse border-border/50">
                  <CardHeader className="space-y-2 pb-4">
                    <div className="h-10 w-10 rounded-lg bg-muted" />
                    <div className="h-4 w-24 bg-muted rounded" />
                  </CardHeader>
                  <CardContent className="h-16 bg-muted/20" />
                </Card>
              ))}
            </div>
          ) : (
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
              {roles.map((role, idx) => {
                // Filtra usuários vinculados a esta role
                const roleUsers = users.filter((u) => u.role_id === role.id);
                const roleColor = getRoleColor(idx);

                return (
                  <Card
                    key={role.id}
                    className="group border border-primary/20 hover:border-primary/40 hover:shadow-md transition-all duration-300 bg-card overflow-hidden flex flex-col justify-between"
                  >
                    <CardHeader className="pb-3 flex flex-row items-start justify-between space-y-0">
                      <div className="flex items-center gap-3">
                        <span
                          className="grid h-10 w-10 place-content-center rounded-lg shadow-inner text-white transition-transform duration-300 group-hover:scale-105"
                          style={{ backgroundColor: roleColor }}
                        >
                          {getRoleIcon(role.name)}
                        </span>
                        <div>
                          <CardTitle className="text-base font-bold text-foreground">
                            {role.name}
                          </CardTitle>
                          <CardDescription className="text-xs line-clamp-1 mt-0.5">
                            {role.description || "Sem descrição definida."}
                          </CardDescription>
                        </div>
                      </div>

                      {/* Menu de Opções no Card */}
                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                          <Button variant="ghost" size="icon" className="h-7 w-7 rounded-full text-muted-foreground hover:text-foreground">
                            <MoreVertical className="h-4 w-4" />
                          </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end" className="w-32">
                          <DropdownMenuItem onClick={() => handleEditRoleClick(role)}>
                            Renomear/Editar
                          </DropdownMenuItem>
                          <DropdownMenuItem
                            onClick={() => handleDeleteRoleClick(role)}
                            className="text-rose-600 focus:text-rose-600 focus:bg-rose-50 dark:focus:bg-rose-950/20"
                          >
                            Excluir Perfil
                          </DropdownMenuItem>
                        </DropdownMenuContent>
                      </DropdownMenu>
                    </CardHeader>

                    <CardContent className="pt-0">
                      {/* Avatares dos Usuários */}
                      <div className="flex items-center justify-between mt-4 border-t pt-3 border-border/40">
                        <div className="flex -space-x-2.5 overflow-hidden">
                          {roleUsers.length === 0 ? (
                            <span className="text-xs text-muted-foreground italic">
                              Sem usuários
                            </span>
                          ) : (
                            roleUsers.slice(0, 4).map((u) => (
                              <div
                                key={u.id}
                                className={`inline-flex h-7 w-7 items-center justify-center rounded-full text-[10px] font-bold ring-2 ring-background shadow-sm ${getAvatarColor(
                                  u.name
                                )} text-white`}
                                title={u.name}
                              >
                                {getInitials(u.name)}
                              </div>
                            ))
                          )}
                          {roleUsers.length > 4 && (
                            <div className="inline-flex h-7 w-7 items-center justify-center rounded-full bg-accent text-[9px] font-bold text-muted-foreground ring-2 ring-background">
                              +{roleUsers.length - 4}
                            </div>
                          )}
                        </div>
                        {roleUsers.length > 0 && (
                          <span className="text-xs font-medium text-muted-foreground">
                            Total {roleUsers.length}
                          </span>
                        )}
                      </div>

                      {/* Botão de Editar Perfil */}
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => handleEditRoleClick(role)}
                        className="w-full mt-4 border-border/80 hover:bg-primary/5 hover:text-primary transition-all text-xs h-8"
                      >
                        Configurar Acessos
                      </Button>
                    </CardContent>
                  </Card>
                );
              })}
            </div>
          )}
        </div>

        {/* Tabela de Usuários */}
        <Card className="border border-primary/20 shadow-sm bg-card mt-8">
          <CardHeader className="pb-4">
            <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
              <div>
                <CardTitle className="text-lg font-bold">Todos os Usuários</CardTitle>
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
                onEditUser={handleEditUser}
                onUserUpdated={handleUserUpdated}
                onUserDeleted={handleUserDeleted}
                currentPage={currentPage}
                totalPages={totalPages}
                onPageChange={setCurrentPage}
              />
            )}
          </CardContent>
        </Card>

        {/* Diálogos (Modais) */}
        
        {/* Criar Perfil */}
        <CreateRoleDialog
          isOpen={isCreateRoleOpen}
          onOpenChange={setIsCreateRoleOpen}
          onRoleCreated={handleRoleCreatedOrUpdated}
        />

        {/* Editar Perfil */}
        <EditRoleDialog
          isOpen={isEditRoleOpen}
          onOpenChange={setIsEditRoleOpen}
          role={selectedRole}
          onRoleUpdated={handleRoleCreatedOrUpdated}
        />

        {/* Editar Usuário */}
        <EditUserDialog
          isOpen={isEditUserOpen}
          onOpenChange={setIsEditUserOpen}
          user={selectedUser}
          roles={roles}
          onUserUpdated={handleUserUpdated}
          onUserChange={setSelectedUser}
        />

        {/* Confirmação de Exclusão de Perfil */}
        {isConfirmDeleteRoleOpen && roleToDelete && (
          <ConfirmDeleteDialog
            open={isConfirmDeleteRoleOpen}
            onOpenChange={setIsConfirmDeleteRoleOpen}
            itemName={`o perfil "${roleToDelete.name}"`}
            onConfirm={handleConfirmDeleteRole}
            title="Excluir Perfil?"
            description={`Tem certeza que deseja excluir permanentemente o perfil "${roleToDelete.name}"? Esta ação não pode ser desfeita e irá falhar caso haja usuários vinculados a ele.`}
            confirmButtonText={isDeletingRole ? "Excluindo..." : "Confirmar Exclusão"}
          />
        )}
      </div>
    </PagePermissionGuard>
  );
}
