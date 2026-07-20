import type { CrmDashboardData } from "@/types/crm";

// ATENÇÃO — dados 100% mockados.
//
// O Dashboard CRM (estatísticas, pipeline de vendas, tarefas da equipe e
// clientes recentes) não tem backend hoje: o microsserviço `api.crm`
// (porta 4009) está apenas planejado, ver
// Documentos/Planejamento/Modulo_1_CRM_Integracoes/plano_crm.md.
//
// Este service fica isolado de propósito para ser fácil de trocar depois:
// quando o api.crm existir, basta reescrever `getDashboardData` para chamar
// os endpoints reais via axios (services/common/api.ts), mantendo o mesmo
// tipo de retorno `CrmDashboardData` (ou ajustá-lo conforme o contrato real).

const MOCK_DASHBOARD_DATA: CrmDashboardData = {
  stats: [
    { label: "Receita", value: "R$ 428.900", delta: "↑ 12,4%", trend: "up" },
    { label: "Novos Leads", value: "86", delta: "↑ 8,1%", trend: "up" },
    { label: "Taxa de Conversão", value: "24,5%", delta: "↓ 1,2%", trend: "down" },
    { label: "Ticket Médio", value: "R$ 3.240", delta: "↑ 4,7%", trend: "up" },
  ],
  pipeline: [
    {
      name: "Novo Lead",
      count: 12,
      leads: [
        { id: 1, name: "Fernanda Lima", company: "Grupo Alfa", value: "R$ 8.200" },
        { id: 2, name: "Bruno Alves", company: "Meridian Corp", value: "R$ 4.500" },
      ],
    },
    {
      name: "Em Contato",
      count: 8,
      leads: [
        { id: 3, name: "Rafael Costa", company: "Vitalis SA", value: "R$ 12.000" },
        { id: 4, name: "Juliana Nunes", company: "Portal Norte", value: "R$ 6.300" },
      ],
    },
    {
      name: "Proposta",
      count: 5,
      leads: [
        { id: 5, name: "Marcos Dias", company: "Ábaco Ltda", value: "R$ 21.500" },
        { id: 6, name: "Patrícia Reis", company: "Fusion Tech", value: "R$ 9.800" },
      ],
    },
    {
      name: "Fechado",
      count: 3,
      leads: [
        { id: 7, name: "Camila Souza", company: "Orion Retail", value: "R$ 34.000" },
        { id: 8, name: "Diego Martins", company: "Sol Nascente", value: "R$ 15.700" },
      ],
    },
  ],
  tasks: [
    {
      id: 1,
      title: "Ligar para lead Grupo Alfa",
      assignee: "Carla Souza",
      due: "Hoje",
      done: false,
    },
    {
      id: 2,
      title: "Enviar proposta Fusion Tech",
      assignee: "Rafael Costa",
      due: "Amanhã",
      done: false,
    },
    {
      id: 3,
      title: "Atualizar contrato Vitalis SA",
      assignee: "Ana Ribeiro",
      due: "20 Jul",
      done: true,
    },
    {
      id: 4,
      title: "Reunião de fechamento Orion",
      assignee: "Bruno Alves",
      due: "22 Jul",
      done: false,
    },
  ],
  recentClients: [
    {
      id: 1,
      name: "Fernanda Lima",
      company: "Grupo Alfa",
      owner: "Carla Souza",
      status: "Lead",
      value: "R$ 8.200",
    },
    {
      id: 2,
      name: "Rafael Costa",
      company: "Vitalis SA",
      owner: "Ana Ribeiro",
      status: "Cliente",
      value: "R$ 12.000",
    },
    {
      id: 3,
      name: "Marcos Dias",
      company: "Ábaco Ltda",
      owner: "Bruno Alves",
      status: "Proposta",
      value: "R$ 21.500",
    },
    {
      id: 4,
      name: "Camila Souza",
      company: "Orion Retail",
      owner: "Carla Souza",
      status: "Cliente",
      value: "R$ 34.000",
    },
  ],
};

export const crmDashboardService = {
  async getDashboardData(): Promise<CrmDashboardData> {
    // pequeno atraso artificial só para exercitar os estados de loading da UI
    await new Promise((resolve) => setTimeout(resolve, 300));
    return MOCK_DASHBOARD_DATA;
  },
};
