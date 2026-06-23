# Módulo 1 - CRM & Integrações Legadas

Este módulo gerencia a comunicação e sincronização de dados entre a **Loja Integrada** (via API REST) e o **SQL Server** legado do cliente, além da estruturação do microsserviço de **CRM** (`api.crm`) alimentado por migrações incrementais do legado.

## 🛠️ Especificações Gerais

1. **api.integrations (Porta 4007):** Sincronismo bidirecional de Pedidos e Estoque.
2. **api.crm (Porta 4009):** Controle de Leads, Pipelines (Kanban), Oportunidades, Atividades, Pós-Venda, Churn e Ciclo de Vida do Produto.
3. **Migração de Dados:** Carga do histórico contendo dados de clientes e transações do SQL Server legado para enriquecer o banco PostgreSQL novo.
