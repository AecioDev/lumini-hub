"use client";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Edit, MoreHorizontal, Trash2 } from "lucide-react";
import { useToast } from "@/components/ui/use-toast";
import { User } from "@/services/auth/user-schema";
import { Role } from "@/services/auth/role-schema";
import UserService from "@/services/auth/user-service";

interface UserTableProps {
  users: User[];
  roles: Role[];
  onEditUser: (user: User) => void;
  onUserUpdated: (updatedUser: User) => void;
  currentPage: number;
  totalPages: number;
  onPageChange: (page: number) => void;
}

export function UserTable({
  users,
  roles,
  onEditUser,
  onUserUpdated,
  currentPage,
  totalPages,
  onPageChange,
}: UserTableProps) {
  const { toast } = useToast();

  const toggleUserStatus = async (id: number) => {
    try {
      const userToToggle = users.find((u) => u.id === id);
      if (!userToToggle) return;

      const updatedUser = await UserService.updateUser(id, {
        is_active: !userToToggle.is_active,
      });

      onUserUpdated(updatedUser);

      toast({
        title: updatedUser.is_active ? "Usuário ativado" : "Usuário desativado",
        description: `O usuário ${updatedUser.name} foi ${
          updatedUser.is_active ? "ativado" : "desativado"
        } com sucesso`,
      });
    } catch (error: any) {
      toast({
        variant: "destructive",
        title: "Erro ao alterar status do usuário",
        description:
          error.response?.data?.message ||
          "Ocorreu um erro ao alterar o status do usuário",
      });
    }
  };

  return (
    <div>
      <div className="rounded-md border">
        <table className="w-full">
          <thead>
            <tr className="border-b bg-muted/50 text-left text-sm font-medium">
              <th className="px-4 py-3">Nome</th>
              <th className="px-4 py-3">Usuário</th>
              <th className="px-4 py-3">Email</th>
              <th className="px-4 py-3">Perfil</th>
              <th className="px-4 py-3">Status</th>
              <th className="px-4 py-3 text-right">Ações</th>
            </tr>
          </thead>
          <tbody>
            {users.map((user) => (
              <tr key={user.id} className="border-b">
                <td className="px-4 py-3">{user.name}</td>
                <td className="px-4 py-3">{user.username}</td>
                <td className="px-4 py-3">{user.email || "-"}</td>
                <td className="px-4 py-3">
                  <span className="rounded-full bg-primary/10 px-2 py-1 text-xs font-medium text-primary">
                    {roles.find((r) => r.id === user.role_id)?.name ||
                      user.role.name}
                  </span>
                </td>
                <td className="px-4 py-3">
                  <span
                    className={`rounded-full px-2 py-1 text-xs font-medium ${
                      user.is_active
                        ? "bg-green-100 text-green-700"
                        : "bg-red-100 text-red-700"
                    }`}
                  >
                    {user.is_active ? "Ativo" : "Inativo"}
                  </span>
                </td>
                <td className="px-4 py-3 text-right">
                  <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                      <Button variant="ghost" size="icon">
                        <MoreHorizontal className="h-4 w-4" />
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end">
                      <DropdownMenuLabel>Ações</DropdownMenuLabel>
                      <DropdownMenuSeparator />
                      <DropdownMenuItem onClick={() => onEditUser(user)}>
                        <Edit className="mr-2 h-4 w-4" />
                        Editar
                      </DropdownMenuItem>
                      <DropdownMenuItem
                        onClick={() => toggleUserStatus(user.id)}
                      >
                        <Trash2 className="mr-2 h-4 w-4" />
                        {user.is_active ? "Desativar" : "Ativar"}
                      </DropdownMenuItem>
                    </DropdownMenuContent>
                  </DropdownMenu>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {totalPages > 1 && (
        <div className="mt-4 flex items-center justify-end space-x-2">
          <Button
            variant="outline"
            size="sm"
            onClick={() => onPageChange(Math.max(currentPage - 1, 1))}
            disabled={currentPage === 1}
          >
            Anterior
          </Button>
          <span className="text-sm text-muted-foreground">
            Página {currentPage} de {totalPages}
          </span>
          <Button
            variant="outline"
            size="sm"
            onClick={() => onPageChange(Math.min(currentPage + 1, totalPages))}
            disabled={currentPage === totalPages}
          >
            Próxima
          </Button>
        </div>
      )}
    </div>
  );
}
