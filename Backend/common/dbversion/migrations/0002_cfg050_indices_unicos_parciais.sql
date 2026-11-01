-- CFG-050: unicidade do banco alinhada ao soft delete.
--
-- Antes: customers.document_number e suppliers.document_number tinham constraint UNIQUE comum, então um
-- registro excluído (soft delete) continuava ocupando o documento e recadastrá-lo falhava no banco; e
-- companies.cnpj não tinha índice único (só o validador conferia), então duas criações simultâneas
-- podiam duplicar o CNPJ.
-- Agora: índices únicos parciais (WHERE deleted_at IS NULL), com os mesmos nomes que o AutoMigrate cria
-- (idx_customers_document_number, idx_suppliers_document_number, idx_companies_cnpj). O AutoMigrate
-- nunca remove a constraint antiga, por isso esta migração a derruba em bancos já existentes.
--
-- As tabelas pertencem ao api.core: num banco novo, em que o api.auth sobe antes, elas ainda não existem
-- e o AutoMigrate do api.core já as cria com os índices certos, então cada bloco só age se a tabela existe.
-- Se já houver CNPJ de empresa duplicado entre as ativas, o CREATE INDEX falha e a migração inteira é
-- desfeita: corrija os duplicados e suba de novo.

DO $$
BEGIN
  IF to_regclass('customers') IS NOT NULL THEN
    ALTER TABLE customers DROP CONSTRAINT IF EXISTS uni_customers_document_number;
    CREATE UNIQUE INDEX IF NOT EXISTS idx_customers_document_number ON customers (document_number) WHERE deleted_at IS NULL;
  END IF;
  IF to_regclass('suppliers') IS NOT NULL THEN
    ALTER TABLE suppliers DROP CONSTRAINT IF EXISTS uni_suppliers_document_number;
    CREATE UNIQUE INDEX IF NOT EXISTS idx_suppliers_document_number ON suppliers (document_number) WHERE deleted_at IS NULL;
  END IF;
  IF to_regclass('companies') IS NOT NULL THEN
    CREATE UNIQUE INDEX IF NOT EXISTS idx_companies_cnpj ON companies (cnpj) WHERE deleted_at IS NULL;
  END IF;
END
$$;
