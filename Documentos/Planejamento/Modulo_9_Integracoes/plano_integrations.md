# Módulo 9 - Integrações (Loja Integrada ↔ SQL Server Legado)

Este módulo cobre o `api.integrations` (porta 4007) — middleware de sincronismo bidirecional entre a **Loja Integrada** (e-commerce) e o **SQL Server** legado do cliente (`FOCCO_ERP`), que o Lumini Hub está substituindo aos poucos.

> **Separado do CRM em 2026-09-13**: até então vivia junto do `api.crm` num único `Modulo_1_CRM_Integracoes`, por prioridade da época — mas Integrações é um microsserviço à parte, já com Fase 1 implementada e rodando, sem relação de código com o CRM. Ver `Modulo_1_CRM/plano_crm.md` pro CRM. O detalhe técnico completo (schema legado, endpoints, fluxo de estoque Oficial → Reserva → cliente) continua documentado em `CLAUDE.md` § "Legacy SQL Server Integration" — este arquivo não duplica, só referencia e complementa com o que falta.

## 🛠️ Especificações Gerais

1. **api.integrations (Porta 4007)**: sincronismo bidirecional de pedidos e estoque entre Loja Integrada e SQL Server legado. **Fase 1** (esqueleto, conexões Postgres/legado, `ConfigService`, handlers/rotas, proxy no gateway, tela Configurações → Integrações) entregue — ver `tasks_integrations.md`. **Fase 2** (domain `NOTAS`/`NOTAS1`/`NOTAS3`/`CADPRO`, cliente HTTP da Loja Integrada, scheduler de polling do `LogAltera`, processamento de `WebhookEvent` em pedido de venda) **escopada em 2026-09-27** (ver seção "🔌 Fase 2 — Integração real com a Loja Integrada" abaixo) — ainda não implementada, mas o bloqueio de chaves de API/App foi resolvido (caminho Personal Token, sem fila de aprovação). Seguem pendentes, e continuam sendo decisão do cliente/diretoria, não deste plano: a política de prioridade de estoque loja física vs. virtual, e a regra de emissão de nota/despacho pro pedido pago.
2. **Migração de dados do legado**: carga do histórico de clientes e transações do SQL Server legado pro Postgres novo — reaproveita a camada de leitura do legado já construída aqui na Fase 1, mas ainda sem escopo formal. Cruza com o CRM (`Modulo_1_CRM/plano_crm.md`), já que é o CRM quem consome esses dados migrados.

---

## 🔌 Fase 2 — Integração real com a Loja Integrada (escopado 2026-09-27)

Escopo levantado em sessão de brainstorming com o usuário, com pesquisa direta na documentação pública da Loja Integrada (Central de Ajuda oficial + um doc de terceiro que lista endpoints — fontes citadas em cada seção). O usuário opera o ERP legado do cliente com um login de **colaborador/administrador** na Loja Integrada, não o login do dono/proprietário original da loja.

### Autenticação — Personal Token (não App Key de parceiro)

A Loja Integrada tem dois modelos de autenticação, não intercambiáveis:

1. **Personal Token** — autosserviço, gerado pelo próprio **dono/proprietário** da loja em `Configurações > Chave para API`, autentica sozinho (`Authorization: Basic <token>`), sem depender de aprovação de ninguém. Exige plano pago ativo. Até 5 tokens ativos por loja; **cada token expira a cada 3 meses (90 dias)** e precisa ser renovado manualmente no painel.
   Fonte: [Como gerar chaves de API e o personal token](https://ajuda.lojaintegrada.com.br/pt-BR/articles/931152-como-gerar-chaves-de-api-e-o-personal-token-chave-de-aplicacao-da-minha-loja)
2. **Chave de API + Chave de Aplicação (App Key)** — modelo para parceiros/integradores multi-loja. A Chave de API da própria loja **sozinha não autentica nada**, só funciona combinada com uma App Key emitida pela equipe da Loja Integrada mediante formulário, com aprovação em **5 a 10 dias úteis** (exige IPs de origem, finalidade da integração etc.).
   Fontes: [Como obter a Chave de Aplicação](https://ajuda.lojaintegrada.com.br/pt-BR/articles/5360466-como-obter-a-chave-de-aplicacao-para-integrar-com-a-loja-integrada), [Como utilizar a API](https://ajuda.lojaintegrada.com.br/pt-BR/articles/12500189-como-utilizar-a-api-da-loja-integrada-para-automacoes-e-integracoes)

**Decisão**: usar o modelo **Personal Token**. O botão de gerar Personal Token aparece desabilitado pro login de colaborador que o usuário usa hoje (a documentação confirma que só o dono/proprietário tem acesso a esse botão, independente de plano) — mas o dono da loja confirmou que topa gerar/renovar o token sempre que solicitado. Isso evita a fila de aprovação de 5-10 dias do modelo de App Key, que não seria mais rápido nem necessário pra uma integração de loja única.

**Consequência no modelo de dados**: `IntegrationConfig` (`Backend/microservices/api.integrations/internal/domain/integration_config.go`) — que hoje é um armazenamento chave/valor com um contrato de API fixo (`ApiIntegrationSettings`/`UpdateIntegrationSettingsRequest`) desenhado assumindo o modelo de App Key — muda assim:

| Campo atual | Campo novo | Observação |
|---|---|---|
| `LiApiKey` | `LiPersonalToken` | renomeado — reflete o que de fato é armazenado e usado (`Authorization: Basic <token>`) |
| `LiAppKey` | *(removido)* | sem uso no caminho Personal Token; reintroduzir é barato se um dia for necessário (é chave/valor por baixo, não migração de schema pesada) |
| — | `LiPersonalTokenIssuedAt` (novo, timestamp) | data em que o token foi cadastrado, usada só pra calcular a expiração (emissão + 90 dias) |
| `LiWebhookSecret` | *(mantido por ora, ver seção "Segredo do webhook" abaixo)* | pode migrar de header pra query string dependendo do que a tela de cadastro de webhook da LI permitir na prática |

**Regra de negócio — aviso de expiração**: a tela Configurações → Integrações deve exibir um aviso quando o Personal Token estiver a **15 dias ou menos** de completar 90 dias desde `LiPersonalTokenIssuedAt`, pra dar tempo de pedir a renovação ao dono da loja antes que a sincronização pare silenciosamente.

### Webhooks assinados

A Loja Integrada só tem duas categorias de webhook, no nível da plataforma: **Produto** (Criado, Editado) e **Pedido** (Criado, Editado — "Editado" cobre qualquer mudança de situação/status: pago, cancelado, enviado, entregue, devolvido, aguardando pagamento; não existem webhooks individuais por status, é o mesmo evento "Pedido Editado" com o status mudando dentro do payload).
Fonte: [Como configurar webhook](https://ajuda.lojaintegrada.com.br/pt-BR/articles/9655071-como-configurar-webhook)

**Decisão**: assinar só a categoria **Pedido** (Criado + Editado). A categoria Produto não é assinada — o fluxo de produto/preço/estoque vai na direção contrária (SQL Server → LI, via polling do `LogAltera`, já descrito no `CLAUDE.md`), então o Lumini Hub não precisa ficar sabendo quando um produto muda na própria LI.

### Regras de negócio — reação por status do pedido

Processamento do `WebhookEvent` (`Backend/microservices/api.integrations/internal/domain/webhook_event.go`) recebido via `Pedido Editado`/`Pedido Criado`, decidido status a status com o usuário:

| Status | Ação do Lumini Hub |
|---|---|
| **Criado** | Cria a nota de venda no SQL Server legado (`NOTAS`/`NOTAS1`/`NOTAS3`, domain a construir na Fase 2) e reserva estoque: Oficial → Reserva (fluxo já descrito no `CLAUDE.md`). |
| **Pago** | Marca o pedido como pré-venda **pendente de faturamento no Caixa** — reaproveitando literalmente o mesmo mecanismo/fila que uma venda de loja física deixada sem faturar na hora (não é um status separado, exclusivo de pedidos LI). **Não** dispara faturamento automático, **não** move estoque além do que já foi reservado na criação. Confirmado pelo usuário: *"acredito que a princípio vamos deixar no caixa"*, tratando o pedido web como qualquer outra pré-venda pendente. |
| **Cancelado** | Devolve automaticamente a quantidade reservada pro estoque Oficial. |
| **Enviado / Entregue** | Só log/registro do evento — nenhuma ação automática no legado (a movimentação real de estoque pro cliente já acontece no momento do faturamento manual, conforme fluxo do `CLAUDE.md`). |
| **Devolvido** | Só log/registro — devolução física/fiscal tratada manualmente pelos operadores da loja, sem automação de estoque (decisão do usuário: *"Devoluções provavelmente vão ser tratados manualmente pelos operadores na loja"*). |

**Pendência explícita (decisão do cliente/diretoria, não deste plano)**: a regra de **emissão de nota fiscal e despacho** a partir dessa pré-venda pendente de faturamento ainda não foi definida — o usuário precisa verificar com o cliente como ele quer proceder. Isso se soma à pendência já conhecida da política de prioridade de estoque loja física vs. virtual quando os dois canais vendem a última unidade ao mesmo tempo (ver "Pendências em aberto" no fim desta seção).

### Configurações necessárias pra gerar a nota no legado

Levantado numa rodada de acompanhamento pedida pelo próprio usuário ao revisar o escopo: criar de fato o pedido/nota (`NOTAS`/`NOTAS1`/`NOTAS3`) a partir de um pedido da LI (status **Criado**, ver tabela acima) exige valores de cabeçalho que não vêm do próprio pedido. Resolvido nesta sessão:

**Um vendedor padrão resolve quase tudo**: o legado exige vendedor obrigatório em toda nota de venda. **Decisão**: configurar um único vendedor genérico (novo campo `CodVendedorPadrao` em `IntegrationConfig`) pra representar pedidos vindos da loja virtual — o cadastro desse vendedor no legado já carrega, por padrão, **centro de custo**, **tipo de nota** (`codtipnot`) e **tabela de preços**, então nenhum desses três vira campo de configuração isolado na tela de Integrações. Se algum outro valor de cabeçalho precisar de mapeamento próprio no futuro, revisita-se então — palavras do usuário: *"Se precisar de outro depois a gente mapeia, mas a princípio são esses."*

**`codlocarm_oficial` fica redundante, `codlocarm_reserva` continua necessário à parte**: o "local de armazenamento" do cadastro do vendedor genérico é o mesmo conceito que `codlocarm_oficial` (já configurado desde a Fase 1) — confirmado pelo usuário (*"Sim é a mesma coisa"*). Por isso `codlocarm_oficial` deixa de ser um campo de configuração isolado, resolvido via o vendedor padrão. Já `codlocarm_reserva` (o local pra onde o estoque vai quando um pedido web reserva a quantidade, fluxo Oficial → Reserva do `CLAUDE.md`) **continua existindo separadamente** — não faz parte do cadastro do vendedor, é específico do fluxo de reserva desta integração.

**Campos de `IntegrationConfig` pra gerar a nota — versão final** (complementa a tabela de autenticação, não repete):

| Campo | Situação | Observação |
|---|---|---|
| `CodVendedorPadrao` | **novo** | código do vendedor genérico do legado; resolve tipo de nota, centro de custo e tabela de preços automaticamente via o cadastro dele |
| `CodLocArmReserva` | já existe (Fase 1) | local pra onde o estoque vai na reserva Oficial → Reserva |
| `CodEmp` | já existe (Fase 1) | sem mudança |
| ~~`CodTipNot`~~ | **removido** | deixa de ser campo de configuração isolado — resolvido pelo cadastro do vendedor padrão acima. Não é mais pendência de decisão do cliente/diretoria (era um erro de classificação de uma versão anterior deste plano): é só um valor de configuração, preenchido automaticamente por já vir do cadastro do vendedor |
| ~~`CodLocArmOficial`~~ | **removido** | redundante com o local de armazenamento do vendedor padrão — confirmado que é a mesma coisa |

**De-para de forma de pagamento**: a Loja Integrada representa a forma de pagamento de um pedido com um código específico por loja (`pagamentos[].forma_pagamento.codigo`, ex.: `"pagseguro"`, `"entrega"` — confirmado pela pesquisa no blueprint público, ainda que descontinuado, da API da LI, fonte: [Loja Integrada API - Versão descontinuada (blueprint)](https://jsapi.apiary.io/apis/lojaintegrada.apib)). **Decisão**: o de-para entre esse código da LI e a forma de pagamento correspondente do legado vive **dentro da própria tabela de formas de pagamento do legado** (um campo lá guarda o vínculo) — não numa tabela nova no Postgres do Lumini Hub (chegou a ser cogitada uma `PaymentMethodMapping` no mesmo espírito do `ProductMapping` já existente, mas o usuário preferiu manter o vínculo do lado do legado). O Lumini Hub só **lê** esse campo ao montar o `NOTAS3`; **é responsabilidade do usuário/cliente garantir que esse campo já existe no legado antes da implementação começar, criando-o se necessário** — palavras do usuário: *"considere que existe um campo pra vincular a forma de pagamento da LI no legado, se não existir eu crio e te passo depois."*

**Parcelamento**: não foi encontrada confirmação de que a Loja Integrada exponha número de parcelas no payload do pedido (pode ficar só dentro do gateway de pagamento, ex.: PagSeguro/Mercado Pago). Em vez de tentar replicar o parcelamento real no `NOTAS3`, a decisão foi criar uma forma de pagamento específica no legado (ex.: "Parcelado LI") que gera um **único título a receber** por pedido pago parcelado — a equipe financeira dá baixa manualmente conforme os recebimentos reais chegam via o gateway de pagamento da própria LI.

### Segredo do webhook (`X-Webhook-Secret`) — a confirmar na prática

O endpoint já existente `POST /integrations/webhooks/loja-integrada` (`Backend/microservices/api.integrations/internal/handlers/webhooks.go`) valida um header `X-Webhook-Secret` antes de aceitar o payload. Não há confirmação pública de que a Loja Integrada permite configurar um header customizado ao cadastrar a URL de callback de um webhook — a documentação de registro de webhook está descrita dentro do fluxo de parceiro/App Key, caminho que este plano decidiu não seguir.

**Decisão**: o usuário vai tentar cadastrar o webhook na prática, com o acesso de colaborador que já tem, e verificar o que a tela realmente oferece. Se não permitir header customizado, o segredo passa a viajar como query string na própria URL (`.../webhooks/loja-integrada?secret=XXXX`), e o handler precisa ser ajustado pra ler de lá em vez do header — **ajuste técnico a confirmar durante a implementação**, não uma decisão de negócio pendente.

### Ambiente de testes

Não há confirmação pública de um ambiente de sandbox/homologação oficial da Loja Integrada — só relatos na comunidade pedindo um, sem resposta de que exista.
Fonte: [Sandbox para desenvolvimento de integração via API](https://comunidade.lojaintegrada.com.br/t/sandbox-para-desenvolvimento-de-integracao-via-api/52714)

**Mitigação**: não é um problema neste caso — a loja do cliente já está no ar, mas ainda sem movimento real de clientes. Os testes podem ser feitos criando pedidos manualmente na própria loja de produção, sem risco de afetar cliente real, sem precisar de cautelas extras de horário.

### Sincronismo de Produtos (Legado → Loja Integrada)

Direção contrária à seção de webhooks acima: o legado (`CADPRO`) é a fonte de verdade do catálogo, e o scheduler de polling do `LogAltera` empurra as mudanças pra Loja Integrada. Esta subseção estava só esboçada como uma tabela de endpoints (pesquisa de terceiros); as regras de negócio abaixo foram levantadas depois, numa rodada de acompanhamento pedida pelo próprio usuário ao revisar o escopo.

**Criação de produto novo**: hoje um produto nasce no legado (`CADPRO`) e o cadastro na Loja Integrada seria, em tese, manual. **Decisão**: o Lumini Hub cria o produto automaticamente na LI (via `POST /produto`) assim que ele for criado no legado, sem duplicar cadastro — e grava o `codint` retornado pela LI de volta em `CADPRO.codint` (é esse campo que faz a ponte entre os dois sistemas, conforme `CLAUDE.md`). **Nota de engenharia** (não é decisão de negócio, registrada aqui só pra rastreabilidade): como detectar "produto novo, ainda sem `codint`" é detalhe de implementação a refinar durante a construção — via evento `LogAltera` tipo `P` cruzado com uma checagem de `codint` vazio, ou uma varredura periódica separada de `CADPRO WHERE codint IS NULL`; não bloqueia o escopo.

**O que sincroniza** — confirmado pelo usuário, sem ressalvas (o `CADPRO` já tem tudo estruturado: categoria, imagens e variações/grade são fonte de verdade completa no legado, não algo cadastrado manualmente só dentro da LI):

| Dado | Sincroniza? | Observação |
|---|---|---|
| Preço de **venda** (`LogAltera.TipoAlt='V'`) | Sim | vai pra `/produto_preco` |
| **Custo** (`LogAltera.TipoAlt='U'`) | Não | informação interna, a loja virtual não expõe custo ao cliente |
| Estoque (`LogAltera.TipoAlt='P'`, ou o tipo que a implementação confirmar) | Sim, valor **pendente** | ver "Pendência de estoque" abaixo |
| Nome, descrição | Sim | via `/produto` |
| Categoria | Sim | via `/produto` ou `/categoria`, a confirmar durante implementação |
| Imagens | Sim | endpoint exato a confirmar durante implementação (a tabela de endpoints abaixo não é 100% oficial) |
| Variações / grade | Sim | via `/grade`, a confirmar durante implementação |

**Pendência de estoque — qual valor vai pra LI**: perguntado diretamente se o Lumini Hub deveria mandar o saldo "Oficial" (físico) cru ou aplicar algum colchão de segurança (reservar uma margem pra reduzir risco de venda concorrente entre loja física e virtual), o usuário preferiu **não decidir agora** — fica registrado como mais uma faceta da pendência já conhecida de prioridade de estoque físico vs. virtual (ver "Pendências em aberto" no fim desta seção). Sem essa decisão, o scheduler não pode ser implementado por completo na parte de estoque — pode ser construído e testado com preço/produto primeiro.

**Frequência do polling**: o usuário pediu "tempo real" (mudou no legado, sobe pra LI). Avaliado com ele que um mecanismo orientado a evento (trigger no SQL Server legado) contrariaria o princípio já registrado no `CLAUDE.md` de não alterar a estrutura do banco legado — a alternativa mais segura é manter a arquitetura de polling já prevista, só que num intervalo curto o suficiente pra parecer instantâneo na prática. **Decisão final**: em vez de travar um número fixo, o intervalo do polling vira **configurável** na tela Configurações → Integrações — novo campo `PollingIntervalSeconds` em `IntegrationConfig`, com **default sugerido de 60 segundos** já preenchido (o usuário pode reduzir ou aumentar depois, sem precisar de deploy novo).

**Tratamento de erro**: se o push de uma mudança (produto, preço ou estoque) pra LI falhar (ex.: API fora do ar, produto removido de lá manualmente), o scheduler faz **retry automático na próxima rodada de polling**, com **log de erro visível na tela de Configurações → Integrações** se a falha persistir — reaproveitando o mecanismo de `SyncLog` que já existe desde a Fase 1 (`Backend/microservices/api.integrations/internal/domain` — ver `tasks_integrations.md`).

**Consequência adicional no modelo de dados**: `IntegrationConfig` ganha mais um campo além dos já registrados na seção de Autenticação acima:

| Campo novo | Tipo | Observação |
|---|---|---|
| `PollingIntervalSeconds` | int | intervalo do scheduler de polling do `LogAltera`, editável em Configurações → Integrações, default 60 |

### Endpoints REST relevantes (produto/preço/estoque, pro scheduler de polling do `LogAltera`)

Levantados via doc de terceiro (não 100% oficial, mas consistente com o que o `api.integrations` já espera fazer) — a confirmar/ajustar durante a implementação do cliente HTTP da Loja Integrada, incluindo os endpoints exatos pra categoria/imagens/variações (não detalhados na fonte usada aqui):

| Recurso | Endpoint | Métodos |
|---|---|---|
| Produto (nome, descrição, categoria) | `/produto` | GET / POST / PUT |
| Preço | `/produto_preco` | GET / PUT |
| Estoque | `/produto_estoque` | GET / PUT |
| Pedido | `/pedido` | GET / PUT |
| Categoria / Marca / Grade (variações) | `/categoria`, `/marca`, `/grade` | GET |

Fonte: [Loja Integrada - Developers LinkApi](https://developers.linkapi.solutions/docs/lojaintegrada)

### Pendências em aberto (decisão do cliente/diretoria, não deste plano)

- [ ] Política de prioridade de estoque loja física vs. virtual quando os dois canais vendem a última unidade ao mesmo tempo — inclui **qual valor de estoque** o Lumini Hub deve enviar pra LI (Oficial cru vs. colchão de segurança), levantado na rodada de sincronismo de produtos.
- [ ] Regra de emissão de nota fiscal e despacho a partir da pré-venda pendente de faturamento gerada por um pedido pago da LI — o usuário precisa alinhar com o cliente.
- [ ] Confirmação prática de como a Loja Integrada permite registrar o segredo do webhook (header customizado vs. query string) — só será resolvido quando o usuário tentar cadastrar o webhook de fato.

> **Nota**: `codtipnot` saiu desta lista em 2026-09-27 — não é mais pendência de decisão do cliente/diretoria, virou campo de configuração normal (resolvido automaticamente pelo cadastro do vendedor padrão, ver "Configurações necessárias pra gerar a nota no legado" acima). Era um erro de classificação de uma versão anterior deste plano.
