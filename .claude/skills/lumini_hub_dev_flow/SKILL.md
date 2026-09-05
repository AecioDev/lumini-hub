---
name: lumini-hub-dev-flow
description: Metodologia de desenvolvimento do Lumini Hub. Regras obrigatórias para leitura do planejamento, histórico de logs e o quadro Kanban/Scrum orquestrado que rastreia o que está em backlog, em andamento, finalizado, em teste, testado por usuário e entregue.
---

# Metodologia de Desenvolvimento e Fluxo de Trabalho - Lumini Hub

Esta Skill descreve as regras obrigatórias de leitura e escrita do planejamento que **todos os agentes de IA** e desenvolvedores devem seguir ao atuar no monorepo do **Lumini Hub**.

---

## 🔍 1. Protocolo de Entrada (Antes de Começar)

Antes de iniciar qualquer alteração no código fonte ou na infraestrutura, o agente **deve obrigatoriamente**:
1. Ler o arquivo explicativo geral em [Documentos/Planejamento/README.md](file:///c:/Projetos/lumini-hub/Documentos/Planejamento/README.md).
2. Localizar a pasta do módulo em que irá atuar (ex: `Modulo_1_CRM_Integracoes/`).
3. Ler o arquivo de planejamento técnico `plano_*.md` da feature correspondente para entender o escopo do banco de dados, regras de negócio e arquivos afetados.
4. Ler o arquivo de controle de tarefas `tasks_*.md` correspondente ao módulo para verificar o andamento e o status atual.

---

## 🗂️ 2. Modo atual: Kanban orquestrado (operador único)

Desde 2026-09-04 o projeto é tocado por **um único operador humano** (sem múltiplos agentes programando em paralelo neste momento), com três subagentes dedicados a esse fluxo:

- `.claude/agents/cocriador-lumini-hub.md` — recebe uma ideia crua de feature, conversa com o usuário (perguntas uma de cada vez, pesquisa na web quando depende de sistema externo) e grava o escopo validado em `plano_*.md`. Nunca decompõe em backlog nem escreve código.
- `.claude/agents/orquestrador-lumini-hub.md` — dono do quadro. Pega um `plano_*.md` já escopado e decompõe em Épico → PBI → Tarefa, recomenda a próxima tarefa e verifica o que foi entregue. Nunca inventa regra de negócio nem escreve código.
- `.claude/agents/revisor-codigo-lumini-hub.md` — acionado pelo orquestrador (ou direto pelo usuário) pra checar se o código de uma tarefa reportada como pronta segue os padrões abaixo. Só lê e relata, nunca edita.

Veja os três arquivos para o fluxo completo.

Toda **feature criada a partir de 2026-09-04** usa este esquema em `tasks_*.md` do módulo correspondente:

- **IDs**: `<PREFIXO>-<N>` para Épico (`CFG`/`CRM`/`EST`/`COM`/`VND`/`FIN`/`FIS`/`CTB`, um por Módulo 0-7), `.N` para PBI, `.N` de novo para Tarefa (ex.: `FIS-1.2.1`).
- **Status é uma tag entre colchetes** logo após o checkbox, uma das seis colunas do quadro: `[BACKLOG]` → `[EM-ANDAMENTO]` → `[FINALIZADO]` → `[EM-TESTE]` → `[TESTADO-USUARIO]` → `[ENTREGUE]` (só esta última usa `[x]`). Se você (agente numa sessão de código) pegar uma tarefa recomendada pelo orquestrador, é esperado que você mesmo avance a tag pra `[EM-ANDAMENTO]` ao começar e `[FINALIZADO]` ao terminar — não precisa esperar o usuário relatar isso de volta ao orquestrador para essas duas transições intermediárias. As transições `[EM-TESTE]` → `[TESTADO-USUARIO]` → `[ENTREGUE]` são do orquestrador, porque dependem de revisão de código e confirmação humana de teste real.
- Status de Épico/PBI nunca é escrito à mão — é sempre a etapa mais atrasada entre as tarefas-filhas.

**`tasks_empresa.md` e `tasks_integrations.md` continuam no esquema antigo de trava** (`[EM EXECUÇÃO POR:]`/`[CONCLUÍDO POR:]`, seção abaixo) — não foram migrados retroativamente, e não há necessidade de migrar salvo pedido explícito.

Isso pode ser revisitado se o projeto voltar a ter múltiplos agentes programando ao mesmo tempo — nesse caso, o protocolo de bloqueio concorrente abaixo (mantido só como referência histórica) volta a valer.

### Protocolo de Bloqueio Concorrente (legado, não usado no momento)

Sistema de **bloqueio textual** no arquivo `tasks_*.md`, usado quando múltiplos agentes trabalhavam em paralelo no mesmo módulo:

- **Tarefa Ocupada:** linha com `[EM EXECUÇÃO POR: AGENTE_ID]` ou `[BLOQUEADO POR: AGENTE_ID]` — não inicie, escolha outra livre (`- [ ]`).
- **Bloqueio de Tarefa:** tarefa livre vira `- [/] [EM EXECUÇÃO POR: NOME_DO_SEU_AGENTE] Descrição da tarefa`, salvo **antes** de tocar em código.
- **Sucesso:** `- [x] [CONCLUÍDO POR: NOME_DO_SEU_AGENTE] Descrição da tarefa`.
- **Interrupção:** volta para `- [ ] Descrição da tarefa`.

---

## 📝 3. Protocolo de Saída (Ao Finalizar o Trabalho)

Ao encerrar o seu turno ou antes de parar a execução do agente, deve-se atualizar a documentação de histórico físico:
1. Criar ou editar o arquivo de log do dia em [Documentos/Planejamento/Historico/log_YYYY-MM-DD.md](file:///c:/Projetos/lumini-hub/Documentos/Planejamento/Historico/) (substituindo `YYYY-MM-DD` pela data local atual).
2. Escrever o registro contendo:
   - **Atividades Realizadas:** O que foi codificado e concluído com sucesso.
   - **Gargalos/Pendências:** Dificuldades encontradas, tarefas bloqueadas ou o que deve ser continuado na próxima sessão.
3. Atualizar as tarefas concluídas no respectivo `tasks_*.md` do módulo.
