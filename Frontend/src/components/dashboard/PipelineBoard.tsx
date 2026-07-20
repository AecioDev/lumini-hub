import { Badge, Card, Typography, theme } from "antd";
import type { PipelineStage } from "@/types/crm";

const { Text } = Typography;

export function PipelineBoard({ stages }: { stages: PipelineStage[] }) {
  const { token } = theme.useToken();

  return (
    <Card title="Pipeline de Vendas">
      <div style={{ display: "grid", gridTemplateColumns: "repeat(4, 1fr)", gap: 12 }}>
        {stages.map((stage) => (
          <div key={stage.name}>
            <div
              style={{
                display: "flex",
                alignItems: "center",
                justifyContent: "space-between",
                marginBottom: 10,
              }}
            >
              <Text style={{ fontSize: 12, fontWeight: 500 }} type="secondary">
                {stage.name}
              </Text>
              <Badge
                count={stage.count}
                color={token.colorFillTertiary}
                style={{ color: token.colorTextTertiary, boxShadow: "none" }}
              />
            </div>

            <div style={{ display: "flex", flexDirection: "column", gap: 8 }}>
              {stage.leads.map((lead) => (
                <Card key={lead.id} size="small" styles={{ body: { padding: 10 } }}>
                  <div style={{ fontSize: 12.5, fontWeight: 500, marginBottom: 2 }}>
                    {lead.name}
                  </div>
                  <div
                    style={{
                      fontSize: 11,
                      color: token.colorTextTertiary,
                      marginBottom: 6,
                    }}
                  >
                    {lead.company}
                  </div>
                  <div style={{ fontSize: 11.5, color: token.colorPrimary, fontWeight: 500 }}>
                    {lead.value}
                  </div>
                </Card>
              ))}
            </div>
          </div>
        ))}
      </div>
    </Card>
  );
}
