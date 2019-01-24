--
-- PostgreSQL database dump
--

-- Dumped from database version 10.5
-- Dumped by pg_dump version 10.5

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: prod; Type: SCHEMA; Schema: -; Owner: postgres
--

CREATE SCHEMA prod;


ALTER SCHEMA prod OWNER TO postgres;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: attribute; Type: TABLE; Schema: prod; Owner: postgres
--

CREATE TABLE prod.attribute (
                              id integer NOT NULL,
                              name text NOT NULL
);


ALTER TABLE prod.attribute OWNER TO postgres;

--
-- Name: attribute_id_seq; Type: SEQUENCE; Schema: prod; Owner: postgres
--

CREATE SEQUENCE prod.attribute_id_seq
  AS integer
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE prod.attribute_id_seq OWNER TO postgres;

--
-- Name: attribute_id_seq; Type: SEQUENCE OWNED BY; Schema: prod; Owner: postgres
--

ALTER SEQUENCE prod.attribute_id_seq OWNED BY prod.attribute.id;


--
-- Name: job; Type: TABLE; Schema: prod; Owner: postgres
--

CREATE TABLE prod.job (
                        id integer NOT NULL,
                        product text NOT NULL,
                        version text NOT NULL,
                        name text NOT NULL,
                        measurement text NOT NULL,
                        "timestamp" timestamp without time zone NOT NULL,
                        value integer NOT NULL,
                        created_on timestamp without time zone NOT NULL,
                        created_by text NOT NULL
);


ALTER TABLE prod.job OWNER TO postgres;

--
-- Name: job_attribute; Type: TABLE; Schema: prod; Owner: postgres
--

CREATE TABLE prod.job_attribute (
                                  id integer NOT NULL,
                                  job_id integer NOT NULL,
                                  attribute_id integer NOT NULL
);


ALTER TABLE prod.job_attribute OWNER TO postgres;

--
-- Name: job_attribute_id_seq; Type: SEQUENCE; Schema: prod; Owner: postgres
--

CREATE SEQUENCE prod.job_attribute_id_seq
  AS integer
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE prod.job_attribute_id_seq OWNER TO postgres;

--
-- Name: job_attribute_id_seq; Type: SEQUENCE OWNED BY; Schema: prod; Owner: postgres
--

ALTER SEQUENCE prod.job_attribute_id_seq OWNED BY prod.job_attribute.id;


--
-- Name: job_id_seq; Type: SEQUENCE; Schema: prod; Owner: postgres
--

CREATE SEQUENCE prod.job_id_seq
  AS integer
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE prod.job_id_seq OWNER TO postgres;

--
-- Name: job_id_seq; Type: SEQUENCE OWNED BY; Schema: prod; Owner: postgres
--

ALTER SEQUENCE prod.job_id_seq OWNED BY prod.job.id;


--
-- Name: attribute id; Type: DEFAULT; Schema: prod; Owner: postgres
--

ALTER TABLE ONLY prod.attribute ALTER COLUMN id SET DEFAULT nextval('prod.attribute_id_seq'::regclass);


--
-- Name: job id; Type: DEFAULT; Schema: prod; Owner: postgres
--

ALTER TABLE ONLY prod.job ALTER COLUMN id SET DEFAULT nextval('prod.job_id_seq'::regclass);


--
-- Name: job_attribute id; Type: DEFAULT; Schema: prod; Owner: postgres
--

ALTER TABLE ONLY prod.job_attribute ALTER COLUMN id SET DEFAULT nextval('prod.job_attribute_id_seq'::regclass);


--
-- Data for Name: attribute; Type: TABLE DATA; Schema: prod; Owner: postgres
--

COPY prod.attribute (id, name) FROM stdin;
\.


--
-- Data for Name: job; Type: TABLE DATA; Schema: prod; Owner: postgres
--

COPY prod.job (id, product, version, name, measurement, "timestamp", value, created_on, created_by) FROM stdin;
\.


--
-- Data for Name: job_attribute; Type: TABLE DATA; Schema: prod; Owner: postgres
--

COPY prod.job_attribute (id, job_id, attribute_id) FROM stdin;
\.


--
-- Name: attribute_id_seq; Type: SEQUENCE SET; Schema: prod; Owner: postgres
--

SELECT pg_catalog.setval('prod.attribute_id_seq', 1, false);


--
-- Name: job_attribute_id_seq; Type: SEQUENCE SET; Schema: prod; Owner: postgres
--

SELECT pg_catalog.setval('prod.job_attribute_id_seq', 1, false);


--
-- Name: job_id_seq; Type: SEQUENCE SET; Schema: prod; Owner: postgres
--

SELECT pg_catalog.setval('prod.job_id_seq', 1, false);


--
-- Name: attribute attribute_pk; Type: CONSTRAINT; Schema: prod; Owner: postgres
--

ALTER TABLE ONLY prod.attribute
  ADD CONSTRAINT attribute_pk PRIMARY KEY (id);


--
-- Name: job_attribute job_attribute_pk; Type: CONSTRAINT; Schema: prod; Owner: postgres
--

ALTER TABLE ONLY prod.job_attribute
  ADD CONSTRAINT job_attribute_pk PRIMARY KEY (id);


--
-- Name: job job_pk; Type: CONSTRAINT; Schema: prod; Owner: postgres
--

ALTER TABLE ONLY prod.job
  ADD CONSTRAINT job_pk PRIMARY KEY (id);


--
-- Name: attribute_name_uindex; Type: INDEX; Schema: prod; Owner: postgres
--

CREATE UNIQUE INDEX attribute_name_uindex ON prod.attribute USING btree (name);


--
-- Name: job_attribute_attribute_id_job_id_uindex; Type: INDEX; Schema: prod; Owner: postgres
--

CREATE UNIQUE INDEX job_attribute_attribute_id_job_id_uindex ON prod.job_attribute USING btree (attribute_id, job_id);


--
-- Name: job_attribute_job_id_attribute_id_uindex; Type: INDEX; Schema: prod; Owner: postgres
--

CREATE UNIQUE INDEX job_attribute_job_id_attribute_id_uindex ON prod.job_attribute USING btree (job_id, attribute_id);


--
-- Name: job_product_version_name_measurement_timestamp_value_index; Type: INDEX; Schema: prod; Owner: postgres
--

CREATE INDEX job_product_version_name_measurement_timestamp_value_index ON prod.job USING btree (product, version, name, measurement, "timestamp", value);


--
-- Name: job_attribute job_attribute_attribute_id_fk; Type: FK CONSTRAINT; Schema: prod; Owner: postgres
--

ALTER TABLE ONLY prod.job_attribute
  ADD CONSTRAINT job_attribute_attribute_id_fk FOREIGN KEY (attribute_id) REFERENCES prod.attribute(id);


--
-- Name: job_attribute job_attribute_job_id_fk; Type: FK CONSTRAINT; Schema: prod; Owner: postgres
--

ALTER TABLE ONLY prod.job_attribute
  ADD CONSTRAINT job_attribute_job_id_fk FOREIGN KEY (job_id) REFERENCES prod.job(id);


--
-- Name: test; Type: SCHEMA; Schema: -; Owner: postgres
--

CREATE SCHEMA test;


ALTER SCHEMA test OWNER TO postgres;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: attribute; Type: TABLE; Schema: test; Owner: postgres
--

CREATE TABLE test.attribute (
                              id integer NOT NULL,
                              name text NOT NULL
);


ALTER TABLE test.attribute OWNER TO postgres;

--
-- Name: attribute_id_seq; Type: SEQUENCE; Schema: test; Owner: postgres
--

CREATE SEQUENCE test.attribute_id_seq
  AS integer
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE test.attribute_id_seq OWNER TO postgres;

--
-- Name: attribute_id_seq; Type: SEQUENCE OWNED BY; Schema: test; Owner: postgres
--

ALTER SEQUENCE test.attribute_id_seq OWNED BY test.attribute.id;


--
-- Name: job; Type: TABLE; Schema: test; Owner: postgres
--

CREATE TABLE test.job (
                        id integer NOT NULL,
                        product text NOT NULL,
                        version text NOT NULL,
                        name text NOT NULL,
                        measurement text NOT NULL,
                        "timestamp" timestamp without time zone NOT NULL,
                        value integer NOT NULL,
                        created_on timestamp without time zone NOT NULL,
                        created_by text NOT NULL
);


ALTER TABLE test.job OWNER TO postgres;

--
-- Name: job_attribute; Type: TABLE; Schema: test; Owner: postgres
--

CREATE TABLE test.job_attribute (
                                  id integer NOT NULL,
                                  job_id integer NOT NULL,
                                  attribute_id integer NOT NULL
);


ALTER TABLE test.job_attribute OWNER TO postgres;

--
-- Name: job_attribute_id_seq; Type: SEQUENCE; Schema: test; Owner: postgres
--

CREATE SEQUENCE test.job_attribute_id_seq
  AS integer
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE test.job_attribute_id_seq OWNER TO postgres;

--
-- Name: job_attribute_id_seq; Type: SEQUENCE OWNED BY; Schema: test; Owner: postgres
--

ALTER SEQUENCE test.job_attribute_id_seq OWNED BY test.job_attribute.id;


--
-- Name: job_id_seq; Type: SEQUENCE; Schema: test; Owner: postgres
--

CREATE SEQUENCE test.job_id_seq
  AS integer
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE test.job_id_seq OWNER TO postgres;

--
-- Name: job_id_seq; Type: SEQUENCE OWNED BY; Schema: test; Owner: postgres
--

ALTER SEQUENCE test.job_id_seq OWNED BY test.job.id;


--
-- Name: attribute id; Type: DEFAULT; Schema: test; Owner: postgres
--

ALTER TABLE ONLY test.attribute ALTER COLUMN id SET DEFAULT nextval('test.attribute_id_seq'::regclass);


--
-- Name: job id; Type: DEFAULT; Schema: test; Owner: postgres
--

ALTER TABLE ONLY test.job ALTER COLUMN id SET DEFAULT nextval('test.job_id_seq'::regclass);


--
-- Name: job_attribute id; Type: DEFAULT; Schema: test; Owner: postgres
--

ALTER TABLE ONLY test.job_attribute ALTER COLUMN id SET DEFAULT nextval('test.job_attribute_id_seq'::regclass);


--
-- Data for Name: attribute; Type: TABLE DATA; Schema: test; Owner: postgres
--

COPY test.attribute (id, name) FROM stdin;
1	a31
2	a32
3	a33
4	a34
5	a35
6	a41
7	a42
8	a43
9	a4Y
10	a51
11	a52
12	a53
13	a5Y
14	a61
15	a62
16	a63
17	a6Y
18	a71
19	a72
20	a73
21	a74
\.


--
-- Data for Name: job; Type: TABLE DATA; Schema: test; Owner: postgres
--

COPY test.job (id, product, version, name, measurement, "timestamp", value, created_on, created_by) FROM stdin;
1	p2	v21			2043-10-20 14:57:42	0	2018-12-08 13:40:23.844943	192.168.1.2
2	p2	v22			2040-03-02 12:23:37	0	2018-12-08 13:40:23.851148	192.168.1.2
3	p2	v22			1976-04-27 00:35:38	0	2018-12-08 13:40:23.853021	192.168.1.2
4	p3	v3			2061-10-19 08:47:38	0	2018-12-08 13:40:23.854698	192.168.1.2
5	p3	v3			1980-06-07 01:57:32	0	2018-12-08 13:40:23.862754	192.168.1.2
6	pY	v3			2022-10-29 00:48:17	0	2018-12-08 13:40:23.868723	192.168.1.2
7	p3	vY			2004-02-10 09:43:04	0	2018-12-08 13:40:23.873054	192.168.1.2
8	p4	v4	n41		2044-10-26 11:29:06	0	2018-12-08 13:40:23.877641	192.168.1.2
9	p4	v4	n42		2021-03-05 22:36:25	0	2018-12-08 13:40:23.883882	192.168.1.2
10	p4	v4	n42		1982-02-22 06:48:50	0	2018-12-08 13:40:23.888439	192.168.1.2
11	pY	v4	n43		2019-02-14 08:31:03	0	2018-12-08 13:40:23.893044	192.168.1.2
12	p4	vY	n44		2064-09-05 08:43:07	0	2018-12-08 13:40:23.897403	192.168.1.2
13	p4	v4	n45		2027-03-10 11:24:33	0	2018-12-08 13:40:23.901954	192.168.1.2
14	p5	v5	n5	m51	2021-06-03 02:34:46	0	2018-12-08 13:40:23.907575	192.168.1.2
15	p5	v5	n5	m52	2037-05-31 11:29:11	0	2018-12-08 13:40:23.915844	192.168.1.2
16	p5	v5	n5	m52	1983-11-05 02:47:13	0	2018-12-08 13:40:23.920224	192.168.1.2
17	pY	v5	n5	m53	2009-04-20 14:46:28	0	2018-12-08 13:40:23.924929	192.168.1.2
18	p5	vY	n5	m54	2056-12-28 06:22:01	0	2018-12-08 13:40:23.929852	192.168.1.2
19	p5	v5	n5	m55	2047-11-14 23:15:18	0	2018-12-08 13:40:23.934222	192.168.1.2
20	p5	v5	nY	m56	1983-04-03 11:11:02	0	2018-12-08 13:40:23.939151	192.168.1.2
21	p6	v6	n6	m6	2017-07-01 06:00:04	27	2018-12-08 13:40:23.943638	192.168.1.2
22	p6	v6	n6	m6	2017-08-19 08:40:23	10	2018-12-08 13:40:23.950213	192.168.1.2
23	p6	v6	n6	m6	2017-10-05 12:57:44	33	2018-12-08 13:40:23.977465	192.168.1.2
24	p6	v6	n6	m6	2017-10-29 16:18:07	48	2018-12-08 13:40:24.078009	192.168.1.2
25	pY	v6	n6	m6	2018-02-08 00:51:55	74	2018-12-08 13:40:24.088311	192.168.1.2
26	p6	vY	n6	m6	2018-02-16 02:16:51	37	2018-12-08 13:40:24.095866	192.168.1.2
27	p6	v6	n6	m6	2018-03-24 04:25:43	92	2018-12-08 13:40:24.101748	192.168.1.2
28	p6	v6	nY	m6	2018-05-25 04:52:34	11	2018-12-08 13:40:24.110942	192.168.1.2
29	p6	v6	n6	mY	2018-07-02 06:34:54	42	2018-12-08 13:40:24.117339	192.168.1.2
30	p6	v6	n6	m6	2018-07-11 07:26:41	12	2018-12-08 13:40:24.12593	192.168.1.2
31	p7	v7	n7	m7	2018-02-08 00:51:55	0	2018-12-08 13:40:24.13256	192.168.1.2
32	p7	v7	n7	m7	2018-02-16 02:16:51	0	2018-12-08 13:40:24.144903	192.168.1.2
\.


--
-- Data for Name: job_attribute; Type: TABLE DATA; Schema: test; Owner: postgres
--

COPY test.job_attribute (id, job_id, attribute_id) FROM stdin;
1	4	1
2	4	2
3	5	2
4	5	3
5	6	3
6	6	4
7	7	4
8	7	5
9	8	6
10	8	7
11	8	8
12	9	6
13	9	7
14	9	8
15	10	6
16	10	7
17	10	8
18	11	6
19	11	7
20	11	8
21	12	6
22	12	7
23	12	8
24	13	6
25	13	9
26	13	8
27	14	10
28	14	11
29	14	12
30	15	10
31	15	11
32	15	12
33	16	10
34	16	11
35	16	12
36	17	10
37	17	11
38	17	12
39	18	10
40	18	11
41	18	12
42	19	10
43	19	13
44	19	12
45	20	10
46	20	11
47	20	12
48	21	14
49	21	15
50	21	16
51	22	14
52	22	15
53	22	16
54	23	14
55	23	15
56	23	16
57	24	14
58	24	15
59	24	16
60	25	14
61	25	15
62	25	16
63	26	14
64	26	15
65	26	16
66	27	14
67	27	17
68	27	16
69	28	14
70	28	15
71	28	16
72	29	14
73	29	15
74	29	16
75	30	14
76	30	15
77	30	16
78	31	18
79	31	19
80	31	20
81	32	18
82	32	19
83	32	21
\.


--
-- Name: attribute_id_seq; Type: SEQUENCE SET; Schema: test; Owner: postgres
--

SELECT pg_catalog.setval('test.attribute_id_seq', 21, true);


--
-- Name: job_attribute_id_seq; Type: SEQUENCE SET; Schema: test; Owner: postgres
--

SELECT pg_catalog.setval('test.job_attribute_id_seq', 83, true);


--
-- Name: job_id_seq; Type: SEQUENCE SET; Schema: test; Owner: postgres
--

SELECT pg_catalog.setval('test.job_id_seq', 32, true);


--
-- Name: attribute attribute_pk; Type: CONSTRAINT; Schema: test; Owner: postgres
--

ALTER TABLE ONLY test.attribute
  ADD CONSTRAINT attribute_pk PRIMARY KEY (id);


--
-- Name: job_attribute job_attribute_pk; Type: CONSTRAINT; Schema: test; Owner: postgres
--

ALTER TABLE ONLY test.job_attribute
  ADD CONSTRAINT job_attribute_pk PRIMARY KEY (id);


--
-- Name: job job_pk; Type: CONSTRAINT; Schema: test; Owner: postgres
--

ALTER TABLE ONLY test.job
  ADD CONSTRAINT job_pk PRIMARY KEY (id);


--
-- Name: attribute_name_uindex; Type: INDEX; Schema: test; Owner: postgres
--

CREATE UNIQUE INDEX attribute_name_uindex ON test.attribute USING btree (name);


--
-- Name: job_attribute_attribute_id_job_id_uindex; Type: INDEX; Schema: test; Owner: postgres
--

CREATE UNIQUE INDEX job_attribute_attribute_id_job_id_uindex ON test.job_attribute USING btree (attribute_id, job_id);


--
-- Name: job_attribute_job_id_attribute_id_uindex; Type: INDEX; Schema: test; Owner: postgres
--

CREATE UNIQUE INDEX job_attribute_job_id_attribute_id_uindex ON test.job_attribute USING btree (job_id, attribute_id);


--
-- Name: job_product_version_name_measurement_timestamp_value_index; Type: INDEX; Schema: test; Owner: postgres
--

CREATE INDEX job_product_version_name_measurement_timestamp_value_index ON test.job USING btree (product, version, name, measurement, "timestamp", value);


--
-- Name: job_attribute job_attribute_attribute_id_fk; Type: FK CONSTRAINT; Schema: test; Owner: postgres
--

ALTER TABLE ONLY test.job_attribute
  ADD CONSTRAINT job_attribute_attribute_id_fk FOREIGN KEY (attribute_id) REFERENCES test.attribute(id);


--
-- Name: job_attribute job_attribute_job_id_fk; Type: FK CONSTRAINT; Schema: test; Owner: postgres
--

ALTER TABLE ONLY test.job_attribute
  ADD CONSTRAINT job_attribute_job_id_fk FOREIGN KEY (job_id) REFERENCES test.job(id);


--
-- PostgreSQL database dump complete
--
