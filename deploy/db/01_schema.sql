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
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;


--
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


--
-- Name: update_stock_after_purchase(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.update_stock_after_purchase() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
DECLARE
    current_stock INTEGER;
BEGIN
    -- Obter estoque atual
    SELECT current_stock INTO current_stock FROM products WHERE id = NEW.product_id;
    
    -- Atualizar estoque
    UPDATE products 
    SET current_stock = current_stock + NEW.quantity
    WHERE id = NEW.product_id;
    
    -- Registrar movimentação de estoque
    INSERT INTO inventory_movements (
        product_id, 
        quantity, 
        previous_stock, 
        new_stock, 
        movement_type, 
        reference_id, 
        reference_type, 
        created_by
    ) VALUES (
        NEW.product_id,
        NEW.quantity,
        current_stock,
        current_stock + NEW.quantity,
        'entrada',
        NEW.purchase_id,
        'compra',
        (SELECT created_by FROM purchases WHERE id = NEW.purchase_id)
    );
    
    RETURN NEW;
END;
$$;


--
-- Name: update_stock_after_sale(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.update_stock_after_sale() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
DECLARE
    current_stock INTEGER;
BEGIN
    -- Obter estoque atual
    SELECT current_stock INTO current_stock FROM products WHERE id = NEW.product_id;
    
    -- Atualizar estoque
    UPDATE products 
    SET current_stock = current_stock - NEW.quantity
    WHERE id = NEW.product_id;
    
    -- Registrar movimentação de estoque
    INSERT INTO inventory_movements (
        product_id, 
        quantity, 
        previous_stock, 
        new_stock, 
        movement_type, 
        reference_id, 
        reference_type, 
        created_by
    ) VALUES (
        NEW.product_id,
        NEW.quantity * -1,
        current_stock,
        current_stock - NEW.quantity,
        'saida',
        NEW.sale_id,
        'venda',
        (SELECT created_by FROM sales WHERE id = NEW.sale_id)
    );
    
    RETURN NEW;
END;
$$;


--
-- Name: update_updated_at_column(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.update_updated_at_column() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: address; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.address (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    street text NOT NULL,
    number text,
    neighborhood text NOT NULL,
    zip_code character varying(8) NOT NULL,
    city_id bigint NOT NULL,
    customer_id bigint,
    supplier_id bigint
);


--
-- Name: address_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.address_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: address_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.address_id_seq OWNED BY public.address.id;


--
-- Name: city; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.city (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    name text NOT NULL,
    ibge_code text NOT NULL,
    state_id bigint NOT NULL
);


--
-- Name: city_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.city_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: city_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.city_id_seq OWNED BY public.city.id;


--
-- Name: companies; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.companies (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    parent_id bigint,
    legal_name character varying(150) NOT NULL,
    trade_name character varying(150),
    cnpj character varying(20),
    is_active boolean DEFAULT true,
    created_by bigint,
    updated_by bigint,
    logo_file bytea,
    logo_mime_type character varying(50)
);


--
-- Name: companies_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.companies_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: companies_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.companies_id_seq OWNED BY public.companies.id;


--
-- Name: company_fiscal_configs; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.company_fiscal_configs (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    company_id bigint NOT NULL,
    certificate_file bytea,
    certificate_password text,
    certificate_expiry timestamp with time zone,
    tax_regime character varying(20) NOT NULL,
    accountant_name character varying(150),
    accountant_document character varying(20),
    accountant_contact character varying(150),
    certificate_subject_name character varying(200),
    certificate_subject_document character varying(20)
);


--
-- Name: company_fiscal_configs_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.company_fiscal_configs_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: company_fiscal_configs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.company_fiscal_configs_id_seq OWNED BY public.company_fiscal_configs.id;


--
-- Name: contact; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.contact (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    customer_id bigint,
    supplier_id bigint,
    type text NOT NULL,
    contact text NOT NULL,
    name text
);


--
-- Name: contact_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.contact_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: contact_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.contact_id_seq OWNED BY public.contact.id;


--
-- Name: country; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.country (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    name text NOT NULL,
    bacen_code text NOT NULL,
    phone_code text NOT NULL
);


--
-- Name: country_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.country_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: country_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.country_id_seq OWNED BY public.country.id;


--
-- Name: customers; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.customers (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    document_type character varying(20),
    document_number character varying(20),
    email character varying(100),
    phone character varying(20),
    address_street character varying(255),
    address_number character varying(20),
    address_complement character varying(100),
    address_neighborhood character varying(100),
    address_city character varying(100),
    address_state character varying(50),
    address_zipcode character varying(20),
    is_active boolean DEFAULT true,
    created_by bigint,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    first_name character varying(100),
    last_name character varying(100),
    person_type character varying(1) DEFAULT 'F'::character varying,
    company_name character varying(100),
    notes character varying(255),
    updated_by bigint
);


--
-- Name: customers_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.customers_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: customers_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.customers_id_seq OWNED BY public.customers.id;


--
-- Name: document; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.document (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    customer_id bigint,
    supplier_id bigint,
    type text NOT NULL,
    number text NOT NULL,
    validate timestamp with time zone,
    emission_date timestamp with time zone,
    department text,
    state_id bigint
);


--
-- Name: document_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.document_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: document_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.document_id_seq OWNED BY public.document.id;


--
-- Name: financial_transactions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.financial_transactions (
    id integer NOT NULL,
    transaction_type character varying(20) NOT NULL,
    amount numeric(15,2) NOT NULL,
    description text,
    reference_id integer,
    reference_type character varying(20),
    transaction_date timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    payment_method_id integer,
    status character varying(20) NOT NULL,
    due_date timestamp with time zone,
    payment_date timestamp with time zone,
    created_by integer,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: financial_transactions_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.financial_transactions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: financial_transactions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.financial_transactions_id_seq OWNED BY public.financial_transactions.id;


--
-- Name: integration_configs; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.integration_configs (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    key character varying(100) NOT NULL,
    value text,
    description character varying(255)
);


--
-- Name: integration_configs_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.integration_configs_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: integration_configs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.integration_configs_id_seq OWNED BY public.integration_configs.id;


--
-- Name: inventory_movements; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.inventory_movements (
    id integer NOT NULL,
    product_id integer,
    quantity integer NOT NULL,
    previous_stock integer NOT NULL,
    new_stock integer NOT NULL,
    movement_type character varying(20) NOT NULL,
    reference_id integer,
    reference_type character varying(20),
    notes text,
    created_by integer,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: inventory_movements_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.inventory_movements_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: inventory_movements_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.inventory_movements_id_seq OWNED BY public.inventory_movements.id;


--
-- Name: measurement_units; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.measurement_units (
    id integer NOT NULL,
    name character varying(50) NOT NULL,
    abbreviation character varying(10) NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: measurement_units_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.measurement_units_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: measurement_units_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.measurement_units_id_seq OWNED BY public.measurement_units.id;


--
-- Name: menu_items; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.menu_items (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    name character varying(100) NOT NULL,
    icon character varying(100),
    href character varying(255),
    parent_id bigint,
    permission_id bigint,
    "position" bigint DEFAULT 0,
    is_active boolean DEFAULT true
);


--
-- Name: menu_items_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.menu_items_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: menu_items_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.menu_items_id_seq OWNED BY public.menu_items.id;


--
-- Name: payment_methods; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.payment_methods (
    id integer NOT NULL,
    name character varying(50) NOT NULL,
    description text,
    is_active boolean DEFAULT true,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: payment_methods_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.payment_methods_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: payment_methods_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.payment_methods_id_seq OWNED BY public.payment_methods.id;


--
-- Name: permissions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.permissions (
    id integer NOT NULL,
    permission character varying(100) NOT NULL,
    description text,
    module character varying(50) NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


--
-- Name: permissions_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.permissions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: permissions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.permissions_id_seq OWNED BY public.permissions.id;


--
-- Name: product_categories; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.product_categories (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    description text,
    parent_id integer,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: product_categories_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.product_categories_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: product_categories_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.product_categories_id_seq OWNED BY public.product_categories.id;


--
-- Name: product_mappings; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.product_mappings (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    cod_pro character varying(30) NOT NULL,
    li_product_id character varying(50),
    li_sku character varying(50),
    last_synced_at timestamp with time zone,
    sync_status character varying(20)
);


--
-- Name: product_mappings_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.product_mappings_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: product_mappings_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.product_mappings_id_seq OWNED BY public.product_mappings.id;


--
-- Name: products; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.products (
    id integer NOT NULL,
    sku character varying(50),
    barcode character varying(50),
    name character varying(255) NOT NULL,
    description text,
    category_id integer,
    unit_id integer,
    cost_price numeric(15,2) NOT NULL,
    selling_price numeric(15,2) NOT NULL,
    min_stock integer DEFAULT 0,
    max_stock integer,
    current_stock integer DEFAULT 0,
    is_active boolean DEFAULT true,
    created_by integer,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: products_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.products_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: products_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.products_id_seq OWNED BY public.products.id;


--
-- Name: purchase_items; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.purchase_items (
    id integer NOT NULL,
    purchase_id integer,
    product_id integer,
    quantity integer NOT NULL,
    unit_price numeric(15,2) NOT NULL,
    total_amount numeric(15,2) NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: purchase_items_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.purchase_items_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: purchase_items_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.purchase_items_id_seq OWNED BY public.purchase_items.id;


--
-- Name: purchases; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.purchases (
    id integer NOT NULL,
    supplier_id integer,
    purchase_date timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    total_amount numeric(15,2) NOT NULL,
    status character varying(20) NOT NULL,
    notes text,
    created_by integer,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: purchases_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.purchases_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: purchases_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.purchases_id_seq OWNED BY public.purchases.id;


--
-- Name: role_permissions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.role_permissions (
    role_id integer NOT NULL,
    permission_id integer NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: roles; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.roles (
    id integer NOT NULL,
    name character varying(50) NOT NULL,
    description text,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


--
-- Name: roles_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.roles_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: roles_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.roles_id_seq OWNED BY public.roles.id;


--
-- Name: sale_items; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.sale_items (
    id integer NOT NULL,
    sale_id integer,
    product_id integer,
    quantity integer NOT NULL,
    unit_price numeric(15,2) NOT NULL,
    discount_percent numeric(5,2) DEFAULT 0,
    discount_amount numeric(15,2) DEFAULT 0,
    tax_percent numeric(5,2) DEFAULT 0,
    tax_amount numeric(15,2) DEFAULT 0,
    total_amount numeric(15,2) NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: sale_items_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.sale_items_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: sale_items_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.sale_items_id_seq OWNED BY public.sale_items.id;


--
-- Name: sales; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.sales (
    id integer NOT NULL,
    customer_id integer,
    sale_date timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    total_amount numeric(15,2) NOT NULL,
    discount_amount numeric(15,2) DEFAULT 0,
    tax_amount numeric(15,2) DEFAULT 0,
    final_amount numeric(15,2) NOT NULL,
    payment_method_id integer,
    status character varying(20) NOT NULL,
    notes text,
    created_by integer,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: sales_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.sales_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: sales_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.sales_id_seq OWNED BY public.sales.id;


--
-- Name: state; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.state (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    name text NOT NULL,
    uf character varying(2) NOT NULL,
    ibge_code text NOT NULL,
    country_id bigint NOT NULL
);


--
-- Name: state_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.state_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: state_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.state_id_seq OWNED BY public.state.id;


--
-- Name: suppliers; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.suppliers (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    document_type character varying(20),
    document_number character varying(20),
    contact_name character varying(100),
    email character varying(100),
    phone character varying(20),
    address_street character varying(255),
    address_number character varying(20),
    address_complement character varying(100),
    address_neighborhood character varying(100),
    address_city character varying(100),
    address_state character varying(50),
    address_zipcode character varying(20),
    is_active boolean DEFAULT true,
    created_by bigint,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    first_name character varying(100),
    last_name character varying(100),
    person_type character varying(1) DEFAULT 'F'::character varying,
    company_name character varying(100),
    notes character varying(255),
    updated_by bigint
);


--
-- Name: suppliers_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.suppliers_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: suppliers_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.suppliers_id_seq OWNED BY public.suppliers.id;


--
-- Name: sync_logs; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.sync_logs (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    direction character varying(20) NOT NULL,
    entity_type character varying(50) NOT NULL,
    reference_id character varying(100),
    status character varying(20) NOT NULL,
    message character varying(500),
    details jsonb
);


--
-- Name: sync_logs_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.sync_logs_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: sync_logs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.sync_logs_id_seq OWNED BY public.sync_logs.id;


--
-- Name: system_logs; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.system_logs (
    id integer NOT NULL,
    user_id bigint,
    action character varying(100) NOT NULL,
    entity_type character varying(50),
    entity_id text,
    details jsonb,
    ip_address character varying(45),
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


--
-- Name: system_logs_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.system_logs_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: system_logs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.system_logs_id_seq OWNED BY public.system_logs.id;


--
-- Name: user_permissions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.user_permissions (
    user_id integer NOT NULL,
    permission_id integer NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.users (
    id integer NOT NULL,
    username character varying(50) NOT NULL,
    password_hash character varying(255) NOT NULL,
    name character varying(100) NOT NULL,
    email character varying(100),
    phone character varying(20),
    role_id bigint,
    is_active boolean DEFAULT true,
    last_login timestamp with time zone,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    company_id bigint,
    active_company_id bigint
);


--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: webhook_events; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.webhook_events (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    source character varying(50) NOT NULL,
    event_type character varying(50),
    headers jsonb,
    payload jsonb,
    status character varying(20) NOT NULL,
    received_at timestamp with time zone,
    processed_at timestamp with time zone,
    error_message character varying(500)
);


--
-- Name: webhook_events_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.webhook_events_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: webhook_events_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.webhook_events_id_seq OWNED BY public.webhook_events.id;


--
-- Name: address id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.address ALTER COLUMN id SET DEFAULT nextval('public.address_id_seq'::regclass);


--
-- Name: city id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.city ALTER COLUMN id SET DEFAULT nextval('public.city_id_seq'::regclass);


--
-- Name: companies id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.companies ALTER COLUMN id SET DEFAULT nextval('public.companies_id_seq'::regclass);


--
-- Name: company_fiscal_configs id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.company_fiscal_configs ALTER COLUMN id SET DEFAULT nextval('public.company_fiscal_configs_id_seq'::regclass);


--
-- Name: contact id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.contact ALTER COLUMN id SET DEFAULT nextval('public.contact_id_seq'::regclass);


--
-- Name: country id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.country ALTER COLUMN id SET DEFAULT nextval('public.country_id_seq'::regclass);


--
-- Name: customers id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.customers ALTER COLUMN id SET DEFAULT nextval('public.customers_id_seq'::regclass);


--
-- Name: document id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.document ALTER COLUMN id SET DEFAULT nextval('public.document_id_seq'::regclass);


--
-- Name: financial_transactions id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.financial_transactions ALTER COLUMN id SET DEFAULT nextval('public.financial_transactions_id_seq'::regclass);


--
-- Name: integration_configs id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.integration_configs ALTER COLUMN id SET DEFAULT nextval('public.integration_configs_id_seq'::regclass);


--
-- Name: inventory_movements id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.inventory_movements ALTER COLUMN id SET DEFAULT nextval('public.inventory_movements_id_seq'::regclass);


--
-- Name: measurement_units id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.measurement_units ALTER COLUMN id SET DEFAULT nextval('public.measurement_units_id_seq'::regclass);


--
-- Name: menu_items id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.menu_items ALTER COLUMN id SET DEFAULT nextval('public.menu_items_id_seq'::regclass);


--
-- Name: payment_methods id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.payment_methods ALTER COLUMN id SET DEFAULT nextval('public.payment_methods_id_seq'::regclass);


--
-- Name: permissions id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.permissions ALTER COLUMN id SET DEFAULT nextval('public.permissions_id_seq'::regclass);


--
-- Name: product_categories id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_categories ALTER COLUMN id SET DEFAULT nextval('public.product_categories_id_seq'::regclass);


--
-- Name: product_mappings id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_mappings ALTER COLUMN id SET DEFAULT nextval('public.product_mappings_id_seq'::regclass);


--
-- Name: products id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.products ALTER COLUMN id SET DEFAULT nextval('public.products_id_seq'::regclass);


--
-- Name: purchase_items id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.purchase_items ALTER COLUMN id SET DEFAULT nextval('public.purchase_items_id_seq'::regclass);


--
-- Name: purchases id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.purchases ALTER COLUMN id SET DEFAULT nextval('public.purchases_id_seq'::regclass);


--
-- Name: roles id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.roles ALTER COLUMN id SET DEFAULT nextval('public.roles_id_seq'::regclass);


--
-- Name: sale_items id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sale_items ALTER COLUMN id SET DEFAULT nextval('public.sale_items_id_seq'::regclass);


--
-- Name: sales id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales ALTER COLUMN id SET DEFAULT nextval('public.sales_id_seq'::regclass);


--
-- Name: state id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.state ALTER COLUMN id SET DEFAULT nextval('public.state_id_seq'::regclass);


--
-- Name: suppliers id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.suppliers ALTER COLUMN id SET DEFAULT nextval('public.suppliers_id_seq'::regclass);


--
-- Name: sync_logs id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sync_logs ALTER COLUMN id SET DEFAULT nextval('public.sync_logs_id_seq'::regclass);


--
-- Name: system_logs id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.system_logs ALTER COLUMN id SET DEFAULT nextval('public.system_logs_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Name: webhook_events id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.webhook_events ALTER COLUMN id SET DEFAULT nextval('public.webhook_events_id_seq'::regclass);


--
-- Name: address address_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.address
    ADD CONSTRAINT address_pkey PRIMARY KEY (id);


--
-- Name: city city_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.city
    ADD CONSTRAINT city_pkey PRIMARY KEY (id);


--
-- Name: companies companies_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.companies
    ADD CONSTRAINT companies_pkey PRIMARY KEY (id);


--
-- Name: company_fiscal_configs company_fiscal_configs_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.company_fiscal_configs
    ADD CONSTRAINT company_fiscal_configs_pkey PRIMARY KEY (id);


--
-- Name: contact contact_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.contact
    ADD CONSTRAINT contact_pkey PRIMARY KEY (id);


--
-- Name: country country_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.country
    ADD CONSTRAINT country_pkey PRIMARY KEY (id);


--
-- Name: customers customers_document_number_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.customers
    ADD CONSTRAINT customers_document_number_key UNIQUE (document_number);


--
-- Name: customers customers_email_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.customers
    ADD CONSTRAINT customers_email_key UNIQUE (email);


--
-- Name: customers customers_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.customers
    ADD CONSTRAINT customers_pkey PRIMARY KEY (id);


--
-- Name: document document_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.document
    ADD CONSTRAINT document_pkey PRIMARY KEY (id);


--
-- Name: financial_transactions financial_transactions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.financial_transactions
    ADD CONSTRAINT financial_transactions_pkey PRIMARY KEY (id);


--
-- Name: integration_configs integration_configs_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.integration_configs
    ADD CONSTRAINT integration_configs_pkey PRIMARY KEY (id);


--
-- Name: inventory_movements inventory_movements_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.inventory_movements
    ADD CONSTRAINT inventory_movements_pkey PRIMARY KEY (id);


--
-- Name: measurement_units measurement_units_abbreviation_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.measurement_units
    ADD CONSTRAINT measurement_units_abbreviation_key UNIQUE (abbreviation);


--
-- Name: measurement_units measurement_units_name_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.measurement_units
    ADD CONSTRAINT measurement_units_name_key UNIQUE (name);


--
-- Name: measurement_units measurement_units_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.measurement_units
    ADD CONSTRAINT measurement_units_pkey PRIMARY KEY (id);


--
-- Name: menu_items menu_items_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.menu_items
    ADD CONSTRAINT menu_items_pkey PRIMARY KEY (id);


--
-- Name: payment_methods payment_methods_name_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.payment_methods
    ADD CONSTRAINT payment_methods_name_key UNIQUE (name);


--
-- Name: payment_methods payment_methods_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.payment_methods
    ADD CONSTRAINT payment_methods_pkey PRIMARY KEY (id);


--
-- Name: permissions permissions_name_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.permissions
    ADD CONSTRAINT permissions_name_key UNIQUE (permission);


--
-- Name: permissions permissions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.permissions
    ADD CONSTRAINT permissions_pkey PRIMARY KEY (id);


--
-- Name: product_categories product_categories_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_categories
    ADD CONSTRAINT product_categories_pkey PRIMARY KEY (id);


--
-- Name: product_mappings product_mappings_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_mappings
    ADD CONSTRAINT product_mappings_pkey PRIMARY KEY (id);


--
-- Name: products products_barcode_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_barcode_key UNIQUE (barcode);


--
-- Name: products products_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_pkey PRIMARY KEY (id);


--
-- Name: products products_sku_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_sku_key UNIQUE (sku);


--
-- Name: purchase_items purchase_items_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.purchase_items
    ADD CONSTRAINT purchase_items_pkey PRIMARY KEY (id);


--
-- Name: purchases purchases_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.purchases
    ADD CONSTRAINT purchases_pkey PRIMARY KEY (id);


--
-- Name: role_permissions role_permissions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.role_permissions
    ADD CONSTRAINT role_permissions_pkey PRIMARY KEY (role_id, permission_id);


--
-- Name: roles roles_name_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_name_key UNIQUE (name);


--
-- Name: roles roles_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_pkey PRIMARY KEY (id);


--
-- Name: sale_items sale_items_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sale_items
    ADD CONSTRAINT sale_items_pkey PRIMARY KEY (id);


--
-- Name: sales sales_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales
    ADD CONSTRAINT sales_pkey PRIMARY KEY (id);


--
-- Name: state state_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.state
    ADD CONSTRAINT state_pkey PRIMARY KEY (id);


--
-- Name: suppliers suppliers_document_number_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.suppliers
    ADD CONSTRAINT suppliers_document_number_key UNIQUE (document_number);


--
-- Name: suppliers suppliers_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.suppliers
    ADD CONSTRAINT suppliers_pkey PRIMARY KEY (id);


--
-- Name: sync_logs sync_logs_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sync_logs
    ADD CONSTRAINT sync_logs_pkey PRIMARY KEY (id);


--
-- Name: system_logs system_logs_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.system_logs
    ADD CONSTRAINT system_logs_pkey PRIMARY KEY (id);


--
-- Name: city uni_city_ibge_code; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.city
    ADD CONSTRAINT uni_city_ibge_code UNIQUE (ibge_code);


--
-- Name: contact uni_contact_contact; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.contact
    ADD CONSTRAINT uni_contact_contact UNIQUE (contact);


--
-- Name: country uni_country_bacen_code; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.country
    ADD CONSTRAINT uni_country_bacen_code UNIQUE (bacen_code);


--
-- Name: country uni_country_name; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.country
    ADD CONSTRAINT uni_country_name UNIQUE (name);


--
-- Name: document uni_document_number; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.document
    ADD CONSTRAINT uni_document_number UNIQUE (number);


--
-- Name: integration_configs uni_integration_configs_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.integration_configs
    ADD CONSTRAINT uni_integration_configs_key UNIQUE (key);


--
-- Name: product_mappings uni_product_mappings_cod_pro; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_mappings
    ADD CONSTRAINT uni_product_mappings_cod_pro UNIQUE (cod_pro);


--
-- Name: state uni_state_ibge_code; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.state
    ADD CONSTRAINT uni_state_ibge_code UNIQUE (ibge_code);


--
-- Name: state uni_state_name; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.state
    ADD CONSTRAINT uni_state_name UNIQUE (name);


--
-- Name: state uni_state_uf; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.state
    ADD CONSTRAINT uni_state_uf UNIQUE (uf);


--
-- Name: user_permissions user_permissions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_permissions
    ADD CONSTRAINT user_permissions_pkey PRIMARY KEY (user_id, permission_id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: users users_username_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_username_key UNIQUE (username);


--
-- Name: webhook_events webhook_events_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.webhook_events
    ADD CONSTRAINT webhook_events_pkey PRIMARY KEY (id);


--
-- Name: idx_address_customer_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_address_customer_id ON public.address USING btree (customer_id);


--
-- Name: idx_address_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_address_deleted_at ON public.address USING btree (deleted_at);


--
-- Name: idx_address_supplier_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_address_supplier_id ON public.address USING btree (supplier_id);


--
-- Name: idx_city_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_city_deleted_at ON public.city USING btree (deleted_at);


--
-- Name: idx_companies_cnpj; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_companies_cnpj ON public.companies USING btree (cnpj) WHERE ((deleted_at IS NULL) AND (cnpj IS NOT NULL));


--
-- Name: idx_companies_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_companies_deleted_at ON public.companies USING btree (deleted_at);


--
-- Name: idx_company_fiscal_configs_company_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_company_fiscal_configs_company_id ON public.company_fiscal_configs USING btree (company_id) WHERE (deleted_at IS NULL);


--
-- Name: idx_company_fiscal_configs_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_company_fiscal_configs_deleted_at ON public.company_fiscal_configs USING btree (deleted_at);


--
-- Name: idx_contact_customer_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_contact_customer_id ON public.contact USING btree (customer_id);


--
-- Name: idx_contact_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_contact_deleted_at ON public.contact USING btree (deleted_at);


--
-- Name: idx_contact_supplier_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_contact_supplier_id ON public.contact USING btree (supplier_id);


--
-- Name: idx_country_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_country_deleted_at ON public.country USING btree (deleted_at);


--
-- Name: idx_customers_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_customers_deleted_at ON public.customers USING btree (deleted_at);


--
-- Name: idx_document_customer_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_document_customer_id ON public.document USING btree (customer_id);


--
-- Name: idx_document_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_document_deleted_at ON public.document USING btree (deleted_at);


--
-- Name: idx_document_supplier_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_document_supplier_id ON public.document USING btree (supplier_id);


--
-- Name: idx_financial_transactions_date; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_financial_transactions_date ON public.financial_transactions USING btree (transaction_date);


--
-- Name: idx_integration_configs_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_integration_configs_deleted_at ON public.integration_configs USING btree (deleted_at);


--
-- Name: idx_inventory_movements_date; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_inventory_movements_date ON public.inventory_movements USING btree (created_at);


--
-- Name: idx_inventory_movements_product; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_inventory_movements_product ON public.inventory_movements USING btree (product_id);


--
-- Name: idx_menu_items_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_menu_items_deleted_at ON public.menu_items USING btree (deleted_at);


--
-- Name: idx_permissions_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_permissions_deleted_at ON public.permissions USING btree (deleted_at);


--
-- Name: idx_product_mappings_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_product_mappings_deleted_at ON public.product_mappings USING btree (deleted_at);


--
-- Name: idx_products_category; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_products_category ON public.products USING btree (category_id);


--
-- Name: idx_products_name; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_products_name ON public.products USING btree (name);


--
-- Name: idx_purchase_items_purchase; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_purchase_items_purchase ON public.purchase_items USING btree (purchase_id);


--
-- Name: idx_purchases_supplier; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_purchases_supplier ON public.purchases USING btree (supplier_id);


--
-- Name: idx_roles_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_roles_deleted_at ON public.roles USING btree (deleted_at);


--
-- Name: idx_sale_items_product; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_sale_items_product ON public.sale_items USING btree (product_id);


--
-- Name: idx_sale_items_sale; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_sale_items_sale ON public.sale_items USING btree (sale_id);


--
-- Name: idx_sales_customer; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_sales_customer ON public.sales USING btree (customer_id);


--
-- Name: idx_sales_date; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_sales_date ON public.sales USING btree (sale_date);


--
-- Name: idx_state_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_state_deleted_at ON public.state USING btree (deleted_at);


--
-- Name: idx_suppliers_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_suppliers_deleted_at ON public.suppliers USING btree (deleted_at);


--
-- Name: idx_sync_logs_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_sync_logs_deleted_at ON public.sync_logs USING btree (deleted_at);


--
-- Name: idx_system_logs_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_system_logs_deleted_at ON public.system_logs USING btree (deleted_at);


--
-- Name: idx_users_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_users_deleted_at ON public.users USING btree (deleted_at);


--
-- Name: idx_webhook_events_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_webhook_events_deleted_at ON public.webhook_events USING btree (deleted_at);


--
-- Name: purchase_items trigger_update_stock_after_purchase; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER trigger_update_stock_after_purchase AFTER INSERT ON public.purchase_items FOR EACH ROW EXECUTE FUNCTION public.update_stock_after_purchase();


--
-- Name: sale_items trigger_update_stock_after_sale; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER trigger_update_stock_after_sale AFTER INSERT ON public.sale_items FOR EACH ROW EXECUTE FUNCTION public.update_stock_after_sale();


--
-- Name: customers update_customers_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER update_customers_updated_at BEFORE UPDATE ON public.customers FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: financial_transactions update_financial_transactions_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER update_financial_transactions_updated_at BEFORE UPDATE ON public.financial_transactions FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: products update_products_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER update_products_updated_at BEFORE UPDATE ON public.products FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: purchases update_purchases_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER update_purchases_updated_at BEFORE UPDATE ON public.purchases FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: sales update_sales_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER update_sales_updated_at BEFORE UPDATE ON public.sales FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: suppliers update_suppliers_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER update_suppliers_updated_at BEFORE UPDATE ON public.suppliers FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: users update_users_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON public.users FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: customers customers_created_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.customers
    ADD CONSTRAINT customers_created_by_fkey FOREIGN KEY (created_by) REFERENCES public.users(id);


--
-- Name: financial_transactions financial_transactions_created_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.financial_transactions
    ADD CONSTRAINT financial_transactions_created_by_fkey FOREIGN KEY (created_by) REFERENCES public.users(id);


--
-- Name: financial_transactions financial_transactions_payment_method_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.financial_transactions
    ADD CONSTRAINT financial_transactions_payment_method_id_fkey FOREIGN KEY (payment_method_id) REFERENCES public.payment_methods(id);


--
-- Name: address fk_address_city; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.address
    ADD CONSTRAINT fk_address_city FOREIGN KEY (city_id) REFERENCES public.city(id);


--
-- Name: city fk_city_state; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.city
    ADD CONSTRAINT fk_city_state FOREIGN KEY (state_id) REFERENCES public.state(id);


--
-- Name: companies fk_companies_children; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.companies
    ADD CONSTRAINT fk_companies_children FOREIGN KEY (parent_id) REFERENCES public.companies(id);


--
-- Name: address fk_customers_addresses; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.address
    ADD CONSTRAINT fk_customers_addresses FOREIGN KEY (customer_id) REFERENCES public.customers(id);


--
-- Name: contact fk_customers_contacts; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.contact
    ADD CONSTRAINT fk_customers_contacts FOREIGN KEY (customer_id) REFERENCES public.customers(id);


--
-- Name: customers fk_customers_created_by; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.customers
    ADD CONSTRAINT fk_customers_created_by FOREIGN KEY (created_by) REFERENCES public.users(id);


--
-- Name: document fk_customers_documents; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.document
    ADD CONSTRAINT fk_customers_documents FOREIGN KEY (customer_id) REFERENCES public.customers(id);


--
-- Name: customers fk_customers_updated_by; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.customers
    ADD CONSTRAINT fk_customers_updated_by FOREIGN KEY (updated_by) REFERENCES public.users(id);


--
-- Name: document fk_document_state; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.document
    ADD CONSTRAINT fk_document_state FOREIGN KEY (state_id) REFERENCES public.state(id);


--
-- Name: menu_items fk_menu_items_children; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.menu_items
    ADD CONSTRAINT fk_menu_items_children FOREIGN KEY (parent_id) REFERENCES public.menu_items(id);


--
-- Name: menu_items fk_menu_items_permission; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.menu_items
    ADD CONSTRAINT fk_menu_items_permission FOREIGN KEY (permission_id) REFERENCES public.permissions(id);


--
-- Name: role_permissions fk_role_permissions_permission; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.role_permissions
    ADD CONSTRAINT fk_role_permissions_permission FOREIGN KEY (permission_id) REFERENCES public.permissions(id);


--
-- Name: role_permissions fk_role_permissions_role; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.role_permissions
    ADD CONSTRAINT fk_role_permissions_role FOREIGN KEY (role_id) REFERENCES public.roles(id);


--
-- Name: users fk_roles_users; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT fk_roles_users FOREIGN KEY (role_id) REFERENCES public.roles(id);


--
-- Name: state fk_state_country; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.state
    ADD CONSTRAINT fk_state_country FOREIGN KEY (country_id) REFERENCES public.country(id);


--
-- Name: address fk_suppliers_addresses; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.address
    ADD CONSTRAINT fk_suppliers_addresses FOREIGN KEY (supplier_id) REFERENCES public.suppliers(id);


--
-- Name: contact fk_suppliers_contacts; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.contact
    ADD CONSTRAINT fk_suppliers_contacts FOREIGN KEY (supplier_id) REFERENCES public.suppliers(id);


--
-- Name: suppliers fk_suppliers_created_by; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.suppliers
    ADD CONSTRAINT fk_suppliers_created_by FOREIGN KEY (created_by) REFERENCES public.users(id);


--
-- Name: document fk_suppliers_documents; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.document
    ADD CONSTRAINT fk_suppliers_documents FOREIGN KEY (supplier_id) REFERENCES public.suppliers(id);


--
-- Name: suppliers fk_suppliers_updated_by; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.suppliers
    ADD CONSTRAINT fk_suppliers_updated_by FOREIGN KEY (updated_by) REFERENCES public.users(id);


--
-- Name: system_logs fk_system_logs_user; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.system_logs
    ADD CONSTRAINT fk_system_logs_user FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: users fk_users_active_company; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT fk_users_active_company FOREIGN KEY (active_company_id) REFERENCES public.companies(id);


--
-- Name: users fk_users_company; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT fk_users_company FOREIGN KEY (company_id) REFERENCES public.companies(id);


--
-- Name: inventory_movements inventory_movements_created_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.inventory_movements
    ADD CONSTRAINT inventory_movements_created_by_fkey FOREIGN KEY (created_by) REFERENCES public.users(id);


--
-- Name: inventory_movements inventory_movements_product_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.inventory_movements
    ADD CONSTRAINT inventory_movements_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.products(id);


--
-- Name: product_categories product_categories_parent_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_categories
    ADD CONSTRAINT product_categories_parent_id_fkey FOREIGN KEY (parent_id) REFERENCES public.product_categories(id);


--
-- Name: products products_category_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.product_categories(id);


--
-- Name: products products_created_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_created_by_fkey FOREIGN KEY (created_by) REFERENCES public.users(id);


--
-- Name: products products_unit_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_unit_id_fkey FOREIGN KEY (unit_id) REFERENCES public.measurement_units(id);


--
-- Name: purchase_items purchase_items_product_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.purchase_items
    ADD CONSTRAINT purchase_items_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.products(id);


--
-- Name: purchase_items purchase_items_purchase_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.purchase_items
    ADD CONSTRAINT purchase_items_purchase_id_fkey FOREIGN KEY (purchase_id) REFERENCES public.purchases(id) ON DELETE CASCADE;


--
-- Name: purchases purchases_created_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.purchases
    ADD CONSTRAINT purchases_created_by_fkey FOREIGN KEY (created_by) REFERENCES public.users(id);


--
-- Name: purchases purchases_supplier_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.purchases
    ADD CONSTRAINT purchases_supplier_id_fkey FOREIGN KEY (supplier_id) REFERENCES public.suppliers(id);


--
-- Name: role_permissions role_permissions_permission_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.role_permissions
    ADD CONSTRAINT role_permissions_permission_id_fkey FOREIGN KEY (permission_id) REFERENCES public.permissions(id) ON DELETE CASCADE;


--
-- Name: role_permissions role_permissions_role_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.role_permissions
    ADD CONSTRAINT role_permissions_role_id_fkey FOREIGN KEY (role_id) REFERENCES public.roles(id) ON DELETE CASCADE;


--
-- Name: sale_items sale_items_product_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sale_items
    ADD CONSTRAINT sale_items_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.products(id);


--
-- Name: sale_items sale_items_sale_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sale_items
    ADD CONSTRAINT sale_items_sale_id_fkey FOREIGN KEY (sale_id) REFERENCES public.sales(id) ON DELETE CASCADE;


--
-- Name: sales sales_created_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales
    ADD CONSTRAINT sales_created_by_fkey FOREIGN KEY (created_by) REFERENCES public.users(id);


--
-- Name: sales sales_customer_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales
    ADD CONSTRAINT sales_customer_id_fkey FOREIGN KEY (customer_id) REFERENCES public.customers(id);


--
-- Name: sales sales_payment_method_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sales
    ADD CONSTRAINT sales_payment_method_id_fkey FOREIGN KEY (payment_method_id) REFERENCES public.payment_methods(id);


--
-- Name: suppliers suppliers_created_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.suppliers
    ADD CONSTRAINT suppliers_created_by_fkey FOREIGN KEY (created_by) REFERENCES public.users(id);


--
-- Name: system_logs system_logs_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.system_logs
    ADD CONSTRAINT system_logs_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: user_permissions user_permissions_permission_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_permissions
    ADD CONSTRAINT user_permissions_permission_id_fkey FOREIGN KEY (permission_id) REFERENCES public.permissions(id) ON DELETE CASCADE;


--
-- Name: user_permissions user_permissions_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_permissions
    ADD CONSTRAINT user_permissions_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: users users_role_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_role_id_fkey FOREIGN KEY (role_id) REFERENCES public.roles(id);


--
-- PostgreSQL database dump complete
--

