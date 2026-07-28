import { Result } from "antd";
import { useParams } from "react-router-dom";
import { PageHeader } from "@/components/common/PageHeader";
import { EditUserForm } from "@/components/users/forms/EditUserForm";

export function EditUserPage() {
  const { id } = useParams<{ id: string }>();
  const userId = Number(id);

  if (!id || Number.isNaN(userId)) {
    return <Result status="404" title="Usuário não encontrado" />;
  }

  return (
    <div>
      <PageHeader
        title="Editar Usuário"
        subtitle="Atualize os dados e permissões do usuário"
        backTo="/settings/users"
      />
      <EditUserForm userId={userId} />
    </div>
  );
}
