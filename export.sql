--
-- PostgreSQL database dump
--

\restrict 1a5OL3c0Vwt9XLHrbC8iRR7392L91guCXTB1w7DcwtZUClcwUBVMuMSqkjDsOFu

-- Dumped from database version 16.14 (Debian 16.14-1.pgdg13+1)
-- Dumped by pg_dump version 16.14 (Debian 16.14-1.pgdg13+1)

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

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: equipment; Type: TABLE; Schema: public; Owner: codera
--

CREATE TABLE public.equipment (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    typ character varying(20) NOT NULL,
    attack_bonus integer DEFAULT 0 NOT NULL,
    defense_bonus integer DEFAULT 0 NOT NULL,
    speed_bonus integer DEFAULT 0 NOT NULL,
    hp_bonus integer DEFAULT 0 NOT NULL,
    CONSTRAINT equipment_typ_check CHECK (((typ)::text = ANY ((ARRAY['weapon'::character varying, 'armor'::character varying, 'accessory'::character varying])::text[])))
);


ALTER TABLE public.equipment OWNER TO codera;

--
-- Name: equipment_id_seq; Type: SEQUENCE; Schema: public; Owner: codera
--

CREATE SEQUENCE public.equipment_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.equipment_id_seq OWNER TO codera;

--
-- Name: equipment_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: codera
--

ALTER SEQUENCE public.equipment_id_seq OWNED BY public.equipment.id;


--
-- Name: helden; Type: TABLE; Schema: public; Owner: codera
--

CREATE TABLE public.helden (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    rolle character varying(20) NOT NULL,
    max_hp integer NOT NULL,
    current_hp integer NOT NULL,
    attack integer DEFAULT 0 NOT NULL,
    defense integer DEFAULT 0 NOT NULL,
    speed integer DEFAULT 0 NOT NULL,
    equipped_weapon_id integer,
    equipped_armor_id integer,
    equipped_accessory_id integer,
    CONSTRAINT helden_current_hp_check CHECK ((current_hp >= 0)),
    CONSTRAINT helden_max_hp_check CHECK ((max_hp > 0)),
    CONSTRAINT helden_rolle_check CHECK (((rolle)::text = ANY ((ARRAY['arkan'::character varying, 'druide'::character varying, 'kleriker'::character varying, 'krieger'::character varying, 'schmied'::character varying, 'infiltrator'::character varying])::text[])))
);


ALTER TABLE public.helden OWNER TO codera;

--
-- Name: helden_id_seq; Type: SEQUENCE; Schema: public; Owner: codera
--

CREATE SEQUENCE public.helden_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.helden_id_seq OWNER TO codera;

--
-- Name: helden_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: codera
--

ALTER SEQUENCE public.helden_id_seq OWNED BY public.helden.id;


--
-- Name: skills; Type: TABLE; Schema: public; Owner: codera
--

CREATE TABLE public.skills (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    beschreibung text,
    rolle character varying(20) NOT NULL,
    damage_min integer DEFAULT 0 NOT NULL,
    damage_max integer DEFAULT 0 NOT NULL,
    heal integer DEFAULT 0 NOT NULL,
    accuracy numeric(3,2) NOT NULL,
    target_type character varying(20) NOT NULL,
    CONSTRAINT skills_accuracy_check CHECK (((accuracy >= 0.0) AND (accuracy <= 1.0))),
    CONSTRAINT skills_check CHECK ((damage_max >= damage_min)),
    CONSTRAINT skills_damage_min_check CHECK ((damage_min >= 0)),
    CONSTRAINT skills_heal_check CHECK ((heal >= 0)),
    CONSTRAINT skills_rolle_check CHECK (((rolle)::text = ANY ((ARRAY['arkan'::character varying, 'druide'::character varying, 'kleriker'::character varying, 'krieger'::character varying, 'schmied'::character varying, 'infiltrator'::character varying])::text[]))),
    CONSTRAINT skills_target_type_check CHECK (((target_type)::text = ANY ((ARRAY['single_enemy'::character varying, 'all_enemies'::character varying, 'single_ally'::character varying, 'all_allies'::character varying, 'self'::character varying])::text[])))
);


ALTER TABLE public.skills OWNER TO codera;

--
-- Name: skills_id_seq; Type: SEQUENCE; Schema: public; Owner: codera
--

CREATE SEQUENCE public.skills_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.skills_id_seq OWNER TO codera;

--
-- Name: skills_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: codera
--

ALTER SEQUENCE public.skills_id_seq OWNED BY public.skills.id;


--
-- Name: equipment id; Type: DEFAULT; Schema: public; Owner: codera
--

ALTER TABLE ONLY public.equipment ALTER COLUMN id SET DEFAULT nextval('public.equipment_id_seq'::regclass);


--
-- Name: helden id; Type: DEFAULT; Schema: public; Owner: codera
--

ALTER TABLE ONLY public.helden ALTER COLUMN id SET DEFAULT nextval('public.helden_id_seq'::regclass);


--
-- Name: skills id; Type: DEFAULT; Schema: public; Owner: codera
--

ALTER TABLE ONLY public.skills ALTER COLUMN id SET DEFAULT nextval('public.skills_id_seq'::regclass);


--
-- Data for Name: equipment; Type: TABLE DATA; Schema: public; Owner: codera
--

COPY public.equipment (id, name, typ, attack_bonus, defense_bonus, speed_bonus, hp_bonus) FROM stdin;
1	Pergament-Stab	weapon	8	0	0	0
2	Runen-Gewand	armor	0	5	0	0
3	Tintenfass-Amulett	accessory	0	0	3	20
4	Datenstrom-Mantel	armor	0	4	0	0
5	Transformations-Kristall	weapon	6	0	0	0
6	Schema-Ring	accessory	0	0	5	10
7	Architekten-Hammer	weapon	7	0	0	0
8	Runen-Plattenpanzer	armor	0	9	0	0
9	Siegelring-der-Stabilit??t	accessory	0	0	1	25
\.


--
-- Data for Name: helden; Type: TABLE DATA; Schema: public; Owner: codera
--

COPY public.helden (id, name, rolle, max_hp, current_hp, attack, defense, speed, equipped_weapon_id, equipped_armor_id, equipped_accessory_id) FROM stdin;
1	Ron	arkan	120	120	18	8	14	1	2	3
2	Lumjan	druide	100	100	14	10	16	5	4	6
3	Florentin	schmied	130	130	16	16	10	7	8	9
\.


--
-- Data for Name: skills; Type: TABLE DATA; Schema: public; Owner: codera
--

COPY public.skills (id, name, beschreibung, rolle, damage_min, damage_max, heal, accuracy, target_type) FROM stdin;
1	Runen-Geschoss	Magisches Geschoss aus reiner Rune.	arkan	12	24	0	0.90	single_enemy
2	Arkaner Bann	Flaechenschaden auf alle Gegner.	arkan	8	16	0	0.85	all_enemies
3	Kl??rende-Annotation	Heilt einen Verbuendeten.	arkan	0	0	20	1.00	single_ally
4	Datenklinge	Schnitt mit gebuendelten Datenstroemen.	druide	10	20	0	0.85	single_enemy
5	Strukturwandel	Hoher Schaden durch Transformation.	druide	14	28	0	0.70	single_enemy
6	Transformative-Regeneration	Heilt sich selbst.	druide	0	0	16	1.00	self
7	Architekten-Schlag	Solider Hammerschlag.	schmied	14	26	0	0.85	single_enemy
8	Schutz-Rune	Buff: +3 Verteidigung fuer alle.	schmied	0	0	0	1.00	all_allies
9	Konstrukt-Schild	Buff: -50% Schaden fuer einen Ally.	schmied	0	0	0	1.00	single_ally
\.


--
-- Name: equipment_id_seq; Type: SEQUENCE SET; Schema: public; Owner: codera
--

SELECT pg_catalog.setval('public.equipment_id_seq', 9, true);


--
-- Name: helden_id_seq; Type: SEQUENCE SET; Schema: public; Owner: codera
--

SELECT pg_catalog.setval('public.helden_id_seq', 3, true);


--
-- Name: skills_id_seq; Type: SEQUENCE SET; Schema: public; Owner: codera
--

SELECT pg_catalog.setval('public.skills_id_seq', 9, true);


--
-- Name: equipment equipment_name_key; Type: CONSTRAINT; Schema: public; Owner: codera
--

ALTER TABLE ONLY public.equipment
    ADD CONSTRAINT equipment_name_key UNIQUE (name);


--
-- Name: equipment equipment_pkey; Type: CONSTRAINT; Schema: public; Owner: codera
--

ALTER TABLE ONLY public.equipment
    ADD CONSTRAINT equipment_pkey PRIMARY KEY (id);


--
-- Name: helden helden_pkey; Type: CONSTRAINT; Schema: public; Owner: codera
--

ALTER TABLE ONLY public.helden
    ADD CONSTRAINT helden_pkey PRIMARY KEY (id);


--
-- Name: skills skills_name_key; Type: CONSTRAINT; Schema: public; Owner: codera
--

ALTER TABLE ONLY public.skills
    ADD CONSTRAINT skills_name_key UNIQUE (name);


--
-- Name: skills skills_pkey; Type: CONSTRAINT; Schema: public; Owner: codera
--

ALTER TABLE ONLY public.skills
    ADD CONSTRAINT skills_pkey PRIMARY KEY (id);


--
-- Name: helden helden_equipped_accessory_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: codera
--

ALTER TABLE ONLY public.helden
    ADD CONSTRAINT helden_equipped_accessory_id_fkey FOREIGN KEY (equipped_accessory_id) REFERENCES public.equipment(id);


--
-- Name: helden helden_equipped_armor_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: codera
--

ALTER TABLE ONLY public.helden
    ADD CONSTRAINT helden_equipped_armor_id_fkey FOREIGN KEY (equipped_armor_id) REFERENCES public.equipment(id);


--
-- Name: helden helden_equipped_weapon_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: codera
--

ALTER TABLE ONLY public.helden
    ADD CONSTRAINT helden_equipped_weapon_id_fkey FOREIGN KEY (equipped_weapon_id) REFERENCES public.equipment(id);


--
-- PostgreSQL database dump complete
--

\unrestrict 1a5OL3c0Vwt9XLHrbC8iRR7392L91guCXTB1w7DcwtZUClcwUBVMuMSqkjDsOFu

