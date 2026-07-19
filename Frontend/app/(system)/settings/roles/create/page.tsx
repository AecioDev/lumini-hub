// c:\Projetos\lumini-hub\Frontend\app\(system)\settings\roles\create\page.tsx
"use client";

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { PagePermissionGuard } from "@/components/layout/PagePermissionGuard";
import { PageHeader } from "@/components/layout/page-header";
import { CreateRoleForm } from "@/components/settings/roles/forms/create-role-form";

export default function CreateRolePage() {
  return (
    <PagePermissionGuard
      requiredPermissions={["admin.create_permissions"]}
      accessDeniedMessage="Esta tela é restrita à equipe de desenvolvimento."
    >
      <div className="space-y-6">
        <PageHeader
          title="Novo Perfil"
          description="Defina o nome, a descrição e as permissões iniciais deste perfil."
        />

        <Card>
          <CardHeader>
            <CardTitle>Dados do Perfil</CardTitle>
            <CardDescription>
              Campos marcados com * são obrigatórios.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <CreateRoleForm />
          </CardContent>
        </Card>
      </div>
    </PagePermissionGuard>
  );
}
