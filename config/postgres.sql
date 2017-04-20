DROP DATABASE IF EXISTS "challeng-stone";


CREATE DATABASE "challenge-stone"
WITH
OWNER = postgres
ENCODING = 'UTF8'
CONNECTION LIMIT = -1;


CREATE TABLE public."User"
(
name character varying(32)[],
password character varying(32)[],
PRIMARY KEY (name)
)
WITH (
OIDS = FALSE
);

ALTER TABLE public."User"
OWNER to postgres;


CREATE TABLE public."Invoice"
(
id serial,
amount numeric(18, 6)[],
document character varying(14)[],
month smallint,
year integer,
is_active boolean,
PRIMARY KEY (id)
)
WITH (
OIDS = FALSE
);

ALTER TABLE public."Invoice"
OWNER to postgres;
