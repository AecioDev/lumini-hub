# Tasks - Onboarding da Primeira Empresa (CompanySelectionGate)

Backlog decomposto em 2026-10-01 pelo orquestrador a partir de `plano_onboarding_empresa.md` (escopo fechado, não alterado aqui). Prefixo `CFG`, esquema novo de Kanban (ver `.claude/skills/lumini_hub_dev_flow/SKILL.md` seção 2). Status de Épico/PBI é sempre a etapa mais atrasada entre as tarefas-filhas. Tudo nasce `[BACKLOG]`.

Relação com `tasks_schema_seed.md` (EPIC CFG-9): **o código deste Épico não depende do seed do admin**. A única ponte é documental (CFG-8.2.1, roteiro de primeiro acesso). Frontend apenas, backend sem mudança.

### EPIC CFG-8: Onboarding da primeira empresa no `CompanySelectionGate`

- [ ] [EM-ANDAMENTO] **EPIC CFG-8**: Num banco com zero empresas, o usuário autorizado (`companies.create`, inclusive ADMIN/DEVELOP pelo bypass) cadastra a empresa direto na tela de seleção e entra no sistema ao salvar; sem permissão vê mensagem e botão sair. Sem rota nova, sem flag, sem mudança de backend. Escopo: `plano_onboarding_empresa.md`.

#### PBI CFG-8.1: Cadastrar a primeira empresa dentro do Gate

- [x] [ENTREGUE] **PBI CFG-8.1**: Ao fim, o impasse "seletor vazio" some: login em banco sem empresas leva ao cadastro e, ao salvar, direto ao sistema.

  - [x] [ENTREGUE] CFG-8.1.1: `AuthContext.tsx` — novo `refreshUser()` (chama `authService.me()` e faz `setUser`), exposto no `AuthContextValue`. (regra 5 do plano)
    - Iniciada em 2026-10-01. Implementação na branch `feature/onboarding-empresa` (criada a partir da `develop`), em worktree `C:\Projetos\lumini-hub-onboarding`.
    - Contexto existente: `AuthContext.tsx` já tem `setUser` local e `setActiveCompany` (que faz `setUser(await authService.setActiveCompany(id))`); `authService.me()` já é usado no bootstrap. Seguir o mesmo padrão com `useCallback`.
    - Critério de aceite:
      - [x] `AuthContextValue` ganha `refreshUser: () => Promise<ApiUserDetail>` (ou `Promise<void>`), implementado com `useCallback` chamando `authService.me()` + `setUser`, e incluído no `useMemo` do value e nas suas dependências.
      - [x] Se `me()` falhar, o erro propaga ao chamador (sem engolir) e `user` não é zerado.
      - [x] Nenhuma mudança de backend nem de outros consumidores; `login`/`logout`/`setActiveCompany` e o bootstrap seguem iguais.
      - [x] `pnpm exec tsc --noEmit` e `pnpm lint` sem novo erro. (único warning react-refresh/only-export-components é anterior)
    - Finalizada em 2026-10-01 (commit 7c3c09f); revisor aprovou sem achados bloqueantes.
    - Fechamento em 2026-10-01: FINALIZADO -> EM-TESTE (revisão aprovada) -> TESTADO-USUARIO (usuário autorizou explicitamente usar o teste pessoal do fluxo completo do Gate em banco vazio, onde `refreshUser()` foi exercitado) -> ENTREGUE.
  - [x] [ENTREGUE] CFG-8.1.2: Extrair a lógica de submit (`companyService.create` + feedback via `useFeedback`/`getApiErrorMessage`) de `CreateCompanyPage.tsx`/`CompanyForm` para componente ou hook reutilizável com `onSuccess`/`onCancel` injetados; `CreateCompanyPage` passa a usá-lo sem mudar seu comportamento (navega para `/settings/companies` ao salvar/cancelar). (regra 6)
    - Iniciada em 2026-10-01, mesma branch `feature/onboarding-empresa` (worktree `C:\Projetos\lumini-hub-onboarding`).
    - Critério de aceite:
      - [x] Submit de criação (`companyService.create` + feedback de sucesso/erro) vive num hook reutilizável que recebe só `onSuccess`, sem `useNavigate` dentro dele. Decisão consciente: `onCancel` não vai no hook; o cancelar fica inline na página porque pertence ao `CompanyForm`, que já recebe `onCancel` por prop.
      - [x] `CreateCompanyPage` usa o hook passando `onSuccess` (e `onCancel` inline) que navegam para `/settings/companies`; comportamento visível idêntico ao atual.
      - [x] Nenhuma mudança de backend nem em `EditCompanyPage`; sem regressão no cadastro/edição de empresa.
      - [x] `tsc -b` e eslint sem novo erro.
    - Finalizada em 2026-10-01 (commits 6befadb + ajuste); revisor aprovou sem bloqueantes. Ressalva do revisor (erro do `onSuccess` aparecendo como "Erro ao criar empresa" após o toast de sucesso) corrigida: o hook separa o try da criação do tratamento do `onSuccess`, com mensagem própria.
    - Fechamento em 2026-10-01: FINALIZADO -> EM-TESTE (revisão aprovada) -> TESTADO-USUARIO (usuário autorizou explicitamente usar o teste pessoal do fluxo completo do Gate, que exercitou o hook de criação) -> ENTREGUE.
  - [x] [ENTREGUE] CFG-8.1.3: `CompanySelectionGate.tsx` — estado local `creating`; com `visible_companies` vazio e `hasPermission("companies.create")`: botão "Cadastrar Empresa" que troca o seletor pelo `CompanyForm` (card mais largo, mesmo shell, sem navegação), Cancelar volta ao estado inicial, sucesso chama `refreshUser()` e o Gate some pela auto-seleção existente do backend; sem permissão: texto "Nenhuma empresa cadastrada, peça a um administrador" + botão sair (`logout()`), sem Select vazio; com 1+ empresas o botão não aparece. (regras 1-5; depende de CFG-8.1.1, CFG-8.1.2)
    - Iniciada em 2026-10-01, mesma branch `feature/onboarding-empresa` (worktree `C:\Projetos\lumini-hub-onboarding`).
    - Finalizada em 2026-10-01 (commit 1edddd4; ajustes das ressalvas 1 e 2 no commit fc3ff33). Revisor aprovou sem bloqueantes. Critérios atendidos: botão "Cadastrar Empresa" troca o seletor pelo `CompanyForm` no mesmo cartão, Cancelar volta, sucesso chama `refreshUser()` e o Gate some; sem permissão, mensagem + botão sair; com 1+ empresas o botão não aparece; `tsc -b` e eslint limpos.
    - Ressalvas do revisor: (1) botão Sair nos três estados do Gate, corrigida; (2) com zero empresas o título vira "Nenhuma empresa cadastrada", sem Select vazio, corrigida; (3) falha do `refreshUser` após criar fica como está, por decisão (a mensagem do hook orienta recarregar).
    - Decisão de UX: com zero empresas o Gate não mostra seletor e "Cadastrar Empresa" é o botão primário.
    - Validação: teste real em banco vazio (`erp_system_onboarding`, local): POST /companies 201 + GET /auth/me 200, sistema entrou direto, CNPJ gravado só em dígitos. O usuário testou pessoalmente visual e funcionamento e confirmou em 2026-10-01. Etapas EM-TESTE/TESTADO-USUARIO não passaram formalmente pelo quadro: mantida em FINALIZADO a pedido; avançar para TESTADO-USUARIO/ENTREGUE exige o usuário confirmar explicitamente que quer fechar com base nessa validação.
    - Fechamento em 2026-10-01: FINALIZADO -> EM-TESTE (revisão aprovada) -> TESTADO-USUARIO (usuário confirmou explicitamente que o teste pessoal em banco vazio vale como esta etapa) -> ENTREGUE após conferência final (`CompanySelectionGate.tsx`, `useCreateCompany.ts` e `refreshUser` em `AuthContext.tsx` existem no worktree da branch `feature/onboarding-empresa`).
    - Observação do revisor (CFG-8.1.1): quem chamar `refreshUser()` (este Gate) deve tratar o erro com try/catch ou `.catch`, para não deixar promise rejeitada sem tratamento.
    - Nota (CFG-8.1.2): o Gate deve passar `onSuccess = refreshUser` (async) ao hook de criação; o hook já trata a falha dele (mensagem própria), então não duplicar o tratamento. Memoizar o callback se for usado como dependência de `useEffect`.
  - [x] [ENTREGUE] CFG-8.1.4: Verificar a regra 7: `CompanyForm` carrega `companyService.list()` (exige `companies.view`); garantir que a falha por falta dessa permissão não gere toast de erro nem trave o formulário no Gate; ajustar `CompanyForm.tsx` só se necessário. (depende de CFG-8.1.3, entregue)
    - Iniciada em 2026-10-01, mesma branch `feature/onboarding-empresa` (worktree `C:\Projetos\lumini-hub-onboarding`).
    - Critério de aceite:
      - [x] Sem `companies.view`, o `CompanyForm` dentro do Gate não dispara toast de erro.
      - [x] Sem `companies.view`, o formulário não trava (continua editável e submetível).
      - [x] O campo "Vinculada a" degrada de forma graciosa (sem opções/desabilitado, sem quebrar).
      - [x] `CompanyForm.tsx` só é alterado se necessário; sem regressão com `companies.view` (cadastro/edição normais). Não foi necessário alterar.
      - [x] `tsc -b` e eslint sem novo erro (nenhum código alterado).
    - Finalizada em 2026-10-01 por verificação estática (leitura de código), sem alteração de código e sem diff para revisar/commitar. Evidências: `CompanyForm.tsx` faz `companyService.list().then(setCompanies).catch(() => setCompanies([]))` ao montar (falha engolida, sem toast); `api.ts` só trata 401 no interceptor, sem toast global; `permission.go` devolve 403 (não 401) sem a permissão, caindo no `.catch`; com lista vazia "Vinculada a" mostra só "Nenhuma — esta é a empresa Matriz" e o submit não depende da lista.
    - Risco residual baixo: NÃO houve teste ao vivo com usuário não-ADMIN que tenha só `companies.create` (sem `companies.view`); o teste do usuário usou ADMIN (bypass). Oferecer teste ao vivo se o usuário quiser; sem ele a tarefa não avança de FINALIZADO.
    - Fechamento em 2026-10-01: sem diff para revisar (EM-TESTE sem revisão de código); usuário autorizou explicitamente fechar pelo teste pessoal (TESTADO-USUARIO), com a ressalva de que esse teste usou ADMIN (bypass) e a 8.1.4 foi verificada só estaticamente -> ENTREGUE. O risco residual baixo permanece registrado acima.

#### PBI CFG-8.2: Roteiro de primeiro acesso documentado

- [ ] [BACKLOG] **PBI CFG-8.2**: O roteiro admin, login, cadastro da empresa fica registrado na doc de deploy. Dependência **apenas documental** do seed do admin.

  - [ ] [BACKLOG] CFG-8.2.1: Documentar o roteiro de primeiro acesso (admin por `BOOTSTRAP_ADMIN_*`, login, cadastro da empresa no Gate) na doc de deploy da OCI, se existir. Pode ser feita junto de CFG-9.5.3 (mesma doc, mesma pendência P5 do `plano_schema_seed.md`). (depende de CFG-8.1.3 e CFG-9.4.1; só documentação, não bloqueia código)

### Pendência herdada do plano (não é tarefa de código)

- Como o admin nasce num banco novo: resolvida pelo EPIC CFG-9 (`tasks_schema_seed.md`); confirmar com o usuário o procedimento usado hoje (SQL manual) apenas se precisar constar no roteiro antes do CFG-9 ficar pronto.
