# Planejamento - Lumini Hub

Esta pasta guarda o planejamento do **Lumini Hub**: o que foi feito, o que falta, o que vai em cada versão e o escopo de cada módulo.

## O que olhar primeiro

- **`BACKLOG.md`**: o painel. Tudo o que está feito e testado e tudo o que falta, cada item com um ID (`CFG-014`, `VND-001`...).
- **`PROXIMA_VERSAO.md`**: o que mudou desde o último envio ao servidor (base do próximo release).
- **`releases/`**: um arquivo `Release_<versão>.md` por versão publicada, lido pelo portal de releases.

## 📂 Estrutura de Diretórios

- **[Historico/](file:///c:/Projetos/lumini-hub/Documentos/Planejamento/Historico/)**: Diário de bordo detalhando as atividades executadas por dia/período de desenvolvimento.
- **[Arquivo/](file:///c:/Projetos/lumini-hub/Documentos/Planejamento/Arquivo/)**: `tasks_*.md` de trabalhos já fechados (somente consulta).
- **[releases/](file:///c:/Projetos/lumini-hub/Documentos/Planejamento/releases/)**: um arquivo por versão publicada.
- **[Modulo_0_Configuracao/](file:///c:/Projetos/lumini-hub/Documentos/Planejamento/Modulo_0_Configuracao/)**: Setup do projeto frontend e backend, e o cadastro base de Empresa (multi-empresa self-referencing, config fiscal e plano de contas) que serve de fundação pros demais módulos.
- **[Modulo_1_CRM/](file:///c:/Projetos/lumini-hub/Documentos/Planejamento/Modulo_1_CRM/)**: CRM do Vendedor — Leads, Pipeline (Kanban), Oportunidades, Atividades, Pós-Venda e Churn.
- **[Modulo_2_Produtos_Estoque/](file:///c:/Projetos/lumini-hub/Documentos/Planejamento/Modulo_2_Produtos_Estoque/)**: Catálogo de Produtos, Imagens e Localização/Saldos de Estoque.
- **[Modulo_3_Compras/](file:///c:/Projetos/lumini-hub/Documentos/Planejamento/Modulo_3_Compras/)**: Pedidos de Compras e Manifesto Eletrônico da SEFAZ.
- **[Modulo_4_Vendas_Caixas/](file:///c:/Projetos/lumini-hub/Documentos/Planejamento/Modulo_4_Vendas_Caixas/)**: Vendas de Balcão (Ambientes/Obra), Projetos CAD, Faturamento Parcial e Controle de Caixas.
- **[Modulo_5_Financeiro/](file:///c:/Projetos/lumini-hub/Documentos/Planejamento/Modulo_5_Financeiro/)**: Contas a Pagar/Receber, Bancos, Caixa e Emissão de Boletos.
- **[Modulo_6_Fiscal/](file:///c:/Projetos/lumini-hub/Documentos/Planejamento/Modulo_6_Fiscal/)**: Certificados Digitais A1, Parametrizações Tributárias e Emissão de NFe/NFSe/NFCe/MDFe via ACBr.
- **[Modulo_7_Contabilidade/](file:///c:/Projetos/lumini-hub/Documentos/Planejamento/Modulo_7_Contabilidade/)**: Plano de Contas, Lançamentos Contábeis de Dupla Partida, Balancete, Razão e DRE.
- **[Modulo_8_MultiTenant/](file:///c:/Projetos/lumini-hub/Documentos/Planejamento/Modulo_8_MultiTenant/)**: Estratégia de isolamento entre clientes/tenants (banco de dados físico separado por cliente) — ortogonal ao Módulo 0 (Empresa), que trata da estrutura dentro do banco de um único cliente.
- **[Modulo_9_Integracoes/](file:///c:/Projetos/lumini-hub/Documentos/Planejamento/Modulo_9_Integracoes/)** (**PAUSADO desde 2026-10-01**: acesso ao legado virá de aplicação separada no servidor do cliente; Fase 1 permanece no código, Fase 2 pausada): Middleware de sincronismo bidirecional (Loja Integrada $\leftrightarrow$ SQL Server legado) — separado do Módulo 1 (CRM) em 2026-09-13, eram dois microsserviços diferentes espremidos no mesmo módulo por prioridade da época.

## Como trabalhar

1. A tarefa nasce numa conversa e recebe um ID no `BACKLOG.md`.
2. Ao terminar, o Claude acrescenta uma linha no `PROXIMA_VERSAO.md` e deixa o item `[ ]` com a nota "(implementado, aguardando teste)".
3. Só quem testou marca `[x]`.
4. Na hora de publicar, gera-se o `releases/Release_<versão>.md` e o `PROXIMA_VERSAO.md` é zerado.

O escopo e as regras de negócio de cada módulo ficam nos `plano_*.md`. Os `tasks_*.md` de trabalhos fechados estão em `Arquivo/`; os diários antigos, em `Historico/`.
