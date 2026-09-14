# Módulo Multi-Tenant - Estratégia de Isolamento entre Clientes da Lumini Hub

Este módulo trata do isolamento de dados entre **clientes/tenants diferentes** da Lumini Hub (cliente A vs. cliente B) — decisão já tomada informalmente: **cada cliente terá um banco de dados físico separado**. O *como* implementar isso (roteamento de conexão, provisionamento, deploy) ainda não foi desenhado; hoje o sistema roda com um único banco, um único cliente.

**Importante — isso é ortogonal ao Módulo 0 / Empresa**: `Company` (Matriz + empresas vinculadas, ver `Modulo_0_Configuracao/plano_empresa.md`) é a estrutura **dentro** do banco de um único tenant (um cliente da Lumini Hub pode ter várias empresas/CNPJs geridas pelo mesmo gestor, todas no mesmo banco). Multi-tenant é o isolamento **entre** bancos de clientes diferentes. Os dois conceitos não se sobrepõem e não devem ser modelados juntos.

---

## ❓ Perguntas em aberto (sem decisão tomada, sem desenho de solução — revisitar quando um segundo cliente real estiver no horizonte)

Registradas em 2026-09-13, movidas de `plano_empresa.md` pra este módulo próprio. Não é escopo de nenhum módulo ainda — desenhar a solução agora seria prematuro.

1. **Onde fica a fronteira entre "banco de controle" e "banco por tenant"?** O `api.auth` (usuários/login) também vira um banco por cliente, ou existe um banco central compartilhado que sabe resolver "esse usuário é de qual tenant" antes de rotear pro banco certo?
2. **Como o login descobre qual banco conectar antes mesmo de autenticar?** Por subdomínio (`clienteA.luminihub.com`), por um campo extra no formulário de login, por um diretório central de tenants consultado primeiro?
3. **Como provisionar um banco novo pra um cliente novo?** Quem roda o `AutoMigrate` e os seeders (menu_items, roles, permissions default) em cada banco novo — processo manual, script, endpoint administrativo?
4. **O que muda na infraestrutura de conexão atual?** Hoje `common/database` mantém um único pool de conexão fixo por microsserviço (config vinda do `.env`) — isso teria que virar um pool resolvido em runtime por request, por tenant.
5. **Deploy é uma instância só servindo todos os tenants** (mesmos processos, trocando a conexão de banco em runtime) **ou uma instância de infraestrutura inteira por cliente** (containers/processos separados)?
6. **Como aplicar uma migration nova em N bancos ao mesmo tempo** quando o schema evoluir, sem um ficar pra trás?
