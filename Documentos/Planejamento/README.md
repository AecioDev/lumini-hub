# Planejamento & Metodologia de Desenvolvimento - Lumini Hub

Esta pasta centraliza os documentos de planejamento e histórico de desenvolvimento do **Lumini Hub**. 
Aqui estruturamos as tarefas, metas de cada funcionalidade e diários de bordo para manter a coesão das entregas, independentemente de quem (humano ou agente de IA) esteja executando as tarefas.

## 📂 Estrutura de Diretórios

- **[Historico/](file:///c:/Projetos/lumini-hub/Documentos/Planejamento/Historico/)**: Diário de bordo detalhando as atividades executadas por dia/período de desenvolvimento.
- **[Modulo_0_Configuracao/](file:///c:/Projetos/lumini-hub/Documentos/Planejamento/Modulo_0_Configuracao/)**: Setup do projeto frontend e backend, e o cadastro base de Empresa (multi-empresa self-referencing, config fiscal e plano de contas) que serve de fundação pros demais módulos.
- **[Modulo_1_CRM_Integracoes/](file:///c:/Projetos/lumini-hub/Documentos/Planejamento/Modulo_1_CRM_Integracoes/)**: Middleware (Loja Integrada $\leftrightarrow$ SQL Server) e o CRM do Vendedor.
- **[Modulo_2_Produtos_Estoque/](file:///c:/Projetos/lumini-hub/Documentos/Planejamento/Modulo_2_Produtos_Estoque/)**: Catálogo de Produtos, Imagens e Localização/Saldos de Estoque.
- **[Modulo_3_Compras/](file:///c:/Projetos/lumini-hub/Documentos/Planejamento/Modulo_3_Compras/)**: Pedidos de Compras e Manifesto Eletrônico da SEFAZ.
- **[Modulo_4_Vendas_Caixas/](file:///c:/Projetos/lumini-hub/Documentos/Planejamento/Modulo_4_Vendas_Caixas/)**: Vendas de Balcão (Ambientes/Obra), Projetos CAD, Faturamento Parcial e Controle de Caixas.
- **[Modulo_5_Financeiro/](file:///c:/Projetos/lumini-hub/Documentos/Planejamento/Modulo_5_Financeiro/)**: Contas a Pagar/Receber, Bancos, Caixa e Emissão de Boletos.
- **[Modulo_6_Fiscal/](file:///c:/Projetos/lumini-hub/Documentos/Planejamento/Modulo_6_Fiscal/)**: Certificados Digitais A1, Parametrizações Tributárias e Emissão de NFe/NFSe/NFCe/MDFe via ACBr.
- **[Modulo_7_Contabilidade/](file:///c:/Projetos/lumini-hub/Documentos/Planejamento/Modulo_7_Contabilidade/)**: Plano de Contas, Lançamentos Contábeis de Dupla Partida, Balancete, Razão e DRE.

## 📜 Regras de Uso
1. **Antes de Iniciar uma Feature:** Crie ou atualize o `plano_*.md` e `tasks_*.md` na pasta do respectivo módulo detalhando as regras.
2. **Durante a Execução:** Vá marcando as tarefas como concluídas `[x]` no arquivo `tasks_*.md`.
3. **Ao Finalizar a Sessão:** Adicione um arquivo de log na pasta `Historico/log_YYYY-MM-DD.md` registrando as alterações realizadas, problemas encontrados e o que está pendente.
