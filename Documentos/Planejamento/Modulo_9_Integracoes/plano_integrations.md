# Módulo 9 - Integrações (Loja Integrada ↔ SQL Server Legado)

Este módulo cobre o `api.integrations` (porta 4007) — middleware de sincronismo bidirecional entre a **Loja Integrada** (e-commerce) e o **SQL Server** legado do cliente (`FOCCO_ERP`), que o Lumini Hub está substituindo aos poucos.

> **Separado do CRM em 2026-09-13**: até então vivia junto do `api.crm` num único `Modulo_1_CRM_Integracoes`, por prioridade da época — mas Integrações é um microsserviço à parte, já com Fase 1 implementada e rodando, sem relação de código com o CRM. Ver `Modulo_1_CRM/plano_crm.md` pro CRM. O detalhe técnico completo (schema legado, endpoints, fluxo de estoque Oficial → Reserva → cliente) continua documentado em `CLAUDE.md` § "Legacy SQL Server Integration" — este arquivo não duplica, só referencia e complementa com o que falta.

## 🛠️ Especificações Gerais

1. **api.integrations (Porta 4007)**: sincronismo bidirecional de pedidos e estoque entre Loja Integrada e SQL Server legado. **Fase 1** (esqueleto, conexões Postgres/legado, `ConfigService`, handlers/rotas, proxy no gateway, tela Configurações → Integrações) entregue — ver `tasks_integrations.md`. **Fase 2** (domain `NOTAS`/`NOTAS1`/`NOTAS3`/`CADPRO`, cliente HTTP da Loja Integrada, scheduler de polling do `LogAltera`, processamento de `WebhookEvent` em pedido de venda) ainda não implementada — bloqueada em pendências externas: chaves de API/App da Loja Integrada, valor de `codtipnot` a usar pra pedidos originados na LI, e a política de prioridade de estoque loja física vs. virtual quando os dois canais vendem a última unidade ao mesmo tempo (a confirmar com a diretoria do cliente).
2. **Migração de dados do legado**: carga do histórico de clientes e transações do SQL Server legado pro Postgres novo — reaproveita a camada de leitura do legado já construída aqui na Fase 1, mas ainda sem escopo formal. Cruza com o CRM (`Modulo_1_CRM/plano_crm.md`), já que é o CRM quem consome esses dados migrados.
