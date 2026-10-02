--
-- PostgreSQL database dump
--

-- Dumped from database version 16.7
-- Dumped by pg_dump version 16.7

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Data for Name: permissions; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (1, 'users.view', 'Visualizar usuários', 'Usuários', '2026-05-21 23:18:26.024052-03', '2026-05-21 23:18:26.024052-03', NULL);
INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (2, 'users.create', 'Criar usuários', 'Usuários', '2026-05-21 23:18:26.024052-03', '2026-05-21 23:18:26.024052-03', NULL);
INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (3, 'users.edit', 'Editar usuários', 'Usuários', '2026-05-21 23:18:26.024052-03', '2026-05-21 23:18:26.024052-03', NULL);
INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (4, 'users.delete', 'Excluir usuários', 'Usuários', '2026-05-21 23:18:26.024052-03', '2026-05-21 23:18:26.024052-03', NULL);
INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (5, 'sales.view', 'Visualizar vendas', 'Vendas', '2026-05-21 23:18:26.026195-03', '2026-05-21 23:18:26.026195-03', NULL);
INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (6, 'sales.create', 'Criar vendas', 'Vendas', '2026-05-21 23:18:26.026195-03', '2026-05-21 23:18:26.026195-03', NULL);
INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (7, 'sales.edit', 'Editar vendas', 'Vendas', '2026-05-21 23:18:26.026195-03', '2026-05-21 23:18:26.026195-03', NULL);
INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (8, 'sales.delete', 'Excluir vendas', 'Vendas', '2026-05-21 23:18:26.026195-03', '2026-05-21 23:18:26.026195-03', NULL);
INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (9, 'sales.reports', 'Gerar relatórios de vendas', 'Vendas', '2026-05-21 23:18:26.026195-03', '2026-05-21 23:18:26.026195-03', NULL);
INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (10, 'purchases.view', 'Visualizar compras', 'Compras', '2026-05-21 23:18:26.026609-03', '2026-05-21 23:18:26.026609-03', NULL);
INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (11, 'purchases.create', 'Criar compras', 'Compras', '2026-05-21 23:18:26.026609-03', '2026-05-21 23:18:26.026609-03', NULL);
INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (12, 'purchases.edit', 'Editar compras', 'Compras', '2026-05-21 23:18:26.026609-03', '2026-05-21 23:18:26.026609-03', NULL);
INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (13, 'purchases.delete', 'Excluir compras', 'Compras', '2026-05-21 23:18:26.026609-03', '2026-05-21 23:18:26.026609-03', NULL);
INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (14, 'purchases.reports', 'Gerar relatórios de compras', 'Compras', '2026-05-21 23:18:26.026609-03', '2026-05-21 23:18:26.026609-03', NULL);
INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (15, 'inventory.view', 'Visualizar estoque', 'Estoque', '2026-05-21 23:18:26.02699-03', '2026-05-21 23:18:26.02699-03', NULL);
INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (16, 'inventory.create', 'Adicionar itens ao estoque', 'Estoque', '2026-05-21 23:18:26.02699-03', '2026-05-21 23:18:26.02699-03', NULL);
INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (17, 'inventory.edit', 'Editar itens do estoque', 'Estoque', '2026-05-21 23:18:26.02699-03', '2026-05-21 23:18:26.02699-03', NULL);
INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (18, 'inventory.delete', 'Remover itens do estoque', 'Estoque', '2026-05-21 23:18:26.02699-03', '2026-05-21 23:18:26.02699-03', NULL);
INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (19, 'inventory.reports', 'Gerar relatórios de estoque', 'Estoque', '2026-05-21 23:18:26.02699-03', '2026-05-21 23:18:26.02699-03', NULL);
INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (20, 'finance.view', 'Visualizar finanças', 'Financeiro', '2026-05-21 23:18:26.027399-03', '2026-05-21 23:18:26.027399-03', NULL);
INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (21, 'finance.create', 'Criar transações financeiras', 'Financeiro', '2026-05-21 23:18:26.027399-03', '2026-05-21 23:18:26.027399-03', NULL);
INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (22, 'finance.edit', 'Editar transações financeiras', 'Financeiro', '2026-05-21 23:18:26.027399-03', '2026-05-21 23:18:26.027399-03', NULL);
INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (23, 'finance.delete', 'Excluir transações financeiras', 'Financeiro', '2026-05-21 23:18:26.027399-03', '2026-05-21 23:18:26.027399-03', NULL);
INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (24, 'finance.reports', 'Gerar relatórios financeiros', 'Financeiro', '2026-05-21 23:18:26.027399-03', '2026-05-21 23:18:26.027399-03', NULL);
INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (36, 'dashboard.view_default', 'Visualizar o dashboard padrão do perfil', 'dashboard', '2026-05-21 23:35:14.70263-03', '2026-05-21 23:35:14.70263-03', NULL);
INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (51, 'integrations.view', 'Visualizar Integrações', 'Integrações', '2026-07-18 01:34:55.258388-03', '2026-07-18 01:34:55.258388-03', NULL);
INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (27, 'orders.view', 'Visualizar pedidos de vendas', 'Vendas', '2026-05-21 23:35:14.697118-03', '2026-05-21 23:35:14.697118-03', NULL);
INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (54, 'companies.view', 'Visualizar empresas', 'Empresas', '2026-07-21 23:57:45.130635-03', '2026-07-21 23:57:45.130635-03', NULL);
INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (55, 'companies.create', 'Criar empresas', 'Empresas', '2026-07-21 23:57:45.130635-03', '2026-07-21 23:57:45.130635-03', NULL);
INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (56, 'companies.edit', 'Editar empresas', 'Empresas', '2026-07-21 23:57:45.130635-03', '2026-07-21 23:57:45.130635-03', NULL);
INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (57, 'companies.delete', 'Excluir empresas', 'Empresas', '2026-07-21 23:57:45.130635-03', '2026-07-21 23:57:45.130635-03', NULL);
INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (58, 'companies.fiscal_config.view', 'Visualizar configuração fiscal de empresas', 'Empresas', '2026-09-05 14:45:25.309069-03', '2026-09-05 14:45:25.309069-03', NULL);
INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (59, 'companies.fiscal_config.create', 'Criar configuração fiscal de empresas', 'Empresas', '2026-09-05 14:45:25.309069-03', '2026-09-05 14:45:25.309069-03', NULL);
INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (37, 'dashboard.admin.view', 'Visualizar o dashboard do administrador', 'Admin', '2026-05-21 23:35:14.705171-03', '2026-05-21 23:35:14.705171-03', NULL);
INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (38, 'dashboard.sales.view', 'Visualizar o dashboard de vendas', 'Vendas', '2026-05-21 23:35:14.706137-03', '2026-05-21 23:35:14.706137-03', NULL);
INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (39, 'dashboard.finance.view', 'Visualizar o dashboard financeiro', 'Financeiro', '2026-05-21 23:35:14.706643-03', '2026-05-21 23:35:14.706643-03', NULL);
INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (40, 'dashboard.manager.view', 'Visualizar o dashboard gerencial', 'Gerencial', '2026-05-21 23:35:14.706643-03', '2026-05-21 23:35:14.706643-03', NULL);
INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (41, 'dashboard.inventory.view', 'Visualizar o dashboard de estoque', 'Estoque', '2026-05-21 23:35:14.708168-03', '2026-05-21 23:35:14.708168-03', NULL);
INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (60, 'companies.fiscal_config.edit', 'Editar configuração fiscal de empresas (inclui upload de certificado)', 'Empresas', '2026-09-05 14:45:25.309069-03', '2026-09-05 14:45:25.309069-03', NULL);
INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (43, 'permissions.view', 'Visualizar Permissões', 'Develop', '2026-05-21 23:35:14.710021-03', '2026-05-21 23:35:14.710021-03', NULL);
INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (44, 'permissions.create', 'Criar Permissões', 'Develop', '2026-05-21 23:35:14.710592-03', '2026-05-21 23:35:14.710592-03', NULL);
INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (45, 'permissions.edit', 'Editar Permissões', 'Develop', '2026-05-21 23:35:14.710592-03', '2026-05-21 23:35:14.710592-03', NULL);
INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (46, 'permissions.delete', 'Excluir Permissões', 'Develop', '2026-05-21 23:35:14.71112-03', '2026-05-21 23:35:14.71112-03', NULL);
INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (47, 'roles.view', 'Visualizar usuários', 'Develop', '2026-05-21 23:35:14.71112-03', '2026-05-21 23:35:14.71112-03', NULL);
INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (48, 'roles.create', 'Criar usuários', 'Develop', '2026-05-21 23:35:14.711757-03', '2026-05-21 23:35:14.711757-03', NULL);
INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (49, 'roles.edit', 'Editar usuários', 'Develop', '2026-05-21 23:35:14.711757-03', '2026-05-21 23:35:14.711757-03', NULL);
INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (50, 'roles.delete', 'Excluir usuários', 'Develop', '2026-05-21 23:35:14.712291-03', '2026-05-21 23:35:14.712291-03', NULL);
INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (28, 'products.view', 'Visualizar produtos', 'Estoque', '2026-05-21 23:35:14.698255-03', '2026-05-21 23:35:14.698255-03', NULL);
INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (29, 'supplier_codes.view', 'Visualizar códigos por fornecedor', 'Estoque', '2026-05-21 23:35:14.698849-03', '2026-05-21 23:35:14.698849-03', NULL);
INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (30, 'stock_locations.view', 'Visualizar locais de estoque', 'Estoque', '2026-05-21 23:35:14.699488-03', '2026-05-21 23:35:14.699488-03', NULL);
INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (31, 'product_location.view', 'Visualizar localização de produtos', 'Estoque', '2026-05-21 23:35:14.699488-03', '2026-05-21 23:35:14.699488-03', NULL);
INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (32, 'prices_promotions.view', 'Visualizar preços e promoções', 'Estoque', '2026-05-21 23:35:14.70006-03', '2026-05-21 23:35:14.70006-03', NULL);
INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (33, 'taxation.view', 'Visualizar tributação', 'Estoque', '2026-05-21 23:35:14.70006-03', '2026-05-21 23:35:14.70006-03', NULL);
INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (25, 'customers.view', 'Visualizar clientes', 'Vendas', '2026-05-21 23:35:14.695158-03', '2026-05-21 23:35:14.695158-03', NULL);
INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (26, 'payment_plans.view', 'Visualizar planos de pagamento', 'Vendas', '2026-05-21 23:35:14.696361-03', '2026-05-21 23:35:14.696361-03', NULL);
INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (34, 'finance.receive_boleto', 'Permissão para receber boleto financeiro', 'Financeiro', '2026-05-21 23:35:14.701545-03', '2026-05-21 23:35:14.701545-03', NULL);
INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (35, 'finance.view_pendencies', 'Visualizar pendências financeiras', 'Financeiro', '2026-05-21 23:35:14.70211-03', '2026-05-21 23:35:14.70211-03', NULL);
INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (52, 'admin.create_permissions', 'Gerenciar perfis e permissões (catálogo)', 'Develop', '2026-07-21 14:55:04.502731-03', '2026-07-21 14:55:04.502731-03', NULL);
INSERT INTO public.permissions (id, permission, description, module, created_at, updated_at, deleted_at) VALUES (62, 'companies.hierarchy.view', 'Visualizar empresas subsidiárias na hierarquia', 'Empresas', '2026-09-05 19:08:35.886228-03', '2026-09-05 19:08:35.886228-03', NULL);


--
-- Data for Name: roles; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.roles (id, name, description, created_at, updated_at, deleted_at) VALUES (6, 'Gerente', 'Acesso gerencial a múltiplos módulos', '2026-05-21 23:18:26.021806-03', '2026-05-21 23:18:26.021806-03', NULL);
INSERT INTO public.roles (id, name, description, created_at, updated_at, deleted_at) VALUES (2, 'Vendas', 'Acesso ao módulo de vendas', '2026-05-21 23:18:26.021806-03', '2026-07-21 15:33:30.792689-03', NULL);
INSERT INTO public.roles (id, name, description, created_at, updated_at, deleted_at) VALUES (1, 'ADMIN', 'Administrador do sistema', '2026-05-21 23:18:26.021806-03', '2026-07-21 16:14:25.192524-03', NULL);
INSERT INTO public.roles (id, name, description, created_at, updated_at, deleted_at) VALUES (5, 'Financeiro', 'Acesso ao módulo financeiro', '2026-05-21 23:18:26.021806-03', '2026-07-21 17:35:30.299084-03', NULL);
INSERT INTO public.roles (id, name, description, created_at, updated_at, deleted_at) VALUES (3, 'Compras', 'Acesso ao módulo de compras', '2026-05-21 23:18:26.021806-03', '2026-07-21 18:29:04.475444-03', NULL);
INSERT INTO public.roles (id, name, description, created_at, updated_at, deleted_at) VALUES (4, 'Estoque', 'Acesso ao módulo de estoque', '2026-05-21 23:18:26.021806-03', '2026-07-21 19:26:06.288742-03', NULL);
INSERT INTO public.roles (id, name, description, created_at, updated_at, deleted_at) VALUES (7, 'DEVELOP', 'Acesso exclusivo para desenvolvedores', '2026-07-20 00:01:44.066526-03', '2026-07-20 00:01:44.066526-03', NULL);


--
-- Data for Name: role_permissions; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (1, 1, '2026-05-21 23:18:26.027782-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (1, 2, '2026-05-21 23:18:26.027782-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (1, 3, '2026-05-21 23:18:26.027782-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (1, 4, '2026-05-21 23:18:26.027782-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (1, 5, '2026-05-21 23:18:26.027782-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (1, 6, '2026-05-21 23:18:26.027782-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (1, 7, '2026-05-21 23:18:26.027782-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (1, 8, '2026-05-21 23:18:26.027782-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (1, 9, '2026-05-21 23:18:26.027782-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (1, 10, '2026-05-21 23:18:26.027782-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (1, 11, '2026-05-21 23:18:26.027782-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (1, 12, '2026-05-21 23:18:26.027782-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (1, 13, '2026-05-21 23:18:26.027782-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (1, 14, '2026-05-21 23:18:26.027782-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (1, 15, '2026-05-21 23:18:26.027782-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (1, 16, '2026-05-21 23:18:26.027782-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (1, 17, '2026-05-21 23:18:26.027782-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (1, 18, '2026-05-21 23:18:26.027782-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (1, 19, '2026-05-21 23:18:26.027782-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (1, 20, '2026-05-21 23:18:26.027782-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (1, 21, '2026-05-21 23:18:26.027782-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (1, 22, '2026-05-21 23:18:26.027782-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (1, 23, '2026-05-21 23:18:26.027782-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (1, 24, '2026-05-21 23:18:26.027782-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (2, 5, '2026-05-21 23:18:26.030414-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (2, 6, '2026-05-21 23:18:26.030414-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (2, 7, '2026-05-21 23:18:26.030414-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (2, 8, '2026-05-21 23:18:26.030414-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (2, 9, '2026-05-21 23:18:26.030414-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (3, 10, '2026-05-21 23:18:26.03126-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (3, 11, '2026-05-21 23:18:26.03126-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (3, 12, '2026-05-21 23:18:26.03126-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (3, 13, '2026-05-21 23:18:26.03126-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (3, 14, '2026-05-21 23:18:26.03126-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (4, 15, '2026-05-21 23:18:26.031774-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (4, 16, '2026-05-21 23:18:26.031774-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (4, 17, '2026-05-21 23:18:26.031774-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (4, 18, '2026-05-21 23:18:26.031774-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (4, 19, '2026-05-21 23:18:26.031774-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (5, 20, '2026-05-21 23:18:26.033904-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (5, 21, '2026-05-21 23:18:26.033904-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (5, 22, '2026-05-21 23:18:26.033904-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (5, 23, '2026-05-21 23:18:26.033904-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (5, 24, '2026-05-21 23:18:26.033904-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (6, 1, '2026-05-21 23:18:26.034533-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (6, 5, '2026-05-21 23:18:26.034533-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (6, 9, '2026-05-21 23:18:26.034533-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (6, 10, '2026-05-21 23:18:26.034533-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (6, 14, '2026-05-21 23:18:26.034533-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (6, 15, '2026-05-21 23:18:26.034533-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (6, 19, '2026-05-21 23:18:26.034533-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (6, 20, '2026-05-21 23:18:26.034533-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (6, 24, '2026-05-21 23:18:26.034533-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (1, 25, '2026-05-21 23:35:14.685726-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (1, 26, '2026-05-21 23:35:14.685726-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (1, 27, '2026-05-21 23:35:14.685726-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (1, 28, '2026-05-21 23:35:14.685726-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (1, 29, '2026-05-21 23:35:14.685726-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (1, 30, '2026-05-21 23:35:14.685726-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (1, 31, '2026-05-21 23:35:14.685726-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (1, 32, '2026-05-21 23:35:14.685726-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (1, 33, '2026-05-21 23:35:14.685726-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (1, 34, '2026-05-21 23:35:14.685726-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (1, 35, '2026-05-21 23:35:14.685726-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (1, 36, '2026-05-21 23:35:14.685726-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (1, 37, '2026-05-21 23:35:14.685726-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (1, 38, '2026-05-21 23:35:14.685726-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (1, 39, '2026-05-21 23:35:14.685726-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (1, 40, '2026-05-21 23:35:14.685726-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (1, 41, '2026-05-21 23:35:14.685726-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (1, 51, '2026-07-20 00:23:03.232952-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (2, 27, '2026-07-21 15:33:30.792653-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (2, 38, '2026-07-21 15:33:30.792653-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (2, 25, '2026-07-21 15:33:30.792653-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (2, 26, '2026-07-21 15:33:30.792653-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (5, 39, '2026-07-21 17:35:30.299336-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (5, 34, '2026-07-21 17:35:30.299336-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (5, 35, '2026-07-21 17:35:30.299336-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (4, 41, '2026-07-21 18:33:14.362829-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (4, 31, '2026-07-21 18:33:14.362829-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (4, 10, '2026-07-21 19:26:06.288906-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (4, 14, '2026-07-21 19:26:06.288906-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (1, 54, '2026-07-21 23:57:45.130635-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (1, 55, '2026-07-21 23:57:45.130635-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (1, 56, '2026-07-21 23:57:45.130635-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (1, 57, '2026-07-21 23:57:45.130635-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (1, 58, '2026-09-05 14:45:25.309069-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (1, 59, '2026-09-05 14:45:25.309069-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (1, 60, '2026-09-05 14:45:25.309069-03');
INSERT INTO public.role_permissions (role_id, permission_id, created_at) VALUES (1, 62, '2026-09-05 19:08:35.889343-03');


--
-- Name: permissions_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.permissions_id_seq', 65, true);


--
-- Name: roles_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.roles_id_seq', 12, true);


--
-- PostgreSQL database dump complete
--

