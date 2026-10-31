-- CFG-041: textos e módulos do catálogo de permissões em bancos já semeados.
--
-- O seed (seed.SyncCatalog) só CRIA permissões ausentes e nunca altera uma existente. Os textos
-- corrigidos no código (CFG-014: roles.* = perfis, "Receber boletos", permissions.*/integrations.*
-- em minúsculas, módulo Dashboard) por isso só valem em banco novo. Este script leva um banco já
-- semeado antes dessa correção para os mesmos textos do catálogo (api.auth/seed/permission_catalog.go).
--
-- Idempotente e conservador: cada linha só é alterada se ainda tiver EXATAMENTE o módulo e a
-- descrição antigos conhecidos. Uma permissão que o operador já ajustou à mão, ou que já está
-- certa, não é tocada; rodar de novo não faz nada. O código da permissão nunca muda.
-- Uso: psql -U <usuario> -d <banco> -f deploy/sql/CFG-041_textos_permissoes.sql
-- (o cliente precisa estar em UTF-8: PGCLIENTENCODING=UTF8).

BEGIN;

UPDATE permissions AS p
SET module = fix.new_module,
    description = fix.new_description
FROM (VALUES
  ('dashboard.view_default', 'dashboard',   'Visualizar o dashboard padrão do perfil', 'Dashboard',    'Visualizar o dashboard padrão do perfil'),
  ('permissions.create',     'Develop',     'Criar Permissões',                        'Develop',      'Criar permissões'),
  ('permissions.delete',     'Develop',     'Excluir Permissões',                      'Develop',      'Excluir permissões'),
  ('permissions.edit',       'Develop',     'Editar Permissões',                       'Develop',      'Editar permissões'),
  ('permissions.view',       'Develop',     'Visualizar Permissões',                   'Develop',      'Visualizar permissões'),
  ('roles.create',           'Develop',     'Criar usuários',                          'Develop',      'Criar perfis'),
  ('roles.delete',           'Develop',     'Excluir usuários',                        'Develop',      'Excluir perfis'),
  ('roles.edit',             'Develop',     'Editar usuários',                         'Develop',      'Editar perfis'),
  ('roles.view',             'Develop',     'Visualizar usuários',                     'Develop',      'Visualizar perfis'),
  ('finance.receive_boleto', 'Financeiro',  'Permissão para receber boleto financeiro', 'Financeiro',  'Receber boletos'),
  ('integrations.edit',      'Integrações', 'Editar Integrações',                      'Integrações',  'Editar integrações'),
  ('integrations.view',      'Integrações', 'Visualizar Integrações',                  'Integrações',  'Visualizar integrações')
) AS fix(code, old_module, old_description, new_module, new_description)
WHERE p.permission = fix.code
  AND p.module = fix.old_module
  AND p.description = fix.old_description
  AND p.deleted_at IS NULL;

COMMIT;
