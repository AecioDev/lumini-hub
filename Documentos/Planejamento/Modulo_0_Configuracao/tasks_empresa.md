# Tasks - Módulo Empresa

Acompanhamento da implementação do cadastro de Empresa e fundação fiscal/contábil (ver `plano_empresa.md`). Protocolo de bloqueio conforme `lumini_hub_dev_flow`.

## Backend — `api.core`

- [x] [CONCLUÍDO POR: Claude] `Company` — domain model, DTOs, repository, validator, service, handler, rotas, AutoMigrate, proxy no gateway, permissões `companies.*` seedadas pro ADMIN, swagger regenerado (checklist de 12 passos da skill `lumini_hub_entity_creation`). Self-referencing com validação de raiz única (`ParentID == nil`) e detecção de ciclo. Sem endpoint `/filter` de propósito — volume esperado de empresas por tenant é baixo. **Renomeado 2026-07-27**: implementado inicialmente em português (`Empresa`/tabela `empresas`), destoava do resto do banco — renomeado pro padrão do projeto (`Company`/tabela `companies`, ver `plano_empresa.md`). **Testado end-to-end em 2026-07-27** via curl através do gateway (login ADMIN → POST cria Matriz → GET/:id → PUT atualiza → DELETE → GET volta vazio) — tudo OK. Tabela antiga `empresas` continua órfã no banco, ainda não dropada.
- [ ] `CompanyFiscalConfig` — domain/repository/service/handler/rotas (certificado A1 em `bytea`, senha criptografada, vencimento, tipo de tributação, dados do contador)
- [ ] `CompanyDocumentIssuanceConfig` — estrutura fina depende do Módulo 6 (Fiscal) ser desenhado; não iniciar antes disso
- [ ] `ChartOfAccounts` — domain/repository/service/handler/rotas, self-referencing, com endpoint/lógica de clonar o template padrão na criação de uma Empresa
- [ ] `CompanyVisualConfig` — domain/repository/service/handler/rotas (logo em `bytea` + mimetype, cor primária/secundária/accent), 1:1 com Company, endpoint de upload de logo separado do PUT de configuração (multipart)
- [ ] Conteúdo do template padrão do Plano de Contas (seed) — bloqueado até o usuário decidir/fornecer o conteúdo (ver `plano_empresa.md`, ele disse que por enquanto pode ser um modelo básico a alinhar depois com o contador)
- [ ] `User.CompanyID` — campo em `api.auth`, nullable, referência por ID puro (sem GORM relation cross-service)
- [ ] Regra de visibilidade hierárquica — permission `companies.hierarchy.view` + lógica de resolver o conjunto de `CompanyID`s visíveis por usuário (Caso 1 master / Caso 2 com hierarquia)
- [ ] Mecanismo de "empresa ativa" do usuário master — endpoint `PUT /me/active-company` (ou nome final) + guarda de estado fora do JWT, reconferindo visibilidade a cada leitura

## Frontend

- [x] [CONCLUÍDO POR: Claude] Tela de Cadastro/Edição de Empresa (lista + full-page + Zod, mesmo padrão de Usuários/Perfis) — `CompaniesTable`, `CompanyForm` (com Select de "Vinculada a" excluindo a própria empresa), `CompaniesListPage`/`CreateCompanyPage`/`EditCompanyPage`, rotas em `router.tsx`.
- [ ] Seletor de empresa ativa no header (visível só pro usuário master, Caso 1)
- [x] [CONCLUÍDO POR: Claude] Item de menu "Empresas" inserido em `menu_item_seeder.go`, dentro de Configurações, gate `companies.view` — precisa reiniciar `api.auth` pra propagar (seeder é idempotente por item).
- [ ] Aba/seção "Identidade Visual" na tela de Empresa — upload de logo + color pickers (primária/secundária/accent), preview ao vivo
- [ ] `AuthContext.tsx` passa a expor a config visual da empresa ativa junto de `user`/`menuItems`; `App.tsx` aplica via `ConfigProvider` sobre os tokens padrão de `src/theme/antd-theme.ts` (fallback pra paleta Lumini Hub se a empresa não tiver configurado nada)

## Fora do escopo por enquanto

- Export de Plano de Contas pra Excel (mencionado pelo usuário como uso futuro, pra mandar pro contador validar) — não é bloqueante da fundação, avaliar quando o módulo de Contabilidade for desenhado.
