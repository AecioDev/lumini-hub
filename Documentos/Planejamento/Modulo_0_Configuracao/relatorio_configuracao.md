# Relatório do Módulo 0 (Configuração) — situação em 30/10/2026

Versão do sistema: **2026.3.4.0** · versão do banco: **3** · publicada na VM de testes e conferida.
Números tirados do código da `develop` (commit `4f62226`), não de documentação antiga.

## 1. Primeiro ponto: não existe um "microsserviço de configurações"

O Módulo 0 hoje está espalhado por dois serviços e por uma biblioteca comum:

| Peça | Onde mora | O que tem do Módulo 0 |
|---|---|---|
| **`api.auth`** (4001) | identidade e acesso | login, usuários, perfis, permissões, menu, favoritos, versão do sistema |
| **`api.core`** (4002) | cadastros base | empresas (matriz e filiais), configuração fiscal e certificado A1, logo |
| **`common/`** | biblioteca compartilhada | RBAC, JWT, paginação, versão do banco, config, middlewares |
| **`api.gateway`** (4000) | porta única | proxy por prefixo, CORS, Swagger |

O `api.core` também guarda Clientes, Fornecedores, Endereços, Contatos e Documentos. Isso não é configuração: são os cadastros base de Vendas e Compras. Contei separado abaixo.

## 2. Tabelas

**Módulo 0 em sentido estrito: 10 tabelas.**

| Serviço | Tabela | Para que serve |
|---|---|---|
| api.auth | `users` | usuários (senha com hash, perfil, empresa, ativo) |
| api.auth | `roles` | perfis (modelos de permissões) |
| api.auth | `permissions` | catálogo de permissões (código + módulo) |
| api.auth | `menu_items` | árvore do menu lateral (autorreferenciada) |
| api.auth | `user_menu_favorites` | favoritos do menu, até 5 por usuário (nova na 2026.3.4.0) |
| api.auth | `role_permissions` | junção perfil ↔ permissão |
| api.auth | `user_permissions` | junção usuário ↔ permissão (é esta que o RBAC confere) |
| api.core | `companies` | empresas, matriz e filiais, logo em `bytea` |
| api.core | `company_fiscal_configs` | regime tributário, contador, certificado A1 |
| common | `schema_versions` | versão do banco (gravada pelo `api.auth` a cada release) |

**Cadastros base que também moram no `api.core` (8):** `customers`, `suppliers`, `address`, `contact`, `document`, `country`, `state`, `city`.

**Fora do módulo:** `api.integrations` tem mais 4 (`integration_configs`, `product_mappings`, `sync_logs`, `webhook_events`), pausado.

**Total no banco hoje: 22 tabelas** (10 do Módulo 0, 8 de cadastros base, 4 de integrações).

> **Achado:** a struct `SystemLog` (tabela `system_logs`, auditoria) existe em `common/database/models.go`, mas **não está em nenhum `AutoMigrate`**. A tabela não é criada em banco nenhum.

## 3. Endpoints

**Módulo 0: 46 endpoints.**

**`api.auth` (34)**

| Grupo | Qtd | Rotas |
|---|---|---|
| `/auth` | 7 | `login`, `refresh-token`, `logout`, `me`, `version`, `PUT active-company`, `PUT favorites` |
| `/users` | 8 | listar, `POST /filter`, buscar por id, criar, editar, excluir, `PUT /:id/password`, `PUT /:id/permissions` |
| `/roles` | 6 | listar, buscar, criar, editar, `PUT /:id/permissions`, excluir |
| `/permissions` | 8 | listar, `POST /filter`, buscar, `by-module`, `modules`, criar, editar, excluir |
| `/menu-items` | 5 | árvore, buscar, criar, editar, excluir (só desenvolvedor) |

**`api.core`, parte de configuração (12)**

| Grupo | Qtd | Rotas |
|---|---|---|
| `/companies` | 8 | listar, buscar, criar, editar, excluir, enviar logo, remover logo, `GET /:id/logo` |
| `/company-fiscal-configs` | 4 | `GET /by-company/:companyId`, criar, editar, `POST /:id/certificate` |

**Fora do módulo:** `/customers` 6 e `/suppliers` 6 (`api.core`); `/integrations` 5 (`api.integrations`). **Total geral: 63.**

**Padrões:** 2 `POST .../filter` no `api.auth` (usuários e permissões) com busca no servidor; os demais retornam a lista inteira. Toda rota exige `RequirePermission`, exceto as de `/auth` e `GET /companies/:id/logo` (identidade de sessão, exige só login).

## 4. Permissões e menu

- **67 permissões no catálogo**, 58 vinculadas ao `ADMIN` e 9 do módulo `Develop` (só o `DEVELOP`).
- Por módulo: Vendas 12, Estoque 12, Develop 9, Compras 9, Financeiro 8, Empresas 8, Usuários 4, Integrações 2, e uma cada em Gerencial, Dashboard e Admin.
- Das 67, só as de Usuários, Empresas, Perfis (`roles.*`), Permissões e Admin pertencem ao Módulo 0. As de Vendas, Estoque, Compras e Financeiro já estão no catálogo para os módulos futuros.
- **30 nós** na árvore do menu semeada (grupos e telas).
- Toda permissão exigida por um `RequirePermission` do código está no catálogo (CFG-013).

## 5. Telas do frontend (Módulo 0)

| Tela | Estado |
|---|---|
| Login | real |
| Esqueci a senha | **só UI**, sem endpoint no backend |
| Usuários: lista, cadastro, edição (com permissões individuais) | real |
| Perfis e Permissões: lista, cadastro e edição de perfil | real |
| Permissões (catálogo): lista, cadastro, edição | real, só DEVELOP |
| Empresas: lista, cadastro, edição (matriz/filial, fiscal, certificado, logo) | real |
| Primeira empresa num banco vazio | real |
| Troca de empresa ativa, barra de status com versão, favoritos do menu | real |
| Administração do menu (`/menu-items`) | **sem tela**, só API |

## 6. Tamanho e testes

| | Código (linhas, sem testes) | Linhas de teste | Funções `Test*` |
|---|---|---|---|
| api.auth | 5.010 | 1.240 | 33 |
| api.core | 3.652 | 80 | 2 |
| api.gateway | 183 | 0 | 0 |
| api.integrations | 1.154 | 0 | 0 |
| common | 2.056 | n/d | 30 |
| Frontend (TS/TSX) | 8.414 | sem executor de testes | 0 |

Os testes que usam Postgres só rodam com `TEST_DATABASE_URL`; sem ela são ignorados. Na última rodada (com banco `lumini_test`), `common`, `api.auth` e `api.core/internal/domain` passaram. O `api.core` está pouco coberto (só o teste de tamanho dos campos).

## 7. Pontos de atenção encontrados

1. **Auditoria desligada.** O `LoggerMiddleware` (grava em `system_logs` toda escrita bem-sucedida) está definido em `common/middlewares/logger.go`, mas **nenhum serviço o registra**, e a tabela não é migrada. Hoje ninguém sabe quem alterou o quê.
2. **Recuperar senha** é só tela. Falta endpoint, e ele depende de um canal de envio (e-mail), que o sistema ainda não tem.
3. **Menu sem tela de administração.** Mudar o menu exige API ou SQL.
4. **Textos de permissões em bancos já semeados** (como a VM): as correções do catálogo (CFG-014) só valem em banco novo. Nos já semeados precisam de um `UPDATE` à parte (pendência do release 2026.3.4.0).
5. **Desvios pontuais do padrão Service → UoW:** o `AuthService` usa o `s.db` direto, e `CompanyService.GetCompanies` entrega `GetDB()` ao `ResolveVisibleCompanyIDs` (exceção documentada, só leitura).
6. **`DB_AUTO_MIGRATE=false`** (banco abaixo da versão exigida bloqueia o serviço) passou nos testes automáticos, mas nunca rodou num servidor real.
7. **`BootstrapAdmin` é transitório**: a ideia é que um aplicativo de provisionamento crie o primeiro admin na fase multi-tenant.
8. **Geografia sem seed** (países, estados, cidades): `Document.StateID` é `uint` simples com FK, e o valor zero viola a constraint. `person_type` ainda é texto livre, sem validação `F`/`J`.
9. **HTTPS e domínio** ainda não existem no ambiente de testes (acesso por IP, HTTP).

## 8. O que mais pode entrar nesse conjunto

O `BACKLOG.md` **não tem nenhum item aberto de CFG**: tudo o que existia está `[x]`. O que segue é derivado das lacunas acima, **não de um plano já decidido**, e cada item precisa da sua decisão. Marquei a origem de cada um.

| Candidato | Onde entraria | Origem |
|---|---|---|
| Ligar a **auditoria** (migrar `system_logs`, registrar o middleware, tela de consulta) | `api.auth` ou `common` | achado nº 1 (código já existe) |
| **Recuperar senha** de ponta a ponta (token, e-mail, endpoint, tela) | `api.auth` | achado nº 2 (a tela já existe) |
| **Tela de administração do menu** (só DEVELOP) | `api.auth` + frontend | achado nº 3 (API já existe) |
| **Geografia sob demanda** (CEP cria país/estado/cidade) e `Document.StateID` opcional | `api.core` | `CLAUDE.md`, "Sem seed de geografia"; entra junto com as telas de Cliente/Fornecedor |
| **Script de correção de textos de permissões** em bancos semeados | `api.auth` | pendência do release |
| **Provisionamento do primeiro admin** pelo app comercial | fora deste sistema | `plano_multitenant.md` |
| Política de senha, bloqueio por tentativas, sessões ativas, 2FA | `api.auth` | **sugestão minha**, nada nos planos |
| Parâmetros gerais do sistema (chave/valor por empresa) | `api.core` | **sugestão minha**; hoje só o `api.integrations` tem `integration_configs` |

O README do planejamento descreve o Módulo 0 como incluindo "plano de contas", enquanto o `BACKLOG.md` o coloca em `CTB-001` (Contabilidade, bloqueado por decisão). Vale acertar uma das duas fontes.

## 9. Recomendação de estrutura

Não criar um `api.config`. A divisão atual segue o que cada serviço guarda: `api.auth` cuida de **quem pode entrar e fazer o quê**, e `api.core` cuida de **dados da empresa e cadastros base**. Qualquer item da seção 8 cabe num dos dois. Se mais tarde quiser um serviço só de configuração, o candidato natural é tirar `companies` e `company_fiscal_configs` do `api.core` (hoje 12 endpoints), mas isso mexeria no gateway, no banco e na multiempresa já testada, sem ganho visível agora.
