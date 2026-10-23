---
name: lumini-hub-dev-flow
description: Como trabalhar no repositório do Lumini Hub. Quem conduz é o usuário; sem agentes, sem quadro Kanban. Use para saber onde fica o planejamento, quando lê-lo e o que não fazer sozinho.
---

# Como trabalhar no Lumini Hub

O usuário conduz o projeto. Cada tarefa nasce numa conversa com ele e nada roda por conta própria.

## Regras

1. **Não crie subagentes** (Agent/Task) a menos que o usuário peça.
2. **Não mexa nos `tasks_*.md` antigos nem escreva `Historico/log_*.md`** a menos que o usuário peça. Os status antigos (`[EM-TESTE]`, `[ENTREGUE]`, etc.) são histórico, não um fluxo a manter.
   **Controle do trabalho = 2 arquivos** em `Documentos/Planejamento/`: `BACKLOG.md` (o que falta; cada item tem ID `PREFIXO-NNN`, ex.: `CFG-014`, com o prefixo do módulo vindo da tabela no topo do arquivo, número sequencial por prefixo e nunca reaproveitado: todo item novo recebe o próximo ID livre e o ID é citado no `PROXIMA_VERSAO.md`; só o usuário marca `[x]`, depois de testar) e `PROXIMA_VERSAO.md` (o que mudou desde o último envio ao servidor). Ao terminar uma tarefa pedida: acrescente uma linha em `PROXIMA_VERSAO.md` e mostre-a; no `BACKLOG.md` mantenha `[ ]` com a nota `(implementado, aguardando teste)`. Quando o usuário pedir o release: leia `PROXIMA_VERSAO.md`, atualize as constantes `Required` (só se houve mudança de banco) e `AppVersion` em `Backend/common/dbversion/dbversion.go`, gere `releases/Release_<versão>.md` e zere `PROXIMA_VERSAO.md`. O release segue o formato de `releases/Release_2026.2.2.0.md` (um app do portal lê esses arquivos, então a estrutura não muda): título `# Lumini Hub — Release <versão>`; cabeçalho com versão do sistema, versão do banco e tema; data de geração e legenda ✅ testado · ⚠️ implementado, ainda não testado · ⏳ pendente; seções 1 Banco de dados, 2 Novidades, 3 Correções, 4 Infraestrutura, 5 Configuração necessária no servidor, 6 Arquivos principais alterados / criados, 7 Pendências conhecidas (tabela). Os ✅/⚠️ vêm do que o usuário realmente testou, nunca de suposição. Versão: `ano.versao_do_banco.sequencia_versao.sequencia_patch`; o 2º número é a versão do banco e sobe uma unidade por publicação que traz mudanças de banco (não por mudança individual), e todo release lista o que mudou no banco. Próxima: `2026.2.2.0`.
3. **Faça o que foi pedido, no escopo pedido.** Se enxergar algo a mais (bug, dívida, melhoria), avise em uma linha e pergunte; não faça sozinho.
4. **Mostre o que fez**: arquivos alterados e o resultado dos testes, em português, de forma curta.
5. **Nunca derrube nem reinicie o backend** do usuário (`run_services.bat`, `go run`); peça a ele.
6. **Navegador e teste visual só se o usuário pedir** (ver `agent-browser`). Primeiro leia código, logs e testes.

## Onde fica o planejamento

`Documentos/Planejamento/` tem uma pasta por módulo (`Modulo_0_Configuracao` ... `Modulo_9_Integracoes`):

- `plano_*.md`: escopo, banco, regras de negócio, arquivos afetados.
- `tasks_*.md`: lista de tarefas e status.
- `Historico/`: diário de sessões antigas.

Leia só quando o usuário apontar o módulo/feature, ou quando precisar da regra de negócio da tarefa. Alguns arquivos são muito grandes (`tasks_empresa.md` tem cerca de 330 KB): leia apenas a seção necessária, nunca o arquivo todo.

## Padrões do código

Estão no `CLAUDE.md` (arquitetura, padrões obrigatórios do backend e do frontend) e nas skills `lumini_hub_backend_architecture` e `lumini_hub_entity_creation`.
