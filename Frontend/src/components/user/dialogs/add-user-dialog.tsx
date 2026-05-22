"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
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
import { UserPlus } from "lucide-react";
import { CreateUserDto, User } from "@/services/auth/user-schema";
import { Role } from "@/services/auth/role-schema";
import UserService from "@/services/auth/user-service";

interface AddUserDialogProps {
  roles: Role[];
  onUserAdded: (user: User) => void;
}

export function AddUserDialog({ roles, onUserAdded }: AddUserDialogProps) {
  const [isOpen, setIsOpen] = useState(false);
  const [newUser, setNewUser] = useState<CreateUserDto>({
    username: "",
    password: "",
    name: "",
    email: "",
    role_id: 0,
  });
  const { toast } = useToast();
  const validRoles = roles.filter((role) => role.id && role.name);

  const resetForm = () => {
    setNewUser({
      username: "",
      password: "",
      name: "",
      email: "",
      role_id: 0,
    });
  };

  const handleAddUser = async () => {
    if (
      !newUser.username ||
      !newUser.password ||
      !newUser.name ||
      !newUser.role_id
    ) {
      toast({
        variant: "destructive",
        title: "Erro ao adicionar usuário",
        description: "Preencha todos os campos obrigatórios",
      });
      return;
    }

    try {
      const createdUser = await UserService.createUser(newUser);
      onUserAdded(createdUser);
      resetForm();
      setIsOpen(false);

      toast({
        title: "Usuário adicionado",
        description: `O usuário ${createdUser.name} foi adicionado com sucesso`,
      });
    } catch (error: any) {
      toast({
        variant: "destructive",
        title: "Erro ao adicionar usuário",
        description:
          error.response?.data?.message ||
          "Ocorreu um erro ao adicionar o usuário",
      });
    }
  };

  return (
    <Dialog open={isOpen} onOpenChange={setIsOpen}>
      <DialogTrigger asChild>
        <Button>
          <UserPlus className="mr-2 h-4 w-4" />
          Novo Usuário
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Adicionar Novo Usuário</DialogTitle>
          <DialogDescription>
            Preencha os dados para criar um novo usuário no sistema
          </DialogDescription>
        </DialogHeader>
        <div className="grid gap-4 py-4">
          <div className="grid gap-2">
            <Label htmlFor="name">Nome Completo</Label>
            <Input
              id="name"
              value={newUser.name}
              onChange={(e) => setNewUser({ ...newUser, name: e.target.value })}
            />
          </div>
          <div className="grid gap-2">
            <Label htmlFor="username">Nome de Usuário</Label>
            <Input
              id="username"
              value={newUser.username}
              placeholder="Nome_de_usuário sem espaços"
              onChange={(e) =>
                setNewUser({ ...newUser, username: e.target.value })
              }
            />
          </div>
          <div className="grid gap-2">
            <Label htmlFor="email">Email</Label>
            <Input
              id="email"
              type="email"
              value={newUser.email}
              placeholder="Informe um email válido!"
              onChange={(e) =>
                setNewUser({ ...newUser, email: e.target.value })
              }
            />
          </div>
          <div className="grid gap-2">
            <Label htmlFor="password">Senha</Label>
            <Input
              id="password"
              type="password"
              placeholder="Mínimo 6 caracteres"
              value={newUser.password}
              onChange={(e) =>
                setNewUser({ ...newUser, password: e.target.value })
              }
            />
          </div>
          <div className="grid gap-2">
            <Label htmlFor="role">Perfil</Label>
            <Select
              value={newUser.role_id ? newUser.role_id.toString() : ""}
              onValueChange={(value) =>
                setNewUser({ ...newUser, role_id: Number.parseInt(value) })
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
          <Button variant="outline" onClick={() => setIsOpen(false)}>
            Cancelar
          </Button>
          <Button onClick={handleAddUser}>Adicionar</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
