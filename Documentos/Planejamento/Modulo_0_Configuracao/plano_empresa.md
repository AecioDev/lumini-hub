# Módulo Empresa - Cadastro Multi-Empresa e Fundação Fiscal/Contábil

> **Onde este arquivo mora**: coloquei dentro de `Modulo_0_Configuracao/` porque é fundação (todo o resto do sistema depende disso), não um módulo de negócio isolado como CRM/Compras/Vendas. Se preferir uma pasta própria (`Modulo_0B_Empresa/` ou renumerar), me avisa — não tem custo trocar agora que ainda é só planejamento.

Este módulo cria o cadastro de **Empresa** — a estrutura organizacional de um único cliente (tenant) da Lumini Hub, que pode ter uma empresa Matriz com outras empresas vinculadas embaixo dela (filial por CNPJ ou não, mas geridas pelo mesmo gestor geral). Junto dele, criamos a base de duas coisas que os Módulos 6 (Fiscal) e 7 (Contabilidade) já pressupõem sem ela existir: **Configuração Fiscal** por empresa (certificados, tributação, numeração de notas) e **Plano de Contas** por empresa.

**Importante**: isso não é multi-tenant. Multi-tenant (cliente A da Lumini vs. cliente B) é banco de dados físico separado, resolvido fora deste módulo. Aqui é a estrutura **dentro** do banco de um único cliente.

---

## 🛠️ Escopo Geral

1. **Cadastro de Empresa** — self-referencing (Matriz → Filiais/Empresas vinculadas), CNPJ, razão social, nome fantasia, dados básicos.
2. **Configuração Fiscal por Empresa** — certificado digital A1, tipo de tributação (Simples/Presumido/Real), dados do contador responsável, contadores/séries de numeração de notas fiscais.
3. **Plano de Contas por Empresa** — estrutura hierárquica de contas contábeis pra lançamentos em partida dobrada (base do Módulo 7 — CMV, DRE, Balancete).
4. **Escopo por empresa em todo o resto do sistema** — regra de visibilidade de dados (usuário só vê o que é da(s) empresa(s) dele) e a decisão de escopo obrigatória em toda entidade nova (já reforçada em `.agents/skills/lumini_hub_entity_creation/SKILL.md`, Step 0).

---

## 📐 Modelagem

### `Empresa`
| Campo | Tipo | Observação |
|---|---|---|
| `ID` | uint | PK |
| `ParentID` | *uint (nullable) | self-referencing, mesmo padrão de `MenuItem.ParentID` (`Backend/microservices/api.auth/internal/domain/menu_item.go`). **Definido 2026-07-21**: a Matriz é simplesmente a empresa com `ParentID == nil` — sem flag explícita `IsMatriz`. Validação a garantir: só pode existir uma empresa raiz por tenant (constraint/checagem de negócio, já que fisicamente nada impede duas linhas com `ParentID == nil` numa tabela normal). |
| `RazaoSocial` | string | |
| `NomeFantasia` | string | |
| `CNPJ` | string | único |
| `IsActive` | bool | |

### `EmpresaConfiguracaoFiscal` (1:1 com Empresa)
**Definido 2026-07-21**: essa tabela guarda o certificado digital A1 + dados cadastrais fiscais básicos (tributação, contador). Numeração/série de notas fica em tabela própria à parte (ver `EmpresaConfiguracaoEmissaoNota` abaixo), não aqui — são conceitos diferentes.

| Campo | Tipo | Observação |
|---|---|---|
| `EmpresaID` | uint | FK |
| `CertificadoArquivo` | bytea | **Definido 2026-07-21**: guardado como bytes na própria coluna (`bytea`) — decisão final, pra reduzir complexidade de leitura/gerenciamento de arquivo em disco. |
| `CertificadoSenha` | string (criptografado) | senha do certificado, **criptografada com chave** (não hash — precisa ser reversível pra uso no ACBrLib), nunca texto puro |
| `CertificadoVencimento` | date | data de validade do certificado, pra alertar o usuário antes de vencer |
| `TipoTributacao` | string/enum | Simples / Presumido / Real |
| `DadosContador` | — | nome, CPF/CRC, contato — campos exatos a definir conforme for evoluindo |

> **Por que bytea**: certificado A1 costuma ser um arquivo pequeno (poucos KB), então o overhead do bytea é irrelevante. Guardar na própria linha significa que o certificado sempre viaja junto do backup do banco (nosso modelo de tenant = banco separado já trata isso automaticamente), sem precisar gerenciar uma pasta/volume compartilhado à parte nem se preocupar com arquivo órfão se a empresa for excluída.

### `EmpresaConfiguracaoEmissaoNota` (1 Empresa : N séries/tipos de documento)
**Novo, a partir do que você descreveu em 2026-07-21** — separado da config fiscal genérica porque numeração de nota é um conceito específico que cresce (múltiplos tipos de documento, séries, CFOPs). Estrutura ainda em esboço, campos a refinar conforme o Módulo 6 (Fiscal) for sendo desenhado de verdade:

| Campo | Tipo | Observação |
|---|---|---|
| `EmpresaID` | uint | FK |
| `TipoDocumento` | enum | orçamento / pedido / NFe / NFSe / NFCe / MDFe |
| `TipoOperacao` | string/enum | a definir junto do Módulo 6 |
| `Serie` | string | |
| `Sequencia` | int | contador atual, incrementado a cada emissão (mesmo espírito do `GetSequencia`/`CADSEQ` do legado, ver `CLAUDE.md`) |
| `CFOP` | string | varia por entrada/saída e por estadual/interestadual — provavelmente mais de um CFOP configurado por combinação, estrutura exata a definir |

> **Config geral, não só fiscal**: como você apontou, `Empresa` vai precisar de uma "gaveta" de configurações que cresce com o sistema (não só fiscal) — outros módulos vão pedir suas próprias configs por empresa conforme forem sendo construídos. Não vou tentar prever isso tudo agora; cada módulo novo pede sua própria tabela de configuração quando chegar a vez dele (mesma lógica do Step 0 de escopo por empresa — perguntar antes de modelar, não assumir).

### `PlanoDeContas` (N:1 com Empresa, self-referencing)
**Definido 2026-07-21**: template global padrão (o mais comumente usado no Brasil), clonado pra cada empresa na criação — facilita a vida do cliente, que já começa operando sem montar plano de contas do zero. Depois de clonado, cada empresa pode customizar o próprio livremente.

| Campo | Tipo | Observação |
|---|---|---|
| `EmpresaID` | uint | FK — cada empresa tem sua própria cópia, clonada do template padrão na criação |
| `ParentID` | *uint (nullable) | self-referencing — hierarquia tipo "3.1.01.001 Receita de Vendas" sob "3.1 Receitas Operacionais" sob "3 Contas de Resultado" |
| `Codigo` | string | ex: `3.1.01.001` |
| `Nome` | string | |
| `Tipo` | enum | Ativo / Passivo / Receita / Despesa / Resultado (a validar terminologia contábil exata) |

---

## 🔐 Regra de Visibilidade (definida em 2026-07-21)

Dois casos, conforme o cadastro do usuário:

**Caso 1 — usuário SEM Empresa vinculada no cadastro.**
Automaticamente "master": visão geral de todas as Empresas, sem precisar de nenhuma permission extra. Inicia sempre pela empresa Matriz (a com `ParentID == nil`). Pode trocar de empresa livremente através de um seletor no header (antd `Select`), sem restrição.

Rastreio da empresa ativa (**definido 2026-07-21**, Opção B): fica guardado no backend, fora do JWT — mesma ideia de variável de sessão que o usuário já usa hoje no ScriptCase. Um endpoint dedicado (ex.: `PUT /me/active-empresa`) grava a escolha; as próximas requisições dessa sessão já sabem qual é a empresa ativa sem o frontend precisar reenviar. Toda leitura desse valor ainda reconfere no backend se o usuário realmente pode ver aquela empresa antes de aplicar como filtro (nunca confiar cegamente no valor guardado). Mecanismo exato de guarda (cookie próprio de sessão vs. tabela) fica pra quando `Empresa` estiver sendo implementada.

**Caso 2 — usuário COM Empresa vinculada no cadastro.**
- Se tiver a permission `empresa.hierarquia.view` (nome definido em 2026-07-21): vê a própria empresa **e** as empresas abaixo dela na hierarquia (filhas, netas, etc.).
- Se não tiver essa permission: vê **apenas** a própria empresa vinculada.

Essa regra vale pra toda entidade "hard-scoped" (bucket 1 do Step 0 da skill de criação de entidade — estoque, vendas, financeiro, etc.): a query sempre filtra pelo conjunto de `EmpresaID`s visíveis calculado a partir dessas duas regras, nunca por uma única empresa fixa.

---

## 🧩 Escopo por Empresa nas demais entidades (referência rápida)

Baldes definidos no Step 0 de `lumini_hub_entity_creation/SKILL.md` — perguntar isso é obrigatório antes de modelar qualquer entidade nova:

1. **Hard-scoped** (`EmpresaID` próprio): estoque, vendas, pedidos de compra, lançamentos financeiros, caixas — nunca podem misturar dado entre empresas.
2. **Global com relacionamento fraco**: `Product` é a referência — catálogo compartilhado entre as empresas do tenant, mas uma tabela `product_empresas` (N:N) controla quais produtos cada empresa pode ver/vender.
3. **Global sem relação**: catálogo de Permissões, tabelas de lookup genéricas.

---

## ❓ Decisões em Aberto

Todas as decisões estruturais de arquitetura foram fechadas em 2026-07-21: `IsMatriz` (é `ParentID == nil`), nome da permission (`empresa.hierarquia.view`), Plano de Contas (template global clonável por empresa), certificado (bytea na própria linha, senha criptografada, campo de vencimento), rastreio de empresa ativa (Opção B, fora do JWT, estilo variável de sessão) e **local do microsserviço: `api.core`** (mesmo serviço de Customers/Suppliers/Addresses/Contacts hoje — não justifica um serviço dedicado só pra isso, sem necessidade de escala própria). `api.auth` referencia `EmpresaID` por ID puro em `User`, sem GORM relation cross-service, mesmo padrão já usado com `Customer`.

Só ficam pendências de **conteúdo**, não de arquitetura — não bloqueiam começar a modelar `Empresa`/`EmpresaConfiguracaoFiscal`, mas bloqueiam `PlanoDeContas` funcionar de verdade:

- [ ] **O conteúdo real do Plano de Contas padrão.** Preciso da lista de contas em si (código + nome + tipo, hierarquia completa) pra ter o que clonar em cada empresa nova. Isso é dado contábil/de negócio, não uma decisão de engenharia — o usuário tem um plano de contas de referência (do contador, de outro sistema, etc.) pra eu usar como seed, ou monto uma sugestão padrão simplificada pra começar e ajustamos depois?
- [ ] Campos exatos de `EmpresaConfiguracaoEmissaoNota` (`TipoOperacao`, estrutura de `CFOP`) — fica pra quando o Módulo 6 (Fiscal) for desenhado de verdade, não bloqueia nada agora.

---

## Próximos Passos

1. Fechar as decisões em aberto acima (ou seguir com os defaults sugeridos e ajustar depois).
2. Criar `tasks_empresa.md` decompondo isso em tarefas (seguindo o protocolo de bloqueio do `lumini_hub_dev_flow`).
3. Seguir o checklist de 12 passos da skill `lumini_hub_entity_creation` pra `Empresa` primeiro (é a base de tudo o resto aqui), depois `EmpresaConfiguracaoFiscal`/`PlanoDeContas`.
