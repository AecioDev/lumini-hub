import { Result, Typography } from "antd";
import { useParams } from "react-router-dom";
import { EditUserForm } from "@/components/users/forms/EditUserForm";

const { Title, Paragraph } = Typography;

export function EditUserPage() {
  const { id } = useParams<{ id: string }>();
  const userId = Number(id);

  if (!id || Number.isNaN(userId)) {
    return <Result status="404" title="Usuário não encontrado" />;
  }

  return (
    <div>
      <Title level={4} style={{ margin: "0 0 4px" }}>
        Editar Usuário
      </Title>
      <Paragraph type="secondary" style={{ marginBottom: 20 }}>
        Atualize os dados e permissões do usuário
      </Paragraph>
      <EditUserForm userId={userId} />
    </div>
  );
}
