---
name: lumini-hub-entity-creation
description: Passo a passo para criar e registrar novas entidades de domínio, repositories, services e rotas no backend do Lumini Hub.
---

# Diretrizes para criar entidades no Lumini Hub

Esta skill define o passo a passo para criar, escrever e registrar uma nova entidade de banco (como Marcas, Produtos etc.) dentro dos microsserviços do backend Go do Lumini Hub.

Execute sempre estes passos em ordem quando o usuário pedir para adicionar ou criar uma entidade nova.

---

## Visão geral do contexto dos microsserviços
Todos os caminhos de arquivo abaixo são relativos à pasta do microsserviço a que a entidade pertence:
*   **Entidades do Core (ERP/Admin/Cadastros)**: ficam em `Backend/microservices/api.core/`
*   **Entidades de Auth/Perfis**: ficam em `Backend/microservices/api.auth/`

---

## Passo 0: Escopo por empresa: pergunte antes de modelar

O Lumini Hub é multiempresa **dentro do banco de um único tenant** (o multi-tenant em si é um banco físico por cliente, tratado fora desta skill). Uma empresa matriz pode ter empresas filiais autorreferenciadas abaixo dela (`Company.ParentID`, mesmo formato do `MenuItem.ParentID`). Quase toda entidade precisa definir como se relaciona com `Company`, e errar isso sai caro de corrigir depois que já existem linhas: **nunca presuma, pergunte ao usuário em qual destes três grupos a nova entidade se encaixa antes de escrever o modelo do Passo 1**:

1. **Restrita à empresa**: a entidade tem a sua própria FK `CompanyID uint` e toda linha pertence a exatamente uma Company. Suposição padrão para tudo que for operacional/transacional e nunca puder se misturar entre empresas: estoque, vendas, pedidos de compra, lançamentos financeiros, notas fiscais, caixas etc. Repositories/services de entidades restritas precisam filtrar pela Company do solicitante (ver a regra de visibilidade abaixo) em toda consulta: isso é tão obrigatório quanto o padrão Unit of Work, não um reforço opcional.
2. **Global com relação de visibilidade**: sem `CompanyID` na própria entidade (o catálogo continua compartilhado/global), mas uma tabela de muitos-para-muitos fraca (`<entity>_companies`) controla quais Companies podem ver/usar cada linha. Caso de referência: `Product`: um catálogo compartilhado, mas uma tabela de junção `product_companies` controla quais produtos os vendedores de determinada Company podem ver/vender.
3. **Global, sem escopo**: nenhuma relação com Company. Reservado para dados de referência realmente comuns a todas as empresas (por exemplo, o catálogo de Permissões, tabelas genéricas de consulta). Justifique de forma explícita por que a entidade está aqui e não no grupo 1 ou 2: deve ser a exceção, não o padrão.

**Regra de visibilidade a aplicar quando a `Company` existir** (confirmada em 2026-07-21): o usuário tem um `CompanyID` padrão no próprio registro. Ele só vê dados da sua Company, A MENOS que tenha uma permissão que dê visibilidade às Companies subsidiárias abaixo da sua na hierarquia autorreferenciada (ou, se o usuário não tiver Company vinculada, essa mesma permissão decide se ele vê tudo): é uma exceção controlada por permissão, não automática só por estar mais acima na árvore.

Se a própria `Company` ainda não existir numa versão do código, faça mesmo assim a pergunta de escopo e registre a resposta no `plano_*.md` da funcionalidade, para que a FK `CompanyID`/a tabela de junção seja adicionada na mesma passada quando a `Company` chegar, em vez de ser esquecida em silêncio.

## Passo 0.5: Entidades de configuração 1:1 não têm Delete

Descoberto em 2026-09-05 ao construir o `CompanyFiscalConfig`: uma linha com soft delete somada a um índice único simples (por exemplo, `CompanyID uniqueIndex`) faz a linha apagada continuar ocupando o valor único, impedindo recriá-la depois (o `ExistsByX` do lado da aplicação ignora corretamente as linhas com soft delete pelo escopo padrão do GORM, mas o *índice* único cru do banco não ignora, a menos que seja um índice parcial com `where:deleted_at IS NULL`). A correção de verdade não é um índice mais esperto (embora seja um seguro barato, e vale colocá-lo mesmo assim): é reconhecer que **um registro de configuração 1:1 intrínseco a um pai (por exemplo, `CompanyFiscalConfig`, `CompanyVisualConfig`) não tem significado de negócio para "excluir"**: enquanto o pai (`Company`) existir, ele sempre tem (ou deveria ter) essa configuração, mesmo que vazia. Apagar a linha não corresponde a nada que um usuário realmente queira: "limpar o certificado" ou "restaurar as cores" são `Update`s com valores vazios/padrão, ou uma ação dedicada, nunca a exclusão da linha.

**Antes do Passo 6 (Handler), decida**: esta entidade é uma lista de verdade (as linhas realmente entram e saem: `Company`, `ChartOfAccounts`, `CompanyDocumentIssuanceConfig`) ou uma configuração 1:1 colada a um pai? Para o segundo caso, pule o Delete por inteiro: sem `Delete<Entity>` no service, sem método no handler, sem rota e sem a permissão `<entity>.delete` no Passo 11: apenas `view`/`create`/`edit`. Não construa "por completude" para remover depois; faça esta pergunta antes do Passo 1, no mesmo espírito do Passo 0.

---

## Passo a passo para criar uma entidade

### Passo 1: Modelo de domínio base
Crie a estrutura do modelo base em `internal/models/<entity>.go`.
*   Defina a struct que representa a tabela do banco com as tags do GORM (por exemplo, `gorm:"primaryKey"`, `gorm:"type:varchar(100)"`, `gorm:"uniqueIndex"`).
*   Crie os payloads de criação e atualização (por exemplo, `Create<Entity>Request`, `Update<Entity>Request`) com tags de validação (por exemplo, `binding:"required,min=2,max=100"`). Todo campo de texto com `gorm:"size:N"` precisa de `max=N` no `binding` (e no schema Zod do formulário); inclua a entidade no `request_size_test.go` do `api.core`.

### Passo 2: DTOs e mapeadores do modelo
Crie a estrutura de DTO em `internal/models/<entity>_dto.go`.
*   Inclua DTOs otimizados para listas e para detalhes (por exemplo, `<Entity>DTO`, `<Entity>DetailDTO`).
*   Implemente métodos/funções conversoras para mapear entidades em DTOs com segurança (por exemplo, `ToDTO()` e `ToDetailDTO()`).

### Passo 3: Interface de repository específica e implementação GORM
Crie o repository em `internal/repository/<entity>_repository.go`.
*   Defina uma interface específica que estenda `commonrepo.Repository[domain.<Entity>]`.
*   Forneça uma implementação em struct com GORM.
*   Se a entidade precisar de relacionamentos pré-carregados (por exemplo, `Addresses` ou `Contacts`), sobrescreva o `FindByID` ou escreva um método de busca próprio que aplique o `preload`.

### Passo 4: Validador de entrada e de negócio
Crie o validador em `internal/validator/<entity>_validator.go`.
*   Implemente validações que garantam as regras de negócio antes da persistência (por exemplo, unicidade de um código ou documento, validação de itens aninhados).

### Passo 5: Camada de service (lógica de negócio)
Crie o service em `internal/service/<entity>_service.go`.
*   Injete a interface do unit of work `repository.UnitOfWork` e o validador.
*   **Regra estrita**: use o Unit of Work (`s.uow.Execute(func(uow repository.UnitOfWork) error { ... })`) para envolver todas as operações de insert, update ou delete dos repositories. Isso garante que as transações sejam confirmadas de forma atômica.

### Passo 6: Camada de handler (controller)
Crie o handler em `internal/handlers/<entity>.go`.
*   Injete o Service.
*   Implemente os endpoints de Create, Update, Delete, GetByID e a busca (método POST `/filter`).
*   Use sempre `utils.SuccessResponse`, `utils.ErrorResponse` ou `utils.ValidationErrorResponse` para devolver os payloads. Num 500 o texto do erro não vai ao cliente (fica no log do servidor); em qualquer outro status o texto de erro do banco é trocado por uma mensagem genérica (e vai para o log), então valide tamanho, duplicidade e formato antes do banco.
*   Documente cada endpoint com as anotações declarativas padrão do Swagger (`@Summary`, `@Tags`, `@Router` etc.).

### Passo 7: Definição das rotas
Crie o arquivo de rotas em `internal/routes/<entity>_routes.go`.
*   Registre os endpoints no grupo do roteador Gin, aplicando os middlewares de Auth ou de permissão que forem necessários.

### Passo 8: Integração no servidor
Ligue o repository, o validator, o service, o handler e o registro de rotas no código de inicialização do próprio microsserviço (hoje o `main.go`: nem o `api.core` nem o `api.auth` têm ainda um pacote `internal/server/` separado, apesar de o `CLAUDE.md` citá-lo como convenção; confira o microsserviço atual antes de presumir que o arquivo existe).

### Passo 9: Migrações automáticas do banco
Adicione a nova struct à chamada `AutoMigrate(...)` do `main.go` do microsserviço (por exemplo, `api.core/main.go` para a entidade `Company`; no `api.auth` é o `prepareDatabase` do `main.go`), na ordem de dependência das chaves estrangeiras. Ela só roda quando `DB_AUTO_MIGRATE` está ligada. Não existe um arquivo `migrate.go` dedicado no `api.core` nem no `api.auth`.

### Passo 10: Seed do banco
Adicione dados padrão de seed (menu, configuração do sistema) pelo pacote `internal/seeder/` do microsserviço, se ele tiver um (o `api.auth` tem, por exemplo o `menu_item_seeder.go`; o `api.core` ainda não: siga o mesmo padrão de um pacote por microsserviço se você for a primeira entidade ali que precisa de dados de seed). O seed precisa ser idempotente e rodar sob a mesma flag `DB_AUTO_MIGRATE`. As permissões **não** entram aqui: ficam no catálogo do Passo 11.

### Passo 11: Segurança e permissões RBAC
Registre as permissões da nova entidade (`<entity>.view`, `<entity>.create`, `<entity>.edit`, `<entity>.delete`) no `PermissionCatalog` (`Backend/microservices/api.auth/seed/permission_catalog.go`), com uma descrição em português completa e o módulo a que pertencem (por exemplo "Criar fornecedores", módulo "Compras"). O próximo boot do `api.auth` as cria e as vincula ao perfil `ADMIN` automaticamente (tudo que não for do módulo `Develop`), sem tocar em vínculos que já existem. Toda permissão que uma rota exige com `RequirePermission("...")` precisa estar no catálogo. Atualize os pisos do teste de âncora (`catalog_test.go`) e o texto do `CLAUDE.md` se as contagens mudarem.

### Passo 12: Testes unitários (recomendado)
Crie testes unitários/de integração em `internal/service/<entity>_service_test.go` ou `internal/validator/<entity>_validator_test.go` para testar a lógica de forma isolada.

---

## Regras importantes para a geração de código por IA
1.  **Não pule camadas**: nunca ligue controllers diretamente a repositories. Implemente sempre a camada de Service e faça as operações passarem pelo Unit of Work.
2.  **Sem GORM direto nos Services**: os services só podem interagir com o `repository.UnitOfWork`.
3.  **Nomes de arquivo estritos**: mantenha os nomes em snake_case (por exemplo, `customer_supplier_repository.go`) para seguir a convenção dos arquivos existentes.
4.  **Regerar o Swagger**: sempre que alterar ou adicionar handlers, regere a documentação do Swagger na pasta do gateway:
    ```powershell
    swag init -g main.go -d ./,../../common,../api.auth,../api.core --parseDependency
    ```
5.  **Permissões no catálogo, não em seed solto**: ao gerar os arquivos de uma entidade nova, você **DEVE** acrescentar as permissões (`<entity>.view`, `<entity>.create`, `<entity>.edit`, `<entity>.delete`) ao `PermissionCatalog` na mesma passada (Passo 11). É assim que elas chegam ao banco e ao `ADMIN`.
6.  **Faça a pergunta de escopo por empresa antes do Passo 1, toda vez**: restrita à empresa (com `CompanyID` próprio), global com relação de visibilidade (tabela de junção) ou global sem escopo: ver o Passo 0 acima. Não adote "sem escopo" como padrão só porque a `Company` ainda não existe no código; registre a resposta pretendida no `plano_*.md` de qualquer forma.
7.  **Entidades de configuração 1:1 pulam o Delete por inteiro** (service, handler, rota e a permissão `.delete`): ver o Passo 0.5. Só as entidades em formato de lista de verdade mantêm o CRUD completo.
