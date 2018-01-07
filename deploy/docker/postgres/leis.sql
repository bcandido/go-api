/* Sequence */
CREATE SEQUENCE lei_id_seq
  INCREMENT 1
  MINVALUE 1
  MAXVALUE 9223372036854775807
  START 1
  CACHE 1;

/* Tables */
CREATE TABLE public.leis
(
  NOME character varying(100) COLLATE pg_catalog."default" NOT NULL,
  ID integer NOT NULL DEFAULT nextval('lei_id_seq'::regclass),
  CONSTRAINT lei_pkey PRIMARY KEY (NOME)
)

WITH (
  OIDS = FALSE
)

TABLESPACE pg_default;

ALTER TABLE public.leis
  OWNER to postgres;