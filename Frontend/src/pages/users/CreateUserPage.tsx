import { Typography } from "antd";
import { CreateUserForm } from "@/components/users/forms/CreateUserForm";

const { Title, Paragraph } = Typography;

export function CreateUserPage() {
  return (
    <div>
      <Title level={4} style={{ margin: "0 0 4px" }}>
        Novo Usuário
      </Title>
      <Paragraph type="secondary" style={{ marginBottom: 20 }}>
        Cadastre um novo membro da equipe e configure suas permissões
      </Paragraph>
      <CreateUserForm />
    </div>
  );
}
