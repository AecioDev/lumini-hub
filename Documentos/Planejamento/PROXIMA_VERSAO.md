# Próxima versão

Alterações feitas desde o último envio ao servidor. O Claude acrescenta uma linha aqui ao terminar cada tarefa. Quando você enviar a versão, peça o release: o Claude lê este arquivo, gera `releases/Release_<versão>.md` e zera esta lista.

- **Versão prevista:** a definir (próxima publicação: patch ou nova sequência se não mexer no banco; versão do banco +1 se mexer)
- **Última versão gerada:** `2026.3.4.0` (versão do banco 3), release em `releases/Release_2026.3.4.0.md`. Já publicada no ambiente de testes e conferida.
- **Formato:** `ano.versao_do_banco.sequencia_versao.sequencia_patch`. A versão do banco sobe uma unidade por **publicação** que contém mudanças de banco (quantas tabelas, campos ou triggers forem, na mesma publicação conta uma só); a seção "Mudanças que afetam o banco" abaixo lista o que mudou nessa versão.

## Mudanças que afetam o banco ou a configuração do servidor

- **Tabela nova `system_logs`** (auditoria, CFG-037): criada pelo `AutoMigrate` do `api.auth` (id, data, usuário, ação, entidade, IP e detalhes em jsonb, com índices). Muda o banco: a próxima versão sobe `dbversion.Required` de 3 para 4. Permissão nova `audit.view` (módulo Admin), criada e vinculada ao ADMIN no boot; o item de menu "Auditoria" é semeado sozinho em bancos já existentes.

## Novidades e melhorias

- **Auditoria (CFG-037, implementado, aguardando teste):** toda escrita bem-sucedida (criar, alterar, excluir) feita nos serviços `api.auth`, `api.core` e `api.integrations` é registrada com usuário, ação, entidade, ID e IP (buscas de lista, login, sessão e webhooks ficam de fora; o corpo da requisição nunca é gravado). Tela nova Configurações → Auditoria, só leitura, com busca por usuário/ação e filtros de entidade, ação e período; endpoints `POST /audit-logs/filter` e `GET /audit-logs/entity-types` (permissão `audit.view`). Cada registro guarda também **o que mudou**: numa alteração os campos que mudaram com o valor de antes e de depois (e, em listas como as permissões de um usuário, o que entrou e o que saiu), numa criação o objeto criado e numa exclusão o objeto excluído, vistos na gaveta "Ver detalhes". Senhas, tokens, hashes e certificados nunca são gravados. O que usuários do Perfil DEVELOP fazem só aparece para o próprio DEVELOP.

## Correções

- O `LoggerMiddleware` existia mas não estava ligado em nenhum serviço, e a tabela `system_logs` nunca era criada: agora é (ver Auditoria).

## Interno (não afeta o servidor)

- **CFG-044:** testes do `api.core` (`internal/service`, precisam de `TEST_DATABASE_URL`): empresa (CNPJ só com dígitos e único, matriz exige CNPJ, filial sem CNPJ, vínculo e ciclo, exclusão, logo, visibilidade por usuário), configuração fiscal (uma por empresa, tributação, certificado A1 com senha criptografada) e cliente/fornecedor (documento único com ou sem máscara, nome mínimo, busca com filtro, ordenação e paginação). Certificado autoassinado só de teste em `testdata/`.
- **CFG-041:** script `deploy/sql/CFG-041_textos_permissoes.sql` corrige os textos e módulos das permissões em bancos semeados antes do CFG-014 (só altera linhas que ainda têm o texto antigo exato, idempotente). Já aplicado no ambiente de testes (12 linhas); não muda a versão do banco. Qualquer outro banco semeado antigo precisa rodá-lo uma vez.
