"use client";

import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useToast } from "@/components/ui/use-toast";
import { Role } from "@/services/auth/role-schema";
import { User } from "@/services/auth/user-schema";
import UserService from "@/services/auth/user-service";

interface EditUserDialogProps {
  isOpen: boolean;
  onOpenChange: (open: boolean) => void;
  user: User | null;
  roles: Role[];
  onUserUpdated: (user: User) => void;
  onUserChange: (user: User) => void;
}

export function EditUserDialog({
  isOpen,
  onOpenChange,
  user,
  roles,
  onUserUpdated,
  onUserChange,
}: EditUserDialogProps) {
  const { toast } = useToast();
  const validRoles = roles.filter((role) => role.id && role.name);

  const handleUpdateUser = async () => {
    if (!user) return;

    try {
      const updatedUser = await UserService.updateUser(user.id, {
        name: user.name,
        email: user.email,
        role_id: user.role_id,
        is_active: user.is_active,
      });

      onUserUpdated(updatedUser);
      onOpenChange(false);

      toast({
        title: "Usuário atualizado",
        description: `O usuário ${updatedUser.name} foi atualizado com sucesso`,
      });
    } catch (error: any) {
      toast({
        variant: "destructive",
        title: "Erro ao atualizar usuário",
        description:
          error.response?.data?.message ||
          "Ocorreu um erro ao atualizar o usuário",
      });
    }
  };

  if (!user) return null;

  return (
    <Dialog open={isOpen} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Editar Usuário</DialogTitle>
          <DialogDescription>
            Atualize as informações do usuário
          </DialogDescription>
        </DialogHeader>
        <div className="grid gap-4 py-4">
          <div className="grid gap-2">
            <Label htmlFor="edit-name">Nome Completo</Label>
            <Input
              id="edit-name"
              value={user.name}
              onChange={(e) => onUserChange({ ...user, name: e.target.value })}
            />
          </div>
          <div className="grid gap-2">
            <Label htmlFor="edit-email">Email</Label>
            <Input
              id="edit-email"
              type="email"
              value={user.email || ""}
              onChange={(e) => onUserChange({ ...user, email: e.target.value })}
            />
          </div>
          <div className="grid gap-2">
            <Label htmlFor="edit-role">Perfil</Label>
            <Select
              value={user.role_id.toString()}
              onValueChange={(value) =>
                onUserChange({
                  ...user,
                  role_id: Number.parseInt(value),
                })
              }
            >
              <SelectTrigger>
                <SelectValue placeholder="Selecione um perfil" />
              </SelectTrigger>
              <SelectContent>
                {validRoles.length > 0 ? (
                  validRoles.map((role) => (
                    <SelectItem key={role.id} value={role.id.toString()}>
                      {role.name}
                    </SelectItem>
                  ))
                ) : (
                  <p>Não há perfis válidos disponíveis.</p>
                )}
              </SelectContent>
            </Select>
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            Cancelar
          </Button>
          <Button onClick={handleUpdateUser}>Salvar Alterações</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
