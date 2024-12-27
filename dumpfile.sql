--
-- PostgreSQL database dump
--

-- Dumped from database version 16.6 (Ubuntu 16.6-0ubuntu0.24.04.1)
-- Dumped by pg_dump version 16.6 (Ubuntu 16.6-0ubuntu0.24.04.1)

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
-- Name: messages_message_id_seq; Type: SEQUENCE; Schema: public; Owner: sergey
--

CREATE SEQUENCE public.messages_message_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.messages_message_id_seq OWNER TO sergey;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: messages; Type: TABLE; Schema: public; Owner: sergey
--

CREATE TABLE public.messages (
    message_id integer DEFAULT nextval('public.messages_message_id_seq'::regclass) NOT NULL,
    subject character varying,
    body character varying,
    send_to character varying,
    status character varying
);


ALTER TABLE public.messages OWNER TO sergey;

--
-- Name: recepients_recepient_id_seq; Type: SEQUENCE; Schema: public; Owner: sergey
--

CREATE SEQUENCE public.recepients_recepient_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.recepients_recepient_id_seq OWNER TO sergey;

--
-- Name: recepient; Type: TABLE; Schema: public; Owner: sergey
--

CREATE TABLE public.recepient (
    recepient_id integer DEFAULT nextval('public.recepients_recepient_id_seq'::regclass) NOT NULL,
    email_adress character varying,
    status character varying,
    message_id integer
);


ALTER TABLE public.recepient OWNER TO sergey;

--
-- Data for Name: messages; Type: TABLE DATA; Schema: public; Owner: sergey
--

COPY public.messages (message_id, subject, body, send_to, status) FROM stdin;
1	Test Subject	This is a test email.	serzh.rybakov.06@gmail.com	sent
2	Test Subject	This is a test email.	serzh.rybakov.06@gmail.com	sent
3	Test Subject	This is a test email.	serzh.rybakov.06@gmail.com	sent
4	Test Subject	This is a test email.	serzh.rybakov.06@gmail.com	sent
5	Test Subject	This is a test email.	serzh.rybakov.06@gmail.com	sent
6	Test Subject	This is a test email.	serzh.rybakov.06@gmail.com	sent
7	Test Subject	This is a test email.	serzh.rybakov.06@gmail.com	sent
\.


--
-- Data for Name: recepient; Type: TABLE DATA; Schema: public; Owner: sergey
--

COPY public.recepient (recepient_id, email_adress, status, message_id) FROM stdin;
1	serzh.rybakov.06@mail.ru	CC serzh.rybakov.06@mail.ru: успешно отправлено	1
2	jopa342@mail.ru	CC jopa342@mail.ru: успешно отправлено	1
3	serzh.rybakov.06@mail.ru	CC serzh.rybakov.06@mail.ru: успешно отправлено	2
4	jopa342@mail.ru	CC jopa342@mail.ru: успешно отправлено	2
5	serzh.rybakov.06@mail.ru	CC serzh.rybakov.06@mail.ru: успешно отправлено	3
6	jopa342@mail.ru	CC jopa342@mail.ru: успешно отправлено	3
7	serzh.rybakov.06@mail.ru	CC serzh.rybakov.06@mail.ru: успешно отправлено	4
8	jopa342@mail.ru	CC jopa342@mail.ru: успешно отправлено	4
9	serzh.rybakov.06@mail.ru	CC serzh.rybakov.06@mail.ru: успешно отправлено	5
10	jopa342@mail.ru	CC jopa342@mail.ru: успешно отправлено	5
11	serzh.rybakov.06@mail.ru	CC serzh.rybakov.06@mail.ru: успешно отправлено	6
12	jopa342@mail.ru	CC jopa342@mail.ru: успешно отправлено	6
13	serzh.rybakov.06@mail.ru	CC serzh.rybakov.06@mail.ru: успешно отправлено	7
14	jopa342@mail.ru	CC jopa342@mail.ru: успешно отправлено	7
\.


--
-- Name: messages_message_id_seq; Type: SEQUENCE SET; Schema: public; Owner: sergey
--

SELECT pg_catalog.setval('public.messages_message_id_seq', 7, true);


--
-- Name: recepients_recepient_id_seq; Type: SEQUENCE SET; Schema: public; Owner: sergey
--

SELECT pg_catalog.setval('public.recepients_recepient_id_seq', 14, true);


--
-- Name: messages messages_pk; Type: CONSTRAINT; Schema: public; Owner: sergey
--

ALTER TABLE ONLY public.messages
    ADD CONSTRAINT messages_pk PRIMARY KEY (message_id);


--
-- Name: recepient recipient_pk; Type: CONSTRAINT; Schema: public; Owner: sergey
--

ALTER TABLE ONLY public.recepient
    ADD CONSTRAINT recipient_pk PRIMARY KEY (recepient_id);


--
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: pg_database_owner
--

GRANT ALL ON SCHEMA public TO sergey;


--
-- Name: DEFAULT PRIVILEGES FOR TABLES; Type: DEFAULT ACL; Schema: public; Owner: postgres
--

ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA public GRANT ALL ON TABLES TO sergey;


--
-- PostgreSQL database dump complete
--

