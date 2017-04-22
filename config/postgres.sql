DROP DATABASE IF EXISTS "challeng-stone";


CREATE DATABASE "challenge-stone"
WITH
OWNER = postgres
ENCODING = 'UTF8'
CONNECTION LIMIT = -1;

CREATE TABLE public."User"
(
name character varying(32) NOT NULL,
password character varying(32) NOT NULL,
PRIMARY KEY (name)
)
WITH (
OIDS = FALSE
);

ALTER TABLE public."User"
OWNER to postgres;

CREATE TABLE public."Invoice"
(
id serial NOT NULL,
amount numeric(18, 6) NOT NULL,
document character varying(14) NOT NULL,
month smallint NOT NULL,
year integer NOT NULL,
is_active boolean NOT NULL,
PRIMARY KEY (id)
)
WITH (
OIDS = FALSE
);

ALTER TABLE public."Invoice"
OWNER to postgres;
