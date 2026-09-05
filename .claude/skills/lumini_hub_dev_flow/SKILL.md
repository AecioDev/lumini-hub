---
name: lumini-hub-dev-flow
description: Metodologia de desenvolvimento do Lumini Hub. Regras obrigatórias para leitura do planejamento, histórico de logs e sistema de bloqueio de tarefas para desenvolvimento cooperativo entre múltiplos agentes.
---

# Metodologia de Desenvolvimento e Fluxo de Trabalho - Lumini Hub

Esta Skill descreve as regras obrigatórias de leitura, escrita e bloqueio concorrente de tarefas que **todos os agentes de IA** e desenvolvedores devem seguir ao atuar no monorepo do **Lumini Hub**.

---

## 🔍 1. Protocolo de Entrada (Antes de Começar)

Antes de iniciar qualquer alteração no código fonte ou na infraestrutura, o agente **deve obrigatoriamente**:
1. Ler o arquivo explicativo geral em [Documentos/Planejamento/README.md](file:///c:/Projetos/lumini-hub/Documentos/Planejamento/README.md).
2. Localizar a pasta do módulo em que irá atuar (ex: `Modulo_1_CRM_Integracoes/`).
3. Ler o arquivo de planejamento técnico `plano_*.md` da feature correspondente para entender o escopo do banco de dados, regras de negócio e arquivos afetados.
4. Ler o arquivo de controle de tarefas `tasks_*.md` correspondente ao módulo para verificar o andamento e o status atual.

---

## 🔒 2. Protocolo de Bloqueio Concorrente (Multi-Agente)

Para permitir que múltiplos agentes e programadores trabalhem de forma cooperativa em paralelo sem gerar conflitos de código ou duplicidade de esforço, adota-se o seguinte sistema de **bloqueio textual** no arquivo `tasks_*.md` do módulo:

### Ações para iniciar uma tarefa:
1. O agente analisa a lista de tarefas no `tasks_*.md`.
2. **Tarefa Ocupada:** Se a linha da tarefa contiver a tag `[EM EXECUÇÃO POR: AGENTE_ID]` ou `[BLOQUEADO POR: AGENTE_ID]`, o agente **NÃO DEVE** iniciar este desenvolvimento. Deve escolher outra tarefa que esteja livre (`- [ ]`).
3. **Bloqueio de Tarefa:** Se a tarefa estiver livre (`- [ ]`), o agente deve editá-la imediatamente no markdown do repositório físico, alterando a linha para:
   ```markdown
   - [/] [EM EXECUÇÃO POR: NOME_DO_SEU_AGENTE] Descrição da tarefa
   ```
4. **Persistência do Bloqueio:** O agente deve salvar o arquivo `tasks_*.md` **antes** de realizar qualquer alteração em arquivos de código do backend ou do frontend.

### Ações ao concluir ou interromper:
- **Sucesso:** Ao finalizar o desenvolvimento daquela tarefa com sucesso, o agente a marca como concluída e remove a tag de execução (ou a substitui por conclusão):
  ```markdown
  - [x] [CONCLUÍDO POR: NOME_DO_SEU_AGENTE] Descrição da tarefa
  ```
- **Interrupção:** Se o agente falhar em terminar ou for interrompido, ele deve remover o bloqueio e devolver a tarefa ao estado livre:
  ```markdown
  - [ ] Descrição da tarefa
  ```

---

## 📝 3. Protocolo de Saída (Ao Finalizar o Trabalho)

Ao encerrar o seu turno ou antes de parar a execução do agente, deve-se atualizar a documentação de histórico físico:
1. Criar ou editar o arquivo de log do dia em [Documentos/Planejamento/Historico/log_YYYY-MM-DD.md](file:///c:/Projetos/lumini-hub/Documentos/Planejamento/Historico/) (substituindo `YYYY-MM-DD` pela data local atual).
2. Escrever o registro contendo:
   - **Atividades Realizadas:** O que foi codificado e concluído com sucesso.
   - **Gargalos/Pendências:** Dificuldades encontradas, tarefas bloqueadas ou o que deve ser continuado na próxima sessão.
3. Atualizar as tarefas concluídas no respectivo `tasks_*.md` do módulo.
