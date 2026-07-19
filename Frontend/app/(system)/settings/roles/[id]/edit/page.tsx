// c:\Projetos\lumini-hub\Frontend\app\(system)\settings\roles\[id]\edit\page.tsx
"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { PagePermissionGuard } from "@/components/layout/PagePermissionGuard";
import { PageHeader } from "@/components/layout/page-header";
import { EditRoleForm } from "@/components/settings/roles/forms/edit-role-form";
import { Role } from "@/services/auth/role-schema";
import RoleService from "@/services/auth/role-service";
import { Loader2 } from "lucide-react";

export default function EditRolePage() {
  const params = useParams<{ id: string }>();
  const roleId = Number(params.id);

  const [role, setRole] = useState<Role | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    RoleService.getRoleById(roleId)
      .then(setRole)
      .finally(() => setIsLoading(false));
  }, [roleId]);

  return (
    <PagePermissionGuard
      requiredPermissions={["admin.create_permissions"]}
      accessDeniedMessage="Esta tela é restrita à equipe de desenvolvimento."
    >
      <div className="space-y-6">
        <PageHeader
          title={`Editar Perfil${role ? `: ${role.name}` : ""}`}
          description="Atualize o nome e a descrição deste perfil."
        />

        <Card>
          <CardHeader>
            <CardTitle>Dados do Perfil</CardTitle>
            <CardDescription>
              Campos marcados com * são obrigatórios.
            </CardDescription>
          </CardHeader>
          <CardContent>
            {isLoading ? (
              <div className="flex flex-col items-center justify-center py-16">
                <Loader2 className="h-8 w-8 animate-spin text-primary mb-3" />
                <p className="text-sm text-muted-foreground">Carregando perfil...</p>
              </div>
            ) : !role ? (
              <div className="flex flex-col items-center justify-center py-16">
                <p className="text-sm text-muted-foreground">Perfil não encontrado.</p>
              </div>
            ) : (
              <EditRoleForm role={role} />
            )}
          </CardContent>
        </Card>
      </div>
    </PagePermissionGuard>
  );
}
