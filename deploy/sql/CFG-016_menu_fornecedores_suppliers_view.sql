-- CFG-016: itens de menu "Fornecedores" (href /suppliers) em bancos já semeados.
--
-- O seed do menu (seeder.SeedMenuItems) pula itens que já existem, então um banco semeado antes
-- de 2026-10-20 continua com esses itens ligados a purchases.view (Compras). Fornecedores é o
-- cadastro (suppliers.*); purchases.* é a rotina de Compras. Bancos novos já nascem certos.
--
-- Idempotente: só altera o que ainda não está em suppliers.view, e não faz nada se a permissão
-- suppliers.view não existir (suba o api.auth uma vez, que o seed a cria).
-- Uso: psql -U <usuario> -d <banco> -f deploy/sql/CFG-016_menu_fornecedores_suppliers_view.sql

BEGIN;

UPDATE menu_items
SET permission_id = (SELECT id FROM permissions WHERE permission = 'suppliers.view' AND deleted_at IS NULL)
WHERE href = '/suppliers'
  AND deleted_at IS NULL
  AND EXISTS (SELECT 1 FROM permissions WHERE permission = 'suppliers.view' AND deleted_at IS NULL)
  AND permission_id IS DISTINCT FROM
      (SELECT id FROM permissions WHERE permission = 'suppliers.view' AND deleted_at IS NULL);

COMMIT;
