"use client";

import { useState } from "react";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Edit2, MoreHorizontal, ShieldAlert, ShieldX, Trash2 } from "lucide-react";
import { useToast } from "@/components/ui/use-toast";
import { User } from "@/services/auth/user-schema";
import { Role } from "@/services/auth/role-schema";
import UserService from "@/services/auth/user-service";
import { ConfirmDeleteDialog } from "@/components/common/confirm-delete-dialog";
import { routes } from "@/config/routes";

interface UserTableProps {
  users: User[];
  roles: Role[];
  onUserUpdated: (updatedUser: User) => void;
  onUserDeleted: (id: number) => void;
  currentPage: number;
  totalPages: number;
  onPageChange: (page: number) => void;
}

export function UserTable({
  users,
  roles,
  onUserUpdated,
  onUserDeleted,
  currentPage,
  totalPages,
  onPageChange,
}: UserTableProps) {
  const { toast } = useToast();
  const [userToDelete, setUserToDelete] = useState<User | null>(null);
  const [isConfirmDeleteOpen, setIsConfirmDeleteOpen] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);

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

  const handleDeleteClick = (user: User) => {
    setUserToDelete(user);
    setIsConfirmDeleteOpen(true);
  };

  const handleConfirmDelete = async () => {
    if (!userToDelete) return;
    setIsDeleting(true);
    try {
      await UserService.deleteUser(userToDelete.id);
      onUserDeleted(userToDelete.id);
      setIsConfirmDeleteOpen(false);
      setUserToDelete(null);
      toast({
        title: "Usuário excluído",
        description: `O usuário ${userToDelete.name} foi excluído do sistema.`,
      });
    } catch (error: any) {
      toast({
        variant: "destructive",
        title: "Erro ao excluir usuário",
        description: error.message || "Não foi possível excluir o usuário.",
      });
    } finally {
      setIsDeleting(false);
    }
  };

  // Mapeia estilos CSS dinâmicos com base no nome do perfil
  const getRoleBadgeStyle = (roleName: string) => {
    const name = roleName.toLowerCase();
    if (name.includes("admin")) return "bg-blue-500/10 text-blue-600 border border-blue-200 dark:border-blue-800/30";
    if (name.includes("gerente") || name.includes("manager")) return "bg-violet-500/10 text-violet-600 border border-violet-200 dark:border-violet-800/30";
    if (name.includes("venda") || name.includes("sales")) return "bg-rose-500/10 text-rose-600 border border-rose-200 dark:border-rose-800/30";
    if (name.includes("suporte") || name.includes("support")) return "bg-amber-500/10 text-amber-600 border border-amber-200 dark:border-amber-800/30";
    if (name.includes("desenvolvedor") || name.includes("developer") || name.includes("dev")) return "bg-emerald-500/10 text-emerald-600 border border-emerald-200 dark:border-emerald-800/30";
    return "bg-slate-500/10 text-slate-600 border border-slate-200 dark:border-slate-800/30";
  };

  // Determina cores para os avatares de forma rotativa baseada no nome
  const getAvatarColor = (name: string) => {
    const colors = [
      "bg-red-500 text-white",
      "bg-orange-500 text-white",
      "bg-amber-500 text-white",
      "bg-green-500 text-white",
      "bg-teal-500 text-white",
      "bg-blue-500 text-white",
      "bg-indigo-500 text-white",
      "bg-purple-500 text-white",
      "bg-pink-500 text-white",
      "bg-rose-500 text-white"
    ];
    let hash = 0;
    for (let i = 0; i < name.length; i++) {
      hash = name.charCodeAt(i) + ((hash << 5) - hash);
    }
    const index = Math.abs(hash) % colors.length;
    return colors[index];
  };

  const getInitials = (name: string) => {
    return name
      .split(" ")
      .map((n) => n[0])
      .slice(0, 2)
      .join("")
      .toUpperCase();
  };

  return (
    <div className="space-y-4">
      <div className="rounded-lg border bg-card overflow-hidden shadow-sm">
        <div className="overflow-x-auto">
          <table className="w-full border-collapse text-left">
            <thead>
              <tr className="border-b bg-muted/40 text-xs font-semibold text-muted-foreground uppercase tracking-wider">
                <th className="px-6 py-4">Usuário</th>
                <th className="px-6 py-4">Nome de Acesso</th>
                <th className="px-6 py-4">Perfil</th>
                <th className="px-6 py-4">Status</th>
                <th className="px-6 py-4 text-right">Ações</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-border text-sm">
              {users.length === 0 ? (
                <tr>
                  <td colSpan={5} className="px-6 py-10 text-center text-muted-foreground">
                    Nenhum usuário cadastrado ou encontrado.
                  </td>
                </tr>
              ) : (
                users.map((user) => {
                  const userRole = roles.find((r) => r.id === user.role_id) || user.role;
                  const roleName = userRole?.name || "Sem Perfil";

                  return (
                    <tr key={user.id} className="hover:bg-muted/30 transition-colors">
                      <td className="px-6 py-4">
                        <div className="flex items-center gap-3">
                          <div className={`h-9 w-9 rounded-full flex items-center justify-center font-bold text-xs shadow-inner ${getAvatarColor(user.name)}`}>
                            {getInitials(user.name)}
                          </div>
                          <div className="flex flex-col">
                            <span className="font-semibold text-foreground">{user.name}</span>
                            <span className="text-xs text-muted-foreground">{user.email || "-"}</span>
                          </div>
                        </div>
                      </td>
                      <td className="px-6 py-4 text-muted-foreground font-mono text-xs">
                        {user.username}
                      </td>
                      <td className="px-6 py-4">
                        <span className={`rounded-full px-2.5 py-0.5 text-xs font-medium ${getRoleBadgeStyle(roleName)}`}>
                          {roleName}
                        </span>
                      </td>
                      <td className="px-6 py-4">
                        <span
                          className={`inline-flex items-center gap-1.5 rounded-full px-2.5 py-0.5 text-xs font-medium ${
                            user.is_active
                              ? "bg-emerald-500/10 text-emerald-600 border border-emerald-200 dark:border-emerald-800/30"
                              : "bg-rose-500/10 text-rose-600 border border-rose-200 dark:border-rose-800/30"
                          }`}
                        >
                          <span className={`h-1.5 w-1.5 rounded-full ${user.is_active ? "bg-emerald-500" : "bg-rose-500"}`} />
                          {user.is_active ? "Ativo" : "Inativo"}
                        </span>
                      </td>
                      <td className="px-6 py-4 text-right">
                        <DropdownMenu>
                          <DropdownMenuTrigger asChild>
                            <Button variant="ghost" size="icon" className="hover:bg-accent rounded-full h-8 w-8">
                              <MoreHorizontal className="h-4 w-4" />
                            </Button>
                          </DropdownMenuTrigger>
                          <DropdownMenuContent align="end" className="w-40">
                            <DropdownMenuLabel>Opções</DropdownMenuLabel>
                            <DropdownMenuSeparator />
                            <DropdownMenuItem asChild>
                              <Link href={routes.settings.users.edit(user.id)}>
                                <Edit2 className="mr-2 h-3.5 w-3.5" />
                                Editar Info
                              </Link>
                            </DropdownMenuItem>
                            <DropdownMenuItem onClick={() => toggleUserStatus(user.id)}>
                              {user.is_active ? (
                                <>
                                  <ShieldX className="mr-2 h-3.5 w-3.5 text-rose-500" />
                                  Desativar
                                </>
                              ) : (
                                <>
                                  <ShieldAlert className="mr-2 h-3.5 w-3.5 text-emerald-500" />
                                  Ativar
                                </>
                              )}
                            </DropdownMenuItem>
                            <DropdownMenuSeparator />
                            <DropdownMenuItem
                              onClick={() => handleDeleteClick(user)}
                              className="text-rose-600 focus:text-rose-600 focus:bg-rose-50 dark:focus:bg-rose-950/20"
                            >
                              <Trash2 className="mr-2 h-3.5 w-3.5" />
                              Excluir
                            </DropdownMenuItem>
                          </DropdownMenuContent>
                        </DropdownMenu>
                      </td>
                    </tr>
                  );
                })
              )}
            </tbody>
          </table>
        </div>
      </div>

      {/* Paginação Estilizada */}
      {totalPages > 1 && (
        <div className="flex items-center justify-between px-2">
          <p className="text-xs text-muted-foreground">
            Página {currentPage} de {totalPages}
          </p>
          <div className="flex items-center space-x-2">
            <Button
              variant="outline"
              size="sm"
              onClick={() => onPageChange(Math.max(currentPage - 1, 1))}
              disabled={currentPage === 1}
              className="h-8 text-xs"
            >
              Anterior
            </Button>
            <Button
              variant="outline"
              size="sm"
              onClick={() => onPageChange(Math.min(currentPage + 1, totalPages))}
              disabled={currentPage === totalPages}
              className="h-8 text-xs"
            >
              Próxima
            </Button>
          </div>
        </div>
      )}

      {/* Diálogo de confirmação de exclusão */}
      {isConfirmDeleteOpen && userToDelete && (
        <ConfirmDeleteDialog
          open={isConfirmDeleteOpen}
          onOpenChange={setIsConfirmDeleteOpen}
          itemName={`o usuário "${userToDelete.name}"`}
          onConfirm={handleConfirmDelete}
          title="Excluir Usuário?"
          description={`Tem certeza que deseja excluir permanentemente o usuário "${userToDelete.name}"? Esta ação removerá o acesso dele ao sistema e não poderá ser desfeita.`}
          confirmButtonText={isDeleting ? "Excluindo..." : "Confirmar Exclusão"}
        />
      )}
    </div>
  );
}
