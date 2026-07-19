// c:\Projetos\lumini-hub\Frontend\app\(system)\settings\users\create\page.tsx
"use client";

import { useEffect, useState } from "react";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { PagePermissionGuard } from "@/components/layout/PagePermissionGuard";
import { PageHeader } from "@/components/layout/page-header";
import { CreateUserForm } from "@/components/user/forms/create-user-form";
import { Role } from "@/services/auth/role-schema";
import RoleService from "@/services/auth/role-service";
import { Loader2 } from "lucide-react";

export default function CreateUserPage() {
  const [roles, setRoles] = useState<Role[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    RoleService.getRoles()
      .then(setRoles)
      .finally(() => setIsLoading(false));
  }, []);

  return (
    <PagePermissionGuard
      requiredPermissions={["users.create"]}
      accessDeniedMessage="Você não tem permissão para cadastrar usuários."
    >
      <div className="space-y-6">
        <PageHeader
          title="Novo Usuário"
          description="Preencha os dados do usuário e selecione as permissões de acesso."
        />

        <Card>
          <CardHeader>
            <CardTitle>Dados do Usuário</CardTitle>
            <CardDescription>
              Campos marcados com * são obrigatórios.
            </CardDescription>
          </CardHeader>
          <CardContent>
            {isLoading ? (
              <div className="flex flex-col items-center justify-center py-16">
                <Loader2 className="h-8 w-8 animate-spin text-primary mb-3" />
                <p className="text-sm text-muted-foreground">Carregando perfis...</p>
              </div>
            ) : (
              <CreateUserForm roles={roles} />
            )}
          </CardContent>
        </Card>
      </div>
    </PagePermissionGuard>
  );
}
