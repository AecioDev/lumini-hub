# Onboarding de Primeiro Acesso - Cadastro da Primeira Empresa

> Escopado em 2026-10-01 (cocriador, decisões fechadas com o usuário no chat). Mora em `Modulo_0_Configuracao/` ao lado de `plano_empresa.md` por ser extensão direta do cadastro de Empresa (fundação), não um módulo de negócio isolado. Ainda **não decomposto em backlog** — isso é do `orquestrador-lumini-hub`.

## Problema (descoberto no deploy na OCI, 2026-10-01)

Num banco novo, com **zero empresas**, o login do usuário master (sem `users.company_id`) trava num seletor "Escolha uma empresa" vazio:

- `ResolveActiveCompany` (`Backend/common/utils/company_visibility.go`, linhas 88-128): usuário `unrestricted` com 0 empresas ativas cai no `return nil, true, nil` (linha 127), isto é, `requiresSelection=true` com lista vazia.
- `ProtectedRoute.tsx` (linha 29) renderiza `CompanySelectionGate` no lugar do `<Outlet/>`, bloqueando qualquer rota, inclusive `/settings/companies/create`, que fica atrás do mesmo `ProtectedRoute`.
- Resultado: impossível sair do impasse pela UI. Foi contornado inserindo uma empresa por SQL.

## Objetivo

No primeiro acesso a um banco sem nenhuma empresa, o usuário autorizado consegue **cadastrar a empresa direto na tela de seleção** e, ao salvar, entrar no sistema. É a empresa do próprio usuário que acessa; **não é multi-tenant** (ver `Modulo_8_MultiTenant/plano_multitenant.md`, futuro, não antecipado aqui).

## Abordagem escolhida

Descartadas pelo usuário: tela/wizard de onboarding novo, flag novo no backend (`requires_company_setup`), rota liberada no `ProtectedRoute`, criação na tela de login. Escolhida: **alternar por estado local dentro do `CompanySelectionGate`**, sem rota nova e **sem mudança no backend**.

## Regras de negócio

1. **Quando o botão "Cadastrar Empresa" aparece:** somente quando **não existe nenhuma empresa cadastrada** (`user.visible_companies` vazio no Gate), **independente de quem está logado**, **e** o usuário tem a permission `companies.create` (checada via `hasPermission("companies.create")` do `AuthContext`; ADMIN e DEVELOP já passam pelo bypass existente). Com 1 ou mais empresas, o botão **não aparece** e vale o fluxo normal (seletor para 2+, auto-seleção para 1).
2. **Sem empresa e sem permissão:** o Gate mostra "Nenhuma empresa cadastrada, peça a um administrador", com botão para **sair** (`logout()` do `AuthContext`). Sem Select vazio.
3. **Ao clicar no botão:** o Gate troca (estado local `creating`) o seletor pelo `CompanyForm` existente, no mesmo shell centralizado, sem sidebar, com o card mais largo. Sem navegação de rota.
4. **Cancelar/Voltar** no formulário volta ao estado inicial do Gate (estado "sem empresa", com o botão). `CompanyForm` já recebe `onCancel` por prop.
5. **Ao salvar com sucesso:** `companyService.create(...)` seguido de `refreshUser()` (`GET /auth/me`). Com exatamente 1 empresa ativa, a regra **já existente** de auto-seleção do backend (`company_visibility.go`, linhas 108-110) devolve `requires_company_selection=false` e o Gate some sozinho: o usuário **entra direto no sistema**, sem passar pelo seletor. Não muda o backend.
6. **Erro ao salvar:** exibe a mensagem via `useFeedback`/`getApiErrorMessage` (mesmo padrão de `CreateCompanyPage`) e permanece no formulário.
7. **Campo "Vinculada a" (`parent_id`):** o `CompanyForm` carrega as empresas via `companyService.list()` (`GET /companies`, exige `companies.view`). Já tem tratamento gracioso (`.catch(() => setCompanies([]))`). Na primeira empresa a lista é vazia de qualquer forma, então o campo não oferece opções e a empresa nasce como raiz (`ParentID == nil`, "Matriz" por definição de `plano_empresa.md`). Ponto a **verificar na implementação**: garantir que a falha por falta de `companies.view` não gere toast de erro nem trave o formulário.
8. **Segurança:** o botão é conveniência de UI; a criação continua protegida no backend por `RequirePermission("companies.create")` em `POST /companies`.
9. **Escopo da empresa criada:** o usuário continua "master" (sem `users.company_id`). Nenhum vínculo é gravado em `users`.

## Dado / DB

Nenhuma mudança. Sem migration, sem novo campo, sem novo endpoint, sem nova permission. Usa `companies` e `POST /companies` como estão.

## Arquivos afetados (todos frontend)

- `Frontend/src/components/companies/CompanySelectionGate.tsx`: estado `creating`; botão "Cadastrar Empresa" (regras 1-2); estado vazio com texto e botão sair; renderização do `CompanyForm` no shell do Gate (regras 3-6).
- `Frontend/src/contexts/AuthContext.tsx`: novo `refreshUser()` (chama `authService.me()` e faz `setUser`), exposto no `AuthContextValue`.
- `Frontend/src/pages/settings/companies/CreateCompanyPage.tsx` e `CompanyForm`: **extrair a lógica de submit** (`companyService.create` + feedback) para um componente/hook reutilizável com `onSuccess`/`onCancel` injetados. Hoje a página faz `navigate("/settings/companies")` ao salvar e cancelar, o que não serve dentro do Gate. `CreateCompanyPage` passa a usar a versão extraída (comportamento dela não muda).
- `Frontend/src/components/companies/forms/CompanyForm.tsx`: só se a verificação da regra 7 exigir ajuste.
- `Frontend/src/routes/ProtectedRoute.tsx` e `router.tsx`: **não mudam**.
- Backend: **nada**.

## Fora de escopo

- Multi-tenant (Módulo 8), qualquer isolamento entre clientes.
- Wizard de primeiro acesso em várias etapas (empresa, depois fiscal, depois usuários).
- Vínculo do usuário à empresa criada (`users.company_id`); o admin segue master.
- Config fiscal, certificado e logo da empresa no onboarding (ficam para a tela normal `/settings/companies`).
- Mudança de backend: nenhum flag novo, nenhuma alteração em `ResolveActiveCompany`.
- Forçar o seletor com 1 empresa só.
- Criação do usuário admin num banco novo (ver pendências).

## Pendências em aberto

- [ ] **Como o usuário admin nasce num banco novo.** Não há seeder do admin nem das permissions `companies.*` em `api.auth/internal/seeder` (só o menu "Empresas"). Confirmar com o usuário como isso é feito hoje (SQL manual? outro seeder?) para documentar no roteiro de deploy. Não bloqueia este escopo, já que ADMIN/DEVELOP passam pelo bypass.
- [ ] Documentar o roteiro de primeiro acesso (admin, depois login, depois cadastro da empresa) na doc de deploy da OCI, se existir.
