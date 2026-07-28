import { PageHeader } from "@/components/common/PageHeader";
import { CreateUserForm } from "@/components/users/forms/CreateUserForm";

export function CreateUserPage() {
  return (
    <div>
      <PageHeader
        title="Novo Usuário"
        subtitle="Cadastre um novo membro da equipe e configure suas permissões"
        backTo="/settings/users"
      />
      <CreateUserForm />
    </div>
  );
}
