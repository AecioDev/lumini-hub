# Módulo Empresa - Cadastro Multi-Empresa e Fundação Fiscal/Contábil

> **Onde este arquivo mora**: coloquei dentro de `Modulo_0_Configuracao/` porque é fundação (todo o resto do sistema depende disso), não um módulo de negócio isolado como CRM/Compras/Vendas. Se preferir uma pasta própria (`Modulo_0B_Empresa/` ou renumerar), me avisa — não tem custo trocar agora que ainda é só planejamento.

Este módulo cria o cadastro de **Empresa** — a estrutura organizacional de um único cliente (tenant) da Lumini Hub, que pode ter uma empresa Matriz com outras empresas vinculadas embaixo dela (filial por CNPJ ou não, mas geridas pelo mesmo gestor geral). Junto dele, criamos a base de duas coisas que os Módulos 6 (Fiscal) e 7 (Contabilidade) já pressupõem sem ela existir: **Configuração Fiscal** por empresa (certificados, tributação, numeração de notas) e **Plano de Contas** por empresa.

**Importante**: isso não é multi-tenant. Multi-tenant (cliente A da Lumini vs. cliente B) é banco de dados físico separado, resolvido fora deste módulo. Aqui é a estrutura **dentro** do banco de um único cliente.

---

## 🛠️ Escopo Geral

1. **Cadastro de Empresa** — self-referencing (Matriz → Filiais/Empresas vinculadas), CNPJ, razão social, nome fantasia, dados básicos.
2. **Configuração Fiscal por Empresa** — certificado digital A1, tipo de tributação (Simples/Presumido/Real), dados do contador responsável, contadores/séries de numeração de notas fiscais.
3. **Plano de Contas por Empresa** — estrutura hierárquica de contas contábeis pra lançamentos em partida dobrada (base do Módulo 7 — CMV, DRE, Balancete).
4. **Escopo por empresa em todo o resto do sistema** — regra de visibilidade de dados (usuário só vê o que é da(s) empresa(s) dele) e a decisão de escopo obrigatória em toda entidade nova (já reforçada em `.claude/skills/lumini_hub_entity_creation/SKILL.md`, Step 0).
5. **Configuração Visual por Empresa** — logo e paleta de cores próprias, aplicadas no frontend pra personalizar a experiência de cada cliente que usa a Lumini Hub (ver `CompanyVisualConfig` abaixo).

---

## 📐 Modelagem

### `Company` (tabela `companies`)
**Renomeado 2026-07-27**: implementado inicialmente como `Empresa`/`empresas` (português), destoando do resto do banco (todas as demais tabelas — `customers`, `suppliers`, `menu_items`, etc. — são em inglês). Renomeado pro padrão do projeto: struct `Company`, tabela `companies`, campos abaixo, rotas `/companies`, permissões `companies.*`. A tabela antiga `empresas` fica órfã no banco até ser dropada manualmente (o `AutoMigrate` só cria/altera, não renomeia).

| Campo | Tipo | Observação |
|---|---|---|
| `ID` | uint | PK |
| `ParentID` | *uint (nullable) | self-referencing, mesmo padrão de `MenuItem.ParentID` (`Backend/microservices/api.auth/internal/domain/menu_item.go`). **Definido 2026-07-21, revisto 2026-09-05**: "Matriz" é simplesmente qualquer empresa com `ParentID == nil` — sem flag explícita `IsMatriz`. **Não é mais única por tenant**: o usuário pode ter vários grupos empresariais de fato independentes (sem holding entre si), cada um com sua própria raiz — removida a validação que bloqueava criar uma segunda empresa sem `parent_id` (`CompanyValidator`/`ExistsRoot`/`ExistsRootExcept`, existiam só até então). Continua validado: ciclo na cadeia de ancestrais (uma empresa não pode acabar sendo pai dela mesma via um vínculo indireto). |
| `LegalName` (`legal_name`) | string | Razão Social |
| `TradeName` (`trade_name`) | string | Nome Fantasia |
| `CNPJ` (`cnpj`) | *string (nullable) | **Renomeado de `TaxID`/`tax_id` em 2026-09-05** — CPF/CNPJ/CEP/RG mantêm o nome brasileiro em vez de uma tradução literal, ver `CLAUDE.md` § "Brazilian document/address fields". Único (índice parcial, ignora linhas soft-deleted e nulas). **Obrigatório só pra Matriz** — empresa vinculada pode ficar sem CNPJ quando a parte fiscal fica a cargo da Matriz (caso real do usuário, decidido 2026-09-05: cadastro de uma "vinculada" só pra personalizar relatórios/logo de uma marca/setor). |
| `IsActive` | bool | |

### `CompanyFiscalConfig` (1:1 com Empresa)
**Definido 2026-07-21**: essa tabela guarda o certificado digital A1 + dados cadastrais fiscais básicos (tributação, contador). Numeração/série de notas fica em tabela própria à parte (ver `CompanyDocumentIssuanceConfig` abaixo), não aqui — são conceitos diferentes.

| Campo | Tipo | Observação |
|---|---|---|
| `CompanyID` | uint | FK |
| `CertificateFile` | bytea | **Definido 2026-07-21**: guardado como bytes na própria coluna (`bytea`) — decisão final, pra reduzir complexidade de leitura/gerenciamento de arquivo em disco. |
| `CertificatePassword` | string (criptografado) | senha do certificado, **criptografada com chave** (não hash — precisa ser reversível pra uso no ACBrLib), nunca texto puro |
| `CertificateExpiry` | date | data de validade do certificado, pra alertar o usuário antes de vencer |
| `TaxRegime` | string/enum | Simples / Presumido / Real |
| `AccountantInfo` | — | nome, CPF/CRC, contato — campos exatos a definir conforme for evoluindo |

> **Por que bytea**: certificado A1 costuma ser um arquivo pequeno (poucos KB), então o overhead do bytea é irrelevante. Guardar na própria linha significa que o certificado sempre viaja junto do backup do banco (nosso modelo de tenant = banco separado já trata isso automaticamente), sem precisar gerenciar uma pasta/volume compartilhado à parte nem se preocupar com arquivo órfão se a empresa for excluída.

### `CompanyDocumentIssuanceConfig` (1 Empresa : N séries/tipos de documento)
**Novo, a partir do que você descreveu em 2026-07-21** — separado da config fiscal genérica porque numeração de nota é um conceito específico que cresce (múltiplos tipos de documento, séries, CFOPs). Estrutura ainda em esboço, campos a refinar conforme o Módulo 6 (Fiscal) for sendo desenhado de verdade:

| Campo | Tipo | Observação |
|---|---|---|
| `CompanyID` | uint | FK |
| `DocumentType` | enum | orçamento / pedido / NFe / NFSe / NFCe / MDFe |
| `OperationType` | string/enum | a definir junto do Módulo 6 |
| `Series` | string | |
| `Sequence` | int | contador atual, incrementado a cada emissão (mesmo espírito do `GetSequencia`/`CADSEQ` do legado, ver `CLAUDE.md`) |
| `CFOP` | string | varia por entrada/saída e por estadual/interestadual — provavelmente mais de um CFOP configurado por combinação, estrutura exata a definir |

> **Config geral, não só fiscal**: como você apontou, `Company` vai precisar de uma "gaveta" de configurações que cresce com o sistema (não só fiscal) — outros módulos vão pedir suas próprias configs por empresa conforme forem sendo construídos. Não vou tentar prever isso tudo agora; cada módulo novo pede sua própria tabela de configuração quando chegar a vez dele (mesma lógica do Step 0 de escopo por empresa — perguntar antes de modelar, não assumir).

### `CompanyVisualConfig` (1:1 com Empresa)
**Definido 2026-07-27**: personalização visual por cliente (a Lumini Hub é usada por várias empresas diferentes, cada uma podendo querer sua própria marca aparecendo no sistema). Tabela própria, não uma coluna solta em `Company` nem misturada nas outras configs (fiscal/emissão de nota) — é um conceito totalmente diferente (visual, não fiscal/contábil) e mais uma "gaveta" que só cresce (pode ganhar campos como tema dark/light forçado, fonte, etc. no futuro).

| Campo | Tipo | Observação |
|---|---|---|
| `CompanyID` | uint | FK |
| `LogoFile` | bytea | mesmo raciocínio do `CertificateFile` em `CompanyFiscalConfig` — arquivo pequeno, guardado na própria linha, viaja junto do backup |
| `LogoMimeType` | string | `image/png`, `image/svg+xml`, `image/jpeg` — necessário pra servir o arquivo corretamente, já que é `bytea` puro |
| `PrimaryColor` | string | hex, ex: `#2563EB` — vira `colorPrimary` no `ConfigProvider` do antd |
| `SecondaryColor` | string (nullable) | hex — uso pontual (ex: `colorInfo`/detalhes), opcional |
| `AccentColor` | string (nullable) | hex — terceira cor pra gradientes/destaques, opcional, mesmo espírito do roxo na paleta da própria Lumini Hub |

> **Aplicação no frontend**: mesmo mecanismo do menu dinâmico (`AuthContext.tsx`) — a config visual da empresa ativa do usuário viaja no bootstrap da sessão (`GET /auth/me`/login/refresh) e substitui os tokens de `src/theme/antd-theme.ts` (`BRAND`) em runtime via `ConfigProvider`. **Fallback**: campos vazios/tabela sem registro para aquela empresa usam a paleta padrão da própria Lumini Hub (não força o cliente a configurar nada pra já ter um visual coerente).

### `CompanyColorPalette` (N:1 com Empresa) — paletas personalizadas

**Definido 2026-09-06**, a partir de feedback do usuário testando a tela de Identidade Visual: além das 6 paletas fixas no código (`PRESET_PALETTES`, frontend), o usuário quer poder salvar combinações próprias de cor pra reaproveitar depois, com um botão "Salvar como paleta personalizada" logo após os 3 color pickers.

- **Escopo por Empresa (Step 0)**: hard-scoped, `CompanyID` próprio — decidido explicitamente pelo usuário ("por empresa, é uma configuração de ADMIN não uma personalização de usuário"). Uma paleta salva editando a Empresa A **não** aparece editando a Empresa B, mesmo dentro do mesmo tenant — cada empresa tem sua própria lista.
- **É lista de verdade, não config 1:1** (diferente de `CompanyVisualConfig`/`CompanyFiscalConfig`) — uma empresa pode ter zero, uma ou várias paletas salvas. Por isso **tem Delete** (Step 0.5 não se aplica aqui).
- Sem Update por enquanto — o fluxo é criar uma nova ou apagar uma antiga, não editar uma existente in-place (MVP; revisitar se virar necessidade real).

| Campo | Tipo | Observação |
|---|---|---|
| `CompanyID` | uint | FK, hard-scoped |
| `Name` | string | nome dado pelo usuário pra identificar a paleta na lista (ex: "Campanha de Verão") |
| `PrimaryColor` | string | hex, obrigatória — mesmo formato de `CompanyVisualConfig.PrimaryColor` |
| `SecondaryColor` | string (nullable) | hex, opcional — espelha o que estava preenchido no formulário no momento de salvar |
| `AccentColor` | string (nullable) | hex, opcional |

> **Permissões**: reaproveita `companies.visual_config.create` (criar paleta) e `companies.visual_config.edit` (apagar paleta) em vez de um catálogo `companies.color_palettes.*` dedicado — a feature é um apêndice da mesma tela/permissão de Identidade Visual, não uma tela própria; criar uma permission nova só pra isso infla o catálogo sem necessidade real de granularidade separada.
> **UI**: lista de paletas salvas exibida ao lado das 6 fixas (mesmo componente de swatch clicável), com uma ação de remover por paleta. Formulário de "nome + salvar" é **inline** (`Input` + `Button` na própria tela), não modal — CLAUDE.md reserva modais pra confirmação de ação destrutiva/usuário, nunca pra entrada de dados.

### `ChartOfAccounts` (N:1 com Empresa, self-referencing)
**Definido 2026-07-21**: template global padrão (o mais comumente usado no Brasil), clonado pra cada empresa na criação — facilita a vida do cliente, que já começa operando sem montar plano de contas do zero. Depois de clonado, cada empresa pode customizar o próprio livremente.

| Campo | Tipo | Observação |
|---|---|---|
| `CompanyID` | uint | FK — cada empresa tem sua própria cópia, clonada do template padrão na criação |
| `ParentID` | *uint (nullable) | self-referencing — hierarquia tipo "3.1.01.001 Receita de Vendas" sob "3.1 Receitas Operacionais" sob "3 Contas de Resultado" |
| `Code` | string | ex: `3.1.01.001` |
| `Name` | string | |
| `Type` | enum | Ativo / Passivo / Receita / Despesa / Resultado (a validar terminologia contábil exata) |

---

## 🚫 Configs 1:1 não têm Delete (definido em 2026-09-05)

Descoberto testando `CompanyFiscalConfig`: um índice único (`company_id`) mais soft delete do GORM permite que uma linha excluída continue "ocupando" o valor único, bloqueando recriar a config da mesma empresa depois. A correção técnica (índice parcial `WHERE deleted_at IS NULL`) resolve o sintoma, mas discutindo com o usuário chegamos numa raiz melhor: **uma configuração 1:1 intrínseca à empresa não tem "excluir" com significado de negócio** — enquanto a empresa existe, ela sempre tem (ou devia ter) aquela configuração associada, nem que seja vazia. Diferente de excluir um registro de uma lista de verdade (`Company`, `ChartOfAccounts`, `CompanyDocumentIssuanceConfig`), onde a linha deixar de existir faz sentido.

**Daqui pra frente**: toda config 1:1 da empresa (`CompanyFiscalConfig`, `CompanyVisualConfig`, e qualquer outra do mesmo formato que surgir) só expõe `view`/`create`/`edit` — sem `delete`, sem botão de excluir na tela, sem permission `.delete` no catálogo. Se um dia for preciso "limpar" um campo específico (ex.: remover o certificado sem apagar o resto da config), isso vira uma ação dedicada (ex.: um endpoint específico) ou um `PUT` que aceita valor vazio — nunca apagar o registro inteiro. Esse princípio também foi registrado na skill `lumini_hub_entity_creation` (Step 0), pra ser considerado antes de modelar qualquer config nova.

---

## 🔐 Regra de Visibilidade (definida em 2026-07-21, ajustada em 2026-09-05)

Dois casos, conforme o cadastro do usuário:

**Caso 1 — usuário SEM Empresa vinculada no cadastro.**
Automaticamente "master": visão geral de todas as Empresas, sem precisar de nenhuma permission extra. Pode trocar de empresa livremente através de um seletor no header (antd `Select`), sem restrição.

**Mecanismo de empresa ativa — implementado em 2026-09-05** (Opção B confirmada em 2026-07-21: guardado no backend, fora do JWT): coluna `users.active_company_id` (nullable, FK pra `companies.id`), gravada via `PUT /auth/active-company` (`common/utils.SetActiveCompany` valida que o usuário realmente enxerga aquela empresa antes de gravar — nunca aceita o ID cego vindo do cliente) e resolvida/revalidada a cada login/refresh/`GET /auth/me` (`common/utils.ResolveActiveCompany` — nunca confia cegamente no valor guardado; se ele não for mais válido, some do payload e o frontend volta a pedir escolha). Quem só enxerga UMA empresa (sem `companies.hierarchy.view`, com Company vinculada) nunca precisa escolher — a ativa é sempre a própria. Quem enxerga mais de uma (master, ou com `companies.hierarchy.view`) e ainda não tem uma empresa ativa válida recebe `requires_company_selection: true` no payload de login/me — o frontend (`ProtectedRoute.tsx` + `CompanySelectionGate.tsx`) bloqueia o acesso ao resto do sistema com uma tela de seleção obrigatória até isso ser resolvido, escolha então persiste entre logins (só pede de novo se a empresa salva deixar de ser válida). O seletor de troca no header (`AuthContext`/`ConfigProvider`) ainda não foi implementado (ver `tasks_empresa.md`, EPIC CFG-6).

**Caso 2 — usuário COM Empresa vinculada no cadastro.**
- Se tiver a permission `companies.hierarchy.view` (nome definido em 2026-07-21, renomeado 2026-07-27 pro padrão `<módulo>.<ação>` em inglês): vê a própria empresa **e** as empresas abaixo dela na hierarquia (filhas, netas, etc.).
- Se não tiver essa permission: vê **apenas** a própria empresa vinculada.

Essa regra vale pra toda entidade "hard-scoped" (bucket 1 do Step 0 da skill de criação de entidade — estoque, vendas, financeiro, etc.): a query sempre filtra pelo conjunto de `CompanyID`s visíveis calculado a partir dessas duas regras, nunca por uma única empresa fixa.

---

## 🧩 Escopo por Empresa nas demais entidades (referência rápida)

Baldes definidos no Step 0 de `lumini_hub_entity_creation/SKILL.md` — perguntar isso é obrigatório antes de modelar qualquer entidade nova:

1. **Hard-scoped** (`CompanyID` próprio): estoque, vendas, pedidos de compra, lançamentos financeiros, caixas — nunca podem misturar dado entre empresas.
2. **Global com relacionamento fraco**: `Product` é a referência — catálogo compartilhado entre as empresas do tenant, mas uma tabela `product_companies` (N:N) controla quais produtos cada empresa pode ver/vender.
3. **Global sem relação**: catálogo de Permissões, tabelas de lookup genéricas.

---

## ❓ Decisões em Aberto

Todas as decisões estruturais de arquitetura foram fechadas em 2026-07-21: `IsMatriz` (é `ParentID == nil`), nome da permission (`companies.hierarchy.view`), Plano de Contas (template global clonável por empresa), certificado (bytea na própria linha, senha criptografada, campo de vencimento), rastreio de empresa ativa (Opção B, fora do JWT, estilo variável de sessão) e **local do microsserviço: `api.core`** (mesmo serviço de Customers/Suppliers/Addresses/Contacts hoje — não justifica um serviço dedicado só pra isso, sem necessidade de escala própria). `api.auth` referencia `CompanyID` por ID puro em `User`, sem GORM relation cross-service, mesmo padrão já usado com `Customer`.

Só ficam pendências de **conteúdo**, não de arquitetura — não bloqueiam começar a modelar `Company`/`CompanyFiscalConfig`, mas bloqueiam `ChartOfAccounts` funcionar de verdade:

- [ ] **O conteúdo real do Plano de Contas padrão.** Preciso da lista de contas em si (código + nome + tipo, hierarquia completa) pra ter o que clonar em cada empresa nova. Isso é dado contábil/de negócio, não uma decisão de engenharia — o usuário tem um plano de contas de referência (do contador, de outro sistema, etc.) pra eu usar como seed, ou monto uma sugestão padrão simplificada pra começar e ajustamos depois?
- [ ] Campos exatos de `CompanyDocumentIssuanceConfig` (`OperationType`, estrutura de `CFOP`) — fica pra quando o Módulo 6 (Fiscal) for desenhado de verdade, não bloqueia nada agora.

---

## Próximos Passos

1. Fechar as decisões em aberto acima (ou seguir com os defaults sugeridos e ajustar depois).
2. Criar `tasks_empresa.md` decompondo isso em tarefas (seguindo o protocolo de bloqueio do `lumini_hub_dev_flow`).
3. Seguir o checklist de 12 passos da skill `lumini_hub_entity_creation` pra `Company` primeiro (é a base de tudo o resto aqui), depois `CompanyFiscalConfig`/`ChartOfAccounts`.
