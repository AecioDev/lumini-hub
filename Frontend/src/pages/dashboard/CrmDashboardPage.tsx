import { useEffect, useState } from "react";
import { Button, Col, Row, Skeleton, Typography } from "antd";
import { crmDashboardService } from "@/services/crm/crm-dashboard-service";
import type { CrmDashboardData } from "@/types/crm";
import { StatCard } from "@/components/dashboard/StatCard";
import { PipelineBoard } from "@/components/dashboard/PipelineBoard";
import { TasksList } from "@/components/dashboard/TasksList";
import { RecentClientsTable } from "@/components/dashboard/RecentClientsTable";
import { useFeedback } from "@/hooks/useFeedback";

const { Title, Paragraph } = Typography;

export function CrmDashboardPage() {
  const [data, setData] = useState<CrmDashboardData | null>(null);
  const [loading, setLoading] = useState(true);
  const feedback = useFeedback();

  useEffect(() => {
    let cancelled = false;
    crmDashboardService
      .getDashboardData()
      .then((result) => {
        if (!cancelled) setData(result);
      })
      .finally(() => {
        if (!cancelled) setLoading(false);
      });
    return () => {
      cancelled = true;
    };
  }, []);

  if (loading || !data) {
    return <Skeleton active paragraph={{ rows: 8 }} />;
  }

  return (
    <div>
      <div
        style={{
          display: "flex",
          alignItems: "center",
          justifyContent: "space-between",
          marginBottom: 20,
          flexWrap: "wrap",
          gap: 12,
        }}
      >
        <div>
          <Title level={4} style={{ margin: "0 0 4px" }}>
            Dashboard CRM
          </Title>
          <Paragraph type="secondary" style={{ margin: 0 }}>
            Visão geral de vendas, leads e desempenho da equipe
          </Paragraph>
        </div>
        <div style={{ display: "flex", gap: 10 }}>
          <Button>Mensal</Button>
          <Button
            type="primary"
            onClick={() => feedback.info("Exportação ainda não implementada.")}
          >
            Exportar
          </Button>
        </div>
      </div>

      <Row gutter={[16, 16]}>
        {data.stats.map((stat) => (
          <Col xs={24} sm={12} lg={6} key={stat.label}>
            <StatCard stat={stat} />
          </Col>
        ))}
      </Row>

      <Row gutter={[16, 16]} style={{ margin: "16px 0" }}>
        <Col xs={24} lg={15}>
          <PipelineBoard stages={data.pipeline} />
        </Col>
        <Col xs={24} lg={9}>
          <TasksList tasks={data.tasks} />
        </Col>
      </Row>

      <RecentClientsTable clients={data.recentClients} />
    </div>
  );
}
