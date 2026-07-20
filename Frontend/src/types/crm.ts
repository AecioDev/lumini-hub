// Tipos do Dashboard CRM. Ainda NÃO existe um backend para isso (api.crm,
// porta 4009, está apenas planejado — ver Documentos/Planejamento/Modulo_1_CRM_Integracoes).
// Estes tipos alimentam o mock isolado em services/crm/crm-dashboard-service.ts
// e devem ser trocados pelos DTOs reais assim que o microsserviço existir.

export type TrendDirection = "up" | "down";

export interface DashboardStat {
  label: string;
  value: string;
  delta: string;
  trend: TrendDirection;
}

export interface PipelineLead {
  id: number;
  name: string;
  company: string;
  value: string;
}

export interface PipelineStage {
  name: string;
  count: number;
  leads: PipelineLead[];
}

export interface TeamTask {
  id: number;
  title: string;
  assignee: string;
  due: string;
  done: boolean;
}

export type ClientStatus = "Lead" | "Proposta" | "Cliente";

export interface RecentClient {
  id: number;
  name: string;
  company: string;
  owner: string;
  status: ClientStatus;
  value: string;
}

export interface CrmDashboardData {
  stats: DashboardStat[];
  pipeline: PipelineStage[];
  tasks: TeamTask[];
  recentClients: RecentClient[];
}
