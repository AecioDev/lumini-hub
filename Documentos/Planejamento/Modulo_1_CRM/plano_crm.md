# Módulo 1 - CRM

Este módulo estrutura o microsserviço de **CRM** (`api.crm`) do Lumini Hub: Leads, Pipeline (Kanban), Oportunidades, Atividades, Pós-Venda, Churn e Ciclo de Vida do Produto.

> **Separado de Integrações em 2026-09-13**: este módulo nasceu junto com `api.integrations` (Loja Integrada ↔ SQL Server legado) num único `Modulo_1_CRM_Integracoes`, por prioridade da época — mas são dois microsserviços com ciclos de vida bem diferentes (Integrações já tem a Fase 1 implementada e rodando; CRM ainda não tem nenhuma linha de código real, só o Dashboard mockado no frontend, ver `Frontend/src/services/crm/crm-dashboard-service.ts`). Separados a pedido do usuário. Ver `Modulo_9_Integracoes/plano_integrations.md` pra tudo que era a parte de sincronismo/legado.

## 🛠️ Especificações Gerais

1. **api.crm (Porta 4009):** Controle de Leads, Pipelines (Kanban), Oportunidades, Atividades, Pós-Venda, Churn e Ciclo de Vida do Produto.
2. **Migração de dados do legado pra enriquecer o CRM**: carga do histórico de clientes/transações do SQL Server legado — reaproveita a camada de leitura do legado já construída em `api.integrations` (ver `Modulo_9_Integracoes/plano_integrations.md`), mas ainda sem escopo formal aqui.
