import { ArrowDownOutlined, ArrowUpOutlined } from "@ant-design/icons";
import { Card, Statistic, Tag } from "antd";
import type { DashboardStat } from "@/types/crm";

export function StatCard({ stat }: { stat: DashboardStat }) {
  const isUp = stat.trend === "up";
  const deltaText = stat.delta.replace(/^[↑↓]\s*/, "");

  return (
    <Card>
      <Statistic
        title={stat.label}
        value={stat.value}
        valueStyle={{ fontSize: 24, fontWeight: 600 }}
      />
      <Tag color={isUp ? "success" : "error"} style={{ marginTop: 10 }}>
        {isUp ? <ArrowUpOutlined /> : <ArrowDownOutlined />} {deltaText}
      </Tag>
    </Card>
  );
}
