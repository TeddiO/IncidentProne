--
-- PostgreSQL database dump
--

-- Dumped from database version 14.5 (Debian 14.5-1.pgdg110+1)
-- Dumped by pg_dump version 14.4

-- Started on 2022-09-25 14:24:50

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
-- TOC entry 6 (class 2615 OID 16385)
-- Name: incidentprone; Type: SCHEMA; Schema: -; Owner: postgres
--

CREATE SCHEMA incidentprone;


ALTER SCHEMA incidentprone OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 211 (class 1259 OID 16393)
-- Name: reportTypes; Type: TABLE; Schema: incidentprone; Owner: postgres
--

CREATE TABLE incidentprone."reportTypes" (
    reason character varying(64) NOT NULL,
    "internalId" integer NOT NULL
);


ALTER TABLE incidentprone."reportTypes" OWNER TO postgres;

--
-- TOC entry 210 (class 1259 OID 16386)
-- Name: reports; Type: TABLE; Schema: incidentprone; Owner: postgres
--

CREATE TABLE incidentprone.reports (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    "reporterName" character varying(255) NOT NULL,
    "issueType" integer NOT NULL,
    "issueSummary" character varying(140) NOT NULL,
    "overallIssue" character varying(4096) NOT NULL,
    resolved boolean NOT NULL,
    created timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    last_updated timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE incidentprone.reports OWNER TO postgres;

--
-- TOC entry 212 (class 1259 OID 16409)
-- Name: sub_reports; Type: TABLE; Schema: incidentprone; Owner: postgres
--

CREATE TABLE incidentprone.sub_reports (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    username character varying(150) NOT NULL,
    message character varying(1024) NOT NULL,
    referenced_issue uuid NOT NULL,
    "time" timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE incidentprone.sub_reports OWNER TO postgres;

--
-- TOC entry 3330 (class 0 OID 16393)
-- Dependencies: 211
-- Data for Name: reportTypes; Type: TABLE DATA; Schema: incidentprone; Owner: postgres
--

COPY incidentprone."reportTypes" (reason, "internalId") FROM stdin;
P1	1
P2	2
P3	3
P4	4
P5	5
Sink is full	6
Someone knocked over milk	7
Is someone in the office?	8
We're out of coffee!	9
\.


--
-- TOC entry 3329 (class 0 OID 16386)
-- Dependencies: 210
-- Data for Name: reports; Type: TABLE DATA; Schema: incidentprone; Owner: postgres
--

COPY incidentprone.reports (id, "reporterName", "issueType", "issueSummary", "overallIssue", resolved, created, last_updated) FROM stdin;
\.


--
-- TOC entry 3331 (class 0 OID 16409)
-- Dependencies: 212
-- Data for Name: sub_reports; Type: TABLE DATA; Schema: incidentprone; Owner: postgres
--

COPY incidentprone.sub_reports (id, username, message, referenced_issue, "time") FROM stdin;
\.


--
-- TOC entry 3184 (class 2606 OID 16399)
-- Name: reportTypes reportTypes_pkey; Type: CONSTRAINT; Schema: incidentprone; Owner: postgres
--

ALTER TABLE ONLY incidentprone."reportTypes"
    ADD CONSTRAINT "reportTypes_pkey" PRIMARY KEY ("internalId");


--
-- TOC entry 3182 (class 2606 OID 16392)
-- Name: reports reports_pkey; Type: CONSTRAINT; Schema: incidentprone; Owner: postgres
--

ALTER TABLE ONLY incidentprone.reports
    ADD CONSTRAINT reports_pkey PRIMARY KEY (id);


--
-- TOC entry 3187 (class 2606 OID 16415)
-- Name: sub_reports sub_reports_pkey; Type: CONSTRAINT; Schema: incidentprone; Owner: postgres
--

ALTER TABLE ONLY incidentprone.sub_reports
    ADD CONSTRAINT sub_reports_pkey PRIMARY KEY (id);


--
-- TOC entry 3180 (class 1259 OID 16405)
-- Name: fki_reportType; Type: INDEX; Schema: incidentprone; Owner: postgres
--

CREATE INDEX "fki_reportType" ON incidentprone.reports USING btree ("issueType");


--
-- TOC entry 3185 (class 1259 OID 16423)
-- Name: fki_report_uuid; Type: INDEX; Schema: incidentprone; Owner: postgres
--

CREATE INDEX fki_report_uuid ON incidentprone.sub_reports USING btree (referenced_issue);


--
-- TOC entry 3188 (class 2606 OID 16400)
-- Name: reports reportType; Type: FK CONSTRAINT; Schema: incidentprone; Owner: postgres
--

ALTER TABLE ONLY incidentprone.reports
    ADD CONSTRAINT "reportType" FOREIGN KEY ("issueType") REFERENCES incidentprone."reportTypes"("internalId") NOT VALID;


--
-- TOC entry 3189 (class 2606 OID 16418)
-- Name: sub_reports report_uuid; Type: FK CONSTRAINT; Schema: incidentprone; Owner: postgres
--

ALTER TABLE ONLY incidentprone.sub_reports
    ADD CONSTRAINT report_uuid FOREIGN KEY (referenced_issue) REFERENCES incidentprone.reports(id) NOT VALID;


-- Completed on 2022-09-25 14:24:51

--
-- PostgreSQL database dump complete
--

