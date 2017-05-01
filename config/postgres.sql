-- DROP DATABASE IF EXISTS "challenge-stone";

-- CREATE USER docker;
-- CREATE IF NOT EXISTS DATABASE docker
-- WITH
-- ENCODING='UTF8'
-- CONNECTION LIMIT = -1;
-- GRANT ALL PRIVILEGES ON DATABASE docker TO docker;

-- CREATE DATABASE "challenge-stone"
-- WITH
-- OWNER = postgres
-- ENCODING = 'UTF8'
-- CONNECTION LIMIT = -1;

CREATE TABLE IF NOT EXISTS public."User"
(
name character varying(32) NOT NULL,
password character varying(32) NOT NULL,
PRIMARY KEY (name)
)
WITH (
OIDS = FALSE
);

-- ALTER TABLE public."User"
-- OWNER to postgres;

CREATE TABLE IF NOT EXISTS public."Invoice"
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

-- ALTER TABLE public."Invoice"
-- OWNER to postgres;
