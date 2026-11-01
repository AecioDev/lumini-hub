# Próxima versão

Alterações feitas desde o último envio ao servidor. O Claude acrescenta uma linha aqui ao terminar cada tarefa. Quando você enviar a versão, peça o release: o Claude lê este arquivo, gera `releases/Release_<versão>.md` e zera esta lista.

- **Versão prevista:** a definir (próxima publicação: patch ou nova sequência se não mexer no banco; versão do banco +1 se mexer)
- **Última versão gerada:** `2026.3.4.0` (versão do banco 3), release em `releases/Release_2026.3.4.0.md`. Já publicada no ambiente de testes e conferida.
- **Formato:** `ano.versao_do_banco.sequencia_versao.sequencia_patch`. A versão do banco sobe uma unidade por **publicação** que contém mudanças de banco (quantas tabelas, campos ou triggers forem, na mesma publicação conta uma só); a seção "Mudanças que afetam o banco" abaixo lista o que mudou nessa versão.

## Mudanças que afetam o banco ou a configuração do servidor

- **Tabela nova `system_logs`** (auditoria, CFG-037): criada pelo `AutoMigrate` do `api.auth` (id, data, usuário, ação, entidade, IP e detalhes em jsonb, com índices). Muda o banco: a próxima versão sobe `dbversion.Required` de 3 para 4. Permissão nova `audit.view` (módulo Admin), criada e vinculada ao ADMIN no boot; o item de menu "Auditoria" é semeado sozinho em bancos já existentes.

- **Índices únicos parciais (CFG-050):** `customers.document_number` e `suppliers.document_number` trocam a constraint única comum por índice único `WHERE deleted_at IS NULL` (`idx_customers_document_number`, `idx_suppliers_document_number`), e `companies.cnpj` ganha o índice único parcial `idx_companies_cnpj` (antes só o validador conferia). Os índices novos são criados pelo `AutoMigrate`; as constraints antigas (`uni_customers_document_number`, `uni_suppliers_document_number`) são derrubadas pela migração SQL `0002` (ver abaixo), sem passo manual. Entra na mesma versão do banco (4) da auditoria.
- **Migrações SQL versionadas (CFG-051):** tabela nova `schema_migrations` (criada pelo `api.auth`). Os arquivos de `Backend/common/dbversion/migrations/` rodam sozinhos no boot do `api.auth` (com `DB_AUTO_MIGRATE=true`), uma vez por banco, e quem tem `DB_AUTO_MIGRATE=false` precisa subir o `api.auth` uma vez com a flag ligada (o `dbversion.Check` recusa o boot com migração pendente). Nesta versão: `0001` (textos das permissões, CFG-041, no-op em banco que já recebeu o script à mão) e `0002` (índices únicos parciais, CFG-050). Se houver CNPJ de empresa ativa duplicado, a `0002` falha, o `api.auth` não sobe e a mensagem diz qual migração: corrija os duplicados e suba de novo.

## Novidades e melhorias

- **Auditoria (CFG-037, implementado, aguardando teste):** toda escrita bem-sucedida (criar, alterar, excluir) feita nos serviços `api.auth`, `api.core` e `api.integrations` é registrada com usuário, ação, entidade, ID e IP (buscas de lista, login, sessão e webhooks ficam de fora; o corpo da requisição nunca é gravado). Tela nova Configurações → Auditoria, só leitura, com busca por usuário/ação e filtros de entidade, ação e período; endpoints `POST /audit-logs/filter` e `GET /audit-logs/entity-types` (permissão `audit.view`). Cada registro guarda também **o que mudou**: numa alteração os campos que mudaram com o valor de antes e de depois (e, em listas como as permissões de um usuário, o que entrou e o que saiu), numa criação o objeto criado e numa exclusão o objeto excluído, vistos na gaveta "Ver detalhes". Senhas, tokens, hashes e certificados nunca são gravados. O que usuários do Perfil DEVELOP fazem só aparece para o próprio DEVELOP.

## Correções

- **CFG-050:** recadastrar o documento de um cliente ou fornecedor excluído, ou o CNPJ de uma empresa excluída, agora funciona (antes passava na validação e falhava no banco com "Não foi possível concluir a operação."); duas criações simultâneas com o mesmo CNPJ de empresa também deixam de duplicar.
- **CFG-049 (segurança):** o cadastro de fornecedor deixou de aceitar `created_by` no corpo da requisição; o autor do registro é sempre o usuário logado (antes, quem chamasse a API podia atribuir a criação a outro usuário). Sem mudança de banco; o campo some do Swagger da criação.
- O `LoggerMiddleware` existia mas não estava ligado em nenhum serviço, e a tabela `system_logs` nunca era criada: agora é (ver Auditoria).

## Interno (não afeta o servidor)

- **CFG-044:** testes do `api.core` (`internal/service`, precisam de `TEST_DATABASE_URL`): empresa (CNPJ só com dígitos e único, matriz exige CNPJ, filial sem CNPJ, vínculo e ciclo, exclusão, logo, visibilidade por usuário), configuração fiscal (uma por empresa, tributação, certificado A1 com senha criptografada) e cliente/fornecedor (documento único com ou sem máscara, nome mínimo, busca com filtro, ordenação e paginação). Certificado autoassinado só de teste em `testdata/`.
- **CFG-041:** o script (hoje a migração `0001`) corrige os textos e módulos das permissões em bancos semeados antes do CFG-014 (só altera linhas que ainda têm o texto antigo exato, idempotente). Já aplicado no ambiente de testes (12 linhas); não muda a versão do banco. Qualquer outro banco semeado antigo precisa rodá-lo uma vez.
