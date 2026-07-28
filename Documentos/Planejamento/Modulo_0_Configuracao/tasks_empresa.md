# Tasks - Módulo Empresa

Acompanhamento da implementação do cadastro de Empresa e fundação fiscal/contábil (ver `plano_empresa.md`). Protocolo de bloqueio conforme `lumini_hub_dev_flow`.

## Backend — `api.core`

- [x] [CONCLUÍDO POR: Claude] `Empresa` — domain model, DTOs, repository, validator, service, handler, rotas, AutoMigrate, proxy no gateway, permissões `empresas.*` seedadas pro ADMIN, swagger regenerado (checklist de 12 passos da skill `lumini_hub_entity_creation`). Self-referencing com validação de raiz única (`ParentID == nil`) e detecção de ciclo. Sem endpoint `/filter` de propósito — volume esperado de empresas por tenant é baixo. Falta testar end-to-end (precisa reiniciar o backend pro `AutoMigrate` criar a tabela `empresas`).
- [ ] `EmpresaConfiguracaoFiscal` — domain/repository/service/handler/rotas (certificado A1 em `bytea`, senha criptografada, vencimento, tipo de tributação, dados do contador)
- [ ] `EmpresaConfiguracaoEmissaoNota` — estrutura fina depende do Módulo 6 (Fiscal) ser desenhado; não iniciar antes disso
- [ ] `PlanoDeContas` — domain/repository/service/handler/rotas, self-referencing, com endpoint/lógica de clonar o template padrão na criação de uma Empresa
- [ ] Conteúdo do template padrão do Plano de Contas (seed) — bloqueado até o usuário decidir/fornecer o conteúdo (ver `plano_empresa.md`, ele disse que por enquanto pode ser um modelo básico a alinhar depois com o contador)
- [ ] `User.EmpresaID` — campo em `api.auth`, nullable, referência por ID puro (sem GORM relation cross-service)
- [ ] Regra de visibilidade hierárquica — permission `empresa.hierarquia.view` + lógica de resolver o conjunto de `EmpresaID`s visíveis por usuário (Caso 1 master / Caso 2 com hierarquia)
- [ ] Mecanismo de "empresa ativa" do usuário master — endpoint `PUT /me/active-empresa` (ou nome final) + guarda de estado fora do JWT, reconferindo visibilidade a cada leitura

## Frontend

- [ ] Tela de Cadastro/Edição de Empresa (seguindo o padrão de páginas full-page + Zod, mesmo de Usuários/Perfis)
- [ ] Seletor de empresa ativa no header (visível só pro usuário master, Caso 1)
- [ ] Item de menu "Empresas" (inserir em `menu_item_seeder.go`, dentro de Configurações, gate `empresas.view`)

## Fora do escopo por enquanto

- Export de Plano de Contas pra Excel (mencionado pelo usuário como uso futuro, pra mandar pro contador validar) — não é bloqueante da fundação, avaliar quando o módulo de Contabilidade for desenhado.
