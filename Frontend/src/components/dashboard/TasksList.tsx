import { Card, List, Typography, theme } from "antd";
import type { TeamTask } from "@/types/crm";

const { Text } = Typography;

export function TasksList({ tasks }: { tasks: TeamTask[] }) {
  const { token } = theme.useToken();

  return (
    <Card title="Tarefas da Equipe">
      <List
        itemLayout="horizontal"
        dataSource={tasks}
        renderItem={(task) => (
          <List.Item style={{ border: "none", padding: "8px 0" }}>
            <div style={{ display: "flex", alignItems: "flex-start", gap: 10, width: "100%" }}>
              <span
                style={{
                  width: 14,
                  height: 14,
                  borderRadius: "50%",
                  marginTop: 3,
                  flexShrink: 0,
                  display: "inline-block",
                  border: `1.5px solid ${
                    task.done ? token.colorSuccess : token.colorTextTertiary
                  }`,
                  background: task.done ? token.colorSuccess : "transparent",
                }}
              />
              <div style={{ flex: 1, minWidth: 0 }}>
                <Text
                  delete={task.done}
                  type={task.done ? "secondary" : undefined}
                  style={{ fontSize: 13 }}
                >
                  {task.title}
                </Text>
                <div style={{ fontSize: 11.5, color: token.colorTextTertiary, marginTop: 2 }}>
                  {task.assignee} · {task.due}
                </div>
              </div>
            </div>
          </List.Item>
        )}
      />
    </Card>
  );
}
